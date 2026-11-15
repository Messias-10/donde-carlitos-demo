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
	"github.com/Messias-10/donde-carlitos-demo/api/internal/limite"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/mensajes"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Config del servidor HTTP. Ningún valor secreto se registra en logs.
type Config struct {
	Produccion         bool
	OrigenesPermitidos []string       // p. ej. http://localhost:3100
	ProxiesConfiables  []netip.Prefix // desde dónde se cree el último valor de X-Forwarded-For
	Tienda             mensajes.Config
	HorarioConfirmado  bool
	WhatsAppSantiago   string // wa.me, solo dígitos
	WhatsAppRegiones   string // "" = no confirmado (la landing usa el de Santiago)
	DiasRecienLlegado  int
	Zona               *time.Location
	SecretoSesion      []byte // vacío = panel deshabilitado
	DuracionSesion     time.Duration
	InactividadSesion  time.Duration
	LoginGlobalMinuto  int           // intentos de login por minuto (todas las IP)
	BcryptConcurrentes int           // comparaciones bcrypt simultáneas
	LoginEsperaInicial time.Duration // espera tras el 1.er fallo (usuario+IP); se duplica en cada fallo
	LoginEsperaTope    time.Duration // tope de la espera
	LoginVentana       time.Duration // tras esto sin fallos, se olvidan
	ModoDemo           bool          // MODO_DEMO: todo pedido nuevo es_demo (wa.me sin número + línea demo)
	StockEstricto      bool          // STOCK_ESTRICTO: si no alcanza el stock al pagar, rechaza (si no, avisa)
	AdminPrincipal     string        // ADMIN_USUARIO: nadie más puede editarlo ni cambiarle la clave
	SeguimientoPorIP   int           // consultas de seguimiento por minuto por IP
	SeguimientoGlobal  int           // consultas de seguimiento por minuto en total
	CostoBcryptDummy   int           // costo del hash señuelo (tests usan bcrypt.MinCost)
	CostoBcryptUsuario int           // costo para claves de usuarios creados en el panel
	ImagenesDir        string        // carpeta (volumen) de imágenes subidas; "" = subida deshabilitada
}

// ConfigPorDefecto con valores seguros para desarrollo.
func ConfigPorDefecto() Config {
	zona, err := time.LoadLocation("America/Santiago")
	if err != nil {
		zona = time.UTC
	}
	return Config{
		OrigenesPermitidos: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		Tienda: mensajes.Config{
			NombreTienda:    "Donde Carlitos Emporio",
			DireccionRetiro: "Alameda 2432, Santiago Centro",
			// Horario de nov 2026 según propuesta/mensajes-whatsapp.md (TT-7707471893363789063). Por confirmar.
			HorarioRetiro: "L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00",
			URLSitio:      "http://localhost:3000",
		},
		WhatsAppSantiago:   "56921145368",
		DiasRecienLlegado:  14,
		Zona:               zona,
		DuracionSesion:     12 * time.Hour,
		InactividadSesion:  2 * time.Hour,
		LoginGlobalMinuto:  30,
		BcryptConcurrentes: 2,
		LoginEsperaInicial: time.Second,
		LoginEsperaTope:    5 * time.Minute,
		LoginVentana:       15 * time.Minute,
		SeguimientoPorIP:   30,
		SeguimientoGlobal:  600,
		CostoBcryptDummy:   auth.CostoBcrypt,
		CostoBcryptUsuario: auth.CostoBcrypt,
	}
}

// Server agrupa dependencias de los handlers.
type Server struct {
	cfg            Config
	st             store.Store
	log            *slog.Logger
	ahora          func() time.Time
	limLoginGlobal *limite.Ventana
	limSegIP       *limite.Ventana
	limSegGlobal   *limite.Ventana
	semBcrypt      limite.Semaforo
	semImagen      limite.Semaforo
	bloqueo        *auth.BloqueoLogin

	dummyOnce sync.Once
	dummyHash string
}

// New crea el servidor.
func New(cfg Config, st store.Store, log *slog.Logger) *Server {
	if cfg.Zona == nil {
		cfg.Zona = time.UTC
	}
	return &Server{
		cfg:            cfg,
		st:             st,
		log:            log,
		ahora:          func() time.Time { return time.Now().UTC() },
		limLoginGlobal: limite.NuevaVentana(cfg.LoginGlobalMinuto, time.Minute, 1),
		limSegIP:       limite.NuevaVentana(cfg.SeguimientoPorIP, time.Minute, 20000),
		limSegGlobal:   limite.NuevaVentana(cfg.SeguimientoGlobal, time.Minute, 1),
		semBcrypt:      limite.NuevoSemaforo(max(cfg.BcryptConcurrentes, 1)),
		semImagen:      limite.NuevoSemaforo(2),
		bloqueo:        auth.NuevoBloqueoLogin(cfg.LoginEsperaInicial, cfg.LoginEsperaTope, cfg.LoginVentana),
	}
}

// SetReloj reemplaza el reloj (tests).
func (s *Server) SetReloj(f func() time.Time) { s.ahora = func() time.Time { return f().UTC() } }

// Handler arma rutas y middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.healthz)

	// Públicos (landing y seguimiento). No hay pedido público en v2.
	mux.HandleFunc("GET /api/config-publica", s.configPublica)
	mux.HandleFunc("GET /api/categorias", s.listarCategorias)
	mux.HandleFunc("GET /api/productos", s.listarProductos)
	mux.HandleFunc("GET /api/productos/{slug}", s.obtenerProducto)
	mux.HandleFunc("GET /api/vitrina", s.vitrina)
	mux.HandleFunc("GET /api/opciones-entrega", s.opcionesEntrega)
	mux.HandleFunc("GET /api/imagenes/{archivo}", s.servirImagen)
	mux.HandleFunc("GET /api/seguimiento/{token}", s.seguimiento)

	// Sesión
	mux.HandleFunc("POST /api/admin/login", s.login)
	mux.HandleFunc("POST /api/admin/logout", s.logout)
	mux.Handle("GET /api/admin/sesion", s.panel(s.sesionActual))
	mux.Handle("POST /api/admin/sesion/clave", s.panel(s.cambiarClavePropia))

	// Panel: admin y vendedor
	mux.Handle("GET /api/admin/hoy", s.panel(s.hoyPanel))
	mux.Handle("GET /api/admin/pedidos", s.panel(s.listarPedidos))
	mux.Handle("POST /api/admin/pedidos", s.panel(s.crearPedido))
	mux.Handle("GET /api/admin/pedidos/{id}", s.panel(s.obtenerPedido))
	mux.Handle("PUT /api/admin/pedidos/{id}", s.panel(s.editarPedido))
	mux.Handle("PATCH /api/admin/pedidos/{id}/estado", s.panel(s.cambiarEstado))
	mux.Handle("GET /api/admin/pedidos/{id}/mensajes", s.panel(s.mensajesPedido))
	mux.Handle("GET /api/admin/clientes", s.panel(s.cliente))
	mux.Handle("GET /api/admin/productos", s.panel(s.adminListarProductos))
	mux.Handle("GET /api/admin/productos/{id}", s.panel(s.adminObtenerProducto))

	// Solo admin
	mux.Handle("GET /api/admin/pedidos.csv", s.soloAdmin(s.exportarCSV))
	mux.Handle("GET /api/admin/resumen", s.soloAdmin(s.resumen))
	mux.Handle("GET /api/admin/ajustes", s.soloAdmin(s.ajustes))
	mux.Handle("POST /api/admin/pedidos/{id}/token", s.soloAdmin(s.regenerarToken))
	mux.Handle("POST /api/admin/productos", s.soloAdmin(s.adminCrearProducto))
	mux.Handle("PUT /api/admin/productos/{id}", s.soloAdmin(s.adminActualizarProducto))
	mux.Handle("DELETE /api/admin/productos/{id}", s.soloAdmin(s.adminEliminarProducto))
	mux.Handle("POST /api/admin/productos/{id}/imagen", s.soloAdmin(s.adminSubirImagen))
	mux.Handle("GET /api/admin/usuarios", s.soloAdmin(s.listarUsuarios))
	mux.Handle("POST /api/admin/usuarios", s.soloAdmin(s.crearUsuario))
	mux.Handle("PATCH /api/admin/usuarios/{id}", s.soloAdmin(s.actualizarUsuario))
	mux.Handle("POST /api/admin/usuarios/{id}/clave", s.soloAdmin(s.regenerarClave))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Recurso no encontrado.", nil)
	})

	var h http.Handler = mux
	h = s.verificarOrigen(h)
	h = limitarBody(h, func(r *http.Request) int64 {
		if esRutaSubida(r) {
			return maxImagen + 64<<10 // imagen + encabezados multipart
		}
		return 64 << 10
	})
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
