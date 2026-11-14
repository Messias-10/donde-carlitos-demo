package httpapi

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "golang.org/x/image/webp" // registra el decodificador webp para image.Decode

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

const (
	maxImagen        = 5 << 20 // 5 MB
	maxLadoImagen    = 5000
	maxPixelesImagen = 16_000_000
	prefijoImagenes  = "/api/imagenes/"
)

// Nombre generado por el servidor: 32 hex + extensión. Nada del nombre original.
var reArchivoImagen = regexp.MustCompile(`^[a-f0-9]{32}\.(jpg|png)$`)

// esRutaSubida indica si la ruta es la de subida de imagen (body más grande).
func esRutaSubida(r *http.Request) bool {
	return r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/admin/productos/") && strings.HasSuffix(r.URL.Path, "/imagen")
}

// tipoPorFirma reconoce el formato por los primeros bytes (no por el nombre ni
// por el Content-Type que manda el navegador).
func tipoPorFirma(b []byte) string {
	switch {
	case len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF:
		return "jpeg"
	case len(b) >= 8 && bytes.Equal(b[:8], []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}):
		return "png"
	case len(b) >= 12 && string(b[:4]) == "RIFF" && string(b[8:12]) == "WEBP":
		return "webp"
	}
	return ""
}

// leerImagenMultipart lee el campo "imagen" (único) de un multipart/form-data.
func leerImagenMultipart(r *http.Request) ([]byte, int, string, string) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "multipart/form-data" {
		return nil, http.StatusUnsupportedMediaType, "tipo_no_soportado", "Envía la imagen como multipart/form-data en el campo \"imagen\"."
	}
	mr, err := r.MultipartReader()
	if err != nil {
		return nil, http.StatusBadRequest, "imagen_invalida", "No se pudo leer la imagen."
	}
	var datos []byte
	for {
		parte, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			var mbe *http.MaxBytesError
			if errors.As(err, &mbe) {
				return nil, http.StatusRequestEntityTooLarge, "cuerpo_muy_grande", "La imagen pesa más de 5 MB."
			}
			return nil, http.StatusBadRequest, "imagen_invalida", "No se pudo leer la imagen."
		}
		if parte.FormName() != "imagen" || datos != nil {
			_ = parte.Close()
			return nil, http.StatusBadRequest, "imagen_invalida", "Envía un solo archivo en el campo \"imagen\"."
		}
		datos, err = io.ReadAll(io.LimitReader(parte, maxImagen+1))
		_ = parte.Close()
		if err != nil {
			var mbe *http.MaxBytesError
			if errors.As(err, &mbe) {
				return nil, http.StatusRequestEntityTooLarge, "cuerpo_muy_grande", "La imagen pesa más de 5 MB."
			}
			return nil, http.StatusBadRequest, "imagen_invalida", "No se pudo leer la imagen."
		}
		if len(datos) > maxImagen {
			return nil, http.StatusRequestEntityTooLarge, "cuerpo_muy_grande", "La imagen pesa más de 5 MB."
		}
	}
	if len(datos) == 0 {
		return nil, http.StatusBadRequest, "imagen_invalida", "Falta el archivo en el campo \"imagen\"."
	}
	return datos, 0, "", ""
}

// recodificar valida dimensiones, decodifica y vuelve a codificar: así se
// eliminan metadatos (EXIF, GPS) y cualquier contenido extra del archivo.
// PNG queda PNG (conserva transparencia); JPEG y WebP quedan JPEG.
func recodificar(datos []byte, tipo string) ([]byte, string, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(datos))
	if err != nil {
		return nil, "", err
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxLadoImagen || cfg.Height > maxLadoImagen || cfg.Width*cfg.Height > maxPixelesImagen {
		return nil, "", errors.New("dimensiones fuera de rango")
	}
	img, _, err := image.Decode(bytes.NewReader(datos))
	if err != nil {
		return nil, "", err
	}
	var out bytes.Buffer
	if tipo == "png" {
		err = (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(&out, img)
		return out.Bytes(), "png", err
	}
	err = jpeg.Encode(&out, img, &jpeg.Options{Quality: 85})
	return out.Bytes(), "jpg", err
}

func nombreImagen(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + "." + ext, nil
}

func escribirAtomico(dir, nombre string, datos []byte) error {
	f, err := os.CreateTemp(dir, ".subida-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer func() { _ = os.Remove(tmp) }()
	if _, err := f.Write(datos); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, nombre))
}

// adminSubirImagen: POST /api/admin/productos/{id}/imagen (multipart, campo "imagen").
func (s *Server) adminSubirImagen(w http.ResponseWriter, r *http.Request, ses store.Sesion) {
	if s.cfg.ImagenesDir == "" {
		responderError(w, http.StatusServiceUnavailable, "no_disponible", "La subida de imágenes no está configurada.", nil)
		return
	}
	id, ok := idDeRuta(r)
	if !ok {
		responderError(w, http.StatusNotFound, "no_encontrado", "Producto no encontrado.", nil)
		return
	}
	if !s.semImagen.Intentar() {
		w.Header().Set("Retry-After", "2")
		responderError(w, http.StatusTooManyRequests, "demasiadas_solicitudes", "Hay otra imagen subiéndose. Intenta de nuevo en unos segundos.", nil)
		return
	}
	defer s.semImagen.Soltar()

	// Más tiempo que el ReadTimeout general: una foto desde el celular puede tardar.
	rc := http.NewResponseController(w)
	_ = rc.SetReadDeadline(time.Now().Add(60 * time.Second))
	_ = rc.SetWriteDeadline(time.Now().Add(60 * time.Second))

	datos, status, codigo, msg := leerImagenMultipart(r)
	if datos == nil {
		responderError(w, status, codigo, msg, nil)
		return
	}
	tipo := tipoPorFirma(datos)
	if tipo == "" {
		responderError(w, http.StatusUnsupportedMediaType, "tipo_no_soportado", "Solo se aceptan imágenes JPG, PNG o WebP.", nil)
		return
	}
	limpio, ext, err := recodificar(datos, tipo)
	if err != nil {
		s.log.Warn("imagen rechazada", "err", err)
		responderError(w, http.StatusBadRequest, "imagen_invalida", "La imagen está dañada o es demasiado grande (máx. 5000 px por lado).", nil)
		return
	}
	nombre, err := nombreImagen(ext)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	if err := escribirAtomico(s.cfg.ImagenesDir, nombre, limpio); err != nil {
		s.errorInterno(w, r, err)
		return
	}
	anterior, err := s.st.ActualizarImagen(r.Context(), id, prefijoImagenes+nombre, s.ahora())
	if err != nil {
		_ = os.Remove(filepath.Join(s.cfg.ImagenesDir, nombre))
		s.responderErrorGuardar(w, r, err)
		return
	}
	// Si la imagen anterior también era subida, se borra el archivo viejo.
	if viejo, ok := strings.CutPrefix(anterior, prefijoImagenes); ok && reArchivoImagen.MatchString(viejo) {
		_ = os.Remove(filepath.Join(s.cfg.ImagenesDir, viejo))
	}
	p, err := s.st.ProductoPorID(r.Context(), id)
	if err != nil {
		s.errorInterno(w, r, err)
		return
	}
	s.log.Info("imagen de producto subida", "id", id, "bytes", len(limpio), "usuario", ses.Usuario)
	responderJSON(w, http.StatusOK, s.productoAdmin(p))
}

// servirImagen: GET /api/imagenes/{archivo}. Solo nombres generados por el
// servidor; Content-Type fijo por extensión, nosniff y CSP 'none' (globales).
func (s *Server) servirImagen(w http.ResponseWriter, r *http.Request) {
	nombre := r.PathValue("archivo")
	if s.cfg.ImagenesDir == "" || !reArchivoImagen.MatchString(nombre) {
		responderError(w, http.StatusNotFound, "no_encontrado", "Imagen no encontrada.", nil)
		return
	}
	f, err := os.Open(filepath.Join(s.cfg.ImagenesDir, nombre))
	if err != nil {
		responderError(w, http.StatusNotFound, "no_encontrado", "Imagen no encontrada.", nil)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		responderError(w, http.StatusNotFound, "no_encontrado", "Imagen no encontrada.", nil)
		return
	}
	ct := "image/jpeg"
	if strings.HasSuffix(nombre, ".png") {
		ct = "image/png"
	}
	h := w.Header()
	h.Set("Content-Type", ct)
	h.Set("Content-Disposition", "inline")
	// El nombre cambia con cada subida: se puede cachear para siempre.
	h.Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeContent(w, r, nombre, info.ModTime(), f)
}
