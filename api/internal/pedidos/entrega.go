package pedidos

// Opcion es un par código / nombre para listas del checkout.
type Opcion struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

// Regiones de Chile (16), de norte a sur.
var Regiones = []Opcion{
	{"arica-parinacota", "Arica y Parinacota"},
	{"tarapaca", "Tarapacá"},
	{"antofagasta", "Antofagasta"},
	{"atacama", "Atacama"},
	{"coquimbo", "Coquimbo"},
	{"valparaiso", "Valparaíso"},
	{"metropolitana", "Metropolitana de Santiago"},
	{"ohiggins", "Libertador General Bernardo O'Higgins"},
	{"maule", "Maule"},
	{"nuble", "Ñuble"},
	{"biobio", "Biobío"},
	{"araucania", "La Araucanía"},
	{"los-rios", "Los Ríos"},
	{"los-lagos", "Los Lagos"},
	{"aysen", "Aysén del General Carlos Ibáñez del Campo"},
	{"magallanes", "Magallanes y de la Antártica Chilena"},
}

// Agencias de envío que el cliente puede elegir. [EJEMPLO]: la lista real la
// define el dueño (research/oportunidades.md, punto 5).
var Agencias = []Opcion{
	{"starken", "Starken"},
	{"chilexpress", "Chilexpress"},
	{"varmontt", "Varmontt"},
	{"pullman-cargo", "Pullman Cargo"},
	{"blue-express", "Blue Express"},
	{"otra", "Otra (indícala en el comentario)"},
}

// AgenciasSonEjemplo marca que la lista de agencias no está confirmada.
const AgenciasSonEjemplo = true

func buscar(lista []Opcion, codigo string) (Opcion, bool) {
	for _, o := range lista {
		if o.Codigo == codigo {
			return o, true
		}
	}
	return Opcion{}, false
}

// RegionPorCodigo busca una región.
func RegionPorCodigo(c string) (Opcion, bool) { return buscar(Regiones, c) }

// AgenciaPorCodigo busca una agencia.
func AgenciaPorCodigo(c string) (Opcion, bool) { return buscar(Agencias, c) }
