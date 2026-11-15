package sqlite_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store/sqlite"
)

var cotizar = pedidos.Cotizador("2026-11-14")

func abrir(t *testing.T) *sqlite.DB {
	t.Helper()
	db, err := sqlite.Abrir(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func sembrarProducto(t *testing.T, db *sqlite.DB, nombre string, precio int64, stock int) store.Producto {
	t.Helper()
	ctx := context.Background()
	cats, _ := db.ListarCategorias(ctx)
	var catID int64
	if len(cats) == 0 {
		c := store.Categoria{Slug: "navidad", Nombre: "Navidad"}
		if err := db.CrearCategoria(ctx, &c, time.Now()); err != nil {
			t.Fatal(err)
		}
		catID = c.ID
	} else {
		catID = cats[0].ID
	}
	cod := strings.ToUpper("C-" + productos.Slugify(nombre))
	if len(cod) > 20 {
		cod = cod[:20]
	}
	p, errs := productos.Entrada{Codigo: cod, Nombre: nombre, CategoriaID: catID, Precio: precio, Stock: stock}.Validar()
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	p.CategoriaNombre = "Navidad"
	p.Busqueda = productos.TextoBusqueda(p)
	if err := db.CrearProducto(ctx, &p, time.Now()); err != nil {
		t.Fatal(err)
	}
	return p
}

var tokens int

func nuevo(t *testing.T, telefono string, items ...store.ItemSolicitado) store.PedidoNuevo {
	t.Helper()
	tokens++
	return store.PedidoNuevo{
		DatosPedido: store.DatosPedido{Canal: "regiones", Nombre: "Cliente Test", Telefono: telefono, Busqueda: "cliente test",
			Entrega: "envio", Comuna: "Chillán", Agencia: "starken", FormaPago: "transferencia", Items: items},
		EstadoInicial: "cotizado", TokenSeguimiento: fmt.Sprintf("tok%019d", tokens), Ahora: time.Now(),
	}
}

func pagar(db *sqlite.DB, id int64, monto int64) (store.Pedido, error) {
	return db.CambiarEstado(context.Background(), store.CambioEstado{PedidoID: id, Nuevo: "pagado", ConfirmoPago: true,
		Pago: &store.Pago{Medio: "transferencia", Monto: monto, Fecha: "2026-11-14"}, Ahora: time.Now()}, pedidos.Validador())
}

func TestNumeracionCorrelativaConcurrente(t *testing.T) {
	db := abrir(t)
	p := sembrarProducto(t, db, "Esfera navideña", 500, 10)
	const n = 60
	var wg sync.WaitGroup
	numeros := make(chan string, n)
	errs := make(chan error, n)
	pedidosN := make([]store.PedidoNuevo, n)
	for i := range pedidosN {
		pedidosN[i] = nuevo(t, "+56900000001", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1})
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(pn store.PedidoNuevo) {
			defer wg.Done()
			ped, err := db.CrearPedido(context.Background(), pn, cotizar)
			if err != nil {
				errs <- err
				return
			}
			numeros <- ped.Numero
		}(pedidosN[i])
	}
	wg.Wait()
	close(numeros)
	close(errs)
	for err := range errs {
		t.Fatalf("error concurrente: %v", err)
	}
	var lista []string
	vistos := map[string]bool{}
	for num := range numeros {
		if vistos[num] {
			t.Fatalf("número duplicado: %s", num)
		}
		vistos[num] = true
		lista = append(lista, num)
	}
	sort.Strings(lista)
	if len(lista) != n || lista[0] != "DC-0001" || lista[n-1] != "DC-0060" {
		t.Fatalf("correlativo con huecos: %s … %s (%d)", lista[0], lista[len(lista)-1], len(lista))
	}
	if act, _ := db.ProductoPorID(context.Background(), p.ID); act.Stock != 10 {
		t.Fatalf("registrar no descuenta stock: %d", act.Stock)
	}
}

func TestTotalesDescuentoEnvioYManual(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Caja premium", 15000, 100)
	n := nuevo(t, "+56900000002", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 2})
	n.Descuento, n.CostoEnvio = 1000, 5000
	ped, err := db.CrearPedido(ctx, n, cotizar)
	if err != nil || ped.Subtotal != 30000 || ped.Total != 34000 || ped.TotalManual {
		t.Fatalf("total calculado: %+v %v", ped, err)
	}
	mal := int64(1)
	n.TotalManual = &mal
	if _, err := db.CrearPedido(ctx, n, cotizar); !errors.Is(err, store.ErrTotalNoCoincide) {
		t.Fatalf("total del cliente distinto: %v", err)
	}
	n.TotalManual, n.Descuento = nil, 999999
	if _, err := db.CrearPedido(ctx, n, cotizar); !errors.Is(err, store.ErrTotalNegativo) {
		t.Fatalf("total negativo: %v", err)
	}
	// Línea libre sin precio: manda el total escrito.
	total := int64(20000)
	n = nuevo(t, "+56900000003", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1}, store.ItemSolicitado{Descripcion: "Bolsas", Cantidad: 1})
	n.TotalManual = &total
	ped, err = db.CrearPedido(ctx, n, cotizar)
	if err != nil || ped.Total != 20000 || !ped.TotalManual || ped.Items[1].PrecioUnitario != nil || ped.Subtotal != 15000 {
		t.Fatalf("total manual: %+v %v", ped, err)
	}
}

// El stock se descuenta al pagar (no al registrar) y pagar en paralelo no sobrevende.
func TestStockAlPagarSinSobreventa(t *testing.T) {
	db := abrir(t)
	p := sembrarProducto(t, db, "Tren navideño", 17000, 5)
	var ids []int64
	for i := 0; i < 20; i++ {
		ped, err := db.CrearPedido(context.Background(), nuevo(t, fmt.Sprintf("+5690000%04d", i), store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1}), cotizar)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, ped.ID)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	ok, sinStock := 0, 0
	for _, id := range ids {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			_, err := pagar(db, id, 17000)
			mu.Lock()
			defer mu.Unlock()
			var es *store.ErrSinStock
			switch {
			case err == nil:
				ok++
			case errors.As(err, &es) && es.ProductoIDs[0] == p.ID:
				sinStock++
			default:
				t.Errorf("error inesperado: %v", err)
			}
		}(id)
	}
	wg.Wait()
	if ok != 5 || sinStock != 15 {
		t.Fatalf("pagados=%d sin stock=%d", ok, sinStock)
	}
	if act, _ := db.ProductoPorID(context.Background(), p.ID); act.Stock != 0 {
		t.Fatalf("stock final %d", act.Stock)
	}
}

func TestFlujoCompletoHistorialYEdicion(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Calcetas navideñas", 1000, 10)
	ped, err := db.CrearPedido(ctx, nuevo(t, "+56900000004", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 4}), cotizar)
	if err != nil {
		t.Fatal(err)
	}
	// Editar antes de pagar: cambia cantidades y precio congelado.
	d := nuevo(t, "+56900000004", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 6}).DatosPedido
	ed, err := db.EditarPedido(ctx, ped.ID, d, cotizar, time.Now())
	if err != nil || ed.Total != 6000 {
		t.Fatalf("editar: %+v %v", ed, err)
	}
	if _, err := pagar(db, ped.ID, 6000); err != nil {
		t.Fatal(err)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 4 {
		t.Fatalf("pagar descuenta: %d", act.Stock)
	}
	if _, err := db.EditarPedido(ctx, ped.ID, d, cotizar, time.Now()); !errors.Is(err, store.ErrNoEditable) {
		t.Fatalf("pagado no se edita: %v", err)
	}
	for _, c := range []store.CambioEstado{
		{Nuevo: "preparando"},
		{Nuevo: "despachado", NSeguimiento: "000123456789"},
	} {
		c.PedidoID, c.Ahora = ped.ID, time.Now()
		if _, err := db.CambiarEstado(ctx, c, pedidos.Validador()); err != nil {
			t.Fatalf("%s: %v", c.Nuevo, err)
		}
	}
	fin, err := db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "cancelado", Motivo: "Se nos agotó", Ahora: time.Now()}, pedidos.Validador())
	if err != nil {
		t.Fatal(err)
	}
	if fin.MotivoCancelacion != "Se nos agotó" || len(fin.Historial) != 5 || fin.Historial[4].EstadoAnterior != "despachado" ||
		fin.NSeguimiento != "000123456789" || fin.Pago == nil || fin.Pago.Monto != 6000 || fin.StockDescontado {
		t.Fatalf("pedido final: %+v", fin)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 10 {
		t.Fatalf("cancelar devuelve: %d", act.Stock)
	}
	por, err := db.PedidoPorToken(ctx, fin.TokenSeguimiento)
	if err != nil || por.ID != fin.ID {
		t.Fatalf("por token: %v", err)
	}
}

func TestFiltrosPedidos(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Lazos", 1000, 100)
	a := nuevo(t, "+56900000005", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1})
	a.Busqueda = "carla demo 56900000005"
	b := nuevo(t, "+56900000006", store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1})
	b.Canal, b.Entrega, b.FormaPago, b.Busqueda = "santiago", "retiro", "en_tienda", "jorge demo 56900000006"
	b.Ahora = time.Now().Add(-48 * time.Hour)
	for _, n := range []store.PedidoNuevo{a, b} {
		if _, err := db.CrearPedido(ctx, n, cotizar); err != nil {
			t.Fatal(err)
		}
	}
	contar := func(f store.FiltroPedidos) int {
		_, n, err := db.ListarPedidos(ctx, f)
		if err != nil {
			t.Fatal(err)
		}
		return n
	}
	for nombre, c := range map[string]struct {
		f    store.FiltroPedidos
		want int
	}{
		"canal":        {store.FiltroPedidos{Canal: "santiago"}, 1},
		"nombre":       {store.FiltroPedidos{Palabras: []string{"carla"}}, 1},
		"numero":       {store.FiltroPedidos{Palabras: []string{"dc-0002"}}, 1},
		"telefono":     {store.FiltroPedidos{Palabras: []string{"0006"}}, 1},
		"tel exacto":   {store.FiltroPedidos{Telefono: "+56900000005"}, 1},
		"desde":        {store.FiltroPedidos{CreadoDesde: time.Now().Add(-time.Hour)}, 1},
		"estado":       {store.FiltroPedidos{Estados: []string{"cotizado"}}, 2},
		"estado_desde": {store.FiltroPedidos{EstadoDesdeAntes: time.Now().Add(-24 * time.Hour)}, 1},
		"inyeccion":    {store.FiltroPedidos{Palabras: []string{"' or 1=1 --"}}, 0},
		"comodin":      {store.FiltroPedidos{Palabras: []string{"%"}}, 0},
	} {
		if got := contar(c.f); got != c.want {
			t.Errorf("%s: %d, quería %d", nombre, got, c.want)
		}
	}
	filas, err := db.ResumenPedidos(ctx, time.Now().Add(-72*time.Hour), time.Now().Add(time.Hour))
	if err != nil || len(filas) != 2 {
		t.Fatalf("resumen: %+v %v", filas, err)
	}
}

func TestBusquedaProductosInyeccion(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	sembrarProducto(t, db, "Telón de Navidad", 5000, 10)
	sembrarProducto(t, db, "Lazos navideños", 1000, 10)
	buscar := func(q string) int {
		_, total, err := db.ListarProductos(ctx, store.FiltroProductos{Palabras: productos.PalabrasBusqueda(q, 5), Limite: 50})
		if err != nil {
			t.Fatalf("q=%q: %v", q, err)
		}
		return total
	}
	if buscar("telon") != 1 || buscar("NAVIDEÑOS") != 1 {
		t.Fatal("búsqueda normal")
	}
	for _, q := range []string{"%", "_", `\`, "' OR '1'='1", "'; DROP TABLE productos; --"} {
		if n := buscar(q); n != 0 {
			t.Errorf("q=%q devolvió %d", q, n)
		}
	}
	if buscar("") != 2 {
		t.Fatal("tabla intacta")
	}
}

func TestUsuariosYSesiones(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	id, err := db.SincronizarAdmin(ctx, "admin", "hash-a", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	v := store.Usuario{Usuario: "vane", Nombre: "Vane", Rol: "vendedor", PasswordHash: "hash-v"}
	if err := db.CrearUsuario(ctx, &v, time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := db.CrearUsuario(ctx, &store.Usuario{Usuario: "vane", Nombre: "Otra", Rol: "vendedor", PasswordHash: "x"}, time.Now()); !errors.Is(err, store.ErrDuplicado) {
		t.Fatalf("duplicado: %v", err)
	}
	ahora := time.Now()
	ses := store.Sesion{TokenHash: "h1", UsuarioID: v.ID, CreadoEn: ahora, ExpiraEn: ahora.Add(time.Hour), UltimoUso: ahora}
	if err := db.CrearSesion(ctx, ses); err != nil {
		t.Fatal(err)
	}
	if s, err := db.SesionPorHash(ctx, "h1"); err != nil || s.Rol != "vendedor" || s.Nombre != "Vane" {
		t.Fatalf("sesión: %+v %v", s, err)
	}
	inactivo := false
	if _, err := db.ActualizarUsuario(ctx, v.ID, store.CambiosUsuario{Activo: &inactivo}, ahora); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SesionPorHash(ctx, "h1"); !errors.Is(err, store.ErrNoEncontrado) {
		t.Fatal("desactivar cierra sus sesiones")
	}
	// El admin del .env se reactiva y recupera rol al arrancar.
	vend := "vendedor"
	if _, err := db.ActualizarUsuario(ctx, id, store.CambiosUsuario{Rol: &vend, Activo: &inactivo}, ahora); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SincronizarAdmin(ctx, "admin", "hash-a", time.Now()); err != nil {
		t.Fatal(err)
	}
	if u, _ := db.UsuarioPorID(ctx, id); u.Rol != "admin" || !u.Activo {
		t.Fatalf("admin recuperado: %+v", u)
	}
}
