// Package seed carga el catálogo de demo y pedidos ficticios cuando la base
// está vacía. Los datos (semilla.json) se generaron desde
// research/productos-semilla.json: cada producto trae su video de origen.
package seed

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Messias-10/donde-carlitos-demo/api/internal/pedidos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/productos"
	"github.com/Messias-10/donde-carlitos-demo/api/internal/store"
)

// Embebido es la semilla incluida en el binario.
//
//go:embed semilla.json
var Embebido []byte

// Semilla es el formato de semilla.json.
type Semilla struct {
	Categorias []struct {
		Slug   string `json:"slug"`
		Nombre string `json:"nombre"`
		Orden  int    `json:"orden"`
	} `json:"categorias"`
	Productos   []ProductoSemilla `json:"productos"`
	PedidosDemo []PedidoSemilla   `json:"pedidos_demo"`
}

// ProductoSemilla es un productos.Entrada con la categoría por slug.
type ProductoSemilla struct {
	productos.Entrada
	Categoria string `json:"categoria"`
}

// PedidoSemilla describe un pedido DEMO y los estados por los que pasa.
type PedidoSemilla struct {
	Nombre     string   `json:"nombre"`
	Telefono   string   `json:"telefono"`
	Entrega    string   `json:"entrega"`
	Region     string   `json:"region"`
	Comuna     string   `json:"comuna"`
	Agencia    string   `json:"agencia"`
	Comentario string   `json:"comentario"`
	HaceHoras  int      `json:"hace_horas"`
	Estados    []string `json:"estados"`
	NotaFinal  string   `json:"nota_final"`
	Items      []struct {
		Slug     string `json:"slug"`
		Cantidad int    `json:"cantidad"`
	} `json:"items"`
}

// Cargar siembra si el catálogo está vacío. Devuelve true si sembró.
func Cargar(ctx context.Context, st store.Store, datos []byte, ahora time.Time, conPedidos bool, log *slog.Logger) (bool, error) {
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

	ids := map[string]int64{}
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
		ids[p.Slug] = p.ID
	}

	n := 0
	if conPedidos {
		for i, pd := range s.PedidosDemo {
			sol := pedidos.Solicitud{Nombre: pd.Nombre, Telefono: pd.Telefono, Entrega: pd.Entrega,
				Region: pd.Region, Comuna: pd.Comuna, Agencia: pd.Agencia, Comentario: pd.Comentario}
			for _, it := range pd.Items {
				id, ok := ids[it.Slug]
				if !ok {
					return false, fmt.Errorf("semilla: pedido %d usa producto %q inexistente", i, it.Slug)
				}
				sol.Items = append(sol.Items, pedidos.ItemSolicitud{ProductoID: id, Cantidad: it.Cantidad})
			}
			nuevo, errs := sol.Validar()
			if len(errs) > 0 {
				return false, fmt.Errorf("semilla: pedido %d inválido: %v", i, errs)
			}
			nuevo.EsDemo = true
			t := ahora.Add(-time.Duration(pd.HaceHoras) * time.Hour)
			nuevo.Ahora = t
			p, err := st.CrearPedido(ctx, nuevo, pedidos.Cotizar)
			if err != nil {
				return false, fmt.Errorf("semilla: pedido %d: %w", i, err)
			}
			for j, estado := range pd.Estados {
				t = t.Add(40 * time.Minute)
				nota := ""
				if j == len(pd.Estados)-1 {
					nota = pd.NotaFinal
				}
				if _, err := st.CambiarEstado(ctx, store.CambioEstado{PedidoID: p.ID, Nuevo: estado, Nota: nota, Ahora: t}, pedidos.Validador(estado)); err != nil {
					return false, fmt.Errorf("semilla: pedido %d -> %s: %w", i, estado, err)
				}
			}
			n++
		}
	}
	log.Info("semilla cargada", "categorias", len(cats), "productos", len(ids), "pedidos_demo", n)
	return true, nil
}
