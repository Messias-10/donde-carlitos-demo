// Package httpapi expone la API REST (net/http de la librería estándar).
// El directorio se llama internal/http; el paquete, httpapi, para no chocar
// con net/http.
package httpapi

import (
	"log/slog"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Config del servidor HTTP. Ningún valor secreto se registra en logs.
type Config struct {
	Produccion         bool
	OrigenesPermitidos []string       // p. ej. http://localhost:3000
	ProxiesConfiables  []netip.Prefix // desde dónde se acepta X-Forwarded-For
	Tienda             pedidos.Tienda
	PedidosDemo        bool // marca los pedidos web como DEMO (mensaje con aviso de prueba)
	DiasRecienLlegado  int
	Zona               *time.Location
	SecretoSesion      []byte // vacío = panel admin deshabilitado
	DuracionSesion     time.Duration
	InactividadSesion  time.Duration
	PedidosPorVentana  int
	VentanaPedidos     time.Duration
	LoginMaxFallos     int
	LoginVentana       time.Duration
	LoginBloqueo       time.Duration
	CostoBcryptDummy   int // costo del hash señuelo (tests usan bcrypt.MinCost)
}

// ConfigPorDefecto con valores seguros para desarrollo.
func ConfigPorDefecto() Config {
	zona, err := time.LoadLocation("America/Santiago")
	if err != nil {
		zona = time.UTC
	}
	return Config{
		OrigenesPermitidos: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		Tienda: pedidos.Tienda{
			Nombre:         "Donde Carlitos Emporio",
			DireccionLocal: "Alameda 2432, Santiago Centro",
			DireccionCorta: "Alameda 2432",
			WhatsApp:       "56921145368",
		},
		PedidosDemo:       true,
		DiasRecienLlegado: 14,
		Zona:              zona,
		DuracionSesion:    12 * time.Hour,
		InactividadSesion: 2 * time.Hour,
		PedidosPorVentana: 10,
		VentanaPedidos:    10 * time.Minute,
		LoginMaxFallos:    5,
		LoginVentana:      15 * time.Minute,
		LoginBloqueo:      15 * time.Minute,
		CostoBcryptDummy:  auth.CostoBcrypt,
	}
}

// Server agrupa dependencias de los handlers.
type Server struct {
	cfg        Config
	st         store.Store
	log        *slog.Logger
	ahora      func() time.Time
	limPedidos *LimitadorVentana
	bloqueo    *auth.BloqueoLogin

	dummyOnce sync.Once
	dummyHash string
}

// New crea el servidor.
func New(cfg Config, st store.Store, log *slog.Logger) *Server {
	if cfg.Zona == nil {
		cfg.Zona = time.UTC
	}
	return &Server{
		cfg:        cfg,
		st:         st,
		log:        log,
		ahora:      time.Now,
		limPedidos: NuevoLimitadorVentana(cfg.PedidosPorVentana, cfg.VentanaPedidos),
		bloqueo:    auth.NuevoBloqueoLogin(cfg.LoginMaxFallos, cfg.LoginVentana, cfg.LoginBloqueo),
	}
}

// SetReloj reemplaza el reloj (tests).
func (s *Server) SetReloj(f func() time.Time) { s.ahora = f }

// Handler arma rutas y middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.healthz)

	// Públicos
	mux.HandleFunc("GET /api/categorias", s.listarCategorias)
	mux.HandleFunc("GET /api/productos", s.listarProductos)
	mux.HandleFunc("GET /api/productos/{slug}", s.obtenerProducto)
	mux.HandleFunc("GET /api/opciones-entrega", s.opcionesEntrega)
	mux.HandleFunc("POST /api/pedidos", s.crearPedido)

	// Admin
	mux.HandleFunc("POST /api/admin/login", s.login)
	mux.HandleFunc("POST /api/admin/logout", s.logout)
	mux.Handle("GET /api/admin/sesion", s.conSesion(s.sesionActual))
	mux.Handle("GET /api/admin/pedidos", s.conSesion(s.adminListarPedidos))
	mux.Handle("GET /api/admin/pedidos/{id}", s.conSesion(s.adminObtenerPedido))
	mux.Handle("PATCH /api/admin/pedidos/{id}/estado", s.conSesion(s.adminCambiarEstado))
	mux.Handle("GET /api/admin/productos", s.conSesion(s.adminListarProductos))
	mux.Handle("POST /api/admin/productos", s.conSesion(s.adminCrearProducto))
	mux.Handle("GET /api/admin/productos/{id}", s.conSesion(s.adminObtenerProducto))
	mux.Handle("PUT /api/admin/productos/{id}", s.conSesion(s.adminActualizarProducto))
	mux.Handle("DELETE /api/admin/productos/{id}", s.conSesion(s.adminEliminarProducto))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Recurso no encontrado.", nil)
	})

	var h http.Handler = mux
	h = s.verificarOrigen(h)
	h = limitarBody(h, 64<<10)
	h = s.cabecerasSeguridad(h)
	h = s.registrar(h)
	h = s.recuperar(h)
	return h
}

// NuevoHTTPServer aplica timeouts y límites de cabeceras.
func NuevoHTTPServer(addr string, h http.Handler, log *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	if err := s.st.Ping(r.Context()); err != nil {
		s.log.Error("healthz: base de datos no responde", "err", err)
		responderError(w, http.StatusServiceUnavailable, "no_disponible", "Servicio no disponible.", nil)
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
}
