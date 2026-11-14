package productos

import (
	"testing"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

func TestDisponibilidad(t *testing.T) {
	for _, c := range []struct {
		stock, umbral int
		want          string
	}{{0, 5, Agotado}, {1, 5, QuedanPocos}, {5, 5, QuedanPocos}, {6, 5, Disponible}, {3, 0, Disponible}} {
		if got := CalcularDisponibilidad(c.stock, c.umbral); got != c.want {
			t.Errorf("%d/%d = %s", c.stock, c.umbral, got)
		}
	}
}

func TestPrecioUnitario(t *testing.T) {
	of, my, min := int64(17000), int64(15000), 6
	p := store.Producto{PrecioDetalle: 20000, PrecioOferta: &of, PrecioMayor: &my, MayorMinimo: &min}
	if v, tipo := PrecioUnitario(p, 1); v != 17000 || tipo != PrecioOferta {
		t.Errorf("1 u: %d %s", v, tipo)
	}
	if v, tipo := PrecioUnitario(p, 6); v != 15000 || tipo != PrecioMayor {
		t.Errorf("6 u: %d %s", v, tipo)
	}
	p.PrecioOferta = nil
	if v, tipo := PrecioUnitario(p, 5); v != 20000 || tipo != PrecioDetalle {
		t.Errorf("5 u sin oferta: %d %s", v, tipo)
	}
}

func TestBadges(t *testing.T) {
	hoy := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	of := int64(1000)
	p := store.Producto{PrecioDetalle: 1500, PrecioOferta: &of, EtiquetaOferta: "remate", FechaLlegada: "2026-11-10", Temporada: "navidad"}
	b := Badges(p, hoy, 14)
	if len(b) != 3 || b[0] != BadgeRemate || b[1] != BadgeRecienLlegado || b[2] != BadgeTemporada {
		t.Fatalf("badges = %v", b)
	}
	p.FechaLlegada = "2026-10-01"
	if EsRecienLlegado(p.FechaLlegada, hoy, 14) {
		t.Fatal("octubre no es recién llegado")
	}
	if EsRecienLlegado("2026-11-20", hoy, 14) {
		t.Fatal("fecha futura (preventa) no cuenta como llegado")
	}
}

func TestNormalizarYSlug(t *testing.T) {
	if got := NormalizarBusqueda("  Telón   NAVIDEÑO  Ñandú "); got != "telon navideno nandu" {
		t.Errorf("normalizar = %q", got)
	}
	if got := Slugify("Telón de Navidad 3 m x 2,5 m"); got != "telon-de-navidad-3-m-x-2-5-m" {
		t.Errorf("slug = %q", got)
	}
	if got := Slugify("Peluche gigante 'patas largas'"); got != "peluche-gigante-patas-largas" {
		t.Errorf("slug = %q", got)
	}
}

func TestValidarEntrada(t *testing.T) {
	activo := true
	ok := Entrada{Codigo: "p001", Nombre: "Farol navideño", CategoriaID: 1, Precio: 1500, Stock: 10, Activo: &activo}
	p, errs := ok.Validar()
	if len(errs) > 0 {
		t.Fatalf("errores: %v", errs)
	}
	if p.Slug != "farol-navideno" || p.Codigo != "P001" || p.UmbralPocos != UmbralPocosDefect || p.EtiquetaOferta != BadgeOferta {
		t.Fatalf("defaults: %+v", p)
	}

	s := func(v string) *string { return &v }
	i := func(v int64) *int64 { return &v }
	n := func(v int) *int { return &v }
	casos := map[string]func(*Entrada){
		"precio":        func(e *Entrada) { e.Precio = 0 },
		"precio_oferta": func(e *Entrada) { e.PrecioOferta = i(1500) },
		"precio_mayor":  func(e *Entrada) { e.PrecioMayor = i(1000) }, // sin mínimo
		"mayor_minimo":  func(e *Entrada) { e.PrecioMayor, e.MayorMinimo = i(1000), n(1) },
		"stock":         func(e *Entrada) { e.Stock = -1 },
		"temporada":     func(e *Entrada) { e.Temporada = s("pascua") },
		"tiktok_url":    func(e *Entrada) { e.TikTokURL = s("https://evil.example/@x/video/1") },
		"imagen_url":    func(e *Entrada) { e.ImagenURL = s("javascript:alert(1)") },
		"slug":          func(e *Entrada) { e.Slug = "Con Espacios" },
		"codigo":        func(e *Entrada) { e.Codigo = "x" },
		"fecha_llegada": func(e *Entrada) { e.FechaLlegada = s("14/11/2026") },
		"nombre":        func(e *Entrada) { e.Nombre = "a\x00b" },
	}
	for campo, mod := range casos {
		e := ok
		mod(&e)
		if _, errs := e.Validar(); errs[campo] == "" {
			t.Errorf("esperaba error en %s, obtuve %v", campo, errs)
		}
	}
	for _, u := range []string{"//evil.example/x.png", "/../etc/passwd", "http://x.cl/a.png"} {
		if EsImagenValida(u) {
			t.Errorf("imagen %q debió rechazarse", u)
		}
	}
	if !EsURLTikTok("https://www.tiktok.com/@mayte_emporio/video/7707649267464293639") {
		t.Error("tiktok válido rechazado")
	}
}
