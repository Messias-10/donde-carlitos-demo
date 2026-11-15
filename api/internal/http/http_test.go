package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	httpapi "github.com/Messias-10/donde-carlitos-demo/api/internal/http"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/seed"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store/sqlite"
)

const origen = "http://localhost:3000"

// claveTest se genera en cada corrida; no hay claves fijas en el código.
type entorno struct {
	srv   *httptest.Server
	clave string
	db    *sqlite.DB
}

func nuevoEntorno(t *testing.T, ajustar func(*httpapi.Config)) *entorno {
	t.Helper()
	return nuevoEntornoEn(t, time.Date(2026, 11, 14, 15, 0, 0, 0, time.UTC), ajustar)
}

func nuevoEntornoEn(t *testing.T, ahora time.Time, ajustar func(*httpapi.Config)) *entorno {
	t.Helper()
	ctx := context.Background()
	db, err := sqlite.Abrir(ctx, filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if _, err := seed.Cargar(ctx, db, seed.Embebido, ahora, log); err != nil {
		t.Fatal(err)
	}
	zona, _ := time.LoadLocation("America/Santiago")
	if _, err := seed.CargarPedidosDemo(ctx, db, seed.PedidosDemo, ahora, zona, log); err != nil {
		t.Fatal(err)
	}

	clave, _ := auth.GenerarClave(24)
	hash, _ := auth.HashClave(clave, bcrypt.MinCost)
	if _, err := db.SincronizarAdmin(ctx, "admin", hash, ahora); err != nil {
		t.Fatal(err)
	}
	secreto, _ := auth.NuevoToken()

	cfg := httpapi.ConfigPorDefecto()
	cfg.SecretoSesion = []byte(secreto)
	cfg.CostoBcryptDummy = bcrypt.MinCost
	cfg.CostoBcryptUsuario = bcrypt.MinCost
	cfg.Tienda.URLSitio = "http://localhost:3100"
	cfg.ProxiesConfiables = []netip.Prefix{netip.MustParsePrefix("127.0.0.0/8")}
	cfg.ImagenesDir = t.TempDir()
	if ajustar != nil {
		ajustar(&cfg)
	}
	s := httpapi.New(cfg, db, log)
	s.SetReloj(func() time.Time { return ahora })
	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)
	return &entorno{srv: ts, clave: clave, db: db}
}

type resp struct {
	status int
	header http.Header
	body   []byte
}

func (r resp) json(t *testing.T, v any) {
	t.Helper()
	if err := json.Unmarshal(r.body, v); err != nil {
		t.Fatalf("json: %v\n%s", err, r.body)
	}
}

func (e *entorno) pedir(t *testing.T, metodo, ruta string, cuerpo any, hdr map[string]string, cookies ...*http.Cookie) resp {
	t.Helper()
	var rd io.Reader
	switch c := cuerpo.(type) {
	case nil:
	case string:
		rd = strings.NewReader(c)
	default:
		b, _ := json.Marshal(c)
		rd = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(metodo, e.srv.URL+ruta, rd)
	if rd != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if metodo != http.MethodGet {
		req.Header.Set("Origin", origen)
	}
	for k, v := range hdr {
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return resp{res.StatusCode, res.Header, b}
}

func (e *entorno) login(t *testing.T) *http.Cookie {
	t.Helper()
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	if r.status != 200 {
		t.Fatalf("login: %d %s", r.status, r.body)
	}
	for _, c := range (&http.Response{Header: r.header}).Cookies() {
		if c.Name == "dc_sesion" {
			if !c.HttpOnly || c.SameSite != http.SameSiteStrictMode {
				t.Fatalf("cookie sin HttpOnly/SameSite=Strict: %+v", c)
			}
			return c
		}
	}
	t.Fatal("login sin cookie")
	return nil
}

type productoResp struct {
	ID             int64    `json:"id"`
	Slug           string   `json:"slug"`
	Precio         int64    `json:"precio"`
	PrecioFinal    int64    `json:"precio_final"`
	PrecioMayor    *int64   `json:"precio_mayor"`
	MayorMinimo    *int     `json:"mayor_minimo"`
	Disponibilidad string   `json:"disponibilidad"`
	Badges         []string `json:"badges"`
}

type pagina struct {
	Items []json.RawMessage `json:"items"`
	Total int               `json:"total"`
}

func (e *entorno) producto(t *testing.T, slug string) productoResp {
	t.Helper()
	r := e.pedir(t, "GET", "/api/productos/"+slug, nil, nil)
	if r.status != 200 {
		t.Fatalf("producto %s: %d", slug, r.status)
	}
	var p productoResp
	r.json(t, &p)
	return p
}

// loginComo entra con cualquier usuario y devuelve la cookie de sesión.
func (e *entorno) loginComo(t *testing.T, usuario, clave string) *http.Cookie {
	t.Helper()
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": usuario, "password": clave}, nil)
	if r.status != 200 {
		t.Fatalf("login %s: %d %s", usuario, r.status, r.body)
	}
	for _, c := range (&http.Response{Header: r.header}).Cookies() {
		if c.Name == "dc_sesion" {
			return c
		}
	}
	t.Fatal("login sin cookie")
	return nil
}

func TestHealthzYCabeceras(t *testing.T) {
	e := nuevoEntorno(t, nil)
	r := e.pedir(t, "GET", "/healthz", nil, nil)
	if r.status != 200 {
		t.Fatalf("healthz %d", r.status)
	}
	for k, v := range map[string]string{"X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer", "X-Frame-Options": "DENY"} {
		if r.header.Get(k) != v {
			t.Errorf("%s = %q", k, r.header.Get(k))
		}
	}
	if !strings.Contains(r.header.Get("Content-Security-Policy"), "frame-ancestors 'none'") || r.header.Get("Strict-Transport-Security") != "" {
		t.Error("CSP o HSTS")
	}
}

func TestPedidoPublicoEliminado(t *testing.T) {
	e := nuevoEntorno(t, nil)
	if r := e.pedir(t, "POST", "/api/pedidos", `{}`, nil); r.status != 404 {
		t.Fatalf("POST /api/pedidos debe no existir: %d", r.status)
	}
}

func TestCatalogoPublicoNoExponeStock(t *testing.T) {
	e := nuevoEntorno(t, nil)
	r := e.pedir(t, "GET", "/api/productos?page_size=48", nil, nil)
	if r.status != 200 || bytes.Contains(r.body, []byte(`"stock"`)) {
		t.Fatalf("%d / la API pública no expone stock", r.status)
	}
	var pg pagina
	r.json(t, &pg)
	if pg.Total != 37 {
		t.Fatalf("total productos = %d", pg.Total)
	}
	if r := e.pedir(t, "GET", "/api/productos?page_size=500", nil, nil); r.status != 400 {
		t.Fatalf("page_size fuera de rango: %d", r.status)
	}
}

func TestBusquedaInyeccionSQLRespuestaSegura(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	for _, q := range []string{"' OR '1'='1", "'; DROP TABLE pedidos; --", "%", "_", "1 UNION SELECT password_hash FROM usuarios"} {
		for _, ruta := range []string{"/api/productos?q=", "/api/admin/pedidos?q="} {
			r := e.pedir(t, "GET", ruta+url.QueryEscape(q), nil, nil, c)
			var pg pagina
			r.json(t, &pg)
			if r.status != 200 || pg.Total != 0 || bytes.Contains(r.body, []byte("$2")) {
				t.Errorf("%s%q: %d total=%d", ruta, q, r.status, pg.Total)
			}
		}
	}
	r := e.pedir(t, "GET", "/api/admin/pedidos?page_size=100", nil, nil, c)
	var pg pagina
	r.json(t, &pg)
	if pg.Total != 20 {
		t.Fatalf("tabla de pedidos intacta: %d", pg.Total)
	}
}

func TestOrigenObligatorioEnEscrituras(t *testing.T) {
	e := nuevoEntorno(t, nil)
	body := map[string]any{"usuario": "admin", "password": "x"}
	for hdr, v := range map[string]string{"Origin": "", "Sec-Fetch-Site": "cross-site"} {
		if r := e.pedir(t, "POST", "/api/admin/login", body, map[string]string{hdr: v}); r.status != 403 {
			t.Fatalf("%s=%q: %d", hdr, v, r.status)
		}
	}
	if r := e.pedir(t, "POST", "/api/admin/login", body, map[string]string{"Origin": "https://evil.example"}); r.status != 403 {
		t.Fatalf("Origin ajeno: %d", r.status)
	}
}

func TestRutasPanelSinSesion401(t *testing.T) {
	e := nuevoEntorno(t, nil)
	rutas := []struct{ m, r string }{
		{"GET", "/api/admin/sesion"}, {"GET", "/api/admin/hoy"}, {"GET", "/api/admin/pedidos"}, {"POST", "/api/admin/pedidos"},
		{"GET", "/api/admin/pedidos.csv"}, {"GET", "/api/admin/pedidos/1"}, {"PUT", "/api/admin/pedidos/1"},
		{"PATCH", "/api/admin/pedidos/1/estado"}, {"GET", "/api/admin/pedidos/1/mensajes"},
		{"GET", "/api/admin/clientes?telefono=900000101"}, {"GET", "/api/admin/resumen"},
		{"GET", "/api/admin/usuarios"}, {"POST", "/api/admin/usuarios"}, {"PATCH", "/api/admin/usuarios/1"},
		{"POST", "/api/admin/usuarios/1/clave"}, {"GET", "/api/admin/productos"}, {"POST", "/api/admin/productos"},
		{"PUT", "/api/admin/productos/1"}, {"DELETE", "/api/admin/productos/1"},
	}
	falsa := &http.Cookie{Name: "dc_sesion", Value: strings.Repeat("A", 43)}
	for _, rt := range rutas {
		for _, cookies := range [][]*http.Cookie{nil, {falsa}} {
			if r := e.pedir(t, rt.m, rt.r, `{}`, nil, cookies...); r.status != 401 {
				t.Errorf("%s %s (cookies=%d) = %d", rt.m, rt.r, len(cookies), r.status)
			}
		}
	}
}

func TestLimiteIntentosLogin(t *testing.T) {
	e := nuevoEntorno(t, nil)
	malo := map[string]string{"usuario": "admin", "password": "incorrecta"}
	for i := 0; i < 5; i++ {
		if r := e.pedir(t, "POST", "/api/admin/login", malo, nil); r.status != 401 {
			t.Fatalf("intento %d: %d", i+1, r.status)
		}
	}
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	if r.status != 429 || r.header.Get("Retry-After") == "" {
		t.Fatalf("tras 5 fallos: %d", r.status)
	}
	r = e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "nadie", "password": "x"}, map[string]string{"X-Forwarded-For": "203.0.113.10"})
	if r.status != 401 || !bytes.Contains(r.body, []byte("credenciales_invalidas")) {
		t.Fatalf("usuario inexistente: %d %s", r.status, r.body)
	}
}

func TestCookieSecureEnProduccion(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.Produccion = true })
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	sc := r.header.Get("Set-Cookie")
	if !strings.HasPrefix(sc, "__Host-dc_sesion=") || !strings.Contains(sc, "Secure") || !strings.Contains(sc, "HttpOnly") ||
		!strings.Contains(sc, "SameSite=Strict") || r.header.Get("Strict-Transport-Security") == "" {
		t.Fatalf("cookie en producción: %s", sc)
	}
}

func itoa(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
