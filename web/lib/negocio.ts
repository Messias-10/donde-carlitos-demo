// Datos públicos del negocio. Fuente de cada dato en propuesta/copy-landing.md.
// Los números de WhatsApp son públicos (no son secretos) y vienen por configuración:
// la landing funciona aunque la API no responda (arquitectura-landing.md §A).

const digitos = (v: string | undefined) => (v || "").replace(/\D/g, "");
const valido = (n: string) => n.length >= 8 && n.length <= 15;

const santiagoEnv = digitos(process.env.NEXT_PUBLIC_WHATSAPP_SANTIAGO || process.env.NEXT_PUBLIC_WHATSAPP_PEDIDOS);
const regionesEnv = digitos(process.env.NEXT_PUBLIC_WHATSAPP_REGIONES);

/** Demo: los mensajes (a), (b) y (c) llevan la primera línea *[PRUEBA · DEMO]*. */
export const ES_DEMO = (process.env.NEXT_PUBLIC_DEMO ?? "true") !== "false";

export const NEGOCIO = {
  nombre: "Donde Carlitos Emporio",
  firma: "Familia Emporio",
  direccion: "Alameda 2432",
  comuna: "Santiago Centro",
  ciudad: "Santiago",
  region: "Región Metropolitana",
  tiktok: { usuario: "@mayte_emporio", url: "https://www.tiktok.com/@mayte_emporio" },
  instagram: { usuario: "@mayte_emporioo", url: "https://www.instagram.com/mayte_emporioo/" },
  mapsUrl: "https://www.google.com/maps/search/?api=1&query=Alameda+2432%2C+Santiago%2C+Chile",
  videoContenedor: "https://www.tiktok.com/@mayte_emporio/video/7441724226857454904",
} as const;

export interface Numeros {
  santiago: string;
  /** null = no confirmado: se usa el de Santiago y se marca "por confirmar". */
  regiones: string | null;
}

export const NUMEROS_ENV: Numeros = {
  santiago: valido(santiagoEnv) ? santiagoEnv : "56921145368",
  regiones: valido(regionesEnv) ? regionesEnv : null,
};

/** "56921145368" -> "+56 9 2114 5368". */
export function telefonoVisible(n: string): string {
  const m = /^56(9)(\d{4})(\d{4})$/.exec(n);
  return m ? `+56 ${m[1]} ${m[2]} ${m[3]}` : `+${n}`;
}

// Horario de temporada publicado en nov 2026 [CONFIRMAR CON EL DUEÑO].
export const HORARIO = {
  porConfirmar: true,
  fuente: "Publicado en noviembre de 2026.",
  semana: { abre: "9:00", cierra: "18:00" },
  finde: { abre: "9:00", cierra: "16:00" },
} as const;

/** Horario que corresponde hoy en Santiago (sin feriados: no hay fuente). */
export function horarioDeHoy(fecha = new Date()): { abre: string; cierra: string } {
  const dia = new Intl.DateTimeFormat("en-US", { weekday: "short", timeZone: "America/Santiago" }).format(fecha);
  return dia === "Sat" || dia === "Sun" ? HORARIO.finde : HORARIO.semana;
}
