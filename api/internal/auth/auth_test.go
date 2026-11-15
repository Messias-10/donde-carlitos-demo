package auth

import (
	"testing"
	"time"
)

func TestBloqueoProgresivo(t *testing.T) {
	b := NuevoBloqueoLogin(time.Second, 5*time.Minute, 15*time.Minute)
	t0 := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	ahora := t0
	var esperas []time.Duration
	for i := 0; i < 11; i++ {
		_, e := b.Fallo("admin|10.0.0.1", ahora)
		esperas = append(esperas, e)
		if bloq, _ := b.Bloqueado("admin|10.0.0.1", ahora); !bloq {
			t.Fatalf("fallo %d sin espera", i+1)
		}
		ahora = ahora.Add(e)
	}
	want := []time.Duration{1, 2, 4, 8, 16, 32, 64, 128, 256, 300, 300}
	for i, w := range want {
		if esperas[i] != w*time.Second {
			t.Fatalf("espera %d = %v, quería %vs", i+1, esperas[i], w)
		}
	}
	// Otra combinación usuario+IP no queda bloqueada.
	if bloq, _ := b.Bloqueado("vane|10.0.0.1", ahora); bloq {
		t.Fatal("otra clave bloqueada")
	}
	// Tras la ventana sin intentos, vuelve a empezar en 1 s.
	if _, e := b.Fallo("admin|10.0.0.1", ahora.Add(20*time.Minute)); e != time.Second {
		t.Fatalf("tras ventana: %v", e)
	}
	b.Exito("admin|10.0.0.1")
	if bloq, _ := b.Bloqueado("admin|10.0.0.1", ahora.Add(20*time.Minute)); bloq {
		t.Fatal("éxito limpia")
	}
}
