package httpapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/mensajes"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

const (
	pageSizePedidosDef = 20
	maxFilasCSV        = 10000
	maxListaHoy        = 50
	maxDiasResumen     = 366
)

// diaChile convierte "AAAA-MM-DD" (día de Chile) en el instante UTC en que empieza.
func (s *Server) diaChile(v string) (time.Time, bool) {
	f, err := time.ParseInLocation(time.DateOnly, v, s.cfg.Zona)
	if err != nil {
		return time.Time{}, false
	}
	return f.UTC(), true
}

// inicioHoy es el comienzo del día de hoy en Chile, en UTC.
func (s *Server) inicioHoy() time.Time {
	h := s.hoy()
	return time.Date(h.Year(), h.Month(), h.Day(), 0, 0, 0, 0, s.cfg.Zona).UTC()
}

// filtroPedidos lee canal, estado, q, desde y hasta.
func (s *Server) filtroPedidos(q url.Values) (store.FiltroPedidos, map[string]string) {
	campos := map[string]string{}
	var f store.FiltroPedidos
	if c := q.Get("canal"); c != "" {
		if !pedidos.EsCanal(c) {
			campos["canal"] = "debe ser santiago o regiones"
		}
		f.Canal = c
	}
	if e := q.Get("estado"); e != "" {
		if !pedidos.EsEstado(e) {
			campos["estado"] = "estado desconocido"
		}
		f.Estados = []string{e}
	}
	if v := q.Get("q"); v != "" {
		limpio, ok := texto.Limpiar(v, false)
		if !ok || texto.Largo(limpio) > maxLargoBusqueda {
			campos["q"] = "máximo 60 caracteres"
		} else {
			sinSep := strings.NewReplacer(" ", "", "+", "", "-", "", "(", "", ")", "", ".", "").Replace(limpio)
			if dig := pedidos.Digitos(sinSep); len(dig) >= 4 && dig == sinSep {
				f.Palabras = []string{dig} // teléfono (o parte)
			} else {
				f.Palabras = productos.PalabrasBusqueda(limpio, maxPalabras)
			}
		}
	}
	if v := q.Get("desde"); v != "" {
		d, ok := s.diaChile(v)
		if !ok {
			campos["desde"] = "fecha AAAA-MM-DD"
		}
		f.CreadoDesde = d
	}
	if v := q.Get("hasta"); v != "" {
		h, ok := s.diaChile(v)
		if !ok {
			campos["hasta"] = "fecha AAAA-MM-DD"
		}
		f.CreadoHasta = h.AddDate(0, 0, 1) // inclusive
	}
	if !f.CreadoDesde.IsZero() && !f.CreadoHasta.IsZero() && !f.CreadoDesde.Before(f.CreadoHasta) {
		campos["hasta"] = "debe ser igual o posterior a desde"
	}
	return f, campos
}

func (s *Server) listarPedidos(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	q := r.URL.Query()
	f, campos := s.filtroPedidos(q)
	page, size, cp := paginacion(q, pageSizePedidosDef, pageSizeAdminMax)
	for k, v := range cp {
		campos[k] = v
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los filtros.", campos)
		return
	}
	f.Limite, f.Desplazamiento = size, (page-1)*size
	lista, total, err := s.st.ListarPedidos(r.Context(), f)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, nuevaPagina(s.pedidosDTO(lista), page, size, total))
}

func nuevoTokenSeguimiento() (string, error) {
	b := make([]byte, 16) // 128 bits; base64url = 22 caracteres. Nunca el DC-xxxx (se adivina).
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// responderErrorPedido traduce errores de dominio y store a respuestas.
func (s *Server) responderErrorPedido(w http.ResponseWriter, r *http.Request, err error) {
	var items *pedidos.ErrItems
	var ti *pedidos.ErrTransicionInvalida
	var req *pedidos.ErrRequisitos
	var sinStock *store.ErrSinStock
	switch {
	case errors.As(err, &items):
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los productos del pedido.", items.Campos)
	case errors.Is(err, store.ErrTotalNegativo):
		responderError(w, http.StatusBadRequest, "validacion", "Revisa el descuento.", map[string]string{"descuento": "el total quedaría negativo"})
	case errors.Is(err, store.ErrTotalNoCoincide):
		responderError(w, http.StatusBadRequest, "validacion", "Revisa el total.", map[string]string{"total": "no coincide con productos − descuento + envío (o falta, si hay líneas sin precio)"})
	case errors.Is(err, store.ErrNoEditable):
		responderError(w, http.StatusConflict, "no_editable", "El pedido ya está pagado o más avanzado: no se puede editar.", nil)
	case errors.Is(err, store.ErrNoEncontrado):
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
	case errors.As(err, &ti):
		responderError(w, http.StatusConflict, "transicion_invalida", "No se puede pasar de \""+ti.De+"\" a \""+ti.A+"\".", nil)
	case errors.As(err, &req):
		responderError(w, http.StatusBadRequest, "validacion", "Faltan datos para este cambio de estado.", req.Campos)
	case errors.As(err, &sinStock):
		det := make([]map[string]any, 0, len(sinStock.ProductoIDs))
		for _, id := range sinStock.ProductoIDs {
			det = append(det, map[string]any{"producto_id": id, "motivo": "stock_insuficiente"})
		}
		responderErrorDetalles(w, http.StatusConflict, "stock_insuficiente",
			"No hay stock suficiente para marcarlo como pagado. Ajusta el stock o el pedido.", det)
	case errors.Is(err, store.ErrConflicto):
		responderError(w, http.StatusConflict, "conflicto", "El pedido cambió mientras lo editabas. Recarga e intenta de nuevo.", nil)
	default:
		s.errorInterno(w, r, err)
	}
}

func (s *Server) leerEntradaPedido(w http.ResponseWriter, r *http.Request) (store.DatosPedido, string, bool) {
	var in pedidos.Entrada
	if err := leerJSON(w, r, &in, 64<<10); err != nil {
		responderErrorDecodificar(w, err)
		return store.DatosPedido{}, "", false
	}
	d, inicial, campos := in.Validar()
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos del pedido.", campos)
		return store.DatosPedido{}, "", false
	}
	return d, inicial, true
}

func (s *Server) crearPedido(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	d, inicial, ok := s.leerEntradaPedido(w, r)
	if !ok {
		return
	}
	token, err := nuevoTokenSeguimiento()
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	p, err := s.st.CrearPedido(r.Context(), store.PedidoNuevo{
		DatosPedido: d, EstadoInicial: inicial, TokenSeguimiento: token,
		RegistradoPor: ses.UsuarioID, EsDemo: s.cfg.ModoDemo, Ahora: s.ahora(),
	}, pedidos.Cotizador(s.hoyISO()))
	if err != nil {
		s.responderErrorPedido(w, r, err)
		return
	}
	s.log.Info("pedido registrado", "numero", p.Numero, "canal", p.Canal, "total", p.Total,
		"telefono", pedidos.EnmascararTelefono(p.Telefono), "usuario", ses.Usuario)
	responderJSON(w, http.StatusCreated, s.pedidoDTO(p, true))
}

func (s *Server) obtenerPedido(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	p, err := s.st.PedidoPorID(r.Context(), id)
	if err != nil {
		s.responderErrorPedido(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, s.pedidoDTO(p, true))
}

func (s *Server) editarPedido(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	d, _, ok := s.leerEntradaPedido(w, r)
	if !ok {
		return
	}
	p, err := s.st.EditarPedido(r.Context(), id, d, pedidos.Cotizador(s.hoyISO()), s.ahora())
	if err != nil {
		s.responderErrorPedido(w, r, err)
		return
	}
	s.log.Info("pedido editado", "numero", p.Numero, "total", p.Total, "usuario", ses.Usuario)
	responderJSON(w, http.StatusOK, s.pedidoDTO(p, true))
}

func (s *Server) cambiarEstado(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	var in pedidos.CambioEntrada
	if err := leerJSON(w, r, &in, 8<<10); err != nil {
		responderErrorDecodificar(w, err)
		return
	}
	c, campos := in.Validar(s.hoy())
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los datos.", campos)
		return
	}
	c.PedidoID, c.UsuarioID, c.Ahora, c.StockEstricto = id, ses.UsuarioID, s.ahora(), s.cfg.StockEstricto
	p, err := s.st.CambiarEstado(r.Context(), c, pedidos.Validador())
	if err != nil {
		s.responderErrorPedido(w, r, err)
		return
	}
	s.log.Info("estado de pedido cambiado", "numero", p.Numero, "estado", p.Estado, "usuario", ses.Usuario)
	d := s.pedidoDTO(p, true)
	if len(p.AvisoStock) > 0 {
		d.AvisoStock = s.avisoStock(p)
		s.log.Warn("stock insuficiente aceptado (STOCK_ESTRICTO=false)", "numero", p.Numero, "productos", len(p.AvisoStock))
	}
	responderJSON(w, http.StatusOK, d)
}

func (s *Server) mensajesPedido(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	p, err := s.st.PedidoPorID(r.Context(), id)
	if err != nil {
		s.responderErrorPedido(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, map[string]any{
		"numero":           p.Numero,
		"telefono_cliente": p.Telefono,
		"mensajes":         mensajes.Para(p, s.cfg.Tienda, s.ahora()),
	})
}

type clienteDTO struct {
	Telefono      string      `json:"telefono"`
	Nombre        *string     `json:"nombre"`
	TotalPedidos  int         `json:"total_pedidos"`
	TotalComprado int64       `json:"total_comprado"`
	Pedidos       []pedidoDTO `json:"pedidos"`
}

func (s *Server) cliente(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	tel, ok := pedidos.NormalizarTelefono(r.URL.Query().Get("telefono"))
	if !ok {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa el teléfono.", map[string]string{"telefono": "número chileno, por ejemplo +56 9 1234 5678"})
		return
	}
	lista, total, err := s.st.ListarPedidos(r.Context(), store.FiltroPedidos{Telefono: tel, Limite: 200})
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	out := clienteDTO{Telefono: tel, TotalPedidos: total, Pedidos: s.pedidosDTO(lista)}
	if len(lista) > 0 {
		out.Nombre = &lista[0].Nombre
	}
	for _, p := range lista {
		if p.Estado != pedidos.EstadoCancelado {
			out.TotalComprado += p.Total
		}
	}
	responderJSON(w, http.StatusOK, out)
}

// ---------------------------------------------------------------- Hoy

type hoyDTO struct {
	Fecha               string         `json:"fecha"`
	Contadores          map[string]int `json:"contadores"`
	EsperandoPagoMas24h []pedidoDTO    `json:"esperando_pago_mas_24h"`
	PagadosSinDespachar []pedidoDTO    `json:"pagados_sin_despachar"`
	ListosRetiro        []pedidoDTO    `json:"listos_retiro"`
	CreadosHoy          []pedidoDTO    `json:"creados_hoy"`
}

func (s *Server) hoyPanel(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	ctx := r.Context()
	ahora := s.ahora()
	inicio := s.inicioHoy()
	consultas := []store.FiltroPedidos{
		{Estados: []string{pedidos.EstadoEsperandoPago}, EstadoDesdeAntes: ahora.Add(-24 * time.Hour)},
		{Estados: []string{pedidos.EstadoPagado, pedidos.EstadoPreparando}},
		{Estados: []string{pedidos.EstadoListoRetiro}},
		{CreadoDesde: inicio},
	}
	var listas [4][]pedidoDTO
	var totales [4]int
	for i, f := range consultas {
		f.Limite = maxListaHoy
		l, t, err := s.st.ListarPedidos(ctx, f)
		if err != nil {
			s.errorInterno(w, r, err)
			return
		}
		listas[i], totales[i] = s.pedidosDTO(l), t
	}
	_, santiago, err := s.st.ListarPedidos(ctx, store.FiltroPedidos{CreadoDesde: inicio, Canal: pedidos.CanalSantiago, Limite: 1})
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	_, regiones, err := s.st.ListarPedidos(ctx, store.FiltroPedidos{CreadoDesde: inicio, Canal: pedidos.CanalRegiones, Limite: 1})
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, hoyDTO{
		Fecha: s.hoyISO(),
		Contadores: map[string]int{
			"esperando_pago_mas_24h": totales[0],
			"pagados_sin_despachar":  totales[1],
			"listos_retiro":          totales[2],
			"creados_hoy_santiago":   santiago,
			"creados_hoy_regiones":   regiones,
		},
		EsperandoPagoMas24h: listas[0], PagadosSinDespachar: listas[1], ListosRetiro: listas[2], CreadosHoy: listas[3],
	})
}

// ---------------------------------------------------------------- Resumen

// filaResumenDTO: grupo con pedidos (no cancelados) y monto vendido (con pago registrado).
type filaResumenDTO struct {
	Clave   *string `json:"clave"`
	Nombre  *string `json:"nombre"`
	Pedidos int     `json:"pedidos"`
	Monto   int64   `json:"monto"`
}

type comunaDTO struct {
	Comuna  string `json:"comuna"`
	Pedidos int    `json:"pedidos"`
	Monto   int64  `json:"monto"`
}

type usuarioResumenDTO struct {
	Usuario     string `json:"usuario"`
	Registrados int    `json:"registrados"`
	Despachados int    `json:"despachados"`
}

type resumenDTO struct {
	Desde      string              `json:"desde"`
	Hasta      string              `json:"hasta"`
	Totales    map[string]int64    `json:"totales"`
	PorCanal   []filaResumenDTO    `json:"por_canal"`
	PorEstado  []filaResumenDTO    `json:"por_estado"`
	PorAgencia []filaResumenDTO    `json:"por_agencia"`
	PorComuna  []comunaDTO         `json:"por_comuna"`
	PorUsuario []usuarioResumenDTO `json:"por_usuario"`
}

// resumen: GET /api/admin/resumen?desde&hasta (días de Chile, por defecto los últimos 30).
//   - pedidos: no cancelados; cancelados aparte.
//   - vendido (y "monto" de cada fila): solo pedidos con pago registrado y no cancelados.
func (s *Server) resumen(w http.ResponseWriter, r *http.Request, _ store.Sesion) {
	q := r.URL.Query()
	hoy := s.hoy()
	desdeStr, hastaStr := q.Get("desde"), q.Get("hasta")
	if hastaStr == "" {
		hastaStr = hoy.Format(time.DateOnly)
	}
	if desdeStr == "" {
		desdeStr = hoy.AddDate(0, 0, -29).Format(time.DateOnly)
	}
	campos := map[string]string{}
	desde, ok1 := s.diaChile(desdeStr)
	hasta, ok2 := s.diaChile(hastaStr)
	switch {
	case !ok1:
		campos["desde"] = "fecha AAAA-MM-DD"
	case !ok2:
		campos["hasta"] = "fecha AAAA-MM-DD"
	case hasta.Before(desde):
		campos["hasta"] = "debe ser igual o posterior a desde"
	case hasta.Sub(desde) > maxDiasResumen*24*time.Hour:
		campos["desde"] = "máximo 366 días"
	}
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa las fechas.", campos)
		return
	}
	finExcl := hasta.AddDate(0, 0, 1)
	filas, err := s.st.ResumenPedidos(r.Context(), desde, finExcl)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	usuarios, err := s.st.ResumenUsuarios(r.Context(), desde, finExcl)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}

	porCanal := map[string]*filaResumenDTO{}
	porEstado := map[string]*filaResumenDTO{}
	porAgencia := map[string]*filaResumenDTO{}
	porComuna := map[string]*comunaDTO{}
	tot := map[string]int64{"pedidos": 0, "pedidos_santiago": 0, "pedidos_regiones": 0, "pedidos_vendidos": 0,
		"vendido": 0, "monto": 0, "ticket_promedio": 0, "cancelados": 0}
	sumar := func(m map[string]*filaResumenDTO, clave, nombre string, n int, vendido int64) {
		x, ok := m[clave]
		if !ok {
			x = &filaResumenDTO{Clave: ptr(clave), Nombre: ptr(nombre)}
			m[clave] = x
		}
		x.Pedidos += n
		x.Monto += vendido
	}
	for _, f := range filas {
		var vendido int64
		if f.Pagado && f.Estado != pedidos.EstadoCancelado {
			vendido = f.Monto
		}
		sumar(porEstado, f.Estado, "", f.Pedidos, vendido)
		if f.Estado == pedidos.EstadoCancelado {
			tot["cancelados"] += int64(f.Pedidos)
			continue
		}
		tot["pedidos"] += int64(f.Pedidos)
		tot["pedidos_"+f.Canal] += int64(f.Pedidos)
		if vendido > 0 || f.Pagado {
			tot["pedidos_vendidos"] += int64(f.Pedidos)
		}
		tot["vendido"] += vendido
		sumar(porCanal, f.Canal, "", f.Pedidos, vendido)
		if f.Entrega == pedidos.EntregaEnvio {
			nombre := ""
			if a, ok := pedidos.AgenciaPorCodigo(f.Agencia); ok {
				nombre = a.Nombre
			}
			sumar(porAgencia, f.Agencia, nombre, f.Pedidos, vendido)
		}
		if f.Comuna != "" {
			c, ok := porComuna[f.Comuna]
			if !ok {
				c = &comunaDTO{Comuna: f.Comuna}
				porComuna[f.Comuna] = c
			}
			c.Pedidos += f.Pedidos
			c.Monto += vendido
		}
	}
	tot["monto"] = tot["vendido"]
	if tot["pedidos_vendidos"] > 0 {
		tot["ticket_promedio"] = tot["vendido"] / tot["pedidos_vendidos"]
	}
	ordenar := func(m map[string]*filaResumenDTO, orden []string) []filaResumenDTO {
		out := []filaResumenDTO{}
		for _, k := range orden {
			if x, ok := m[k]; ok {
				out = append(out, *x)
				delete(m, k)
			}
		}
		for _, x := range m {
			out = append(out, *x)
		}
		return out
	}
	agencias := make([]string, 0, len(pedidos.Agencias))
	for _, a := range pedidos.Agencias {
		agencias = append(agencias, a.Codigo)
	}
	comunas := make([]comunaDTO, 0, len(porComuna))
	for _, c := range porComuna {
		comunas = append(comunas, *c)
	}
	sort.Slice(comunas, func(i, j int) bool {
		if comunas[i].Pedidos != comunas[j].Pedidos {
			return comunas[i].Pedidos > comunas[j].Pedidos
		}
		return comunas[i].Comuna < comunas[j].Comuna
	})
	if len(comunas) > 20 {
		comunas = comunas[:20]
	}
	porUsuario := make([]usuarioResumenDTO, 0, len(usuarios))
	for _, u := range usuarios {
		porUsuario = append(porUsuario, usuarioResumenDTO{Usuario: u.Nombre, Registrados: u.Registrados, Despachados: u.Despachados})
	}
	responderJSON(w, http.StatusOK, resumenDTO{
		Desde: desdeStr, Hasta: hastaStr, Totales: tot,
		PorCanal:   ordenar(porCanal, []string{pedidos.CanalSantiago, pedidos.CanalRegiones}),
		PorEstado:  ordenar(porEstado, pedidos.Estados),
		PorAgencia: ordenar(porAgencia, append(agencias, "")),
		PorComuna:  comunas,
		PorUsuario: porUsuario,
	})
}

// ---------------------------------------------------------------- Ajustes

// ajustes: GET /api/admin/ajustes (solo admin, solo lectura). Nunca devuelve
// el texto de DATOS_TRANSFERENCIA, solo si está configurado.
func (s *Server) ajustes(w http.ResponseWriter, _ *http.Request, _ store.Sesion) {
	responderJSON(w, http.StatusOK, map[string]any{
		"whatsapp_santiago":                s.cfg.WhatsAppSantiago,
		"whatsapp_regiones":                ptr(s.cfg.WhatsAppRegiones),
		"direccion_retiro":                 s.cfg.Tienda.DireccionRetiro,
		"horario_retiro":                   s.cfg.Tienda.HorarioRetiro,
		"url_sitio":                        s.cfg.Tienda.URLSitio,
		"datos_transferencia_configurados": strings.TrimSpace(s.cfg.Tienda.DatosTransferencia) != "",
		"modo_demo":                        s.cfg.ModoDemo,
	})
}

// regenerarToken: POST /api/admin/pedidos/{id}/token (solo admin). El link anterior deja de servir.
func (s *Server) regenerarToken(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Pedido no encontrado.", nil)
		return
	}
	token, err := nuevoTokenSeguimiento()
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	p, err := s.st.RegenerarToken(r.Context(), id, token, s.ahora())
	if err != nil {
		s.responderErrorPedido(w, r, err)
		return
	}
	s.log.Info("link de seguimiento regenerado", "numero", p.Numero, "usuario", ses.Usuario)
	responderJSON(w, http.StatusOK, s.pedidoDTO(p, true))
}

// ---------------------------------------------------------------- CSV

// celdaSegura evita CSV injection: un texto que empieza con = + - @ tab o CR
// se antepone con ' para que la planilla no lo interprete como fórmula.
func celdaSegura(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}

var (
	nombreCanal   = map[string]string{"santiago": "Santiago", "regiones": "Regiones"}
	nombreEntrega = map[string]string{"retiro": "Retiro", "envio": "Despacho"}
	nombreEstado  = map[string]string{"cotizado": "Cotizado", "esperando_pago": "Esperando pago", "pagado": "Pagado",
		"preparando": "Preparando", "despachado": "Despachado", "listo_retiro": "Listo para retiro",
		"entregado": "Entregado", "cancelado": "Cancelado"}
	nombrePago = map[string]string{"transferencia": "Transferencia", "en_tienda": "En tienda"}
)

// fechaChile: dd-mm-aaaa hh:mm en America/Santiago.
func (s *Server) fechaChile(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(s.cfg.Zona).Format("02-01-2006 15:04")
}

// exportarCSV: columnas de propuesta/copy-landing.md P5, en orden.
func (s *Server) exportarCSV(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	f, campos := s.filtroPedidos(r.URL.Query())
	if len(campos) > 0 {
		responderError(w, http.StatusBadRequest, "validacion", "Revisa los filtros.", campos)
		return
	}
	f.Limite, f.ConHistorial = maxFilasCSV, true
	lista, _, err := s.st.ListarPedidos(r.Context(), f)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="pedidos-`+s.hoy().Format("20060102")+`.csv"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("\xEF\xBB\xBF")) // BOM: Excel abre bien las tildes
	cw := csv.NewWriter(w)
	cw.Comma = ';'
	_ = cw.Write([]string{"numero", "fecha_registro", "canal", "entrega", "estado", "nombre", "telefono", "comuna", "agencia",
		"nro_seguimiento", "productos", "descuento", "costo_envio", "total", "pago", "registrado_por", "fecha_pagado",
		"fecha_despachado", "despachado_por", "fecha_entregado", "motivo_cancelacion"})
	n := func(v int64) string { return strconv.FormatInt(v, 10) }
	for _, p := range lista {
		var prods []string
		for _, it := range p.Items {
			prods = append(prods, strconv.Itoa(it.Cantidad)+" x "+it.Descripcion)
		}
		var reg, fPagado, fDesp, despPor, fEnt string
		if p.RegistradoPor != nil {
			reg = p.RegistradoPor.Nombre
		}
		for _, h := range p.Historial {
			switch h.EstadoNuevo {
			case pedidos.EstadoPagado:
				fPagado = s.fechaChile(h.CreadoEn)
			case pedidos.EstadoDespachado:
				fDesp = s.fechaChile(h.CreadoEn)
				if h.Usuario != nil {
					despPor = h.Usuario.Nombre
				}
			case pedidos.EstadoEntregado:
				fEnt = s.fechaChile(h.CreadoEn)
				if fPagado == "" && p.Pago != nil { // pagó al retirar
					fPagado = fEnt
				}
			}
		}
		agencia := p.Agencia
		if a, ok := pedidos.AgenciaPorCodigo(p.Agencia); ok {
			agencia = a.Nombre
		}
		_ = cw.Write([]string{
			p.Numero, s.fechaChile(p.CreadoEn), nombreCanal[p.Canal], nombreEntrega[p.Entrega], nombreEstado[p.Estado],
			celdaSegura(p.Nombre), pedidos.Digitos(p.Telefono), celdaSegura(p.Comuna), celdaSegura(agencia),
			celdaSegura(p.NSeguimiento), celdaSegura(strings.Join(prods, " | ")), n(p.Descuento), n(p.CostoEnvio), n(p.Total),
			nombrePago[p.FormaPago], celdaSegura(reg), fPagado, fDesp, celdaSegura(despPor), fEnt, celdaSegura(p.MotivoCancelacion),
		})
	}
	cw.Flush()
	s.log.Info("CSV exportado", "filas", len(lista), "usuario", ses.Usuario)
}
