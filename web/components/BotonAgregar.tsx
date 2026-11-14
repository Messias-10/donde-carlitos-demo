"use client";

import { useState } from "react";
import { agregar, MAX_CANTIDAD, type ProductoParaCarrito } from "@/lib/carrito";
import { linkWhatsApp, msjAgotado } from "@/lib/whatsapp";
import { IconoMas, IconoMenos, IconoWhatsApp } from "./Iconos";
import { avisar } from "./Toast";

/** "Agregar" / "¡Yo quiero mi promo!" / "¿Vuelve? Pregúntanos" (copy-landing.md §15.1). */
export function BotonAgregar({
  producto,
  agotado,
  promo,
  conCantidad = false,
  ancho = true,
}: {
  producto: ProductoParaCarrito;
  agotado: boolean;
  promo: boolean;
  conCantidad?: boolean;
  ancho?: boolean;
}) {
  const [cantidad, setCantidad] = useState(1);
  const [error, setError] = useState("");

  if (agotado) {
    return (
      <a
        href={linkWhatsApp(msjAgotado(producto.nombre, producto.codigo))}
        target="_blank"
        rel="noopener noreferrer"
        className={`inline-flex min-h-11 items-center justify-center gap-2 rounded-[10px] border-2 border-negro bg-blanco px-3 text-[14px] font-extrabold text-negro ${ancho ? "w-full" : ""}`}
      >
        <IconoWhatsApp size={18} className="text-whatsapp-osc" />
        ¿Vuelve? Pregúntanos
      </a>
    );
  }

  function onAgregar() {
    const ok = agregar(producto, cantidad);
    if (!ok) {
      setError("Tu pedido ya tiene 50 productos distintos. Envíalo y arma otro.");
      return;
    }
    setError("");
    avisar(`¡Listo! Agregaste ${producto.nombre}.`);
  }

  return (
    <div className={conCantidad ? "flex flex-col gap-3" : ""}>
      {conCantidad ? (
        <div className="flex items-center gap-3">
          <span id={`cant-${producto.id}`} className="text-[15px] font-semibold">
            Cantidad
          </span>
          <div className="flex items-center rounded-[10px] border-2 border-negro bg-blanco" role="group" aria-labelledby={`cant-${producto.id}`}>
            <button
              type="button"
              className="grid size-11 place-items-center disabled:opacity-30"
              onClick={() => setCantidad((c) => Math.max(1, c - 1))}
              disabled={cantidad <= 1}
              aria-label="Restar uno"
            >
              <IconoMenos size={18} />
            </button>
            <input
              type="number"
              inputMode="numeric"
              min={1}
              max={MAX_CANTIDAD}
              value={cantidad}
              onChange={(e) => {
                const n = Math.floor(Number(e.target.value));
                setCantidad(Number.isFinite(n) ? Math.min(Math.max(1, n), MAX_CANTIDAD) : 1);
              }}
              className="h-11 w-14 border-x-2 border-negro bg-transparent text-center text-[16px] font-bold tabular-nums [appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none"
              aria-label="Cantidad"
            />
            <button
              type="button"
              className="grid size-11 place-items-center disabled:opacity-30"
              onClick={() => setCantidad((c) => Math.min(MAX_CANTIDAD, c + 1))}
              disabled={cantidad >= MAX_CANTIDAD}
              aria-label="Sumar uno"
            >
              <IconoMas size={18} />
            </button>
          </div>
        </div>
      ) : null}
      <button
        type="button"
        onClick={onAgregar}
        className={`inline-flex min-h-11 items-center justify-center rounded-[10px] px-3 text-[15px] font-extrabold transition-transform active:scale-95 ${
          promo ? "bg-magenta text-blanco" : "bg-negro text-oro"
        } ${ancho ? "w-full" : ""}`}
      >
        {promo ? "¡Yo quiero mi promo!" : "Agregar"}
        <span className="sr-only"> {producto.nombre}</span>
      </button>
      {error ? (
        <p className="mt-2 text-[13px] font-semibold text-rojo" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  );
}
