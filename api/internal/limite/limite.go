// Package limite tiene las piezas de limitación en memoria: un LRU de
// capacidad fija (para que un atacante con muchas claves no llene la memoria
// ni deje afuera a los usuarios legítimos), un limitador por ventana fija y un
// semáforo no bloqueante. Suficiente para una sola instancia (demo).
package limite

import (
	"container/list"
	"sync"
	"time"
)

// LRU es un mapa de capacidad fija que desaloja la entrada usada hace más
// tiempo. No es seguro para uso concurrente: lo protege quien lo usa.
type LRU[V any] struct {
	cap int
	ll  *list.List
	m   map[string]*list.Element
}

type nodo[V any] struct {
	clave string
	val   V
}

// NuevoLRU crea un LRU con capacidad >= 1.
func NuevoLRU[V any](capacidad int) *LRU[V] {
	if capacidad < 1 {
		capacidad = 1
	}
	return &LRU[V]{cap: capacidad, ll: list.New(), m: make(map[string]*list.Element, capacidad)}
}

// Obtener devuelve un puntero al valor y lo marca como recién usado.
func (c *LRU[V]) Obtener(clave string) (*V, bool) {
	e, ok := c.m[clave]
	if !ok {
		return nil, false
	}
	c.ll.MoveToFront(e)
	return &e.Value.(*nodo[V]).val, true
}

// Poner inserta o reemplaza; si está lleno, desaloja la más antigua.
func (c *LRU[V]) Poner(clave string, v V) *V {
	if e, ok := c.m[clave]; ok {
		n := e.Value.(*nodo[V])
		n.val = v
		c.ll.MoveToFront(e)
		return &n.val
	}
	if c.ll.Len() >= c.cap {
		if viejo := c.ll.Back(); viejo != nil {
			c.ll.Remove(viejo)
			delete(c.m, viejo.Value.(*nodo[V]).clave)
		}
	}
	n := &nodo[V]{clave: clave, val: v}
	c.m[clave] = c.ll.PushFront(n)
	return &n.val
}

// Borrar quita una clave.
func (c *LRU[V]) Borrar(clave string) {
	if e, ok := c.m[clave]; ok {
		c.ll.Remove(e)
		delete(c.m, clave)
	}
}

// Len cuenta las entradas.
func (c *LRU[V]) Len() int { return c.ll.Len() }

// Ventana permite como máximo n eventos por clave dentro de una ventana fija.
type Ventana struct {
	mu  sync.Mutex
	n   int
	dur time.Duration
	c   *LRU[cupo]
}

type cupo struct {
	inicio time.Time
	usadas int
}

// NuevaVentana crea un limitador de n eventos por dur, con capacidad claves en memoria.
func NuevaVentana(n int, dur time.Duration, capacidad int) *Ventana {
	return &Ventana{n: n, dur: dur, c: NuevoLRU[cupo](capacidad)}
}

// Permitir consume un cupo. Si no hay, devuelve false y cuánto falta.
func (v *Ventana) Permitir(clave string, ahora time.Time) (bool, time.Duration) {
	v.mu.Lock()
	defer v.mu.Unlock()
	c, ok := v.c.Obtener(clave)
	if !ok || ahora.Sub(c.inicio) >= v.dur {
		v.c.Poner(clave, cupo{inicio: ahora, usadas: 1})
		return v.n >= 1, 0
	}
	if c.usadas >= v.n {
		return false, v.dur - ahora.Sub(c.inicio)
	}
	c.usadas++
	return true, 0
}

// Semaforo limita cuántas operaciones caras corren a la vez, sin esperar.
type Semaforo chan struct{}

// NuevoSemaforo con n cupos.
func NuevoSemaforo(n int) Semaforo { return make(Semaforo, n) }

// Intentar toma un cupo si hay; si no, devuelve false de inmediato.
func (s Semaforo) Intentar() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

// Soltar devuelve el cupo.
func (s Semaforo) Soltar() { <-s }
