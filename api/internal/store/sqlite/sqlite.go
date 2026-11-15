// Package sqlite implementa store.Store sobre SQLite con modernc.org/sqlite
// (Go puro, sin CGO). Todo el SQL usa parámetros; las únicas partes que se
// arman con texto son listas de "?" y cláusulas ORDER BY de una lista blanca.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite" // driver "sqlite"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/migrations"
)

// formatoTiempo tiene ancho fijo para que el orden de texto sea el orden temporal.
const formatoTiempo = "2006-01-02T15:04:05.000Z"

func ts(t time.Time) string { return t.UTC().Format(formatoTiempo) }

func parseTS(s string) time.Time {
	t, err := time.Parse(formatoTiempo, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// DB es la implementación SQLite de store.Store.
type DB struct {
	db *sql.DB
}

var _ store.Store = (*DB)(nil)

// Abrir abre (o crea) la base en ruta y aplica las migraciones pendientes.
func Abrir(ctx context.Context, ruta string) (*DB, error) {
	q := url.Values{}
	q.Add("_pragma", "foreign_keys(1)")
	q.Add("_pragma", "journal_mode(WAL)")
	q.Add("_pragma", "busy_timeout(10000)")
	q.Add("_pragma", "synchronous(NORMAL)")
	q.Add("_pragma", "temp_store(memory)")
	// BEGIN IMMEDIATE: las transacciones de escritura toman el lock al inicio,
	// así dos pedidos simultáneos se serializan en vez de fallar con SQLITE_BUSY.
	q.Set("_txlock", "immediate")
	dsn := "file:" + ruta + "?" + q.Encode()

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	db.SetConnMaxIdleTime(5 * time.Minute)

	d := &DB{db: db}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	if err := d.migrar(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return d, nil
}

// Ping verifica la conexión.
func (d *DB) Ping(ctx context.Context) error { return d.db.PingContext(ctx) }

// Close cierra la base.
func (d *DB) Close() error { return d.db.Close() }

func (d *DB) migrar(ctx context.Context) error {
	if _, err := d.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migraciones (
		version TEXT PRIMARY KEY,
		aplicada_en TEXT NOT NULL
	)`); err != nil {
		return fmt.Errorf("crear schema_migraciones: %w", err)
	}
	entradas, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("leer migraciones: %w", err)
	}
	var nombres []string
	for _, e := range entradas {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			nombres = append(nombres, e.Name())
		}
	}
	sort.Strings(nombres)

	for _, nombre := range nombres {
		var n int
		if err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migraciones WHERE version = ?`, nombre).Scan(&n); err != nil {
			return fmt.Errorf("consultar migración %s: %w", nombre, err)
		}
		if n > 0 {
			continue
		}
		contenido, err := fs.ReadFile(migrations.FS, nombre)
		if err != nil {
			return fmt.Errorf("leer %s: %w", nombre, err)
		}
		tx, err := d.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(contenido)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("aplicar %s: %w", nombre, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migraciones (version, aplicada_en) VALUES (?, ?)`, nombre, ts(time.Now())); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("registrar %s: %w", nombre, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit %s: %w", nombre, err)
		}
	}
	return nil
}

// mapearError traduce errores de restricciones de SQLite a errores del store.
func mapearError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "UNIQUE constraint failed"):
		return fmt.Errorf("%w: %v", store.ErrDuplicado, err)
	case strings.Contains(msg, "FOREIGN KEY constraint failed"):
		return fmt.Errorf("%w: %v", store.ErrReferenciaInvalida, err)
	}
	return err
}

// escaparLike neutraliza los comodines de LIKE. Se usa con ESCAPE '\'.
func escaparLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

// marcadores devuelve "?, ?, ?" (n veces). Solo genera signos, nunca datos.
func marcadores(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimSuffix(strings.Repeat("?, ", n), ", ")
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullInt64(p *int64) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullInt(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullTS(t *time.Time) any {
	if t == nil {
		return nil
	}
	return ts(*t)
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ---------------------------------------------------------------- catálogo

// ListarCategorias devuelve las categorías con su conteo de productos visibles.
func (d *DB) ListarCategorias(ctx context.Context) ([]store.Categoria, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT c.id, c.slug, c.nombre, c.orden,
		       (SELECT COUNT(*) FROM productos p
		         WHERE p.categoria_id = c.id AND p.activo = 1 AND p.eliminado_en IS NULL)
		  FROM categorias c
		 ORDER BY c.orden, c.nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []store.Categoria
	for rows.Next() {
		var c store.Categoria
		if err := rows.Scan(&c.ID, &c.Slug, &c.Nombre, &c.Orden, &c.TotalProductos); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// CrearCategoria inserta una categoría.
func (d *DB) CrearCategoria(ctx context.Context, c *store.Categoria, ahora time.Time) error {
	res, err := d.db.ExecContext(ctx, `INSERT INTO categorias (slug, nombre, orden, creado_en) VALUES (?, ?, ?, ?)`,
		c.Slug, c.Nombre, c.Orden, ts(ahora))
	if err != nil {
		return mapearError(err)
	}
	c.ID, err = res.LastInsertId()
	return err
}

// ExisteCategoria indica si la categoría existe.
func (d *DB) ExisteCategoria(ctx context.Context, id int64) (bool, error) {
	var n int
	err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categorias WHERE id = ?`, id).Scan(&n)
	return n > 0, err
}

const columnasProducto = `p.id, p.slug, p.codigo, p.nombre, p.descripcion, p.busqueda,
	p.categoria_id, c.slug, c.nombre,
	p.precio_detalle, p.precio_oferta, p.etiqueta_oferta, p.precio_mayor, p.mayor_minimo, p.unidad_mayor,
	p.stock, p.umbral_pocos, p.activo, p.temporada, p.fecha_llegada,
	p.imagen_url, p.tiktok_url, p.fuente_url, p.fuente_fecha, p.es_ejemplo, p.creado_en, p.actualizado_en,
	p.mayor_minimo_confirmado, p.oferta_hasta, p.precio_actualizado_en,
	p.precio_es_promo, p.nota_promo, p.en_vitrina, p.orden_vitrina`

const desdeProducto = ` FROM productos p JOIN categorias c ON c.id = p.categoria_id `

type escaneable interface {
	Scan(dest ...any) error
}

func escanearProducto(s escaneable) (store.Producto, error) {
	var (
		p                                  store.Producto
		oferta, mayor, minimo              sql.NullInt64
		temporada, llegada, imagen, tiktok sql.NullString
		unidad, fuente, fuenteFecha        sql.NullString
		ofertaHasta, precioAct, notaPromo  sql.NullString
		ordenVitrina                       sql.NullInt64
		activo, ejemplo, mayorConf         int
		esPromo, enVitrina                 int
		creado, actualizado                string
	)
	err := s.Scan(&p.ID, &p.Slug, &p.Codigo, &p.Nombre, &p.Descripcion, &p.Busqueda,
		&p.CategoriaID, &p.CategoriaSlug, &p.CategoriaNombre,
		&p.PrecioDetalle, &oferta, &p.EtiquetaOferta, &mayor, &minimo, &unidad,
		&p.Stock, &p.UmbralPocos, &activo, &temporada, &llegada,
		&imagen, &tiktok, &fuente, &fuenteFecha, &ejemplo, &creado, &actualizado,
		&mayorConf, &ofertaHasta, &precioAct,
		&esPromo, &notaPromo, &enVitrina, &ordenVitrina)
	if err != nil {
		return p, err
	}
	if oferta.Valid {
		v := oferta.Int64
		p.PrecioOferta = &v
	}
	if mayor.Valid {
		v := mayor.Int64
		p.PrecioMayor = &v
	}
	if minimo.Valid {
		v := int(minimo.Int64)
		p.MayorMinimo = &v
	}
	p.Activo = activo == 1
	p.EsEjemplo = ejemplo == 1
	p.Temporada = temporada.String
	p.FechaLlegada = llegada.String
	p.ImagenURL = imagen.String
	p.TikTokURL = tiktok.String
	p.UnidadMayor = unidad.String
	p.FuenteURL = fuente.String
	p.FuenteFecha = fuenteFecha.String
	p.MayorMinimoConfirmado = mayorConf == 1
	p.OfertaHasta = ofertaHasta.String
	p.PrecioEsPromo, p.NotaPromo, p.EnVitrina = esPromo == 1, notaPromo.String, enVitrina == 1
	if ordenVitrina.Valid {
		v := int(ordenVitrina.Int64)
		p.OrdenVitrina = &v
	}
	if precioAct.Valid {
		t := parseTS(precioAct.String)
		p.PrecioActualizadoEn = &t
	}
	p.CreadoEn = parseTS(creado)
	p.ActualizadoEn = parseTS(actualizado)
	return p, nil
}

// precioFinalSQL es el precio que se cobra hoy: la oferta solo si sigue vigente.
// Lleva un parámetro (la fecha de hoy).
const precioFinalSQL = `(CASE WHEN p.precio_oferta IS NOT NULL AND (p.oferta_hasta IS NULL OR p.oferta_hasta >= ?) THEN p.precio_oferta ELSE p.precio_detalle END)`

// hoyFiltro evita que un filtro sin fecha deje fuera ofertas vigentes.
func hoyFiltro(h string) string {
	if h == "" {
		return "0000-00-00"
	}
	return h
}

// ordenes es la lista blanca de ORDER BY. Los agotados van al final en "relevancia".
var ordenes = map[string]string{
	"relevancia":  `(p.stock = 0) ASC, p.creado_en DESC, p.id DESC`,
	"precio_asc":  precioFinalSQL + ` ASC, p.id ASC`,
	"precio_desc": precioFinalSQL + ` DESC, p.id DESC`,
	"recientes":   `p.fecha_llegada IS NULL, p.fecha_llegada DESC, p.id DESC`,
	"vitrina":     `p.orden_vitrina IS NULL, p.orden_vitrina, p.fecha_llegada DESC, p.id DESC`,
}

// ListarProductos aplica filtros y paginación. Devuelve la página y el total.
func (d *DB) ListarProductos(ctx context.Context, f store.FiltroProductos) ([]store.Producto, int, error) {
	conds := []string{"p.eliminado_en IS NULL"}
	var args []any
	if !f.IncluirInactivos {
		conds = append(conds, "p.activo = 1")
	}
	if f.CategoriaSlug != "" {
		conds = append(conds, "c.slug = ?")
		args = append(args, f.CategoriaSlug)
	}
	for _, palabra := range f.Palabras {
		conds = append(conds, `p.busqueda LIKE ? ESCAPE '\'`)
		args = append(args, "%"+escaparLike(palabra)+"%")
	}
	if f.SoloVitrina {
		conds = append(conds, "p.en_vitrina = 1")
	}
	if f.SoloOferta {
		conds = append(conds, "p.precio_oferta IS NOT NULL AND (p.oferta_hasta IS NULL OR p.oferta_hasta >= ?)")
		args = append(args, hoyFiltro(f.Hoy))
	}
	if f.Temporada != "" {
		conds = append(conds, "p.temporada = ?")
		args = append(args, f.Temporada)
	}
	if f.LlegadaDesde != "" {
		conds = append(conds, "p.fecha_llegada IS NOT NULL AND p.fecha_llegada >= ?")
		args = append(args, f.LlegadaDesde)
	}
	where := " WHERE " + strings.Join(conds, " AND ")

	var total int
	if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*)"+desdeProducto+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orden, ok := ordenes[f.Orden]
	if !ok {
		orden = ordenes["relevancia"]
	}
	if strings.Contains(orden, precioFinalSQL) {
		args = append(args, hoyFiltro(f.Hoy))
	}
	if len(f.Palabras) > 0 && (f.Orden == "" || f.Orden == "relevancia") {
		orden = `(CASE WHEN p.busqueda LIKE ? ESCAPE '\' THEN 0 ELSE 1 END), instr(p.busqueda, ?), ` + orden
		args = append(args, escaparLike(f.Palabras[0])+"%", f.Palabras[0])
	}
	limite := f.Limite
	if limite <= 0 {
		limite = 24
	}
	q := "SELECT " + columnasProducto + desdeProducto + where + " ORDER BY " + orden + " LIMIT ? OFFSET ?"
	rows, err := d.db.QueryContext(ctx, q, append(args, limite, f.Desplazamiento)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []store.Producto
	for rows.Next() {
		p, err := escanearProducto(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

// ProductoPorSlug busca por slug entre los no eliminados.
func (d *DB) ProductoPorSlug(ctx context.Context, slug string, incluirInactivos bool) (store.Producto, error) {
	q := "SELECT " + columnasProducto + desdeProducto + " WHERE p.slug = ? AND p.eliminado_en IS NULL"
	if !incluirInactivos {
		q += " AND p.activo = 1"
	}
	p, err := escanearProducto(d.db.QueryRowContext(ctx, q, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return p, store.ErrNoEncontrado
	}
	return p, err
}

// ProductoPorID busca por ID entre los no eliminados (activos o no).
func (d *DB) ProductoPorID(ctx context.Context, id int64) (store.Producto, error) {
	q := "SELECT " + columnasProducto + desdeProducto + " WHERE p.id = ? AND p.eliminado_en IS NULL"
	p, err := escanearProducto(d.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return p, store.ErrNoEncontrado
	}
	return p, err
}

// CrearProducto inserta un producto y completa su ID.
func (d *DB) CrearProducto(ctx context.Context, p *store.Producto, ahora time.Time) error {
	res, err := d.db.ExecContext(ctx, `
		INSERT INTO productos (slug, codigo, nombre, descripcion, busqueda, categoria_id,
			precio_detalle, precio_oferta, etiqueta_oferta, precio_mayor, mayor_minimo, unidad_mayor,
			stock, umbral_pocos, activo, temporada, fecha_llegada, imagen_url, tiktok_url,
			fuente_url, fuente_fecha, es_ejemplo, creado_en, actualizado_en,
			mayor_minimo_confirmado, oferta_hasta, precio_actualizado_en,
			precio_es_promo, nota_promo, en_vitrina, orden_vitrina)
		VALUES (`+marcadores(31)+`)`,
		p.Slug, p.Codigo, p.Nombre, p.Descripcion, p.Busqueda, p.CategoriaID,
		p.PrecioDetalle, nullInt64(p.PrecioOferta), p.EtiquetaOferta, nullInt64(p.PrecioMayor), nullInt(p.MayorMinimo), nullStr(p.UnidadMayor),
		p.Stock, p.UmbralPocos, boolInt(p.Activo), nullStr(p.Temporada), nullStr(p.FechaLlegada),
		nullStr(p.ImagenURL), nullStr(p.TikTokURL), nullStr(p.FuenteURL), nullStr(p.FuenteFecha),
		boolInt(p.EsEjemplo), ts(ahora), ts(ahora),
		boolInt(p.MayorMinimoConfirmado), nullStr(p.OfertaHasta), nullTS(p.PrecioActualizadoEn),
		boolInt(p.PrecioEsPromo), nullStr(p.NotaPromo), boolInt(p.EnVitrina), nullInt(p.OrdenVitrina))
	if err != nil {
		return mapearError(err)
	}
	p.ID, err = res.LastInsertId()
	return err
}

// ActualizarProducto reemplaza todos los campos editables.
func (d *DB) ActualizarProducto(ctx context.Context, p *store.Producto, ahora time.Time) error {
	res, err := d.db.ExecContext(ctx, `
		UPDATE productos SET slug = ?, codigo = ?, nombre = ?, descripcion = ?, busqueda = ?,
			categoria_id = ?, precio_detalle = ?, precio_oferta = ?, etiqueta_oferta = ?,
			precio_mayor = ?, mayor_minimo = ?, unidad_mayor = ?, stock = ?, umbral_pocos = ?, activo = ?,
			temporada = ?, fecha_llegada = ?, imagen_url = ?, tiktok_url = ?,
			fuente_url = ?, fuente_fecha = ?, es_ejemplo = ?, actualizado_en = ?,
			mayor_minimo_confirmado = ?, oferta_hasta = ?, precio_actualizado_en = ?,
			precio_es_promo = ?, nota_promo = ?, en_vitrina = ?, orden_vitrina = ?
		WHERE id = ? AND eliminado_en IS NULL`,
		p.Slug, p.Codigo, p.Nombre, p.Descripcion, p.Busqueda,
		p.CategoriaID, p.PrecioDetalle, nullInt64(p.PrecioOferta), p.EtiquetaOferta,
		nullInt64(p.PrecioMayor), nullInt(p.MayorMinimo), nullStr(p.UnidadMayor), p.Stock, p.UmbralPocos, boolInt(p.Activo),
		nullStr(p.Temporada), nullStr(p.FechaLlegada), nullStr(p.ImagenURL), nullStr(p.TikTokURL),
		nullStr(p.FuenteURL), nullStr(p.FuenteFecha), boolInt(p.EsEjemplo),
		ts(ahora), boolInt(p.MayorMinimoConfirmado), nullStr(p.OfertaHasta), nullTS(p.PrecioActualizadoEn),
		boolInt(p.PrecioEsPromo), nullStr(p.NotaPromo), boolInt(p.EnVitrina), nullInt(p.OrdenVitrina), p.ID)
	if err != nil {
		return mapearError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return store.ErrNoEncontrado
	}
	return nil
}

// ActualizarImagen cambia solo imagen_url y devuelve la anterior.
func (d *DB) ActualizarImagen(ctx context.Context, id int64, url string, ahora time.Time) (string, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()
	var anterior sql.NullString
	err = tx.QueryRowContext(ctx, `SELECT imagen_url FROM productos WHERE id = ? AND eliminado_en IS NULL`, id).Scan(&anterior)
	if errors.Is(err, sql.ErrNoRows) {
		return "", store.ErrNoEncontrado
	}
	if err != nil {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE productos SET imagen_url = ?, actualizado_en = ? WHERE id = ?`, url, ts(ahora), id); err != nil {
		return "", err
	}
	return anterior.String, tx.Commit()
}

// ContarVitrina cuenta productos en vitrina, sin contar excluirID.
func (d *DB) ContarVitrina(ctx context.Context, excluirID int64) (int, error) {
	var n int
	err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM productos WHERE en_vitrina = 1 AND eliminado_en IS NULL AND id <> ?`, excluirID).Scan(&n)
	return n, err
}

// EliminarProducto hace borrado lógico: los pedidos antiguos siguen apuntando a él.
func (d *DB) EliminarProducto(ctx context.Context, id int64, ahora time.Time) error {
	res, err := d.db.ExecContext(ctx, `UPDATE productos SET eliminado_en = ?, activo = 0, actualizado_en = ? WHERE id = ? AND eliminado_en IS NULL`,
		ts(ahora), ts(ahora), id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return store.ErrNoEncontrado
	}
	return nil
}

// CatalogoVacio indica si aún no hay categorías (para decidir si sembrar).
func (d *DB) CatalogoVacio(ctx context.Context) (bool, error) {
	var n int
	err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categorias`).Scan(&n)
	return n == 0, err
}
