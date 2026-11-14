"use client";

// Carrito en localStorage. Solo productos y cantidades: nada de datos personales.
import { useSyncExternalStore } from "react";
import { precioUnitario, type DatosPrecio } from "./precio";
import type { TipoPrecio } from "./tipos";

const CLAVE = "dc_carrito_v1";
const VERSION = 1;
export const MAX_CANTIDAD = 100; // igual que la API (pedidos.MaxCantidadPorItem)
export const MAX_LINEAS = 50; // igual que la API

export interface LineaCarrito extends DatosPrecio {
  id: number;
  slug: string;
  codigo: string;
  nombre: string;
  unidad_mayor: string | null;
  mayor_minimo_confirmado?: boolean;
  categoria: string;
  cantidad: number;
}

interface Estado {
  v: number;
  items: LineaCarrito[];
}

const VACIO: Estado = { v: VERSION, items: [] };
let estado: Estado = VACIO;
let cargado = false;
const oyentes = new Set<() => void>();

function esLinea(x: unknown): x is LineaCarrito {
  if (!x || typeof x !== "object") return false;
  const l = x as Record<string, unknown>;
  return (
    Number.isInteger(l.id) &&
    typeof l.slug === "string" &&
    typeof l.codigo === "string" &&
    typeof l.nombre === "string" &&
    Number.isInteger(l.cantidad) &&
    typeof l.precio === "number" &&
    typeof l.precio_final === "number"
  );
}

function leer(): Estado {
  try {
    const raw = window.localStorage.getItem(CLAVE);
    if (!raw) return VACIO;
    const data = JSON.parse(raw) as Partial<Estado>;
    if (data.v !== VERSION || !Array.isArray(data.items)) return VACIO;
    const items = data.items
      .filter(esLinea)
      .map((l) => ({ ...l, cantidad: Math.min(Math.max(1, l.cantidad), MAX_CANTIDAD) }))
      .slice(0, MAX_LINEAS);
    return { v: VERSION, items };
  } catch {
    return VACIO;
  }
}

function asegurarCarga() {
  if (!cargado && typeof window !== "undefined") {
    estado = leer();
    cargado = true;
  }
}

function guardar(nuevo: Estado) {
  estado = nuevo;
  try {
    window.localStorage.setItem(CLAVE, JSON.stringify(nuevo));
  } catch {
    // Sin almacenamiento (modo privado): el carrito vive solo en memoria.
  }
  oyentes.forEach((f) => f());
}

function suscribir(f: () => void) {
  oyentes.add(f);
  const alCambiarOtraPestana = (e: StorageEvent) => {
    if (e.key === CLAVE) {
      estado = leer();
      f();
    }
  };
  window.addEventListener("storage", alCambiarOtraPestana);
  return () => {
    oyentes.delete(f);
    window.removeEventListener("storage", alCambiarOtraPestana);
  };
}

function snapshot(): Estado {
  asegurarCarga();
  return estado;
}

function snapshotServidor(): Estado {
  return VACIO;
}

export function useCarrito(): LineaCarrito[] {
  return useSyncExternalStore(suscribir, snapshot, snapshotServidor).items;
}

export type ProductoParaCarrito = Omit<LineaCarrito, "cantidad">;

export function agregar(p: ProductoParaCarrito, cantidad = 1): boolean {
  asegurarCarga();
  const items = [...estado.items];
  const i = items.findIndex((l) => l.id === p.id);
  if (i >= 0) {
    items[i] = { ...items[i], ...p, cantidad: Math.min(items[i].cantidad + cantidad, MAX_CANTIDAD) };
  } else {
    if (items.length >= MAX_LINEAS) return false;
    items.push({ ...p, cantidad: Math.min(Math.max(1, cantidad), MAX_CANTIDAD) });
  }
  guardar({ v: VERSION, items });
  return true;
}

export function cambiarCantidad(id: number, cantidad: number) {
  asegurarCarga();
  const c = Math.min(Math.max(1, Math.floor(cantidad) || 1), MAX_CANTIDAD);
  guardar({ v: VERSION, items: estado.items.map((l) => (l.id === id ? { ...l, cantidad: c } : l)) });
}

export function quitar(id: number) {
  asegurarCarga();
  guardar({ v: VERSION, items: estado.items.filter((l) => l.id !== id) });
}

export function quitarVarios(ids: number[]) {
  asegurarCarga();
  const s = new Set(ids);
  guardar({ v: VERSION, items: estado.items.filter((l) => !s.has(l.id)) });
}

export function vaciar() {
  guardar({ v: VERSION, items: [] });
}

export interface LineaCotizada extends LineaCarrito {
  precioAplicado: number;
  tipo: TipoPrecio;
  subtotal: number;
}

export function cotizar(items: LineaCarrito[]): { lineas: LineaCotizada[]; total: number; unidades: number } {
  let total = 0;
  let unidades = 0;
  const lineas = items.map((l) => {
    const { precio, tipo } = precioUnitario(l, l.cantidad);
    const subtotal = precio * l.cantidad;
    total += subtotal;
    unidades += l.cantidad;
    return { ...l, precioAplicado: precio, tipo, subtotal };
  });
  return { lineas, total, unidades };
}
