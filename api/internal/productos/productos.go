// Package productos tiene las reglas del catálogo: precio que aplica según la
// cantidad, disponibilidad pública, badges, normalización para búsqueda y
// validación de lo que edita el admin.
package productos

import (
	"strings"
	"time"
	"unicode"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Disponibilidad pública. El stock exacto nunca sale en la API pública.
const (
	Disponible  = "disponible"
	QuedanPocos = "quedan_pocos"
	Agotado     = "agotado"
)

// Tipos de precio aplicado a una línea de pedido.
const (
	PrecioDetalle = "detalle"
	PrecioOferta  = "oferta"
	PrecioMayor   = "mayor"
)

// Badges posibles.
const (
	BadgeOferta        = "oferta"
	BadgeRemate        = "remate"
	BadgeRecienLlegado = "recien_llegado"
	BadgeTemporada     = "temporada"
)

// Temporadas válidas.
var Temporadas = []string{"navidad", "escolar", "fiestas_patrias", "halloween", "invierno", "verano"}

// EsTemporada indica si t es una temporada válida.
func EsTemporada(t string) bool {
	for _, v := range Temporadas {
		if v == t {
			return true
		}
	}
	return false
}

// CalcularDisponibilidad traduce stock interno a la etiqueta pública.
func CalcularDisponibilidad(stock, umbral int) string {
	switch {
	case stock <= 0:
		return Agotado
	case stock <= umbral:
		return QuedanPocos
	default:
		return Disponible
	}
}

// OfertaVigente: hay precio de oferta y no venció. hoy es AAAA-MM-DD (Chile);
// oferta_hasta es el último día en que vale.
func OfertaVigente(p store.Producto, hoy string) bool {
	if p.PrecioOferta == nil || *p.PrecioOferta <= 0 || *p.PrecioOferta >= p.PrecioDetalle {
		return false
	}
	return p.OfertaHasta == "" || hoy <= p.OfertaHasta
}

// PrecioFinal es el precio al detalle que se cobra hoy (oferta si está vigente).
func PrecioFinal(p store.Producto, hoy string) int64 {
	if OfertaVigente(p, hoy) {
		return *p.PrecioOferta
	}
	return p.PrecioDetalle
}

// PrecioUnitario devuelve el precio por unidad que corresponde a la cantidad:
// el menor entre detalle, oferta vigente y (solo si cantidad >= mayor_minimo)
// precio por mayor.
func PrecioUnitario(p store.Producto, cantidad int, hoy string) (int64, string) {
	precio, tipo := p.PrecioDetalle, PrecioDetalle
	if f := PrecioFinal(p, hoy); f < precio {
		precio, tipo = f, PrecioOferta
	}
	if p.PrecioMayor != nil && p.MayorMinimo != nil && cantidad >= *p.MayorMinimo && *p.PrecioMayor > 0 && *p.PrecioMayor < precio {
		precio, tipo = *p.PrecioMayor, PrecioMayor
	}
	return precio, tipo
}

// Badges calcula las etiquetas del producto. hoy es la fecha local (America/Santiago).
func Badges(p store.Producto, hoy time.Time, diasNuevo int) []string {
	out := []string{}
	if OfertaVigente(p, hoy.Format(time.DateOnly)) {
		if p.EtiquetaOferta == BadgeRemate {
			out = append(out, BadgeRemate)
		} else {
			out = append(out, BadgeOferta)
		}
	}
	if EsRecienLlegado(p.FechaLlegada, hoy, diasNuevo) {
		out = append(out, BadgeRecienLlegado)
	}
	if p.Temporada != "" {
		out = append(out, BadgeTemporada)
	}
	return out
}

// EsRecienLlegado: la fecha de llegada está entre hoy-diasNuevo y hoy.
func EsRecienLlegado(fecha string, hoy time.Time, diasNuevo int) bool {
	if fecha == "" {
		return false
	}
	desde := LlegadaDesde(hoy, diasNuevo)
	hasta := hoy.Format(time.DateOnly)
	return fecha >= desde && fecha <= hasta
}

// LlegadaDesde es la fecha mínima (YYYY-MM-DD) para contar como recién llegado.
func LlegadaDesde(hoy time.Time, diasNuevo int) string {
	return hoy.AddDate(0, 0, -diasNuevo).Format(time.DateOnly)
}

var sinTilde = map[rune]rune{
	'á': 'a', 'à': 'a', 'ä': 'a', 'â': 'a',
	'é': 'e', 'è': 'e', 'ë': 'e', 'ê': 'e',
	'í': 'i', 'ì': 'i', 'ï': 'i', 'î': 'i',
	'ó': 'o', 'ò': 'o', 'ö': 'o', 'ô': 'o',
	'ú': 'u', 'ù': 'u', 'ü': 'u', 'û': 'u',
	'ñ': 'n', 'ç': 'c',
}

// NormalizarBusqueda pasa a minúsculas, quita tildes y deja un solo espacio.
// Así "Niño" encuentra "nino" y viceversa.
func NormalizarBusqueda(s string) string {
	var b strings.Builder
	espacio := false
	for _, r := range strings.ToLower(s) {
		if v, ok := sinTilde[r]; ok {
			r = v
		}
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			if !espacio && b.Len() > 0 {
				b.WriteByte(' ')
			}
			espacio = true
			continue
		}
		espacio = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

// PalabrasBusqueda divide la búsqueda normalizada en hasta max palabras.
func PalabrasBusqueda(q string, max int) []string {
	campos := strings.Fields(NormalizarBusqueda(q))
	if len(campos) > max {
		campos = campos[:max]
	}
	return campos
}

// TextoBusqueda arma el campo de búsqueda del producto.
func TextoBusqueda(p store.Producto) string {
	return NormalizarBusqueda(strings.Join([]string{p.Nombre, p.Codigo, p.Descripcion, p.CategoriaNombre}, " "))
}

// Slugify genera un slug ASCII: "Telón de Navidad 3 m" -> "telon-de-navidad-3-m".
func Slugify(s string) string {
	var b strings.Builder
	guion := false
	for _, r := range NormalizarBusqueda(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			guion = false
		default:
			if !guion && b.Len() > 0 {
				b.WriteByte('-')
				guion = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 80 {
		out = strings.Trim(out[:80], "-")
	}
	return out
}
