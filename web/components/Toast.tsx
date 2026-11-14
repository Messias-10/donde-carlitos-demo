"use client";

import Link from "next/link";
import { useEffect, useSyncExternalStore } from "react";

// Aviso "¡Listo! Agregaste {nombre}." (copy-landing.md §15.1).
type Aviso = { id: number; texto: string } | null;

let aviso: Aviso = null;
let contador = 0;
const oyentes = new Set<() => void>();

export function avisar(texto: string) {
  contador += 1;
  aviso = { id: contador, texto };
  oyentes.forEach((f) => f());
}

function cerrar() {
  aviso = null;
  oyentes.forEach((f) => f());
}

function suscribir(f: () => void) {
  oyentes.add(f);
  return () => oyentes.delete(f);
}

/** Cuenta de "agregados" en esta visita (para el pop del contador). */
export function useContadorAvisos(): number {
  return useSyncExternalStore(
    suscribir,
    () => contador,
    () => 0,
  );
}

export function Toast() {
  const a = useSyncExternalStore(
    suscribir,
    () => aviso,
    () => null,
  );
  useEffect(() => {
    if (!a) return;
    const t = window.setTimeout(cerrar, 3500);
    return () => window.clearTimeout(t);
  }, [a]);

  return (
    <div
      aria-live="polite"
      role="status"
      className="pointer-events-none fixed inset-x-0 bottom-[92px] z-50 flex justify-center px-4"
    >
      {a ? (
        <div
          key={a.id}
          className="pointer-events-auto flex w-full max-w-sm animate-subir items-center justify-between gap-3 rounded-[12px] bg-negro px-4 py-3 text-blanco shadow-lg ring-1 ring-oro/40"
        >
          <p className="text-[14px] leading-snug">{a.texto}</p>
          <Link
            href="/carrito"
            className="shrink-0 rounded-[10px] bg-oro px-3 py-2 text-[13px] font-extrabold text-negro"
            onClick={cerrar}
          >
            Ver mi pedido
          </Link>
        </div>
      ) : null}
    </div>
  );
}
