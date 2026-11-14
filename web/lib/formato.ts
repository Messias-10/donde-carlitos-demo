// Moneda chilena sin Intl (según el motor, Intl puede omitir el punto de miles).
export function clp(n: number): string {
  const v = Math.round(n);
  return (v < 0 ? "-$" : "$") + miles(Math.abs(v));
}

/** Solo la parte numérica, para pintar el "$" más chico. */
export function miles(n: number): string {
  return Math.round(n).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
}

const RE_FECHA = /^\d{4}-\d{2}-\d{2}/;

/** "2026-11-13" -> "13/11" */
export function ddmm(iso?: string | null): string {
  if (!iso || !RE_FECHA.test(iso)) return "";
  return `${iso.slice(8, 10)}/${iso.slice(5, 7)}`;
}

/** "2026-11-13" -> "13/11/2026" */
export function ddmmaaaa(iso?: string | null): string {
  if (!iso || !RE_FECHA.test(iso)) return "";
  return `${iso.slice(8, 10)}/${iso.slice(5, 7)}/${iso.slice(0, 4)}`;
}

export function plural(n: number, uno: string, varios: string): string {
  return n === 1 ? uno : varios;
}

/** "hace 5 min", "hace 2 h", "hace 3 días". */
export function haceCuanto(iso: string, ahora = Date.now()): string {
  const t = Date.parse(iso);
  if (Number.isNaN(t)) return "";
  const min = Math.max(0, Math.round((ahora - t) / 60000));
  if (min < 1) return "recién";
  if (min < 60) return `hace ${min} min`;
  const h = Math.round(min / 60);
  if (h < 24) return `hace ${h} h`;
  const d = Math.round(h / 24);
  return `hace ${d} ${plural(d, "día", "días")}`;
}

/** Fecha y hora local de Santiago: "14/11 15:32". */
export function fechaHora(iso: string): string {
  const t = new Date(iso);
  if (Number.isNaN(t.getTime())) return "";
  return new Intl.DateTimeFormat("es-CL", {
    timeZone: "America/Santiago",
    day: "2-digit",
    month: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(t);
}
