package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
)

// Corre en un directorio temporal: nunca toca el api/.env real.
func TestSetupAdminEscribeArchivos(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("APP_ENV=development\nADMIN_PASSWORD_HASH_B64=\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := correr(dir, "admin", false); err != nil {
		t.Fatal(err)
	}
	env, _ := os.ReadFile(filepath.Join(dir, ".env"))
	clave, _ := os.ReadFile(filepath.Join(dir, archivoClave))
	lineas := strings.Split(strings.TrimSpace(string(env)), "\n")

	if valor(lineas, "APP_ENV") != "development" {
		t.Fatal("perdió variables existentes")
	}
	if strings.Count(string(env), "ADMIN_PASSWORD_HASH_B64=") != 1 {
		t.Fatal("variable duplicada")
	}
	if strings.Contains(string(env), "$") {
		t.Fatal("el .env no debe tener '$' (docker compose lo interpola)")
	}
	hash, err := base64.StdEncoding.DecodeString(valor(lineas, "ADMIN_PASSWORD_HASH_B64"))
	if err != nil || !auth.VerificarClave(string(hash), strings.TrimSpace(string(clave))) {
		t.Fatal("el hash no corresponde a la clave escrita")
	}
	if len(valor(lineas, "SESSION_SECRET")) != 64 {
		t.Fatal("SESSION_SECRET debe ser 32 bytes en hex")
	}
	if len(strings.TrimSpace(string(clave))) != largoClave {
		t.Fatal("largo de clave")
	}
	// Sin -rotar no pisa una clave existente.
	if err := correr(dir, "admin", false); err == nil {
		t.Fatal("debió negarse a sobrescribir")
	}
	if err := correr(dir, "admin", true); err != nil {
		t.Fatal(err)
	}
	clave2, _ := os.ReadFile(filepath.Join(dir, archivoClave))
	if string(clave2) == string(clave) {
		t.Fatal("rotar debe generar otra clave")
	}
}
