import type { Producto } from "@/lib/tipos";

const TEMPORADA: Record<string, string> = {
  navidad: "NAVIDAD",
  escolar: "ESCOLAR",
  halloween: "HALLOWEEN",
  fiestas_patrias: "FIESTAS PATRIAS",
  invierno: "INVIERNO",
  verano: "VERANO",
};

interface Badge {
  texto: string;
  clase: string;
}

/** Badges de estilo-carlitos §8.2, en orden de prioridad. */
export function badgesDe(p: Producto): Badge[] {
  const out: Badge[] = [];
  if (p.badges.includes("remate")) {
    out.push({ texto: "REMATE", clase: "bg-negro text-oro ring-2 ring-oro font-display tracking-wide" });
  } else if (p.badges.includes("oferta")) {
    out.push({ texto: "OFERTA", clase: "bg-magenta text-blanco font-extrabold" });
  }
  if (p.badges.includes("recien_llegado")) {
    out.push({ texto: "¡LLEGÓ!", clase: "bg-oro text-negro font-extrabold" });
  }
  if (p.temporada && TEMPORADA[p.temporada]) {
    out.push({ texto: TEMPORADA[p.temporada], clase: "bg-rojo text-blanco font-extrabold" });
  }
  if (p.precio_mayor) {
    out.push({ texto: "PRECIO X MAYOR", clase: "bg-blanco text-negro ring-1 ring-negro font-extrabold" });
  }
  return out;
}

export function Badges({ producto, max = 2 }: { producto: Producto; max?: number }) {
  const lista = badgesDe(producto).slice(0, max);
  if (lista.length === 0) return null;
  return (
    <ul className="flex flex-wrap gap-1.5" aria-label="Etiquetas">
      {lista.map((b) => (
        <li key={b.texto} className={`rounded-full px-2.5 py-1 text-[12px] leading-none ${b.clase}`}>
          {b.texto}
        </li>
      ))}
    </ul>
  );
}
