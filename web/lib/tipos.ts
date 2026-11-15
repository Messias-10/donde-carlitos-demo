// Tipos de la API (api/internal/http/dto.go; docs/openapi.yaml).

export type Disponibilidad = "disponible" | "quedan_pocos" | "agotado";
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
  oferta_vigente?: boolean;
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

export interface ErrorApi {
  error: {
    codigo: string;
    mensaje: string;
    campos?: Record<string, string>;
    detalles?: unknown;
  };
}

// ---------------------------------------------------------------- v2: landing y seguimiento

export interface ConfigPublica {
  nombre_tienda: string;
  whatsapp_santiago: string;
  whatsapp_regiones: string | null;
  whatsapp_regiones_confirmado: boolean;
  direccion: string;
  horario: string;
  horario_confirmado: boolean;
  site_url: string;
}

export type Canal = "santiago" | "regiones";
export type EstadoPedido =
  | "cotizado"
  | "esperando_pago"
  | "pagado"
  | "preparando"
  | "despachado"
  | "listo_retiro"
  | "entregado"
  | "cancelado";
export type TipoPrecio = "detalle" | "oferta" | "mayor" | "libre";
export type FormaPago = "transferencia" | "en_tienda";
export type MedioPago = "transferencia" | "efectivo" | "debito" | "credito" | "otro";
export type Entrega = "retiro" | "envio";

export interface Seguimiento {
  numero: string;
  estado: EstadoPedido;
  entrega: Entrega;
  agencia_nombre: string | null;
  n_seguimiento: string | null;
  linea_tiempo: { estado: EstadoPedido; fecha: string }[];
}

// ---------------------------------------------------------------- v2: panel

export interface UsuarioRef {
  id: number;
  nombre: string;
}

export interface ItemPedido {
  producto_id: number | null;
  codigo: string | null;
  descripcion: string;
  cantidad: number;
  precio_unitario: number | null;
  tipo_precio: TipoPrecio;
  subtotal: number | null;
}

export interface CambioEstado {
  estado_anterior: EstadoPedido | null;
  estado_nuevo: EstadoPedido;
  usuario: UsuarioRef | null;
  nota: string;
  creado_en: string;
}

export interface Pedido {
  id: number;
  numero: string;
  canal: Canal;
  estado: EstadoPedido;
  estado_desde: string;
  es_demo: boolean;
  nombre: string;
  telefono: string;
  entrega: Entrega;
  region: string | null;
  region_nombre: string | null;
  comuna: string | null;
  agencia: string | null;
  agencia_nombre: string | null;
  n_seguimiento: string | null;
  items: ItemPedido[];
  forma_pago: FormaPago;
  subtotal: number;
  descuento: number;
  costo_envio: number;
  envio_por_pagar: boolean;
  total: number;
  total_manual: boolean;
  motivo_cancelacion?: string;
  nota_interna?: string;
  pago: { medio: MedioPago; monto: number; fecha: string; verificado_por: UsuarioRef | null } | null;
  stock_descontado?: boolean;
  registrado_por: UsuarioRef | null;
  seguimiento_url: string;
  transiciones_permitidas: EstadoPedido[];
  whatsapp_cliente_url: string;
  creado_en: string;
  actualizado_en: string;
  historial?: CambioEstado[];
}

export interface MensajeCliente {
  clave:
    | "cotizacion"
    | "datos_transferencia"
    | "recordatorio_pago"
    | "pago_recibido"
    | "despachado"
    | "listo_retiro"
    | "cancelado"
    | "abrir_chat";
  titulo: string;
  texto: string;
  whatsapp_url: string;
  sugerido: boolean;
}

export interface Hoy {
  fecha: string;
  contadores: {
    esperando_pago_mas_24h: number;
    pagados_sin_despachar: number;
    listos_retiro: number;
    creados_hoy_santiago: number;
    creados_hoy_regiones: number;
  };
  esperando_pago_mas_24h: Pedido[];
  pagados_sin_despachar: Pedido[];
  listos_retiro: Pedido[];
  creados_hoy: Pedido[];
}

export interface Cliente {
  telefono: string;
  nombre: string | null;
  total_pedidos: number;
  total_comprado: number;
  pedidos: Pedido[];
}

export interface FilaResumen {
  clave: string | null;
  nombre: string | null;
  pedidos: number;
  monto: number;
}

export interface Resumen {
  desde: string;
  hasta: string;
  totales: { pedidos: number; monto: number; cancelados: number };
  por_canal: FilaResumen[];
  por_estado: FilaResumen[];
  por_agencia: FilaResumen[];
}

export type Rol = "admin" | "vendedor";

export interface Usuario {
  id: number;
  usuario: string;
  nombre: string;
  rol: Rol;
  activo: boolean;
  creado_en: string;
}

export interface UsuarioConClave {
  usuario: Usuario;
  clave_temporal: string;
}

export interface SesionPanel {
  id: number;
  usuario: string;
  nombre: string;
  rol: Rol;
  expira_en: string;
}
