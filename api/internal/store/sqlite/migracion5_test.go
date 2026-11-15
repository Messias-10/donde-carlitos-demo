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

// Una base v2 (0001–0004) con pedidos y la semilla cargada se actualiza a 2.1 (0005).
func TestMigracion0005SobreBaseV2(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "v2.db")
	raw, err := sql.Open("sqlite", "file:"+ruta+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec(`CREATE TABLE schema_migraciones (version TEXT PRIMARY KEY, aplicada_en TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for _, v := range []string{"0001_inicial.sql", "0002_stock_al_confirmar.sql", "0003_campos_qa.sql", "0004_panel_pedidos.sql"} {
		b, _ := fs.ReadFile(migrations.FS, v)
		if _, err := raw.Exec(string(b)); err != nil {
			t.Fatalf("%s: %v", v, err)
		}
		if _, err := raw.Exec(`INSERT INTO schema_migraciones VALUES (?, '2026-11-14T00:00:00.000Z')`, v); err != nil {
			t.Fatal(err)
		}
	}
	ts := "2026-11-14T00:00:00.000Z"
	for _, q := range []string{
		`INSERT INTO categorias (id, slug, nombre, creado_en) VALUES (1, 'navidad', 'Navidad', '` + ts + `')`,
		`INSERT INTO productos (id, slug, codigo, nombre, categoria_id, precio_detalle, stock, unidad_mayor, precio_mayor, mayor_minimo, creado_en, actualizado_en)
		 VALUES (1, 'esfera', 'P002', 'Esfera', 1, 500, 5, NULL, NULL, NULL, '` + ts + `', '` + ts + `'),
		        (2, 'scooter', 'P003', 'Scooter', 1, 3500, 34, 'embalaje de 16 (surtido)', 3000, 16, '` + ts + `', '` + ts + `'),
		        (3, 'tren', 'P012', 'Tren', 1, 20000, 37, NULL, NULL, NULL, '` + ts + `', '` + ts + `')`,
		`INSERT INTO usuarios (id, usuario, nombre, password_hash, rol, creado_en, actualizado_en) VALUES (2, 'regiones1', 'Regiones 1', 'h', 'vendedor', '` + ts + `', '` + ts + `')`,
		`INSERT INTO pedidos (id, secuencia, numero, token_seguimiento, canal, cliente_nombre, cliente_telefono, entrega, forma_pago,
		   subtotal, total, estado, estado_desde, creado_en, actualizado_en)
		 VALUES (1, 1, 'DC-0001', 'tokentokentokentoken01', 'santiago', 'A', '+56900000001', 'retiro', 'en_tienda',
		   1000, 1000, 'entregado', '2026-11-10T12:00:00.000Z', '` + ts + `', '` + ts + `')`,
		`INSERT INTO items_pedido (pedido_id, producto_id, codigo, descripcion, cantidad, precio_unitario, tipo_precio, subtotal)
		 VALUES (1, 1, 'P002', 'Esfera', 2, 500, 'detalle', 1000)`,
	} {
		if _, err := raw.Exec(q); err != nil {
			t.Fatalf("%v\n%s", err, q)
		}
	}
	_ = raw.Close()

	db, err := sqlite.Abrir(context.Background(), ruta) // aplica 0005
	if err != nil {
		t.Fatalf("migración 0005: %v", err)
	}
	defer db.Close()
	ctx := context.Background()
	p, err := db.PedidoPorID(ctx, 1)
	if err != nil || p.TokenVenceEn == nil || !p.TokenVenceEn.Equal(time.Date(2026, 12, 10, 12, 0, 0, 0, time.UTC)) || len(p.Items) != 1 {
		t.Fatalf("pedido migrado: %+v %v", p.TokenVenceEn, err)
	}
	esfera, _ := db.ProductoPorID(ctx, 1)
	scooter, _ := db.ProductoPorID(ctx, 2)
	tren, _ := db.ProductoPorID(ctx, 3)
	if !esfera.PrecioEsPromo || esfera.NotaPromo == "" || scooter.UnidadMayor != "embalaje surtido de 16" ||
		!tren.EnVitrina || tren.OrdenVitrina == nil || *tren.OrdenVitrina != 1 || tren.Stock != 300 {
		t.Fatalf("productos migrados: %+v / %+v / %+v", esfera, scooter, tren)
	}
	u, _ := db.UsuarioPorUsuario(ctx, "regiones1")
	if u.CanalPorDefecto != "regiones" || !u.DebeCambiarClave {
		t.Fatalf("usuario demo: %+v", u)
	}
}
