import { NEGOCIO } from "./negocio";

// Mensajes de consulta: propuesta/mensaje-pedido-whatsapp.md §5 (texto exacto).
// El mensaje del PEDIDO no se arma aquí: viene de la API (mensaje_whatsapp).
export const MSJ_CONSULTA = "Hola Donde Carlitos Emporio, vengo de la web y tengo una consulta: ";

export function msjAgotado(nombre: string, codigo: string): string {
  return `Hola Donde Carlitos Emporio, vengo de la web. ¿Les queda ${limpiar(nombre)} (cód. ${limpiar(codigo)})?`;
}

export function msjAdminCliente(nombreCliente: string, numero: string): string {
  return `Hola ${limpiar(nombreCliente)}, te escribimos de Donde Carlitos Emporio por tu pedido ${numero}.`;
}

/** Limpieza de texto libre (plantilla §3). */
export function limpiar(s = ""): string {
  return s
    .replace(/[\r\n\t]+/g, " ")
    .replace(/[*_~`]/g, "")
    .replace(/\s{2,}/g, " ")
    .trim();
}

/** Link wa.me al número de pedidos, con el texto codificado. */
export function linkWhatsApp(texto?: string, numero: string = NEGOCIO.whatsapp): string {
  const n = numero.replace(/\D/g, "");
  return texto ? `https://wa.me/${n}?text=${encodeURIComponent(texto)}` : `https://wa.me/${n}`;
}

/** Link al chat con un cliente desde el panel (solo dígitos del teléfono). */
export function linkCliente(telefono: string, texto: string): string {
  return linkWhatsApp(texto, telefono);
}
