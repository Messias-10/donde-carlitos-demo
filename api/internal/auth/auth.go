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

	"github.com/Messias-10/donde-carlitos-demo/api/internal/limite"
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

// BloqueoLogin aplica una espera progresiva por clave (usuario + IP): tras cada
// fallo hay que esperar 1 s, 2 s, 4 s, … hasta un tope (5 min). Los fallos se
// olvidan tras una ventana sin intentos o con un login correcto. Guarda como
// máximo 10.000 claves (LRU): un flood de claves nuevas desaloja las más
// antiguas en vez de llenar la memoria.
type BloqueoLogin struct {
	mu      sync.Mutex
	inicial time.Duration
	tope    time.Duration
	ventana time.Duration
	c       *limite.LRU[estadoLogin]
}

type estadoLogin struct {
	fallos int
	ultimo time.Time
	hasta  time.Time
}

// NuevoBloqueoLogin: espera inicial (1 s), tope (5 min) y ventana de olvido (15 min).
func NuevoBloqueoLogin(inicial, tope, ventana time.Duration) *BloqueoLogin {
	return &BloqueoLogin{inicial: inicial, tope: tope, ventana: ventana, c: limite.NuevoLRU[estadoLogin](10000)}
}

// Bloqueado indica si la clave debe esperar y cuánto falta.
func (b *BloqueoLogin) Bloqueado(clave string, ahora time.Time) (bool, time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.c.Obtener(clave)
	if !ok || !ahora.Before(e.hasta) {
		return false, 0
	}
	return true, e.hasta.Sub(ahora)
}

// Fallo registra un intento fallido y devuelve cuántos fallos seguidos lleva y
// cuánto hay que esperar para el próximo intento.
func (b *BloqueoLogin) Fallo(clave string, ahora time.Time) (int, time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.c.Obtener(clave)
	if !ok {
		e = b.c.Poner(clave, estadoLogin{})
	}
	if !e.ultimo.IsZero() && ahora.Sub(e.ultimo) > b.ventana {
		e.fallos = 0
	}
	e.fallos++
	e.ultimo = ahora
	espera := b.inicial
	for i := 1; i < e.fallos && espera < b.tope; i++ {
		espera *= 2
	}
	if espera > b.tope {
		espera = b.tope
	}
	e.hasta = ahora.Add(espera)
	return e.fallos, espera
}

// Exito limpia los fallos de la clave.
func (b *BloqueoLogin) Exito(clave string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.c.Borrar(clave)
}
