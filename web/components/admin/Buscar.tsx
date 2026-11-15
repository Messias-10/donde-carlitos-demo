"use client";

import { useState } from "react";
import { ESTADOS, ETIQUETA } from "@/lib/estados";
import { clp } from "@/lib/formato";
import type { Cliente, Pagina, Pedido } from "@/lib/tipos";
import { Aviso, BotonNuevoPedido, INPUT, TarjetaPedido } from "./ui";
import { useApi } from "./useApi";

interface Filtro {
  q: string;
  canal: string;
  estado: string;
  desde: string;
  hasta: string;
  page: number;
}

const VACIO: Filtro = { q: "", canal: "", estado: "", desde: "", hasta: "", page: 1 };

function query(f: Filtro): string {
  const p = new URLSearchParams({ page: String(f.page), page_size: "20" });
  if (f.q.trim()) p.set("q", f.q.trim().slice(0, 60));
  if (f.canal) p.set("canal", f.canal);
  if (f.estado) p.set("estado", f.estado);
  if (f.desde) p.set("desde", f.desde);
  if (f.hasta) p.set("hasta", f.hasta);
  return p.toString();
}

/** Solo dígitos (con espacios, + o -): es un teléfono. */
const esTelefono = (q: string) => /^[\d\s+-]{8,}$/.test(q.trim()) && q.replace(/\D/g, "").length >= 8;

function FichaCliente({ telefono }: { telefono: string }) {
  const { data } = useApi<Cliente>(`/api/admin/clientes?telefono=${encodeURIComponent(telefono)}`);
  if (!data || data.total_pedidos === 0) return null;
  return (
    <p className="rounded-[12px] bg-st-prep-bg p-3 text-[15px] text-st-prep-fg">
      <strong>{data.nombre || "Cliente sin nombre"}</strong> · {data.total_pedidos} {data.total_pedidos === 1 ? "pedido" : "pedidos"} · compró{" "}
      {clp(data.total_comprado)}
      {data.pedidos[0] ? ` · último: ${data.pedidos[0].numero}` : ""}
    </p>
  );
}

export function Buscar() {
  const [borrador, setBorrador] = useState<Filtro>(VACIO);
  const [f, setF] = useState<Filtro>(VACIO);
  const { data, error, cargando } = useApi<Pagina<Pedido>>(`/api/admin/pedidos?${query(f)}`);
  const set = (k: keyof Filtro, v: string) => setBorrador((b) => ({ ...b, [k]: v }));

  return (
    <div className="space-y-4">
      <h1 className="titular text-[28px]">Buscar</h1>
      <form
        role="search"
        className="space-y-3 rounded-[12px] bg-blanco p-4"
        onSubmit={(e) => {
          e.preventDefault();
          setF({ ...borrador, page: 1 });
        }}
      >
        <div>
          <label htmlFor="b-q" className="text-[14px] font-bold">
            Teléfono, n.º de pedido o nombre
          </label>
          <input id="b-q" type="search" value={borrador.q} onChange={(e) => set("q", e.target.value)} maxLength={60} className={INPUT} />
        </div>
        <details>
          <summary className="min-h-11 py-2 text-[15px] font-bold underline underline-offset-2">Más filtros</summary>
          <div className="mt-2 grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="b-canal" className="text-[14px] font-bold">
                Canal
              </label>
              <select id="b-canal" value={borrador.canal} onChange={(e) => set("canal", e.target.value)} className={INPUT}>
                <option value="">Todos</option>
                <option value="santiago">Santiago</option>
                <option value="regiones">Regiones</option>
              </select>
            </div>
            <div>
              <label htmlFor="b-estado" className="text-[14px] font-bold">
                Estado
              </label>
              <select id="b-estado" value={borrador.estado} onChange={(e) => set("estado", e.target.value)} className={INPUT}>
                <option value="">Todos</option>
                {ESTADOS.map((e) => (
                  <option key={e} value={e}>
                    {ETIQUETA[e]}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label htmlFor="b-desde" className="text-[14px] font-bold">
                Desde
              </label>
              <input id="b-desde" type="date" value={borrador.desde} onChange={(e) => set("desde", e.target.value)} className={INPUT} />
            </div>
            <div>
              <label htmlFor="b-hasta" className="text-[14px] font-bold">
                Hasta
              </label>
              <input id="b-hasta" type="date" value={borrador.hasta} onChange={(e) => set("hasta", e.target.value)} className={INPUT} />
            </div>
          </div>
        </details>
        <button type="submit" className="flex min-h-12 w-full items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-blanco">
          Buscar
        </button>
      </form>

      {esTelefono(f.q) ? <FichaCliente telefono={f.q.trim()} /> : null}
      {error ? <Aviso>{error}</Aviso> : null}

      {data ? (
        data.items.length === 0 ? (
          <p className="rounded-[12px] bg-blanco p-5 text-center text-[15px]">No encontramos pedidos con eso.</p>
        ) : (
          <>
            <p className="text-[14px] text-gris">
              {data.total} {data.total === 1 ? "pedido" : "pedidos"}
            </p>
            <ul className={`space-y-2 ${cargando ? "opacity-60" : ""}`} aria-busy={cargando}>
              {data.items.map((p) => (
                <li key={p.id}>
                  <TarjetaPedido p={p} />
                </li>
              ))}
            </ul>
            {data.total_pages > 1 ? (
              <div className="flex items-center justify-center gap-3">
                <button
                  type="button"
                  disabled={f.page <= 1}
                  onClick={() => setF((x) => ({ ...x, page: x.page - 1 }))}
                  className="h-11 rounded-[10px] bg-blanco px-4 font-bold ring-1 ring-negro/15 disabled:opacity-40"
                >
                  Anterior
                </button>
                <span className="text-[14px] text-gris">
                  {data.page} de {data.total_pages}
                </span>
                <button
                  type="button"
                  disabled={f.page >= data.total_pages}
                  onClick={() => setF((x) => ({ ...x, page: x.page + 1 }))}
                  className="h-11 rounded-[10px] bg-negro px-4 font-bold text-blanco disabled:opacity-40"
                >
                  Siguiente
                </button>
              </div>
            ) : null}
          </>
        )
      ) : (
        <div className="h-40 animate-pulse rounded-[12px] bg-blanco" aria-hidden />
      )}
      <BotonNuevoPedido />
    </div>
  );
}
