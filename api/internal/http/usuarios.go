package httpapi

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

// largoClaveTemporal: 16 caracteres de 66 símbolos ≈ 96 bits.
const largoClaveTemporal = 16

var reUsuario = regexp.MustCompile(`^[a-z][a-z0-9._-]{2,31}$`)

type usuarioDTO struct {
	ID               int64     `json:"id"`
	Usuario          string    `json:"usuario"`
	Nombre           string    `json:"nombre"`
	Rol              string    `json:"rol"`
	Activo           bool      `json:"activo"`
	CanalPorDefecto  *string   `json:"canal_por_defecto"`
	DebeCambiarClave bool      `json:"debe_cambiar_clave"`
	CreadoEn         time.Time `json:"creado_en"`
}

func usuarioAPI(u store.Usuario) usuarioDTO {
	return usuarioDTO{ID: u.ID, Usuario: u.Usuario, Nombre: u.Nombre, Rol: u.Rol, Activo: u.Activo,
		CanalPorDefecto: ptr(u.CanalPorDefecto), DebeCambiarClave: u.DebeCambiarClave, CreadoEn: u.CreadoEn.UTC()}
}

func esCanalUsuario(c string) bool { return c == "" || c == "santiago" || c == "regiones" }

// protegerPrincipal: nadie más que el propio admin principal (ADMIN_USUARIO de
// api/.env) puede editarlo ni cambiarle la clave (seguridad B8).
func (s *Server) protegerPrincipal(w http.ResponseWriter, r *http.Request, objetivo int64, ses store.Sesion) bool {
	if s.cfg.AdminPrincipal == "" || objetivo == ses.UsuarioID {
		return true
	}
	u, err := s.st.UsuarioPorID(r.Context(), objetivo)
	if err == nil && u.Usuario == s.cfg.AdminPrincipal {
		responderError(w, http.StatusForbidden, "admin_principal_protegido", "El administrador principal solo se edita a sí mismo.", nil)
		return false
	}
	return true
}

type usuarioConClaveDTO struct {
	Usuario       usuarioDTO `json:"usuario"`
	ClaveTemporal string     `json:"clave_temporal"`
}

func esRol(r string) bool { return r == store.RolAdmin || r == store.RolVendedor }

func validarNombreVisible(n string) (string, bool) {
	limpio, ok := texto.Limpiar(n, false)
	return limpio, ok && texto.Largo(limpio) >= 2 && texto.Largo(limpio) <= 60
}

func (s *Server) listarUsuarios(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	us, err := s.st.ListarUsuarios(r.Context())
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	out := make([]usuarioDTO, 0, len(us))
	for _, u := range us {
		out = append(out, usuarioAPI(u))
	}
	responderJSON(w, http.StatusOK, map[string]any{"items": out})
}

// claveNueva genera la clave y su hash. La clave solo viaja en la respuesta.
func (s *Server) claveNueva() (string, string, error) {
	clave, err := auth.GenerarClave(largoClaveTemporal)
	if err != nil {
		return "", "", err
	}
	hash, err := auth.HashClave(clave, s.cfg.CostoBcryptUsuario)
	return clave, hash, err
}

func (s *Server) crearUsuario(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	var in struct {
		Usuario         string `json:"usuario"`
		Nombre          string `json:"nombre"`
		Rol             string `json:"rol"`
		CanalPorDefecto string `json:"canal_por_defecto"`
	}
	if err := leerJSON(w, r, &in, 4<<10); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	campos := map[string]string{}
	if !reUsuario.MatchString(in.Usuario) {
		campos["usuario"] = "3 a 32 caracteres: minúsculas, números, . _ -"
	}
	nombre, ok := validarNombreVisible(in.Nombre)
	if !ok {
		campos["nombre"] = "2 a 60 caracteres"
	}
	if !esRol(in.Rol) {
		campos["rol"] = "admin o vendedor"
	}
	if !esCanalUsuario(in.CanalPorDefecto) {
		campos["canal_por_defecto"] = "santiago, regiones o vacío"
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del usuario.", campos)
		return
	}
	clave, hash, err := s.claveNueva()
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	u := store.Usuario{Usuario: in.Usuario, Nombre: nombre, Rol: in.Rol, PasswordHash: hash,
		CanalPorDefecto: in.CanalPorDefecto, DebeCambiarClave: true}
	if err := s.st.CrearUsuario(r.Context(), &u, s.ahora()); err != nil {
		if errors.Is(err, store.ErrDuplicado) {
			responderError(w, http.StatusConflict, "duplicado", "Ya existe un usuario con ese nombre de usuario.", nil)
			return
		}
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("usuario creado", "usuario", u.Usuario, "rol", u.Rol, "por", ses.Usuario)
	responderJSON(w, http.StatusCreated, usuarioConClaveDTO{Usuario: usuarioAPI(u), ClaveTemporal: clave})
}

func (s *Server) actualizarUsuario(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Usuario no encontrado.", nil)
		return
	}
	var in struct {
		Nombre          *string `json:"nombre"`
		Rol             *string `json:"rol"`
		Activo          *bool   `json:"activo"`
		CanalPorDefecto *string `json:"canal_por_defecto"`
	}
	if err := leerJSON(w, r, &in, 4<<10); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	campos := map[string]string{}
	var c store.CambiosUsuario
	if in.Nombre != nil {
		n, ok := validarNombreVisible(*in.Nombre)
		if !ok {
			campos["nombre"] = "2 a 60 caracteres"
		}
		c.Nombre = &n
	}
	if in.Rol != nil {
		if !esRol(*in.Rol) {
			campos["rol"] = "admin o vendedor"
		}
		c.Rol = in.Rol
	}
	c.Activo = in.Activo
	if in.CanalPorDefecto != nil {
		if !esCanalUsuario(*in.CanalPorDefecto) {
			campos["canal_por_defecto"] = "santiago, regiones o vacío"
		}
		c.CanalPorDefecto = in.CanalPorDefecto
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del usuario.", campos)
		return
	}
	if !s.protegerPrincipal(w, r, id, ses) {
		return
	}
	if id == ses.UsuarioID && ((c.Activo != nil && !*c.Activo) || (c.Rol != nil && *c.Rol != store.RolAdmin)) {
		responderError(w, http.StatusConflict, "no_permitido_sobre_si_mismo", "No puedes desactivarte ni quitarte el rol de admin.", nil)
		return
	}
	u, err := s.st.ActualizarUsuario(r.Context(), id, c, s.ahora())
	if errors.Is(err, store.ErrNoEncontrado) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Usuario no encontrado.", nil)
		return
	}
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("usuario actualizado", "usuario", u.Usuario, "activo", u.Activo, "rol", u.Rol, "por", ses.Usuario)
	responderJSON(w, http.StatusOK, usuarioAPI(u))
}

func (s *Server) regenerarClave(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Usuario no encontrado.", nil)
		return
	}
	if !s.protegerPrincipal(w, r, id, ses) {
		return
	}
	u, err := s.st.UsuarioPorID(r.Context(), id)
	if errors.Is(err, store.ErrNoEncontrado) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Usuario no encontrado.", nil)
		return
	}
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	clave, hash, err := s.claveNueva()
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	// Clave temporal: al entrar deberá cambiarla (B9). Se cierran sus sesiones.
	if err := s.st.CambiarClave(r.Context(), id, hash, true, "", s.ahora()); err != nil {
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("clave regenerada", "usuario", u.Usuario, "por", ses.Usuario)
	u.DebeCambiarClave = true
	responderJSON(w, http.StatusOK, usuarioConClaveDTO{Usuario: usuarioAPI(u), ClaveTemporal: clave})
}
