// Package migrations embebe el SQL versionado del esquema. Los archivos se
// aplican en orden lexicográfico (0001_, 0002_, ...) y nunca se editan una vez
// publicados: un cambio de esquema es un archivo nuevo.
package migrations

import "embed"

// FS contiene los archivos *.sql de este directorio.
//
//go:embed *.sql
var FS embed.FS
