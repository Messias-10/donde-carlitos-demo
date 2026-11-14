// Package auth maneja claves (bcrypt), tokens de sesión y el bloqueo temporal
// tras intentos fallidos de login.
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// CostoBcrypt para claves de admin.
const CostoBcrypt = 12

// LargoToken en bytes aleatorios (256 bits).
const LargoToken = 32

// alfabeto sin caracteres ambiguos (0/O, 1/l/I) para que la clave se pueda tipear.
const alfabeto = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789-_.!@#%+="

// GenerarClave devuelve una clave aleatoria de n caracteres (crypto/rand).
// Con n=24 y 66 símbolos son ~145 bits de entropía.
func GenerarClave(n int) (string, error) {
	if n < 16 {
		return "", errors.New("auth: la clave debe tener al menos 16 caracteres")
	}
	var b strings.Builder
	max := big.NewInt(int64(len(alfabeto)))
	for i := 0; i < n; i++ {
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b.WriteByte(alfabeto[k.Int64()])
	}
	return b.String(), nil
}

// HashClave genera el hash bcrypt con el costo indicado.
func HashClave(clave string, costo int) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(clave), costo)
	return string(h), err
}

// VerificarClave compara en tiempo constante (bcrypt).
func VerificarClave(hash, clave string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(clave)) == nil
}

// EsHashBcrypt hace una verificación de forma.
func EsHashBcrypt(h string) bool {
	_, err := bcrypt.Cost([]byte(h))
	return err == nil
}

// NuevoToken devuelve un token aleatorio en base64url sin relleno (43 caracteres).
func NuevoToken() (string, error) {
	b := make([]byte, LargoToken)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// TokenConForma valida el formato antes de tocar la base.
func TokenConForma(t string) bool {
	if len(t) != 43 {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(t)
	return err == nil
}

// HashToken es lo que se guarda en la tabla sesiones: HMAC-SHA256(secreto, token).
// Si alguien lee la base, no puede usar las sesiones sin el secreto.
func HashToken(secreto []byte, token string) string {
	m := hmac.New(sha256.New, secreto)
	m.Write([]byte(token))
	return hex.EncodeToString(m.Sum(nil))
}

// BloqueoLogin cuenta fallos por clave (IP) y bloquea temporalmente.
type BloqueoLogin struct {
	mu          sync.Mutex
	maxFallos   int
	ventana     time.Duration
	bloqueo     time.Duration
	maxEntradas int
	m           map[string]*estadoLogin
}

type estadoLogin struct {
	fallos int
	desde  time.Time
	hasta  time.Time
}

// NuevoBloqueoLogin: tras maxFallos dentro de ventana, bloquea por bloqueo.
func NuevoBloqueoLogin(maxFallos int, ventana, bloqueo time.Duration) *BloqueoLogin {
	return &BloqueoLogin{maxFallos: maxFallos, ventana: ventana, bloqueo: bloqueo, maxEntradas: 10000, m: map[string]*estadoLogin{}}
}

// Bloqueado indica si la clave está bloqueada y cuánto falta.
func (b *BloqueoLogin) Bloqueado(clave string, ahora time.Time) (bool, time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.m[clave]
	if !ok || !ahora.Before(e.hasta) {
		return false, 0
	}
	return true, e.hasta.Sub(ahora)
}

// Fallo registra un intento fallido. Devuelve true si con este quedó bloqueada.
func (b *BloqueoLogin) Fallo(clave string, ahora time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.m[clave]
	if !ok {
		if len(b.m) >= b.maxEntradas {
			b.purgar(ahora)
		}
		e = &estadoLogin{desde: ahora}
		b.m[clave] = e
	}
	if ahora.Sub(e.desde) > b.ventana {
		e.fallos, e.desde = 0, ahora
	}
	e.fallos++
	if e.fallos >= b.maxFallos {
		e.hasta = ahora.Add(b.bloqueo)
		e.fallos, e.desde = 0, ahora
		return true
	}
	return false
}

// Exito limpia los fallos de la clave.
func (b *BloqueoLogin) Exito(clave string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.m, clave)
}

func (b *BloqueoLogin) purgar(ahora time.Time) {
	for k, e := range b.m {
		if !ahora.Before(e.hasta) && ahora.Sub(e.desde) > b.ventana {
			delete(b.m, k)
		}
	}
}
