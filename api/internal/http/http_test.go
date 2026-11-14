package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	httpapi "github.com/Messias-10/donde-carlitos-demo/api/internal/http"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/seed"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store/sqlite"
)

const origen = "http://localhost:3000"

// claveTest se genera en cada corrida; no hay claves fijas en el código.
type entorno struct {
	srv   *httptest.Server
	clave string
	db    *sqlite.DB
}

func nuevoEntorno(t *testing.T, ajustar func(*httpapi.Config)) *entorno {
	t.Helper()
	return nuevoEntornoEn(t, time.Date(2026, 11, 14, 15, 0, 0, 0, time.UTC), ajustar)
}

func nuevoEntornoEn(t *testing.T, ahora time.Time, ajustar func(*httpapi.Config)) *entorno {
	t.Helper()
	ctx := context.Background()
	db, err := sqlite.Abrir(ctx, filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if _, err := seed.Cargar(ctx, db, seed.Embebido, ahora, true, log); err != nil {
		t.Fatal(err)
	}

	clave, _ := auth.GenerarClave(24)
	hash, _ := auth.HashClave(clave, bcrypt.MinCost)
	if _, err := db.SincronizarAdmin(ctx, "admin", hash, ahora); err != nil {
		t.Fatal(err)
	}
	secreto, _ := auth.NuevoToken()

	cfg := httpapi.ConfigPorDefecto()
	cfg.SecretoSesion = []byte(secreto)
	cfg.CostoBcryptDummy = bcrypt.MinCost
	cfg.ProxiesConfiables = []netip.Prefix{netip.MustParsePrefix("127.0.0.0/8")}
	cfg.ImagenesDir = t.TempDir()
	if ajustar != nil {
		ajustar(&cfg)
	}
	s := httpapi.New(cfg, db, log)
	s.SetReloj(func() time.Time { return ahora })
	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)
	return &entorno{srv: ts, clave: clave, db: db}
}

type resp struct {
	status int
	header http.Header
	body   []byte
}

func (r resp) json(t *testing.T, v any) {
	t.Helper()
	if err := json.Unmarshal(r.body, v); err != nil {
		t.Fatalf("json: %v\n%s", err, r.body)
	}
}

func (e *entorno) pedir(t *testing.T, metodo, ruta string, cuerpo any, hdr map[string]string, cookies ...*http.Cookie) resp {
	t.Helper()
	var rd io.Reader
	switch c := cuerpo.(type) {
	case nil:
	case string:
		rd = strings.NewReader(c)
	default:
		b, _ := json.Marshal(c)
		rd = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(metodo, e.srv.URL+ruta, rd)
	if rd != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if metodo != http.MethodGet {
		req.Header.Set("Origin", origen)
	}
	for k, v := range hdr {
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return resp{res.StatusCode, res.Header, b}
}

func (e *entorno) login(t *testing.T) *http.Cookie {
	t.Helper()
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	if r.status != 200 {
		t.Fatalf("login: %d %s", r.status, r.body)
	}
	for _, c := range (&http.Response{Header: r.header}).Cookies() {
		if c.Name == "dc_sesion" {
			if !c.HttpOnly || c.SameSite != http.SameSiteStrictMode {
				t.Fatalf("cookie sin HttpOnly/SameSite=Strict: %+v", c)
			}
			return c
		}
	}
	t.Fatal("login sin cookie")
	return nil
}

type productoResp struct {
	ID             int64    `json:"id"`
	Slug           string   `json:"slug"`
	Precio         int64    `json:"precio"`
	PrecioFinal    int64    `json:"precio_final"`
	PrecioMayor    *int64   `json:"precio_mayor"`
	MayorMinimo    *int     `json:"mayor_minimo"`
	Disponibilidad string   `json:"disponibilidad"`
	Badges         []string `json:"badges"`
}

type pagina struct {
	Items []json.RawMessage `json:"items"`
	Total int               `json:"total"`
}

func (e *entorno) producto(t *testing.T, slug string) productoResp {
	t.Helper()
	r := e.pedir(t, "GET", "/api/productos/"+slug, nil, nil)
	if r.status != 200 {
		t.Fatalf("producto %s: %d", slug, r.status)
	}
	var p productoResp
	r.json(t, &p)
	return p
}

func TestHealthzYCabeceras(t *testing.T) {
	e := nuevoEntorno(t, nil)
	r := e.pedir(t, "GET", "/healthz", nil, nil)
	if r.status != 200 {
		t.Fatalf("healthz %d", r.status)
	}
	for k, v := range map[string]string{
		"X-Content-Type-Options": "nosniff",
		"Referrer-Policy":        "no-referrer",
		"X-Frame-Options":        "DENY",
	} {
		if r.header.Get(k) != v {
			t.Errorf("%s = %q", k, r.header.Get(k))
		}
	}
	if !strings.Contains(r.header.Get("Content-Security-Policy"), "frame-ancestors 'none'") {
		t.Error("falta CSP frame-ancestors")
	}
	if r.header.Get("Strict-Transport-Security") != "" {
		t.Error("HSTS solo en producción")
	}
}

func TestCatalogoPublicoNoExponeStock(t *testing.T) {
	e := nuevoEntorno(t, nil)
	r := e.pedir(t, "GET", "/api/productos?page_size=48", nil, nil)
	if r.status != 200 {
		t.Fatalf("%d %s", r.status, r.body)
	}
	if bytes.Contains(r.body, []byte(`"stock"`)) || bytes.Contains(r.body, []byte(`"umbral_pocos"`)) {
		t.Fatal("la API pública no debe exponer stock")
	}
	var pg pagina
	r.json(t, &pg)
	if pg.Total != 37 {
		t.Fatalf("total productos = %d", pg.Total)
	}
	r = e.pedir(t, "GET", "/api/productos?categoria=navidad&oferta=true", nil, nil)
	r.json(t, &pg)
	if pg.Total == 0 {
		t.Fatal("debería haber ofertas en Navidad")
	}
	r = e.pedir(t, "GET", "/api/productos?page_size=500", nil, nil)
	if r.status != 400 {
		t.Fatalf("page_size fuera de rango: %d", r.status)
	}
	if e.producto(t, "repisa-plegable-metalica-con-ruedas-negra").Disponibilidad != "agotado" {
		t.Fatal("producto con stock 0 debe verse agotado")
	}
}

func TestBusquedaInyeccionSQLRespuestaSegura(t *testing.T) {
	e := nuevoEntorno(t, nil)
	for _, q := range []string{"' OR '1'='1", "'; DROP TABLE productos; --", "%", "_", `\%`, "1 UNION SELECT password_hash FROM admins"} {
		r := e.pedir(t, "GET", "/api/productos?q="+url.QueryEscape(q), nil, nil)
		if r.status != 200 {
			t.Fatalf("q=%q status %d", q, r.status)
		}
		var pg pagina
		r.json(t, &pg)
		if pg.Total != 0 {
			t.Errorf("q=%q devolvió %d resultados", q, pg.Total)
		}
		if bytes.Contains(r.body, []byte("$2")) || bytes.Contains(bytes.ToLower(r.body), []byte("sql")) {
			t.Errorf("q=%q filtró datos internos: %s", q, r.body)
		}
	}
	r := e.pedir(t, "GET", "/api/productos?q="+url.QueryEscape("telón"), nil, nil)
	var pg pagina
	r.json(t, &pg)
	if pg.Total != 1 {
		t.Fatalf("búsqueda normal 'telón' = %d", pg.Total)
	}
	r = e.pedir(t, "GET", "/api/productos?q="+strings.Repeat("a", 61), nil, nil)
	if r.status != 400 {
		t.Fatalf("búsqueda larga: %d", r.status)
	}
}

func TestPedidoRecalculaTotalIgnorandoAlCliente(t *testing.T) {
	e := nuevoEntorno(t, nil)
	squishy := e.producto(t, "squishy-viral-que-cambia-de-color-con-el-sol") // $4.000, x mayor $3.500 desde 12
	farol := e.producto(t, "farol-navideno-modelo-mas-vendido")              // $1.500, oferta $1.000

	body := map[string]any{
		"nombre": "José Muñoz", "telefono": "+56 9 1111 2222", "entrega": "retiro",
		"items": []map[string]any{
			{"producto_id": squishy.ID, "cantidad": 12, "precio_unitario": 1},
			{"producto_id": farol.ID, "cantidad": 2, "precio_unitario": 1},
		},
		"total": 1,
	}
	r := e.pedir(t, "POST", "/api/pedidos", body, nil)
	if r.status != 201 {
		t.Fatalf("crear pedido: %d %s", r.status, r.body)
	}
	var out struct {
		Numero      string `json:"numero"`
		Total       int64  `json:"total"`
		Estado      string `json:"estado"`
		Mensaje     string `json:"mensaje_whatsapp"`
		WhatsAppURL string `json:"whatsapp_url"`
	}
	r.json(t, &out)
	want := int64(12*3500 + 2*1000)
	if out.Total != want {
		t.Fatalf("total = %d, quería %d (el del cliente era 1)", out.Total, want)
	}
	if out.Estado != "nuevo" || out.Numero != "DC-0007" { // 6 pedidos DEMO en la semilla
		t.Fatalf("numero/estado = %s %s", out.Numero, out.Estado)
	}
	if !strings.Contains(out.Mensaje, "$44.000") || !strings.HasPrefix(out.WhatsAppURL, "https://wa.me/56921145368?text=") {
		t.Fatalf("mensaje: %s", out.Mensaje)
	}

	if !strings.HasPrefix(out.Mensaje, "*[PEDIDO DE PRUEBA · DEMO]*\nHola Donde Carlitos Emporio") {
		t.Fatalf("en modo demo el mensaje debe partir con el aviso: %q", out.Mensaje[:60])
	}

	// 11 unidades no alcanzan el mínimo por mayor (12): se cobra al detalle.
	body["items"] = []map[string]any{{"producto_id": squishy.ID, "cantidad": 11}}
	r = e.pedir(t, "POST", "/api/pedidos", body, nil)
	r.json(t, &out)
	if out.Total != 11*4000 {
		t.Fatalf("bajo el mínimo por mayor: total %d", out.Total)
	}

	// El total guardado en la base también es el recalculado.
	c := e.login(t)
	r = e.pedir(t, "GET", "/api/admin/pedidos?estado=nuevo", nil, nil, c)
	if !bytes.Contains(r.body, []byte(`"total":44000`)) {
		t.Fatalf("total guardado incorrecto: %s", r.body)
	}
}

func TestPedidoRechazos(t *testing.T) {
	e := nuevoEntorno(t, nil)
	agotado := e.producto(t, "repisa-plegable-metalica-con-ruedas-negra")
	ok := e.producto(t, "lazos-navidenos")
	base := func(id int64, cant int) map[string]any {
		return map[string]any{"nombre": "Ana Pérez", "telefono": "912345678", "entrega": "retiro",
			"items": []map[string]any{{"producto_id": id, "cantidad": cant}}}
	}
	if r := e.pedir(t, "POST", "/api/pedidos", base(agotado.ID, 1), nil); r.status != 409 || !bytes.Contains(r.body, []byte("agotado")) {
		t.Fatalf("agotado: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", base(999999, 1), nil); r.status != 409 {
		t.Fatalf("inexistente: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", base(ok.ID, 0), nil); r.status != 400 {
		t.Fatalf("cantidad 0: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", base(ok.ID, 10000), nil); r.status != 400 {
		t.Fatalf("cantidad enorme: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", base(ok.ID, 101), nil); r.status != 400 {
		t.Fatalf("más de 100 por línea: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", base(ok.ID, 60), nil); r.status != 409 || !bytes.Contains(r.body, []byte("stock_insuficiente")) {
		t.Fatalf("más que el stock: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", `{"nombre":"x","campo_raro":1}`, nil); r.status != 400 {
		t.Fatalf("campo desconocido: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", "nombre=x", map[string]string{"Content-Type": "application/x-www-form-urlencoded"}); r.status != 415 {
		t.Fatalf("content-type: %d", r.status)
	}
	grande := `{"nombre":"` + strings.Repeat("a", 20<<10) + `"}`
	if r := e.pedir(t, "POST", "/api/pedidos", grande, nil); r.status != 413 {
		t.Fatalf("cuerpo grande: %d", r.status)
	}
}

func TestOrigenObligatorioEnEscrituras(t *testing.T) {
	e := nuevoEntorno(t, nil)
	body := map[string]any{"usuario": "admin", "password": "x"}
	if r := e.pedir(t, "POST", "/api/admin/login", body, map[string]string{"Origin": ""}); r.status != 403 {
		t.Fatalf("sin Origin: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/admin/login", body, map[string]string{"Origin": "https://evil.example"}); r.status != 403 {
		t.Fatalf("Origin ajeno: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/pedidos", body, map[string]string{"Sec-Fetch-Site": "cross-site"}); r.status != 403 {
		t.Fatalf("cross-site: %d", r.status)
	}
}

func TestRutasAdminSinSesion401(t *testing.T) {
	e := nuevoEntorno(t, nil)
	rutas := []struct{ m, r string }{
		{"GET", "/api/admin/sesion"},
		{"GET", "/api/admin/pedidos"},
		{"GET", "/api/admin/pedidos/1"},
		{"PATCH", "/api/admin/pedidos/1/estado"},
		{"GET", "/api/admin/productos"},
		{"POST", "/api/admin/productos"},
		{"GET", "/api/admin/productos/1"},
		{"PUT", "/api/admin/productos/1"},
		{"DELETE", "/api/admin/productos/1"},
	}
	falsa := &http.Cookie{Name: "dc_sesion", Value: strings.Repeat("A", 43)}
	for _, rt := range rutas {
		for _, cookies := range [][]*http.Cookie{nil, {falsa}} {
			r := e.pedir(t, rt.m, rt.r, `{}`, nil, cookies...)
			if r.status != 401 {
				t.Errorf("%s %s (cookies=%d) = %d, quería 401", rt.m, rt.r, len(cookies), r.status)
			}
		}
	}
}

func TestLimiteIntentosLogin(t *testing.T) {
	e := nuevoEntorno(t, nil)
	malo := map[string]string{"usuario": "admin", "password": "incorrecta"}
	for i := 0; i < 5; i++ {
		if r := e.pedir(t, "POST", "/api/admin/login", malo, nil); r.status != 401 {
			t.Fatalf("intento %d: %d", i+1, r.status)
		}
	}
	// Bloqueado: ni la clave correcta entra.
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	if r.status != 429 || r.header.Get("Retry-After") == "" {
		t.Fatalf("tras 5 fallos: %d", r.status)
	}
	// Otra IP (vía proxy confiable) no queda bloqueada.
	r = e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, map[string]string{"X-Forwarded-For": "203.0.113.9"})
	if r.status != 200 {
		t.Fatalf("otra IP: %d", r.status)
	}
	// Usuario inexistente responde igual que clave mala.
	r = e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "nadie", "password": "x"}, map[string]string{"X-Forwarded-For": "203.0.113.10"})
	if r.status != 401 || !bytes.Contains(r.body, []byte("credenciales_invalidas")) {
		t.Fatalf("usuario inexistente: %d %s", r.status, r.body)
	}
}

func TestLimitePedidosPorIP(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.PedidosPorVentana = 3 })
	// Cuerpos inválidos no consumen cupo (A1).
	for i := 0; i < 15; i++ {
		if r := e.pedir(t, "POST", "/api/pedidos", `{}`, nil); r.status != 400 {
			t.Fatalf("inválido %d: %d", i, r.status)
		}
	}
	lazos := e.producto(t, "lazos-navidenos")
	for i := 0; i < 3; i++ {
		if r := e.pedir(t, "POST", "/api/pedidos", pedidoValido(lazos.ID, i), nil); r.status != 201 {
			t.Fatalf("pedido %d: %d %s", i, r.status, r.body)
		}
	}
	if r := e.pedir(t, "POST", "/api/pedidos", pedidoValido(lazos.ID, 9), nil); r.status != 429 || r.header.Get("Retry-After") == "" {
		t.Fatalf("4º pedido: %d", r.status)
	}
}

func TestFlujoAdmin(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)

	// Lista con filtro por estado.
	r := e.pedir(t, "GET", "/api/admin/pedidos?estado=nuevo", nil, nil, c)
	var pg struct {
		Items []struct {
			ID                     int64    `json:"id"`
			Numero                 string   `json:"numero"`
			EsDemo                 bool     `json:"es_demo"`
			TransicionesPermitidas []string `json:"transiciones_permitidas"`
		} `json:"items"`
		Total int `json:"total"`
	}
	r.json(t, &pg)
	if pg.Total != 1 || !pg.Items[0].EsDemo {
		t.Fatalf("pedidos nuevos demo: %s", r.body)
	}
	id := pg.Items[0].ID
	ruta := "/api/admin/pedidos/" + itoa(id) + "/estado"

	if r := e.pedir(t, "PATCH", ruta, map[string]string{"estado": "entregado"}, nil, c); r.status != 409 {
		t.Fatalf("transición inválida: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "PATCH", ruta, map[string]string{"estado": "confirmado", "nota": "Confirmado por WhatsApp"}, nil, c); r.status != 200 || !bytes.Contains(r.body, []byte(`"historial"`)) {
		t.Fatalf("confirmar: %d %s", r.status, r.body)
	}

	// Editar producto: precio, oferta, stock y activo.
	r = e.pedir(t, "GET", "/api/admin/productos?q=lazos", nil, nil, c)
	var prods struct {
		Items []map[string]any `json:"items"`
	}
	r.json(t, &prods)
	if len(prods.Items) != 1 || prods.Items[0]["stock"] == nil {
		t.Fatalf("admin debe ver stock: %s", r.body)
	}
	p := prods.Items[0]
	cat := p["categoria"].(map[string]any)
	edit := map[string]any{
		"slug": p["slug"], "codigo": p["codigo"], "nombre": p["nombre"], "descripcion": p["descripcion"],
		"categoria_id": cat["id"], "precio": 1200, "precio_oferta": 900, "etiqueta_oferta": "remate",
		"stock": 0, "activo": false, "temporada": "navidad",
	}
	r = e.pedir(t, "PUT", "/api/admin/productos/"+itoa(int64(p["id"].(float64))), edit, nil, c)
	if r.status != 200 || !bytes.Contains(r.body, []byte(`"precio_oferta":900`)) {
		t.Fatalf("editar: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "GET", "/api/productos/"+p["slug"].(string), nil, nil); r.status != 404 {
		t.Fatalf("inactivo debe ocultarse del público: %d", r.status)
	}

	// Crear con datos inválidos y válidos; luego borrar.
	if r := e.pedir(t, "POST", "/api/admin/productos", map[string]any{"codigo": "X1", "nombre": "A", "categoria_id": 1, "precio": -5}, nil, c); r.status != 400 {
		t.Fatalf("crear inválido: %d", r.status)
	}
	nuevo := map[string]any{"codigo": "P900", "nombre": "Guirnalda LED [EJEMPLO]", "categoria_id": cat["id"], "precio": 2990, "stock": 8, "es_ejemplo": true}
	r = e.pedir(t, "POST", "/api/admin/productos", nuevo, nil, c)
	if r.status != 201 {
		t.Fatalf("crear: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "POST", "/api/admin/productos", nuevo, nil, c); r.status != 409 {
		t.Fatalf("duplicado: %d", r.status)
	}
	var creado struct {
		ID int64 `json:"id"`
	}
	r.json(t, &creado)
	if r := e.pedir(t, "DELETE", "/api/admin/productos/"+itoa(creado.ID), nil, nil, c); r.status != 204 {
		t.Fatalf("borrar: %d", r.status)
	}

	// Logout invalida la sesión.
	if r := e.pedir(t, "POST", "/api/admin/logout", nil, nil, c); r.status != 204 {
		t.Fatalf("logout: %d", r.status)
	}
	if r := e.pedir(t, "GET", "/api/admin/sesion", nil, nil, c); r.status != 401 {
		t.Fatalf("sesión tras logout: %d", r.status)
	}
}

func TestCookieSecureEnProduccion(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.Produccion = true })
	r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "admin", "password": e.clave}, nil)
	sc := r.header.Get("Set-Cookie")
	if !strings.HasPrefix(sc, "__Host-dc_sesion=") || !strings.Contains(sc, "Secure") || !strings.Contains(sc, "HttpOnly") || !strings.Contains(sc, "SameSite=Strict") {
		t.Fatalf("cookie en producción: %s", sc)
	}
	if r.header.Get("Strict-Transport-Security") == "" {
		t.Fatal("falta HSTS en producción")
	}
}

func itoa(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
