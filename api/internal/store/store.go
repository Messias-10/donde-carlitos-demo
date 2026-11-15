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
	ErrNoEditable         = errors.New("store: el pedido ya no se puede editar")
	ErrTotalNegativo      = errors.New("store: el total quedaría negativo")
	ErrTotalNoCoincide    = errors.New("store: el total no coincide con las líneas, el descuento y el envío")
)

// ErrSinStock indica qué productos no alcanzan al marcar un pedido como pagado.
// errors.Is(err, ErrStockInsuficiente) es verdadero.
type ErrSinStock struct {
	ProductoIDs []int64
}

func (e *ErrSinStock) Error() string { return ErrStockInsuficiente.Error() }

// Unwrap permite errors.Is(err, ErrStockInsuficiente).
func (e *ErrSinStock) Unwrap() error { return ErrStockInsuficiente }

// EfectoStock es lo que un cambio de estado hace con el stock.
type EfectoStock int

// Efectos posibles.
const (
	StockSinCambio EfectoStock = iota
	StockDescontar             // al pasar a pagado: se descuenta lo pedido
	StockDevolver              // al cancelar un pedido con stock descontado: se devuelve
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
	// MayorMinimoConfirmado: el mínimo salió de la fuente (true) o es [EJEMPLO] (false).
	MayorMinimoConfirmado bool
	OfertaHasta           string     // AAAA-MM-DD (hora de Chile) o "" = sin vencimiento
	PrecioActualizadoEn   *time.Time // último cambio de precio hecho en el admin (nil = precio de la fuente)
	PrecioEsPromo         bool       // el precio es de promo por comentario (no de detalle)
	NotaPromo             string     // p. ej. "Comenta 'yo quiero mi promoción' en el video"
	EnVitrina             bool       // aparece en la vitrina de la landing (máx. 8)
	OrdenVitrina          *int       // 1..99; nil = al final
	Stock                 int
	UmbralPocos           int
	Activo                bool
	Temporada             string // "" = sin temporada
	FechaLlegada          string // YYYY-MM-DD o ""
	ImagenURL             string
	TikTokURL             string
	FuenteURL             string // procedencia del precio
	FuenteFecha           string // YYYY-MM-DD
	EsEjemplo             bool   // precio o producto sin fuente real
	CreadoEn              time.Time
	ActualizadoEn         time.Time
}

// FiltroProductos para listados. Busqueda ya viene normalizada por el dominio.
type FiltroProductos struct {
	CategoriaSlug    string
	Palabras         []string // cada palabra debe aparecer (AND)
	SoloOferta       bool
	Temporada        string
	LlegadaDesde     string // YYYY-MM-DD; "" = sin filtro
	Hoy              string // YYYY-MM-DD (Chile): para saber qué ofertas siguen vigentes
	IncluirInactivos bool
	SoloVitrina      bool
	Orden            string // "relevancia" | "precio_asc" | "precio_desc" | "recientes" | "vitrina"
	Limite           int
	Desplazamiento   int
}

// ---------------------------------------------------------------- pedidos (v2: panel)

// ItemSolicitado es una línea como la registra el equipo: del catálogo
// (ProductoID > 0) o libre (ProductoID == 0, con Descripcion y PrecioUnitario).
type ItemSolicitado struct {
	ProductoID     int64
	Cantidad       int
	Descripcion    string // solo libre
	PrecioUnitario *int64 // solo libre; nil = línea sin precio
}

// ItemPedido es una línea con el precio congelado.
type ItemPedido struct {
	ProductoID     int64  // 0 = línea libre
	Codigo         string // "" en línea libre
	Descripcion    string
	Cantidad       int
	PrecioUnitario *int64 // nil = línea libre sin precio
	TipoPrecio     string // "detalle" | "oferta" | "mayor" | "libre"
	Subtotal       *int64 // nil si no hay precio
}

// Cotizador recibe los productos del catálogo usados (por ID) y las líneas
// pedidas, y devuelve las líneas con precio congelado y el subtotal. Corre
// DENTRO de la transacción que crea o edita el pedido.
type Cotizador func(productos map[int64]Producto, items []ItemSolicitado) ([]ItemPedido, int64, error)

// DatosPedido son los datos editables, ya validados por el dominio.
type DatosPedido struct {
	Canal         string // "santiago" | "regiones"
	Nombre        string
	Telefono      string // E.164
	Busqueda      string // número se agrega en el store
	Entrega       string // "retiro" | "envio"
	Region        string
	Comuna        string
	Agencia       string
	FormaPago     string // "transferencia" | "en_tienda"
	Items         []ItemSolicitado
	Descuento     int64
	CostoEnvio    int64
	EnvioPorPagar bool
	// TotalManual: total escrito por quien atiende. Obligatorio si hay líneas
	// sin precio; si todas tienen precio y viene, debe coincidir con el cálculo.
	TotalManual *int64
	NotaInterna string
}

// PedidoNuevo agrega lo que se fija al crear.
type PedidoNuevo struct {
	DatosPedido
	EstadoInicial    string
	TokenSeguimiento string
	RegistradoPor    int64 // 0 = sistema (semilla)
	EsDemo           bool
	Ahora            time.Time
}

// UsuarioRef identifica a quien hizo algo (para mostrar su nombre).
type UsuarioRef struct {
	ID     int64
	Nombre string
}

// Pago registrado al pasar a "pagado".
type Pago struct {
	Medio         string
	Monto         int64
	Fecha         string // AAAA-MM-DD
	VerificadoPor *UsuarioRef
}

// Pedido completo.
type Pedido struct {
	ID                int64
	Secuencia         int64
	Numero            string
	TokenSeguimiento  string
	Canal             string
	Nombre            string
	Telefono          string
	Entrega           string
	Region            string
	Comuna            string
	Agencia           string
	NSeguimiento      string
	FormaPago         string
	Subtotal          int64
	Descuento         int64
	CostoEnvio        int64
	EnvioPorPagar     bool
	Total             int64
	TotalManual       bool
	MotivoCancelacion string
	NotaInterna       string
	Estado            string
	EstadoDesde       time.Time
	StockDescontado   bool
	Pago              *Pago
	RegistradoPor     *UsuarioRef
	EsDemo            bool
	TokenVenceEn      *time.Time // el link de seguimiento vence (30 días tras entregado o cancelado)
	CreadoEn          time.Time
	ActualizadoEn     time.Time
	Items             []ItemPedido
	Historial         []CambioRegistrado
	// AvisoStock: productos que no alcanzaban al descontar en modo no estricto
	// (el stock quedó en 0). Solo lo llena CambiarEstado; no se guarda.
	AvisoStock []int64
}

// CambioRegistrado es una fila de historial_estados.
type CambioRegistrado struct {
	EstadoAnterior string // "" para la creación
	EstadoNuevo    string
	Usuario        *UsuarioRef
	Nota           string
	CreadoEn       time.Time
}

// CambioEstado es la solicitud de transición, con los datos que algunas exigen.
type CambioEstado struct {
	PedidoID     int64
	Nuevo        string
	UsuarioID    int64 // 0 = sistema (semilla)
	Nota         string
	Pago         *Pago  // para "pagado" (VerificadoPor lo pone el store con UsuarioID)
	ConfirmoPago bool   // para "pagado": quien atiende confirmó que la plata está en la cuenta
	Agencia      string // para "despachado" ("" = mantener la del pedido)
	NSeguimiento string // para "despachado"
	Comuna       string // para "despachado" ("" = mantener la del pedido)
	Motivo       string // para "cancelado"
	// StockEstricto: si el stock no alcanza al descontar, rechaza (ErrSinStock).
	// Si es false, descuenta hasta 0 y devuelve el aviso en Pedido.AvisoStock.
	StockEstricto bool
	Ahora         time.Time
}

// ValidadorCambio decide, con el pedido leído dentro de la transacción, si la
// transición es válida y qué pasa con el stock.
type ValidadorCambio func(actual Pedido, c CambioEstado) (EfectoStock, error)

// FiltroPedidos para búsquedas, la pantalla Hoy, el CSV y el historial de un cliente.
type FiltroPedidos struct {
	Canal            string
	Estados          []string
	Palabras         []string // en "busqueda" (AND)
	Telefono         string   // exacto (E.164)
	CreadoDesde      time.Time
	CreadoHasta      time.Time // exclusivo
	EstadoDesdeAntes time.Time // estado_desde < este instante
	Limite           int
	Desplazamiento   int
	ConHistorial     bool
}

// FilaResumen es un grupo (canal, estado, agencia) con su conteo y monto.
type FilaResumen struct {
	Canal   string
	Estado  string
	Entrega string
	Agencia string
	Comuna  string
	Pagado  bool // tiene pago registrado
	Pedidos int
	Monto   int64
}

// FilaUsuario: pedidos registrados y despachados por usuario en un rango.
type FilaUsuario struct {
	UsuarioID   int64
	Nombre      string
	Registrados int
	Despachados int
}

// ---------------------------------------------------------------- usuarios

// Roles del panel.
const (
	RolAdmin    = "admin"
	RolVendedor = "vendedor"
)

// Usuario del panel.
type Usuario struct {
	ID               int64
	Usuario          string
	Nombre           string
	PasswordHash     string
	Rol              string
	Activo           bool
	CanalPorDefecto  string // "" | "santiago" | "regiones"
	DebeCambiarClave bool   // clave temporal: hay que cambiarla antes de usar el panel
	CreadoEn         time.Time
}

// CambiosUsuario: nil = no cambia.
type CambiosUsuario struct {
	Nombre          *string
	Rol             *string
	Activo          *bool
	CanalPorDefecto *string // "" = sin canal
}

// Sesion de un usuario activo. El token nunca se guarda: solo su HMAC.
type Sesion struct {
	TokenHash        string
	UsuarioID        int64
	Usuario          string
	Nombre           string
	Rol              string
	CanalPorDefecto  string
	DebeCambiarClave bool
	CreadoEn         time.Time
	ExpiraEn         time.Time
	UltimoUso        time.Time
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
	// ActualizarImagen cambia solo imagen_url y devuelve la anterior.
	ActualizarImagen(ctx context.Context, id int64, url string, ahora time.Time) (anterior string, err error)
	CrearProducto(ctx context.Context, p *Producto, ahora time.Time) error
	ActualizarProducto(ctx context.Context, p *Producto, ahora time.Time) error
	EliminarProducto(ctx context.Context, id int64, ahora time.Time) error

	// Pedidos
	CrearPedido(ctx context.Context, n PedidoNuevo, cotizar Cotizador) (Pedido, error)
	// EditarPedido reemplaza datos e ítems si el pedido sigue en cotizado o
	// esperando_pago (si no, ErrNoEditable).
	EditarPedido(ctx context.Context, id int64, d DatosPedido, cotizar Cotizador, ahora time.Time) (Pedido, error)
	ListarPedidos(ctx context.Context, f FiltroPedidos) ([]Pedido, int, error)
	PedidoPorID(ctx context.Context, id int64) (Pedido, error)
	PedidoPorToken(ctx context.Context, token string) (Pedido, error)
	CambiarEstado(ctx context.Context, c CambioEstado, validar ValidadorCambio) (Pedido, error)
	ResumenPedidos(ctx context.Context, desde, hasta time.Time) ([]FilaResumen, error)
	ResumenUsuarios(ctx context.Context, desde, hasta time.Time) ([]FilaUsuario, error)
	// RegenerarToken pone un token de seguimiento nuevo (vence a los 30 días si el pedido ya terminó).
	RegenerarToken(ctx context.Context, id int64, token string, ahora time.Time) (Pedido, error)
	// ContarVitrina cuenta productos en vitrina (no eliminados) sin contar excluirID.
	ContarVitrina(ctx context.Context, excluirID int64) (int, error)

	// Usuarios y sesiones
	SincronizarAdmin(ctx context.Context, usuario, passwordHash string, ahora time.Time) (int64, error)
	UsuarioPorUsuario(ctx context.Context, usuario string) (Usuario, error)
	UsuarioPorID(ctx context.Context, id int64) (Usuario, error)
	ListarUsuarios(ctx context.Context) ([]Usuario, error)
	CrearUsuario(ctx context.Context, u *Usuario, ahora time.Time) error
	// ActualizarUsuario aplica cambios; si se desactiva, borra sus sesiones.
	ActualizarUsuario(ctx context.Context, id int64, c CambiosUsuario, ahora time.Time) (Usuario, error)
	// CambiarClave reemplaza el hash, fija debe_cambiar_clave y borra sus
	// sesiones salvo conservarSesion (el hash de la sesión actual, o "").
	CambiarClave(ctx context.Context, id int64, passwordHash string, debeCambiar bool, conservarSesion string, ahora time.Time) error
	CrearSesion(ctx context.Context, s Sesion) error
	// SesionPorHash solo encuentra sesiones de usuarios activos.
	SesionPorHash(ctx context.Context, tokenHash string) (Sesion, error)
	TocarSesion(ctx context.Context, tokenHash string, ahora time.Time) error
	EliminarSesion(ctx context.Context, tokenHash string) error
	PurgarSesiones(ctx context.Context, ahora time.Time) (int64, error)

	// Semilla
	CatalogoVacio(ctx context.Context) (bool, error)
	PedidosVacio(ctx context.Context) (bool, error)
	// FijarCorrelativo pone el contador de pedidos (solo la semilla, con la tabla vacía).
	FijarCorrelativo(ctx context.Context, valor int64) error
}
