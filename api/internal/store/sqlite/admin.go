package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// SincronizarAdmin crea el admin o actualiza su hash con el que viene del
// entorno. Si el hash cambió (rotación de clave), cierra todas sus sesiones.
func (d *DB) SincronizarAdmin(ctx context.Context, usuario, passwordHash string, ahora time.Time) (int64, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var id int64
	var actual string
	err = tx.QueryRowContext(ctx, `SELECT id, password_hash FROM admins WHERE usuario = ?`, usuario).Scan(&id, &actual)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		res, err := tx.ExecContext(ctx, `INSERT INTO admins (usuario, password_hash, creado_en, actualizado_en) VALUES (?, ?, ?, ?)`,
			usuario, passwordHash, ts(ahora), ts(ahora))
		if err != nil {
			return 0, mapearError(err)
		}
		if id, err = res.LastInsertId(); err != nil {
			return 0, err
		}
	case err != nil:
		return 0, err
	case actual != passwordHash:
		if _, err := tx.ExecContext(ctx, `UPDATE admins SET password_hash = ?, actualizado_en = ? WHERE id = ?`, passwordHash, ts(ahora), id); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM sesiones WHERE admin_id = ?`, id); err != nil {
			return 0, err
		}
	}
	return id, tx.Commit()
}

// AdminPorUsuario busca un admin.
func (d *DB) AdminPorUsuario(ctx context.Context, usuario string) (store.Admin, error) {
	var a store.Admin
	err := d.db.QueryRowContext(ctx, `SELECT id, usuario, password_hash FROM admins WHERE usuario = ?`, usuario).
		Scan(&a.ID, &a.Usuario, &a.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return a, store.ErrNoEncontrado
	}
	return a, err
}

// CrearSesion guarda una sesión (solo el hash del token).
func (d *DB) CrearSesion(ctx context.Context, s store.Sesion) error {
	_, err := d.db.ExecContext(ctx, `INSERT INTO sesiones (token_hash, admin_id, creado_en, expira_en, ultimo_uso) VALUES (?, ?, ?, ?, ?)`,
		s.TokenHash, s.AdminID, ts(s.CreadoEn), ts(s.ExpiraEn), ts(s.UltimoUso))
	return mapearError(err)
}

// SesionPorHash busca una sesión y el usuario dueño.
func (d *DB) SesionPorHash(ctx context.Context, tokenHash string) (store.Sesion, error) {
	var s store.Sesion
	var creado, expira, uso string
	err := d.db.QueryRowContext(ctx, `
		SELECT s.token_hash, s.admin_id, a.usuario, s.creado_en, s.expira_en, s.ultimo_uso
		  FROM sesiones s JOIN admins a ON a.id = s.admin_id
		 WHERE s.token_hash = ?`, tokenHash).Scan(&s.TokenHash, &s.AdminID, &s.Usuario, &creado, &expira, &uso)
	if errors.Is(err, sql.ErrNoRows) {
		return s, store.ErrNoEncontrado
	}
	s.CreadoEn, s.ExpiraEn, s.UltimoUso = parseTS(creado), parseTS(expira), parseTS(uso)
	return s, err
}

// TocarSesion actualiza el último uso (para el vencimiento por inactividad).
func (d *DB) TocarSesion(ctx context.Context, tokenHash string, ahora time.Time) error {
	_, err := d.db.ExecContext(ctx, `UPDATE sesiones SET ultimo_uso = ? WHERE token_hash = ?`, ts(ahora), tokenHash)
	return err
}

// EliminarSesion borra una sesión (logout).
func (d *DB) EliminarSesion(ctx context.Context, tokenHash string) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM sesiones WHERE token_hash = ?`, tokenHash)
	return err
}

// PurgarSesiones borra las sesiones vencidas.
func (d *DB) PurgarSesiones(ctx context.Context, ahora time.Time) (int64, error) {
	res, err := d.db.ExecContext(ctx, `DELETE FROM sesiones WHERE expira_en < ?`, ts(ahora))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
