package httpapi

import (
	"errors"
	"net/http"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

const maxBodyProducto = 32 << 10

func (s *Server) adminListarProductos(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	f, page, size, campos := s.filtroProductos(r.URL.Query(), pageSizeAdminMax)
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los filtros.", campos)
		return
	}
	f.IncluirInactivos = true
	lista, total, err := s.st.ListarProductos(r.Context(), f)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	out := make([]productoAdminDTO, 0, len(lista))
	for _, p := range lista {
		out = append(out, s.productoAdmin(p))
	}
	responderJSON(w, http.StatusOK, nuevaPagina(out, page, size, total))
}

func (s *Server) adminObtenerProducto(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	p, err := s.st.ProductoPorID(r.Context(), id)
	if errors.Is(err, store.ErrNoEncontrado) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, s.productoAdmin(p))
}

// validarProducto lee y valida el cuerpo, incluida la existencia de la categoría.
func (s *Server) validarProducto(w http.ResponseWriter, r *http.Request) (store.Producto, bool) {
	var in productos.Entrada
	if err := leerJSON(w, r, &in, maxBodyProducto); err != nil {
		responderErrorDecodificar(w, err)
		return store.Producto{}, false
	}
	p, campos := in.Validar()
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del producto.", campos)
		return store.Producto{}, false
	}
	cats, err := s.st.ListarCategorias(r.Context())
	if err != nil {
		s.errorInterno(w, r, err)
		return store.Producto{}, false
	}
	for _, c := range cats {
		if c.ID == p.CategoriaID {
			p.CategoriaNombre = c.Nombre
			p.Busqueda = productos.TextoBusqueda(p)
			return p, true
		}
	}
	responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del producto.", map[string]string{"categoria_id": "la categoría no existe"})
	return store.Producto{}, false
}

func (s *Server) responderErrorGuardar(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, store.ErrDuplicado):
		responderError(w, http.StatusConflict, "duplicado", "Ya existe un producto con ese slug o código.", nil)
	case errors.Is(err, store.ErrReferenciaInvalida):
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del producto.", map[string]string{"categoria_id": "la categoría no existe"})
	case errors.Is(err, store.ErrNoEncontrado):
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
	default:
		s.errorInterno(w, r, err)
	}
}

func (s *Server) adminCrearProducto(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	p, ok := s.validarProducto(w, r)
	if !ok {
		return
	}
	ahora := s.ahora()
	p.PrecioActualizadoEn = &ahora // el precio lo fijó el admin
	if err := s.st.CrearProducto(r.Context(), &p, ahora); err != nil {
		s.responderErrorGuardar(w, r, err)
		return
	}
	creado, err := s.st.ProductoPorID(r.Context(), p.ID)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("producto creado", "id", creado.ID, "slug", creado.Slug, "usuario", ses.Usuario)
	responderJSON(w, http.StatusCreated, s.productoAdmin(creado))
}

func (s *Server) adminActualizarProducto(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	p, ok := s.validarProducto(w, r)
	if !ok {
		return
	}
	p.ID = id
	actual, err := s.st.ProductoPorID(r.Context(), id)
	if err != nil {
		s.responderErrorGuardar(w, r, err)
		return
	}
	ahora := s.ahora()
	p.PrecioActualizadoEn = actual.PrecioActualizadoEn
	if preciosCambiaron(actual, p) {
		p.PrecioActualizadoEn = &ahora
	}
	if err := s.st.ActualizarProducto(r.Context(), &p, ahora); err != nil {
		s.responderErrorGuardar(w, r, err)
		return
	}
	act, err := s.st.ProductoPorID(r.Context(), id)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("producto actualizado", "id", id, "usuario", ses.Usuario)
	responderJSON(w, http.StatusOK, s.productoAdmin(act))
}

func igualPtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// preciosCambiaron: precio_detalle, precio_oferta o precio_mayor distintos.
func preciosCambiaron(antes, despues store.Producto) bool {
	return antes.PrecioDetalle != despues.PrecioDetalle ||
		!igualPtr(antes.PrecioOferta, despues.PrecioOferta) ||
		!igualPtr(antes.PrecioMayor, despues.PrecioMayor)
}

func (s *Server) adminEliminarProducto(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	if err := s.st.EliminarProducto(r.Context(), id, s.ahora()); err != nil {
		s.responderErrorGuardar(w, r, err)
		return
	}
	s.log.Info("producto eliminado", "id", id, "usuario", ses.Usuario)
	w.WriteHeader(http.StatusNoContent)
}
