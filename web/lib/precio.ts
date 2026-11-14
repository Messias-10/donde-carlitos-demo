import type { TipoPrecio } from "./tipos";

export interface DatosPrecio {
  precio: number;
  precio_final: number;
  precio_mayor: number | null;
  mayor_minimo: number | null;
}

/**
 * Misma regla que la API (productos.PrecioUnitario): el menor entre detalle,
 * oferta y, si se llega al mínimo, precio por mayor. Es solo un estimado:
 * el total que vale es el que devuelve POST /api/pedidos.
 */
export function precioUnitario(p: DatosPrecio, cantidad: number): { precio: number; tipo: TipoPrecio } {
  let precio = p.precio;
  let tipo: TipoPrecio = "detalle";
  if (p.precio_final > 0 && p.precio_final < precio) {
    precio = p.precio_final;
    tipo = "oferta";
  }
  if (p.precio_mayor && p.mayor_minimo && cantidad >= p.mayor_minimo && p.precio_mayor < precio) {
    precio = p.precio_mayor;
    tipo = "mayor";
  }
  return { precio, tipo };
}

export const ETIQUETA_TIPO: Record<TipoPrecio, string> = {
  detalle: "detalle",
  oferta: "oferta",
  mayor: "x mayor",
};
