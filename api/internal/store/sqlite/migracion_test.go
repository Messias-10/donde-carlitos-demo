package sqlite_test

import (
	"context"
	"database/sql"
	"io/fs"
	"path/filepath"
	"testing"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store/sqlite"
	"github.com/Messias-10/donde-carlitos-demo/api/migrations"
)

// Una base v1.1 (0001–0003 con admin, sesión y pedidos viejos) se actualiza a
// v2 sin perder el catálogo ni el admin, y con el stock devuelto.
func TestMigracionDesdeV11(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "v11.db")
	raw, err := sql.Open("sqlite", "file:"+ruta+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec(`CREATE TABLE schema_migraciones (version TEXT PRIMARY KEY, aplicada_en TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for _, v := range []string{"0001_inicial.sql", "0002_stock_al_confirmar.sql", "0003_campos_qa.sql"} {
		b, err := fs.ReadFile(migrations.FS, v)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := raw.Exec(string(b)); err != nil {
			t.Fatalf("%s: %v", v, err)
		}
		if _, err := raw.Exec(`INSERT INTO schema_migraciones VALUES (?, '2026-11-14T00:00:00.000Z')`, v); err != nil {
			t.Fatal(err)
		}
	}
	ahora := "2026-11-14T00:00:00.000Z"
	for _, q := range []string{
		`INSERT INTO categorias (id, slug, nombre, creado_en) VALUES (1, 'navidad', 'Navidad', '` + ahora + `')`,
		`INSERT INTO productos (id, slug, codigo, nombre, categoria_id, precio_detalle, stock, creado_en, actualizado_en)
		 VALUES (1, 'tren', 'P012', 'Tren', 1, 20000, 8, '` + ahora + `', '` + ahora + `')`,
		`INSERT INTO admins (id, usuario, password_hash, creado_en, actualizado_en) VALUES (1, 'admin', 'h', '` + ahora + `', '` + ahora + `')`,
		`INSERT INTO sesiones (token_hash, admin_id, creado_en, expira_en, ultimo_uso) VALUES ('s1', 1, '` + ahora + `', '2099-01-01T00:00:00.000Z', '` + ahora + `')`,
		`INSERT INTO pedidos (id, secuencia, numero, cliente_nombre, cliente_telefono, entrega, total, estado, creado_en, actualizado_en)
		 VALUES (1, 1, 'DC-0001', 'A', '+56900000001', 'retiro', 40000, 'confirmado', '` + ahora + `', '` + ahora + `'),
		        (2, 2, 'DC-0002', 'B', '+56900000002', 'retiro', 20000, 'nuevo', '` + ahora + `', '` + ahora + `')`,
		`INSERT INTO items_pedido (pedido_id, producto_id, codigo, nombre, cantidad, precio_unitario, tipo_precio, subtotal)
		 VALUES (1, 1, 'P012', 'Tren', 2, 20000, 'detalle', 40000), (2, 1, 'P012', 'Tren', 1, 20000, 'detalle', 20000)`,
		`UPDATE contadores SET valor = 2 WHERE nombre = 'pedidos'`,
	} {
		if _, err := raw.Exec(q); err != nil {
			t.Fatalf("%v\n%s", err, q)
		}
	}
	_ = raw.Close()

	db, err := sqlite.Abrir(context.Background(), ruta) // aplica 0004
	if err != nil {
		t.Fatalf("migración 0004: %v", err)
	}
	defer db.Close()
	ctx := context.Background()
	if p, err := db.ProductoPorID(ctx, 1); err != nil || p.Stock != 10 {
		t.Fatalf("stock devuelto del pedido confirmado (8+2): %+v %v", p.Stock, err)
	}
	u, err := db.UsuarioPorUsuario(ctx, "admin")
	if err != nil || u.Rol != "admin" || !u.Activo || u.Nombre != "Administrador" {
		t.Fatalf("admin migrado: %+v %v", u, err)
	}
	if s, err := db.SesionPorHash(ctx, "s1"); err != nil || s.UsuarioID != 1 {
		t.Fatalf("sesión migrada: %v", err)
	}
	if vacio, _ := db.PedidosVacio(ctx); !vacio {
		t.Fatal("los pedidos v1 (ficticios) se borran")
	}
	if _, err := db.SincronizarAdmin(ctx, "admin", "h2", time.Now()); err != nil {
		t.Fatal(err)
	}
}
