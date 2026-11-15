package httpapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

const (
	pageSizePublicoMax = 48
	pageSizePublicoDef = 24
	maxLargoBusqueda   = 60
	maxPalabras        = 5
)

func (s *Server) listarCategorias(w http.ResponseWriter, r *http.Request) {
	cats, err := s.st.ListarCategorias(r.Context())
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	out := make([]categoriaDTO, 0, len(cats))
	for _, c := range cats {
		out = append(out, categoriaDTO{ID: c.ID, Slug: c.Slug, Nombre: c.Nombre, TotalProductos: c.TotalProductos})
	}
	responderJSON(w, http.StatusOK, map[string]any{"items": out})
}

// paginacion lee page y page_size. Valores no numéricos o fuera de rango -> error.
func paginacion(q url.Values, def, max int) (page, size int, campos map[string]string) {
	campos = map[string]string{}
	page, size = 1, def
	if v := q.Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 10000 {
			campos["page"] = "debe ser un entero desde 1"
		} else {
			page = n
		}
	}
	if v := q.Get("page_size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > max {
			campos["page_size"] = "debe ser un entero entre 1 y " + strconv.Itoa(max)
		} else {
			size = n
		}
	}
	return page, size, campos
}

func leerBool(q url.Values, k string, campos map[string]string) bool {
	switch q.Get(k) {
	case "":
		return false
	case "true", "1":
		return true
	case "false", "0":
		return false
	}
	campos[k] = "debe ser true o false"
	return false
}

// filtroProductos interpreta los query params comunes a público y admin.
func (s *Server) filtroProductos(q url.Values, max int) (store.FiltroProductos, int, int, map[string]string) {
	page, size, campos := paginacion(q, pageSizePublicoDef, max)
	f := store.FiltroProductos{Orden: "relevancia", Hoy: s.hoyISO()}

	if v := q.Get("categoria"); v != "" {
		if len(v) > 80 || productos.Slugify(v) != v {
			campos["categoria"] = "slug de categoría inválido"
		}
		f.CategoriaSlug = v
	}
	if v := q.Get("q"); v != "" {
		limpio, ok := texto.Limpiar(v, false)
		if !ok || texto.Largo(limpio) > maxLargoBusqueda {
			campos["q"] = "búsqueda de máximo 60 caracteres"
		} else {
			f.Palabras = productos.PalabrasBusqueda(limpio, maxPalabras)
		}
	}
	f.SoloOferta = leerBool(q, "oferta", campos)
	if leerBool(q, "recien_llegado", campos) {
		f.LlegadaDesde = productos.LlegadaDesde(s.hoy(), s.cfg.DiasRecienLlegado)
	}
	if v := q.Get("temporada"); v != "" {
		if !productos.EsTemporada(v) {
			campos["temporada"] = "temporada desconocida"
		}
		f.Temporada = v
	}
	if v := q.Get("orden"); v != "" {
		switch v {
		case "relevancia", "precio_asc", "precio_desc", "recientes":
			f.Orden = v
		default:
			campos["orden"] = "debe ser relevancia, precio_asc, precio_desc o recientes"
		}
	}
	f.Limite = size
	f.Desplazamiento = (page - 1) * size
	return f, page, size, campos
}

func (s *Server) listarProductos(w http.ResponseWriter, r *http.Request) {
	f, page, size, campos := s.filtroProductos(r.URL.Query(), pageSizePublicoMax)
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los filtros.", campos)
		return
	}
	lista, total, err := s.st.ListarProductos(r.Context(), f)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	out := make([]productoPublicoDTO, 0, len(lista))
	for _, p := range lista {
		out = append(out, s.productoPublico(p))
	}
	responderJSON(w, http.StatusOK, nuevaPagina(out, page, size, total))
}

func (s *Server) obtenerProducto(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if len(slug) > 80 || productos.Slugify(slug) != slug {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	p, err := s.st.ProductoPorSlug(r.Context(), slug, false)
	if errors.Is(err, store.ErrNoEncontrado) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, s.productoPublico(p))
}

func (s *Server) opcionesEntrega(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, map[string]any{
		"retiro": map[string]string{
			"direccion": s.cfg.Tienda.DireccionRetiro,
		},
		"regiones":            pedidos.Regiones,
		"agencias":            pedidos.Agencias,
		"agencias_es_ejemplo": pedidos.AgenciasSonEjemplo,
	})
}

func segundos(d time.Duration) int {
	n := int(d.Seconds())
	if n < 1 {
		n = 1
	}
	return n
}

type configPublicaDTO struct {
	NombreTienda               string  `json:"nombre_tienda"`
	WhatsAppSantiago           string  `json:"whatsapp_santiago"`
	WhatsAppRegiones           *string `json:"whatsapp_regiones"`
	WhatsAppRegionesConfirmado bool    `json:"whatsapp_regiones_confirmado"`
	Direccion                  string  `json:"direccion"`
	Horario                    string  `json:"horario"`
	HorarioConfirmado          bool    `json:"horario_confirmado"`
	SiteURL                    string  `json:"site_url"`
}

// configPublica: solo datos públicos para la landing.
func (s *Server) configPublica(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, configPublicaDTO{
		NombreTienda:               s.cfg.Tienda.NombreTienda,
		WhatsAppSantiago:           s.cfg.WhatsAppSantiago,
		WhatsAppRegiones:           ptr(s.cfg.WhatsAppRegiones),
		WhatsAppRegionesConfirmado: s.cfg.WhatsAppRegiones != "",
		Direccion:                  s.cfg.Tienda.DireccionRetiro,
		Horario:                    s.cfg.Tienda.HorarioRetiro,
		HorarioConfirmado:          s.cfg.HorarioConfirmado,
		SiteURL:                    s.cfg.Tienda.URLSitio,
	})
}
