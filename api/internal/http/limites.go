package httpapi

import (
	"sync"
	"time"
)

// LimitadorVentana permite como máximo n solicitudes por clave dentro de una
// ventana fija. En memoria: suficiente para una sola instancia (demo).
type LimitadorVentana struct {
	mu          sync.Mutex
	n           int
	ventana     time.Duration
	maxEntradas int
	m           map[string]*ventana
}

type ventana struct {
	inicio time.Time
	usadas int
}

// NuevoLimitadorVentana crea un limitador.
func NuevoLimitadorVentana(n int, v time.Duration) *LimitadorVentana {
	return &LimitadorVentana{n: n, ventana: v, maxEntradas: 20000, m: map[string]*ventana{}}
}

// Permitir consume un cupo. Si no hay, devuelve false y cuánto esperar.
func (l *LimitadorVentana) Permitir(clave string, ahora time.Time) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.m[clave]
	if !ok || ahora.Sub(v.inicio) >= l.ventana {
		if !ok && len(l.m) >= l.maxEntradas {
			for k, e := range l.m {
				if ahora.Sub(e.inicio) >= l.ventana {
					delete(l.m, k)
				}
			}
			if len(l.m) >= l.maxEntradas {
				// Bajo un ataque con muchas IPs se falla cerrado.
				return false, l.ventana
			}
		}
		l.m[clave] = &ventana{inicio: ahora, usadas: 1}
		return true, 0
	}
	if v.usadas >= l.n {
		return false, l.ventana - ahora.Sub(v.inicio)
	}
	v.usadas++
	return true, 0
}
