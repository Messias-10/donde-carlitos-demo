"use client";

import { useState } from "react";
import type { Canal, Hoy as DatosHoy, Pedido } from "@/lib/tipos";
import { Aviso, BotonNuevoPedido, hace, horasDesde, TarjetaPedido } from "./ui";
import { useApi } from "./useApi";

type Filtro = "todos" | Canal;

function Bloque({
  titulo,
  pedidos,
  total,
  alerta,
  detalle,
  accion,
}: {
  titulo: string;
  pedidos: Pedido[];
  total: number;
  alerta?: boolean;
  detalle: (p: Pedido) => string;
  accion?: string;
}) {
  if (pedidos.length === 0) return null;
  return (
    <section className="space-y-2">
      <h2 className="titular text-[20px]">
        {titulo} ({total})
      </h2>
      <ul className="space-y-2">
        {pedidos.map((p) => (
          <li key={p.id}>
            <TarjetaPedido p={p} alerta={alerta} detalle={detalle(p)} accion={accion} />
          </li>
        ))}
      </ul>
    </section>
  );
}

export function Hoy() {
  const { data, error, cargando, recargar } = useApi<DatosHoy>("/api/admin/hoy");
  const [filtro, setFiltro] = useState<Filtro>("todos");

  const f = (l: Pedido[]) => (filtro === "todos" ? l : l.filter((p) => p.canal === filtro));

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between gap-3">
        <h1 className="titular text-[30px]">Hoy</h1>
        <button type="button" onClick={recargar} className="min-h-11 px-2 text-[14px] font-bold underline underline-offset-2">
          Actualizar
        </button>
      </div>

      <div role="group" aria-label="Canal" className="grid grid-cols-3 gap-1 rounded-[12px] bg-blanco p-1 ring-1 ring-negro/10">
        {(
          [
            ["todos", "Todos"],
            ["santiago", "Santiago"],
            ["regiones", "Regiones"],
          ] as const
        ).map(([v, t]) => (
          <button
            key={v}
            type="button"
            aria-pressed={filtro === v}
            onClick={() => setFiltro(v)}
            className={`h-11 rounded-[10px] text-[15px] font-bold ${filtro === v ? "bg-negro text-blanco" : "text-negro"}`}
          >
            {t}
          </button>
        ))}
      </div>

      {error ? <Aviso>{error}</Aviso> : null}
      {cargando && !data ? <div className="h-40 animate-pulse rounded-[12px] bg-blanco" aria-hidden /> : null}

      {data ? (
        <>
          <p className="rounded-[12px] bg-blanco p-3 text-[15px] ring-1 ring-negro/5">
            Hoy: <strong>{data.contadores.creados_hoy_santiago + data.contadores.creados_hoy_regiones}</strong> nuevos (Santiago{" "}
            {data.contadores.creados_hoy_santiago} · Regiones {data.contadores.creados_hoy_regiones}) ·{" "}
            <strong>{data.contadores.pagados_sin_despachar}</strong> por despachar · <strong>{data.contadores.listos_retiro}</strong>{" "}
            listos para retiro
          </p>

          <Bloque
            titulo="Esperando pago hace más de 24 h"
            pedidos={f(data.esperando_pago_mas_24h)}
            total={data.contadores.esperando_pago_mas_24h}
            alerta
            detalle={(p) => `esperando pago ${hace(p.estado_desde)}`}
            accion="Recordar pago o cancelar"
          />
          <Bloque
            titulo="Pagados sin despachar"
            pedidos={f(data.pagados_sin_despachar)
              .slice()
              .sort((a, b) => Date.parse(a.estado_desde) - Date.parse(b.estado_desde))}
            total={data.contadores.pagados_sin_despachar}
            alerta
            detalle={(p) => `${p.estado === "pagado" ? "pagado" : "preparando"} ${hace(p.estado_desde)}`}
            accion="Marcar despachado"
          />
          <Bloque
            titulo="Listos para retiro"
            pedidos={f(data.listos_retiro)}
            total={data.contadores.listos_retiro}
            detalle={(p) => `listo hace ${Math.max(0, Math.floor(horasDesde(p.estado_desde) / 24))} días`}
            accion="Marcar entregado"
          />
          <Bloque
            titulo="Creados hoy"
            pedidos={f(data.creados_hoy)}
            total={data.creados_hoy.length}
            detalle={(p) => hace(p.creado_en)}
          />

          {f(data.esperando_pago_mas_24h).length +
            f(data.pagados_sin_despachar).length +
            f(data.listos_retiro).length +
            f(data.creados_hoy).length ===
          0 ? (
            <p className="rounded-[12px] bg-blanco p-5 text-center text-[16px] font-semibold">Nada pendiente por aquí. ¡Bien!</p>
          ) : null}
        </>
      ) : null}
      <BotonNuevoPedido />
    </div>
  );
}
