// Llamadas del panel desde el servidor de Next (Server Components). Reenvía la
// cookie de sesión que trae la solicitud del navegador: la API valida la sesión
// igual que si llamara el navegador. Solo GET (las escrituras van desde el
// navegador, con Origin, por el rewrite /api/*).
import { cookies } from "next/headers";

const API = (process.env.API_INTERNAL_URL || "http://api:8080").replace(/\/+$/, "");
const NOMBRES = ["dc_sesion", "__Host-dc_sesion"];

export type Resultado<T> = { ok: true; data: T } | { ok: false; status: number; codigo: string; mensaje: string };

export async function getPanel<T>(ruta: string): Promise<Resultado<T>> {
  const jar = await cookies();
  const cookie = NOMBRES.flatMap((n) => {
    const v = jar.get(n)?.value;
    return v && /^[A-Za-z0-9_.~+/=-]{1,300}$/.test(v) ? [`${n}=${v}`] : [];
  }).join("; ");
  if (!cookie) return { ok: false, status: 401, codigo: "no_autenticado", mensaje: "Inicia sesión para continuar." };
  try {
    const res = await fetch(API + ruta, {
      cache: "no-store",
      headers: { Accept: "application/json", Cookie: cookie },
      signal: AbortSignal.timeout(8000),
    });
    const data = (await res.json().catch(() => null)) as unknown;
    if (!res.ok) {
      const e = (data as { error?: { codigo?: string; mensaje?: string } } | null)?.error;
      return { ok: false, status: res.status, codigo: e?.codigo ?? "interno", mensaje: e?.mensaje ?? "No pudimos cargar esto." };
    }
    return { ok: true, data: data as T };
  } catch {
    return { ok: false, status: 0, codigo: "sin_conexion", mensaje: "No pudimos conectar con el servidor. Intenta de nuevo." };
  }
}
