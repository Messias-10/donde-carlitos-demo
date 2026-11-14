"use client";

// Llamadas del navegador a /api/* (mismo origen; la cookie httpOnly viaja sola).
import type { ErrorApi } from "./tipos";

export class ErrorCliente extends Error {
  constructor(
    public status: number,
    public codigo: string,
    mensaje: string,
    public campos?: Record<string, string>,
  ) {
    super(mensaje);
  }
}

export async function api<T>(ruta: string, opciones: { method?: string; body?: unknown } = {}): Promise<T> {
  let res: Response;
  try {
    res = await fetch(ruta, {
      method: opciones.method ?? "GET",
      headers: opciones.body !== undefined ? { "Content-Type": "application/json", Accept: "application/json" } : { Accept: "application/json" },
      body: opciones.body !== undefined ? JSON.stringify(opciones.body) : undefined,
      credentials: "same-origin",
      cache: "no-store",
    });
  } catch {
    throw new ErrorCliente(0, "sin_conexion", "Parece que no hay internet. Revisa tu conexión.");
  }
  if (res.status === 204) return undefined as T;
  const data = (await res.json().catch(() => null)) as unknown;
  if (!res.ok) {
    const e = (data as ErrorApi | null)?.error;
    throw new ErrorCliente(res.status, e?.codigo ?? "interno", e?.mensaje ?? "Ocurrió un error. Intenta de nuevo en un rato.", e?.campos);
  }
  return data as T;
}
