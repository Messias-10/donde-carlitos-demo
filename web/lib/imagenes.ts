// Recortes de producto con sufijo de hash (QA ronda 2, I-6): cambiar el nombre invalida
// la caché de /_next/image si se rehace un recorte. La semilla de la API apunta a
// /productos/tt-<videoId>.webp; este mapa la traduce al archivo real.
const MAPA: Record<string, string> = {
  "7366788833616088325": "2a302df0",
  "7418593012223184133": "64bcdad6",
  "7592989161653685522": "77d4598f",
  "7706661076712082695": "386a3690",
  "7706661259298540806": "8d150e79",
  "7706661393260367110": "f0b02687",
  "7706669405505129735": "fe4d36af",
  "7706670593428573458": "76bf0a6a",
  "7706913997378948370": "bf1d9781",
  "7706993060130950408": "f4f55a04",
  "7707040126286695732": "9e17d67b",
  "7707047766052048134": "42c57cbc",
  "7707047875322121480": "648e5cfb",
  "7707047904761990408": "adb6c749",
  "7707053973999471879": "30cfe4cb",
  "7707296451398864136": "edad1a76",
  "7707300664988142866": "0256c8d8",
  "7707411807194926343": "a1672ae5",
  "7707415244787731720": "dbe71251",
  "7707634987562061074": "68192ea7",
  "7707647176440171783": "33511b59",
  "7707649267464293639": "d45c5146",
};

const RE_SEMILLA = /^\/productos\/tt-(\d{8,25})\.webp$/;

/** "/productos/tt-<id>.webp" -> "/productos/tt-<id>-<hash>.webp" si hay recorte; si no, igual. */
export function rutaImagen(url: string | null): string | null {
  if (!url) return null;
  const m = RE_SEMILLA.exec(url);
  if (m && MAPA[m[1]]) return `/productos/tt-${m[1]}-${MAPA[m[1]]}.webp`;
  return url;
}
