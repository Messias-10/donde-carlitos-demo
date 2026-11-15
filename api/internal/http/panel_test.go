package httpapi_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	httpapi "github.com/Messias-10/donde-carlitos-demo/api/internal/http"
)

type pedidoResp struct {
	ID                     int64    `json:"id"`
	Numero                 string   `json:"numero"`
	Canal                  string   `json:"canal"`
	Estado                 string   `json:"estado"`
	EsDemo                 bool     `json:"es_demo"`
	Nombre                 string   `json:"nombre"`
	Telefono               string   `json:"telefono"`
	Subtotal               int64    `json:"subtotal"`
	Total                  int64    `json:"total"`
	TotalManual            bool     `json:"total_manual"`
	StockDescontado        bool     `json:"stock_descontado"`
	SeguimientoURL         string   `json:"seguimiento_url"`
	TransicionesPermitidas []string `json:"transiciones_permitidas"`
	WhatsAppClienteURL     string   `json:"whatsapp_cliente_url"`
	RegistradoPor          *struct {
		Nombre string `json:"nombre"`
	} `json:"registrado_por"`
	Historial []struct {
		EstadoNuevo string `json:"estado_nuevo"`
		Usuario     *struct {
			Nombre string `json:"nombre"`
		} `json:"usuario"`
	} `json:"historial"`
	Items []struct {
		ProductoID     *int64 `json:"producto_id"`
		PrecioUnitario *int64 `json:"precio_unitario"`
		TipoPrecio     string `json:"tipo_precio"`
	} `json:"items"`
}

func (e *entorno) buscarPedido(t *testing.T, c *http.Cookie, numero string) pedidoResp {
	t.Helper()
	var pg struct {
		Items []pedidoResp `json:"items"`
	}
	e.pedir(t, "GET", "/api/admin/pedidos?q="+numero, nil, nil, c).json(t, &pg)
	if len(pg.Items) != 1 {
		t.Fatalf("buscar %s: %d resultados", numero, len(pg.Items))
	}
	return pg.Items[0]
}

func TestHoyConLaSemillaDelGuion(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	r := e.pedir(t, "GET", "/api/admin/hoy", nil, nil, c)
	var hoy struct {
		Fecha               string         `json:"fecha"`
		Contadores          map[string]int `json:"contadores"`
		EsperandoPagoMas24h []pedidoResp   `json:"esperando_pago_mas_24h"`
		PagadosSinDespachar []pedidoResp   `json:"pagados_sin_despachar"`
		ListosRetiro        []pedidoResp   `json:"listos_retiro"`
	}
	r.json(t, &hoy)
	nums := func(ps []pedidoResp) string {
		var s []string
		for _, p := range ps {
			s = append(s, p.Numero)
		}
		return strings.Join(s, ",")
	}
	if hoy.Fecha != "2026-11-14" || nums(hoy.EsperandoPagoMas24h) != "DC-0109,DC-0101" ||
		nums(hoy.PagadosSinDespachar) != "DC-0102" || nums(hoy.ListosRetiro) != "DC-0104" ||
		hoy.Contadores["creados_hoy_regiones"] != 2 || hoy.Contadores["creados_hoy_santiago"] != 0 {
		t.Fatalf("Hoy: %s", r.body)
	}
}

func TestRegistrarPedidoYFlujo(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	squishy := e.producto(t, "squishy-viral-que-cambia-de-color-con-el-sol")
	body := map[string]any{
		"canal": "regiones", "nombre": "Prueba Panel", "telefono": "+56 9 1111 2222", "comuna": "Talca", "region": "maule",
		"items": []map[string]any{
			{"producto_id": squishy.ID, "cantidad": 12},
			{"descripcion": "Bolsas de regalo surtidas", "cantidad": 2, "precio_unitario": 500},
		},
		"descuento": 2000, "costo_envio": 5000,
	}
	r := e.pedir(t, "POST", "/api/admin/pedidos", body, nil, c)
	if r.status != 201 {
		t.Fatalf("crear: %d %s", r.status, r.body)
	}
	var p pedidoResp
	r.json(t, &p)
	if p.Numero != "DC-0110" || p.Total != 12*3500+1000-2000+5000 || p.Estado != "cotizado" || p.RegistradoPor == nil ||
		p.Items[0].TipoPrecio != "mayor" || p.Items[1].ProductoID != nil || p.EsDemo {
		t.Fatalf("pedido creado: %s", r.body)
	}
	if !regexp.MustCompile(`^http://localhost:3100/seguimiento/[A-Za-z0-9_-]{22}$`).MatchString(p.SeguimientoURL) ||
		strings.Contains(p.SeguimientoURL, "DC-") || p.WhatsAppClienteURL != "https://wa.me/56911112222" {
		t.Fatalf("links: %s %s", p.SeguimientoURL, p.WhatsAppClienteURL)
	}
	ruta := fmt.Sprintf("/api/admin/pedidos/%d/estado", p.ID)

	// Pagado exige confirmar que la plata está en la cuenta.
	pago := map[string]any{"medio": "transferencia", "monto": p.Total}
	if r := e.pedir(t, "PATCH", ruta, map[string]any{"estado": "pagado", "pago": pago}, nil, c); r.status != 400 || !bytes.Contains(r.body, []byte("confirmo_pago_recibido")) {
		t.Fatalf("pagado sin confirmar: %d %s", r.status, r.body)
	}
	r = e.pedir(t, "PATCH", ruta, map[string]any{"estado": "pagado", "pago": pago, "confirmo_pago_recibido": true}, nil, c)
	r.json(t, &p)
	if r.status != 200 || !p.StockDescontado || p.Historial[len(p.Historial)-1].Usuario == nil {
		t.Fatalf("pagado: %d %s", r.status, r.body)
	}
	// Ya pagado: no se edita.
	if r := e.pedir(t, "PUT", fmt.Sprintf("/api/admin/pedidos/%d", p.ID), body, nil, c); r.status != 409 || !bytes.Contains(r.body, []byte("no_editable")) {
		t.Fatalf("editar pagado: %d", r.status)
	}
	e.pedir(t, "PATCH", ruta, map[string]any{"estado": "preparando"}, nil, c)
	if r := e.pedir(t, "PATCH", ruta, map[string]any{"estado": "despachado"}, nil, c); r.status != 400 ||
		!bytes.Contains(r.body, []byte("agencia")) || !bytes.Contains(r.body, []byte("n_seguimiento")) {
		t.Fatalf("despachado sin datos: %d %s", r.status, r.body)
	}
	r = e.pedir(t, "PATCH", ruta, map[string]any{"estado": "despachado", "agencia": "starken", "n_seguimiento": "000123456789"}, nil, c)
	if r.status != 200 {
		t.Fatalf("despachado: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "PATCH", ruta, map[string]any{"estado": "cotizado"}, nil, c); r.status != 409 {
		t.Fatalf("transición inválida: %d", r.status)
	}

	// Mensajes: el que toca es el de despacho, hacia el teléfono del cliente.
	r = e.pedir(t, "GET", fmt.Sprintf("/api/admin/pedidos/%d/mensajes", p.ID), nil, nil, c)
	var ms struct {
		Mensajes []struct {
			Clave       string `json:"clave"`
			Texto       string `json:"texto"`
			WhatsAppURL string `json:"whatsapp_url"`
			Sugerido    bool   `json:"sugerido"`
		} `json:"mensajes"`
	}
	r.json(t, &ms)
	var sug string
	for _, m := range ms.Mensajes {
		if m.Sugerido {
			sug = m.Clave
			if !strings.HasPrefix(m.WhatsAppURL, "https://wa.me/56911112222?text=") || !strings.Contains(m.Texto, "*N.º de seguimiento:* 000123456789") {
				t.Fatalf("mensaje despacho: %+v", m)
			}
		}
	}
	if sug != "despachado" {
		t.Fatalf("sugerido = %q", sug)
	}

	// Línea sin precio exige total escrito.
	body["items"] = []map[string]any{{"descripcion": "Algo sin precio", "cantidad": 1}}
	if r := e.pedir(t, "POST", "/api/admin/pedidos", body, nil, c); r.status != 400 || !bytes.Contains(r.body, []byte(`"total"`)) {
		t.Fatalf("sin precio y sin total: %d %s", r.status, r.body)
	}
	body["total"] = 9990
	if r := e.pedir(t, "POST", "/api/admin/pedidos", body, nil, c); r.status != 201 || !bytes.Contains(r.body, []byte(`"total_manual":true`)) {
		t.Fatalf("total manual: %d %s", r.status, r.body)
	}
}

func TestMensajesDemoSinNumero(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	p := e.buscarPedido(t, c, "DC-0101")
	if !p.EsDemo || p.WhatsAppClienteURL != "https://wa.me/" {
		t.Fatalf("demo: %+v", p)
	}
	r := e.pedir(t, "GET", fmt.Sprintf("/api/admin/pedidos/%d/mensajes", p.ID), nil, nil, c)
	if bytes.Contains(r.body, []byte("wa.me/569")) || !bytes.Contains(r.body, []byte(`"clave":"recordatorio_pago","titulo":"Recordar pago"`)) {
		t.Fatalf("mensajes demo: %s", r.body)
	}
	if !bytes.Contains(r.body, []byte(`"sugerido":true`)) || bytes.Contains(r.body, []byte("datos_transferencia")) {
		t.Fatalf("sin DATOS_TRANSFERENCIA no se ofrece d2: %s", r.body)
	}
}

func TestSeguimientoPublicoSinDatosPersonales(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.SeguimientoPorIP = 3 })
	c := e.login(t)
	p := e.buscarPedido(t, c, "DC-0103")
	token := p.SeguimientoURL[strings.LastIndex(p.SeguimientoURL, "/")+1:]
	r := e.pedir(t, "GET", "/api/seguimiento/"+token, nil, nil)
	if r.status != 200 || r.header.Get("X-Robots-Tag") == "" || !bytes.Contains(r.body, []byte(`"n_seguimiento":"000123456789"`)) ||
		!bytes.Contains(r.body, []byte(`"agencia_nombre":"Starken"`)) {
		t.Fatalf("seguimiento: %d %s", r.status, r.body)
	}
	for _, prohibido := range []string{"Jorge", "+569", "56900000103", "La Serena", "total", "items", "nota", "registrado"} {
		if bytes.Contains(r.body, []byte(prohibido)) {
			t.Errorf("el seguimiento filtra %q: %s", prohibido, r.body)
		}
	}
	for _, tk := range []string{"DC-0103", strings.Repeat("a", 22), "../../etc"} {
		if r := e.pedir(t, "GET", "/api/seguimiento/"+tk, nil, nil); r.status != 404 {
			t.Errorf("token %q: %d", tk, r.status)
		}
	}
	if r := e.pedir(t, "GET", "/api/seguimiento/"+token, nil, nil); r.status != 429 {
		t.Fatalf("límite por IP: %d", r.status)
	}
}

func TestRolesYUsuarios(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	r := e.pedir(t, "POST", "/api/admin/usuarios", map[string]string{"usuario": "vane", "nombre": "Vane (regiones)", "rol": "vendedor"}, nil, c)
	var creado struct {
		Usuario struct {
			ID int64 `json:"id"`
		} `json:"usuario"`
		ClaveTemporal string `json:"clave_temporal"`
	}
	r.json(t, &creado)
	if r.status != 201 || len(creado.ClaveTemporal) != 16 {
		t.Fatalf("crear vendedor: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "POST", "/api/admin/usuarios", map[string]string{"usuario": "vane", "nombre": "Otra", "rol": "vendedor"}, nil, c); r.status != 409 {
		t.Fatalf("duplicado: %d", r.status)
	}
	lista := e.pedir(t, "GET", "/api/admin/usuarios", nil, nil, c)
	if bytes.Contains(lista.body, []byte(creado.ClaveTemporal)) || bytes.Contains(lista.body, []byte("password")) || !bytes.Contains(lista.body, []byte("Regiones 1")) {
		t.Fatalf("la lista no muestra claves ni hashes: %s", lista.body)
	}

	v := e.loginComo(t, "vane", creado.ClaveTemporal)
	if r := e.pedir(t, "GET", "/api/admin/pedidos", nil, nil, v); r.status != 200 {
		t.Fatalf("vendedor ve pedidos: %d", r.status)
	}
	if r := e.pedir(t, "GET", "/api/admin/productos", nil, nil, v); r.status != 200 {
		t.Fatalf("vendedor ve productos: %d", r.status)
	}
	for _, rt := range []struct{ m, r string }{
		{"POST", "/api/admin/productos"}, {"PUT", "/api/admin/productos/1"}, {"DELETE", "/api/admin/productos/1"},
		{"GET", "/api/admin/usuarios"}, {"POST", "/api/admin/usuarios"}, {"GET", "/api/admin/resumen"}, {"GET", "/api/admin/pedidos.csv"},
	} {
		if r := e.pedir(t, rt.m, rt.r, `{}`, nil, v); r.status != 403 || !bytes.Contains(r.body, []byte("sin_permiso")) {
			t.Errorf("vendedor %s %s: %d", rt.m, rt.r, r.status)
		}
	}

	// Desactivar: su sesión se corta y no vuelve a entrar.
	ruta := fmt.Sprintf("/api/admin/usuarios/%d", creado.Usuario.ID)
	if r := e.pedir(t, "PATCH", ruta, map[string]any{"activo": false}, nil, c); r.status != 200 {
		t.Fatalf("desactivar: %d", r.status)
	}
	if r := e.pedir(t, "GET", "/api/admin/sesion", nil, nil, v); r.status != 401 {
		t.Fatalf("sesión del desactivado: %d", r.status)
	}
	if r := e.pedir(t, "POST", "/api/admin/login", map[string]string{"usuario": "vane", "password": creado.ClaveTemporal}, nil); r.status != 401 {
		t.Fatalf("login desactivado: %d", r.status)
	}
	// Un admin no se desactiva a sí mismo.
	var ses struct {
		ID int64 `json:"id"`
	}
	e.pedir(t, "GET", "/api/admin/sesion", nil, nil, c).json(t, &ses)
	if r := e.pedir(t, "PATCH", fmt.Sprintf("/api/admin/usuarios/%d", ses.ID), map[string]any{"activo": false}, nil, c); r.status != 409 {
		t.Fatalf("auto-desactivar: %d", r.status)
	}
	// Clave nueva (se muestra una vez) y reactivar.
	e.pedir(t, "PATCH", ruta, map[string]any{"activo": true}, nil, c)
	r = e.pedir(t, "POST", ruta+"/clave", nil, nil, c)
	var nueva struct {
		ClaveTemporal string `json:"clave_temporal"`
	}
	r.json(t, &nueva)
	if nueva.ClaveTemporal == "" || nueva.ClaveTemporal == creado.ClaveTemporal {
		t.Fatal("clave regenerada")
	}
	e.loginComo(t, "vane", nueva.ClaveTemporal)
}

func TestCSVProtegidoContraInyeccion(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	body := map[string]any{"canal": "regiones", "nombre": "Ana", "telefono": "912345678", "comuna": "-1 de prueba",
		"items": []map[string]any{{"descripcion": "Algo", "cantidad": 1, "precio_unitario": 1000}}}
	if r := e.pedir(t, "POST", "/api/admin/pedidos", body, nil, c); r.status != 201 {
		t.Fatalf("crear: %d %s", r.status, r.body)
	}
	r := e.pedir(t, "GET", "/api/admin/pedidos.csv", nil, nil, c)
	if r.status != 200 || !strings.HasPrefix(r.header.Get("Content-Type"), "text/csv") ||
		!strings.Contains(r.header.Get("Content-Disposition"), "attachment") || !bytes.HasPrefix(r.body, []byte("\xEF\xBB\xBF")) {
		t.Fatalf("csv: %d %v", r.status, r.header)
	}
	cr := csv.NewReader(bytes.NewReader(r.body[3:]))
	cr.Comma = ';'
	filas, err := cr.ReadAll()
	if err != nil || len(filas) != 22 || filas[0][0] != "numero" {
		t.Fatalf("filas %d: %v", len(filas), err)
	}
	comuna := -1
	for i, h := range filas[0] {
		if h == "comuna" {
			comuna = i
		}
	}
	if filas[1][comuna] != "'-1 de prueba" {
		t.Fatalf("celda sin escapar: %q", filas[1][comuna])
	}
	for _, f := range filas[1:] {
		for _, celda := range f {
			if celda != "" && strings.ContainsRune("=+@", rune(celda[0])) {
				t.Fatalf("celda peligrosa: %q", celda)
			}
		}
	}
}

func TestClientesResumenYConfigPublica(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	r := e.pedir(t, "GET", "/api/admin/clientes?telefono=9%200000%200101", nil, nil, c)
	if r.status != 200 || !bytes.Contains(r.body, []byte(`"total_pedidos":1`)) || !bytes.Contains(r.body, []byte(`"nombre":"Pedro Soto Demo"`)) {
		t.Fatalf("cliente: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "GET", "/api/admin/clientes?telefono=abc", nil, nil, c); r.status != 400 {
		t.Fatalf("teléfono inválido: %d", r.status)
	}
	r = e.pedir(t, "GET", "/api/admin/resumen", nil, nil, c)
	var res struct {
		Totales  map[string]int64 `json:"totales"`
		PorCanal []struct {
			Clave   string `json:"clave"`
			Pedidos int    `json:"pedidos"`
		} `json:"por_canal"`
		PorAgencia []struct {
			Clave string `json:"clave"`
		} `json:"por_agencia"`
	}
	r.json(t, &res)
	if r.status != 200 || res.Totales["pedidos"] != 18 || res.Totales["cancelados"] != 2 || len(res.PorCanal) != 2 || res.PorCanal[0].Clave != "santiago" {
		t.Fatalf("resumen: %s", r.body)
	}
	if r := e.pedir(t, "GET", "/api/admin/resumen?desde=2026-11-20&hasta=2026-11-01", nil, nil, c); r.status != 400 {
		t.Fatalf("rango inválido: %d", r.status)
	}
	r = e.pedir(t, "GET", "/api/config-publica", nil, nil)
	if !bytes.Contains(r.body, []byte(`"whatsapp_santiago":"56921145368"`)) || !bytes.Contains(r.body, []byte(`"whatsapp_regiones":null`)) ||
		!bytes.Contains(r.body, []byte(`"whatsapp_regiones_confirmado":false`)) || bytes.Contains(r.body, []byte("transfer")) {
		t.Fatalf("config pública: %s", r.body)
	}
}

func TestBusquedaPedidos(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	contar := func(q string) int {
		var pg pagina
		e.pedir(t, "GET", "/api/admin/pedidos?"+q, nil, nil, c).json(t, &pg)
		return pg.Total
	}
	for q, want := range map[string]int{
		"q=DC-0104":                       1,
		"q=camila":                        1,
		"q=diaz":                          1, // sin tildes
		"q=%2B56%209%200000%200102":       1,
		"canal=santiago":                  6,
		"estado=esperando_pago":           3,
		"canal=regiones&estado=cancelado": 2,
		"desde=2026-11-14":                2,
	} {
		if got := contar(q); got != want {
			t.Errorf("%s: %d, quería %d", q, got, want)
		}
	}
	if r := e.pedir(t, "GET", "/api/admin/pedidos?canal=marte", nil, nil, c); r.status != 400 {
		t.Fatalf("canal inválido: %d", r.status)
	}
}
