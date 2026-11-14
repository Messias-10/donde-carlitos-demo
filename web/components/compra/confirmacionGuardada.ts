"use client";

import type { PedidoCreado } from "@/lib/tipos";
import type { DatosConfirmacion } from "./Confirmacion";

// La confirmación se guarda en sessionStorage (solo esta pestaña; se borra al cerrarla)
// para que no se pierda al recargar o al volver de WhatsApp. La URL solo lleva el número.
const CLAVE = "dc_confirmacion_v1";
export const RE_NUMERO = /^DC-\d{4,8}$/;

export function guardarConfirmacion(p: PedidoCreado) {
  const d: DatosConfirmacion = {
    numero: p.numero,
    total: p.total,
    entrega: p.entrega,
    es_demo: p.es_demo,
    mensaje_whatsapp: p.mensaje_whatsapp,
    whatsapp_url: p.whatsapp_url,
  };
  try {
    window.sessionStorage.setItem(CLAVE, JSON.stringify(d));
  } catch {
    // Sin almacenamiento: la página /pedido/listo muestra la versión corta.
  }
}

export function leerConfirmacion(numero: string): DatosConfirmacion | null {
  try {
    const raw = window.sessionStorage.getItem(CLAVE);
    if (!raw) return null;
    const d = JSON.parse(raw) as DatosConfirmacion;
    if (d.numero !== numero || typeof d.mensaje_whatsapp !== "string" || typeof d.whatsapp_url !== "string") return null;
    return d;
  } catch {
    return null;
  }
}
