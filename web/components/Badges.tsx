import type { Producto } from "@/lib/tipos";

/** Un solo badge por tarjeta de vitrina (estilo-carlitos §8.2), en orden de prioridad. */
export function badgeVitrina(p: Producto): { texto: string; clase: string } | null {
  const vigente = p.precio_final < p.precio;
  // B-2: precio de promo por comentario: no es precio de detalle.
  if (p.precio_es_promo) return { texto: "PROMO DEL VIDEO", clase: "bg-magenta text-blanco font-extrabold" };
  if (vigente && p.etiqueta_oferta === "remate") {
    return { texto: "REMATE", clase: "bg-negro text-oro ring-2 ring-oro font-display tracking-wide" };
  }
  if (vigente) return { texto: "OFERTA", clase: "bg-magenta text-blanco font-extrabold" };
  if (p.badges.includes("recien_llegado")) return { texto: "¡LLEGÓ!", clase: "bg-oro text-negro font-display tracking-wide" };
  if (p.precio_mayor) return { texto: "PRECIO X MAYOR", clase: "bg-blanco text-negro ring-1 ring-negro font-extrabold" };
  return null;
}

export function BadgeVitrina({ producto }: { producto: Producto }) {
  const b = badgeVitrina(producto);
  if (!b) return null;
  return <span className={`inline-block rounded-full px-2.5 py-1 text-[12px] leading-none ${b.clase}`}>{b.texto}</span>;
}
