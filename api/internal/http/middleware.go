package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"runtime/debug"
	"slices"
	"strings"
	"time"
)

// ---------------------------------------------------------------- respuestas

type cuerpoError struct {
	Error detalleError `json:"error"`
}

type detalleError struct {
	Codigo   string            `json:"codigo"`
	Mensaje  string            `json:"mensaje"`
	Campos   map[string]string `json:"campos,omitempty"`
	Detalles any               `json:"detalles,omitempty"`
}

func responderJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// responderError nunca incluye detalles internos (errores de SQL, rutas, etc.).
func responderError(w http.ResponseWriter, status int, codigo, mensaje string, campos map[string]string) {
	responderJSON(w, status, cuerpoError{Error: detalleError{Codigo: codigo, Mensaje: mensaje, Campos: campos}})
}

func responderErrorDetalles(w http.ResponseWriter, status int, codigo, mensaje string, detalles any) {
	responderJSON(w, status, cuerpoError{Error: detalleError{Codigo: codigo, Mensaje: mensaje, Detalles: detalles}})
}

func (s *Server) errorInterno(w http.ResponseWriter, r *http.Request, err error) {
	s.log.Error("error interno", "metodo", r.Method, "ruta", r.URL.Path, "err", err)
	responderError(w, http.StatusInternalServerError, "interno", "Ocurrió un error. Intenta de nuevo en un rato.", nil)
}

// errDecodificar distingue los errores al leer JSON.
type errDecodificar struct {
	status  int
	codigo  string
	mensaje string
}

func (e *errDecodificar) Error() string { return e.mensaje }

// leerJSON exige Content-Type JSON, limita el tamaño, rechaza campos
// desconocidos y datos extra después del objeto.
func leerJSON(w http.ResponseWriter, r *http.Request, dst any, max int64) error {
	ct := r.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil || mt != "application/json" {
		return &errDecodificar{http.StatusUnsupportedMediaType, "tipo_no_soportado", "El cuerpo debe ser JSON (Content-Type: application/json)."}
	}
	r.Body = http.MaxBytesReader(w, r.Body, max)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			return &errDecodificar{http.StatusRequestEntityTooLarge, "cuerpo_muy_grande", "La solicitud es demasiado grande."}
		}
		return &errDecodificar{http.StatusBadRequest, "json_invalido", "El JSON no es válido o trae campos desconocidos."}
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return &errDecodificar{http.StatusBadRequest, "json_invalido", "El cuerpo debe tener un solo objeto JSON."}
	}
	return nil
}

func responderErrorDecodificar(w http.ResponseWriter, err error) {
	var e *errDecodificar
	if errors.As(err, &e) {
		responderError(w, e.status, e.codigo, e.mensaje, nil)
		return
	}
	responderError(w, http.StatusBadRequest, "json_invalido", "Solicitud inválida.", nil)
}

// ---------------------------------------------------------------- middleware

type escritorConEstado struct {
	http.ResponseWriter
	status int
}

func (e *escritorConEstado) WriteHeader(code int) {
	e.status = code
	e.ResponseWriter.WriteHeader(code)
}

func (s *Server) recuperar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				s.log.Error("panic", "valor", v, "stack", string(debug.Stack()))
				responderError(w, http.StatusInternalServerError, "interno", "Ocurrió un error. Intenta de nuevo en un rato.", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// registrar escribe un log por solicitud. No registra cuerpos, cookies ni la
// IP completa (se enmascara el último tramo).
func (s *Server) registrar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		ew := &escritorConEstado{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ew, r)
		if r.URL.Path == "/healthz" && ew.status == http.StatusOK {
			return
		}
		s.log.Info("http",
			"metodo", r.Method,
			"ruta", r.URL.Path,
			"estado", ew.status,
			"ms", time.Since(inicio).Milliseconds(),
			"ip", enmascararIP(s.ipCliente(r)),
		)
	})
}

func (s *Server) cabecerasSeguridad(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
		h.Set("Cache-Control", "no-store")
		if s.cfg.Produccion {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

func limitarBody(next http.Handler, max int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, max)
		}
		next.ServeHTTP(w, r)
	})
}

// verificarOrigen: en métodos que cambian datos exige un header Origin de la
// lista permitida (defensa CSRF además de SameSite=Strict). Rechaza también
// Sec-Fetch-Site: cross-site.
func (s *Server) verificarOrigen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		origen := r.Header.Get("Origin")
		if origen == "" || !slices.Contains(s.cfg.OrigenesPermitidos, origen) || r.Header.Get("Sec-Fetch-Site") == "cross-site" {
			s.log.Warn("origen rechazado", "metodo", r.Method, "ruta", r.URL.Path, "origen", recortar(origen, 80))
			responderError(w, http.StatusForbidden, "origen_no_permitido", "Solicitud no permitida.", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recortar(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// ---------------------------------------------------------------- IP cliente

// ipCliente usa RemoteAddr y, solo si viene de un proxy confiable (el
// contenedor web), el X-Forwarded-For recorrido de derecha a izquierda.
func (s *Server) ipCliente(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	remota, err := netip.ParseAddr(host)
	if err != nil {
		return "desconocida"
	}
	remota = remota.Unmap()
	if !s.esProxyConfiable(remota) {
		return remota.String()
	}
	var partes []string
	for _, v := range r.Header.Values("X-Forwarded-For") {
		partes = append(partes, strings.Split(v, ",")...)
	}
	for i := len(partes) - 1; i >= 0; i-- {
		a, err := netip.ParseAddr(strings.TrimSpace(partes[i]))
		if err != nil {
			break
		}
		a = a.Unmap()
		if !s.esProxyConfiable(a) {
			return a.String()
		}
	}
	return remota.String()
}

func (s *Server) esProxyConfiable(a netip.Addr) bool {
	for _, p := range s.cfg.ProxiesConfiables {
		if p.Contains(a) {
			return true
		}
	}
	return false
}

// claveLimite agrupa IPv6 por /64 para que rotar direcciones no evada el límite.
func claveLimite(ip string) string {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return ip
	}
	if a.Is6() {
		p, err := a.Prefix(64)
		if err == nil {
			return p.String()
		}
	}
	return a.String()
}

func enmascararIP(ip string) string {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return "?"
	}
	if a.Is4() {
		p, _ := a.Prefix(24)
		return p.String()
	}
	p, _ := a.Prefix(48)
	return p.String()
}

type claveCtx int

const ctxSesion claveCtx = 1

func conValor(r *http.Request, k claveCtx, v any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), k, v))
}
