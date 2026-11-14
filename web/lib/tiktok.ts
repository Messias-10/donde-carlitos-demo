const RE_ID = /\/video\/(\d{8,25})/;

/** Id numérico de un link de TikTok (null si es un link corto sin id). */
export function idTikTok(url: string | null): string | null {
  if (!url) return null;
  const m = RE_ID.exec(url);
  return m ? m[1] : null;
}
