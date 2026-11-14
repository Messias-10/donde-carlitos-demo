"use client";

import Link from "next/link";
import { useState } from "react";
import { api } from "@/lib/api-cliente";
import { clp, fechaHora } from "@/lib/formato";
import { ETIQUETA_TIPO } from "@/lib/precio";
import type { EstadoPedido, PedidoAdmin } from "@/lib/tipos";
import { msjAdminCliente } from "@/lib/whatsapp";
import { IconoFlechaIzq, IconoWhatsApp } from "../Iconos";
import { ACCION, COLOR_ESTADO, ETIQUETA_ESTADO } from "./estados";
import { useManejarError } from "./Sesion";
import { useApi } from "./useApi";

function linkCliente(p: PedidoAdmin): string | null {
  // whatsapp_cliente_url = https://wa.me/<dígitos>. Solo se acepta esa forma.
  if (!/^https:\/\/wa\.me\/\d{8,15}$/.test(p.whatsapp_cliente_url)) return null;
  return `${p.whatsapp_cliente_url}?text=${encodeURIComponent(msjAdminCliente(p.nombre, p.numero))}`;
}

export function DetallePedido({ id }: { id: string }) {
  const [error, setError] = useState("");
  const [nota, setNota] = useState("");
  const [guardando, setGuardando] = useState<EstadoPedido | null>(null);
  const [confirmarCancelar, setConfirmarCancelar] = useState(false);
  const [aviso, setAviso] = useState("");
  const manejar = useManejarError();

  const { data: p, error: errorCarga, reemplazar } = useApi<PedidoAdmin>(`/api/admin/pedidos/${encodeURIComponent(id)}`);

  async function cambiar(estado: EstadoPedido) {
    if (!p) return;
    setGuardando(estado);
    setError("");
    setAviso("");
    try {
      const act = await api<PedidoAdmin>(`/api/admin/pedidos/${p.id}/estado`, {
        method: "PATCH",
        body: { estado, nota: nota.trim() },
      });
      reemplazar(act);
      setNota("");
      setConfirmarCancelar(false);
      setAviso(`Pedido marcado como "${ETIQUETA_ESTADO[act.estado]}".`);
    } catch (e) {
      setError(manejar(e));
    } finally {
      setGuardando(null);
    }
  }

  if (!p) {
    return errorCarga ? (
      <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 font-semibold text-rojo">
        {errorCarga}
      </p>
    ) : (
      <div className="h-60 animate-pulse rounded-[12px] bg-blanco/70" aria-hidden />
    );
  }

  const wa = linkCliente(p);
  const acciones = p.transiciones_permitidas.filter((e) => e !== "cancelado");
  const puedeCancelar = p.transiciones_permitidas.includes("cancelado");

  return (
    <div className="space-y-4">
      <Link href="/admin" className="inline-flex min-h-11 items-center gap-1 text-[14px] font-bold">
        <IconoFlechaIzq size={18} /> Pedidos
      </Link>

      <div className="rounded-[12px] bg-blanco p-4">
        <div className="flex items-center justify-between gap-2">
          <h1 className="font-display text-[30px] leading-none">{p.numero}</h1>
          <span className={`rounded-full px-3 py-1 text-[13px] font-extrabold ${COLOR_ESTADO[p.estado]}`}>{ETIQUETA_ESTADO[p.estado]}</span>
        </div>
        {p.es_demo ? <p className="mt-1 text-[12px] font-semibold text-gris">Pedido de prueba (demo)</p> : null}
        <p className="mt-1 text-[13px] text-gris">Creado el {fechaHora(p.creado_en)}</p>

        <dl className="mt-4 grid grid-cols-[auto_1fr] gap-x-4 gap-y-1 text-[15px]">
          <dt className="text-gris">Cliente</dt>
          <dd className="font-semibold">{p.nombre}</dd>
          <dt className="text-gris">WhatsApp</dt>
          <dd className="font-semibold tabular-nums">{p.telefono}</dd>
          <dt className="text-gris">Entrega</dt>
          <dd className="font-semibold">
            {p.entrega === "retiro"
              ? `Retiro en tienda${p.comuna ? ` (viene desde ${p.comuna})` : ""}`
              : `Envío a ${p.comuna ?? ""}, ${p.region_nombre ?? p.region ?? ""} por ${p.agencia_nombre ?? p.agencia ?? ""}`}
          </dd>
        </dl>
        {p.comentario ? (
          <p className="mt-3 whitespace-pre-line rounded-[10px] bg-crema p-3 text-[15px]">
            <span className="font-bold">Comentario: </span>
            {p.comentario}
          </p>
        ) : null}
        {wa ? (
          <a
            href={wa}
            target="_blank"
            rel="noopener noreferrer"
            className="mt-4 flex h-12 items-center justify-center gap-2 rounded-[12px] bg-whatsapp text-[15px] font-extrabold text-negro"
          >
            <IconoWhatsApp size={20} /> Abrir chat con el cliente
          </a>
        ) : null}
      </div>

      <div className="rounded-[12px] bg-blanco p-4">
        <h2 className="text-[17px] font-extrabold">Productos</h2>
        <ul className="mt-2 divide-y divide-negro/10">
          {p.items.map((it) => (
            <li key={it.producto_id} className="flex justify-between gap-3 py-2 text-[15px]">
              <span className="min-w-0">
                <span className="font-semibold">
                  {it.cantidad} x {it.nombre}
                </span>
                <span className="block text-[13px] text-gris">
                  Cód. {it.codigo} · {clp(it.precio_unitario)} c/u ({ETIQUETA_TIPO[it.tipo_precio]})
                </span>
              </span>
              <span className="shrink-0 font-semibold tabular-nums">{clp(it.subtotal)}</span>
            </li>
          ))}
        </ul>
        <p className="mt-2 flex justify-between border-t-2 border-negro pt-2 text-[17px] font-extrabold">
          <span>Total</span>
          <span>{clp(p.total)}</span>
        </p>
      </div>

      {acciones.length > 0 || puedeCancelar ? (
        <div className="space-y-3 rounded-[12px] bg-blanco p-4">
          <h2 className="text-[17px] font-extrabold">Cambiar estado</h2>
          <div>
            <label htmlFor="nota" className="text-[14px] font-bold">
              Nota (opcional)
            </label>
            <input
              id="nota"
              value={nota}
              onChange={(e) => setNota(e.target.value)}
              maxLength={200}
              placeholder="Ej: pagó por transferencia, retira el lunes"
              className="mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 px-3 text-[16px] focus:border-negro"
            />
          </div>
          {acciones.map((e) => (
            <button
              key={e}
              type="button"
              disabled={guardando !== null}
              onClick={() => void cambiar(e)}
              className="flex h-12 w-full items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-oro disabled:opacity-60"
            >
              {guardando === e ? "Guardando…" : ACCION[e]}
            </button>
          ))}
          {puedeCancelar ? (
            confirmarCancelar ? (
              <div className="rounded-[10px] border-2 border-rojo p-3">
                <p className="text-[15px] font-semibold">¿Cancelar el pedido {p.numero}? Las unidades vuelven al stock.</p>
                <div className="mt-3 flex gap-2">
                  <button
                    type="button"
                    disabled={guardando !== null}
                    onClick={() => void cambiar("cancelado")}
                    className="h-11 flex-1 rounded-[10px] bg-rojo font-extrabold text-blanco disabled:opacity-60"
                  >
                    Sí, cancelar
                  </button>
                  <button type="button" onClick={() => setConfirmarCancelar(false)} className="h-11 flex-1 rounded-[10px] font-bold ring-1 ring-negro/20">
                    No
                  </button>
                </div>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setConfirmarCancelar(true)}
                className="h-11 w-full rounded-[12px] text-[15px] font-bold text-rojo underline underline-offset-2"
              >
                Cancelar pedido
              </button>
            )
          ) : null}
        </div>
      ) : null}

      <div aria-live="polite">
        {aviso ? <p className="rounded-[10px] bg-verde-stock/10 p-3 text-[15px] font-semibold text-verde-stock">{aviso}</p> : null}
        {error ? (
          <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 text-[15px] font-semibold text-rojo">
            {error}
          </p>
        ) : null}
      </div>

      {p.historial && p.historial.length > 0 ? (
        <div className="rounded-[12px] bg-blanco p-4">
          <h2 className="text-[17px] font-extrabold">Historial</h2>
          <ul className="mt-2 space-y-1 text-[14px]">
            {p.historial.map((h, i) => (
              <li key={i}>
                {fechaHora(h.creado_en)} · {h.estado_anterior ? ETIQUETA_ESTADO[h.estado_anterior] : "—"} →{" "}
                {ETIQUETA_ESTADO[h.estado_nuevo]} · {h.admin ?? "web"}
                {h.nota ? <span className="block pl-3 text-gris">{h.nota}</span> : null}
              </li>
            ))}
          </ul>
        </div>
      ) : null}
    </div>
  );
}
