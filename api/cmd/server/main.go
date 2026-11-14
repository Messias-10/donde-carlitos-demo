// Comando server: API de Donde Carlitos Emporio.
//
//	server               levanta la API (configuración por variables de entorno)
//	server -healthcheck  consulta /healthz local y sale con 0 (ok) o 1 (para Docker)
package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata" // zona America/Santiago en imágenes sin tzdata

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	httpapi "github.com/Messias-10/donde-carlitos-demo/api/internal/http"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/seed"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store/sqlite"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "consulta /healthz y termina")
	flag.Parse()

	addr := env("LISTEN_ADDR", ":8080")
	if *healthcheck {
		os.Exit(chequear(addr))
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: nivelLog(env("LOG_LEVEL", "info"))}))
	if err := correr(addr, log); err != nil {
		log.Error("el servidor terminó con error", "err", err)
		os.Exit(1)
	}
}

func correr(addr string, log *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, admin, err := cargarConfig()
	if err != nil {
		return err
	}

	rutaDB := env("DB_PATH", "/data/carlitos.db")
	st, err := sqlite.Abrir(ctx, rutaDB)
	if err != nil {
		return err
	}
	defer st.Close()

	ahora := time.Now()
	if admin.hash != "" && len(cfg.SecretoSesion) >= 32 {
		if _, err := st.SincronizarAdmin(ctx, admin.usuario, admin.hash, ahora); err != nil {
			return fmt.Errorf("sincronizar admin: %w", err)
		}
		log.Info("panel admin habilitado", "usuario", admin.usuario)
	} else {
		cfg.SecretoSesion = nil
		log.Warn("panel admin deshabilitado: falta ADMIN_PASSWORD_HASH_B64 o SESSION_SECRET; corre setup-admin")
	}

	if envBool("SEED_DEMO", true) {
		if _, err := seed.Cargar(ctx, st, seed.Embebido, ahora, envBool("SEED_PEDIDOS_DEMO", true), log); err != nil {
			return err
		}
	}

	srv := httpapi.New(cfg, st, log)
	hs := httpapi.NuevoHTTPServer(addr, srv.Handler(), log)

	go func() {
		t := time.NewTicker(30 * time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if n, err := st.PurgarSesiones(ctx, time.Now()); err != nil {
					log.Warn("purga de sesiones falló", "err", err)
				} else if n > 0 {
					log.Info("sesiones vencidas borradas", "n", n)
				}
			}
		}
	}()

	errSrv := make(chan error, 1)
	go func() {
		log.Info("API escuchando", "addr", addr, "produccion", cfg.Produccion)
		errSrv <- hs.ListenAndServe()
	}()

	select {
	case err := <-errSrv:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Info("apagando…")
	}
	sctx, scancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer scancel()
	return hs.Shutdown(sctx)
}

type datosAdmin struct {
	usuario string
	hash    string
}

// cargarConfig lee variables de entorno. Los secretos nunca se registran.
func cargarConfig() (httpapi.Config, datosAdmin, error) {
	cfg := httpapi.ConfigPorDefecto()
	cfg.Produccion = env("APP_ENV", "development") == "production"
	cfg.PedidosDemo = envBool("PEDIDOS_DEMO", true)

	if v := env("ALLOWED_ORIGINS", ""); v != "" {
		cfg.OrigenesPermitidos = nil
		for _, o := range strings.Split(v, ",") {
			if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
				cfg.OrigenesPermitidos = append(cfg.OrigenesPermitidos, o)
			}
		}
	}
	if v := env("TRUSTED_PROXIES", ""); v != "" {
		for _, c := range strings.Split(v, ",") {
			p, err := netip.ParsePrefix(strings.TrimSpace(c))
			if err != nil {
				return cfg, datosAdmin{}, fmt.Errorf("TRUSTED_PROXIES: %q no es un CIDR válido", c)
			}
			cfg.ProxiesConfiables = append(cfg.ProxiesConfiables, p)
		}
	}
	if v := env("WHATSAPP_NUMERO", ""); v != "" {
		if _, err := strconv.ParseUint(v, 10, 64); err != nil || len(v) < 10 || len(v) > 15 {
			return cfg, datosAdmin{}, errors.New("WHATSAPP_NUMERO: solo dígitos, formato wa.me (56912345678)")
		}
		cfg.Tienda.WhatsApp = v
	}
	if v := env("TIENDA_DIRECCION", ""); v != "" {
		cfg.Tienda.DireccionLocal = v
	}
	var err error
	if cfg.DiasRecienLlegado, err = envInt("DIAS_RECIEN_LLEGADO", cfg.DiasRecienLlegado, 0, 365); err != nil {
		return cfg, datosAdmin{}, err
	}
	if cfg.PedidosPorVentana, err = envInt("PEDIDOS_POR_IP", cfg.PedidosPorVentana, 1, 10000); err != nil {
		return cfg, datosAdmin{}, err
	}
	min, err := envInt("PEDIDOS_VENTANA_MIN", int(cfg.VentanaPedidos.Minutes()), 1, 1440)
	if err != nil {
		return cfg, datosAdmin{}, err
	}
	cfg.VentanaPedidos = time.Duration(min) * time.Minute

	if v := os.Getenv("SESSION_SECRET"); v != "" {
		b, err := hex.DecodeString(v)
		if err != nil || len(b) < 32 {
			return cfg, datosAdmin{}, errors.New("SESSION_SECRET: debe ser hex de al menos 32 bytes (lo genera setup-admin)")
		}
		cfg.SecretoSesion = b
	}
	admin := datosAdmin{usuario: env("ADMIN_USUARIO", "admin")}
	if v := os.Getenv("ADMIN_PASSWORD_HASH_B64"); v != "" {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil || !auth.EsHashBcrypt(string(b)) {
			return cfg, datosAdmin{}, errors.New("ADMIN_PASSWORD_HASH_B64: no es un hash bcrypt en base64 (lo genera setup-admin)")
		}
		admin.hash = string(b)
	}
	return cfg, admin, nil
}

func chequear(addr string) int {
	host := addr
	if strings.HasPrefix(host, ":") {
		host = "127.0.0.1" + host
	}
	c := &http.Client{Timeout: 3 * time.Second}
	res, err := c.Get("http://" + host + "/healthz")
	if err != nil {
		return 1
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

func env(k, def string) string {
	if v, ok := os.LookupEnv(k); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return def
}

func envBool(k string, def bool) bool {
	switch strings.ToLower(env(k, "")) {
	case "1", "true", "si", "sí", "yes":
		return true
	case "0", "false", "no":
		return false
	}
	return def
}

func envInt(k string, def, min, max int) (int, error) {
	v := env(k, "")
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < min || n > max {
		return 0, fmt.Errorf("%s: debe ser un entero entre %d y %d", k, min, max)
	}
	return n, nil
}

func nivelLog(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}
