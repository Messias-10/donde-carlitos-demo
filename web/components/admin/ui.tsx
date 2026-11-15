"use client";

import { useEffect, useRef } from "react";

export {
  Aviso,
  BotonNuevoPedido,
  ChipEstado,
  EtiquetaCanal,
  hace,
  horasDesde,
  TarjetaPedido,
  telefonoEnmascarado,
  telefonoLegible,
} from "./piezas";

export const INPUT =
  "mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px] focus:border-negro aria-[invalid=true]:border-rojo";
export const BTN_PRIMARIO =
  "flex min-h-12 w-full items-center justify-center rounded-[12px] bg-negro px-4 text-[16px] font-extrabold text-blanco disabled:opacity-60";
export const BTN_SECUNDARIO =
  "flex min-h-12 w-full items-center justify-center rounded-[12px] border-2 border-negro bg-blanco px-4 text-[16px] font-bold text-negro disabled:opacity-60";

/** Diálogo modal (confirmaciones sensibles, estilo-carlitos §8.7). */
export function Dialogo({
  abierto,
  onCerrar,
  titulo,
  children,
}: {
  abierto: boolean;
  onCerrar: () => void;
  titulo: string;
  children: React.ReactNode;
}) {
  const ref = useRef<HTMLDialogElement>(null);
  useEffect(() => {
    const d = ref.current;
    if (!d) return;
    if (abierto && !d.open) d.showModal();
    if (!abierto && d.open) d.close();
  }, [abierto]);
  return (
    <dialog
      ref={ref}
      aria-labelledby="dialogo-titulo"
      onClose={onCerrar}
      className="m-0 mt-auto w-full max-w-none rounded-t-[20px] bg-blanco p-0 text-negro backdrop:bg-negro/60 sm:m-auto sm:max-w-md sm:rounded-[20px]"
    >
      <div className="space-y-4 px-5 pb-[max(20px,env(safe-area-inset-bottom))] pt-5">
        <h2 id="dialogo-titulo" className="text-[20px] font-extrabold leading-snug">
          {titulo}
        </h2>
        {children}
      </div>
    </dialog>
  );
}
