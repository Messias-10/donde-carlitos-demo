package httpapi

import (
	"errors"
	"net/http"
	"strconv"
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
	responderJSON(w, http.StatusOK, sesionDTO{ID: u.ID, Usuario: u.Usuario, Nombre: u.Nombre, Rol: u.Rol, ExpiraEn: ses.ExpiraEn})
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
	ID       int64     `json:"id"`
	Usuario  string    `json:"usuario"`
	Nombre   string    `json:"nombre"`
	Rol      string    `json:"rol"`
	ExpiraEn time.Time `json:"expira_en"`
}

func (s *Server) sesionActual(w http.ResponseWriter, _ *http.Request, ses store.Sesion) {
	responderJSON(w, http.StatusOK, sesionDTO{ID: ses.UsuarioID, Usuario: ses.Usuario, Nombre: ses.Nombre, Rol: ses.Rol, ExpiraEn: ses.ExpiraEn.UTC()})
}

func idDeRuta(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}
