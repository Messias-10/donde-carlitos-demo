// Llamadas a la API desde el servidor de Next (SSR). Solo se importa en
// Server Components. El navegador, en cambio, llama a /api/* por el mismo
// origen y los rewrites de next.config.ts lo mandan al servicio Go.
import { existsSync } from "node:fs";
import path from "node:path";
import { rutaImagen } from "./imagenes";
import type { ConfigPublica, Pagina, Producto, ProductoVista, Seguimiento } from "./tipos";

const API = (process.env.API_INTERNAL_URL || "http://api:8080").replace(/\/+$/, "");

export class ErrorApiServidor extends Error {
  constructor(public status: number) {
    super(`API respondió ${status}`);
  }
}

async function getJSON<T>(ruta: string, extra: Record<string, string> = {}): Promise<T> {
  const res = await fetch(API + ruta, {
    cache: "no-store",
    headers: { Accept: "application/json", ...extra },
    signal: AbortSignal.timeout(8000),
  });
  if (!res.ok) throw new ErrorApiServidor(res.status);
  return (await res.json()) as T;
}

// ---------------------------------------------------------------- imágenes

const RE_RUTA_LOCAL = /^\/[A-Za-z0-9._~/-]{1,300}$/;
const existe = new Map<string, boolean>();

const RE_IMAGEN_API = /^\/api\/imagenes\/[A-Za-z0-9._-]{1,120}$/;

/**
 * Fotos subidas desde el panel: las sirve la API en /api/imagenes/{archivo} (mismo origen).
 * Fotos de la semilla: se usan si el archivo existe en public/. Si no, null (se muestra la placa).
 */
export function imagenLocal(url: string | null): string | null {
  if (url && RE_IMAGEN_API.test(url) && !url.includes("..")) return url;
  if (!url || !RE_RUTA_LOCAL.test(url) || url.includes("..") || url.startsWith("//")) return null;
  let ok = existe.get(url);
  if (ok === undefined) {
    ok = existsSync(path.join(process.cwd(), "public", url));
    existe.set(url, ok);
  }
  return ok ? url : null;
}

function vista(p: Producto): ProductoVista {
  return { ...p, imagen: imagenLocal(rutaImagen(p.imagen_url)) };
}

// ---------------------------------------------------------------- catálogo

export interface FiltroCatalogo {
  q?: string;
  categoria?: string;
  oferta?: boolean;
  recien_llegado?: boolean;
  temporada?: string;
  orden?: "relevancia" | "precio_asc" | "precio_desc" | "recientes";
  page?: number;
  page_size?: number;
}

export async function listarProductos(f: FiltroCatalogo = {}): Promise<Pagina<ProductoVista>> {
  const q = new URLSearchParams();
  if (f.q) q.set("q", f.q);
  if (f.categoria) q.set("categoria", f.categoria);
  if (f.oferta) q.set("oferta", "true");
  if (f.recien_llegado) q.set("recien_llegado", "true");
  if (f.temporada) q.set("temporada", f.temporada);
  if (f.orden && f.orden !== "relevancia") q.set("orden", f.orden);
  if (f.page && f.page > 1) q.set("page", String(f.page));
  if (f.page_size) q.set("page_size", String(f.page_size));
  const qs = q.toString();
  const pag = await getJSON<Pagina<Producto>>("/api/productos" + (qs ? `?${qs}` : ""));
  return { ...pag, items: pag.items.map(vista) };
}


// ---------------------------------------------------------------- v2


export async function configPublica(): Promise<ConfigPublica> {
  return getJSON<ConfigPublica>("/api/config-publica");
}

/**
 * null = token inválido, vencido o inexistente (la API responde 404 igual en los tres casos).
 * Lanza ErrorApiServidor(429) si hay demasiadas consultas. Se reenvía la IP del cliente
 * (X-Forwarded-For) para que el límite de la API no lo compartan todos (seguridad M4).
 */
export async function seguimiento(token: string, ipCliente?: string): Promise<Seguimiento | null> {
  if (!/^[A-Za-z0-9_-]{22}$/.test(token)) return null;
  const extra: Record<string, string> = {};
  if (ipCliente && /^[0-9a-fA-F:.]{2,45}$/.test(ipCliente)) extra["X-Forwarded-For"] = ipCliente;
  try {
    return await getJSON<Seguimiento>(`/api/seguimiento/${token}`, extra);
  } catch (e) {
    if (e instanceof ErrorApiServidor && e.status === 404) return null;
    throw e;
  }
}

/** GET /api/vitrina: hasta 8 productos que el dueño marcó, en su orden. */
export async function vitrinaPublica(): Promise<ProductoVista[]> {
  const r = await getJSON<Producto[] | { items: Producto[] }>("/api/vitrina");
  const items = Array.isArray(r) ? r : r.items;
  return items.slice(0, 8).map(vista);
}
