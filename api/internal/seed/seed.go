// Package seed carga el catálogo del demo (semilla.json, generado desde
// research/productos-semilla.json) y pedidos DEMO ficticios para el panel
// (pedidos_demo.json).
package seed

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/auth"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Embebido es el catálogo incluido en el binario.
//
//go:embed semilla.json
var Embebido []byte

// PedidosDemo son los pedidos ficticios del panel.
//
//go:embed pedidos_demo.json
var PedidosDemo []byte

// Semilla es el formato de semilla.json.
type Semilla struct {
	Categorias []struct {
		Slug   string `json:"slug"`
		Nombre string `json:"nombre"`
		Orden  int    `json:"orden"`
	} `json:"categorias"`
	Productos []ProductoSemilla `json:"productos"`
}

// ProductoSemilla es un productos.Entrada con la categoría por slug.
type ProductoSemilla struct {
	productos.Entrada
	Categoria string `json:"categoria"`
}

// Cargar siembra el catálogo si no hay categorías. Devuelve true si sembró.
func Cargar(ctx context.Context, st store.Store, datos []byte, ahora time.Time, log *slog.Logger) (bool, error) {
	vacio, err := st.CatalogoVacio(ctx)
	if err != nil || !vacio {
		return false, err
	}
	var s Semilla
	if err := json.Unmarshal(datos, &s); err != nil {
		return false, fmt.Errorf("semilla: json: %w", err)
	}
	cats := map[string]store.Categoria{}
	for _, c := range s.Categorias {
		cat := store.Categoria{Slug: c.Slug, Nombre: c.Nombre, Orden: c.Orden}
		if err := st.CrearCategoria(ctx, &cat, ahora); err != nil {
			return false, fmt.Errorf("semilla: categoría %s: %w", c.Slug, err)
		}
		cats[c.Slug] = cat
	}
	for _, ps := range s.Productos {
		cat, ok := cats[ps.Categoria]
		if !ok {
			return false, fmt.Errorf("semilla: producto %s: categoría %q no existe", ps.Slug, ps.Categoria)
		}
		ps.Entrada.CategoriaID = cat.ID
		p, errs := ps.Entrada.Validar()
		if len(errs) > 0 {
			return false, fmt.Errorf("semilla: producto %s inválido: %v", ps.Slug, errs)
		}
		p.CategoriaNombre = cat.Nombre
		p.Busqueda = productos.TextoBusqueda(p)
		if err := st.CrearProducto(ctx, &p, ahora); err != nil {
			return false, fmt.Errorf("semilla: producto %s: %w", ps.Slug, err)
		}
	}
	log.Info("catálogo sembrado", "categorias", len(cats), "productos", len(s.Productos))
	return true, nil
}

func token() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// PedidosSemilla es el formato de pedidos_demo.json.
type PedidosSemilla struct {
	CorrelativoInicial int64 `json:"correlativo_inicial"`
	Usuarios           []struct {
		Usuario string `json:"usuario"`
		Nombre  string `json:"nombre"`
	} `json:"usuarios"`
	Pedidos []PedidoSemilla `json:"pedidos"`
}

// PedidoSemilla es un pedido DEMO y los estados por los que pasa.
type PedidoSemilla struct {
	Canal         string `json:"canal"`
	Nombre        string `json:"nombre"`
	Telefono      string `json:"telefono"`
	Entrega       string `json:"entrega"`
	Region        string `json:"region"`
	Comuna        string `json:"comuna"`
	Agencia       string `json:"agencia"`
	FormaPago     string `json:"forma_pago"`
	Descuento     int64  `json:"descuento"`
	CostoEnvio    int64  `json:"costo_envio"`
	EnvioPorPagar bool   `json:"envio_por_pagar"`
	Total         *int64 `json:"total"`
	NotaInterna   string `json:"nota_interna"`
	RegistradoPor string `json:"registrado_por"`
	HaceHoras     int    `json:"hace_horas"`
	Items         []struct {
		Slug           string `json:"slug"`
		Descripcion    string `json:"descripcion"`
		Cantidad       int    `json:"cantidad"`
		PrecioUnitario *int64 `json:"precio_unitario"`
	} `json:"items"`
	Estados []struct {
		Estado       string `json:"estado"`
		HaceHoras    int    `json:"hace_horas"`
		Nota         string `json:"nota"`
		NSeguimiento string `json:"n_seguimiento"`
		Pago         string `json:"pago"` // medio; monto = total del pedido
		Motivo       string `json:"motivo"`
	} `json:"estados"`
}

// CargarPedidosDemo siembra usuarios demo (sin clave utilizable) y pedidos
// DEMO si la tabla de pedidos está vacía. Pasan por las mismas reglas que el
// panel: validación, precios, estados y stock.
func CargarPedidosDemo(ctx context.Context, st store.Store, datos []byte, ahora time.Time, zona *time.Location, log *slog.Logger) (int, error) {
	vacio, err := st.PedidosVacio(ctx)
	if err != nil || !vacio {
		return 0, err
	}
	var doc PedidosSemilla
	if err := json.Unmarshal(datos, &doc); err != nil {
		return 0, fmt.Errorf("pedidos demo: json: %w", err)
	}

	// Usuarios demo para que el historial muestre quién atendió. Su clave es
	// aleatoria y se descarta: nadie puede entrar con ellos hasta que un admin
	// les genere una (POST /api/admin/usuarios/{id}/clave).
	ids := map[string]int64{}
	for _, u := range doc.Usuarios {
		if ex, err := st.UsuarioPorUsuario(ctx, u.Usuario); err == nil {
			ids[u.Usuario] = ex.ID
			continue
		}
		descartable, err := auth.GenerarClave(32)
		if err != nil {
			return 0, err
		}
		hash, err := auth.HashClave(descartable, bcrypt.MinCost)
		if err != nil {
			return 0, err
		}
		nu := store.Usuario{Usuario: u.Usuario, Nombre: u.Nombre, Rol: store.RolVendedor, PasswordHash: hash}
		if err := st.CrearUsuario(ctx, &nu, ahora); err != nil {
			return 0, fmt.Errorf("pedidos demo: usuario %s: %w", u.Usuario, err)
		}
		ids[u.Usuario] = nu.ID
	}

	if doc.CorrelativoInicial > 0 {
		if err := st.FijarCorrelativo(ctx, doc.CorrelativoInicial); err != nil {
			return 0, err
		}
	}

	n := 0
	for i, pd := range doc.Pedidos {
		in := pedidos.Entrada{Canal: pd.Canal, Nombre: pd.Nombre, Telefono: pd.Telefono, Entrega: pd.Entrega,
			Region: pd.Region, Comuna: pd.Comuna, Agencia: pd.Agencia, FormaPago: pd.FormaPago,
			Descuento: pd.Descuento, CostoEnvio: pd.CostoEnvio, EnvioPorPagar: pd.EnvioPorPagar,
			Total: pd.Total, NotaInterna: pd.NotaInterna}
		for _, it := range pd.Items {
			if it.Slug == "" {
				in.Items = append(in.Items, pedidos.ItemEntrada{Descripcion: it.Descripcion, Cantidad: it.Cantidad, PrecioUnitario: it.PrecioUnitario})
				continue
			}
			p, err := st.ProductoPorSlug(ctx, it.Slug, true)
			if err != nil {
				return n, fmt.Errorf("pedidos demo %d: producto %q: %w", i, it.Slug, err)
			}
			in.Items = append(in.Items, pedidos.ItemEntrada{ProductoID: p.ID, Cantidad: it.Cantidad})
		}
		d, inicial, errs := in.Validar()
		if len(errs) > 0 {
			return n, fmt.Errorf("pedidos demo %d inválido: %v", i, errs)
		}
		tok, err := token()
		if err != nil {
			return n, err
		}
		usuario := ids[pd.RegistradoPor]
		creado := ahora.Add(-time.Duration(pd.HaceHoras) * time.Hour)
		p, err := st.CrearPedido(ctx, store.PedidoNuevo{DatosPedido: d, EstadoInicial: inicial, TokenSeguimiento: tok,
			RegistradoPor: usuario, EsDemo: true, Ahora: creado}, pedidos.Cotizador(creado.In(zona).Format(time.DateOnly)))
		if err != nil {
			return n, fmt.Errorf("pedidos demo %d: %w", i, err)
		}
		t := creado
		for _, e := range pd.Estados {
			if e.HaceHoras > 0 {
				t = ahora.Add(-time.Duration(e.HaceHoras) * time.Hour)
			} else {
				t = t.Add(30 * time.Minute)
			}
			c := store.CambioEstado{PedidoID: p.ID, Nuevo: e.Estado, UsuarioID: usuario, Nota: e.Nota,
				NSeguimiento: e.NSeguimiento, Motivo: e.Motivo, Ahora: t}
			if e.Pago != "" {
				c.Pago = &store.Pago{Medio: e.Pago, Monto: p.Total, Fecha: t.In(zona).Format(time.DateOnly)}
				c.ConfirmoPago = true
			}
			if p, err = st.CambiarEstado(ctx, c, pedidos.Validador()); err != nil {
				return n, fmt.Errorf("pedidos demo %d -> %s: %w", i, e.Estado, err)
			}
		}
		n++
	}
	log.Info("pedidos DEMO sembrados", "pedidos", n, "usuarios_demo", len(ids))
	return n, nil
}
