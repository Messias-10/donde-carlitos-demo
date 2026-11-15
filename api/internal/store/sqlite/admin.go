package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// SincronizarAdmin crea o actualiza el admin que viene de api/.env (setup-admin):
// hash del entorno, rol admin y activo. Es la vía de recuperación del panel.
// Si el hash cambió, cierra sus sesiones.
func (d *DB) SincronizarAdmin(ctx context.Context, usuario, passwordHash string, ahora time.Time) (int64, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var id int64
	var actual string
	err = tx.QueryRowContext(ctx, `SELECT id, password_hash FROM usuarios WHERE usuario = ?`, usuario).Scan(&id, &actual)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		res, err := tx.ExecContext(ctx, `INSERT INTO usuarios (usuario, nombre, password_hash, rol, activo, creado_en, actualizado_en)
			VALUES (?, 'Administrador', ?, 'admin', 1, ?, ?)`, usuario, passwordHash, ts(ahora), ts(ahora))
		if err != nil {
			return 0, mapearError(err)
		}
		if id, err = res.LastInsertId(); err != nil {
			return 0, err
		}
	case err != nil:
		return 0, err
	default:
		if _, err := tx.ExecContext(ctx, `UPDATE usuarios SET password_hash = ?, rol = 'admin', activo = 1, debe_cambiar_clave = 0, actualizado_en = ? WHERE id = ?`,
			passwordHash, ts(ahora), id); err != nil {
			return 0, err
		}
		if actual != passwordHash {
			if _, err := tx.ExecContext(ctx, `DELETE FROM sesiones WHERE usuario_id = ?`, id); err != nil {
				return 0, err
			}
		}
	}
	return id, tx.Commit()
}

const columnasUsuario = `id, usuario, nombre, password_hash, rol, activo, COALESCE(canal_por_defecto, ''), debe_cambiar_clave, creado_en`

func escanearUsuario(s escaneable) (store.Usuario, error) {
	var u store.Usuario
	var activo, debe int
	var creado string
	err := s.Scan(&u.ID, &u.Usuario, &u.Nombre, &u.PasswordHash, &u.Rol, &activo, &u.CanalPorDefecto, &debe, &creado)
	u.Activo, u.DebeCambiarClave, u.CreadoEn = activo == 1, debe == 1, parseTS(creado)
	if errors.Is(err, sql.ErrNoRows) {
		return u, store.ErrNoEncontrado
	}
	return u, err
}

// UsuarioPorUsuario busca por nombre de usuario (activo o no).
func (d *DB) UsuarioPorUsuario(ctx context.Context, usuario string) (store.Usuario, error) {
	return escanearUsuario(d.db.QueryRowContext(ctx, `SELECT `+columnasUsuario+` FROM usuarios WHERE usuario = ?`, usuario))
}

// UsuarioPorID busca por ID.
func (d *DB) UsuarioPorID(ctx context.Context, id int64) (store.Usuario, error) {
	return escanearUsuario(d.db.QueryRowContext(ctx, `SELECT `+columnasUsuario+` FROM usuarios WHERE id = ?`, id))
}

// ListarUsuarios devuelve todos, por nombre.
func (d *DB) ListarUsuarios(ctx context.Context) ([]store.Usuario, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT `+columnasUsuario+` FROM usuarios ORDER BY activo DESC, nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []store.Usuario
	for rows.Next() {
		u, err := escanearUsuario(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// CrearUsuario inserta un usuario (activo).
func (d *DB) CrearUsuario(ctx context.Context, u *store.Usuario, ahora time.Time) error {
	res, err := d.db.ExecContext(ctx, `INSERT INTO usuarios (usuario, nombre, password_hash, rol, activo, canal_por_defecto, debe_cambiar_clave, creado_en, actualizado_en)
		VALUES (?, ?, ?, ?, 1, ?, ?, ?, ?)`, u.Usuario, u.Nombre, u.PasswordHash, u.Rol, nullStr(u.CanalPorDefecto), boolInt(u.DebeCambiarClave), ts(ahora), ts(ahora))
	if err != nil {
		return mapearError(err)
	}
	u.ID, err = res.LastInsertId()
	u.Activo, u.CreadoEn = true, ahora
	return err
}

// ActualizarUsuario aplica los cambios; desactivar borra sus sesiones.
func (d *DB) ActualizarUsuario(ctx context.Context, id int64, c store.CambiosUsuario, ahora time.Time) (store.Usuario, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return store.Usuario{}, err
	}
	defer func() { _ = tx.Rollback() }()
	sets := []string{"actualizado_en = ?"}
	args := []any{ts(ahora)}
	if c.Nombre != nil {
		sets = append(sets, "nombre = ?")
		args = append(args, *c.Nombre)
	}
	if c.Rol != nil {
		sets = append(sets, "rol = ?")
		args = append(args, *c.Rol)
	}
	if c.Activo != nil {
		sets = append(sets, "activo = ?")
		args = append(args, boolInt(*c.Activo))
	}
	if c.CanalPorDefecto != nil {
		sets = append(sets, "canal_por_defecto = ?")
		args = append(args, nullStr(*c.CanalPorDefecto))
	}
	res, err := tx.ExecContext(ctx, "UPDATE usuarios SET "+strings.Join(sets, ", ")+" WHERE id = ?", append(args, id)...)
	if err != nil {
		return store.Usuario{}, mapearError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return store.Usuario{}, store.ErrNoEncontrado
	}
	if c.Activo != nil && !*c.Activo {
		if _, err := tx.ExecContext(ctx, `DELETE FROM sesiones WHERE usuario_id = ?`, id); err != nil {
			return store.Usuario{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return store.Usuario{}, err
	}
	return d.UsuarioPorID(ctx, id)
}

// CambiarClave reemplaza el hash, fija debe_cambiar_clave y cierra las
// sesiones del usuario salvo conservarSesion.
func (d *DB) CambiarClave(ctx context.Context, id int64, passwordHash string, debeCambiar bool, conservarSesion string, ahora time.Time) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.ExecContext(ctx, `UPDATE usuarios SET password_hash = ?, debe_cambiar_clave = ?, actualizado_en = ? WHERE id = ?`,
		passwordHash, boolInt(debeCambiar), ts(ahora), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return store.ErrNoEncontrado
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM sesiones WHERE usuario_id = ? AND token_hash <> ?`, id, conservarSesion); err != nil {
		return err
	}
	return tx.Commit()
}

// CrearSesion guarda una sesión (solo el hash del token).
func (d *DB) CrearSesion(ctx context.Context, s store.Sesion) error {
	_, err := d.db.ExecContext(ctx, `INSERT INTO sesiones (token_hash, usuario_id, creado_en, expira_en, ultimo_uso) VALUES (?, ?, ?, ?, ?)`,
		s.TokenHash, s.UsuarioID, ts(s.CreadoEn), ts(s.ExpiraEn), ts(s.UltimoUso))
	return mapearError(err)
}

// SesionPorHash busca una sesión de un usuario activo.
func (d *DB) SesionPorHash(ctx context.Context, tokenHash string) (store.Sesion, error) {
	var s store.Sesion
	var creado, expira, uso string
	err := d.db.QueryRowContext(ctx, `
		SELECT s.token_hash, s.usuario_id, u.usuario, u.nombre, u.rol, COALESCE(u.canal_por_defecto, ''), u.debe_cambiar_clave,
		       s.creado_en, s.expira_en, s.ultimo_uso
		  FROM sesiones s JOIN usuarios u ON u.id = s.usuario_id
		 WHERE s.token_hash = ? AND u.activo = 1`, tokenHash).
		Scan(&s.TokenHash, &s.UsuarioID, &s.Usuario, &s.Nombre, &s.Rol, &s.CanalPorDefecto, &s.DebeCambiarClave, &creado, &expira, &uso)
	if errors.Is(err, sql.ErrNoRows) {
		return s, store.ErrNoEncontrado
	}
	s.CreadoEn, s.ExpiraEn, s.UltimoUso = parseTS(creado), parseTS(expira), parseTS(uso)
	return s, err
}

// TocarSesion actualiza el último uso.
func (d *DB) TocarSesion(ctx context.Context, tokenHash string, ahora time.Time) error {
	_, err := d.db.ExecContext(ctx, `UPDATE sesiones SET ultimo_uso = ? WHERE token_hash = ?`, ts(ahora), tokenHash)
	return err
}

// EliminarSesion borra una sesión (logout).
func (d *DB) EliminarSesion(ctx context.Context, tokenHash string) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM sesiones WHERE token_hash = ?`, tokenHash)
	return err
}

// PurgarSesiones borra las vencidas.
func (d *DB) PurgarSesiones(ctx context.Context, ahora time.Time) (int64, error) {
	res, err := d.db.ExecContext(ctx, `DELETE FROM sesiones WHERE expira_en < ?`, ts(ahora))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
