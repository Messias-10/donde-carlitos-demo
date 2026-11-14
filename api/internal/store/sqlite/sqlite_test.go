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
			ped, err := db.CrearPedido(context.Background(), pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1}), pedidos.Cotizar)
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
	if act.Stock != 100000-n {
		t.Fatalf("stock = %d", act.Stock)
	}
}

func TestStockNoSeSobrevendeEnConcurrencia(t *testing.T) {
	db := abrir(t)
	p := sembrarProducto(t, db, "Tren navideño", 17000, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex
	ok, rechazados := 0, 0
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := db.CrearPedido(context.Background(), pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1}), pedidos.Cotizar)
			mu.Lock()
			defer mu.Unlock()
			var nd *pedidos.ErrNoDisponibles
			switch {
			case err == nil:
				ok++
			case errors.As(err, &nd) || errors.Is(err, store.ErrStockInsuficiente):
				rechazados++
			default:
				t.Errorf("error inesperado: %v", err)
			}
		}()
	}
	wg.Wait()
	if ok != 5 || rechazados != 15 {
		t.Fatalf("ok=%d rechazados=%d", ok, rechazados)
	}
	// Un pedido rechazado no consume número.
	ped, err := db.CrearPedido(context.Background(), pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 1}), pedidos.Cotizar)
	if err == nil || ped.Numero != "" {
		t.Fatal("debió rechazarse por stock")
	}
}

func TestCambioEstadoHistorialYStock(t *testing.T) {
	db := abrir(t)
	ctx := context.Background()
	p := sembrarProducto(t, db, "Calcetas navideñas", 1000, 10)
	ped, err := db.CrearPedido(ctx, pedidoNuevo(store.ItemSolicitado{ProductoID: p.ID, Cantidad: 4}), pedidos.Cotizar)
	if err != nil {
		t.Fatal(err)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 6 {
		t.Fatalf("stock reservado = %d", act.Stock)
	}
	// inválida: nuevo -> entregado
	_, err = db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "entregado", Ahora: time.Now()}, pedidos.Validador("entregado"))
	var ti *pedidos.ErrTransicionInvalida
	if !errors.As(err, &ti) {
		t.Fatalf("quería transición inválida: %v", err)
	}
	// válidas
	if _, err := db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "confirmado", Ahora: time.Now()}, pedidos.Validador("confirmado")); err != nil {
		t.Fatal(err)
	}
	fin, err := db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "cancelado", Nota: "no confirmó", Ahora: time.Now()}, pedidos.Validador("cancelado"))
	if err != nil {
		t.Fatal(err)
	}
	if fin.Estado != "cancelado" || len(fin.Historial) != 3 || fin.Historial[2].EstadoAnterior != "confirmado" || fin.Historial[2].Nota != "no confirmó" {
		t.Fatalf("historial: %+v", fin.Historial)
	}
	if act, _ := db.ProductoPorID(ctx, p.ID); act.Stock != 10 {
		t.Fatalf("stock tras cancelar = %d (debía volver a 10)", act.Stock)
	}
	// cancelado es final
	if _, err := db.CambiarEstado(ctx, store.CambioEstado{PedidoID: ped.ID, Nuevo: "confirmado", Ahora: time.Now()}, pedidos.Validador("confirmado")); !errors.As(err, &ti) {
		t.Fatalf("cancelado debe ser final: %v", err)
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
