package httpapi_test

// Correcciones de QA ronda 2 y seguridad ronda 2 (contrato 2.1.0).

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
	"time"

	httpapi "github.com/Messias-10/donde-carlitos-demo/api/internal/http"
)

// B-1a: de pagado se puede despachar directo (transición directa, sin
// "preparando" implícito). M-7: un despachado no se cancela.
func TestPagadoADespachadoDirecto(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	p := e.buscarPedido(t, c, "DC-0102") // pagado en la semilla
	ruta := fmt.Sprintf("/api/admin/pedidos/%d/estado", p.ID)
	if !strings.Contains(strings.Join(p.TransicionesPermitidas, ","), "despachado") {
		t.Fatalf("transiciones de pagado: %v", p.TransicionesPermitidas)
	}
	r := e.pedir(t, "PATCH", ruta, map[string]any{"estado": "despachado", "n_seguimiento": "VT-998877"}, nil, c)
	r.json(t, &p)
	if r.status != 200 || p.Estado != "despachado" || p.Historial[len(p.Historial)-2].EstadoNuevo != "pagado" {
		t.Fatalf("despacho directo: %d %s", r.status, r.body)
	}
	if strings.Contains(strings.Join(p.TransicionesPermitidas, ","), "cancelado") {
		t.Fatal("despachado no se cancela")
	}
	if r := e.pedir(t, "PATCH", ruta, map[string]any{"estado": "cancelado"}, nil, c); r.status != 409 {
		t.Fatalf("cancelar despachado: %d", r.status)
	}
}

// B-3: con MODO_DEMO, todo pedido nuevo es demo (wa.me sin número + línea demo).
func TestModoDemo(t *testing.T) {
	body := map[string]any{"canal": "santiago", "telefono": "912345678",
		"items": []map[string]any{{"descripcion": "Algo", "cantidad": 1, "precio_unitario": 1000}}}
	for _, demo := range []bool{true, false} {
		e := nuevoEntorno(t, func(c *httpapi.Config) { c.ModoDemo = demo })
		c := e.login(t)
		r := e.pedir(t, "POST", "/api/admin/pedidos", body, nil, c)
		var p pedidoResp
		r.json(t, &p)
		if p.EsDemo != demo || (demo && p.WhatsAppClienteURL != "https://wa.me/") {
			t.Fatalf("modo demo=%v: %s", demo, r.body)
		}
		m := e.pedir(t, "GET", fmt.Sprintf("/api/admin/pedidos/%d/mensajes", p.ID), nil, nil, c)
		if demo != bytes.Contains(m.body, []byte("PRUEBA · DEMO")) || demo == bytes.Contains(m.body, []byte("wa.me/56912345678")) {
			t.Fatalf("mensajes modo demo=%v: %s", demo, m.body)
		}
	}
}

// I-1: stock no estricto por defecto: se acepta, queda en 0 y avisa. Estricto: 409.
func TestStockNoEstrictoConAviso(t *testing.T) {
	for _, estricto := range []bool{false, true} {
		e := nuevoEntorno(t, func(c *httpapi.Config) { c.StockEstricto = estricto })
		c := e.login(t)
		esfera := e.producto(t, "esfera-navidena-con-luces")
		r := e.pedir(t, "POST", "/api/admin/pedidos", map[string]any{"canal": "santiago", "telefono": "912345678", "forma_pago": "transferencia",
			"items": []map[string]any{{"producto_id": esfera.ID, "cantidad": 9999}}}, nil, c)
		var p pedidoResp
		r.json(t, &p)
		r = e.pedir(t, "PATCH", fmt.Sprintf("/api/admin/pedidos/%d/estado", p.ID), map[string]any{"estado": "pagado",
			"pago": map[string]any{"medio": "transferencia", "monto": p.Total}, "confirmo_pago_recibido": true}, nil, c)
		if estricto {
			if r.status != 409 || !bytes.Contains(r.body, []byte("stock_insuficiente")) {
				t.Fatalf("estricto: %d %s", r.status, r.body)
			}
			continue
		}
		if r.status != 200 || !bytes.Contains(r.body, []byte(`"aviso_stock":{"mensaje"`)) || !bytes.Contains(r.body, []byte("Esfera navideña")) {
			t.Fatalf("no estricto: %d %s", r.status, r.body)
		}
		if d := e.producto(t, "esfera-navidena-con-luces").Disponibilidad; d != "agotado" {
			t.Fatalf("el stock queda en 0: %s", d)
		}
	}
}

// B-2 e I-2: promo por comentario y vitrina elegida por el dueño.
func TestPromoYVitrina(t *testing.T) {
	e := nuevoEntorno(t, nil)
	r := e.pedir(t, "GET", "/api/productos/esfera-navidena-con-luces", nil, nil)
	if !bytes.Contains(r.body, []byte(`"precio_es_promo":true`)) || !bytes.Contains(r.body, []byte(`yo quiero mi promoción`)) {
		t.Fatalf("promo: %s", r.body)
	}
	r = e.pedir(t, "GET", "/api/vitrina", nil, nil)
	var v struct {
		Items []struct {
			Codigo       string `json:"codigo"`
			OrdenVitrina int    `json:"orden_vitrina"`
		} `json:"items"`
		Elegidos bool `json:"elegidos"`
	}
	r.json(t, &v)
	var codigos []string
	for _, it := range v.Items {
		codigos = append(codigos, it.Codigo)
	}
	if !v.Elegidos || strings.Join(codigos, ",") != "P012,P010,P016,P001,P005,P020,P002,P027" {
		t.Fatalf("vitrina: %v %v", v.Elegidos, codigos)
	}
	// Una novena en vitrina: 409.
	c := e.login(t)
	var adm struct {
		Items []map[string]any `json:"items"`
	}
	e.pedir(t, "GET", "/api/admin/productos?q=lazos", nil, nil, c).json(t, &adm)
	p := adm.Items[0]
	put := map[string]any{"slug": p["slug"], "codigo": p["codigo"], "nombre": p["nombre"], "categoria_id": p["categoria"].(map[string]any)["id"],
		"precio": p["precio"], "stock": p["stock"], "activo": true, "en_vitrina": true, "orden_vitrina": 9}
	if r := e.pedir(t, "PUT", fmt.Sprintf("/api/admin/productos/%d", int64(p["id"].(float64))), put, nil, c); r.status != 409 || !bytes.Contains(r.body, []byte("vitrina_llena")) {
		t.Fatalf("vitrina llena: %d %s", r.status, r.body)
	}
	// M-11: unidad por mayor sin paréntesis.
	put["en_vitrina"], put["orden_vitrina"] = false, nil
	put["precio_mayor"], put["mayor_minimo"], put["unidad_mayor"] = 800, 12, "caja (12)"
	if r := e.pedir(t, "PUT", fmt.Sprintf("/api/admin/productos/%d", int64(p["id"].(float64))), put, nil, c); r.status != 400 || !bytes.Contains(r.body, []byte("unidad_mayor")) {
		t.Fatalf("unidad con paréntesis: %d %s", r.status, r.body)
	}
	// M-2: al buscar "squishy", primero los que empiezan con esa palabra.
	var lista struct {
		Items []struct {
			Nombre string `json:"nombre"`
		} `json:"items"`
	}
	e.pedir(t, "GET", "/api/productos?q=squishy", nil, nil).json(t, &lista)
	if len(lista.Items) < 3 || !strings.HasPrefix(lista.Items[0].Nombre, "Squishy") || !strings.HasPrefix(lista.Items[1].Nombre, "Squishy") {
		t.Fatalf("orden de búsqueda: %+v", lista.Items)
	}
}

// B-1b e I-3: resumen con comunas, personas y "vendido" solo con pago.
func TestResumenComunaUsuarioYVendido(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	r := e.pedir(t, "GET", "/api/admin/resumen", nil, nil, c)
	var res struct {
		Totales   map[string]int64 `json:"totales"`
		PorComuna []struct {
			Comuna  string `json:"comuna"`
			Pedidos int    `json:"pedidos"`
			Monto   int64  `json:"monto"`
		} `json:"por_comuna"`
		PorUsuario []struct {
			Usuario     string `json:"usuario"`
			Registrados int    `json:"registrados"`
			Despachados int    `json:"despachados"`
		} `json:"por_usuario"`
	}
	r.json(t, &res)
	// 18 pedidos no cancelados; los de la semilla con pago registrado son 13 (DC-0090..0100 menos 0095, y 0102, 0103, 0107).
	if res.Totales["pedidos"] != 18 || res.Totales["cancelados"] != 2 || res.Totales["pedidos_vendidos"] != 13 ||
		res.Totales["vendido"] == 0 || res.Totales["ticket_promedio"] != res.Totales["vendido"]/13 || res.Totales["monto"] != res.Totales["vendido"] {
		t.Fatalf("totales: %v", res.Totales)
	}
	if len(res.PorComuna) == 0 || res.PorComuna[0].Pedidos < res.PorComuna[len(res.PorComuna)-1].Pedidos {
		t.Fatalf("por comuna: %+v", res.PorComuna)
	}
	var reg1 *struct {
		Usuario     string `json:"usuario"`
		Registrados int    `json:"registrados"`
		Despachados int    `json:"despachados"`
	}
	for i := range res.PorUsuario {
		if res.PorUsuario[i].Usuario == "Regiones 1" {
			reg1 = &res.PorUsuario[i]
		}
	}
	if reg1 == nil || reg1.Registrados != 14 || reg1.Despachados != 7 {
		t.Fatalf("por usuario: %+v", res.PorUsuario)
	}
}

// B-1c: ajustes de solo lectura, sin el texto de los datos bancarios.
func TestAjustes(t *testing.T) {
	e := nuevoEntorno(t, func(c *httpapi.Config) { c.Tienda.DatosTransferencia = "Banco secreto 123"; c.ModoDemo = true })
	c := e.login(t)
	r := e.pedir(t, "GET", "/api/admin/ajustes", nil, nil, c)
	if r.status != 200 || bytes.Contains(r.body, []byte("Banco secreto")) || !bytes.Contains(r.body, []byte(`"datos_transferencia_configurados":true`)) ||
		!bytes.Contains(r.body, []byte(`"modo_demo":true`)) || !bytes.Contains(r.body, []byte(`"whatsapp_regiones":null`)) {
		t.Fatalf("ajustes: %d %s", r.status, r.body)
	}
}

// I-4: CSV con las columnas de P5, hora de Chile, nombres y teléfono en dígitos.
func TestCSVColumnasP5(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	r := e.pedir(t, "GET", "/api/admin/pedidos.csv?q=DC-0103", nil, nil, c)
	cr := csv.NewReader(bytes.NewReader(r.body[3:]))
	cr.Comma = ';'
	filas, err := cr.ReadAll()
	if err != nil || len(filas) != 2 {
		t.Fatalf("csv: %v %d", err, len(filas))
	}
	want := "numero;fecha_registro;canal;entrega;estado;nombre;telefono;comuna;agencia;nro_seguimiento;productos;descuento;costo_envio;total;pago;registrado_por;fecha_pagado;fecha_despachado;despachado_por;fecha_entregado;motivo_cancelacion"
	if strings.Join(filas[0], ";") != want {
		t.Fatalf("columnas:\n%s", strings.Join(filas[0], ";"))
	}
	f := map[string]string{}
	for i, h := range filas[0] {
		f[h] = filas[1][i]
	}
	// Semilla a 2026-11-14 15:00 UTC; DC-0103 creado hace 50 h = 2026-11-12 13:00 UTC = 10:00 en Chile.
	if f["numero"] != "DC-0103" || f["fecha_registro"] != "12-11-2026 10:00" || f["canal"] != "Regiones" || f["entrega"] != "Despacho" ||
		f["estado"] != "Despachado" || f["telefono"] != "56900000103" || f["agencia"] != "Starken" ||
		f["productos"] != "2 x Tren navideño para árbol (grande)" || f["pago"] != "Transferencia" ||
		f["despachado_por"] != "Regiones 1" || f["fecha_despachado"] != "13-11-2026 14:00" || f["fecha_pagado"] == "" {
		t.Fatalf("fila: %v", f)
	}
}

// B8: nadie más que el admin principal puede editarlo o cambiarle la clave.
func TestAdminPrincipalProtegido(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	var ses struct {
		ID int64 `json:"id"`
	}
	e.pedir(t, "GET", "/api/admin/sesion", nil, nil, c).json(t, &ses)
	r := e.pedir(t, "POST", "/api/admin/usuarios", map[string]string{"usuario": "otro", "nombre": "Otro admin", "rol": "admin"}, nil, c)
	var nuevo struct {
		ClaveTemporal string `json:"clave_temporal"`
	}
	r.json(t, &nuevo)
	o := e.loginComo(t, "otro", nuevo.ClaveTemporal)
	e.pedir(t, "POST", "/api/admin/sesion/clave", map[string]string{"clave_actual": nuevo.ClaveTemporal, "clave_nueva": "clave del otro admin"}, nil, o)
	for _, rt := range []struct {
		m, r string
		b    any
	}{
		{"POST", fmt.Sprintf("/api/admin/usuarios/%d/clave", ses.ID), nil},
		{"PATCH", fmt.Sprintf("/api/admin/usuarios/%d", ses.ID), map[string]any{"rol": "vendedor"}},
	} {
		if r := e.pedir(t, rt.m, rt.r, rt.b, nil, o); r.status != 403 || !bytes.Contains(r.body, []byte("admin_principal_protegido")) {
			t.Fatalf("%s %s: %d %s", rt.m, rt.r, r.status, r.body)
		}
	}
	// El principal sí se edita a sí mismo (nombre, canal).
	if r := e.pedir(t, "PATCH", fmt.Sprintf("/api/admin/usuarios/%d", ses.ID), map[string]any{"nombre": "Dueño", "canal_por_defecto": "santiago"}, nil, c); r.status != 200 {
		t.Fatalf("auto-edición: %d", r.status)
	}
	if r := e.pedir(t, "GET", "/api/admin/sesion", nil, nil, c); !bytes.Contains(r.body, []byte(`"canal_por_defecto":"santiago"`)) {
		t.Fatalf("canal por defecto en sesión (M-15): %s", r.body)
	}
}

// B10: el link vence 30 días después de entregado y se puede regenerar.
func TestSeguimientoVenceYSeRegenera(t *testing.T) {
	e := nuevoEntorno(t, nil)
	c := e.login(t)
	p := e.buscarPedido(t, c, "DC-0107") // entregado hace 52 h
	tok := func(u string) string { return u[strings.LastIndex(u, "/")+1:] }
	if r := e.pedir(t, "GET", "/api/seguimiento/"+tok(p.SeguimientoURL), nil, nil); r.status != 200 {
		t.Fatalf("antes de vencer: %d", r.status)
	}
	e.avanzar(29 * 24 * time.Hour)
	if r := e.pedir(t, "GET", "/api/seguimiento/"+tok(p.SeguimientoURL), nil, nil); r.status != 404 {
		t.Fatalf("vencido: %d", r.status)
	}
	c = e.login(t) // la sesión anterior venció con el salto de reloj
	r := e.pedir(t, "POST", fmt.Sprintf("/api/admin/pedidos/%d/token", p.ID), nil, nil, c)
	var n pedidoResp
	r.json(t, &n)
	if r.status != 200 || n.SeguimientoURL == p.SeguimientoURL {
		t.Fatalf("regenerar: %d %s", r.status, r.body)
	}
	if r := e.pedir(t, "GET", "/api/seguimiento/"+tok(n.SeguimientoURL), nil, nil); r.status != 200 {
		t.Fatalf("token nuevo: %d", r.status)
	}
}
