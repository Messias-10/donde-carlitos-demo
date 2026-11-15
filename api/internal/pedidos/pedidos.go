// Package pedidos tiene las reglas del pedido registrado en el panel (v2):
// validación de lo que escribe el equipo, cotización con precios de la base,
// líneas libres (con o sin precio), descuento, costo de envío y la máquina de
// estados (pago, despacho, stock).
package pedidos

import (
	"fmt"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

// Límites y valores.
const (
	MaxItems           = 50
	MaxCantidad        = 10000
	MaxMonto           = 100_000_000
	MaxLargoNombre     = 80
	MaxLargoComuna     = 60
	MaxLargoNota       = 1000
	MaxLargoDescLibre  = 120
	MaxLargoMotivo     = 200
	CanalSantiago      = "santiago"
	CanalRegiones      = "regiones"
	EntregaRetiro      = "retiro"
	EntregaEnvio       = "envio"
	PagoTransferencia  = "transferencia"
	PagoEnTienda       = "en_tienda"
	TipoLibre          = "libre"
	TipoManual         = "manual" // línea del catálogo con precio editado por quien atiende
	errCantidadFueraRg = "debe estar entre 1 y 10.000"
)

// EsCanal indica si c es un canal válido.
func EsCanal(c string) bool { return c == CanalSantiago || c == CanalRegiones }

// ItemEntrada es una línea tal como la escribe el equipo.
type ItemEntrada struct {
	ProductoID     int64  `json:"producto_id"`
	Descripcion    string `json:"descripcion"`
	Cantidad       int    `json:"cantidad"`
	PrecioUnitario *int64 `json:"precio_unitario"`
}

// Entrada es el cuerpo de POST y PUT /api/admin/pedidos.
type Entrada struct {
	Canal         string        `json:"canal"`
	Nombre        string        `json:"nombre"`
	Telefono      string        `json:"telefono"`
	Entrega       string        `json:"entrega"`
	Region        string        `json:"region"`
	Comuna        string        `json:"comuna"`
	Agencia       string        `json:"agencia"`
	FormaPago     string        `json:"forma_pago"`
	Items         []ItemEntrada `json:"items"`
	Descuento     int64         `json:"descuento"`
	CostoEnvio    int64         `json:"costo_envio"`
	EnvioPorPagar bool          `json:"envio_por_pagar"`
	Total         *int64        `json:"total"`
	NotaInterna   string        `json:"nota_interna"`
	EstadoInicial string        `json:"estado_inicial"`
}

// Digitos deja solo los dígitos (para buscar por teléfono).
func Digitos(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Validar normaliza y valida. Devuelve los datos, el estado inicial y errores por campo.
// Por defecto: Santiago → retiro y pago en tienda; Regiones → envío y transferencia.
func (e Entrada) Validar() (store.DatosPedido, string, map[string]string) {
	errs := map[string]string{}
	var d store.DatosPedido

	d.Canal = strings.TrimSpace(e.Canal)
	if !EsCanal(d.Canal) {
		errs["canal"] = "debe ser santiago o regiones"
	}

	// El nombre es opcional (a veces el cliente no lo da); si viene, se valida.
	nombre, ok := texto.Limpiar(e.Nombre, false)
	switch {
	case !ok || (nombre != "" && !texto.EsNombrePersona(nombre)):
		errs["nombre"] = "solo letras, espacios, apóstrofo, punto y guion"
	case texto.Largo(nombre) > MaxLargoNombre:
		errs["nombre"] = "máximo 80 letras"
	}
	d.Nombre = nombre

	tel, ok := NormalizarTelefono(e.Telefono)
	if !ok {
		errs["telefono"] = "número chileno, por ejemplo +56 9 1234 5678"
	}
	d.Telefono = tel

	d.Entrega = strings.TrimSpace(e.Entrega)
	if d.Entrega == "" {
		d.Entrega = EntregaRetiro
		if d.Canal == CanalRegiones {
			d.Entrega = EntregaEnvio
		}
	}
	comuna, okComuna := texto.Limpiar(e.Comuna, false)
	if !okComuna || (comuna != "" && (!texto.EsLugar(comuna) || texto.Largo(comuna) < 2 || texto.Largo(comuna) > MaxLargoComuna)) {
		errs["comuna"] = "comuna o ciudad de 2 a 60 letras"
	}
	d.Comuna = comuna
	if r := strings.TrimSpace(e.Region); r != "" {
		if _, ok := RegionPorCodigo(r); !ok {
			errs["region"] = "elige una región de la lista"
		}
		d.Region = r
	}
	if a := strings.TrimSpace(e.Agencia); a != "" {
		if _, ok := AgenciaPorCodigo(a); !ok {
			errs["agencia"] = "elige una agencia de la lista"
		}
		d.Agencia = a
	}
	switch d.Entrega {
	case EntregaRetiro:
		d.Region, d.Agencia = "", "" // no aplican
	case EntregaEnvio:
	default:
		errs["entrega"] = "debe ser retiro o envio"
	}

	d.FormaPago = strings.TrimSpace(e.FormaPago)
	if d.FormaPago == "" {
		d.FormaPago = PagoEnTienda
		if d.Canal == CanalRegiones {
			d.FormaPago = PagoTransferencia
		}
	}
	if d.FormaPago != PagoTransferencia && d.FormaPago != PagoEnTienda {
		errs["forma_pago"] = "transferencia o en_tienda"
	}
	if d.FormaPago == PagoEnTienda && d.Entrega == EntregaEnvio {
		errs["forma_pago"] = "con envío el pago es por transferencia"
	}

	if e.Descuento < 0 || e.Descuento > MaxMonto {
		errs["descuento"] = "monto en pesos, 0 o más"
	}
	if e.CostoEnvio < 0 || e.CostoEnvio > MaxMonto {
		errs["costo_envio"] = "monto en pesos, 0 o más"
	}
	if e.CostoEnvio > 0 && d.Entrega != EntregaEnvio {
		errs["costo_envio"] = "solo para envío"
	}
	d.Descuento, d.CostoEnvio = e.Descuento, e.CostoEnvio
	// "Por pagar" solo tiene sentido en envío y sin costo de envío cobrado.
	d.EnvioPorPagar = e.EnvioPorPagar && d.Entrega == EntregaEnvio && e.CostoEnvio == 0
	if e.Total != nil && (*e.Total < 0 || *e.Total > MaxMonto) {
		errs["total"] = "monto en pesos, 0 o más"
	}
	d.TotalManual = e.Total

	nota, ok := texto.Limpiar(e.NotaInterna, true)
	if !ok || texto.Largo(nota) > MaxLargoNota {
		errs["nota_interna"] = "máximo 1000 caracteres"
	}
	d.NotaInterna = nota

	sinPrecio := false
	switch {
	case len(e.Items) == 0:
		errs["items"] = "el pedido no tiene productos"
	case len(e.Items) > MaxItems:
		errs["items"] = fmt.Sprintf("máximo %d líneas", MaxItems)
	default:
		pos := map[int64]int{}
		for i, it := range e.Items {
			campo := fmt.Sprintf("items[%d]", i)
			if it.Cantidad < 1 || it.Cantidad > MaxCantidad {
				errs[campo+".cantidad"] = errCantidadFueraRg
				continue
			}
			if it.ProductoID > 0 {
				if strings.TrimSpace(it.Descripcion) != "" {
					errs[campo] = "una línea del catálogo no lleva descripcion"
					continue
				}
				// Precio editado a mano (opcional): se valida y no se junta con otras líneas.
				if it.PrecioUnitario != nil {
					if *it.PrecioUnitario < 0 || *it.PrecioUnitario > MaxMonto {
						errs[campo+".precio_unitario"] = "precio en pesos, 0 o más"
						continue
					}
					v := *it.PrecioUnitario
					d.Items = append(d.Items, store.ItemSolicitado{ProductoID: it.ProductoID, Cantidad: it.Cantidad, PrecioUnitario: &v})
					continue
				}
				if j, ok := pos[it.ProductoID]; ok { // se juntan líneas del mismo producto
					d.Items[j].Cantidad += it.Cantidad
					if d.Items[j].Cantidad > MaxCantidad {
						errs[campo+".cantidad"] = errCantidadFueraRg
					}
					continue
				}
				pos[it.ProductoID] = len(d.Items)
				d.Items = append(d.Items, store.ItemSolicitado{ProductoID: it.ProductoID, Cantidad: it.Cantidad})
				continue
			}
			if it.ProductoID < 0 {
				errs[campo+".producto_id"] = "producto inválido"
				continue
			}
			desc, ok := texto.Limpiar(it.Descripcion, false)
			if !ok || texto.Largo(desc) < 2 || texto.Largo(desc) > MaxLargoDescLibre {
				errs[campo+".descripcion"] = "describe la línea (2 a 120 caracteres)"
				continue
			}
			if it.PrecioUnitario != nil && (*it.PrecioUnitario < 0 || *it.PrecioUnitario > MaxMonto) {
				errs[campo+".precio_unitario"] = "precio en pesos, 0 o más"
				continue
			}
			sinPrecio = sinPrecio || it.PrecioUnitario == nil
			d.Items = append(d.Items, store.ItemSolicitado{Descripcion: desc, Cantidad: it.Cantidad, PrecioUnitario: it.PrecioUnitario})
		}
	}
	if sinPrecio && e.Total == nil {
		errs["total"] = "hay líneas sin precio: escribe el total"
	}

	inicial := strings.TrimSpace(e.EstadoInicial)
	if inicial == "" {
		inicial = EstadoCotizado
	}
	if inicial != EstadoCotizado && inicial != EstadoEsperandoPago {
		errs["estado_inicial"] = "debe ser cotizado o esperando_pago"
	}

	if len(errs) > 0 {
		return store.DatosPedido{}, "", errs
	}
	d.Busqueda = strings.TrimSpace(productos.NormalizarBusqueda(d.Nombre) + " " + Digitos(d.Telefono))
	return d, inicial, nil
}

// ErrItems son líneas de catálogo que no se pueden usar (producto inexistente o inactivo).
type ErrItems struct {
	Campos map[string]string
}

func (e *ErrItems) Error() string {
	return fmt.Sprintf("pedidos: %d línea(s) inválidas", len(e.Campos))
}

// Cotizador devuelve el store.Cotizador para la fecha hoy (AAAA-MM-DD, Chile).
// Catálogo: precio de la base (el menor entre detalle, oferta vigente y, si se
// alcanza el mínimo, precio por mayor), congelado. Libre: el precio escrito, o
// sin precio. El subtotal suma solo las líneas con precio. No mira el stock:
// se descuenta (y se verifica) al pasar a "pagado".
func Cotizador(hoy string) store.Cotizador {
	return func(prods map[int64]store.Producto, items []store.ItemSolicitado) ([]store.ItemPedido, int64, error) {
		campos := map[string]string{}
		out := make([]store.ItemPedido, 0, len(items))
		var subtotal int64
		for i, it := range items {
			if it.ProductoID == 0 {
				l := store.ItemPedido{Descripcion: it.Descripcion, Cantidad: it.Cantidad, TipoPrecio: TipoLibre}
				if it.PrecioUnitario != nil {
					precio, sub := *it.PrecioUnitario, *it.PrecioUnitario*int64(it.Cantidad)
					l.PrecioUnitario, l.Subtotal = &precio, &sub
					subtotal += sub
				}
				out = append(out, l)
				continue
			}
			p, ok := prods[it.ProductoID]
			if !ok || !p.Activo {
				campos[fmt.Sprintf("items[%d].producto_id", i)] = "el producto no existe o está inactivo"
				continue
			}
			precio, tipo := productos.PrecioUnitario(p, it.Cantidad, hoy)
			if it.PrecioUnitario != nil { // precio editado por quien atiende
				precio, tipo = *it.PrecioUnitario, TipoManual
			}
			sub := precio * int64(it.Cantidad)
			subtotal += sub
			out = append(out, store.ItemPedido{ProductoID: p.ID, Codigo: p.Codigo, Descripcion: p.Nombre,
				Cantidad: it.Cantidad, PrecioUnitario: &precio, TipoPrecio: tipo, Subtotal: &sub})
		}
		if len(campos) > 0 {
			return nil, 0, &ErrItems{Campos: campos}
		}
		return out, subtotal, nil
	}
}

// ---------------------------------------------------------------- cambio de estado

// MediosPago válidos (cómo se pagó de verdad).
var MediosPago = []string{"transferencia", "efectivo", "debito", "credito", "otro"}

// PagoEntrada es el pago que se informa al pasar a "pagado".
type PagoEntrada struct {
	Medio string `json:"medio"`
	Monto int64  `json:"monto"`
	Fecha string `json:"fecha"`
}

// CambioEntrada es el cuerpo de PATCH /api/admin/pedidos/{id}/estado.
type CambioEntrada struct {
	Estado               string       `json:"estado"`
	Nota                 string       `json:"nota"`
	Pago                 *PagoEntrada `json:"pago"`
	ConfirmoPagoRecibido bool         `json:"confirmo_pago_recibido"`
	Agencia              string       `json:"agencia"`
	NSeguimiento         string       `json:"n_seguimiento"`
	Comuna               string       `json:"comuna"`
	Motivo               string       `json:"motivo"`
}

// Validar revisa forma y valores (los requisitos que dependen del pedido los
// revisa Validador dentro de la transacción).
func (c CambioEntrada) Validar(hoy time.Time) (store.CambioEstado, map[string]string) {
	errs := map[string]string{}
	var out store.CambioEstado
	if !EsEstado(c.Estado) {
		errs["estado"] = "estado desconocido"
	}
	out.Nuevo = c.Estado
	nota, ok := texto.Limpiar(c.Nota, false)
	if !ok || texto.Largo(nota) > 200 {
		errs["nota"] = "máximo 200 caracteres"
	}
	out.Nota = nota
	out.ConfirmoPago = c.ConfirmoPagoRecibido
	if c.Pago != nil {
		if c.Estado != EstadoPagado && c.Estado != EstadoEntregado {
			errs["pago"] = "solo al pasar a pagado (o a entregado si pagó al retirar)"
		}
		medioOK := false
		for _, m := range MediosPago {
			medioOK = medioOK || m == c.Pago.Medio
		}
		if !medioOK {
			errs["pago.medio"] = "transferencia, efectivo, debito, credito u otro"
		}
		if c.Pago.Monto < 1 || c.Pago.Monto > MaxMonto {
			errs["pago.monto"] = "monto en pesos mayor que 0"
		}
		fecha := strings.TrimSpace(c.Pago.Fecha)
		if fecha == "" {
			fecha = hoy.Format(time.DateOnly)
		} else if f, err := time.Parse(time.DateOnly, fecha); err != nil || f.After(hoy) {
			errs["pago.fecha"] = "fecha AAAA-MM-DD, no futura"
		}
		out.Pago = &store.Pago{Medio: c.Pago.Medio, Monto: c.Pago.Monto, Fecha: fecha}
	}
	if a := strings.TrimSpace(c.Agencia); a != "" {
		if _, ok := AgenciaPorCodigo(a); !ok {
			errs["agencia"] = "elige una agencia de la lista"
		}
		out.Agencia = a
	}
	if n := strings.TrimSpace(c.NSeguimiento); n != "" {
		limpio, ok := texto.Limpiar(n, false)
		if !ok || texto.Largo(limpio) > 60 || !esCodigoSeguimiento(limpio) {
			errs["n_seguimiento"] = "letras, números, espacios y guiones (máx. 60)"
		}
		out.NSeguimiento = limpio
	}
	if co := strings.TrimSpace(c.Comuna); co != "" {
		limpio, ok := texto.Limpiar(co, false)
		if !ok || !texto.EsLugar(limpio) || texto.Largo(limpio) < 2 || texto.Largo(limpio) > MaxLargoComuna {
			errs["comuna"] = "comuna o ciudad de 2 a 60 letras"
		}
		out.Comuna = limpio
	}
	if m := strings.TrimSpace(c.Motivo); m != "" {
		if c.Estado != EstadoCancelado {
			errs["motivo"] = "solo al cancelar"
		}
		limpio, ok := texto.Limpiar(m, false)
		if !ok || texto.Largo(limpio) > MaxLargoMotivo {
			errs["motivo"] = "máximo 200 caracteres"
		}
		out.Motivo = limpio
	}
	if len(errs) > 0 {
		return store.CambioEstado{}, errs
	}
	return out, nil
}

func esCodigoSeguimiento(s string) bool {
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9', r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r == ' ', r == '-', r == '.', r == '/':
		default:
			return false
		}
	}
	return s != ""
}
