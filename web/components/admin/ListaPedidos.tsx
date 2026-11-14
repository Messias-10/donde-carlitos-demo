"use client";

import Link from "next/link";
import { useState } from "react";
import { clp, haceCuanto } from "@/lib/formato";
import type { EstadoPedido, Pagina, PedidoAdmin } from "@/lib/tipos";
import { COLOR_ESTADO, ETIQUETA_ESTADO, PESTANAS } from "./estados";
import { useApi } from "./useApi";

export function ListaPedidos() {
  const [estado, setEstado] = useState<EstadoPedido | "">("nuevo");
  const [page, setPage] = useState(1);
  const q = new URLSearchParams({ page: String(page), page_size: "20" });
  if (estado) q.set("estado", estado);
  const lista = useApi<Pagina<PedidoAdmin>>(`/api/admin/pedidos?${q}`);
  const conteo = useApi<Pagina<PedidoAdmin>>("/api/admin/pedidos?estado=nuevo&page_size=1");
  const datos = lista.data ?? null;
  const nuevos = conteo.data?.total ?? null;
  const error = lista.error ?? "";
  const cargando = lista.cargando;
  const cargar = () => {
    lista.recargar();
    conteo.recargar();
  };

  return (
    <div>
      <div className="flex items-center justify-between gap-3">
        <h1 className="titular text-[28px]">Pedidos</h1>
        <button type="button" onClick={cargar} className="min-h-11 px-2 text-[14px] font-bold underline underline-offset-2">
          Actualizar
        </button>
      </div>

      <nav aria-label="Estados" className="scroll-fila -mx-4 mt-3 flex gap-2 overflow-x-auto px-4 pb-1">
        {PESTANAS.map((t) => {
          const activo = t.estado === estado;
          return (
            <button
              key={t.texto}
              type="button"
              aria-pressed={activo}
              onClick={() => {
                setEstado(t.estado);
                setPage(1);
              }}
              className={`h-11 shrink-0 rounded-full px-4 text-[14px] font-bold ${activo ? "bg-negro text-oro" : "bg-blanco text-negro ring-1 ring-negro/15"}`}
            >
              {t.texto}
              {t.estado === "nuevo" && nuevos !== null ? ` (${nuevos})` : ""}
            </button>
          );
        })}
      </nav>

      {error ? (
        <p role="alert" className="mt-4 rounded-[10px] border-2 border-rojo p-3 text-[15px] font-semibold text-rojo">
          {error}
        </p>
      ) : null}

      {cargando && !datos ? <div className="mt-4 h-40 animate-pulse rounded-[12px] bg-blanco/70" aria-hidden /> : null}

      {datos ? (
        datos.items.length === 0 ? (
          <p className="mt-6 rounded-[12px] bg-blanco p-5 text-center text-[15px]">No hay pedidos en este estado.</p>
        ) : (
          <ul className={`mt-4 space-y-3 ${cargando ? "opacity-60" : ""}`} aria-busy={cargando}>
            {datos.items.map((p) => (
              <li key={p.id}>
                <Link href={`/admin/pedidos/${p.id}`} className="block rounded-[12px] bg-blanco p-4 ring-1 ring-negro/5 hover:ring-oro">
                  <div className="flex items-center justify-between gap-2">
                    <span className="font-display text-[22px] leading-none">{p.numero}</span>
                    <span className={`rounded-full px-2.5 py-1 text-[12px] font-extrabold ${COLOR_ESTADO[p.estado]}`}>
                      {ETIQUETA_ESTADO[p.estado]}
                    </span>
                  </div>
                  <p className="mt-2 text-[15px] font-semibold">{p.nombre}</p>
                  <p className="text-[14px] text-gris">
                    {p.entrega === "retiro" ? "Retiro" : `Envío a ${p.comuna ?? ""}, ${p.region_nombre ?? p.region ?? ""}`}
                  </p>
                  <div className="mt-2 flex items-center justify-between text-[14px]">
                    <strong>{clp(p.total)}</strong>
                    <span className="text-gris">{haceCuanto(p.creado_en)}</span>
                  </div>
                  {p.es_demo ? <p className="mt-1 text-[12px] font-semibold text-gris">Pedido de prueba (demo)</p> : null}
                </Link>
              </li>
            ))}
          </ul>
        )
      ) : null}

      {datos && datos.total_pages > 1 ? (
        <div className="mt-5 flex items-center justify-center gap-3">
          <button
            type="button"
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
            className="h-11 rounded-[10px] bg-blanco px-4 font-bold ring-1 ring-negro/15 disabled:opacity-40"
          >
            Anterior
          </button>
          <span className="text-[14px] text-gris">
            {datos.page} de {datos.total_pages}
          </span>
          <button
            type="button"
            disabled={page >= datos.total_pages}
            onClick={() => setPage((p) => p + 1)}
            className="h-11 rounded-[10px] bg-negro px-4 font-bold text-oro disabled:opacity-40"
          >
            Siguiente
          </button>
        </div>
      ) : null}
    </div>
  );
}
