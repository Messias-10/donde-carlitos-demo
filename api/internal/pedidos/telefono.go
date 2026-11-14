package pedidos

import "strings"

// NormalizarTelefono acepta formatos chilenos comunes ("+56 9 1234 5678",
// "56912345678", "912345678", "9 1234-5678") y devuelve E.164: +56912345678.
// Acepta celulares (9XXXXXXXX) y fijos de 9 dígitos (2XXXXXXXX, 3X..., etc.).
func NormalizarTelefono(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if len(s) > 25 {
		return "", false
	}
	var d strings.Builder
	for i, r := range s {
		switch {
		case r >= '0' && r <= '9':
			d.WriteRune(r)
		case r == '+' && i == 0:
		case r == ' ' || r == '-' || r == '(' || r == ')' || r == '.':
		default:
			return "", false
		}
	}
	dig := d.String()
	switch {
	case len(dig) == 11 && strings.HasPrefix(dig, "56"):
		dig = dig[2:]
	case len(dig) == 9:
	default:
		return "", false
	}
	if dig[0] < '2' || dig[0] > '9' {
		return "", false
	}
	return "+56" + dig, true
}

// EnmascararTelefono deja visibles solo el prefijo y los 4 últimos dígitos,
// para los logs: +56912345678 -> +569****5678.
func EnmascararTelefono(tel string) string {
	if len(tel) < 9 {
		return "****"
	}
	return tel[:4] + strings.Repeat("*", len(tel)-8) + tel[len(tel)-4:]
}
