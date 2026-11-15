"use client";

import Link from "next/link";
import { useEffect, useRef } from "react";
import { CHIP, ETIQUETA } from "@/lib/estados";
import { clp } from "@/lib/formato";
import type { Canal, EstadoPedido, Pedido } from "@/lib/tipos";

export const INPUT =
  "mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px] focus:border-negro aria-[invalid=true]:border-rojo";
export const BTN_PRIMARIO =
  "flex min-h-12 w-full items-center justify-center rounded-[12px] bg-negro px-4 text-[16px] font-extrabold text-blanco disabled:opacity-60";
export const BTN_SECUNDARIO =
  "flex min-h-12 w-full items-center justify-center rounded-[12px] border-2 border-negro bg-blanco px-4 text-[16px] font-bold text-negro disabled:opacity-60";

export function ChipEstado({ estado, grande = false }: { estado: EstadoPedido; grande?: boolean }) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 whitespace-nowrap rounded-full font-extrabold ${grande ? "px-3 py-1.5 text-[14px]" : "px-2.5 py-1 text-[12px]"} ${CHIP[estado]}`}
    >
      <span className="inline-block size-2 rounded-full bg-current" aria-hidden />
      {ETIQUETA[estado]}
    </span>
  );
}

/** Mini etiqueta de canal con los colores de los botones de canal (estilo-carlitos §8.7). */
export function EtiquetaCanal({ canal }: { canal: Canal }) {
  return canal === "santiago" ? (
    <span className="inline-block whitespace-nowrap rounded-full border-2 border-whatsapp px-2 py-px text-[12px] font-bold text-negro">
      Santiago
    </span>
  ) : (
    <span className="inline-block whitespace-nowrap rounded-full bg-whatsapp-osc px-2 py-0.5 text-[12px] font-bold text-blanco">Regiones</span>
  );
}

/** En listas el teléfono va enmascarado: +56 9 •••• 5368. */
export function telefonoEnmascarado(t: string): string {
  const m = /^\+56(\d)(\d{4})(\d{4})$/.exec(t);
  return m ? `+56 ${m[1]} •••• ${m[3]}` : t.replace(/\d(?=\d{4})/g, "•");
}

/** En la ficha va completo: +56 9 1234 5678. */
export function telefonoLegible(t: string): string {
  const m = /^\+56(\d)(\d{4})(\d{4})$/.exec(t);
  return m ? `+56 ${m[1]} ${m[2]} ${m[3]}` : t;
}

export function horasDesde(iso: string, ahora = Date.now()): number {
  const t = Date.parse(iso);
  return Number.isNaN(t) ? 0 : Math.max(0, Math.floor((ahora - t) / 3_600_000));
}

export function hace(iso: string): string {
  const h = horasDesde(iso);
  if (h < 1) return "hace menos de 1 h";
  if (h < 48) return `hace ${h} h`;
  return `hace ${Math.floor(h / 24)} días`;
}

/** Tarjeta de pedido para listas (Hoy, Buscar, Cliente). */
export function TarjetaPedido({
  p,
  alerta = false,
  detalle,
  accion,
}: {
  p: Pedido;
  alerta?: boolean;
  detalle?: string;
  accion?: string;
}) {
  return (
    <Link
      href={`/admin/pedidos/${p.id}`}
      className={`block rounded-[12px] bg-blanco p-4 ring-1 ring-negro/5 hover:ring-negro/30 ${alerta ? "border-l-4 border-rojo" : ""}`}
    >
      <div className="flex flex-wrap items-center gap-2">
        <span className="font-display text-[22px] leading-none">{p.numero}</span>
        <ChipEstado estado={p.estado} />
        <EtiquetaCanal canal={p.canal} />
        {p.es_demo ? <span className="text-[12px] font-semibold text-gris">demo</span> : null}
      </div>
      <p className="mt-2 text-[15px] font-semibold">{p.nombre || "Sin nombre"}</p>
      <p className="text-[14px] text-gris">
        {telefonoEnmascarado(p.telefono)}
        {p.entrega === "envio" ? ` · Envío${p.comuna ? ` a ${p.comuna}` : ""}${p.agencia_nombre ? ` · ${p.agencia_nombre}` : ""}` : " · Retiro"}
      </p>
      <div className="mt-2 flex items-center justify-between gap-2 text-[14px]">
        <strong>{clp(p.total)}</strong>
        <span className={alerta ? "font-bold text-rojo" : "text-gris"}>{detalle ?? hace(p.creado_en)}</span>
      </div>
      {accion ? <p className="mt-2 text-[14px] font-bold underline underline-offset-2">{accion}</p> : null}
    </Link>
  );
}

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

export function Aviso({ tipo = "error", children }: { tipo?: "error" | "ok" | "info"; children: React.ReactNode }) {
  const c =
    tipo === "error"
      ? "border-2 border-rojo text-rojo"
      : tipo === "ok"
        ? "bg-st-pagado-bg text-st-pagado-fg"
        : "bg-blanco text-negro ring-1 ring-negro/10";
  return (
    <p role={tipo === "error" ? "alert" : "status"} className={`rounded-[10px] p-3 text-[15px] font-semibold ${c}`}>
      {children}
    </p>
  );
}

export function BotonNuevoPedido() {
  return (
    <Link
      href="/admin/pedidos/nuevo"
      className="fixed bottom-[max(16px,env(safe-area-inset-bottom))] right-4 z-30 flex h-14 items-center rounded-full bg-oro px-5 text-[16px] font-extrabold text-negro shadow-[0_8px_24px_-8px_rgba(0,0,0,0.5)]"
    >
      + Nuevo pedido
    </Link>
  );
}
