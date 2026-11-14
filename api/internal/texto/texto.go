// Package texto valida y limpia texto libre que llega de clientes y admins.
package texto

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Limpiar valida UTF-8, normaliza saltos de línea, quita espacios sobrantes y
// rechaza caracteres de control y de reordenamiento bidireccional (que sirven
// para disfrazar texto en el panel). Si multilinea es false, tampoco admite \n.
// Devuelve ok=false si el texto trae algo no permitido.
func Limpiar(s string, multilinea bool) (string, bool) {
	if !utf8.ValidString(s) {
		return "", false
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	for _, r := range s {
		if r == '\n' {
			if !multilinea {
				return "", false
			}
			continue
		}
		if r == '\t' {
			continue
		}
		if unicode.IsControl(r) || esBidi(r) || r == utf8.RuneError {
			return "", false
		}
	}
	s = strings.ReplaceAll(s, "\t", " ")
	if multilinea {
		lineas := strings.Split(s, "\n")
		for i, l := range lineas {
			lineas[i] = colapsar(l)
		}
		s = strings.Join(lineas, "\n")
		for strings.Contains(s, "\n\n\n") {
			s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
		}
		return strings.TrimSpace(s), true
	}
	return colapsar(s), true
}

func colapsar(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func esBidi(r rune) bool {
	return (r >= 0x202A && r <= 0x202E) || (r >= 0x2066 && r <= 0x2069) || r == 0x200E || r == 0x200F || r == 0x061C
}

// Largo cuenta runas.
func Largo(s string) int { return utf8.RuneCountInString(s) }

// EsNombrePersona: letras (con tildes y ñ), espacios, apóstrofo, punto y guion.
func EsNombrePersona(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.Is(unicode.Mn, r), r == ' ', r == '\'', r == '’', r == '.', r == '-':
		default:
			return false
		}
	}
	return true
}

// EsLugar: como un nombre pero admite dígitos y paréntesis (p. ej. "Puerto Montt", "Santiago (centro)").
func EsLugar(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.Is(unicode.Mn, r), unicode.IsDigit(r),
			r == ' ', r == '\'', r == '’', r == '.', r == '-', r == '(', r == ')', r == ',':
		default:
			return false
		}
	}
	return true
}
