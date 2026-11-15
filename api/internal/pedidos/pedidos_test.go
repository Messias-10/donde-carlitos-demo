package pedidos

import (
	"errors"
	"testing"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

func i64(v int64) *int64 { return &v }
func iptr(v int) *int    { return &v }

func TestTransiciones(t *testing.T) {
	casos := []struct {
		entrega, pago, de, a string
		ok                   bool
	}{
		{EntregaEnvio, PagoTransferencia, EstadoCotizado, EstadoEsperandoPago, true},
		{EntregaEnvio, PagoTransferencia, EstadoCotizado, EstadoPagado, true},
		{EntregaEnvio, PagoTransferencia, EstadoEsperandoPago, EstadoPagado, true},
		{EntregaEnvio, PagoTransferencia, EstadoPagado, EstadoPreparando, true},
		{EntregaEnvio, PagoTransferencia, EstadoPreparando, EstadoDespachado, true},
		{EntregaEnvio, PagoTransferencia, EstadoDespachado, EstadoEntregado, true},
		{EntregaRetiro, PagoEnTienda, EstadoCotizado, EstadoPreparando, true}, // paga al retirar
		{EntregaRetiro, PagoEnTienda, EstadoPreparando, EstadoListoRetiro, true},
		{EntregaRetiro, PagoEnTienda, EstadoListoRetiro, EstadoEntregado, true},
		{EntregaEnvio, PagoTransferencia, EstadoPagado, EstadoDespachado, true},   // directo (QA B-1a)
		{EntregaRetiro, PagoTransferencia, EstadoPagado, EstadoListoRetiro, true}, // directo
		{EntregaEnvio, PagoTransferencia, EstadoPreparando, EstadoCancelado, true},
		{EntregaRetiro, PagoEnTienda, EstadoListoRetiro, EstadoCancelado, true},
		// inválidas
		{EntregaEnvio, PagoTransferencia, EstadoCotizado, EstadoPreparando, false}, // transferencia: primero pagar
		{EntregaEnvio, PagoTransferencia, EstadoPreparando, EstadoListoRetiro, false},
		{EntregaRetiro, PagoEnTienda, EstadoPreparando, EstadoDespachado, false},
		{EntregaEnvio, PagoTransferencia, EstadoEsperandoPago, EstadoCotizado, false},
		{EntregaEnvio, PagoTransferencia, EstadoCotizado, EstadoEntregado, false},
		{EntregaEnvio, PagoTransferencia, EstadoEntregado, EstadoCancelado, false},
		{EntregaEnvio, PagoTransferencia, EstadoDespachado, EstadoCancelado, false}, // C3: lo despachado no se cancela
		{EntregaEnvio, PagoTransferencia, EstadoPagado, EstadoListoRetiro, false},
		{EntregaEnvio, PagoTransferencia, EstadoCancelado, EstadoCotizado, false},
		{EntregaEnvio, PagoTransferencia, EstadoCotizado, "inventado", false},
	}
	for _, c := range casos {
		if got := TransicionPermitida(c.entrega, c.pago, c.de, c.a); got != c.ok {
			t.Errorf("%s/%s %s->%s = %v, quería %v", c.entrega, c.pago, c.de, c.a, got, c.ok)
		}
	}
}

func TestValidadorRequisitosYStock(t *testing.T) {
	v := Validador()
	envio := store.Pedido{Entrega: EntregaEnvio, FormaPago: PagoTransferencia}
	pago := &store.Pago{Medio: "transferencia", Monto: 1000, Fecha: "2026-11-14"}

	// pagado: exige pago y confirmación; descuenta.
	envio.Estado = EstadoEsperandoPago
	var req *ErrRequisitos
	if _, err := v(envio, store.CambioEstado{Nuevo: EstadoPagado, Pago: pago}); !errors.As(err, &req) || req.Campos["confirmo_pago_recibido"] == "" {
		t.Fatalf("pagado sin confirmación: %v", err)
	}
	if ef, err := v(envio, store.CambioEstado{Nuevo: EstadoPagado, Pago: pago, ConfirmoPago: true}); err != nil || ef != store.StockDescontar {
		t.Fatalf("pagado: %v %v", ef, err)
	}
	// despachado: agencia, n.º y comuna.
	envio.Estado = EstadoPreparando
	if _, err := v(envio, store.CambioEstado{Nuevo: EstadoDespachado}); !errors.As(err, &req) || len(req.Campos) != 3 {
		t.Fatalf("despachado sin datos: %v", err)
	}
	envio.Agencia, envio.Comuna = "starken", "Chillán"
	if _, err := v(envio, store.CambioEstado{Nuevo: EstadoDespachado, NSeguimiento: "123"}); err != nil {
		t.Fatalf("despachado con datos: %v", err)
	}
	// retiro pagando en tienda: preparar descuenta; entregar exige el pago.
	retiro := store.Pedido{Entrega: EntregaRetiro, FormaPago: PagoEnTienda, Estado: EstadoCotizado}
	if ef, err := v(retiro, store.CambioEstado{Nuevo: EstadoPreparando}); err != nil || ef != store.StockDescontar {
		t.Fatalf("preparar sin pago: %v %v", ef, err)
	}
	retiro.Estado, retiro.StockDescontado = EstadoListoRetiro, true
	if _, err := v(retiro, store.CambioEstado{Nuevo: EstadoEntregado}); !errors.As(err, &req) {
		t.Fatalf("entregar sin pago registrado: %v", err)
	}
	if ef, err := v(retiro, store.CambioEstado{Nuevo: EstadoEntregado, Pago: pago}); err != nil || ef != store.StockSinCambio {
		t.Fatalf("entregar con pago: %v %v", ef, err)
	}
	// cancelar devuelve solo si estaba descontado.
	if ef, _ := v(retiro, store.CambioEstado{Nuevo: EstadoCancelado}); ef != store.StockDevolver {
		t.Fatal("cancelar descontado debe devolver")
	}
	retiro.Estado, retiro.StockDescontado = EstadoCotizado, false
	if ef, _ := v(retiro, store.CambioEstado{Nuevo: EstadoCancelado}); ef != store.StockSinCambio {
		t.Fatal("cancelar sin descontar no toca stock")
	}
	// pagado -> preparando no descuenta dos veces.
	envio.Estado, envio.StockDescontado = EstadoPagado, true
	if ef, _ := v(envio, store.CambioEstado{Nuevo: EstadoPreparando}); ef != store.StockSinCambio {
		t.Fatal("no debe descontar dos veces")
	}
	var ti *ErrTransicionInvalida
	if _, err := v(envio, store.CambioEstado{Nuevo: EstadoEntregado}); !errors.As(err, &ti) {
		t.Fatalf("transición inválida: %v", err)
	}
}

func TestCotizador(t *testing.T) {
	prods := map[int64]store.Producto{
		1: {ID: 1, Codigo: "P001", Nombre: "Squishy", PrecioDetalle: 4000, PrecioMayor: i64(3500), MayorMinimo: iptr(12), Activo: true},
		2: {ID: 2, Codigo: "P012", Nombre: "Tren", PrecioDetalle: 20000, PrecioOferta: i64(17000), OfertaHasta: "2026-11-15", Activo: true},
		3: {ID: 3, PrecioDetalle: 1000, Activo: false},
	}
	items := []store.ItemSolicitado{{ProductoID: 1, Cantidad: 11}, {ProductoID: 1, Cantidad: 12}, {ProductoID: 2, Cantidad: 1},
		{Descripcion: "Bolsas surtidas", Cantidad: 2, PrecioUnitario: i64(500)}, {Descripcion: "Regalo", Cantidad: 1}}
	out, sub, err := Cotizador("2026-11-14")(prods, items)
	if err != nil {
		t.Fatal(err)
	}
	if *out[0].PrecioUnitario != 4000 || *out[1].PrecioUnitario != 3500 || out[1].TipoPrecio != "mayor" {
		t.Fatal("el precio por mayor solo desde el mínimo")
	}
	if *out[2].PrecioUnitario != 17000 || out[4].PrecioUnitario != nil || out[4].TipoPrecio != TipoLibre {
		t.Fatalf("oferta o línea sin precio: %+v", out)
	}
	if sub != 11*4000+12*3500+17000+1000 {
		t.Fatalf("subtotal %d", sub)
	}
	out, _, _ = Cotizador("2026-11-16")(prods, []store.ItemSolicitado{{ProductoID: 2, Cantidad: 1}})
	if *out[0].PrecioUnitario != 20000 {
		t.Fatal("oferta vencida")
	}
	var ei *ErrItems
	if _, _, err := Cotizador("2026-11-14")(prods, []store.ItemSolicitado{{ProductoID: 3, Cantidad: 1}, {ProductoID: 9, Cantidad: 1}}); !errors.As(err, &ei) || len(ei.Campos) != 2 {
		t.Fatalf("inactivo/inexistente: %v", err)
	}
}

func TestEntradaValidar(t *testing.T) {
	e := Entrada{Canal: CanalRegiones, Nombre: " Pedro  Soto ", Telefono: "9 0000 0101",
		Items: []ItemEntrada{{ProductoID: 1, Cantidad: 2}, {ProductoID: 1, Cantidad: 3}}}
	d, inicial, errs := e.Validar()
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if d.Entrega != EntregaEnvio || d.FormaPago != PagoTransferencia || inicial != EstadoCotizado ||
		d.Nombre != "Pedro Soto" || d.Telefono != "+56900000101" || len(d.Items) != 1 || d.Items[0].Cantidad != 5 {
		t.Fatalf("defaults regiones: %+v", d)
	}
	s := Entrada{Canal: CanalSantiago, Telefono: "912345678", Items: []ItemEntrada{{ProductoID: 1, Cantidad: 1}}}
	if d, _, errs := s.Validar(); len(errs) > 0 || d.Entrega != EntregaRetiro || d.FormaPago != PagoEnTienda || d.Nombre != "" {
		t.Fatalf("defaults santiago (nombre opcional): %+v %v", d, errs)
	}

	malos := map[string]func(*Entrada){
		"canal":                func(e *Entrada) { e.Canal = "temuco" },
		"telefono":             func(e *Entrada) { e.Telefono = "123" },
		"nombre":               func(e *Entrada) { e.Nombre = "<script>" },
		"forma_pago":           func(e *Entrada) { e.FormaPago = PagoEnTienda }, // envío + en tienda
		"descuento":            func(e *Entrada) { e.Descuento = -1 },
		"items":                func(e *Entrada) { e.Items = nil },
		"total":                func(e *Entrada) { e.Items = []ItemEntrada{{Descripcion: "Algo sin precio", Cantidad: 1}} },
		"items[0]":             func(e *Entrada) { e.Items = []ItemEntrada{{ProductoID: 1, Cantidad: 1, Descripcion: "x"}} },
		"items[0].descripcion": func(e *Entrada) { e.Items = []ItemEntrada{{Cantidad: 1, PrecioUnitario: i64(1)}} },
		"items[0].precio_unitario": func(e *Entrada) {
			e.Items = []ItemEntrada{{Descripcion: "Algo", Cantidad: 1, PrecioUnitario: i64(-5)}}
		},
		"region":         func(e *Entrada) { e.Region = "marte" },
		"agencia":        func(e *Entrada) { e.Agencia = "dhl-falso" },
		"estado_inicial": func(e *Entrada) { e.EstadoInicial = "pagado" },
		"nota_interna":   func(e *Entrada) { e.NotaInterna = "hola‮oculto" },
	}
	for campo, mod := range malos {
		x := e
		mod(&x)
		if _, _, errs := x.Validar(); errs[campo] == "" {
			t.Errorf("esperaba error en %q, obtuve %v", campo, errs)
		}
	}
	// Línea sin precio con total escrito: válido.
	x := e
	x.Items = []ItemEntrada{{Descripcion: "Bolsas de regalo surtidas", Cantidad: 1}}
	x.Total = i64(20000)
	if _, _, errs := x.Validar(); len(errs) > 0 {
		t.Fatal(errs)
	}
}

func TestCambioEntradaValidar(t *testing.T) {
	hoy := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	c, errs := CambioEntrada{Estado: EstadoPagado, Pago: &PagoEntrada{Medio: "transferencia", Monto: 98000}, ConfirmoPagoRecibido: true}.Validar(hoy)
	if len(errs) > 0 || c.Pago.Fecha != "2026-11-14" || !c.ConfirmoPago {
		t.Fatalf("%+v %v", c, errs)
	}
	for _, in := range []CambioEntrada{
		{Estado: "x"},
		{Estado: EstadoPagado, Pago: &PagoEntrada{Medio: "bitcoin", Monto: 1}},
		{Estado: EstadoPagado, Pago: &PagoEntrada{Medio: "efectivo", Monto: 0}},
		{Estado: EstadoPagado, Pago: &PagoEntrada{Medio: "efectivo", Monto: 1, Fecha: "2026-12-31"}},
		{Estado: EstadoPreparando, Pago: &PagoEntrada{Medio: "efectivo", Monto: 1}},
		{Estado: EstadoDespachado, NSeguimiento: "<b>1</b>"},
		{Estado: EstadoDespachado, Agencia: "dhl-falso"},
		{Estado: EstadoPreparando, Motivo: "no corresponde"},
	} {
		if _, errs := in.Validar(hoy); len(errs) == 0 {
			t.Errorf("debió rechazarse: %+v", in)
		}
	}
}

func TestTelefono(t *testing.T) {
	for in, want := range map[string]string{
		"+56 9 2114 5368": "+56921145368", "56921145368": "+56921145368", "921145368": "+56921145368",
		"9 2114-5368": "+56921145368", "(2) 2345 6789": "+56223456789",
	} {
		if got, ok := NormalizarTelefono(in); !ok || got != want {
			t.Errorf("%q -> %q %v", in, got, ok)
		}
	}
	for _, in := range []string{"", "12345678", "+1 555 1234567", "9123456789", "91234567a", "+56 0 1234 5678"} {
		if _, ok := NormalizarTelefono(in); ok {
			t.Errorf("%q debió rechazarse", in)
		}
	}
	if m := EnmascararTelefono("+56912345678"); m != "+569****5678" {
		t.Errorf("máscara = %s", m)
	}
}

// M-5: línea del catálogo con precio editado: tipo "manual", no se junta con otras.
func TestPrecioManualEnLineaDeCatalogo(t *testing.T) {
	e := Entrada{Canal: CanalRegiones, Telefono: "912345678",
		Items: []ItemEntrada{{ProductoID: 1, Cantidad: 2, PrecioUnitario: i64(3000)}, {ProductoID: 1, Cantidad: 1}}}
	d, _, errs := e.Validar()
	if len(errs) > 0 || len(d.Items) != 2 {
		t.Fatalf("%v %+v", errs, d.Items)
	}
	prods := map[int64]store.Producto{1: {ID: 1, Codigo: "P001", Nombre: "Squishy", PrecioDetalle: 4000, Activo: true}}
	out, sub, err := Cotizador("2026-11-14")(prods, d.Items)
	if err != nil || out[0].TipoPrecio != TipoManual || *out[0].PrecioUnitario != 3000 || out[1].TipoPrecio != "detalle" || sub != 10000 {
		t.Fatalf("%+v %d %v", out, sub, err)
	}
	e.Items = []ItemEntrada{{ProductoID: 1, Cantidad: 1, PrecioUnitario: i64(-1)}}
	if _, _, errs := e.Validar(); errs["items[0].precio_unitario"] == "" {
		t.Fatal("precio manual negativo")
	}
}
