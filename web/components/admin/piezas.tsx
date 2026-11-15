// Piezas del panel sin estado: sirven en Server Components (Hoy) y en componentes cliente.
import Link from "next/link";
import type { ReactNode } from "react";
import { CHIP, ETIQUETA } from "@/lib/estados";
import { clp } from "@/lib/formato";
import type { Canal, EstadoPedido, Pedido } from "@/lib/tipos";

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

/** "hace menos de 1 h", "hace 5 h", "hace 1 día", "hace 3 días" (QA M-4). */
export function hace(iso: string): string {
  const h = horasDesde(iso);
  if (h < 1) return "hace menos de 1 h";
  if (h < 48) return `hace ${h} h`;
  const d = Math.floor(h / 24);
  return `hace ${d} ${d === 1 ? "día" : "días"}`;
}

/** Tarjeta de pedido para listas. Las acciones van fuera del link (no se anidan links). */
export function TarjetaPedido({
  p,
  alerta = false,
  detalle,
  acciones,
}: {
  p: Pedido;
  alerta?: boolean;
  detalle?: string;
  acciones?: ReactNode;
}) {
  return (
    <div className={`overflow-hidden rounded-[12px] bg-blanco ring-1 ring-negro/5 ${alerta ? "border-l-4 border-rojo" : ""}`}>
      <Link href={`/admin/pedidos/${p.id}`} className="block p-4 hover:bg-panel">
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
      </Link>
      {acciones ? <div className="flex flex-wrap gap-2 border-t border-negro/10 px-4 py-3">{acciones}</div> : null}
    </div>
  );
}

export function Aviso({ tipo = "error", children }: { tipo?: "error" | "ok" | "info"; children: ReactNode }) {
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

/** Link de acción secundaria dentro de una tarjeta (44 px). */
export function AccionTarjeta({ href, children, peligro = false }: { href: string; children: ReactNode; peligro?: boolean }) {
  return (
    <Link
      href={href}
      className={`inline-flex min-h-11 items-center rounded-[10px] border-2 px-3 text-[14px] font-bold ${peligro ? "border-rojo text-rojo" : "border-negro text-negro"}`}
    >
      {children}
    </Link>
  );
}
