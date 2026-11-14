// Tipos de la API (api/internal/http/dto.go; docs/openapi.yaml).

export type Disponibilidad = "disponible" | "quedan_pocos" | "agotado";
export type TipoPrecio = "detalle" | "oferta" | "mayor";
export type EstadoPedido = "nuevo" | "confirmado" | "listo" | "despachado" | "entregado" | "cancelado";
export type Temporada = "navidad" | "escolar" | "fiestas_patrias" | "halloween" | "invierno" | "verano";

export interface Categoria {
  id: number;
  slug: string;
  nombre: string;
  total_productos: number;
}

export interface Producto {
  id: number;
  slug: string;
  codigo: string;
  nombre: string;
  descripcion: string;
  categoria: { id: number; slug: string; nombre: string };
  precio: number;
  precio_oferta: number | null;
  precio_final: number;
  etiqueta_oferta: "oferta" | "remate" | null;
  precio_mayor: number | null;
  mayor_minimo: number | null;
  /** false = el mínimo por mayor es [EJEMPLO]: se muestra "mínimo por confirmar". */
  mayor_minimo_confirmado?: boolean;
  unidad_mayor: string | null;
  /** Fecha (UTC) en que la tienda cambió el precio desde el panel; null = precio del video. */
  precio_actualizado_en?: string | null;
  /** Fecha AAAA-MM-DD hasta la que vale la oferta (la API ya excluye las vencidas). */
  oferta_hasta?: string | null;
  disponibilidad: Disponibilidad;
  badges: string[];
  temporada: Temporada | null;
  fecha_llegada: string | null;
  imagen_url: string | null;
  tiktok_url: string | null;
  fuente: { url: string; fecha: string | null } | null;
  es_ejemplo: boolean;
}

/** Producto con la imagen ya verificada en el servidor (null = no hay archivo). */
export interface ProductoVista extends Producto {
  imagen: string | null;
}

export interface ProductoAdmin extends Producto {
  stock: number;
  umbral_pocos: number;
  activo: boolean;
  creado_en: string;
  actualizado_en: string;
}

export interface Pagina<T> {
  items: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface Opcion {
  codigo: string;
  nombre: string;
}

export interface OpcionesEntrega {
  retiro: { direccion: string };
  regiones: Opcion[];
  agencias: Opcion[];
  agencias_es_ejemplo: boolean;
}

export interface ItemPedido {
  producto_id: number;
  codigo: string;
  nombre: string;
  cantidad: number;
  precio_unitario: number;
  tipo_precio: TipoPrecio;
  subtotal: number;
}

/** Datos de entrega, planos dentro del pedido (DatosEntregaPedido). */
export interface DatosEntrega {
  entrega: "retiro" | "envio";
  region: string | null;
  region_nombre: string | null;
  comuna: string | null;
  agencia: string | null;
  agencia_nombre: string | null;
}

export interface PedidoCreado extends DatosEntrega {
  numero: string;
  estado: EstadoPedido;
  nombre: string;
  comentario: string;
  total: number;
  es_demo: boolean;
  items: ItemPedido[];
  mensaje_whatsapp: string;
  whatsapp_url: string;
  creado_en: string;
}

export interface CambioEstado {
  estado_anterior: EstadoPedido | null;
  estado_nuevo: EstadoPedido;
  admin: string | null;
  nota: string;
  creado_en: string;
}

export interface PedidoAdmin extends DatosEntrega {
  id: number;
  numero: string;
  estado: EstadoPedido;
  es_demo: boolean;
  nombre: string;
  telefono: string;
  comentario: string;
  total: number;
  items: ItemPedido[];
  transiciones_permitidas: EstadoPedido[];
  whatsapp_cliente_url: string;
  creado_en: string;
  actualizado_en: string;
  historial?: CambioEstado[];
}

export interface ErrorApi {
  error: {
    codigo: string;
    mensaje: string;
    campos?: Record<string, string>;
    detalles?: unknown;
  };
}
