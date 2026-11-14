package pedidos

import (
	"fmt"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Estados del pedido.
const (
	EstadoNuevo      = "nuevo"
	EstadoConfirmado = "confirmado"
	EstadoListo      = "listo"      // listo para retiro (solo entrega = retiro)
	EstadoDespachado = "despachado" // entregado a la agencia (solo entrega = envio)
	EstadoEntregado  = "entregado"
	EstadoCancelado  = "cancelado"
)

// Estados en orden, para listas y validación.
var Estados = []string{EstadoNuevo, EstadoConfirmado, EstadoListo, EstadoDespachado, EstadoEntregado, EstadoCancelado}

// EsEstado indica si e es un estado conocido.
func EsEstado(e string) bool {
	for _, v := range Estados {
		if v == e {
			return true
		}
	}
	return false
}

// nuevo -> confirmado -> listo (retiro) | despachado (envío) -> entregado.
// Se puede cancelar desde nuevo, confirmado o listo. entregado y cancelado son finales.
var transiciones = map[string][]string{
	EstadoNuevo:      {EstadoConfirmado, EstadoCancelado},
	EstadoConfirmado: {EstadoListo, EstadoDespachado, EstadoCancelado},
	EstadoListo:      {EstadoEntregado, EstadoCancelado},
	EstadoDespachado: {EstadoEntregado},
}

// TransicionesDesde lista los estados a los que puede pasar un pedido.
func TransicionesDesde(entrega, de string) []string {
	out := []string{}
	for _, a := range transiciones[de] {
		if a == EstadoListo && entrega != EntregaRetiro {
			continue
		}
		if a == EstadoDespachado && entrega != EntregaEnvio {
			continue
		}
		out = append(out, a)
	}
	return out
}

// TransicionPermitida indica si se puede pasar de "de" a "a".
func TransicionPermitida(entrega, de, a string) bool {
	for _, v := range TransicionesDesde(entrega, de) {
		if v == a {
			return true
		}
	}
	return false
}

// ErrTransicionInvalida se devuelve al intentar un cambio de estado no permitido.
type ErrTransicionInvalida struct {
	De, A string
}

func (e *ErrTransicionInvalida) Error() string {
	return fmt.Sprintf("no se puede pasar de %q a %q", e.De, e.A)
}

// Validador devuelve el store.ValidadorCambio para pasar al estado "a".
// Al cancelar se devuelve el stock que el pedido había reservado.
func Validador(a string) store.ValidadorCambio {
	return func(actual store.Pedido) (bool, error) {
		if !TransicionPermitida(actual.Entrega, actual.Estado, a) {
			return false, &ErrTransicionInvalida{De: actual.Estado, A: a}
		}
		return a == EstadoCancelado, nil
	}
}
