package httpapi

import (
	"time"

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
	ID             int64        `json:"id"`
	Slug           string       `json:"slug"`
	Codigo         string       `json:"codigo"`
	Nombre         string       `json:"nombre"`
	Descripcion    string       `json:"descripcion"`
	Categoria      categoriaRef `json:"categoria"`
	Precio         int64        `json:"precio"`
	PrecioOferta   *int64       `json:"precio_oferta"`
	PrecioFinal    int64        `json:"precio_final"`
	EtiquetaOferta *string      `json:"etiqueta_oferta"`
	PrecioMayor    *int64       `json:"precio_mayor"`
	MayorMinimo    *int         `json:"mayor_minimo"`
	UnidadMayor    *string      `json:"unidad_mayor"`
	Disponibilidad string       `json:"disponibilidad"`
	Badges         []string     `json:"badges"`
	Temporada      *string      `json:"temporada"`
	FechaLlegada   *string      `json:"fecha_llegada"`
	ImagenURL      *string      `json:"imagen_url"`
	TikTokURL      *string      `json:"tiktok_url"`
	Fuente         *fuenteDTO   `json:"fuente"`
	EsEjemplo      bool         `json:"es_ejemplo"`
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

func (s *Server) productoPublico(p store.Producto) productoPublicoDTO {
	d := productoPublicoDTO{
		ID:             p.ID,
		Slug:           p.Slug,
		Codigo:         p.Codigo,
		Nombre:         p.Nombre,
		Descripcion:    p.Descripcion,
		Categoria:      categoriaRef{ID: p.CategoriaID, Slug: p.CategoriaSlug, Nombre: p.CategoriaNombre},
		Precio:         p.PrecioDetalle,
		PrecioOferta:   p.PrecioOferta,
		PrecioFinal:    productos.PrecioFinal(p),
		PrecioMayor:    p.PrecioMayor,
		MayorMinimo:    p.MayorMinimo,
		UnidadMayor:    ptr(p.UnidadMayor),
		Disponibilidad: productos.CalcularDisponibilidad(p.Stock, p.UmbralPocos),
		Badges:         productos.Badges(p, s.hoy(), s.cfg.DiasRecienLlegado),
		Temporada:      ptr(p.Temporada),
		FechaLlegada:   ptr(p.FechaLlegada),
		ImagenURL:      ptr(p.ImagenURL),
		TikTokURL:      ptr(p.TikTokURL),
		EsEjemplo:      p.EsEjemplo,
	}
	if p.PrecioOferta != nil {
		d.EtiquetaOferta = ptr(p.EtiquetaOferta)
	}
	if p.FuenteURL != "" {
		d.Fuente = &fuenteDTO{URL: p.FuenteURL, Fecha: ptr(p.FuenteFecha)}
	}
	return d
}

func (s *Server) productoAdmin(p store.Producto) productoAdminDTO {
	return productoAdminDTO{
		productoPublicoDTO: s.productoPublico(p),
		Stock:              p.Stock,
		UmbralPocos:        p.UmbralPocos,
		Activo:             p.Activo,
		CreadoEn:           p.CreadoEn,
		ActualizadoEn:      p.ActualizadoEn,
	}
}

type itemDTO struct {
	ProductoID     int64  `json:"producto_id"`
	Codigo         string `json:"codigo"`
	Nombre         string `json:"nombre"`
	Cantidad       int    `json:"cantidad"`
	PrecioUnitario int64  `json:"precio_unitario"`
	TipoPrecio     string `json:"tipo_precio"`
	Subtotal       int64  `json:"subtotal"`
}

// entregaDTO va embebido (campos planos), como pide propuesta/mensaje-pedido-whatsapp.md.
type entregaDTO struct {
	Entrega       string  `json:"entrega"`
	Region        *string `json:"region"`
	RegionNombre  *string `json:"region_nombre"`
	Comuna        *string `json:"comuna"`
	Agencia       *string `json:"agencia"`
	AgenciaNombre *string `json:"agencia_nombre"`
}

func itemsDTO(items []store.ItemPedido) []itemDTO {
	out := make([]itemDTO, 0, len(items))
	for _, it := range items {
		out = append(out, itemDTO(it))
	}
	return out
}

func entrega(p store.Pedido) entregaDTO {
	e := entregaDTO{Entrega: p.Entrega, Region: ptr(p.Region), Comuna: ptr(p.Comuna), Agencia: ptr(p.Agencia)}
	if r, ok := pedidos.RegionPorCodigo(p.Region); ok {
		e.RegionNombre = &r.Nombre
	}
	if a, ok := pedidos.AgenciaPorCodigo(p.Agencia); ok {
		e.AgenciaNombre = &a.Nombre
	}
	return e
}

type pedidoCreadoDTO struct {
	Numero string `json:"numero"`
	Estado string `json:"estado"`
	Nombre string `json:"nombre"`
	entregaDTO
	Comentario      string    `json:"comentario"`
	Total           int64     `json:"total"`
	EsDemo          bool      `json:"es_demo"`
	Items           []itemDTO `json:"items"`
	MensajeWhatsApp string    `json:"mensaje_whatsapp"`
	WhatsAppURL     string    `json:"whatsapp_url"`
	CreadoEn        time.Time `json:"creado_en"`
}

type cambioDTO struct {
	EstadoAnterior *string   `json:"estado_anterior"`
	EstadoNuevo    string    `json:"estado_nuevo"`
	Admin          *string   `json:"admin"`
	Nota           string    `json:"nota"`
	CreadoEn       time.Time `json:"creado_en"`
}

type pedidoAdminDTO struct {
	ID       int64  `json:"id"`
	Numero   string `json:"numero"`
	Estado   string `json:"estado"`
	EsDemo   bool   `json:"es_demo"`
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
	entregaDTO
	Comentario             string      `json:"comentario"`
	Total                  int64       `json:"total"`
	Items                  []itemDTO   `json:"items"`
	TransicionesPermitidas []string    `json:"transiciones_permitidas"`
	WhatsAppClienteURL     string      `json:"whatsapp_cliente_url"`
	CreadoEn               time.Time   `json:"creado_en"`
	ActualizadoEn          time.Time   `json:"actualizado_en"`
	Historial              []cambioDTO `json:"historial,omitempty"`
}

func pedidoAdmin(p store.Pedido, conHistorial bool) pedidoAdminDTO {
	d := pedidoAdminDTO{
		ID:                     p.ID,
		Numero:                 p.Numero,
		Estado:                 p.Estado,
		EsDemo:                 p.EsDemo,
		Nombre:                 p.Nombre,
		Telefono:               p.Telefono,
		entregaDTO:             entrega(p),
		Comentario:             p.Comentario,
		Total:                  p.Total,
		Items:                  itemsDTO(p.Items),
		TransicionesPermitidas: pedidos.TransicionesDesde(p.Entrega, p.Estado),
		WhatsAppClienteURL:     "https://wa.me/" + trimMas(p.Telefono),
		CreadoEn:               p.CreadoEn,
		ActualizadoEn:          p.ActualizadoEn,
	}
	if conHistorial {
		d.Historial = make([]cambioDTO, 0, len(p.Historial))
		for _, c := range p.Historial {
			d.Historial = append(d.Historial, cambioDTO{
				EstadoAnterior: ptr(c.EstadoAnterior),
				EstadoNuevo:    c.EstadoNuevo,
				Admin:          ptr(c.AdminUsuario),
				Nota:           c.Nota,
				CreadoEn:       c.CreadoEn,
			})
		}
	}
	return d
}

func trimMas(t string) string {
	if len(t) > 0 && t[0] == '+' {
		return t[1:]
	}
	return t
}
