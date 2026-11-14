package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store/sqlite"
)

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
	ahora := time.Now()
	cats, _ := db.ListarCategorias(ctx)
	var catID int64
	if len(cats) == 0 {
		c := store.Categoria{Slug: "navidad", Nombre: "Navidad"}
		if err := db.CrearCategoria(ctx, &c, ahora); err != nil {
			t.Fatal(err)
		}
		catID = c.ID
	} else {
		catID = cats[0].ID
	}
	e := productos.Entrada{Nombre: nombre, CategoriaID: catID, Precio: precio, Stock: stock}
	e.Codigo = "C-" + productos.Slugify(nombre)
	if len(e.Codigo) > 20 {
		e.Codigo = e.Codigo[:20]
	}
	e.Codigo = upper(e.Codigo)
	p, errs := e.Validar()
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	p.CategoriaNombre = "Navidad"
	p.Busqueda = productos.TextoBusqueda(p)
	if err := db.CrearProducto(ctx, &p, ahora); err != nil {
		t.Fatal(err)
	}
	return p
}

func upper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

func pedidoNuevo(items ...store.ItemSolicitado) store.PedidoNuevo {
	return store.PedidoNuevo{Nombre: "Cliente Test", Telefono: "+56900000001", Entrega: "retiro", Items: items, Ahora: time.Now()}
}

func TestNumeracionCorrelativaConcurrente(t *testing.T) {
	db := abrir(t)
	p := sembrarProducto(t, db, "Esfera navideña", 500, 100000)

	const n = 60
	var wg sync.WaitGroup
	numeros := make(chan string, n)
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ped, err := db.CrearPedido(context.Background(), pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1}), pedidos.Cotizador("2026-11-14"))
			if err != nil {
				errs <- err
				return
			}
			numeros <- ped.Numero
		}()
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
		t.Fatalf("correlativo con huecos: primero %s, último %s, total %d", lista[0], lista[len(lista)-1], len(lista))
	}
	act, _ := db.ProductoPorID(context.Background(), p.ID)
	if act.Stock != 100000 {
		t.Fatalf("crear un pedido no debe descontar stock; quedó %d", act.Stock)
	}
}

func confirmar(db *sqlite.DB, id int64) (store.Pedido, error) {
	return db.CambiarEstado(context.Background(), store.CambioEstado{PedidoID: id, Nuevo: "confirmado", Ahora: time.Now()}, pedidos.Validador("confirmado"))
}

// Un pedido anónimo ya no agota la tienda: el stock se descuenta al confirmar,
// y confirmar en paralelo no sobrevende.
func TestStockSeDescuentaAlConfirmarSinSobreventa(t *testing.T) {
	db := abrir(t)
	p := sembrarProducto(t, db, "Tren navideño", 17000, 5)
	var ids []int64
	for i := 0; i < 20; i++ {
		n := pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1})
		n.Telefono = "+5690000" + string(rune('1'+i/10)) + string(rune('0'+i%10)) + "00"
		ped, err := db.CrearPedido(context.Background(), n, pedidos.Cotizador("2026-11-14"))
		if err != nil {
			t.Fatalf("pedido %d: %v", i, err)
		}
		ids = append(ids, ped.ID)
	}
	if act, _ := db.ProductoPorID(context.Background(), p.ID); act.Stock != 5 {
		t.Fatalf("20 pedidos nuevos no deben tocar el stock: %d", act.Stock)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	ok, sinStock := 0, 0
	for _, id := range ids {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			_, err := confirmar(db, id)
			mu.Lock()
			defer mu.Unlock()
			var es *store.ErrSinStock
			switch {
			case err == nil:
				ok++
			case errors.As(err, &es) && errors.Is(err, store.ErrStockInsuficiente) && es.ProductoIDs[0] == p.ID:
				sinStock++
			default:
				t.Errorf("error inesperado: %v", err)
			}
		}(id)
	}
	wg.Wait()
	if ok != 5 || sinStock != 15 {
		t.Fatalf("confirmados=%d sin stock=%d", ok, sinStock)
	}
	act, _ := db.ProductoPorID(context.Background(), p.ID)
	if act.Stock != 0 {
		t.Fatalf("stock final = %d", act.Stock)
	}
	// Un pedido que no se pudo confirmar sigue en "nuevo".
	lista, _, _ := db.ListarPedidos(context.Background(), store.FiltroPedidos{Estado: "nuevo", Limite: 50})
	if len(lista) != 15 {
		t.Fatalf("pendientes = %d", len(lista))
	}
}

func TestCambioEstadoHistorialYStock(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Calcetas navideñas", 1000, 10)
	ped, err := db.CrearPedido(ctx, pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 4}), pedidos.Cotizador("2026-11-14"))
	if err != nil {
		t.Fatal(err)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 10 {
		t.Fatalf("en nuevo no se descuenta: %d", act.Stock)
	}
	// inválida: nuevo -> entregado
	_, err = db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "entregado", Ahora: time.Now()}, pedidos.Validador("entregado"))
	var ti *pedidos.ErrTransicionInvalida
	if !errors.As(err, &ti) {
		t.Fatalf("quería transición inválida: %v", err)
	}
	if _, err := confirmar(db, ped.ID); err != nil {
		t.Fatal(err)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 6 {
		t.Fatalf("al confirmar se descuenta: %d", act.Stock)
	}
	fin, err := db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "cancelado", Nota: "no confirmó", Ahora: time.Now()}, pedidos.Validador("cancelado"))
	if err != nil {
		t.Fatal(err)
	}
	if fin.Estado != "cancelado" || len(fin.Historial) != 3 || fin.Historial[2].EstadoAnterior != "confirmado" || fin.Historial[2].Nota != "no confirmó" {
		t.Fatalf("historial: %+v", fin.Historial)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 10 {
		t.Fatalf("stock tras cancelar un confirmado = %d (debía volver a 10)", act.Stock)
	}
	// cancelado es final
	if _, err := confirmar(db, ped.ID); !errors.As(err, &ti) {
		t.Fatalf("cancelado debe ser final: %v", err)
	}
	// Cancelar un "nuevo" no toca el stock.
	ped2, _ := db.CrearPedido(ctx, pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 3}), pedidos.Cotizador("2026-11-14"))
	if _, err := db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped2.ID, Nuevo: "cancelado", Ahora: time.Now()}, pedidos.Validador("cancelado")); err != nil {
		t.Fatal(err)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 10 {
		t.Fatalf("cancelar un nuevo cambió el stock: %d", act.Stock)
	}
}

func TestTopePedidosAbiertosPorTelefono(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Lazos navideños", 1000, 100)
	nuevo := func() store.PedidoNuevo {
		n := pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1})
		n.MaxAbiertos = 3
		return n
	}
	var primero store.Pedido
	for i := 0; i < 3; i++ {
		ped, err := db.CrearPedido(ctx, nuevo(), pedidos.Cotizador("2026-11-14"))
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			primero = ped
		}
	}
	if _, err := db.CrearPedido(ctx, nuevo(), pedidos.Cotizador("2026-11-14")); !errors.Is(err, store.ErrDemasiadosAbiertos) {
		t.Fatalf("4º abierto: %v", err)
	}
	// Al confirmarse uno, se libera un cupo.
	if _, err := confirmar(db, primero.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.CrearPedido(ctx, nuevo(), pedidos.Cotizador("2026-11-14")); err != nil {
		t.Fatalf("tras confirmar uno: %v", err)
	}
	otro := nuevo()
	otro.Telefono = "+56900000099"
	if _, err := db.CrearPedido(ctx, otro, pedidos.Cotizador("2026-11-14")); err != nil {
		t.Fatalf("otro teléfono: %v", err)
	}
}

func TestOfertaVencidaNoCuentaEnFiltroNiOrden(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Tren grande", 20000, 10)
	of := int64(17000)
	p.PrecioOferta, p.OfertaHasta = &of, "2026-11-15"
	if err := db.ActualizarProducto(ctx, &p, time.Now()); err != nil {
		t.Fatal(err)
	}
	for hoy, want := range map[string]int{"2026-11-15": 1, "2026-11-16": 0} {
		_, total, err := db.ListarProductos(ctx, store.FiltroProductos{SoloOferta: true, Hoy: hoy, Orden: "precio_asc", Limite: 10})
		if err != nil || total != want {
			t.Errorf("hoy=%s ofertas=%d err=%v, quería %d", hoy, total, err, want)
		}
	}
}
func TestBusquedaConComodinesEInyeccion(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	sembrarProducto(t, db, "Telón de Navidad", 5000, 10)
	sembrarProducto(t, db, "Lazos navideños", 1000, 10)

	buscar := func(q string) int {
		t.Helper()
		lista, total, err := db.ListarProductos(ctx, store.FiltroProductos{Palabras: productos.PalabrasBusqueda(q, 5), Limite: 50})
		if err != nil {
			t.Fatalf("q=%q: %v", q, err)
		}
		if total != len(lista) {
			t.Fatalf("total %d != %d", total, len(lista))
		}
		return total
	}
	if n := buscar("telon"); n != 1 {
		t.Fatalf("telon -> %d", n)
	}
	if n := buscar("NAVIDEÑOS"); n != 1 {
		t.Fatalf("NAVIDEÑOS -> %d", n)
	}
	for _, q := range []string{"%", "_", `\`, "' OR '1'='1", "'; DROP TABLE productos; --", `" OR 1=1 --`, "%' --"} {
		if n := buscar(q); n != 0 {
			t.Errorf("q=%q devolvió %d resultados", q, n)
		}
	}
	if n := buscar(""); n != 2 {
		t.Fatalf("la tabla debe seguir intacta, hay %d", n)
	}
}
