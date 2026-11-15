package pedidos

import (
	"fmt"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Estados del pedido (v2).
const (
	EstadoCotizado      = "cotizado"
	EstadoEsperandoPago = "esperando_pago"
	EstadoPagado        = "pagado"
	EstadoPreparando    = "preparando"
	EstadoDespachado    = "despachado"   // solo envío; exige agencia y n.º de seguimiento
	EstadoListoRetiro   = "listo_retiro" // solo retiro
	EstadoEntregado     = "entregado"
	EstadoCancelado     = "cancelado"
)

// Estados en orden.
var Estados = []string{EstadoCotizado, EstadoEsperandoPago, EstadoPagado, EstadoPreparando,
	EstadoDespachado, EstadoListoRetiro, EstadoEntregado, EstadoCancelado}

// EsEstado indica si e es un estado conocido.
func EsEstado(e string) bool {
	for _, v := range Estados {
		if v == e {
			return true
		}
	}
	return false
}

// EsFinal: entregado y cancelado.
func EsFinal(e string) bool { return e == EstadoEntregado || e == EstadoCancelado }

// cotizado -> esperando_pago | pagado -> preparando -> despachado | listo_retiro -> entregado.
// Con forma de pago "en_tienda" (retiro, paga al retirar) se puede pasar a
// preparando sin pagar. Cualquier estado no final puede pasar a cancelado.
var transiciones = map[string][]string{
	EstadoCotizado:      {EstadoEsperandoPago, EstadoPagado, EstadoPreparando},
	EstadoEsperandoPago: {EstadoPagado, EstadoPreparando},
	EstadoPagado:        {EstadoPreparando},
	EstadoPreparando:    {EstadoDespachado, EstadoListoRetiro},
	EstadoDespachado:    {EstadoEntregado},
	EstadoListoRetiro:   {EstadoEntregado},
}

// TransicionesDesde lista los estados a los que puede pasar un pedido.
func TransicionesDesde(entrega, formaPago, de string) []string {
	out := []string{}
	for _, a := range transiciones[de] {
		if a == EstadoPreparando && (de == EstadoCotizado || de == EstadoEsperandoPago) && formaPago != PagoEnTienda {
			continue // sin pagar solo se prepara si paga al retirar
		}
		if a == EstadoDespachado && entrega != EntregaEnvio {
			continue
		}
		if a == EstadoListoRetiro && entrega != EntregaRetiro {
			continue
		}
		out = append(out, a)
	}
	if EsEstado(de) && !EsFinal(de) {
		out = append(out, EstadoCancelado)
	}
	return out
}

// TransicionPermitida indica si se puede pasar de "de" a "a".
func TransicionPermitida(entrega, formaPago, de, a string) bool {
	for _, v := range TransicionesDesde(entrega, formaPago, de) {
		if v == a {
			return true
		}
	}
	return false
}

// ErrTransicionInvalida: cambio de estado no permitido.
type ErrTransicionInvalida struct{ De, A string }

func (e *ErrTransicionInvalida) Error() string {
	return fmt.Sprintf("no se puede pasar de %q a %q", e.De, e.A)
}

// ErrRequisitos: faltan datos que la transición exige (pago, agencia, n.º de seguimiento).
type ErrRequisitos struct{ Campos map[string]string }

func (e *ErrRequisitos) Error() string { return "pedidos: faltan datos para el cambio de estado" }

// Validador es el store.ValidadorCambio: revisa la transición con el pedido
// leído dentro de la transacción y decide el efecto en el stock.
//   - a pagado: exige pago y la confirmación de que la plata está en la cuenta.
//   - a despachado: exige agencia, n.º de seguimiento y comuna (del pedido o del cambio).
//   - a entregado sin pago registrado (pagó al retirar): exige el pago.
//   - stock: se descuenta al pasar a pagado o a preparando (lo que ocurra
//     primero) y se devuelve al cancelar si estaba descontado.
func Validador() store.ValidadorCambio {
	return func(actual store.Pedido, c store.CambioEstado) (store.EfectoStock, error) {
		if !TransicionPermitida(actual.Entrega, actual.FormaPago, actual.Estado, c.Nuevo) {
			return store.StockSinCambio, &ErrTransicionInvalida{De: actual.Estado, A: c.Nuevo}
		}
		campos := map[string]string{}
		switch c.Nuevo {
		case EstadoPagado:
			if c.Pago == nil {
				campos["pago"] = "indica medio y monto del pago"
			}
			if !c.ConfirmoPago {
				campos["confirmo_pago_recibido"] = "confirma que revisaste que la plata está en la cuenta"
			}
		case EstadoDespachado:
			if c.Agencia == "" && actual.Agencia == "" {
				campos["agencia"] = "indica la agencia"
			}
			if c.NSeguimiento == "" {
				campos["n_seguimiento"] = "indica el n.º de seguimiento de la agencia"
			}
			if c.Comuna == "" && actual.Comuna == "" {
				campos["comuna"] = "indica la comuna de destino"
			}
		case EstadoEntregado:
			if actual.Pago == nil && c.Pago == nil {
				campos["pago"] = "registra cómo pagó al retirar (medio y monto)"
			}
		}
		if c.Pago != nil && c.Nuevo != EstadoPagado && (c.Nuevo != EstadoEntregado || actual.Pago != nil) {
			campos["pago"] = "el pago se registra al marcar pagado (o al entregar si pagó al retirar)"
		}
		if len(campos) > 0 {
			return store.StockSinCambio, &ErrRequisitos{Campos: campos}
		}
		switch {
		case (c.Nuevo == EstadoPagado || c.Nuevo == EstadoPreparando) && !actual.StockDescontado:
			return store.StockDescontar, nil
		case c.Nuevo == EstadoCancelado && actual.StockDescontado:
			return store.StockDevolver, nil
		}
		return store.StockSinCambio, nil
	}
}
