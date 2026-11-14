package pedidos

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Tienda son los datos públicos que van en el mensaje.
type Tienda struct {
	Nombre         string // "Donde Carlitos Emporio"
	DireccionLocal string // "Alameda 2432, Santiago Centro" (opciones de entrega)
	DireccionCorta string // "Alameda 2432" (mensaje de WhatsApp, según la plantilla)
	WhatsApp       string // formato wa.me, solo dígitos: 56921145368
}

// Límites de la plantilla (propuesta/mensaje-pedido-whatsapp.md §3).
const (
	maxLineasMensaje     = 25
	maxComentarioMensaje = 300
)

// FormatoCLP: 1234567 -> "$1.234.567" (punto de miles, sin decimales).
func FormatoCLP(n int64) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}
	if neg {
		return "-$" + b.String()
	}
	return "$" + b.String()
}

// limpiarMensaje aplica la limpieza de la plantilla: \r \n \t -> espacio,
// borra * _ ~ ` (rompen el formato de WhatsApp), junta espacios y recorta.
func limpiarMensaje(s string) string {
	s = strings.Map(func(r rune) rune {
		switch r {
		case '\r', '\n', '\t':
			return ' '
		case '*', '_', '~', '`':
			return -1
		}
		return r
	}, s)
	return strings.Join(strings.Fields(s), " ")
}

func etiquetaTipo(tipo string) string {
	switch tipo {
	case productos.PrecioMayor:
		return "x mayor"
	case productos.PrecioOferta:
		return "oferta"
	}
	return "detalle"
}

// MensajeWhatsApp arma el texto EXACTO de propuesta/mensaje-pedido-whatsapp.md §2
// (LF, una línea en blanco entre bloques, sin salto final). Cualquier cambio de
// la plantilla se refleja aquí y en su test (ejemplos A, B y C).
func MensajeWhatsApp(p store.Pedido, t Tienda) string {
	var out []string
	if p.EsDemo {
		out = append(out, "*[PEDIDO DE PRUEBA · DEMO]*")
	}
	out = append(out, "Hola "+t.Nombre+", les envío mi pedido de la web 🛍️", "")
	out = append(out, "*Pedido "+p.Numero+"*", "")
	out = append(out, "*Productos*")
	for i, l := range p.Items {
		if i == maxLineasMensaje {
			break
		}
		out = append(out, "• "+strconv.Itoa(l.Cantidad)+" x "+limpiarMensaje(l.Nombre))
		out = append(out, "  Cód. "+limpiarMensaje(l.Codigo)+" · "+FormatoCLP(l.PrecioUnitario)+" c/u ("+etiquetaTipo(l.TipoPrecio)+") = "+FormatoCLP(l.Subtotal))
	}
	if resto := len(p.Items) - maxLineasMensaje; resto > 0 {
		s := "s"
		if resto == 1 {
			s = ""
		}
		out = append(out, "• … y "+strconv.Itoa(resto)+" producto"+s+" más (detalle completo en el pedido "+p.Numero+")")
	}
	out = append(out, "")
	total := "*Total: " + FormatoCLP(p.Total) + "*"
	if p.Entrega == EntregaEnvio {
		total += " (no incluye envío)"
	}
	out = append(out, total, "")

	if p.Entrega == EntregaEnvio {
		region := p.Region
		if r, ok := RegionPorCodigo(p.Region); ok {
			region = r.Nombre
		}
		agencia := p.Agencia
		if a, ok := AgenciaPorCodigo(p.Agencia); ok {
			agencia = a.Nombre
		}
		out = append(out, "📦 *Entrega:* Envío por agencia "+limpiarMensaje(agencia))
		out = append(out, "*Región:* "+limpiarMensaje(region))
		out = append(out, "*Comuna:* "+limpiarMensaje(p.Comuna))
	} else {
		out = append(out, "📍 *Entrega:* Retiro en tienda, "+t.DireccionCorta)
		if c := limpiarMensaje(p.Comuna); c != "" {
			out = append(out, "*Comuna:* "+c)
		}
	}
	out = append(out, "*Nombre:* "+limpiarMensaje(p.Nombre))
	if c := limpiarMensaje(p.Comentario); c != "" {
		if r := []rune(c); len(r) > maxComentarioMensaje {
			c = strings.TrimRight(string(r[:maxComentarioMensaje-1]), " ") + "…"
		}
		out = append(out, "💬 *Comentario:* "+c)
	}
	out = append(out, "", "Quedo esperando su confirmación por este chat.")
	return strings.Join(out, "\n")
}

// URLWhatsApp arma el link wa.me con el texto codificado como
// encodeURIComponent (espacios como %20, no como "+").
func URLWhatsApp(numero, textoMensaje string) string {
	return "https://wa.me/" + numero + "?text=" + strings.ReplaceAll(url.QueryEscape(textoMensaje), "+", "%20")
}
