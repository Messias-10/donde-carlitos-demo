package httpapi_test

// Tests de las correcciones de la auditoría de seguridad (A1, A2, M1) y de QA
// ronda 1 (I-1, I-2, I-3, I-4, M-10, M-15).

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	httpapi "github.com/Messias-10/donde-carlitos-demo/api/internal/http"
)

func TestIPSoloUltimoValorXFFDesdeProxyConfiable(t *testing.T) {
	e := nuevoEntorno(t, nil)
	malo := map[string]string{"usuario": "admin", "password": "incorrecta"}
	for i := 0; i < 5; i++ {
		// El cliente cambia el primer valor; el último lo pone el proxy y no cambia.
		xff := fmt.Sprintf("198.51.100.%d, 203.0.113.50", i)
		e.pedir(t, "POST", "/api/admin/login", malo, map[string]string{"X-Forwarded-For": xff})
	}
	r := e.pedir(t, "POST", "/api/admin/login", malo, map[string]string{"X-Forwarded-For": "198.51.100.99, 203.0.113.50"})
	if r.status != 429 {
		t.Fatalf("rotar valores anteriores de XFF no debe evadir el bloqueo: %d", r.status)
	}
}

func TestXFFIgnoradoSinProxyConfiable(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.ProxiesConfiables = nil })
	malo := map[string]string{"usuario": "admin", "password": "incorrecta"}
	for i := 0; i < 5; i++ {
		e.pedir(t, "POST", "/api/admin/login", malo, map[string]string{"X-Forwarded-For": fmt.Sprintf("203.0.113.%d", i)})
	}
	r := e.pedir(t, "POST", "/api/admin/login", malo, map[string]string{"X-Forwarded-For": "203.0.113.200"})
	if r.status != 429 {
		t.Fatalf("sin proxy confiable se limita por RemoteAddr: %d", r.status)
	}
}

func TestLoginTopeGlobalPorMinuto(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.LoginGlobalMinuto = 3 })
	malo := map[string]string{"usuario": "admin", "password": "x"}
	for i := 0; i < 3; i++ {
		e.pedir(t, "POST", "/api/admin/login", malo, map[string]string{"X-Forwarded-For": fmt.Sprintf("203.0.113.%d", i)})
	}
	r := e.pedir(t, "POST", "/api/admin/login", malo, map[string]string{"X-Forwarded-For": "203.0.113.99"})
	if r.status != 429 || r.header.Get("Retry-After") == "" {
		t.Fatalf("tope global de login: %d", r.status)
	}
}

type productoQA struct {
	Codigo                string     `json:"codigo"`
	Precio                int64      `json:"precio"`
	PrecioOferta          *int64     `json:"precio_oferta"`
	PrecioFinal           int64      `json:"precio_final"`
	OfertaHasta           *string    `json:"oferta_hasta"`
	OfertaVigente         bool       `json:"oferta_vigente"`
	MayorMinimoConfirmado bool       `json:"mayor_minimo_confirmado"`
	PrecioActualizadoEn   *time.Time `json:"precio_actualizado_en"`
	Disponibilidad        string     `json:"disponibilidad"`
	Badges                []string   `json:"badges"`
}

func (e *entorno) productoQA(t *testing.T, slug string) productoQA {
	t.Helper()
	var p productoQA
	e.pedir(t, "GET", "/api/productos/"+slug, nil, nil).json(t, &p)
	return p
}

func TestCamposQA(t *testing.T) {
	e := nuevoEntorno(t, nil)
	// I-1: sin urgencia falsa en la semilla.
	var pg struct {
		Items []productoQA `json:"items"`
	}
	e.pedir(t, "GET", "/api/productos?page_size=48", nil, nil).json(t, &pg)
	for _, p := range pg.Items {
		if p.Disponibilidad == "quedan_pocos" || (p.Disponibilidad == "agotado" && p.Codigo != "P036") {
			t.Errorf("%s: disponibilidad %s", p.Codigo, p.Disponibilidad)
		}
		if p.PrecioActualizadoEn != nil {
			t.Errorf("%s: precio_actualizado_en debe ser null en la semilla", p.Codigo)
		}
	}
	// I-2
	if !e.productoQA(t, "squishy-viral-que-cambia-de-color-con-el-sol").MayorMinimoConfirmado {
		t.Error("P001: mínimo del video")
	}
	if e.productoQA(t, "camino-de-mesa-navideno").MayorMinimoConfirmado {
		t.Error("P009: mínimo [EJEMPLO]")
	}
	// M-15: la oferta del tren vale hasta el domingo 15.
	tren := e.productoQA(t, "tren-navideno-para-arbol-grande")
	if !tren.OfertaVigente || tren.PrecioFinal != 17000 || tren.OfertaHasta == nil || *tren.OfertaHasta != "2026-11-15" {
		t.Fatalf("tren el 14: %+v", tren)
	}

	// I-3: el admin cambia el precio -> precio_actualizado_en en UTC.
	c := e.login(t)
	r := e.pedir(t, "GET", "/api/admin/productos?q=cambia+de+color", nil, nil, c)
	var lista struct {
		Items []map[string]any `json:"items"`
	}
	r.json(t, &lista)
	adm := lista.Items[0]
	id := int64(adm["id"].(float64))
	cat := adm["categoria"].(map[string]any)
	put := map[string]any{"slug": adm["slug"], "codigo": adm["codigo"], "nombre": adm["nombre"], "descripcion": adm["descripcion"],
		"categoria_id": cat["id"], "precio": 3990, "precio_mayor": 3500, "mayor_minimo": 12, "unidad_mayor": "caja de 12",
		"mayor_minimo_confirmado": true, "stock": 40, "activo": true}
	r = e.pedir(t, "PUT", fmt.Sprintf("/api/admin/productos/%d", id), put, nil, c)
	if r.status != 200 || !bytes.Contains(r.body, []byte(`"precio_actualizado_en":"2026-11-14T15:00:00Z"`)) {
		t.Fatalf("precio_actualizado_en tras cambiar precio: %d %s", r.status, r.body)
	}
	p := e.productoQA(t, adm["slug"].(string))
	if p.PrecioActualizadoEn == nil || p.PrecioActualizadoEn.Location() != time.UTC || !p.MayorMinimoConfirmado {
		t.Fatalf("público: %+v", p)
	}
	// oferta_hasta sin precio_oferta es inválido.
	put["oferta_hasta"] = "2026-12-01"
	if r := e.pedir(t, "PUT", fmt.Sprintf("/api/admin/productos/%d", id), put, nil, c); r.status != 400 {
		t.Fatalf("oferta_hasta sin oferta: %d", r.status)
	}
}

func TestOfertaVencidaDejaDeAplicar(t *testing.T) {
	e := nuevoEntornoEn(t, time.Date(2026, 11, 16, 15, 0, 0, 0, time.UTC), nil)
	tren := e.productoQA(t, "tren-navideno-para-arbol-grande")
	if tren.OfertaVigente || tren.PrecioOferta != nil || tren.PrecioFinal != 20000 {
		t.Fatalf("oferta vencida sigue aplicando: %+v", tren)
	}
	for _, b := range tren.Badges {
		if b == "oferta" || b == "remate" {
			t.Fatalf("badge de oferta vencida: %v", tren.Badges)
		}
	}
}

func pngPrueba(t *testing.T) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for x := 0; x < 8; x++ {
		img.Set(x, x, color.RGBA{226, 179, 60, 255})
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func (e *entorno) subir(t *testing.T, id int64, nombre string, datos []byte, c *http.Cookie) resp {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, _ := mw.CreateFormFile("imagen", nombre)
	fw.Write(datos)
	mw.Close()
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/admin/productos/%d/imagen", e.srv.URL, id), &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Origin", origen)
	if c != nil {
		req.AddCookie(c)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var b bytes.Buffer
	b.ReadFrom(res.Body)
	return resp{res.StatusCode, res.Header, b.Bytes()}
}

func TestSubidaDeImagen(t *testing.T) {
	e := nuevoEntorno(t, nil)
	if r := e.subir(t, 1, "a.png", pngPrueba(t), nil); r.status != 401 {
		t.Fatalf("sin sesión: %d", r.status)
	}
	c := e.login(t)

	r := e.subir(t, 1, "foto.png", pngPrueba(t), c)
	if r.status != 200 {
		t.Fatalf("subida: %d %s", r.status, r.body)
	}
	var p struct {
		ImagenURL string `json:"imagen_url"`
	}
	json.Unmarshal(r.body, &p)
	if !strings.HasPrefix(p.ImagenURL, "/api/imagenes/") || !strings.HasSuffix(p.ImagenURL, ".png") || strings.Contains(p.ImagenURL, "foto") {
		t.Fatalf("imagen_url: %s", p.ImagenURL)
	}
	g := e.pedir(t, "GET", p.ImagenURL, nil, nil)
	if g.status != 200 || g.header.Get("Content-Type") != "image/png" || g.header.Get("X-Content-Type-Options") != "nosniff" ||
		!strings.Contains(g.header.Get("Content-Security-Policy"), "default-src 'none'") {
		t.Fatalf("servir: %d %v", g.status, g.header)
	}
	if _, err := png.Decode(bytes.NewReader(g.body)); err != nil {
		t.Fatalf("la imagen servida no es PNG válido: %v", err)
	}

	// Rechazos: HTML disfrazado, PNG truncado, demasiado grande.
	if r := e.subir(t, 1, "x.png", []byte("<html><script>alert(1)</script></html>"), c); r.status != 415 {
		t.Fatalf("html disfrazado: %d", r.status)
	}
	if r := e.subir(t, 1, "x.png", pngPrueba(t)[:20], c); r.status != 400 {
		t.Fatalf("png truncado: %d", r.status)
	}
	if r := e.subir(t, 1, "x.jpg", append([]byte{0xFF, 0xD8, 0xFF}, make([]byte, 6<<20)...), c); r.status != 413 {
		t.Fatalf("más de 5 MB: %d", r.status)
	}
	if r := e.subir(t, 99999, "a.png", pngPrueba(t), c); r.status != 404 {
		t.Fatalf("producto inexistente: %d", r.status)
	}
	// Solo nombres generados por el servidor.
	for _, ruta := range []string{"/api/imagenes/..%2f..%2fetc%2fpasswd", "/api/imagenes/foto.png", "/api/imagenes/" + strings.Repeat("a", 32) + ".svg"} {
		if r := e.pedir(t, "GET", ruta, nil, nil); r.status != 404 {
			t.Errorf("%s: %d", ruta, r.status)
		}
	}
}

func TestTimestampsUTC(t *testing.T) {
	santiago, _ := time.LoadLocation("America/Santiago")
	e := nuevoEntornoEn(t, time.Date(2026, 11, 14, 12, 0, 0, 0, santiago), nil)
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	if !bytes.Contains(r.body, []byte(`"expira_en":"2026-11-15T03:00:00Z"`)) {
		t.Fatalf("expira_en en UTC: %s", r.body)
	}
	c := e.login(t)
	r = e.pedir(t, "POST", "/api/admin/pedidos", map[string]any{"canal": "santiago", "telefono": "900000555",
		"items": []map[string]any{{"descripcion": "Algo", "cantidad": 1, "precio_unitario": 1000}}}, nil, c)
	if r.status != 201 || !bytes.Contains(r.body, []byte(`"creado_en":"2026-11-14T15:00:00Z"`)) {
		t.Fatalf("creado_en en UTC: %d %s", r.status, r.body)
	}
}
