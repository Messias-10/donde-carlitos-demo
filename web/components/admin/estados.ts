import type { EstadoPedido } from "@/lib/tipos";

// copy-landing.md §16.
export const ETIQUETA_ESTADO: Record<EstadoPedido, string> = {
  nuevo: "Nuevo",
  confirmado: "Confirmado",
  listo: "Listo para retiro",
  despachado: "Despachado",
  entregado: "Entregado",
  cancelado: "Cancelado",
};

export const PESTANAS: { estado: EstadoPedido | ""; texto: string }[] = [
  { estado: "nuevo", texto: "Nuevos" },
  { estado: "confirmado", texto: "Confirmados" },
  { estado: "listo", texto: "Listos para retiro" },
  { estado: "despachado", texto: "Despachados" },
  { estado: "entregado", texto: "Entregados" },
  { estado: "cancelado", texto: "Cancelados" },
  { estado: "", texto: "Todos" },
];

export const ACCION: Record<EstadoPedido, string> = {
  nuevo: "Nuevo",
  confirmado: "Confirmar",
  listo: "Marcar listo para retiro",
  despachado: "Marcar despachado",
  entregado: "Marcar entregado",
  cancelado: "Cancelar pedido",
};

export const COLOR_ESTADO: Record<EstadoPedido, string> = {
  nuevo: "bg-fluor text-negro",
  confirmado: "bg-oro text-negro",
  listo: "bg-negro text-oro",
  despachado: "bg-negro text-oro",
  entregado: "bg-verde-stock text-blanco",
  cancelado: "bg-rojo text-blanco",
};
