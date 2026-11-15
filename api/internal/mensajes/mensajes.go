// Package mensajes arma los mensajes de WhatsApp que el equipo le manda al
// cliente desde el panel: plantillas (d1)–(d8) de
// propuesta/mensajes-whatsapp.md. El texto vive en plantillas.go (un solo
// archivo) y se prueba byte a byte contra los ejemplos de ese documento.
package mensajes

import (
	"strconv"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Config son los datos del negocio que usan las plantillas (variables de entorno).
type Config struct {
	NombreTienda       string // "Donde Carlitos Emporio"
	DireccionRetiro    string // DIRECCION_RETIRO
	HorarioRetiro      string // HORARIO_RETIRO ("" = se omite la línea 🕘)
	URLSitio           string // URL_SITIO, base del link de seguimiento
	DatosTransferencia string // DATOS_TRANSFERENCIA ("" = se omite el bloque y no se ofrece d2)
}

// Linea de producto para los mensajes. PrecioUnitario nil = línea sin precio.
type Linea struct {
	Cantidad       int
	Descripcion    string
	PrecioUnitario *int64
}

// Vista son las variables de §4 de la plantilla, ya mapeadas desde el pedido.
type Vista struct {
	Numero          string
	Canal           string // santiago | regiones
	Entrega         string // retiro | despacho  (la API guarda "envio"; aquí se usa el nombre de la plantilla)
	Nombre          string // nombre completo; las plantillas usan el primero
	Lineas          []Linea
	Descuento       int64
	CostoEnvio      int64
	EnvioPorPagar   bool
	Total           int64
	Comuna          string
	Agencia         string // nombre visible
	NroSeguimiento  string
	Pago            string // transferencia | en_tienda
	Pagado          bool
	Motivo          string
	LinkSeguimiento string // "" = sin línea de link
	EsDemo          bool
}

// Mensaje listo para mandar.
type Mensaje struct {
	Clave       string `json:"clave"`
	Titulo      string `json:"titulo"`
	Texto       string `json:"texto"`
	WhatsAppURL string `json:"whatsapp_url"`
	Sugerido    bool   `json:"sugerido"`
}

// CLP: separa miles con punto a mano (sin locale): 12500 -> "$12.500".
// Negativos: "-$7.500".
func CLP(n int64) string {
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

// Limpiar (regla 8): \r \n \t -> espacio, borra * _ ~ `, junta espacios y recorta.
func Limpiar(s string) string {
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

// primerNombre (regla 9): la primera palabra del nombre limpio.
func primerNombre(s string) string {
	if f := strings.Fields(Limpiar(s)); len(f) > 0 {
		return f[0]
	}
	return ""
}

func hola(nombre string) string {
	if n := primerNombre(nombre); n != "" {
		return "Hola " + n
	}
	return "Hola"
}

// lineasDatos (regla 10): respeta los saltos, quita \r, recorta y borra vacías.
func lineasDatos(s string) []string {
	var out []string
	for _, l := range strings.Split(strings.ReplaceAll(s, "\r", ""), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// armar (regla 2): bloques separados por exactamente una línea en blanco; las
// líneas vacías no se escriben y un bloque vacío desaparece.
func armar(bloques ...[]string) string {
	var partes []string
	for _, b := range bloques {
		var ls []string
		for _, l := range b {
			if l != "" {
				ls = append(ls, l)
			}
		}
		if len(ls) > 0 {
			partes = append(partes, strings.Join(ls, "\n"))
		}
	}
	return strings.Join(partes, "\n\n")
}

// si devuelve s solo si cond.
func si(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

// encodeURIComponent igual que en JavaScript (deja sin codificar A–Z a–z 0–9 - _ . ! ~ * ' ( )).
func encodeURIComponent(s string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '.', c == '!', c == '~', c == '*', c == '\'', c == '(', c == ')':
			b.WriteByte(c)
		default:
			b.WriteByte('%')
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&15])
		}
	}
	return b.String()
}

// URLWhatsApp: https://wa.me/{numero}?text=… ; sin número (pedidos DEMO) el
// usuario elige el contacto y nunca se abre un chat con un teléfono ficticio.
func URLWhatsApp(numero, texto string) string {
	return "https://wa.me/" + numero + "?text=" + encodeURIComponent(texto)
}

// SeguimientoURL: {URL_SITIO}/seguimiento/{token}.
func SeguimientoURL(urlSitio, token string) string {
	return strings.TrimRight(urlSitio, "/") + "/seguimiento/" + token
}

// NuevaVista mapea un pedido del store a las variables de la plantilla.
func NuevaVista(p store.Pedido, c Config) Vista {
	v := Vista{
		Numero: p.Numero, Canal: p.Canal, Nombre: p.Nombre,
		Descuento: p.Descuento, CostoEnvio: p.CostoEnvio, EnvioPorPagar: p.EnvioPorPagar,
		Total: p.Total, Comuna: p.Comuna, NroSeguimiento: p.NSeguimiento, Pago: p.FormaPago,
		Pagado: p.Pago != nil, Motivo: p.MotivoCancelacion, EsDemo: p.EsDemo,
		Entrega: "retiro",
	}
	if p.Entrega == pedidos.EntregaEnvio {
		v.Entrega = "despacho"
	}
	if a, ok := pedidos.AgenciaPorCodigo(p.Agencia); ok {
		v.Agencia = a.Nombre
	}
	if p.TokenSeguimiento != "" && c.URLSitio != "" {
		v.LinkSeguimiento = SeguimientoURL(c.URLSitio, p.TokenSeguimiento)
	}
	for _, it := range p.Items {
		v.Lineas = append(v.Lineas, Linea{Cantidad: it.Cantidad, Descripcion: it.Descripcion, PrecioUnitario: it.PrecioUnitario})
	}
	return v
}

// Para devuelve los mensajes que aplican al pedido en su estado actual, con el
// que corresponde marcado como sugerido. ahora sirve para el recordatorio de
// pago (esperando pago hace más de 24 h).
func Para(p store.Pedido, c Config, ahora time.Time) []Mensaje {
	v := NuevaVista(p, c)
	destino := strings.TrimPrefix(p.Telefono, "+")
	if p.EsDemo {
		destino = "" // regla del demo: wa.me sin número
	}
	atrasado := p.Estado == pedidos.EstadoEsperandoPago && ahora.Sub(p.EstadoDesde) > 24*time.Hour
	_, hayDatos := DatosTransferencia(v, c)
	var out []Mensaje
	agregar := func(clave, titulo, texto string, sugerido bool) {
		out = append(out, Mensaje{Clave: clave, Titulo: titulo, Texto: texto, WhatsAppURL: URLWhatsApp(destino, texto), Sugerido: sugerido})
	}
	e := p.Estado
	if e == pedidos.EstadoCotizado || e == pedidos.EstadoEsperandoPago {
		agregar("cotizacion", "Enviar cotización", Cotizacion(v, c), e == pedidos.EstadoCotizado || (!atrasado && !hayDatos))
		if hayDatos && v.Pago == pedidos.PagoTransferencia {
			d, _ := DatosTransferencia(v, c)
			agregar("datos_transferencia", "Enviar datos de transferencia", d, e == pedidos.EstadoEsperandoPago && !atrasado)
		}
	}
	if e == pedidos.EstadoEsperandoPago {
		agregar("recordatorio_pago", "Recordar pago", RecordatorioPago(v), atrasado)
	}
	if v.Pagado && (e == pedidos.EstadoPagado || e == pedidos.EstadoPreparando) {
		agregar("pago_recibido", "Avisar pago recibido", PagoRecibido(v), true)
	}
	if e == pedidos.EstadoDespachado {
		agregar("despachado", "Avisar despacho", Despachado(v), true)
	}
	if e == pedidos.EstadoListoRetiro {
		agregar("listo_retiro", "Avisar que está listo para retiro", ListoRetiro(v, c), true)
	}
	if e == pedidos.EstadoCancelado {
		agregar("cancelado", "Avisar cancelación", Cancelado(v), true)
	}
	agregar("abrir_chat", "Abrir chat", AbrirChat(v), e == pedidos.EstadoEntregado)
	return out
}
