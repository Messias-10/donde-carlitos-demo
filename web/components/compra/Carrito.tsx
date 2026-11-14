"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import { cambiarCantidad, cotizar, MAX_CANTIDAD, quitar, useCarrito } from "@/lib/carrito";
import { clp } from "@/lib/formato";
import { ETIQUETA_TIPO } from "@/lib/precio";
import { IconoMas, IconoMenos } from "../Iconos";
import { MinimoPorConfirmar } from "../Precio";

const noop = () => () => {};

/** true solo en el navegador ya hidratado (el carrito vive en localStorage). */
export function useHidratado(): boolean {
  return useSyncExternalStore(
    noop,
    () => true,
    () => false,
  );
}

export function PasosMini() {
  return (
    <ol className="grid grid-cols-3 gap-2 text-center text-[12px] font-semibold">
      {["Elige", "Envía", "Te confirmamos"].map((t, i) => (
        <li key={t} className={`rounded-[10px] p-2 ${i === 0 ? "bg-negro text-oro" : "bg-blanco text-negro ring-1 ring-negro/10"}`}>
          <span className="block font-display text-[20px] leading-none">{i + 1}</span>
          {t}
        </li>
      ))}
    </ol>
  );
}

export function Carrito() {
  const items = useCarrito();
  const hidratado = useHidratado();
  const { lineas, total } = cotizar(items);

  if (!hidratado) {
    return <div className="h-64 animate-pulse rounded-[12px] bg-blanco/60" aria-hidden />;
  }

  if (lineas.length === 0) {
    return (
      <div className="rounded-[16px] bg-blanco p-6 text-center">
        <p className="text-[17px] font-bold">Tu pedido está vacío. ¡Mira lo que llegó!</p>
        <Link href="/#ofertas" className="mt-4 inline-flex h-12 items-center rounded-[12px] bg-negro px-5 font-extrabold text-oro">
          Ver ofertas
        </Link>
      </div>
    );
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_360px] lg:items-start">
      <ul className="space-y-3">
        {lineas.map((l) => {
          const falta = l.precio_mayor && l.mayor_minimo && l.cantidad < l.mayor_minimo ? l.mayor_minimo - l.cantidad : 0;
          return (
            <li key={l.id} className="rounded-[12px] bg-blanco p-4 ring-1 ring-negro/5">
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <Link href={`/producto/${l.slug}`} className="text-[15px] font-bold leading-snug hover:underline">
                    {l.nombre}
                  </Link>
                  <p className="text-[12px] text-gris">Cód. {l.codigo}</p>
                </div>
                <p className="shrink-0 font-display text-[22px] leading-none">{clp(l.subtotal)}</p>
              </div>
              <p className="mt-2 text-[14px] text-negro/80">
                {l.cantidad} x {clp(l.precioAplicado)} ({ETIQUETA_TIPO[l.tipo]})
              </p>
              {falta > 0 ? (
                <p className="mt-1 text-[13px] font-semibold text-oro-texto">
                  Te faltan {falta} para el precio por mayor ({clp(l.precio_mayor!)} c/u).
                  {l.mayor_minimo_confirmado === false ? <MinimoPorConfirmar /> : null}
                </p>
              ) : l.tipo === "mayor" ? (
                <p className="mt-1 text-[13px] font-semibold text-verde-stock">Precio por mayor aplicado.</p>
              ) : null}
              <div className="mt-3 flex items-center justify-between gap-3">
                <div className="flex items-center rounded-[10px] border-2 border-negro" role="group" aria-label={`Cantidad de ${l.nombre}`}>
                  <button
                    type="button"
                    className="grid size-11 place-items-center disabled:opacity-30"
                    onClick={() => cambiarCantidad(l.id, l.cantidad - 1)}
                    disabled={l.cantidad <= 1}
                    aria-label="Restar uno"
                  >
                    <IconoMenos size={18} />
                  </button>
                  <input
                    type="number"
                    inputMode="numeric"
                    min={1}
                    max={MAX_CANTIDAD}
                    value={l.cantidad}
                    onChange={(e) => cambiarCantidad(l.id, Number(e.target.value))}
                    className="h-11 w-14 border-x-2 border-negro bg-transparent text-center text-[16px] font-bold tabular-nums [appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none"
                    aria-label="Cantidad"
                  />
                  <button
                    type="button"
                    className="grid size-11 place-items-center disabled:opacity-30"
                    onClick={() => cambiarCantidad(l.id, l.cantidad + 1)}
                    disabled={l.cantidad >= MAX_CANTIDAD}
                    aria-label="Sumar uno"
                  >
                    <IconoMas size={18} />
                  </button>
                </div>
                <button type="button" onClick={() => quitar(l.id)} className="min-h-11 px-2 text-[14px] font-bold text-rojo underline underline-offset-2">
                  Quitar
                </button>
              </div>
            </li>
          );
        })}
      </ul>

      <aside className="space-y-4 rounded-[16px] bg-negro p-5 text-blanco lg:sticky lg:top-24">
        <div className="flex items-end justify-between gap-3">
          <p className="text-[15px]">
            Total productos <span className="block text-[12px] text-gris-claro">(no incluye envío)</span>
          </p>
          <p className="placa-precio text-[34px]">{clp(total)}</p>
        </div>
        <p className="text-[13px] text-gris-claro">La web no cobra. El pago lo coordinamos contigo por WhatsApp o en el local.</p>
        <Link href="/pedido" className="flex h-[52px] items-center justify-center rounded-[14px] bg-oro text-[17px] font-extrabold text-negro">
          Continuar
        </Link>
        <Link href="/catalogo" className="flex h-11 items-center justify-center text-[15px] font-bold text-oro underline underline-offset-2">
          Seguir comprando
        </Link>
        <PasosMini />
      </aside>
    </div>
  );
}
