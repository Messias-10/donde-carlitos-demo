// Estados de pedido: etiquetas (copy-landing.md §12 y P4) y colores (estilo-carlitos §8.6).
import type { EstadoPedido } from "./tipos";

export const ESTADOS: EstadoPedido[] = [
  "cotizado",
  "esperando_pago",
  "pagado",
  "preparando",
  "despachado",
  "listo_retiro",
  "entregado",
  "cancelado",
];

/** Etiqueta del panel. */
export const ETIQUETA: Record<EstadoPedido, string> = {
  cotizado: "Cotizado",
  esperando_pago: "Esperando pago",
  pagado: "Pagado",
  preparando: "Preparando",
  despachado: "Despachado",
  listo_retiro: "Listo para retiro",
  entregado: "Entregado",
  cancelado: "Cancelado",
};

/** Botón de acción siguiente (copy-landing.md P4). */
export const ACCION: Record<EstadoPedido, string> = {
  cotizado: "Volver a cotizado",
  esperando_pago: "Marcar esperando pago",
  pagado: "Marcar pagado",
  preparando: "Marcar preparando",
  despachado: "Marcar despachado",
  listo_retiro: "Marcar listo para retiro",
  entregado: "Marcar entregado",
  cancelado: "Cancelar pedido",
};

/** Título y explicación de la página de seguimiento. */
export const SEGUIMIENTO: Record<EstadoPedido, { titulo: string; texto: (x: { agencia?: string | null; direccion: string }) => string }> = {
  cotizado: { titulo: "Cotizado", texto: () => "Te mandamos la cotización por WhatsApp." },
  esperando_pago: { titulo: "Esperando tu pago", texto: () => "Cuando transfieras, mándanos el comprobante por el mismo chat." },
  pagado: { titulo: "Pago recibido", texto: () => "Ya confirmamos tu pago. ¡Gracias!" },
  preparando: { titulo: "Preparando tu pedido", texto: () => "Estamos armando tu pedido." },
  despachado: {
    titulo: "En camino",
    texto: ({ agencia }) => (agencia ? `Tu pedido salió por ${agencia}.` : "Tu pedido ya salió."),
  },
  listo_retiro: { titulo: "Listo para retiro", texto: ({ direccion }) => `Te esperamos en ${direccion}.` },
  entregado: { titulo: "Entregado", texto: () => "¡Gracias por comprar en Donde Carlitos Emporio!" },
  cancelado: { titulo: "Cancelado", texto: () => "Este pedido se canceló. Si tienes dudas, escríbenos." },
};

/** Clases del chip de estado (fondo + texto, contraste AA). */
export const CHIP: Record<EstadoPedido, string> = {
  cotizado: "bg-st-cotizado-bg text-st-cotizado-fg",
  esperando_pago: "bg-st-espera-bg text-st-espera-fg",
  pagado: "bg-st-pagado-bg text-st-pagado-fg",
  preparando: "bg-st-prep-bg text-st-prep-fg",
  despachado: "bg-st-desp-bg text-st-desp-fg",
  listo_retiro: "bg-st-retiro-bg text-st-retiro-fg",
  entregado: "bg-st-entregado-bg text-st-entregado-fg",
  cancelado: "bg-st-cancel-bg text-st-cancel-fg",
};

export const PUNTO: Record<EstadoPedido, string> = {
  cotizado: "bg-st-cotizado-fg",
  esperando_pago: "bg-st-espera-fg",
  pagado: "bg-st-pagado-fg",
  preparando: "bg-st-prep-fg",
  despachado: "bg-st-desp-fg",
  listo_retiro: "bg-st-retiro-fg",
  entregado: "bg-st-entregado-bg",
  cancelado: "bg-st-cancel-fg",
};
