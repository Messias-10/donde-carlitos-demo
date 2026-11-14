package productos

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/texto"
)

// Límites de valores.
const (
	PrecioMaximo      = 10_000_000
	StockMaximo       = 1_000_000
	MayorMinimoTope   = 10_000
	UmbralPocosDefect = 5
)

var (
	reSlug     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	reCodigo   = regexp.MustCompile(`^[A-Z0-9][A-Z0-9-]{1,19}$`)
	reRutaImg  = regexp.MustCompile(`^/[A-Za-z0-9._~/-]{1,300}$`)
	hostTikTok = map[string]bool{"www.tiktok.com": true, "tiktok.com": true, "vm.tiktok.com": true, "m.tiktok.com": true}
)

// Entrada es el cuerpo de POST/PUT /api/admin/productos. PUT reemplaza todo:
// un campo opcional omitido queda en null.
type Entrada struct {
	Slug           string `json:"slug"`
	Codigo         string `json:"codigo"`
	Nombre         string `json:"nombre"`
	Descripcion    string `json:"descripcion"`
	CategoriaID    int64  `json:"categoria_id"`
	Precio         int64  `json:"precio"`
	PrecioOferta   *int64 `json:"precio_oferta"`
	EtiquetaOferta string `json:"etiqueta_oferta"`
	PrecioMayor    *int64 `json:"precio_mayor"`
	MayorMinimo    *int   `json:"mayor_minimo"`
	UnidadMayor    string `json:"unidad_mayor"`
	// MayorMinimoConfirmado: true solo si el mínimo salió de la fuente (video).
	MayorMinimoConfirmado bool    `json:"mayor_minimo_confirmado"`
	OfertaHasta           *string `json:"oferta_hasta"`
	Stock                 int     `json:"stock"`
	UmbralPocos           *int    `json:"umbral_pocos"`
	Activo                *bool   `json:"activo"`
	Temporada             *string `json:"temporada"`
	FechaLlegada          *string `json:"fecha_llegada"`
	ImagenURL             *string `json:"imagen_url"`
	TikTokURL             *string `json:"tiktok_url"`
	FuenteURL             *string `json:"fuente_url"`
	FuenteFecha           *string `json:"fuente_fecha"`
	EsEjemplo             bool    `json:"es_ejemplo"`
}

func val(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}

// Validar devuelve el producto listo para guardar o un mapa campo -> motivo.
// No verifica que la categoría exista (eso lo hace el handler contra la base).
func (e Entrada) Validar() (store.Producto, map[string]string) {
	errs := map[string]string{}
	var p store.Producto

	nombre, ok := texto.Limpiar(e.Nombre, false)
	switch {
	case !ok:
		errs["nombre"] = "tiene caracteres no permitidos"
	case texto.Largo(nombre) < 2 || texto.Largo(nombre) > 120:
		errs["nombre"] = "debe tener entre 2 y 120 caracteres"
	}
	p.Nombre = nombre

	desc, ok := texto.Limpiar(e.Descripcion, true)
	switch {
	case !ok:
		errs["descripcion"] = "tiene caracteres no permitidos"
	case texto.Largo(desc) > 2000:
		errs["descripcion"] = "máximo 2000 caracteres"
	}
	p.Descripcion = desc

	p.Slug = strings.TrimSpace(e.Slug)
	if p.Slug == "" {
		p.Slug = Slugify(nombre)
	}
	if len(p.Slug) > 80 || !reSlug.MatchString(p.Slug) {
		errs["slug"] = "solo minúsculas, números y guiones (máx. 80)"
	}

	p.Codigo = strings.ToUpper(strings.TrimSpace(e.Codigo))
	if !reCodigo.MatchString(p.Codigo) {
		errs["codigo"] = "2 a 20 caracteres: letras mayúsculas, números y guiones"
	}

	if e.CategoriaID <= 0 {
		errs["categoria_id"] = "es obligatoria"
	}
	p.CategoriaID = e.CategoriaID

	if e.Precio <= 0 || e.Precio > PrecioMaximo {
		errs["precio"] = "debe ser un entero en pesos entre 1 y 10.000.000"
	}
	p.PrecioDetalle = e.Precio

	if e.PrecioOferta != nil {
		if *e.PrecioOferta <= 0 || *e.PrecioOferta >= e.Precio {
			errs["precio_oferta"] = "debe ser mayor que 0 y menor que el precio normal"
		}
		v := *e.PrecioOferta
		p.PrecioOferta = &v
	}
	if f := val(e.OfertaHasta); f != "" {
		if !esFecha(f) {
			errs["oferta_hasta"] = "formato AAAA-MM-DD"
		}
		if e.PrecioOferta == nil {
			errs["oferta_hasta"] = "solo se usa junto con precio_oferta"
		}
		p.OfertaHasta = f
	}
	p.EtiquetaOferta = strings.TrimSpace(e.EtiquetaOferta)
	if p.EtiquetaOferta == "" {
		p.EtiquetaOferta = BadgeOferta
	}
	if p.EtiquetaOferta != BadgeOferta && p.EtiquetaOferta != BadgeRemate {
		errs["etiqueta_oferta"] = "debe ser oferta o remate"
	}

	switch {
	case e.PrecioMayor == nil && e.MayorMinimo == nil:
		if strings.TrimSpace(e.UnidadMayor) != "" {
			errs["unidad_mayor"] = "solo se usa junto con precio_mayor"
		}
	case e.PrecioMayor == nil || e.MayorMinimo == nil:
		errs["precio_mayor"] = "precio_mayor y mayor_minimo van juntos"
	default:
		if *e.PrecioMayor <= 0 || *e.PrecioMayor >= e.Precio {
			errs["precio_mayor"] = "debe ser mayor que 0 y menor que el precio al detalle (es por unidad)"
		}
		if *e.MayorMinimo < 2 || *e.MayorMinimo > MayorMinimoTope {
			errs["mayor_minimo"] = "debe estar entre 2 y 10.000 unidades"
		}
		pm, mm := *e.PrecioMayor, *e.MayorMinimo
		p.PrecioMayor, p.MayorMinimo = &pm, &mm
		unidad, ok := texto.Limpiar(e.UnidadMayor, false)
		if !ok || texto.Largo(unidad) > 60 {
			errs["unidad_mayor"] = "máximo 60 caracteres, sin caracteres de control"
		}
		p.UnidadMayor = unidad
		p.MayorMinimoConfirmado = e.MayorMinimoConfirmado
	}

	if e.Stock < 0 || e.Stock > StockMaximo {
		errs["stock"] = "debe estar entre 0 y 1.000.000"
	}
	p.Stock = e.Stock

	p.UmbralPocos = UmbralPocosDefect
	if e.UmbralPocos != nil {
		if *e.UmbralPocos < 0 || *e.UmbralPocos > 1000 {
			errs["umbral_pocos"] = "debe estar entre 0 y 1000"
		}
		p.UmbralPocos = *e.UmbralPocos
	}

	p.Activo = true
	if e.Activo != nil {
		p.Activo = *e.Activo
	}

	if t := val(e.Temporada); t != "" {
		if !EsTemporada(t) {
			errs["temporada"] = "temporada desconocida"
		}
		p.Temporada = t
	}

	if f := val(e.FechaLlegada); f != "" {
		if !esFecha(f) {
			errs["fecha_llegada"] = "formato AAAA-MM-DD"
		}
		p.FechaLlegada = f
	}
	if f := val(e.FuenteFecha); f != "" {
		if !esFecha(f) {
			errs["fuente_fecha"] = "formato AAAA-MM-DD"
		}
		p.FuenteFecha = f
	}

	if u := val(e.ImagenURL); u != "" {
		if !EsImagenValida(u) {
			errs["imagen_url"] = "debe ser una ruta /... del sitio o una URL https"
		}
		p.ImagenURL = u
	}
	if u := val(e.TikTokURL); u != "" {
		if !EsURLTikTok(u) {
			errs["tiktok_url"] = "debe ser un link https de tiktok.com"
		}
		p.TikTokURL = u
	}
	if u := val(e.FuenteURL); u != "" {
		if !esHTTPS(u) {
			errs["fuente_url"] = "debe ser una URL https"
		}
		p.FuenteURL = u
	}
	p.EsEjemplo = e.EsEjemplo

	if len(errs) > 0 {
		return store.Producto{}, errs
	}
	p.Busqueda = TextoBusqueda(p)
	return p, nil
}

func esFecha(s string) bool {
	_, err := time.Parse(time.DateOnly, s)
	return err == nil
}

// EsImagenValida acepta rutas locales del front (/productos/x.webp) o https.
func EsImagenValida(u string) bool {
	if strings.HasPrefix(u, "/") {
		return !strings.HasPrefix(u, "//") && !strings.Contains(u, "..") && reRutaImg.MatchString(u)
	}
	return esHTTPS(u)
}

func esHTTPS(u string) bool {
	if len(u) > 500 {
		return false
	}
	pu, err := url.Parse(u)
	return err == nil && pu.Scheme == "https" && pu.Host != "" && pu.User == nil
}

// EsURLTikTok acepta solo links https a tiktok.com.
func EsURLTikTok(u string) bool {
	if !esHTTPS(u) {
		return false
	}
	pu, _ := url.Parse(u)
	return hostTikTok[strings.ToLower(pu.Hostname())] && pu.Port() == ""
}
