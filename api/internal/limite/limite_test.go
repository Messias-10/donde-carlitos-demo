package limite

import (
	"strconv"
	"testing"
	"time"
)

func TestLRUDesalojaLaMasAntigua(t *testing.T) {
	c := NuevoLRU[int](3)
	c.Poner("a", 1)
	c.Poner("b", 2)
	c.Poner("c", 3)
	c.Obtener("a") // "a" pasa a ser la más reciente
	c.Poner("d", 4)
	if _, ok := c.Obtener("b"); ok {
		t.Fatal("b debió desalojarse")
	}
	for _, k := range []string{"a", "c", "d"} {
		if _, ok := c.Obtener(k); !ok {
			t.Fatalf("%s no debió desalojarse", k)
		}
	}
	if c.Len() != 3 {
		t.Fatalf("len = %d", c.Len())
	}
}

func TestVentana(t *testing.T) {
	t0 := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	v := NuevaVentana(2, time.Minute, 100)
	if ok, _ := v.Permitir("x", t0); !ok {
		t.Fatal("1")
	}
	if ok, _ := v.Permitir("x", t0); !ok {
		t.Fatal("2")
	}
	if ok, espera := v.Permitir("x", t0.Add(10*time.Second)); ok || espera != 50*time.Second {
		t.Fatalf("3: %v %v", ok, espera)
	}
	if ok, _ := v.Permitir("x", t0.Add(time.Minute)); !ok {
		t.Fatal("ventana nueva")
	}
}

// Llenar el limitador con claves nuevas no deja afuera a una clave legítima nueva
// (antes se fallaba cerrado al llegar al máximo).
func TestVentanaLlenaNoBloqueaATodos(t *testing.T) {
	t0 := time.Now()
	v := NuevaVentana(1, time.Hour, 1000)
	for i := 0; i < 20050; i++ {
		v.Permitir("atacante-"+strconv.Itoa(i), t0)
	}
	if ok, _ := v.Permitir("cliente-legitimo", t0); !ok {
		t.Fatal("un cliente nuevo debe poder pasar con el limitador lleno")
	}
	if v.c.Len() != 1000 {
		t.Fatalf("memoria acotada: %d", v.c.Len())
	}
}

func TestSemaforo(t *testing.T) {
	s := NuevoSemaforo(2)
	if !s.Intentar() || !s.Intentar() {
		t.Fatal("dos cupos")
	}
	if s.Intentar() {
		t.Fatal("tercero debe fallar")
	}
	s.Soltar()
	if !s.Intentar() {
		t.Fatal("tras soltar")
	}
}
