package mensajes

import (
	"net/url"
	"strings"
	"testing"
)

// Ejemplos de propuesta/mensajes-whatsapp.md §6, byte a byte.

func i(v int64) *int64 { return &v }

var cfgEj = Config{
	NombreTienda:    "Donde Carlitos Emporio",
	DireccionRetiro: "Alameda 2432, Santiago Centro",
	HorarioRetiro:   "L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00",
	URLSitio:        "https://tienda.example",
}

const datosMarcador = "[DATOS QUE CONFIGURA EL DUEÑO EN EL PANEL]"

func comparar(t *testing.T, nombre, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s distinto.\n--- obtuve ---\n%q\n--- quería ---\n%q", nombre, got, want)
	}
	if strings.HasSuffix(got, "\n") || strings.Contains(got, "\n\n\n") || strings.Contains(got, " \n") {
		t.Errorf("%s: formato (salto final, doble línea en blanco o espacio al final)", nombre)
	}
}

func TestCLP(t *testing.T) {
	for n, want := range map[int64]string{0: "$0", 500: "$500", 4000: "$4.000", 12500: "$12.500", 187500: "$187.500", 1250000: "$1.250.000", -7500: "-$7.500"} {
		if got := CLP(n); got != want {
			t.Errorf("%d -> %s", n, got)
		}
	}
}

var (
	pedro = Vista{Numero: "DC-0101", Canal: "regiones", Entrega: "despacho", Nombre: "Pedro Soto",
		Lineas:        []Linea{{12, "Squishy viral que cambia de color con el sol", i(3500)}, {4, "Dinosaurio gigante con ruedas", i(14000)}},
		EnvioPorPagar: true, Total: 98000, Comuna: "Chillán", Agencia: "Starken", Pago: "transferencia",
		LinkSeguimiento: "https://tienda.example/seguimiento/k7Qm2xVb9TfL4nRpA1c3Dw", EsDemo: true}
	marcela = Vista{Numero: "DC-0102", Canal: "regiones", Entrega: "despacho", Nombre: "Marcela Fuentes",
		Lineas:    []Linea{{15, "Caja sorpresa squishy premium", i(12500)}},
		Descuento: 7500, CostoEnvio: 8000, Total: 188000, Comuna: "Concepción", Agencia: "Varmontt", Pago: "transferencia"}
	camila = Vista{Numero: "DC-0104", Canal: "santiago", Entrega: "retiro", Nombre: "Camila Rojas",
		Lineas: []Linea{{3, "Pulpo viral (mueve ojos y tentáculos)", i(4000)}, {1, "Telón de Navidad 3 m x 2,5 m", i(5000)}, {1, "Bolsas de regalo surtidas", nil}},
		Total:  20000, Pago: "en_tienda"}
)

func TestD1Cotizacion(t *testing.T) {
	comparar(t, "d1-A", Cotizacion(pedro, cfgEj), `*[PRUEBA · DEMO]*
Hola Pedro, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *DC-0101*:

• 12 x Squishy viral que cambia de color con el sol · $3.500 c/u = $42.000
• 4 x Dinosaurio gigante con ruedas · $14.000 c/u = $56.000

*Total: $98.000*
El envío lo pagas tú en la agencia, al recibirlo (por pagar).

📦 Despacho a Chillán por Starken.

Si está todo bien, respóndenos *OK* y te mandamos los datos para pagar.

Sigue tu pedido aquí: https://tienda.example/seguimiento/k7Qm2xVb9TfL4nRpA1c3Dw

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.`)

	conDatos := cfgEj
	conDatos.DatosTransferencia = datosMarcador
	comparar(t, "d1-B", Cotizacion(marcela, conDatos), `Hola Marcela, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *DC-0102*:

• 15 x Caja sorpresa squishy premium · $12.500 c/u = $187.500
Descuento: -$7.500
Envío por Varmontt: $8.000

*Total: $188.000*

📦 Despacho a Concepción por Varmontt.

*Para pagar, transfiere a:*
[DATOS QUE CONFIGURA EL DUEÑO EN EL PANEL]

Después mándanos el comprobante por este chat y te confirmamos.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.`)

	comparar(t, "d1-C", Cotizacion(camila, cfgEj), `Hola Camila, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *DC-0104*:

• 3 x Pulpo viral (mueve ojos y tentáculos) · $4.000 c/u = $12.000
• 1 x Telón de Navidad 3 m x 2,5 m · $5.000 c/u = $5.000
• 1 x Bolsas de regalo surtidas

*Total: $20.000*

📍 Retiro en Alameda 2432, Santiago Centro.

Lo pagas al retirar en el local.`)

	// Sin agencia todavía, y más de 25 líneas.
	sinAg := pedro
	sinAg.Agencia, sinAg.EsDemo, sinAg.LinkSeguimiento = "", false, ""
	sinAg.Lineas = nil
	for k := 0; k < 27; k++ {
		sinAg.Lineas = append(sinAg.Lineas, Linea{1, "Lazo", i(1000)})
	}
	got := Cotizacion(sinAg, cfgEj)
	if !strings.Contains(got, "📦 Despacho a Chillán. ¿Por qué agencia te lo mandamos?") ||
		strings.Count(got, "• 1 x Lazo") != 25 || !strings.Contains(got, "• … y 2 productos más") {
		t.Fatalf("sin agencia / tope de 25 líneas:\n%s", got)
	}
}

func TestD2DatosTransferencia(t *testing.T) {
	if _, ok := DatosTransferencia(pedro, cfgEj); ok {
		t.Fatal("sin DATOS_TRANSFERENCIA no se ofrece d2")
	}
	c := cfgEj
	c.DatosTransferencia = datosMarcador
	got, _ := DatosTransferencia(pedro, c)
	comparar(t, "d2-A", got, `*[PRUEBA · DEMO]*
Hola Pedro, estos son los datos para pagar tu pedido *DC-0101* por *$98.000*:

[DATOS QUE CONFIGURA EL DUEÑO EN EL PANEL]

Cuando transfieras, mándanos el comprobante por este chat.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.`)

	c.DatosTransferencia = "  [Banco: lo escribe el dueño]\r\n\r\n[Cuenta: lo escribe el dueño]\n[RUT y nombre: lo escribe el dueño]  "
	luis := Vista{Numero: "DC-0106", Nombre: "Luis Araya", Total: 20000}
	got, _ = DatosTransferencia(luis, c)
	comparar(t, "d2-B", got, `Hola Luis, estos son los datos para pagar tu pedido *DC-0106* por *$20.000*:

[Banco: lo escribe el dueño]
[Cuenta: lo escribe el dueño]
[RUT y nombre: lo escribe el dueño]

Cuando transfieras, mándanos el comprobante por este chat.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.`)
}

func TestD3PagoRecibido(t *testing.T) {
	m := marcela
	m.LinkSeguimiento = "https://tienda.example/seguimiento/R2pX8vLq0WcN5sTe7YbH1g"
	comparar(t, "d3-A", PagoRecibido(m), `Hola Marcela, ¡recibimos tu pago! Muchas gracias.
Pedido *DC-0102* · $188.000

Ahora lo preparamos y te avisamos por aquí cuando salga, con el número de seguimiento de la agencia.

Sigue tu pedido aquí: https://tienda.example/seguimiento/R2pX8vLq0WcN5sTe7YbH1g`)

	andrea := Vista{Numero: "DC-0107", Entrega: "retiro", Nombre: "Andrea Pérez", Total: 12000, EsDemo: true}
	comparar(t, "d3-B", PagoRecibido(andrea), `*[PRUEBA · DEMO]*
Hola Andrea, ¡recibimos tu pago! Muchas gracias.
Pedido *DC-0107* · $12.000

Ahora lo preparamos y te avisamos por aquí cuando esté listo para retiro.`)
}

func TestD4Despachado(t *testing.T) {
	jorge := Vista{Numero: "DC-0103", Entrega: "despacho", Nombre: "Jorge Muñoz", Agencia: "Starken",
		NroSeguimiento: "000123456789", Comuna: "La Serena", LinkSeguimiento: "https://tienda.example/seguimiento/Zt4uB9kQ2mVx7LcP0aRf3w"}
	comparar(t, "d4-A", Despachado(jorge), `Hola Jorge, ¡tu pedido *DC-0103* ya va en camino! 📦

*Agencia:* Starken
*N.º de seguimiento:* 000123456789
*Destino:* La Serena

Con ese número lo puedes seguir en Starken.
Y aquí ves el estado de tu pedido: https://tienda.example/seguimiento/Zt4uB9kQ2mVx7LcP0aRf3w

Si llega algún problema con tu caja, escríbenos por este mismo chat.`)

	m := marcela
	m.EsDemo, m.NroSeguimiento = true, "VT-998877"
	comparar(t, "d4-B", Despachado(m), `*[PRUEBA · DEMO]*
Hola Marcela, ¡tu pedido *DC-0102* ya va en camino! 📦

*Agencia:* Varmontt
*N.º de seguimiento:* VT-998877
*Destino:* Concepción

Con ese número lo puedes seguir en Varmontt.

Si llega algún problema con tu caja, escríbenos por este mismo chat.`)
}

func TestD5ListoRetiro(t *testing.T) {
	comparar(t, "d5-A", ListoRetiro(camila, cfgEj), `Hola Camila, ¡tu pedido *DC-0104* está listo para retiro! 🛍️

📍 Alameda 2432, Santiago Centro
🕘 L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00

*A pagar al retirar:* $20.000
Muestra tu número de pedido en el local: *DC-0104*`)

	andrea := Vista{Numero: "DC-0107", Entrega: "retiro", Nombre: "Andrea Pérez", Total: 12000, Pagado: true, EsDemo: true,
		LinkSeguimiento: "https://tienda.example/seguimiento/Hn6sD1qW8eYk3PzM5uJc0v"}
	comparar(t, "d5-B", ListoRetiro(andrea, cfgEj), `*[PRUEBA · DEMO]*
Hola Andrea, ¡tu pedido *DC-0107* está listo para retiro! 🛍️

📍 Alameda 2432, Santiago Centro
🕘 L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00

Ya está pagado.
Muestra tu número de pedido en el local: *DC-0107*

Sigue tu pedido aquí: https://tienda.example/seguimiento/Hn6sD1qW8eYk3PzM5uJc0v`)

	sinHorario := cfgEj
	sinHorario.HorarioRetiro = ""
	if strings.Contains(ListoRetiro(camila, sinHorario), "🕘") {
		t.Fatal("sin HORARIO_RETIRO se omite la línea 🕘")
	}
}

func TestD6Cancelado(t *testing.T) {
	rodrigo := Vista{Numero: "DC-0108", Nombre: "Rodrigo Vera", Motivo: "Se nos agotó el producto."}
	comparar(t, "d6-A", Cancelado(rodrigo), `Hola Rodrigo, tu pedido *DC-0108* quedó cancelado.
Motivo: Se nos agotó el producto.

Si quieres armar otro pedido, escríbenos por aquí cuando quieras.`)

	sofia := Vista{Numero: "DC-0120", Nombre: "Sofía", Pagado: true, EsDemo: true}
	comparar(t, "d6-B", Cancelado(sofia), `*[PRUEBA · DEMO]*
Hola Sofía, tu pedido *DC-0120* quedó cancelado.

Te escribimos por este chat para coordinar la devolución de tu pago.

Si quieres armar otro pedido, escríbenos por aquí cuando quieras.`)
}

func TestD7RecordatorioPago(t *testing.T) {
	p := pedro
	p.EsDemo = false
	comparar(t, "d7-A", RecordatorioPago(p), `Hola Pedro, te escribimos de *Donde Carlitos Emporio* por tu pedido *DC-0101* por *$98.000*.

Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.
Si cambiaste de idea, avísanos y lo cancelamos, sin problema.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.`)

	sinNombre := Vista{Numero: "DC-0109", Total: 34000, EsDemo: true}
	comparar(t, "d7-B", RecordatorioPago(sinNombre), `*[PRUEBA · DEMO]*
Hola, te escribimos de *Donde Carlitos Emporio* por tu pedido *DC-0109* por *$34.000*.

Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.
Si cambiaste de idea, avísanos y lo cancelamos, sin problema.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.`)
}

func TestD8AbrirChat(t *testing.T) {
	comparar(t, "d8", AbrirChat(Vista{Numero: "DC-0103", Nombre: "Jorge Muñoz"}),
		"Hola Jorge, te escribimos de *Donde Carlitos Emporio* por tu pedido *DC-0103*.")
}

func TestLimpiezaTextoLibre(t *testing.T) {
	v := Vista{Numero: "DC-0001", Nombre: "  *Ana*\n  Pérez ", Motivo: "sin\tstock _ya_ `x`"}
	got := Cancelado(v)
	if !strings.HasPrefix(got, "Hola Ana, tu pedido") || !strings.Contains(got, "Motivo: sin stock ya x") {
		t.Fatalf("limpieza: %q", got)
	}
}

func TestURLWhatsApp(t *testing.T) {
	// Igual que encodeURIComponent: el ejemplo a1 del documento.
	a1 := "Hola Donde Carlitos Emporio 🛍️\nVengo de la web. Quiero hacer un pedido para *retiro en tienda (Santiago)*.\nLo que busco es: "
	want := "https://wa.me/56921145368?text=Hola%20Donde%20Carlitos%20Emporio%20%F0%9F%9B%8D%EF%B8%8F%0AVengo%20de%20la%20web.%20Quiero%20hacer%20un%20pedido%20para%20*retiro%20en%20tienda%20(Santiago)*.%0ALo%20que%20busco%20es%3A%20"
	if got := URLWhatsApp("56921145368", a1); got != want {
		t.Fatalf("url:\n%s\n%s", got, want)
	}
	// Demo: sin número.
	u := URLWhatsApp("", "hola ñandú")
	if !strings.HasPrefix(u, "https://wa.me/?text=") {
		t.Fatal(u)
	}
	pu, _ := url.Parse(u)
	if pu.Query().Get("text") != "hola ñandú" {
		t.Fatal("decodifica distinto")
	}
}
