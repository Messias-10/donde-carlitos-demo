"use client";

import { useEffect } from "react";
import { agregar, type ProductoParaCarrito } from "@/lib/carrito";
import { avisar } from "./Toast";

function leer(raw: string | undefined): ProductoParaCarrito | null {
  if (!raw) return null;
  try {
    const p = JSON.parse(raw) as ProductoParaCarrito;
    if (!Number.isInteger(p.id) || typeof p.nombre !== "string" || typeof p.precio !== "number") return null;
    return p;
  } catch {
    return null;
  }
}

/** Un solo listener para todos los botones [data-agregar] de la página. */
export function AgregarDelegado() {
  useEffect(() => {
    function onClick(e: MouseEvent) {
      const b = (e.target as Element | null)?.closest<HTMLElement>("[data-agregar]");
      if (!b) return;
      const p = leer(b.dataset.agregar);
      if (!p) return;
      if (agregar(p, 1)) avisar(`¡Listo! Agregaste ${p.nombre}.`);
      else avisar("Tu pedido ya tiene 50 productos distintos. Envíalo y arma otro.");
    }
    document.addEventListener("click", onClick);
    return () => document.removeEventListener("click", onClick);
  }, []);
  return null;
}
