// Plantillas (a), (b) y (c) de propuesta/mensajes-whatsapp.md, texto exacto.
// Las (d1)–(d8) del panel las arma la API (GET /api/admin/pedidos/{id}/mensajes).
import { clp } from "./formato";
import type { Canal } from "./tipos";

const DEMO = "*[PRUEBA · DEMO]*";

/** Limpieza de texto libre (§2, regla 8). */
export function limpiar(s = ""): string {
  return s
    .replace(/[\r\n\t]+/g, " ")
    .replace(/[*_~`]/g, "")
    .replace(/\s{2,}/g, " ")
    .trim();
}

function conDemo(lineas: string[], demo: boolean): string {
  return (demo ? [DEMO, ...lineas] : lineas).join("\n");
}

/** (a) Pedidos Santiago / retiro en tienda. Termina en ": " a propósito. */
export function mensajeSantiago(demo: boolean): string {
  return conDemo(
    [
      "Hola Donde Carlitos Emporio 🛍️",
      "Vengo de la web. Quiero hacer un pedido para *retiro en tienda (Santiago)*.",
      "Lo que busco es: ",
    ],
    demo,
  );
}

/** (b) Pedidos a regiones. Termina en ": " a propósito. */
export function mensajeRegiones(demo: boolean): string {
  return conDemo(
    [
      "Hola Donde Carlitos Emporio 📦",
      "Vengo de la web. Soy de *regiones* y quiero cotizar un pedido con envío.",
      "Mi comuna y lo que busco: ",
    ],
    demo,
  );
}

export interface DatosVitrina {
  nombre: string;
  codigo: string;
  precio: number;
  precio_mayor: number | null;
  unidad_mayor: string | null;
  /** "dd/mm" o vacío. */
  fecha_precio: string;
}

/** (c) "Pedir este por WhatsApp" desde la vitrina. */
export function mensajeProducto(p: DatosVitrina, canal: Canal, demo: boolean): string {
  const lineas = ["Hola Donde Carlitos Emporio 🛍️", "Vengo de la web y me interesa este producto:", ""];
  lineas.push(`*${limpiar(p.nombre)}*`);
  lineas.push(p.fecha_precio ? `Precio publicado el ${p.fecha_precio}: ${clp(p.precio)}` : `Precio publicado: ${clp(p.precio)}`);
  if (p.precio_mayor) {
    const u = limpiar(p.unidad_mayor ?? "");
    lineas.push(`Por mayor: ${clp(p.precio_mayor)} c/u${u ? ` (${u})` : ""}`);
  }
  lineas.push(`Ref. ${limpiar(p.codigo)}`, "");
  if (canal === "santiago") {
    lineas.push("Es para *retiro en tienda (Santiago)*.", "¿Les queda? Quiero esta cantidad: ");
  } else {
    lineas.push("Es con *envío a regiones*.", "¿Les queda? Mi comuna y la cantidad: ");
  }
  return conDemo(lineas, demo);
}

/** Mensaje de la página de seguimiento (copy-landing.md §12 / arquitectura §B). */
export function mensajeSeguimiento(numero: string, demo: boolean): string {
  return conDemo([`Hola, les escribo por mi pedido ${limpiar(numero)}.`], demo);
}

/** Link wa.me con el texto codificado (encodeURIComponent). */
export function linkWhatsApp(numero: string, texto: string): string {
  return `https://wa.me/${numero.replace(/\D/g, "")}?text=${encodeURIComponent(texto)}`;
}
