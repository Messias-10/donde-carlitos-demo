package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// cargarProductos lee (dentro de la transacción) los productos de catálogo usados.
func cargarProductos(ctx context.Context, q consultor, items []store.ItemSolicitado) (map[int64]store.Producto, error) {
	var ids []any
	vistos := map[int64]bool{}
	for _, it := range items {
		if it.ProductoID > 0 && !vistos[it.ProductoID] {
			vistos[it.ProductoID] = true
			ids = append(ids, it.ProductoID)
		}
	}
	out := make(map[int64]store.Producto, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := q.QueryContext(ctx, "SELECT "+columnasProducto+desdeProducto+
		" WHERE p.eliminado_en IS NULL AND p.id IN ("+marcadores(len(ids))+")", ids...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p, err := escanearProducto(rows)
		if err != nil {
			return nil, err
		}
		out[p.ID] = p
	}
	return out, rows.Err()
}

func insertarItems(ctx context.Context, tx *sql.Tx, pedidoID int64, items []store.ItemPedido) error {
	for _, it := range items {
		var prod, cod any
		if it.ProductoID > 0 {
			prod, cod = it.ProductoID, nullStr(it.Codigo)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO items_pedido (pedido_id, producto_id, codigo, descripcion, cantidad, precio_unitario, tipo_precio, subtotal)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			pedidoID, prod, cod, it.Descripcion, it.Cantidad, nullPtr(it.PrecioUnitario), it.TipoPrecio, nullPtr(it.Subtotal)); err != nil {
			return mapearError(err)
		}
	}
	return nil
}

// calcularTotal: con todas las líneas con precio, total = subtotal - descuento
// + costo de envío (si viene un total escrito, debe coincidir). Si hay líneas
// sin precio, manda el total escrito por quien atiende (obligatorio).
func calcularTotal(items []store.ItemPedido, subtotal int64, d store.DatosPedido) (int64, bool, error) {
	sinPrecio := false
	for _, it := range items {
		sinPrecio = sinPrecio || it.PrecioUnitario == nil
	}
	if sinPrecio {
		if d.TotalManual == nil || *d.TotalManual < 0 {
			return 0, true, store.ErrTotalNoCoincide
		}
		return *d.TotalManual, true, nil
	}
	total := subtotal - d.Descuento + d.CostoEnvio
	if total < 0 {
		return 0, false, store.ErrTotalNegativo
	}
	if d.TotalManual != nil && *d.TotalManual != total {
		return 0, false, store.ErrTotalNoCoincide
	}
	return total, false, nil
}

func nullPtr(p *int64) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullID(id int64) any {
	if id <= 0 {
		return nil
	}
	return id
}

// CrearPedido, en UNA transacción (BEGIN IMMEDIATE): cotiza con los precios
// actuales, toma el siguiente correlativo e inserta pedido, líneas e historial.
// No toca el stock (se descuenta al pasar a "pagado").
func (d *DB) CrearPedido(ctx context.Context, n store.PedidoNuevo, cotizar store.Cotizador) (store.Pedido, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return store.Pedido{}, err
	}
	defer func() { _ = tx.Rollback() }()

	prods, err := cargarProductos(ctx, tx, n.Items)
	if err != nil {
		return store.Pedido{}, err
	}
	items, subtotal, err := cotizar(prods, n.Items)
	if err != nil {
		return store.Pedido{}, err
	}
	var sec int64
	if err := tx.QueryRowContext(ctx,
		`UPDATE contadores SET valor = valor + 1 WHERE nombre = 'pedidos' RETURNING valor`).Scan(&sec); err != nil {
		return store.Pedido{}, fmt.Errorf("correlativo: %w", err)
	}
	total, manual, err := calcularTotal(items, subtotal, n.DatosPedido)
	if err != nil {
		return store.Pedido{}, err
	}
	numero := store.NumeroPedido(sec)
	ahora := ts(n.Ahora)
	res, err := tx.ExecContext(ctx, `
		INSERT INTO pedidos (secuencia, numero, token_seguimiento, canal, cliente_nombre, cliente_telefono, busqueda,
			entrega, region, comuna, agencia, forma_pago, subtotal, descuento, costo_envio, envio_por_pagar,
			total, total_manual, nota_interna, estado, estado_desde, registrado_por, es_demo, creado_en, actualizado_en)
		VALUES (`+marcadores(25)+`)`,
		sec, numero, n.TokenSeguimiento, n.Canal, n.Nombre, n.Telefono, strings.ToLower(numero)+" "+n.Busqueda,
		n.Entrega, nullStr(n.Region), nullStr(n.Comuna), nullStr(n.Agencia), n.FormaPago, subtotal, n.Descuento, n.CostoEnvio,
		boolInt(n.EnvioPorPagar), total, boolInt(manual), n.NotaInterna,
		n.EstadoInicial, ahora, nullID(n.RegistradoPor), boolInt(n.EsDemo), ahora, ahora)
	if err != nil {
		return store.Pedido{}, mapearError(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return store.Pedido{}, err
	}
	if err := insertarItems(ctx, tx, id, items); err != nil {
		return store.Pedido{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO historial_estados (pedido_id, estado_anterior, estado_nuevo, usuario_id, nota, creado_en)
		VALUES (?, NULL, ?, ?, '', ?)`, id, n.EstadoInicial, nullID(n.RegistradoPor), ahora); err != nil {
		return store.Pedido{}, err
	}
	if err := tx.Commit(); err != nil {
		return store.Pedido{}, err
	}
	return d.PedidoPorID(ctx, id)
}

// EditarPedido reemplaza datos y líneas (precios recalculados y congelados de
// nuevo). Solo en cotizado o esperando_pago, sin stock descontado.
func (d *DB) EditarPedido(ctx context.Context, id int64, dp store.DatosPedido, cotizar store.Cotizador, ahora time.Time) (store.Pedido, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return store.Pedido{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var estado, numero string
	var descontado int
	err = tx.QueryRowContext(ctx, `SELECT estado, numero, stock_descontado FROM pedidos WHERE id = ?`, id).Scan(&estado, &numero, &descontado)
	if errors.Is(err, sql.ErrNoRows) {
		return store.Pedido{}, store.ErrNoEncontrado
	}
	if err != nil {
		return store.Pedido{}, err
	}
	if (estado != "cotizado" && estado != "esperando_pago") || descontado == 1 {
		return store.Pedido{}, store.ErrNoEditable
	}
	prods, err := cargarProductos(ctx, tx, dp.Items)
	if err != nil {
		return store.Pedido{}, err
	}
	items, subtotal, err := cotizar(prods, dp.Items)
	if err != nil {
		return store.Pedido{}, err
	}
	total, manual, err := calcularTotal(items, subtotal, dp)
	if err != nil {
		return store.Pedido{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE pedidos SET canal = ?, cliente_nombre = ?, cliente_telefono = ?, busqueda = ?, entrega = ?,
			region = ?, comuna = ?, agencia = ?, forma_pago = ?, subtotal = ?, descuento = ?, costo_envio = ?,
			envio_por_pagar = ?, total = ?, total_manual = ?, nota_interna = ?, actualizado_en = ?
		WHERE id = ?`,
		dp.Canal, dp.Nombre, dp.Telefono, strings.ToLower(numero)+" "+dp.Busqueda, dp.Entrega,
		nullStr(dp.Region), nullStr(dp.Comuna), nullStr(dp.Agencia), dp.FormaPago, subtotal, dp.Descuento, dp.CostoEnvio,
		boolInt(dp.EnvioPorPagar), total, boolInt(manual), dp.NotaInterna, ts(ahora), id); err != nil {
		return store.Pedido{}, mapearError(err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM items_pedido WHERE pedido_id = ?`, id); err != nil {
		return store.Pedido{}, err
	}
	if err := insertarItems(ctx, tx, id, items); err != nil {
		return store.Pedido{}, err
	}
	if err := tx.Commit(); err != nil {
		return store.Pedido{}, err
	}
	return d.PedidoPorID(ctx, id)
}

const columnasPedido = `p.id, p.secuencia, p.numero, p.token_seguimiento, p.canal, p.cliente_nombre, p.cliente_telefono,
	p.entrega, p.region, p.comuna, p.agencia, p.n_seguimiento, p.forma_pago, p.subtotal, p.descuento, p.costo_envio,
	p.envio_por_pagar, p.total, p.total_manual, p.motivo_cancelacion, p.nota_interna, p.estado, p.estado_desde, p.stock_descontado, p.pago_medio, p.pago_monto, p.pago_fecha,
	p.pago_verificado_por, uv.nombre, p.registrado_por, ur.nombre, p.es_demo, p.creado_en, p.actualizado_en,
	p.token_vence_en`

const desdePedido = ` FROM pedidos p
	LEFT JOIN usuarios uv ON uv.id = p.pago_verificado_por
	LEFT JOIN usuarios ur ON ur.id = p.registrado_por `

func escanearPedido(s escaneable) (store.Pedido, error) {
	var (
		p                                        store.Pedido
		region, comuna, agencia, nseg            sql.NullString
		pagoMedio, pagoFecha, uvNombre, urNombre sql.NullString
		pagoMonto, uvID, urID                    sql.NullInt64
		descontado, demo, porPagar, manual       int
		estadoDesde, creado, actualizado         string
		vence                                    sql.NullString
	)
	err := s.Scan(&p.ID, &p.Secuencia, &p.Numero, &p.TokenSeguimiento, &p.Canal, &p.Nombre, &p.Telefono,
		&p.Entrega, &region, &comuna, &agencia, &nseg, &p.FormaPago, &p.Subtotal, &p.Descuento, &p.CostoEnvio,
		&porPagar, &p.Total, &manual, &p.MotivoCancelacion, &p.NotaInterna, &p.Estado, &estadoDesde, &descontado, &pagoMedio, &pagoMonto, &pagoFecha,
		&uvID, &uvNombre, &urID, &urNombre, &demo, &creado, &actualizado, &vence)
	if err != nil {
		return p, err
	}
	p.Region, p.Comuna, p.Agencia, p.NSeguimiento = region.String, comuna.String, agencia.String, nseg.String
	p.StockDescontado, p.EsDemo = descontado == 1, demo == 1
	p.EnvioPorPagar, p.TotalManual = porPagar == 1, manual == 1
	if vence.Valid {
		v := parseTS(vence.String)
		p.TokenVenceEn = &v
	}
	p.EstadoDesde, p.CreadoEn, p.ActualizadoEn = parseTS(estadoDesde), parseTS(creado), parseTS(actualizado)
	if pagoMedio.Valid {
		p.Pago = &store.Pago{Medio: pagoMedio.String, Monto: pagoMonto.Int64, Fecha: pagoFecha.String}
		if uvID.Valid {
			p.Pago.VerificadoPor = &store.UsuarioRef{ID: uvID.Int64, Nombre: uvNombre.String}
		}
	}
	if urID.Valid {
		p.RegistradoPor = &store.UsuarioRef{ID: urID.Int64, Nombre: urNombre.String}
	}
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
		SELECT pedido_id, producto_id, codigo, descripcion, cantidad, precio_unitario, tipo_precio, subtotal
		  FROM items_pedido WHERE pedido_id IN (`+marcadores(len(args))+`) ORDER BY id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pid int64
		var prod, precio, sub sql.NullInt64
		var cod sql.NullString
		var it store.ItemPedido
		if err := rows.Scan(&pid, &prod, &cod, &it.Descripcion, &it.Cantidad, &precio, &it.TipoPrecio, &sub); err != nil {
			return nil, err
		}
		it.ProductoID, it.Codigo = prod.Int64, cod.String
		if precio.Valid {
			v, s := precio.Int64, sub.Int64
			it.PrecioUnitario, it.Subtotal = &v, &s
		}
		out[pid] = append(out[pid], it)
	}
	return out, rows.Err()
}

func cargarHistorial(ctx context.Context, q consultor, pedidoIDs []int64) (map[int64][]store.CambioRegistrado, error) {
	out := map[int64][]store.CambioRegistrado{}
	if len(pedidoIDs) == 0 {
		return out, nil
	}
	args := make([]any, len(pedidoIDs))
	for i, id := range pedidoIDs {
		args[i] = id
	}
	rows, err := q.QueryContext(ctx, `
		SELECT h.pedido_id, h.estado_anterior, h.estado_nuevo, h.usuario_id, u.nombre, h.nota, h.creado_en
		  FROM historial_estados h LEFT JOIN usuarios u ON u.id = h.usuario_id
		 WHERE h.pedido_id IN (`+marcadores(len(args))+`) ORDER BY h.creado_en, h.id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pid int64
		var anterior, nombre sql.NullString
		var uid sql.NullInt64
		var c store.CambioRegistrado
		var creado string
		if err := rows.Scan(&pid, &anterior, &c.EstadoNuevo, &uid, &nombre, &c.Nota, &creado); err != nil {
			return nil, err
		}
		c.EstadoAnterior, c.CreadoEn = anterior.String, parseTS(creado)
		if uid.Valid {
			c.Usuario = &store.UsuarioRef{ID: uid.Int64, Nombre: nombre.String}
		}
		out[pid] = append(out[pid], c)
	}
	return out, rows.Err()
}

func (d *DB) completar(ctx context.Context, ps []store.Pedido, conHistorial bool) error {
	ids := make([]int64, len(ps))
	for i := range ps {
		ids[i] = ps[i].ID
	}
	items, err := cargarItems(ctx, d.db, ids)
	if err != nil {
		return err
	}
	var hist map[int64][]store.CambioRegistrado
	if conHistorial {
		if hist, err = cargarHistorial(ctx, d.db, ids); err != nil {
			return err
		}
	}
	for i := range ps {
		ps[i].Items = items[ps[i].ID]
		if conHistorial {
			ps[i].Historial = hist[ps[i].ID]
		}
	}
	return nil
}

func (d *DB) pedidoDonde(ctx context.Context, cond string, arg any) (store.Pedido, error) {
	p, err := escanearPedido(d.db.QueryRowContext(ctx, "SELECT "+columnasPedido+desdePedido+" WHERE "+cond, arg))
	if errors.Is(err, sql.ErrNoRows) {
		return p, store.ErrNoEncontrado
	}
	if err != nil {
		return p, err
	}
	ps := []store.Pedido{p}
	if err := d.completar(ctx, ps, true); err != nil {
		return p, err
	}
	return ps[0], nil
}

// PedidoPorID devuelve el pedido con líneas e historial.
func (d *DB) PedidoPorID(ctx context.Context, id int64) (store.Pedido, error) {
	return d.pedidoDonde(ctx, "p.id = ?", id)
}

// PedidoPorToken busca por token de seguimiento.
func (d *DB) PedidoPorToken(ctx context.Context, token string) (store.Pedido, error) {
	return d.pedidoDonde(ctx, "p.token_seguimiento = ?", token)
}

// ListarPedidos aplica filtros y paginación (del más nuevo al más antiguo).
func (d *DB) ListarPedidos(ctx context.Context, f store.FiltroPedidos) ([]store.Pedido, int, error) {
	conds := []string{"1 = 1"}
	var args []any
	if f.Canal != "" {
		conds = append(conds, "p.canal = ?")
		args = append(args, f.Canal)
	}
	if len(f.Estados) > 0 {
		conds = append(conds, "p.estado IN ("+marcadores(len(f.Estados))+")")
		for _, e := range f.Estados {
			args = append(args, e)
		}
	}
	for _, palabra := range f.Palabras {
		conds = append(conds, `p.busqueda LIKE ? ESCAPE '\'`)
		args = append(args, "%"+escaparLike(palabra)+"%")
	}
	if f.Telefono != "" {
		conds = append(conds, "p.cliente_telefono = ?")
		args = append(args, f.Telefono)
	}
	if !f.CreadoDesde.IsZero() {
		conds = append(conds, "p.creado_en >= ?")
		args = append(args, ts(f.CreadoDesde))
	}
	if !f.CreadoHasta.IsZero() {
		conds = append(conds, "p.creado_en < ?")
		args = append(args, ts(f.CreadoHasta))
	}
	if !f.EstadoDesdeAntes.IsZero() {
		conds = append(conds, "p.estado_desde < ?")
		args = append(args, ts(f.EstadoDesdeAntes))
	}
	where := " WHERE " + strings.Join(conds, " AND ")
	var total int
	if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pedidos p"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	limite := f.Limite
	if limite <= 0 {
		limite = 20
	}
	rows, err := d.db.QueryContext(ctx, "SELECT "+columnasPedido+desdePedido+where+
		" ORDER BY p.creado_en DESC, p.id DESC LIMIT ? OFFSET ?", append(args, limite, f.Desplazamiento)...)
	if err != nil {
		return nil, 0, err
	}
	var out []store.Pedido
	for rows.Next() {
		p, err := escanearPedido(rows)
		if err != nil {
			rows.Close()
			return nil, 0, err
		}
		out = append(out, p)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := d.completar(ctx, out, f.ConHistorial); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// CambiarEstado aplica una transición validada por el dominio y su efecto en
// el stock: al pasar a pagado descuenta (UPDATE condicional; si no alcanza,
// revierte todo y devuelve *store.ErrSinStock) y al cancelar un pedido con stock
// descontado lo devuelve. Guarda pago, agencia/n.º de seguimiento e historial.
func (d *DB) CambiarEstado(ctx context.Context, c store.CambioEstado, validar store.ValidadorCambio) (store.Pedido, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return store.Pedido{}, err
	}
	defer func() { _ = tx.Rollback() }()

	actual, err := escanearPedido(tx.QueryRowContext(ctx, "SELECT "+columnasPedido+desdePedido+" WHERE p.id = ?", c.PedidoID))
	if errors.Is(err, sql.ErrNoRows) {
		return store.Pedido{}, store.ErrNoEncontrado
	}
	if err != nil {
		return store.Pedido{}, err
	}
	efecto, err := validar(actual, c)
	if err != nil {
		return store.Pedido{}, err
	}
	ahora := ts(c.Ahora)
	var avisoStock []int64

	if efecto != store.StockSinCambio {
		items, err := cargarItems(ctx, tx, []int64{c.PedidoID})
		if err != nil {
			return store.Pedido{}, err
		}
		var faltan, avisos []int64
		for _, it := range items[c.PedidoID] {
			if it.ProductoID == 0 {
				continue // línea libre: no hay stock
			}
			if efecto == store.StockDevolver {
				if _, err := tx.ExecContext(ctx, `UPDATE productos SET stock = stock + ?, actualizado_en = ? WHERE id = ?`,
					it.Cantidad, ahora, it.ProductoID); err != nil {
					return store.Pedido{}, err
				}
				continue
			}
			r, err := tx.ExecContext(ctx, `UPDATE productos SET stock = stock - ?, actualizado_en = ? WHERE id = ? AND stock >= ?`,
				it.Cantidad, ahora, it.ProductoID, it.Cantidad)
			if err != nil {
				return store.Pedido{}, err
			}
			if n, err := r.RowsAffected(); err != nil {
				return store.Pedido{}, err
			} else if n != 1 {
				if c.StockEstricto {
					faltan = append(faltan, it.ProductoID)
					continue
				}
				// No estricto (demo): se acepta, el stock queda en 0 y se avisa.
				if _, err := tx.ExecContext(ctx, `UPDATE productos SET stock = 0, actualizado_en = ? WHERE id = ?`, ahora, it.ProductoID); err != nil {
					return store.Pedido{}, err
				}
				avisos = append(avisos, it.ProductoID)
			}
		}
		if len(faltan) > 0 {
			return store.Pedido{}, &store.ErrSinStock{ProductoIDs: faltan}
		}
		avisoStock = avisos
	}

	sets := []string{"estado = ?", "estado_desde = ?", "actualizado_en = ?"}
	args := []any{c.Nuevo, ahora, ahora}
	switch efecto {
	case store.StockDescontar:
		sets = append(sets, "stock_descontado = 1")
	case store.StockDevolver:
		sets = append(sets, "stock_descontado = 0")
	}
	if c.Pago != nil {
		sets = append(sets, "pago_medio = ?", "pago_monto = ?", "pago_fecha = ?", "pago_verificado_por = ?")
		args = append(args, c.Pago.Medio, c.Pago.Monto, c.Pago.Fecha, nullID(c.UsuarioID))
	}
	if c.Agencia != "" {
		sets = append(sets, "agencia = ?")
		args = append(args, c.Agencia)
	}
	if c.NSeguimiento != "" {
		sets = append(sets, "n_seguimiento = ?")
		args = append(args, c.NSeguimiento)
	}
	if c.Comuna != "" {
		sets = append(sets, "comuna = ?")
		args = append(args, c.Comuna)
	}
	if c.Motivo != "" {
		sets = append(sets, "motivo_cancelacion = ?")
		args = append(args, c.Motivo)
	}
	if c.Nuevo == "entregado" || c.Nuevo == "cancelado" {
		sets = append(sets, "token_vence_en = ?")
		args = append(args, ts(c.Ahora.Add(30*24*time.Hour)))
	}
	args = append(args, c.PedidoID, actual.Estado)
	res, err := tx.ExecContext(ctx, "UPDATE pedidos SET "+strings.Join(sets, ", ")+" WHERE id = ? AND estado = ?", args...)
	if err != nil {
		return store.Pedido{}, mapearError(err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return store.Pedido{}, err
	} else if n != 1 {
		return store.Pedido{}, store.ErrConflicto
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO historial_estados (pedido_id, estado_anterior, estado_nuevo, usuario_id, nota, creado_en)
		VALUES (?, ?, ?, ?, ?, ?)`, c.PedidoID, actual.Estado, c.Nuevo, nullID(c.UsuarioID), strings.TrimSpace(c.Nota), ahora); err != nil {
		return store.Pedido{}, err
	}
	if err := tx.Commit(); err != nil {
		return store.Pedido{}, err
	}
	p, err := d.PedidoPorID(ctx, c.PedidoID)
	p.AvisoStock = avisoStock
	return p, err
}

// RegenerarToken pone un token nuevo. Si el pedido ya terminó, vence en 30 días desde ahora.
func (d *DB) RegenerarToken(ctx context.Context, id int64, token string, ahora time.Time) (store.Pedido, error) {
	res, err := d.db.ExecContext(ctx, `
		UPDATE pedidos SET token_seguimiento = ?, actualizado_en = ?,
		       token_vence_en = CASE WHEN estado IN ('entregado', 'cancelado') THEN ? ELSE NULL END
		 WHERE id = ?`, token, ts(ahora), ts(ahora.Add(30*24*time.Hour)), id)
	if err != nil {
		return store.Pedido{}, mapearError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return store.Pedido{}, store.ErrNoEncontrado
	}
	return d.PedidoPorID(ctx, id)
}

// ResumenUsuarios: por usuario, pedidos registrados (por creado_en) y despachos
// (por la fecha del cambio a despachado) en [desde, hasta).
func (d *DB) ResumenUsuarios(ctx context.Context, desde, hasta time.Time) ([]store.FilaUsuario, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT u.id, u.nombre,
		       (SELECT COUNT(*) FROM pedidos p WHERE p.registrado_por = u.id AND p.creado_en >= ? AND p.creado_en < ?),
		       (SELECT COUNT(*) FROM historial_estados h WHERE h.usuario_id = u.id AND h.estado_nuevo = 'despachado'
		                                                   AND h.creado_en >= ? AND h.creado_en < ?)
		  FROM usuarios u
		 ORDER BY u.nombre`, ts(desde), ts(hasta), ts(desde), ts(hasta))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []store.FilaUsuario
	for rows.Next() {
		var f store.FilaUsuario
		if err := rows.Scan(&f.UsuarioID, &f.Nombre, &f.Registrados, &f.Despachados); err != nil {
			return nil, err
		}
		if f.Registrados > 0 || f.Despachados > 0 {
			out = append(out, f)
		}
	}
	return out, rows.Err()
}

// ResumenPedidos agrupa por canal, estado, entrega y agencia en [desde, hasta).
func (d *DB) ResumenPedidos(ctx context.Context, desde, hasta time.Time) ([]store.FilaResumen, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT canal, estado, entrega, COALESCE(agencia, ''), COALESCE(comuna, ''), pago_medio IS NOT NULL,
		       COUNT(*), COALESCE(SUM(total), 0)
		  FROM pedidos
		 WHERE creado_en >= ? AND creado_en < ?
		 GROUP BY canal, estado, entrega, COALESCE(agencia, ''), COALESCE(comuna, ''), pago_medio IS NOT NULL
		 ORDER BY canal, estado`, ts(desde), ts(hasta))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []store.FilaResumen
	for rows.Next() {
		var f store.FilaResumen
		if err := rows.Scan(&f.Canal, &f.Estado, &f.Entrega, &f.Agencia, &f.Comuna, &f.Pagado, &f.Pedidos, &f.Monto); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// FijarCorrelativo pone el contador (la semilla lo usa con la tabla vacía para
// que los pedidos DEMO tengan los números del guion: DC-0101…).
func (d *DB) FijarCorrelativo(ctx context.Context, valor int64) error {
	_, err := d.db.ExecContext(ctx, `UPDATE contadores SET valor = ? WHERE nombre = 'pedidos'`, valor)
	return err
}

// PedidosVacio indica si no hay pedidos (para sembrar los DEMO).
func (d *DB) PedidosVacio(ctx context.Context) (bool, error) {
	var n int
	err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pedidos`).Scan(&n)
	return n == 0, err
}
