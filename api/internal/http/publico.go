package httpapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/limite"
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
	maxBodyPedido      = 16 << 10
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
			"direccion": s.cfg.Tienda.DireccionLocal,
		},
		"regiones":            pedidos.Regiones,
		"agencias":            pedidos.Agencias,
		"agencias_es_ejemplo": pedidos.AgenciasSonEjemplo,
	})
}

func (s *Server) crearPedido(w http.ResponseWriter, r *http.Request) {
	// Primero se valida: un cuerpo inválido no consume cupo de nadie.
	var sol pedidos.Solicitud
	if err := leerJSON(w, r, &sol, maxBodyPedido); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	nuevo, campos := sol.Validar()
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del pedido.", campos)
		return
	}
	ahora := s.ahora()
	nuevo.Ahora = ahora
	nuevo.EsDemo = s.cfg.PedidosDemo
	nuevo.MaxAbiertos = s.cfg.MaxAbiertosPorTel

	// Límites: por IP, por teléfono y global por hora.
	for _, lim := range []struct {
		v     *limite.Ventana
		clave string
	}{
		{s.limPedidosIP, "ip:" + claveLimite(s.ipCliente(r))},
		{s.limPedidosTel, "tel:" + nuevo.Telefono},
		{s.limPedidosGlobal, "global"},
	} {
		if ok, espera := lim.v.Permitir(lim.clave, ahora); !ok {
			if lim.v == s.limPedidosGlobal {
				s.log.Warn("tope global de pedidos por hora alcanzado")
			}
			w.Header().Set("Retry-After", strconv.Itoa(segundos(espera)))
			responderError(w, http.StatusTooManyRequests, "demasiadas_solicitudes", "Hay muchos pedidos seguidos. Espera unos minutos o escríbenos por WhatsApp.", nil)
			return
		}
	}

	p, err := s.st.CrearPedido(r.Context(), nuevo, pedidos.Cotizador(s.hoyISO()))
	var nd *pedidos.ErrNoDisponibles
	switch {
	case errors.As(err, &nd):
		responderErrorDetalles(w, http.StatusConflict, "productos_no_disponibles",
			"Algunos productos ya no están disponibles en esa cantidad.", nd.Detalles)
		return
	case errors.Is(err, store.ErrDemasiadosAbiertos):
		responderError(w, http.StatusConflict, "pedidos_abiertos",
			"Ya tienes pedidos esperando confirmación. Escríbenos por WhatsApp para confirmarlos antes de hacer otro.", nil)
		return
	case err != nil:
		s.errorInterno(w, r, err)
		return
	}

	if sol.Total != nil && *sol.Total != p.Total {
		s.log.Info("total del navegador distinto al recalculado; se usa el del servidor",
			"numero", p.Numero, "total_cliente", *sol.Total, "total", p.Total)
	}
	s.log.Info("pedido creado",
		"numero", p.Numero, "total", p.Total, "items", len(p.Items),
		"entrega", p.Entrega, "telefono", pedidos.EnmascararTelefono(p.Telefono))

	msg := pedidos.MensajeWhatsApp(p, s.cfg.Tienda)
	responderJSON(w, http.StatusCreated, pedidoCreadoDTO{
		Numero:          p.Numero,
		Estado:          p.Estado,
		Nombre:          p.Nombre,
		entregaDTO:      entrega(p),
		Comentario:      p.Comentario,
		Total:           p.Total,
		EsDemo:          p.EsDemo,
		Items:           itemsDTO(p.Items),
		MensajeWhatsApp: msg,
		WhatsAppURL:     pedidos.URLWhatsApp(s.cfg.Tienda.WhatsApp, msg),
		CreadoEn:        p.CreadoEn,
	})
}

func segundos(d time.Duration) int {
	n := int(d.Seconds())
	if n < 1 {
		n = 1
	}
	return n
}
