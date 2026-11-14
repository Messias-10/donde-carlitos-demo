package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

const (
	maxBodyLogin         = 4 << 10
	maxBodyProducto      = 32 << 10
	maxBodyEstado        = 4 << 10
	pageSizeAdminMax     = 100
	pageSizeAdminPedidos = 20
)

// nombreCookie: en producción usa el prefijo __Host- (exige Secure, Path=/ y sin Domain).
func (s *Server) nombreCookie() string {
	if s.cfg.Produccion {
		return "__Host-dc_sesion"
	}
	return "dc_sesion"
}

func (s *Server) adminHabilitado() bool { return len(s.cfg.SecretoSesion) >= 32 }

func (s *Server) cookieSesion(valor string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     s.nombreCookie(),
		Value:    valor,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   s.cfg.Produccion,
		SameSite: http.SameSiteStrictMode,
	}
}

// hashSeñuelo se compara cuando el usuario no existe, para que el tiempo de
// respuesta no revele si el usuario es válido.
func (s *Server) hashSenuelo() string {
	s.dummyOnce.Do(func() {
		clave, err := auth.GenerarClave(24)
		if err == nil {
			s.dummyHash, _ = auth.HashClave(clave, s.cfg.CostoBcryptDummy)
		}
	})
	return s.dummyHash
}

type entradaLogin struct {
	Usuario  string `json:"usuario"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.adminHabilitado() {
		responderError(w, http.StatusServiceUnavailable, "admin_no_configurado", "El panel aún no está configurado.", nil)
		return
	}
	ip := claveLimite(s.ipCliente(r))
	ahora := s.ahora()
	if bloq, espera := s.bloqueo.Bloqueado(ip, ahora); bloq {
		w.Header().Set("Retry-After", strconv.Itoa(segundos(espera)))
		responderError(w, http.StatusTooManyRequests, "login_bloqueado", "Demasiados intentos. Espera unos minutos.", nil)
		return
	}

	var in entradaLogin
	if err := leerJSON(w, r, &in, maxBodyLogin); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	// Tope global por minuto y como mucho N bcrypt a la vez: un flood de logins
	// (aunque venga de muchas IPs) no puede acaparar la CPU.
	if ok, espera := s.limLoginGlobal.Permitir("global", ahora); !ok {
		w.Header().Set("Retry-After", strconv.Itoa(segundos(espera)))
		responderError(w, http.StatusTooManyRequests, "demasiadas_solicitudes", "Demasiados intentos. Espera un minuto.", nil)
		return
	}
	if !s.semBcrypt.Intentar() {
		w.Header().Set("Retry-After", "1")
		responderError(w, http.StatusTooManyRequests, "demasiadas_solicitudes", "Hay mucha actividad. Intenta de nuevo en un segundo.", nil)
		return
	}
	defer s.semBcrypt.Soltar()

	hash := s.hashSenuelo()
	var admin store.Admin
	var err error
	valido := len(in.Usuario) >= 1 && len(in.Usuario) <= 64 && len(in.Password) >= 1 && len(in.Password) <= 72
	if valido {
		admin, err = s.st.AdminPorUsuario(r.Context(), in.Usuario)
		switch {
		case err == nil:
			hash = admin.PasswordHash
		case errors.Is(err, store.ErrNoEncontrado):
			valido = false
		default:
			s.errorInterno(w, r, err)
			return
		}
	}
	// Siempre se corre bcrypt, exista o no el usuario.
	okClave := auth.VerificarClave(hash, in.Password)
	if !valido || !okClave {
		if s.bloqueo.Fallo(ip, ahora) {
			s.log.Warn("login bloqueado por intentos fallidos", "ip", enmascararIP(s.ipCliente(r)))
		} else {
			s.log.Warn("login fallido", "ip", enmascararIP(s.ipCliente(r)))
		}
		responderError(w, http.StatusUnauthorized, "credenciales_invalidas", "Usuario o clave incorrectos.", nil)
		return
	}
	s.bloqueo.Exito(ip)

	token, err := auth.NuevoToken()
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	ses := store.Sesion{
		TokenHash: auth.HashToken(s.cfg.SecretoSesion, token),
		AdminID:   admin.ID,
		CreadoEn:  ahora,
		ExpiraEn:  ahora.Add(s.cfg.DuracionSesion),
		UltimoUso: ahora,
	}
	if err := s.st.CrearSesion(r.Context(), ses); err != nil {
		s.errorInterno(w, r, err)
		return
	}
	http.SetCookie(w, s.cookieSesion(token, int(s.cfg.DuracionSesion.Seconds())))
	s.log.Info("login correcto", "usuario", admin.Usuario)
	responderJSON(w, http.StatusOK, map[string]any{"usuario": admin.Usuario, "expira_en": ses.ExpiraEn})
}

// logout es idempotente: siempre responde 204 y borra la cookie.
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(s.nombreCookie()); err == nil && s.adminHabilitado() && auth.TokenConForma(c.Value) {
		if err := s.st.EliminarSesion(r.Context(), auth.HashToken(s.cfg.SecretoSesion, c.Value)); err != nil {
			s.log.Error("logout: no se pudo borrar la sesión", "err", err)
		}
	}
	http.SetCookie(w, s.cookieSesion("", -1))
	w.WriteHeader(http.StatusNoContent)
}

type handlerConSesion func(w http.ResponseWriter, r *http.Request, ses store.Sesion)

// conSesion exige una sesión válida; si no, 401.
func (s *Server) conSesion(h handlerConSesion) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAut := func() {
			responderError(w, http.StatusUnauthorized, "no_autenticado", "Inicia sesión para continuar.", nil)
		}
		if !s.adminHabilitado() {
			noAut()
			return
		}
		c, err := r.Cookie(s.nombreCookie())
		if err != nil || !auth.TokenConForma(c.Value) {
			noAut()
			return
		}
		hash := auth.HashToken(s.cfg.SecretoSesion, c.Value)
		ses, err := s.st.SesionPorHash(r.Context(), hash)
		if errors.Is(err, store.ErrNoEncontrado) {
			noAut()
			return
		}
		if err != nil {
			s.errorInterno(w, r, err)
			return
		}
		ahora := s.ahora()
		if !ahora.Before(ses.ExpiraEn) || ahora.Sub(ses.UltimoUso) > s.cfg.InactividadSesion {
			_ = s.st.EliminarSesion(r.Context(), hash)
			http.SetCookie(w, s.cookieSesion("", -1))
			noAut()
			return
		}
		if ahora.Sub(ses.UltimoUso) > time.Minute {
			if err := s.st.TocarSesion(r.Context(), hash, ahora); err != nil {
				s.log.Warn("no se pudo actualizar la sesión", "err", err)
			}
		}
		h(w, conValor(r, ctxSesion, ses), ses)
	})
}

func (s *Server) sesionActual(w http.ResponseWriter, _ *http.Request, ses store.Sesion) {
	responderJSON(w, http.StatusOK, map[string]any{"usuario": ses.Usuario, "expira_en": ses.ExpiraEn})
}

func idDeRuta(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}

// ---------------------------------------------------------------- pedidos

func (s *Server) adminListarPedidos(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	q := r.URL.Query()
	page, size, campos := paginacion(q, pageSizeAdminPedidos, pageSizeAdminMax)
	estado := q.Get("estado")
	if estado != "" && !pedidos.EsEstado(estado) {
		campos["estado"] = "estado desconocido"
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los filtros.", campos)
		return
	}
	lista, total, err := s.st.ListarPedidos(r.Context(), store.FiltroPedidos{Estado: estado, Limite: size, Desplazamiento: (page - 1) * size})
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	out := make([]pedidoAdminDTO, 0, len(lista))
	for _, p := range lista {
		out = append(out, pedidoAdmin(p, false))
	}
	responderJSON(w, http.StatusOK, nuevaPagina(out, page, size, total))
}

func (s *Server) adminObtenerPedido(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	p, err := s.st.PedidoPorID(r.Context(), id)
	if errors.Is(err, store.ErrNoEncontrado) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, pedidoAdmin(p, true))
}

type entradaEstado struct {
	Estado string `json:"estado"`
	Nota   string `json:"nota"`
}

func (s *Server) adminCambiarEstado(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	var in entradaEstado
	if err := leerJSON(w, r, &in, maxBodyEstado); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	campos := map[string]string{}
	if !pedidos.EsEstado(in.Estado) {
		campos["estado"] = "estado desconocido"
	}
	nota, okNota := texto.Limpiar(in.Nota, false)
	if !okNota || texto.Largo(nota) > 200 {
		campos["nota"] = "máximo 200 caracteres"
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos.", campos)
		return
	}

	p, err := s.st.CambiarEstado(r.Context(), store.CambioEstado{
		PedidoID: id, Nuevo: in.Estado, AdminID: ses.AdminID, Nota: nota, Ahora: s.ahora(),
	}, pedidos.Validador(in.Estado))
	var ti *pedidos.ErrTransicionInvalida
	var sinStock *store.ErrSinStock
	switch {
	case errors.As(err, &sinStock):
		det := make([]pedidos.Indisponible, 0, len(sinStock.ProductoIDs))
		for _, id := range sinStock.ProductoIDs {
			det = append(det, pedidos.Indisponible{ProductoID: id, Motivo: pedidos.MotivoStockInsufic})
		}
		responderErrorDetalles(w, http.StatusConflict, "stock_insuficiente",
			"No hay stock suficiente para confirmar este pedido. Ajusta el stock o cancela el pedido.", det)
		return
	case errors.Is(err, store.ErrNoEncontrado):
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	case errors.As(err, &ti):
		responderError(w, http.StatusConflict, "transicion_invalida", "No se puede pasar de \""+ti.De+"\" a \""+ti.A+"\".", nil)
		return
	case errors.Is(err, store.ErrConflicto):
		responderError(w, http.StatusConflict, "conflicto", "El pedido cambió mientras lo editabas. Recarga e intenta de nuevo.", nil)
		return
	case err != nil:
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("estado de pedido cambiado", "numero", p.Numero, "estado", p.Estado, "usuario", ses.Usuario)
	responderJSON(w, http.StatusOK, pedidoAdmin(p, true))
}

// ---------------------------------------------------------------- productos

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
