// Package pedidos tiene las reglas del pedido: validación de lo que manda el
// cliente, cotización con precios de la base (nunca los del navegador),
// máquina de estados y el texto del mensaje de WhatsApp.
package pedidos

import (
	"fmt"
	"strings"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

// Límites del pedido.
const (
	MaxItems             = 50
	MaxCantidadPorItem   = 100 // por línea del pedido web; el mayor mínimo real de la semilla es 100 (pijama)
	MaxLargoComentario   = 500
	MaxLargoNombre       = 80
	MaxLargoComuna       = 60
	EntregaRetiro        = "retiro"
	EntregaEnvio         = "envio"
	MotivoNoDisponible   = "no_disponible"
	MotivoAgotado        = "agotado"
	MotivoStockInsufic   = "stock_insuficiente"
	errCantidadFueraRang = "debe estar entre 1 y 100"
)

// ItemSolicitud es una línea tal como la manda el navegador.
type ItemSolicitud struct {
	ProductoID int64 `json:"producto_id"`
	Cantidad   int   `json:"cantidad"`
	// Precio que el navegador cree que aplica. Se acepta pero se IGNORA.
	PrecioUnitario *int64 `json:"precio_unitario,omitempty"`
}

// Solicitud es el cuerpo de POST /api/pedidos.
type Solicitud struct {
	Nombre     string          `json:"nombre"`
	Telefono   string          `json:"telefono"`
	Entrega    string          `json:"entrega"`
	Region     string          `json:"region"`
	Comuna     string          `json:"comuna"`
	Agencia    string          `json:"agencia"`
	Comentario string          `json:"comentario"`
	Items      []ItemSolicitud `json:"items"`
	// Total que calculó el navegador. Se acepta pero se IGNORA: el servidor recalcula.
	Total *int64 `json:"total,omitempty"`
}

// Validar normaliza y valida. Devuelve los datos limpios o errores por campo.
func (s Solicitud) Validar() (store.PedidoNuevo, map[string]string) {
	errs := map[string]string{}
	var n store.PedidoNuevo

	nombre, ok := texto.Limpiar(s.Nombre, false)
	switch {
	case !ok || !texto.EsNombrePersona(nombre):
		errs["nombre"] = "escribe tu nombre solo con letras"
	case texto.Largo(nombre) < 2 || texto.Largo(nombre) > MaxLargoNombre:
		errs["nombre"] = "debe tener entre 2 y 80 letras"
	}
	n.Nombre = nombre

	tel, ok := NormalizarTelefono(s.Telefono)
	if !ok {
		errs["telefono"] = "usa un número chileno, por ejemplo +56 9 1234 5678"
	}
	n.Telefono = tel

	n.Entrega = strings.TrimSpace(s.Entrega)
	comuna, okComuna := texto.Limpiar(s.Comuna, false)
	comunaInvalida := !okComuna || (comuna != "" && (!texto.EsLugar(comuna) || texto.Largo(comuna) < 2 || texto.Largo(comuna) > MaxLargoComuna))
	if comunaInvalida {
		errs["comuna"] = "escribe la comuna o ciudad (2 a 60 letras)"
	}
	switch n.Entrega {
	case EntregaRetiro:
		// Retiro en tienda: comuna opcional, sin región ni agencia.
		n.Comuna = comuna
	case EntregaEnvio:
		region := strings.TrimSpace(s.Region)
		if _, ok := RegionPorCodigo(region); !ok {
			errs["region"] = "elige una región de la lista"
		}
		if comuna == "" && !comunaInvalida {
			errs["comuna"] = "escribe la comuna o ciudad de destino"
		}
		agencia := strings.TrimSpace(s.Agencia)
		if _, ok := AgenciaPorCodigo(agencia); !ok {
			errs["agencia"] = "elige una agencia de la lista"
		}
		n.Region, n.Comuna, n.Agencia = region, comuna, agencia
	default:
		errs["entrega"] = "debe ser retiro o envio"
	}

	comentario, ok := texto.Limpiar(s.Comentario, true)
	switch {
	case !ok:
		errs["comentario"] = "tiene caracteres no permitidos"
	case texto.Largo(comentario) > MaxLargoComentario:
		errs["comentario"] = "máximo 500 caracteres"
	}
	n.Comentario = comentario

	switch {
	case len(s.Items) == 0:
		errs["items"] = "el pedido no tiene productos"
	case len(s.Items) > MaxItems:
		errs["items"] = fmt.Sprintf("máximo %d productos distintos por pedido", MaxItems)
	default:
		// Se juntan líneas repetidas del mismo producto, conservando el orden.
		pos := map[int64]int{}
		for i, it := range s.Items {
			campo := fmt.Sprintf("items[%d]", i)
			if it.ProductoID <= 0 {
				errs[campo+".producto_id"] = "producto inválido"
				continue
			}
			if it.Cantidad < 1 || it.Cantidad > MaxCantidadPorItem {
				errs[campo+".cantidad"] = errCantidadFueraRang
				continue
			}
			if j, ok := pos[it.ProductoID]; ok {
				n.Items[j].Cantidad += it.Cantidad
				if n.Items[j].Cantidad > MaxCantidadPorItem {
					errs[campo+".cantidad"] = errCantidadFueraRang
				}
				continue
			}
			pos[it.ProductoID] = len(n.Items)
			n.Items = append(n.Items, store.ItemSolicitado{ProductoID: it.ProductoID, Cantidad: it.Cantidad})
		}
	}

	if len(errs) > 0 {
		return store.PedidoNuevo{}, errs
	}
	return n, nil
}

// Indisponible describe por qué no se puede pedir un producto.
type Indisponible struct {
	ProductoID int64  `json:"producto_id"`
	Motivo     string `json:"motivo"` // no_disponible | agotado | stock_insuficiente
}

// ErrNoDisponibles se devuelve cuando uno o más productos no se pueden pedir.
type ErrNoDisponibles struct {
	Detalles []Indisponible
}

func (e *ErrNoDisponibles) Error() string {
	return fmt.Sprintf("pedidos: %d producto(s) no disponibles", len(e.Detalles))
}

// Cotizador devuelve el store.Cotizador del dominio para la fecha hoy
// (AAAA-MM-DD, Chile; decide si una oferta sigue vigente). Precio y stock salen
// SOLO de la base: lo que el navegador diga del precio o del total no se usa.
// Verifica que el stock alcance, pero no lo descuenta (eso pasa al confirmar).
func Cotizador(hoy string) store.Cotizador {
	return func(prods map[int64]store.Producto, items []store.ItemSolicitado) ([]store.ItemPedido, int64, error) {
		return cotizar(prods, items, hoy)
	}
}

func cotizar(prods map[int64]store.Producto, items []store.ItemSolicitado, hoy string) ([]store.ItemPedido, int64, error) {
	var malos []Indisponible
	out := make([]store.ItemPedido, 0, len(items))
	var total int64
	for _, it := range items {
		p, ok := prods[it.ProductoID]
		switch {
		case !ok || !p.Activo:
			malos = append(malos, Indisponible{it.ProductoID, MotivoNoDisponible})
			continue
		case p.Stock <= 0:
			malos = append(malos, Indisponible{it.ProductoID, MotivoAgotado})
			continue
		case p.Stock < it.Cantidad:
			malos = append(malos, Indisponible{it.ProductoID, MotivoStockInsufic})
			continue
		}
		precio, tipo := productos.PrecioUnitario(p, it.Cantidad, hoy)
		sub := precio * int64(it.Cantidad)
		total += sub
		out = append(out, store.ItemPedido{
			ProductoID:     p.ID,
			Codigo:         p.Codigo,
			Nombre:         p.Nombre,
			Cantidad:       it.Cantidad,
			PrecioUnitario: precio,
			TipoPrecio:     tipo,
			Subtotal:       sub,
		})
	}
	if len(malos) > 0 {
		return nil, 0, &ErrNoDisponibles{Detalles: malos}
	}
	return out, total, nil
}
