// "Hoy" como Server Component (QA I-8): los datos llegan en el HTML; solo
// "Despachar" y "Actualizar" son islas cliente. El filtro de canal es un link (?canal=).
import Link from "next/link";
import type { ReactNode } from "react";
import { getPanel } from "@/lib/api-panel";
import type { Canal, Hoy as DatosHoy, OpcionesEntrega, Pedido } from "@/lib/tipos";
import { Actualizar, Despachar } from "./Despachar";
import { AccionTarjeta, Aviso, BotonNuevoPedido, hace, horasDesde, TarjetaPedido } from "./piezas";

type Filtro = "todos" | Canal;

function Bloque({
  titulo,
  total,
  pedidos,
  alerta,
  detalle,
  acciones,
}: {
  titulo: string;
  total: number;
  pedidos: Pedido[];
  alerta?: boolean;
  detalle: (p: Pedido) => string;
  acciones?: (p: Pedido) => ReactNode;
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
            <TarjetaPedido p={p} alerta={alerta} detalle={detalle(p)} acciones={acciones?.(p)} />
          </li>
        ))}
      </ul>
    </section>
  );
}

function dias(iso: string): string {
  const d = Math.floor(horasDesde(iso) / 24);
  return d < 1 ? "hoy" : `hace ${d} ${d === 1 ? "día" : "días"}`;
}

export async function Hoy({ canal }: { canal: Filtro }) {
  const [r, opc] = await Promise.all([getPanel<DatosHoy>("/api/admin/hoy"), getPanel<OpcionesEntrega>("/api/opciones-entrega")]);
  if (!r.ok) {
    // Sin sesión el layout ya muestra el login.
    return r.status === 401 || r.status === 403 ? null : <Aviso>{r.mensaje}</Aviso>;
  }
  const data = r.data;
  const agencias = opc.ok ? opc.data.agencias : [];
  const f = (l: Pedido[]) => (canal === "todos" ? l : l.filter((p) => p.canal === canal));
  const cotizadosHoy = data.creados_hoy.filter((p) => p.estado === "cotizado");

  const vacio =
    f(data.esperando_pago_mas_24h).length + f(data.pagados_sin_despachar).length + f(data.listos_retiro).length + f(cotizadosHoy).length ===
    0;

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between gap-3">
        <h1 className="titular text-[30px]">Hoy</h1>
        <Actualizar />
      </div>

      <nav aria-label="Canal" className="grid grid-cols-3 gap-1 rounded-[12px] bg-blanco p-1 ring-1 ring-negro/10">
        {(
          [
            ["todos", "Todos"],
            ["santiago", "Santiago"],
            ["regiones", "Regiones"],
          ] as const
        ).map(([v, t]) => (
          <Link
            key={v}
            href={v === "todos" ? "/admin" : `/admin?canal=${v}`}
            aria-current={canal === v ? "page" : undefined}
            className={`flex h-11 items-center justify-center rounded-[10px] text-[15px] font-bold ${canal === v ? "bg-negro text-blanco" : "text-negro"}`}
          >
            {t}
          </Link>
        ))}
      </nav>

      <p className="rounded-[12px] bg-blanco p-3 text-[15px] ring-1 ring-negro/5">
        Hoy: <strong>{data.contadores.creados_hoy_santiago + data.contadores.creados_hoy_regiones}</strong> nuevos (Santiago{" "}
        {data.contadores.creados_hoy_santiago} · Regiones {data.contadores.creados_hoy_regiones}) ·{" "}
        <strong>{data.contadores.pagados_sin_despachar}</strong> por despachar · <strong>{data.contadores.listos_retiro}</strong> listos
        para retiro
      </p>

      <Bloque
        titulo="Esperando pago hace más de 24 h"
        total={data.contadores.esperando_pago_mas_24h}
        pedidos={f(data.esperando_pago_mas_24h)}
        alerta
        detalle={(p) => `esperando pago ${hace(p.estado_desde)}`}
        acciones={(p) => (
          <>
            <AccionTarjeta href={`/admin/pedidos/${p.id}#mensajes`}>Recordar pago</AccionTarjeta>
            <AccionTarjeta href={`/admin/pedidos/${p.id}?accion=cancelado`} peligro>
              Cancelar
            </AccionTarjeta>
          </>
        )}
      />
      <Bloque
        titulo="Pagados sin despachar"
        total={data.contadores.pagados_sin_despachar}
        pedidos={f(data.pagados_sin_despachar)
          .slice()
          .sort((a, b) => Date.parse(a.estado_desde) - Date.parse(b.estado_desde))}
        alerta
        detalle={(p) => `${p.estado === "pagado" ? "pagado" : "preparando"} ${hace(p.estado_desde)}`}
        acciones={(p) =>
          p.entrega === "envio" ? (
            <Despachar id={p.id} numero={p.numero} agencia={p.agencia} comuna={p.comuna} agencias={agencias} />
          ) : (
            <AccionTarjeta href={`/admin/pedidos/${p.id}?accion=listo_retiro`}>Marcar listo para retiro</AccionTarjeta>
          )
        }
      />
      <Bloque
        titulo="Listos para retiro"
        total={data.contadores.listos_retiro}
        pedidos={f(data.listos_retiro)}
        detalle={(p) => `listo ${dias(p.estado_desde)}`}
        acciones={(p) => <AccionTarjeta href={`/admin/pedidos/${p.id}?accion=entregado`}>Marcar entregado</AccionTarjeta>}
      />
      <Bloque titulo="Cotizados hoy" total={cotizadosHoy.length} pedidos={f(cotizadosHoy)} detalle={(p) => hace(p.creado_en)} />

      {vacio ? <p className="rounded-[12px] bg-blanco p-5 text-center text-[16px] font-semibold">Nada pendiente por aquí. ¡Bien!</p> : null}
      {/* Espacio para que "+ Nuevo pedido" no tape la última tarjeta (QA M-8). */}
      <div className="h-16" aria-hidden />
      <BotonNuevoPedido />
    </div>
  );
}
