package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

const (
	maxBodyLogin     = 4 << 10
	pageSizeAdminMax = 100
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

	var in entradaLogin
	if err := leerJSON(w, r, &in, maxBodyLogin); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	// Espera progresiva por usuario + IP: un tercero que prueba claves de "admin"
	// no deja afuera al resto del equipo, y el tope es de minutos, no de horas.
	clave := strings.ToLower(strings.TrimSpace(in.Usuario)) + "|" + ip
	if bloq, espera := s.bloqueo.Bloqueado(clave, ahora); bloq {
		w.Header().Set("Retry-After", strconv.Itoa(segundos(espera)))
		responderError(w, http.StatusTooManyRequests, "login_bloqueado", "Demasiados intentos. Espera un momento y vuelve a intentar.", nil)
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
	var u store.Usuario
	var err error
	valido := len(in.Usuario) >= 1 && len(in.Usuario) <= 64 && len(in.Password) >= 1 && len(in.Password) <= 72
	if valido {
		u, err = s.st.UsuarioPorUsuario(r.Context(), in.Usuario)
		switch {
		case err == nil:
			hash = u.PasswordHash
			valido = u.Activo // desactivado: misma respuesta que clave incorrecta
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
		n, espera := s.bloqueo.Fallo(clave, ahora)
		if n >= 5 {
			s.log.Warn("ALERTA: intentos de login fallidos repetidos", "fallos", n, "espera_s", int(espera.Seconds()),
				"usuario", recortar(in.Usuario, 32), "ip", enmascararIP(s.ipCliente(r)))
		} else {
			s.log.Warn("login fallido", "fallos", n, "ip", enmascararIP(s.ipCliente(r)))
		}
		responderError(w, http.StatusUnauthorized, "credenciales_invalidas", "Usuario o clave incorrectos.", nil)
		return
	}
	s.bloqueo.Exito(clave)

	token, err := auth.NuevoToken()
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	ses := store.Sesion{
		TokenHash: auth.HashToken(s.cfg.SecretoSesion, token),
		UsuarioID: u.ID,
		CreadoEn:  ahora,
		ExpiraEn:  ahora.Add(s.cfg.DuracionSesion),
		UltimoUso: ahora,
	}
	if err := s.st.CrearSesion(r.Context(), ses); err != nil {
		s.errorInterno(w, r, err)
		return
	}
	http.SetCookie(w, s.cookieSesion(token, int(s.cfg.DuracionSesion.Seconds())))
	s.log.Info("login correcto", "usuario", u.Usuario, "rol", u.Rol)
	responderJSON(w, http.StatusOK, sesionDTO{ID: u.ID, Usuario: u.Usuario, Nombre: u.Nombre, Rol: u.Rol,
		CanalPorDefecto: ptr(u.CanalPorDefecto), DebeCambiarClave: u.DebeCambiarClave, ExpiraEn: ses.ExpiraEn})
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
		// Clave temporal: hasta cambiarla solo puede ver su sesión, cambiar la clave o salir.
		if ses.DebeCambiarClave && r.URL.Path != "/api/admin/sesion" && r.URL.Path != "/api/admin/sesion/clave" {
			responderError(w, http.StatusForbidden, "debe_cambiar_clave", "Cambia tu clave temporal para seguir.", nil)
			return
		}
		h(w, conValor(r, ctxSesion, ses), ses)
	})
}

// panel: cualquier usuario activo (admin o vendedor).
func (s *Server) panel(h handlerConSesion) http.Handler { return s.conSesion(h) }

// soloAdmin: además exige rol admin (si no, 403 sin_permiso).
func (s *Server) soloAdmin(h handlerConSesion) http.Handler {
	return s.conSesion(func(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
		if ses.Rol != store.RolAdmin {
			responderError(w, http.StatusForbidden, "sin_permiso", "Esta sección es solo para administradores.", nil)
			return
		}
		h(w, r, ses)
	})
}

type sesionDTO struct {
	ID               int64     `json:"id"`
	Usuario          string    `json:"usuario"`
	Nombre           string    `json:"nombre"`
	Rol              string    `json:"rol"`
	CanalPorDefecto  *string   `json:"canal_por_defecto"`
	DebeCambiarClave bool      `json:"debe_cambiar_clave"`
	ExpiraEn         time.Time `json:"expira_en"`
}

func (s *Server) sesionActual(w http.ResponseWriter, _ *http.Request, ses store.Sesion) {
	responderJSON(w, http.StatusOK, sesionDTO{ID: ses.UsuarioID, Usuario: ses.Usuario, Nombre: ses.Nombre, Rol: ses.Rol,
		CanalPorDefecto: ptr(ses.CanalPorDefecto), DebeCambiarClave: ses.DebeCambiarClave, ExpiraEn: ses.ExpiraEn.UTC()})
}

// MinLargoClave para claves elegidas por el usuario.
const MinLargoClave = 12

// cambiarClavePropia: POST /api/admin/sesion/clave {clave_actual, clave_nueva}.
// Verifica la actual, guarda la nueva (bcrypt), quita debe_cambiar_clave y
// cierra las demás sesiones del usuario (la actual sigue).
func (s *Server) cambiarClavePropia(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	var in struct {
		ClaveActual string `json:"clave_actual"`
		ClaveNueva  string `json:"clave_nueva"`
	}
	if err := leerJSON(w, r, &in, maxBodyLogin); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	campos := map[string]string{}
	if len([]rune(in.ClaveNueva)) < MinLargoClave || len(in.ClaveNueva) > 72 {
		campos["clave_nueva"] = "entre 12 y 72 caracteres"
	} else if in.ClaveNueva == in.ClaveActual {
		campos["clave_nueva"] = "debe ser distinta de la actual"
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa la clave nueva.", campos)
		return
	}
	if !s.semBcrypt.Intentar() {
		w.Header().Set("Retry-After", "1")
		responderError(w, http.StatusTooManyRequests, "demasiadas_solicitudes", "Hay mucha actividad. Intenta de nuevo en un segundo.", nil)
		return
	}
	defer s.semBcrypt.Soltar()
	u, err := s.st.UsuarioPorID(r.Context(), ses.UsuarioID)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	if !auth.VerificarClave(u.PasswordHash, in.ClaveActual) {
		s.log.Warn("cambio de clave con clave actual incorrecta", "usuario", ses.Usuario)
		responderError(w, http.StatusBadRequest, "validacion", "La clave actual no es correcta.", map[string]string{"clave_actual": "no coincide"})
		return
	}
	hash, err := auth.HashClave(in.ClaveNueva, s.cfg.CostoBcryptUsuario)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	if err := s.st.CambiarClave(r.Context(), ses.UsuarioID, hash, false, ses.TokenHash, s.ahora()); err != nil {
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("clave cambiada por su dueño", "usuario", ses.Usuario)
	w.WriteHeader(http.StatusNoContent)
}

func idDeRuta(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}
