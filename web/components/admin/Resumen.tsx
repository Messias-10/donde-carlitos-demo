"use client";

import { useState } from "react";
import { ETIQUETA } from "@/lib/estados";
import { clp } from "@/lib/formato";
import type { EstadoPedido, FilaResumen, Resumen as DatosResumen } from "@/lib/tipos";
import { SoloAdmin } from "./Sesion";
import { Aviso, INPUT } from "./ui";
import { useApi } from "./useApi";

function dia(offset = 0): string {
  const d = new Date(Date.now() + offset * 86_400_000);
  return new Intl.DateTimeFormat("en-CA", { timeZone: "America/Santiago" }).format(d);
}

function lunesDeEstaSemana(): string {
  const hoy = dia();
  const d = new Date(`${hoy}T12:00:00Z`);
  const dow = (d.getUTCDay() + 6) % 7; // lunes = 0
  d.setUTCDate(d.getUTCDate() - dow);
  return d.toISOString().slice(0, 10);
}

function Tabla({ titulo, filas, nombre }: { titulo: string; filas: FilaResumen[]; nombre: (f: FilaResumen) => string }) {
  if (filas.length === 0) return null;
  return (
    <section className="rounded-[12px] bg-blanco p-4">
      <h2 className="text-[17px] font-extrabold">{titulo}</h2>
      <table className="mt-2 w-full text-[15px]">
        <thead>
          <tr className="text-left text-[13px] text-gris">
            <th className="py-1 font-semibold">Nombre</th>
            <th className="py-1 text-right font-semibold">Pedidos</th>
            <th className="py-1 text-right font-semibold">Monto</th>
          </tr>
        </thead>
        <tbody>
          {filas.map((f, i) => (
            <tr key={`${f.clave}-${i}`} className="border-t border-negro/10">
              <td className="py-2">{nombre(f)}</td>
              <td className="py-2 text-right tabular-nums">{f.pedidos}</td>
              <td className="py-2 text-right tabular-nums">{clp(f.monto)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

function Contenido() {
  const [rango, setRango] = useState({ desde: dia(), hasta: dia() });
  const [borrador, setBorrador] = useState(rango);
  const q = new URLSearchParams(rango).toString();
  const { data, error } = useApi<DatosResumen>(`/api/admin/resumen?${q}`);
  const [avisoCsv, setAvisoCsv] = useState(false);

  const presets: [string, { desde: string; hasta: string }][] = [
    ["Hoy", { desde: dia(), hasta: dia() }],
    ["Esta semana", { desde: lunesDeEstaSemana(), hasta: dia() }],
    ["Últimos 30 días", { desde: dia(-29), hasta: dia() }],
  ];
  const ticket = data && data.totales.pedidos > 0 ? Math.round(data.totales.monto / data.totales.pedidos) : 0;

  return (
    <div className="space-y-4">
      <h1 className="titular text-[28px]">Resumen</h1>
      <div className="flex flex-wrap gap-2">
        {presets.map(([t, r]) => (
          <button
            key={t}
            type="button"
            aria-pressed={rango.desde === r.desde && rango.hasta === r.hasta}
            onClick={() => {
              setRango(r);
              setBorrador(r);
            }}
            className={`min-h-11 rounded-full px-4 text-[15px] font-bold ${
              rango.desde === r.desde && rango.hasta === r.hasta ? "bg-negro text-blanco" : "bg-blanco ring-1 ring-negro/15"
            }`}
          >
            {t}
          </button>
        ))}
      </div>
      <details className="rounded-[12px] bg-blanco p-4">
        <summary className="min-h-11 py-2 text-[15px] font-bold underline underline-offset-2">Elegir fechas</summary>
        <form
          className="mt-2 grid grid-cols-2 gap-3"
          onSubmit={(e) => {
            e.preventDefault();
            setRango(borrador);
          }}
        >
          <div>
            <label htmlFor="r-desde" className="text-[14px] font-bold">
              Desde
            </label>
            <input id="r-desde" type="date" value={borrador.desde} onChange={(e) => setBorrador((b) => ({ ...b, desde: e.target.value }))} className={INPUT} />
          </div>
          <div>
            <label htmlFor="r-hasta" className="text-[14px] font-bold">
              Hasta
            </label>
            <input id="r-hasta" type="date" value={borrador.hasta} onChange={(e) => setBorrador((b) => ({ ...b, hasta: e.target.value }))} className={INPUT} />
          </div>
          <button type="submit" className="col-span-2 min-h-12 rounded-[12px] bg-negro font-extrabold text-blanco">
            Ver
          </button>
        </form>
      </details>

      {error ? <Aviso>{error}</Aviso> : null}
      {data ? (
        <>
          <div className="grid grid-cols-2 gap-3">
            {[
              ["Pedidos", String(data.totales.pedidos)],
              ["Vendido (sin cancelados)", clp(data.totales.monto)],
              ["Ticket promedio", clp(ticket)],
              ["Cancelados", String(data.totales.cancelados)],
            ].map(([t, v]) => (
              <div key={t} className="rounded-[12px] bg-blanco p-4">
                <p className="text-[13px] font-semibold text-gris">{t}</p>
                <p className="mt-1 font-display text-[28px] leading-none">{v}</p>
              </div>
            ))}
          </div>
          <Tabla titulo="Por canal" filas={data.por_canal} nombre={(f) => (f.clave === "regiones" ? "Regiones" : f.clave === "santiago" ? "Santiago" : (f.clave ?? "—"))} />
          <Tabla titulo="Por estado" filas={data.por_estado} nombre={(f) => (f.clave ? (ETIQUETA[f.clave as EstadoPedido] ?? f.clave) : "—")} />
          <Tabla titulo="Agencias más usadas" filas={data.por_agencia} nombre={(f) => f.nombre ?? "Sin agencia asignada"} />
        </>
      ) : null}

      <section className="space-y-2 rounded-[12px] bg-blanco p-4">
        <h2 className="text-[17px] font-extrabold">Exportar</h2>
        {avisoCsv ? (
          <>
            <Aviso tipo="info">Este archivo trae nombres y teléfonos de clientes. Guárdalo bien y no lo compartas.</Aviso>
            <a
              href={`/api/admin/pedidos.csv?${q}`}
              download="pedidos.csv"
              className="flex min-h-12 items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-blanco"
            >
              Descargar CSV ({rango.desde} a {rango.hasta})
            </a>
          </>
        ) : (
          <button
            type="button"
            onClick={() => setAvisoCsv(true)}
            className="flex min-h-12 w-full items-center justify-center rounded-[12px] border-2 border-negro text-[16px] font-bold"
          >
            Exportar a Excel (CSV)
          </button>
        )}
      </section>
    </div>
  );
}

export function Resumen() {
  return (
    <SoloAdmin>
      <Contenido />
    </SoloAdmin>
  );
}
