// Llamadas a la API desde el servidor de Next (SSR). Solo se importa en
// Server Components. El navegador, en cambio, llama a /api/* por el mismo
// origen y los rewrites de next.config.ts lo mandan al servicio Go.
import { existsSync } from "node:fs";
import path from "node:path";
import type { Categoria, OpcionesEntrega, Pagina, Producto, ProductoVista } from "./tipos";

const API = (process.env.API_INTERNAL_URL || "http://api:8080").replace(/\/+$/, "");

export class ErrorApiServidor extends Error {
  constructor(public status: number) {
    super(`API respondió ${status}`);
  }
}

async function getJSON<T>(ruta: string): Promise<T> {
  const res = await fetch(API + ruta, {
    cache: "no-store",
    headers: { Accept: "application/json" },
    signal: AbortSignal.timeout(8000),
  });
  if (!res.ok) throw new ErrorApiServidor(res.status);
  return (await res.json()) as T;
}

// ---------------------------------------------------------------- imágenes

const RE_RUTA_LOCAL = /^\/[A-Za-z0-9._~/-]{1,300}$/;
const existe = new Map<string, boolean>();

/** Devuelve la ruta si el archivo existe en public/; si no, null (se muestra la placa). */
export function imagenLocal(url: string | null): string | null {
  if (!url || !RE_RUTA_LOCAL.test(url) || url.includes("..") || url.startsWith("//")) return null;
  let ok = existe.get(url);
  if (ok === undefined) {
    ok = existsSync(path.join(process.cwd(), "public", url));
    existe.set(url, ok);
  }
  return ok ? url : null;
}

function vista(p: Producto): ProductoVista {
  return { ...p, imagen: imagenLocal(p.imagen_url) };
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

/** Todo el catálogo activo (para elegir los productos de la landing). */
export async function todosLosProductos(): Promise<ProductoVista[]> {
  const primera = await listarProductos({ page_size: 48 });
  if (primera.total_pages <= 1) return primera.items;
  const resto = await Promise.all(
    Array.from({ length: Math.min(primera.total_pages, 6) - 1 }, (_, i) => listarProductos({ page_size: 48, page: i + 2 })),
  );
  return primera.items.concat(...resto.map((p) => p.items));
}

export async function obtenerProducto(slug: string): Promise<ProductoVista | null> {
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug) || slug.length > 80) return null;
  try {
    return vista(await getJSON<Producto>(`/api/productos/${slug}`));
  } catch (e) {
    if (e instanceof ErrorApiServidor && e.status === 404) return null;
    throw e;
  }
}

export async function listarCategorias(): Promise<Categoria[]> {
  const r = await getJSON<{ items: Categoria[] }>("/api/categorias");
  return r.items;
}

export async function opcionesEntrega(): Promise<OpcionesEntrega> {
  return getJSON<OpcionesEntrega>("/api/opciones-entrega");
}
