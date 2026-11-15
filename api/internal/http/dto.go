package httpapi

import (
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/mensajes"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Respuestas JSON. Los nombres coinciden con docs/openapi.yaml.

type categoriaDTO struct {
	ID             int64  `json:"id"`
	Slug           string `json:"slug"`
	Nombre         string `json:"nombre"`
	TotalProductos int    `json:"total_productos"`
}

type categoriaRef struct {
	ID     int64  `json:"id"`
	Slug   string `json:"slug"`
	Nombre string `json:"nombre"`
}

type fuenteDTO struct {
	URL   string  `json:"url"`
	Fecha *string `json:"fecha"`
}

type productoPublicoDTO struct {
	ID                    int64        `json:"id"`
	Slug                  string       `json:"slug"`
	Codigo                string       `json:"codigo"`
	Nombre                string       `json:"nombre"`
	Descripcion           string       `json:"descripcion"`
	Categoria             categoriaRef `json:"categoria"`
	Precio                int64        `json:"precio"`
	PrecioOferta          *int64       `json:"precio_oferta"`
	PrecioFinal           int64        `json:"precio_final"`
	EtiquetaOferta        *string      `json:"etiqueta_oferta"`
	PrecioMayor           *int64       `json:"precio_mayor"`
	MayorMinimo           *int         `json:"mayor_minimo"`
	UnidadMayor           *string      `json:"unidad_mayor"`
	MayorMinimoConfirmado bool         `json:"mayor_minimo_confirmado"`
	OfertaHasta           *string      `json:"oferta_hasta"`
	OfertaVigente         bool         `json:"oferta_vigente"`
	PrecioActualizadoEn   *time.Time   `json:"precio_actualizado_en"`
	Disponibilidad        string       `json:"disponibilidad"`
	Badges                []string     `json:"badges"`
	Temporada             *string      `json:"temporada"`
	FechaLlegada          *string      `json:"fecha_llegada"`
	ImagenURL             *string      `json:"imagen_url"`
	TikTokURL             *string      `json:"tiktok_url"`
	Fuente                *fuenteDTO   `json:"fuente"`
	EsEjemplo             bool         `json:"es_ejemplo"`
}

type productoAdminDTO struct {
	productoPublicoDTO
	Stock         int       `json:"stock"`
	UmbralPocos   int       `json:"umbral_pocos"`
	Activo        bool      `json:"activo"`
	CreadoEn      time.Time `json:"creado_en"`
	ActualizadoEn time.Time `json:"actualizado_en"`
}

type paginaDTO[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func nuevaPagina[T any](items []T, page, size, total int) paginaDTO[T] {
	if items == nil {
		items = []T{}
	}
	tp := 0
	if size > 0 {
		tp = (total + size - 1) / size
	}
	return paginaDTO[T]{Items: items, Page: page, PageSize: size, Total: total, TotalPages: tp}
}

func ptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *Server) hoy() time.Time {
	return s.ahora().In(s.cfg.Zona)
}

// hoyISO es la fecha de hoy en Chile (AAAA-MM-DD): decide si una oferta sigue vigente.
func (s *Server) hoyISO() string { return s.hoy().Format(time.DateOnly) }

func utcPtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	u := t.UTC()
	return &u
}

// productoPublico: si la oferta venció, precio_oferta y etiqueta_oferta van en null.
func (s *Server) productoPublico(p store.Producto) productoPublicoDTO {
	hoy := s.hoyISO()
	vigente := productos.OfertaVigente(p, hoy)
	d := productoPublicoDTO{
		ID:                    p.ID,
		Slug:                  p.Slug,
		Codigo:                p.Codigo,
		Nombre:                p.Nombre,
		Descripcion:           p.Descripcion,
		Categoria:             categoriaRef{ID: p.CategoriaID, Slug: p.CategoriaSlug, Nombre: p.CategoriaNombre},
		Precio:                p.PrecioDetalle,
		PrecioFinal:           productos.PrecioFinal(p, hoy),
		PrecioMayor:           p.PrecioMayor,
		MayorMinimo:           p.MayorMinimo,
		UnidadMayor:           ptr(p.UnidadMayor),
		MayorMinimoConfirmado: p.MayorMinimoConfirmado,
		OfertaHasta:           ptr(p.OfertaHasta),
		OfertaVigente:         vigente,
		PrecioActualizadoEn:   utcPtr(p.PrecioActualizadoEn),
		Disponibilidad:        productos.CalcularDisponibilidad(p.Stock, p.UmbralPocos),
		Badges:                productos.Badges(p, s.hoy(), s.cfg.DiasRecienLlegado),
		Temporada:             ptr(p.Temporada),
		FechaLlegada:          ptr(p.FechaLlegada),
		ImagenURL:             ptr(p.ImagenURL),
		TikTokURL:             ptr(p.TikTokURL),
		EsEjemplo:             p.EsEjemplo,
	}
	if vigente {
		d.PrecioOferta = p.PrecioOferta
		d.EtiquetaOferta = ptr(p.EtiquetaOferta)
	}
	if p.FuenteURL != "" {
		d.Fuente = &fuenteDTO{URL: p.FuenteURL, Fecha: ptr(p.FuenteFecha)}
	}
	return d
}

// productoAdmin muestra precio_oferta guardado aunque haya vencido (para editarlo);
// oferta_vigente dice si hoy se aplica.
func (s *Server) productoAdmin(p store.Producto) productoAdminDTO {
	pub := s.productoPublico(p)
	pub.PrecioOferta = p.PrecioOferta
	if p.PrecioOferta != nil {
		pub.EtiquetaOferta = ptr(p.EtiquetaOferta)
	}
	return productoAdminDTO{
		productoPublicoDTO: pub,
		Stock:              p.Stock,
		UmbralPocos:        p.UmbralPocos,
		Activo:             p.Activo,
		CreadoEn:           p.CreadoEn.UTC(),
		ActualizadoEn:      p.ActualizadoEn.UTC(),
	}
}

// ---------------------------------------------------------------- pedidos (v2)

type itemDTO struct {
	ProductoID     *int64  `json:"producto_id"`
	Codigo         *string `json:"codigo"`
	Descripcion    string  `json:"descripcion"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario *int64  `json:"precio_unitario"`
	TipoPrecio     string  `json:"tipo_precio"`
	Subtotal       *int64  `json:"subtotal"`
}

type usuarioRefDTO struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

func refDTO(u *store.UsuarioRef) *usuarioRefDTO {
	if u == nil {
		return nil
	}
	return &usuarioRefDTO{ID: u.ID, Nombre: u.Nombre}
}

type pagoDTO struct {
	Medio         string         `json:"medio"`
	Monto         int64          `json:"monto"`
	Fecha         string         `json:"fecha"`
	VerificadoPor *usuarioRefDTO `json:"verificado_por"`
}

type cambioDTO struct {
	EstadoAnterior *string        `json:"estado_anterior"`
	EstadoNuevo    string         `json:"estado_nuevo"`
	Usuario        *usuarioRefDTO `json:"usuario"`
	Nota           string         `json:"nota"`
	CreadoEn       time.Time      `json:"creado_en"`
}

type pedidoDTO struct {
	ID                     int64          `json:"id"`
	Numero                 string         `json:"numero"`
	Canal                  string         `json:"canal"`
	Estado                 string         `json:"estado"`
	EstadoDesde            time.Time      `json:"estado_desde"`
	EsDemo                 bool           `json:"es_demo"`
	Nombre                 string         `json:"nombre"`
	Telefono               string         `json:"telefono"`
	Entrega                string         `json:"entrega"`
	Region                 *string        `json:"region"`
	RegionNombre           *string        `json:"region_nombre"`
	Comuna                 *string        `json:"comuna"`
	Agencia                *string        `json:"agencia"`
	AgenciaNombre          *string        `json:"agencia_nombre"`
	NSeguimiento           *string        `json:"n_seguimiento"`
	Items                  []itemDTO      `json:"items"`
	FormaPago              string         `json:"forma_pago"`
	Subtotal               int64          `json:"subtotal"`
	Descuento              int64          `json:"descuento"`
	CostoEnvio             int64          `json:"costo_envio"`
	EnvioPorPagar          bool           `json:"envio_por_pagar"`
	Total                  int64          `json:"total"`
	TotalManual            bool           `json:"total_manual"`
	MotivoCancelacion      string         `json:"motivo_cancelacion"`
	NotaInterna            string         `json:"nota_interna"`
	Pago                   *pagoDTO       `json:"pago"`
	StockDescontado        bool           `json:"stock_descontado"`
	RegistradoPor          *usuarioRefDTO `json:"registrado_por"`
	SeguimientoURL         string         `json:"seguimiento_url"`
	TransicionesPermitidas []string       `json:"transiciones_permitidas"`
	WhatsAppClienteURL     string         `json:"whatsapp_cliente_url"`
	CreadoEn               time.Time      `json:"creado_en"`
	ActualizadoEn          time.Time      `json:"actualizado_en"`
	Historial              []cambioDTO    `json:"historial,omitempty"`
}

func (s *Server) pedidoDTO(p store.Pedido, conHistorial bool) pedidoDTO {
	d := pedidoDTO{
		ID: p.ID, Numero: p.Numero, Canal: p.Canal, Estado: p.Estado, EstadoDesde: p.EstadoDesde.UTC(), EsDemo: p.EsDemo,
		Nombre: p.Nombre, Telefono: p.Telefono, Entrega: p.Entrega,
		Region: ptr(p.Region), Comuna: ptr(p.Comuna), Agencia: ptr(p.Agencia), NSeguimiento: ptr(p.NSeguimiento),
		Items:     make([]itemDTO, 0, len(p.Items)),
		FormaPago: p.FormaPago, Subtotal: p.Subtotal, Descuento: p.Descuento, CostoEnvio: p.CostoEnvio,
		EnvioPorPagar: p.EnvioPorPagar, Total: p.Total, TotalManual: p.TotalManual, MotivoCancelacion: p.MotivoCancelacion,
		NotaInterna: p.NotaInterna, StockDescontado: p.StockDescontado, RegistradoPor: refDTO(p.RegistradoPor),
		SeguimientoURL:         mensajes.SeguimientoURL(s.cfg.Tienda.URLSitio, p.TokenSeguimiento),
		TransicionesPermitidas: pedidos.TransicionesDesde(p.Entrega, p.FormaPago, p.Estado),
		WhatsAppClienteURL:     waCliente(p),
		CreadoEn:               p.CreadoEn.UTC(), ActualizadoEn: p.ActualizadoEn.UTC(),
	}
	if r, ok := pedidos.RegionPorCodigo(p.Region); ok {
		d.RegionNombre = &r.Nombre
	}
	if a, ok := pedidos.AgenciaPorCodigo(p.Agencia); ok {
		d.AgenciaNombre = &a.Nombre
	}
	for _, it := range p.Items {
		i := itemDTO{Descripcion: it.Descripcion, Cantidad: it.Cantidad, PrecioUnitario: it.PrecioUnitario, TipoPrecio: it.TipoPrecio, Subtotal: it.Subtotal}
		if it.ProductoID > 0 {
			id := it.ProductoID
			i.ProductoID, i.Codigo = &id, ptr(it.Codigo)
		}
		d.Items = append(d.Items, i)
	}
	if p.Pago != nil {
		d.Pago = &pagoDTO{Medio: p.Pago.Medio, Monto: p.Pago.Monto, Fecha: p.Pago.Fecha, VerificadoPor: refDTO(p.Pago.VerificadoPor)}
	}
	if conHistorial {
		d.Historial = make([]cambioDTO, 0, len(p.Historial))
		for _, c := range p.Historial {
			d.Historial = append(d.Historial, cambioDTO{EstadoAnterior: ptr(c.EstadoAnterior), EstadoNuevo: c.EstadoNuevo,
				Usuario: refDTO(c.Usuario), Nota: c.Nota, CreadoEn: c.CreadoEn.UTC()})
		}
	}
	return d
}

// waCliente: chat con el cliente; en pedidos DEMO, sin número (regla de la plantilla).
func waCliente(p store.Pedido) string {
	if p.EsDemo {
		return "https://wa.me/"
	}
	return "https://wa.me/" + strings.TrimPrefix(p.Telefono, "+")
}

func (s *Server) pedidosDTO(ps []store.Pedido) []pedidoDTO {
	out := make([]pedidoDTO, 0, len(ps))
	for _, p := range ps {
		out = append(out, s.pedidoDTO(p, false))
	}
	return out
}
