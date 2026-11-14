// Comando setup-admin: crea (o rota) la clave del panel admin.
//
//  1. Genera una clave aleatoria fuerte (crypto/rand).
//  2. Escribe en <dir>/.env el hash bcrypt (en base64, para que docker compose
//     no interprete los "$" del hash) y un secreto de sesión aleatorio.
//  3. Escribe la clave en texto plano en <dir>/ADMIN_PASSWORD.local.txt.
//
// No imprime ningún secreto. Conserva las demás variables del .env.
//
//	setup-admin -dir /work            primera vez
//	setup-admin -dir /work -rotar     cambia la clave (cierra las sesiones abiertas)
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
)

const (
	archivoEnv   = ".env"
	archivoClave = "ADMIN_PASSWORD.local.txt"
	largoClave   = 24
)

var reUsuario = regexp.MustCompile(`^[a-z][a-z0-9._-]{2,31}$`)

func main() {
	dir := flag.String("dir", ".", "carpeta donde escribir .env y ADMIN_PASSWORD.local.txt (api/)")
	usuario := flag.String("usuario", "admin", "usuario del panel")
	rotar := flag.Bool("rotar", false, "reemplaza una clave ya existente")
	flag.Parse()

	if err := correr(*dir, *usuario, *rotar); err != nil {
		fmt.Fprintln(os.Stderr, "setup-admin:", err)
		os.Exit(1)
	}
	fmt.Printf("Listo. Se actualizó %s (hash de la clave y secreto de sesión) y se guardó la clave en %s.\n",
		filepath.Join(*dir, archivoEnv), filepath.Join(*dir, archivoClave))
	fmt.Println("La clave no se muestra aquí: abre ese archivo tú. Reinicia la API para que tome los cambios.")
}

func correr(dir, usuario string, rotar bool) error {
	if !reUsuario.MatchString(usuario) {
		return errors.New("usuario inválido (3 a 32 caracteres: minúsculas, números, . _ -)")
	}
	rutaEnv := filepath.Join(dir, archivoEnv)
	rutaClave := filepath.Join(dir, archivoClave)

	lineas, err := leerLineas(rutaEnv)
	if err != nil {
		return err
	}
	if valor(lineas, "ADMIN_PASSWORD_HASH_B64") != "" && !rotar {
		return errors.New("ya hay una clave configurada en .env; usa -rotar para cambiarla")
	}

	clave, err := auth.GenerarClave(largoClave)
	if err != nil {
		return err
	}
	hash, err := auth.HashClave(clave, auth.CostoBcrypt)
	if err != nil {
		return err
	}
	secreto := make([]byte, 32)
	if _, err := rand.Read(secreto); err != nil {
		return err
	}

	lineas = fijar(lineas, "ADMIN_USUARIO", usuario)
	lineas = fijar(lineas, "ADMIN_PASSWORD_HASH_B64", base64.StdEncoding.EncodeToString([]byte(hash)))
	lineas = fijar(lineas, "SESSION_SECRET", hex.EncodeToString(secreto))

	// Primero la clave: si falla, el .env queda como estaba.
	if err := escribirAtomico(rutaClave, []byte(clave+"\n")); err != nil {
		return fmt.Errorf("escribir %s: %w", archivoClave, err)
	}
	if err := escribirAtomico(rutaEnv, []byte(strings.Join(lineas, "\n")+"\n")); err != nil {
		return fmt.Errorf("escribir %s: %w", archivoEnv, err)
	}
	return nil
}

func leerLineas(ruta string) ([]string, error) {
	b, err := os.ReadFile(ruta)
	if errors.Is(err, os.ErrNotExist) {
		return []string{"# Generado por setup-admin. No se sube a git. Ver api/.env.example."}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("leer .env: %w", err)
	}
	var out []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		out = append(out, strings.TrimRight(sc.Text(), "\r"))
	}
	return out, sc.Err()
}

func valor(lineas []string, clave string) string {
	for _, l := range lineas {
		if k, v, ok := strings.Cut(l, "="); ok && strings.TrimSpace(k) == clave {
			return strings.Trim(strings.TrimSpace(v), `"'`)
		}
	}
	return ""
}

// fijar reemplaza la variable si existe (sin duplicarla) o la agrega al final.
func fijar(lineas []string, clave, v string) []string {
	out := make([]string, 0, len(lineas)+1)
	puesto := false
	for _, l := range lineas {
		if k, _, ok := strings.Cut(l, "="); ok && strings.TrimSpace(k) == clave {
			if !puesto {
				out = append(out, clave+"="+v)
				puesto = true
			}
			continue
		}
		out = append(out, l)
	}
	if !puesto {
		out = append(out, clave+"="+v)
	}
	return out
}

// escribirAtomico escribe en un temporal con permisos 0600 y lo renombra.
func escribirAtomico(ruta string, datos []byte) error {
	f, err := os.CreateTemp(filepath.Dir(ruta), ".setup-admin-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer func() { _ = os.Remove(tmp) }()
	if err := f.Chmod(0o600); err != nil {
		_ = f.Close()
		return err
	}
	if _, err := f.Write(datos); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, ruta)
}
