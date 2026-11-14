// Package store define el modelo persistido y la interfaz de acceso a datos.
// La implementación actual es SQLite (paquete store/sqlite); cambiar a
// PostgreSQL significa escribir otra implementación de Store, sin tocar los
// handlers ni las reglas de negocio.
package store

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Errores de la capa de datos. Las implementaciones los envuelven o devuelven tal cual.
var (
	ErrNoEncontrado       = errors.New("store: no encontrado")
	ErrDuplicado          = errors.New("store: valor duplicado")
	ErrReferenciaInvalida = errors.New("store: referencia inválida")
	ErrStockInsuficiente  = errors.New("store: stock insuficiente")
	ErrConflicto          = errors.New("store: el registro cambió mientras se actualizaba")
)

// NumeroPedido da el formato público del correlativo: DC-0001, DC-0002, ...
func NumeroPedido(secuencia int64) string {
	return fmt.Sprintf("DC-%04d", secuencia)
}

// Categoria de productos.
type Categoria struct {
	ID             int64
	Slug           string
	Nombre         string
	Orden          int
	TotalProductos int // productos activos y no eliminados (solo en listados)
}

// Producto tal como se guarda. Los campos opcionales usan puntero (precios)
// o string vacío (textos) para representar NULL.
type Producto struct {
	ID              int64
	Slug            string
	Codigo          string
	Nombre          string
	Descripcion     string
	Busqueda        string // texto normalizado para buscar; lo calcula el dominio
	CategoriaID     int64
	CategoriaSlug   string
	CategoriaNombre string
	PrecioDetalle   int64
	PrecioOferta    *int64
	EtiquetaOferta  string // "oferta" | "remate"
	PrecioMayor     *int64 // siempre por unidad
	MayorMinimo     *int   // desde cuántas unidades aplica PrecioMayor
	UnidadMayor     string // etiqueta: "caja de 12", "embalaje de 16"
	Stock           int
	UmbralPocos     int
	Activo          bool
	Temporada       string // "" = sin temporada
	FechaLlegada    string // YYYY-MM-DD o ""
	ImagenURL       string
	TikTokURL       string
	FuenteURL       string // procedencia del precio
	FuenteFecha     string // YYYY-MM-DD
	EsEjemplo       bool   // precio o producto sin fuente real
	CreadoEn        time.Time
	ActualizadoEn   time.Time
}

// FiltroProductos para listados. Busqueda ya viene normalizada por el dominio.
type FiltroProductos struct {
	CategoriaSlug    string
	Palabras         []string // cada palabra debe aparecer (AND)
	SoloOferta       bool
	Temporada        string
	LlegadaDesde     string // YYYY-MM-DD; "" = sin filtro
	IncluirInactivos bool
	Orden            string // "relevancia" | "precio_asc" | "precio_desc" | "recientes"
	Limite           int
	Desplazamiento   int
}

// ItemSolicitado es lo que pide el cliente: producto y cantidad. Nada más.
type ItemSolicitado struct {
	ProductoID int64
	Cantidad   int
}

// ItemPedido es una línea del pedido con el precio congelado al momento de pedir.
type ItemPedido struct {
	ProductoID     int64
	Codigo         string
	Nombre         string
	Cantidad       int
	PrecioUnitario int64
	TipoPrecio     string // "detalle" | "oferta" | "mayor"
	Subtotal       int64
}

// PedidoNuevo son los datos ya validados para crear un pedido.
type PedidoNuevo struct {
	Nombre     string
	Telefono   string // E.164, p. ej. +56912345678
	Entrega    string // "retiro" | "envio"
	Region     string // código de región ("" si retiro)
	Comuna     string
	Agencia    string // código de agencia ("" si retiro)
	Comentario string
	Items      []ItemSolicitado
	EsDemo     bool
	Ahora      time.Time
}

// Cotizador recibe los productos vigentes (por ID) y lo solicitado, y devuelve
// las líneas con precio y el total. Corre DENTRO de la transacción que crea el
// pedido, así precio y stock se leen y se descuentan de forma atómica.
type Cotizador func(productos map[int64]Producto, items []ItemSolicitado) ([]ItemPedido, int64, error)

// Pedido completo.
type Pedido struct {
	ID            int64
	Secuencia     int64
	Numero        string
	Nombre        string
	Telefono      string
	Entrega       string
	Region        string
	Comuna        string
	Agencia       string
	Comentario    string
	Total         int64
	Estado        string
	EsDemo        bool
	CreadoEn      time.Time
	ActualizadoEn time.Time
	Items         []ItemPedido
	Historial     []CambioRegistrado
}

// CambioRegistrado es una fila de historial_estados.
type CambioRegistrado struct {
	EstadoAnterior string // "" para la creación
	EstadoNuevo    string
	AdminUsuario   string // "" si fue el sistema o la creación
	Nota           string
	CreadoEn       time.Time
}

// CambioEstado es la solicitud de transición.
type CambioEstado struct {
	PedidoID int64
	Nuevo    string
	AdminID  int64 // 0 = sistema (semilla)
	Nota     string
	Ahora    time.Time
}

// ValidadorCambio decide, con el pedido leído dentro de la transacción, si la
// transición es válida y si hay que devolver el stock reservado.
type ValidadorCambio func(actual Pedido) (restaurarStock bool, err error)

// FiltroPedidos para el panel.
type FiltroPedidos struct {
	Estado         string // "" = todos
	Limite         int
	Desplazamiento int
}

// Admin del panel.
type Admin struct {
	ID           int64
	Usuario      string
	PasswordHash string
}

// Sesion de admin. El token nunca se guarda: solo su HMAC.
type Sesion struct {
	TokenHash string
	AdminID   int64
	Usuario   string
	CreadoEn  time.Time
	ExpiraEn  time.Time
	UltimoUso time.Time
}

// Store es el contrato de persistencia.
type Store interface {
	Ping(ctx context.Context) error
	Close() error

	// Catálogo
	ListarCategorias(ctx context.Context) ([]Categoria, error)
	CrearCategoria(ctx context.Context, c *Categoria, ahora time.Time) error
	ExisteCategoria(ctx context.Context, id int64) (bool, error)
	ListarProductos(ctx context.Context, f FiltroProductos) ([]Producto, int, error)
	ProductoPorSlug(ctx context.Context, slug string, incluirInactivos bool) (Producto, error)
	ProductoPorID(ctx context.Context, id int64) (Producto, error)
	CrearProducto(ctx context.Context, p *Producto, ahora time.Time) error
	ActualizarProducto(ctx context.Context, p *Producto, ahora time.Time) error
	EliminarProducto(ctx context.Context, id int64, ahora time.Time) error

	// Pedidos
	CrearPedido(ctx context.Context, n PedidoNuevo, cotizar Cotizador) (Pedido, error)
	ListarPedidos(ctx context.Context, f FiltroPedidos) ([]Pedido, int, error)
	PedidoPorID(ctx context.Context, id int64) (Pedido, error)
	CambiarEstado(ctx context.Context, c CambioEstado, validar ValidadorCambio) (Pedido, error)

	// Admin y sesiones
	SincronizarAdmin(ctx context.Context, usuario, passwordHash string, ahora time.Time) (int64, error)
	AdminPorUsuario(ctx context.Context, usuario string) (Admin, error)
	CrearSesion(ctx context.Context, s Sesion) error
	SesionPorHash(ctx context.Context, tokenHash string) (Sesion, error)
	TocarSesion(ctx context.Context, tokenHash string, ahora time.Time) error
	EliminarSesion(ctx context.Context, tokenHash string) error
	PurgarSesiones(ctx context.Context, ahora time.Time) (int64, error)

	// Semilla
	CatalogoVacio(ctx context.Context) (bool, error)
}
