package httpapi

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// 16 bytes en base64url sin relleno = 22 caracteres.
var reTokenSeguimiento = regexp.MustCompile(`^[A-Za-z0-9_-]{22}$`)

type hitoDTO struct {
	Estado string    `json:"estado"`
	Fecha  time.Time `json:"fecha"`
}

// seguimientoDTO: SOLO lo que el cliente necesita para saber dónde va su pedido.
// Nunca nombre, teléfono, comuna, región, montos, ítems ni notas.
type seguimientoDTO struct {
	Numero        string    `json:"numero"`
	Estado        string    `json:"estado"`
	Entrega       string    `json:"entrega"`
	AgenciaNombre *string   `json:"agencia_nombre"`
	NSeguimiento  *string   `json:"n_seguimiento"`
	LineaTiempo   []hitoDTO `json:"linea_tiempo"`
}

// seguimiento: GET /api/seguimiento/{token}. 404 genérico y límite por IP y global.
func (s *Server) seguimiento(w http.ResponseWriter, r *http.Request) {
	// La página de seguimiento no debe indexarse (Referrer-Policy: no-referrer ya va en todas).
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	ahora := s.ahora()
	// Solo cuentan las consultas fallidas (404): quien tiene un link válido
	// nunca queda bloqueado por otros. El límite por IP solo aplica si la IP
	// es confiable (hay proxy de borde); si no, todo el tráfico compartiría
	// la IP del contenedor web (seguridad M4 / QA I-5). Con 128 bits de token,
	// el límite global basta para frenar la enumeración.
	noEncontrado := func() {
		limites := []func() (bool, time.Duration){func() (bool, time.Duration) { return s.limSegGlobal.Permitir("global", ahora) }}
		if len(s.cfg.ProxiesConfiables) > 0 {
			limites = append(limites, func() (bool, time.Duration) { return s.limSegIP.Permitir(claveLimite(s.ipCliente(r)), ahora) })
		}
		for _, lim := range limites {
			if ok, espera := lim(); !ok {
				w.Header().Set("Retry-After", strconv.Itoa(segundos(espera)))
				responderError(w, http.StatusTooManyRequests, "demasiadas_solicitudes", "Demasiadas consultas. Intenta en un minuto.", nil)
				return
			}
		}
		responderError(w, http.StatusNotFound, "no_encontrado", "No encontramos ese pedido. Revisa el link que te mandamos.", nil)
	}
	token := r.PathValue("token")
	if !reTokenSeguimiento.MatchString(token) {
		noEncontrado()
		return
	}
	p, err := s.st.PedidoPorToken(r.Context(), token)
	if errors.Is(err, store.ErrNoEncontrado) {
		noEncontrado()
		return
	}
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	// Vence 30 días después de entregado o cancelado: mismo 404 genérico.
	if p.TokenVenceEn != nil && !ahora.Before(*p.TokenVenceEn) {
		noEncontrado()
		return
	}
	out := seguimientoDTO{Numero: p.Numero, Estado: p.Estado, Entrega: p.Entrega, LineaTiempo: []hitoDTO{}}
	if p.Entrega == pedidos.EntregaEnvio {
		if a, ok := pedidos.AgenciaPorCodigo(p.Agencia); ok {
			out.AgenciaNombre = &a.Nombre
		}
		out.NSeguimiento = ptr(p.NSeguimiento)
	}
	for _, h := range p.Historial {
		out.LineaTiempo = append(out.LineaTiempo, hitoDTO{Estado: h.EstadoNuevo, Fecha: h.CreadoEn.UTC()})
	}
	responderJSON(w, http.StatusOK, out)
}
