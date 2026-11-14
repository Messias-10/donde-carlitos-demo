package pedidos

import (
	"errors"
	"net/url"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

func i64(v int64) *int64 { return &v }
func iptr(v int) *int    { return &v }

func TestTransiciones(t *testing.T) {
	casos := []struct {
		entrega, de, a string
		ok             bool
	}{
		{EntregaRetiro, EstadoNuevo, EstadoConfirmado, true},
		{EntregaRetiro, EstadoNuevo, EstadoCancelado, true},
		{EntregaRetiro, EstadoConfirmado, EstadoListo, true},
		{EntregaEnvio, EstadoConfirmado, EstadoDespachado, true},
		{EntregaRetiro, EstadoListo, EstadoEntregado, true},
		{EntregaEnvio, EstadoDespachado, EstadoEntregado, true},
		{EntregaRetiro, EstadoListo, EstadoCancelado, true},
		// inválidas
		{EntregaRetiro, EstadoNuevo, EstadoEntregado, false},
		{EntregaRetiro, EstadoNuevo, EstadoListo, false},
		{EntregaRetiro, EstadoConfirmado, EstadoDespachado, false}, // retiro no se despacha
		{EntregaEnvio, EstadoConfirmado, EstadoListo, false},       // envío no queda "listo para retiro"
		{EntregaEnvio, EstadoDespachado, EstadoCancelado, false},
		{EntregaRetiro, EstadoEntregado, EstadoNuevo, false},
		{EntregaRetiro, EstadoCancelado, EstadoConfirmado, false},
		{EntregaRetiro, EstadoConfirmado, EstadoNuevo, false},
		{EntregaRetiro, EstadoNuevo, "inventado", false},
	}
	for _, c := range casos {
		if got := TransicionPermitida(c.entrega, c.de, c.a); got != c.ok {
			t.Errorf("%s %s->%s = %v, quería %v", c.entrega, c.de, c.a, got, c.ok)
		}
	}
	efectos := []struct {
		de, a string
		want  store.EfectoStock
	}{
		{EstadoNuevo, EstadoConfirmado, store.StockDescontar},
		{EstadoNuevo, EstadoCancelado, store.StockSinCambio}, // nunca se descontó
		{EstadoConfirmado, EstadoCancelado, store.StockDevolver},
		{EstadoListo, EstadoCancelado, store.StockDevolver},
		{EstadoConfirmado, EstadoListo, store.StockSinCambio},
		{EstadoListo, EstadoEntregado, store.StockSinCambio},
	}
	for _, c := range efectos {
		got, err := Validador(c.a)(store.Pedido{Entrega: EntregaRetiro, Estado: c.de})
		if err != nil || got != c.want {
			t.Errorf("%s->%s: efecto %v err %v, quería %v", c.de, c.a, got, err, c.want)
		}
	}
	_, err := Validador(EstadoEntregado)(store.Pedido{Entrega: EntregaRetiro, Estado: EstadoNuevo})
	var ti *ErrTransicionInvalida
	if !errors.As(err, &ti) {
		t.Fatalf("quería ErrTransicionInvalida, obtuve %v", err)
	}
}

func TestCotizarIgnoraPrecioDelCliente(t *testing.T) {
	prods := map[int64]store.Producto{
		1: {ID: 1, Codigo: "P001", Nombre: "Squishy", PrecioDetalle: 4000, PrecioMayor: i64(3500), MayorMinimo: iptr(12), Stock: 100, Activo: true},
		2: {ID: 2, Codigo: "P002", Nombre: "Farol", PrecioDetalle: 1500, PrecioOferta: i64(1000), Stock: 10, Activo: true},
	}
	items, total, err := cotizar(prods, []store.ItemSolicitado{{ProductoID: 1, Cantidad: 12}, {ProductoID: 2, Cantidad: 2}}, "2026-11-14")
	if err != nil {
		t.Fatal(err)
	}
	if total != 12*3500+2*1000 {
		t.Fatalf("total = %d", total)
	}
	if items[0].TipoPrecio != "mayor" || items[1].TipoPrecio != "oferta" {
		t.Fatalf("tipos = %s %s", items[0].TipoPrecio, items[1].TipoPrecio)
	}
	// 11 unidades no alcanzan el mínimo por mayor.
	items, _, _ = cotizar(prods, []store.ItemSolicitado{{ProductoID: 1, Cantidad: 11}}, "2026-11-14")
	if items[0].PrecioUnitario != 4000 {
		t.Fatalf("bajo el mínimo debe ser detalle, fue %d", items[0].PrecioUnitario)
	}
}

func TestCotizarRechazaNoDisponibles(t *testing.T) {
	prods := map[int64]store.Producto{
		1: {ID: 1, PrecioDetalle: 1000, Stock: 0, Activo: true},
		2: {ID: 2, PrecioDetalle: 1000, Stock: 5, Activo: false},
		3: {ID: 3, PrecioDetalle: 1000, Stock: 2, Activo: true},
	}
	_, _, err := cotizar(prods, []store.ItemSolicitado{{ProductoID: 1, Cantidad: 1}, {ProductoID: 2, Cantidad: 1}, {ProductoID: 3, Cantidad: 3}, {ProductoID: 99, Cantidad: 1}}, "2026-11-14")
	var nd *ErrNoDisponibles
	if !errors.As(err, &nd) || len(nd.Detalles) != 4 {
		t.Fatalf("quería 4 no disponibles, obtuve %v", err)
	}
	motivos := []string{nd.Detalles[0].Motivo, nd.Detalles[1].Motivo, nd.Detalles[2].Motivo, nd.Detalles[3].Motivo}
	want := []string{MotivoAgotado, MotivoNoDisponible, MotivoStockInsufic, MotivoNoDisponible}
	for i := range want {
		if motivos[i] != want[i] {
			t.Errorf("motivo %d = %s, quería %s", i, motivos[i], want[i])
		}
	}
}

func TestValidarSolicitud(t *testing.T) {
	base := Solicitud{Nombre: "  José   Muñoz ", Telefono: "+56 9 1234 5678", Entrega: "retiro",
		Items: []ItemSolicitud{{ProductoID: 1, Cantidad: 2}, {ProductoID: 1, Cantidad: 3}}}
	n, errs := base.Validar()
	if len(errs) > 0 {
		t.Fatalf("errores inesperados: %v", errs)
	}
	if n.Nombre != "José Muñoz" || n.Telefono != "+56912345678" || len(n.Items) != 1 || n.Items[0].Cantidad != 5 {
		t.Fatalf("normalización: %+v", n)
	}

	malos := []struct {
		campo string
		mod   func(*Solicitud)
	}{
		{"nombre", func(s *Solicitud) { s.Nombre = "<script>" }},
		{"telefono", func(s *Solicitud) { s.Telefono = "12345" }},
		{"entrega", func(s *Solicitud) { s.Entrega = "dron" }},
		{"region", func(s *Solicitud) { s.Entrega, s.Comuna, s.Agencia = "envio", "Temuco", "starken" }},
		{"agencia", func(s *Solicitud) {
			s.Entrega, s.Comuna, s.Region, s.Agencia = "envio", "Temuco", "araucania", "dhl-falso"
		}},
		{"items", func(s *Solicitud) { s.Items = nil }},
		{"items[0].cantidad", func(s *Solicitud) { s.Items = []ItemSolicitud{{ProductoID: 1, Cantidad: 0}} }},
		{"items[0].cantidad", func(s *Solicitud) { s.Items = []ItemSolicitud{{ProductoID: 1, Cantidad: 101}} }},
		{"items[1].cantidad", func(s *Solicitud) {
			s.Items = []ItemSolicitud{{ProductoID: 1, Cantidad: 60}, {ProductoID: 1, Cantidad: 50}}
		}},
		{"comentario", func(s *Solicitud) { s.Comentario = "hola‮oculto" }},
	}
	for _, m := range malos {
		s := base
		m.mod(&s)
		if _, errs := s.Validar(); errs[m.campo] == "" {
			t.Errorf("esperaba error en %q, obtuve %v", m.campo, errs)
		}
	}
}

func TestTelefono(t *testing.T) {
	for in, want := range map[string]string{
		"+56 9 2114 5368": "+56921145368",
		"56921145368":     "+56921145368",
		"921145368":       "+56921145368",
		"9 2114-5368":     "+56921145368",
		"(2) 2345 6789":   "+56223456789",
	} {
		if got, ok := NormalizarTelefono(in); !ok || got != want {
			t.Errorf("%q -> %q %v", in, got, ok)
		}
	}
	for _, in := range []string{"", "12345678", "+1 555 1234567", "056912345678", "9123456789", "91234567a", "+56 0 1234 5678"} {
		if _, ok := NormalizarTelefono(in); ok {
			t.Errorf("%q debió rechazarse", in)
		}
	}
	if m := EnmascararTelefono("+56912345678"); m != "+569****5678" {
		t.Errorf("máscara = %s", m)
	}
}

func TestFormatoCLP(t *testing.T) {
	for n, want := range map[int64]string{0: "$0", 500: "$500", 1000: "$1.000", 25990: "$25.990", 1234567: "$1.234.567"} {
		if got := FormatoCLP(n); got != want {
			t.Errorf("%d -> %s", n, got)
		}
	}
}

var tiendaTest = Tienda{Nombre: "Donde Carlitos Emporio", DireccionLocal: "Alameda 2432, Santiago Centro", DireccionCorta: "Alameda 2432", WhatsApp: "56921145368"}

// Ejemplos A, B y C de propuesta/mensaje-pedido-whatsapp.md §6, byte a byte.
func TestMensajeWhatsAppPlantilla(t *testing.T) {
	casos := []struct {
		nombre string
		p      store.Pedido
		want   string
	}{
		{"A retiro con comentario", store.Pedido{
			Numero: "DC-0007", Nombre: "Camila Rojas", Entrega: EntregaRetiro, Total: 34600,
			Comentario: "Paso el sábado en la mañana.",
			Items: []store.ItemPedido{
				{Codigo: "JUG-0001", Nombre: "Squishy viral que cambia de color con el sol", Cantidad: 2, PrecioUnitario: 4000, TipoPrecio: "detalle", Subtotal: 8000},
				{Codigo: "NAV-0012", Nombre: "Tren navideño para árbol (grande)", Cantidad: 1, PrecioUnitario: 17000, TipoPrecio: "oferta", Subtotal: 17000},
				{Codigo: "NAV-0009", Nombre: "Camino de mesa navideño", Cantidad: 12, PrecioUnitario: 800, TipoPrecio: "mayor", Subtotal: 9600},
			}},
			"Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️\n\n*Pedido DC-0007*\n\n*Productos*\n" +
				"• 2 x Squishy viral que cambia de color con el sol\n  Cód. JUG-0001 · $4.000 c/u (detalle) = $8.000\n" +
				"• 1 x Tren navideño para árbol (grande)\n  Cód. NAV-0012 · $17.000 c/u (oferta) = $17.000\n" +
				"• 12 x Camino de mesa navideño\n  Cód. NAV-0009 · $800 c/u (x mayor) = $9.600\n\n" +
				"*Total: $34.600*\n\n📍 *Entrega:* Retiro en tienda, Alameda 2432\n*Nombre:* Camila Rojas\n" +
				"💬 *Comentario:* Paso el sábado en la mañana.\n\nQuedo esperando su confirmación por este chat."},
		{"B envío demo", store.Pedido{
			Numero: "DC-0008", Nombre: "Pedro Soto", Entrega: EntregaEnvio, Region: "nuble", Comuna: "Chillán", Agencia: "starken",
			Total: 98000, EsDemo: true, Comentario: "Es para mi negocio. ¿Lo pueden despachar el lunes?",
			Items: []store.ItemPedido{
				{Codigo: "JUG-0001", Nombre: "Squishy viral que cambia de color con el sol", Cantidad: 12, PrecioUnitario: 3500, TipoPrecio: "mayor", Subtotal: 42000},
				{Codigo: "JUG-0010", Nombre: "Dinosaurio gigante con ruedas", Cantidad: 4, PrecioUnitario: 14000, TipoPrecio: "mayor", Subtotal: 56000},
			}},
			"*[PEDIDO DE PRUEBA · DEMO]*\nHola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️\n\n*Pedido DC-0008*\n\n*Productos*\n" +
				"• 12 x Squishy viral que cambia de color con el sol\n  Cód. JUG-0001 · $3.500 c/u (x mayor) = $42.000\n" +
				"• 4 x Dinosaurio gigante con ruedas\n  Cód. JUG-0010 · $14.000 c/u (x mayor) = $56.000\n\n" +
				"*Total: $98.000* (no incluye envío)\n\n📦 *Entrega:* Envío por agencia Starken\n*Región:* Ñuble\n*Comuna:* Chillán\n" +
				"*Nombre:* Pedro Soto\n💬 *Comentario:* Es para mi negocio. ¿Lo pueden despachar el lunes?\n\n" +
				"Quedo esperando su confirmación por este chat."},
		{"C limpieza y comuna", store.Pedido{
			Numero: "DC-0009", Nombre: "  Ana  *Pérez*\n", Entrega: EntregaRetiro, Comuna: "Maipú", Total: 500,
			Items: []store.ItemPedido{{Codigo: "NAV-0002", Nombre: "Esfera navideña con luces", Cantidad: 1, PrecioUnitario: 500, TipoPrecio: "detalle", Subtotal: 500}}},
			"Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️\n\n*Pedido DC-0009*\n\n*Productos*\n" +
				"• 1 x Esfera navideña con luces\n  Cód. NAV-0002 · $500 c/u (detalle) = $500\n\n*Total: $500*\n\n" +
				"📍 *Entrega:* Retiro en tienda, Alameda 2432\n*Comuna:* Maipú\n*Nombre:* Ana Pérez\n\n" +
				"Quedo esperando su confirmación por este chat."},
	}
	for _, c := range casos {
		got := MensajeWhatsApp(c.p, tiendaTest)
		if got != c.want {
			t.Errorf("%s: distinto.\n--- obtuve ---\n%q\n--- quería ---\n%q", c.nombre, got, c.want)
		}
		if !utf8.ValidString(got) || strings.HasSuffix(got, "\n") || strings.Contains(got, "\n\n\n") || strings.Contains(got, " \n") {
			t.Errorf("%s: formato (UTF-8, salto final, doble línea en blanco o espacio al final)", c.nombre)
		}
		// El link decodifica exactamente al mismo texto (tildes, ñ, emojis, saltos).
		u := URLWhatsApp(tiendaTest.WhatsApp, got)
		pu, err := url.Parse(u)
		if err != nil || !strings.HasPrefix(u, "https://wa.me/56921145368?text=") || strings.Contains(u, "+") || pu.Query().Get("text") != got {
			t.Errorf("%s: url mal codificada", c.nombre)
		}
	}
}

func TestMensajeWhatsAppLimites(t *testing.T) {
	p := store.Pedido{Numero: "DC-0100", Nombre: "Luz", Entrega: EntregaRetiro, Total: 1,
		Comentario: strings.Repeat("ñ", 350) + "\n😄 *fin*"}
	for i := 0; i < 27; i++ {
		p.Items = append(p.Items, store.ItemPedido{Codigo: "P001", Nombre: "Lazo", Cantidad: 1, PrecioUnitario: 1000, TipoPrecio: "detalle", Subtotal: 1000})
	}
	got := MensajeWhatsApp(p, tiendaTest)
	if strings.Count(got, "• 1 x Lazo") != 25 || !strings.Contains(got, "• … y 2 productos más (detalle completo en el pedido DC-0100)") {
		t.Fatal("tope de 25 líneas")
	}
	linea := got[strings.Index(got, "💬 *Comentario:* "):]
	linea = strings.TrimPrefix(linea[:strings.Index(linea, "\n")], "💬 *Comentario:* ")
	if utf8.RuneCountInString(linea) != 300 || !strings.HasSuffix(linea, "…") {
		t.Fatalf("comentario recortado a %d runas", utf8.RuneCountInString(linea))
	}
	p.Items = p.Items[:26]
	if !strings.Contains(MensajeWhatsApp(p, tiendaTest), "• … y 1 producto más (") {
		t.Fatal("singular")
	}
}
