package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// CrearPedido, en UNA transacción (BEGIN IMMEDIATE):
//  1. lee los productos pedidos con su precio y stock actuales;
//  2. llama al cotizador del dominio (recalcula precios y total, valida stock);
//  3. toma el siguiente correlativo del contador;
//  4. inserta pedido, líneas e historial y descuenta (reserva) el stock.
//
// Si algo falla, se revierte todo y el correlativo no se consume.
func (d *DB) CrearPedido(ctx context.Context, n store.PedidoNuevo, cotizar store.Cotizador) (store.Pedido, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return store.Pedido{}, err
	}
	defer func() { _ = tx.Rollback() }()

	ids := make([]any, 0, len(n.Items))
	vistos := map[int64]bool{}
	for _, it := range n.Items {
		if !vistos[it.ProductoID] {
			vistos[it.ProductoID] = true
			ids = append(ids, it.ProductoID)
		}
	}
	productos := make(map[int64]store.Producto, len(ids))
	if len(ids) > 0 {
		q := "SELECT " + columnasProducto + desdeProducto +
			" WHERE p.eliminado_en IS NULL AND p.id IN (" + marcadores(len(ids)) + ")"
		rows, err := tx.QueryContext(ctx, q, ids...)
		if err != nil {
			return store.Pedido{}, err
		}
		for rows.Next() {
			p, err := escanearProducto(rows)
			if err != nil {
				rows.Close()
				return store.Pedido{}, err
			}
			productos[p.ID] = p
		}
		if err := rows.Close(); err != nil {
			return store.Pedido{}, err
		}
	}

	items, total, err := cotizar(productos, n.Items)
	if err != nil {
		return store.Pedido{}, err
	}

	var secuencia int64
	if err := tx.QueryRowContext(ctx,
		`UPDATE contadores SET valor = valor + 1 WHERE nombre = 'pedidos' RETURNING valor`).Scan(&secuencia); err != nil {
		return store.Pedido{}, fmt.Errorf("correlativo: %w", err)
	}
	numero := store.NumeroPedido(secuencia)
	ahora := ts(n.Ahora)

	res, err := tx.ExecContext(ctx, `
		INSERT INTO pedidos (secuencia, numero, cliente_nombre, cliente_telefono, entrega,
			region, comuna, agencia, comentario, total, estado, es_demo, creado_en, actualizado_en)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'nuevo', ?, ?, ?)`,
		secuencia, numero, n.Nombre, n.Telefono, n.Entrega,
		nullStr(n.Region), nullStr(n.Comuna), nullStr(n.Agencia), n.Comentario, total,
		boolInt(n.EsDemo), ahora, ahora)
	if err != nil {
		return store.Pedido{}, mapearError(err)
	}
	pedidoID, err := res.LastInsertId()
	if err != nil {
		return store.Pedido{}, err
	}

	for _, it := range items {
		// Descuento condicional: si otro pedido se llevó el stock, no queda negativo.
		r, err := tx.ExecContext(ctx,
			`UPDATE productos SET stock = stock - ?, actualizado_en = ? WHERE id = ? AND stock >= ?`,
			it.Cantidad, ahora, it.ProductoID, it.Cantidad)
		if err != nil {
			return store.Pedido{}, err
		}
		if filas, err := r.RowsAffected(); err != nil {
			return store.Pedido{}, err
		} else if filas != 1 {
			return store.Pedido{}, store.ErrStockInsuficiente
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO items_pedido (pedido_id, producto_id, codigo, nombre, cantidad, precio_unitario, tipo_precio, subtotal)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			pedidoID, it.ProductoID, it.Codigo, it.Nombre, it.Cantidad, it.PrecioUnitario, it.TipoPrecio, it.Subtotal); err != nil {
			return store.Pedido{}, mapearError(err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO historial_estados (pedido_id, estado_anterior, estado_nuevo, admin_id, nota, creado_en)
		VALUES (?, NULL, 'nuevo', NULL, '', ?)`, pedidoID, ahora); err != nil {
		return store.Pedido{}, err
	}
	if err := tx.Commit(); err != nil {
		return store.Pedido{}, err
	}
	return d.PedidoPorID(ctx, pedidoID)
}

const columnasPedido = `id, secuencia, numero, cliente_nombre, cliente_telefono, entrega,
	region, comuna, agencia, comentario, total, estado, es_demo, creado_en, actualizado_en`

func escanearPedido(s escaneable) (store.Pedido, error) {
	var (
		p                       store.Pedido
		region, comuna, agencia sql.NullString
		demo                    int
		creado, actualizado     string
	)
	err := s.Scan(&p.ID, &p.Secuencia, &p.Numero, &p.Nombre, &p.Telefono, &p.Entrega,
		&region, &comuna, &agencia, &p.Comentario, &p.Total, &p.Estado, &demo, &creado, &actualizado)
	if err != nil {
		return p, err
	}
	p.Region, p.Comuna, p.Agencia = region.String, comuna.String, agencia.String
	p.EsDemo = demo == 1
	p.CreadoEn, p.ActualizadoEn = parseTS(creado), parseTS(actualizado)
	return p, nil
}

type consultor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func cargarItems(ctx context.Context, q consultor, pedidoIDs []int64) (map[int64][]store.ItemPedido, error) {
	out := map[int64][]store.ItemPedido{}
	if len(pedidoIDs) == 0 {
		return out, nil
	}
	args := make([]any, len(pedidoIDs))
	for i, id := range pedidoIDs {
		args[i] = id
	}
	rows, err := q.QueryContext(ctx, `
		SELECT pedido_id, producto_id, codigo, nombre, cantidad, precio_unitario, tipo_precio, subtotal
		  FROM items_pedido WHERE pedido_id IN (`+marcadores(len(args))+`) ORDER BY id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pid int64
		var it store.ItemPedido
		if err := rows.Scan(&pid, &it.ProductoID, &it.Codigo, &it.Nombre, &it.Cantidad, &it.PrecioUnitario, &it.TipoPrecio, &it.Subtotal); err != nil {
			return nil, err
		}
		out[pid] = append(out[pid], it)
	}
	return out, rows.Err()
}

// PedidoPorID devuelve el pedido con líneas e historial.
func (d *DB) PedidoPorID(ctx context.Context, id int64) (store.Pedido, error) {
	p, err := escanearPedido(d.db.QueryRowContext(ctx, "SELECT "+columnasPedido+" FROM pedidos WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return p, store.ErrNoEncontrado
	}
	if err != nil {
		return p, err
	}
	items, err := cargarItems(ctx, d.db, []int64{id})
	if err != nil {
		return p, err
	}
	p.Items = items[id]

	rows, err := d.db.QueryContext(ctx, `
		SELECT h.estado_anterior, h.estado_nuevo, a.usuario, h.nota, h.creado_en
		  FROM historial_estados h LEFT JOIN admins a ON a.id = h.admin_id
		 WHERE h.pedido_id = ? ORDER BY h.creado_en, h.id`, id)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var anterior, usuario sql.NullString
		var c store.CambioRegistrado
		var creado string
		if err := rows.Scan(&anterior, &c.EstadoNuevo, &usuario, &c.Nota, &creado); err != nil {
			return p, err
		}
		c.EstadoAnterior, c.AdminUsuario, c.CreadoEn = anterior.String, usuario.String, parseTS(creado)
		p.Historial = append(p.Historial, c)
	}
	return p, rows.Err()
}

// ListarPedidos devuelve pedidos (con líneas, sin historial), del más nuevo al más antiguo.
func (d *DB) ListarPedidos(ctx context.Context, f store.FiltroPedidos) ([]store.Pedido, int, error) {
	where := ""
	var args []any
	if f.Estado != "" {
		where = " WHERE estado = ?"
		args = append(args, f.Estado)
	}
	var total int
	if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pedidos"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	limite := f.Limite
	if limite <= 0 {
		limite = 20
	}
	rows, err := d.db.QueryContext(ctx,
		"SELECT "+columnasPedido+" FROM pedidos"+where+" ORDER BY creado_en DESC, id DESC LIMIT ? OFFSET ?",
		append(args, limite, f.Desplazamiento)...)
	if err != nil {
		return nil, 0, err
	}
	var out []store.Pedido
	var ids []int64
	for rows.Next() {
		p, err := escanearPedido(rows)
		if err != nil {
			rows.Close()
			return nil, 0, err
		}
		out = append(out, p)
		ids = append(ids, p.ID)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	items, err := cargarItems(ctx, d.db, ids)
	if err != nil {
		return nil, 0, err
	}
	for i := range out {
		out[i].Items = items[out[i].ID]
	}
	return out, total, nil
}

// CambiarEstado aplica una transición validada por el dominio, guarda el
// historial y, si corresponde (cancelación), devuelve el stock reservado.
func (d *DB) CambiarEstado(ctx context.Context, c store.CambioEstado, validar store.ValidadorCambio) (store.Pedido, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return store.Pedido{}, err
	}
	defer func() { _ = tx.Rollback() }()

	actual, err := escanearPedido(tx.QueryRowContext(ctx, "SELECT "+columnasPedido+" FROM pedidos WHERE id = ?", c.PedidoID))
	if errors.Is(err, sql.ErrNoRows) {
		return store.Pedido{}, store.ErrNoEncontrado
	}
	if err != nil {
		return store.Pedido{}, err
	}
	restaurar, err := validar(actual)
	if err != nil {
		return store.Pedido{}, err
	}
	ahora := ts(c.Ahora)
	res, err := tx.ExecContext(ctx, `UPDATE pedidos SET estado = ?, actualizado_en = ? WHERE id = ? AND estado = ?`,
		c.Nuevo, ahora, c.PedidoID, actual.Estado)
	if err != nil {
		return store.Pedido{}, err
	}
	if n, err := res.RowsAffected(); err != nil {
		return store.Pedido{}, err
	} else if n != 1 {
		return store.Pedido{}, store.ErrConflicto
	}

	if restaurar {
		items, err := cargarItems(ctx, tx, []int64{c.PedidoID})
		if err != nil {
			return store.Pedido{}, err
		}
		for _, it := range items[c.PedidoID] {
			if _, err := tx.ExecContext(ctx, `UPDATE productos SET stock = stock + ?, actualizado_en = ? WHERE id = ?`,
				it.Cantidad, ahora, it.ProductoID); err != nil {
				return store.Pedido{}, err
			}
		}
	}

	var admin any
	if c.AdminID > 0 {
		admin = c.AdminID
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO historial_estados (pedido_id, estado_anterior, estado_nuevo, admin_id, nota, creado_en)
		VALUES (?, ?, ?, ?, ?, ?)`, c.PedidoID, actual.Estado, c.Nuevo, admin, strings.TrimSpace(c.Nota), ahora); err != nil {
		return store.Pedido{}, err
	}
	if err := tx.Commit(); err != nil {
		return store.Pedido{}, err
	}
	return d.PedidoPorID(ctx, c.PedidoID)
}
