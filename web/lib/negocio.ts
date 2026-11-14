// Datos públicos del negocio. Fuente de cada dato en propuesta/copy-landing.md.
// Lo marcado porConfirmar se muestra con la etiqueta "dato por confirmar".

const numeroEnv = (process.env.NEXT_PUBLIC_WHATSAPP_PEDIDOS || "").replace(/\D/g, "");

export const NEGOCIO = {
  nombre: "Donde Carlitos Emporio",
  firma: "Familia Emporio",
  direccion: "Alameda 2432",
  comuna: "Santiago Centro",
  ciudad: "Santiago",
  region: "Región Metropolitana",
  // Número de destino de los pedidos en formato wa.me. Configurable; no es secreto.
  whatsapp: numeroEnv.length >= 8 && numeroEnv.length <= 15 ? numeroEnv : "56921145368",
  // Número oficial del negocio, para mostrar (bios oficiales).
  telefonoVisible: "+56 9 2114 5368",
  telefonoE164: "+56921145368",
  tiktok: { usuario: "@mayte_emporio", url: "https://www.tiktok.com/@mayte_emporio" },
  instagram: { usuario: "@mayte_emporioo", url: "https://www.instagram.com/mayte_emporioo/" },
  mapsUrl: "https://www.google.com/maps/search/?api=1&query=Alameda+2432%2C+Santiago%2C+Chile",
  mapsEmbed: "https://maps.google.com/maps?q=Alameda%202432%2C%20Santiago%2C%20Chile&z=16&output=embed",
} as const;

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
