"use client";

import { ErrorCliente } from "@/lib/api-cliente";
import type { ProductoAdmin } from "@/lib/tipos";

const MAX_LADO = 1280;
const CALIDAD = 0.85;
export const MAX_BYTES_ORIGINAL = 25 * 1024 * 1024;

/**
 * Achica la foto del celular antes de subirla (una foto de 4–8 MB queda en
 * ~200 KB). Si el navegador no puede decodificarla, se sube tal cual.
 */
export async function prepararFoto(archivo: File): Promise<Blob> {
  if (!archivo.type.startsWith("image/")) throw new Error("tipo");
  try {
    const bitmap = await createImageBitmap(archivo, { imageOrientation: "from-image" });
    const escala = Math.min(1, MAX_LADO / Math.max(bitmap.width, bitmap.height));
    const w = Math.round(bitmap.width * escala);
    const h = Math.round(bitmap.height * escala);
    const canvas = document.createElement("canvas");
    canvas.width = w;
    canvas.height = h;
    const ctx = canvas.getContext("2d");
    if (!ctx) return archivo;
    ctx.drawImage(bitmap, 0, 0, w, h);
    bitmap.close();
    const blob = await new Promise<Blob | null>((ok) => canvas.toBlob(ok, "image/jpeg", CALIDAD));
    return blob ?? archivo;
  } catch {
    return archivo;
  }
}

/** POST /api/admin/productos/{id}/imagen (multipart, campo "imagen"). */
export async function subirFoto(id: number, foto: Blob): Promise<ProductoAdmin | null> {
  const fd = new FormData();
  fd.append("imagen", foto, foto.type === "image/jpeg" ? "foto.jpg" : "foto");
  let res: Response;
  try {
    res = await fetch(`/api/admin/productos/${id}/imagen`, { method: "POST", body: fd, credentials: "same-origin" });
  } catch {
    throw new ErrorCliente(0, "sin_conexion", "No se pudo subir la foto: revisa tu conexión.");
  }
  const data = (await res.json().catch(() => null)) as (ProductoAdmin & { error?: { codigo: string; mensaje: string } }) | null;
  if (!res.ok) {
    throw new ErrorCliente(res.status, data?.error?.codigo ?? "interno", data?.error?.mensaje ?? "No se pudo subir la foto.");
  }
  return data && typeof data.id === "number" ? data : null;
}
