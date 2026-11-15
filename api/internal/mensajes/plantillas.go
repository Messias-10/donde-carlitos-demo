package mensajes

import "strconv"

// Plantillas (d1)–(d8) de propuesta/mensajes-whatsapp.md §6 — ÚNICO lugar con
// el texto de los mensajes del panel. Si el documento cambia, se cambia aquí y
// en plantillas_test.go (ejemplos byte a byte).

const (
	lineaDemo      = "*[PRUEBA · DEMO]*"
	avisoSeguridad = "Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas."
	maxLineas      = 25
)

func demo(v Vista) string { return si(v.EsDemo, lineaDemo) }

// Cotizacion es (d1) Cotización / resumen del pedido.
func Cotizacion(v Vista, c Config) string {
	ag, co := Limpiar(v.Agencia), Limpiar(v.Comuna)
	var datos []string
	if v.Pago == "transferencia" {
		datos = lineasDatos(c.DatosTransferencia)
	}
	var items []string
	for i, l := range v.Lineas {
		if i == maxLineas {
			break
		}
		base := "• " + strconv.Itoa(l.Cantidad) + " x " + Limpiar(l.Descripcion)
		if l.PrecioUnitario != nil {
			base += " · " + CLP(*l.PrecioUnitario) + " c/u = " + CLP(int64(l.Cantidad)*(*l.PrecioUnitario))
		}
		items = append(items, base)
	}
	if resto := len(v.Lineas) - maxLineas; resto > 0 {
		items = append(items, "• … y "+strconv.Itoa(resto)+" producto"+si(resto != 1, "s")+" más")
	}
	envio := v.CostoEnvio > 0
	items = append(items,
		si(v.Descuento > 0, "Descuento: -"+CLP(v.Descuento)),
		si(envio, "Envío"+si(ag != "", " por "+ag)+": "+CLP(v.CostoEnvio)))

	var entrega string
	if v.Entrega == "despacho" {
		entrega = "📦 Despacho" + si(co != "", " a "+co)
		if ag != "" {
			entrega += " por " + ag + "."
		} else {
			entrega += ". ¿Por qué agencia te lo mandamos?"
		}
	} else {
		entrega = "📍 Retiro en " + c.DireccionRetiro + "."
	}

	var bloqueDatos []string
	if len(datos) > 0 {
		bloqueDatos = append([]string{"*Para pagar, transfiere a:*"}, datos...)
	}
	var pago string
	switch {
	case v.Pago == "en_tienda":
		pago = "Lo pagas al retirar en el local."
	case len(datos) > 0:
		pago = "Después mándanos el comprobante por este chat y te confirmamos."
	default:
		pago = "Si está todo bien, respóndenos *OK* y te mandamos los datos para pagar."
	}
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", te escribimos de *Donde Carlitos Emporio* 🛍️",
			"Este es el resumen de tu pedido *" + v.Numero + "*:"},
		items,
		[]string{"*Total: " + CLP(v.Total) + "*",
			si(v.Entrega == "despacho" && v.EnvioPorPagar && !envio, "El envío lo pagas tú en la agencia, al recibirlo (por pagar).")},
		[]string{entrega},
		bloqueDatos,
		[]string{pago},
		[]string{si(v.LinkSeguimiento != "", "Sigue tu pedido aquí: "+v.LinkSeguimiento)},
		[]string{si(v.Pago == "transferencia", avisoSeguridad)},
	)
}

// DatosTransferencia es (d2). ok = false si DATOS_TRANSFERENCIA está vacía (no se ofrece).
func DatosTransferencia(v Vista, c Config) (string, bool) {
	d := lineasDatos(c.DatosTransferencia)
	if len(d) == 0 {
		return "", false
	}
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", estos son los datos para pagar tu pedido *" + v.Numero + "* por *" + CLP(v.Total) + "*:"},
		d,
		[]string{"Cuando transfieras, mándanos el comprobante por este chat."},
		[]string{avisoSeguridad},
	), true
}

// PagoRecibido es (d3).
func PagoRecibido(v Vista) string {
	siguiente := "Ahora lo preparamos y te avisamos por aquí cuando esté listo para retiro."
	if v.Entrega == "despacho" {
		siguiente = "Ahora lo preparamos y te avisamos por aquí cuando salga, con el número de seguimiento de la agencia."
	}
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", ¡recibimos tu pago! Muchas gracias.", "Pedido *" + v.Numero + "* · " + CLP(v.Total)},
		[]string{siguiente},
		[]string{si(v.LinkSeguimiento != "", "Sigue tu pedido aquí: "+v.LinkSeguimiento)},
	)
}

// Despachado es (d4). Requiere agencia, n.º de seguimiento y comuna.
func Despachado(v Vista) string {
	ag := Limpiar(v.Agencia)
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", ¡tu pedido *" + v.Numero + "* ya va en camino! 📦"},
		[]string{"*Agencia:* " + ag, "*N.º de seguimiento:* " + Limpiar(v.NroSeguimiento), "*Destino:* " + Limpiar(v.Comuna)},
		[]string{"Con ese número lo puedes seguir en " + ag + ".",
			si(v.LinkSeguimiento != "", "Y aquí ves el estado de tu pedido: "+v.LinkSeguimiento)},
		[]string{"Si llega algún problema con tu caja, escríbenos por este mismo chat."},
	)
}

// ListoRetiro es (d5).
func ListoRetiro(v Vista, c Config) string {
	horario := Limpiar(c.HorarioRetiro)
	pago := "*A pagar al retirar:* " + CLP(v.Total)
	if v.Pagado {
		pago = "Ya está pagado."
	}
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", ¡tu pedido *" + v.Numero + "* está listo para retiro! 🛍️"},
		[]string{"📍 " + c.DireccionRetiro, si(horario != "", "🕘 "+horario)},
		[]string{pago, "Muestra tu número de pedido en el local: *" + v.Numero + "*"},
		[]string{si(v.LinkSeguimiento != "", "Sigue tu pedido aquí: "+v.LinkSeguimiento)},
	)
}

// Cancelado es (d6).
func Cancelado(v Vista) string {
	motivo := Limpiar(v.Motivo)
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", tu pedido *" + v.Numero + "* quedó cancelado.", si(motivo != "", "Motivo: "+motivo)},
		[]string{si(v.Pagado, "Te escribimos por este chat para coordinar la devolución de tu pago.")},
		[]string{"Si quieres armar otro pedido, escríbenos por aquí cuando quieras."},
	)
}

// RecordatorioPago es (d7), para "Esperando pago" de más de 24 h.
func RecordatorioPago(v Vista) string {
	return armar(
		[]string{demo(v), hola(v.Nombre) + ", te escribimos de *Donde Carlitos Emporio* por tu pedido *" + v.Numero + "* por *" + CLP(v.Total) + "*."},
		[]string{"Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.",
			"Si cambiaste de idea, avísanos y lo cancelamos, sin problema."},
		[]string{avisoSeguridad},
	)
}

// AbrirChat es (d8).
func AbrirChat(v Vista) string {
	return armar([]string{demo(v), hola(v.Nombre) + ", te escribimos de *Donde Carlitos Emporio* por tu pedido *" + v.Numero + "*."})
}
