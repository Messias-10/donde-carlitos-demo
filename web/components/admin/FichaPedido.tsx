"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import { ACCION, ETIQUETA } from "@/lib/estados";
import { clp, fechaHora } from "@/lib/formato";
import type { EstadoPedido, MedioPago, MensajeCliente, OpcionesEntrega, Pedido } from "@/lib/tipos";
import { CopiarTexto } from "../CopiarTexto";
import { IconoFlechaIzq, IconoWhatsApp } from "../Iconos";
import { useManejarError } from "./Sesion";
import { Aviso, BTN_PRIMARIO, BTN_SECUNDARIO, ChipEstado, Dialogo, EtiquetaCanal, INPUT, telefonoLegible } from "./ui";
import { useApi } from "./useApi";

const MEDIOS: [MedioPago, string][] = [
  ["transferencia", "Transferencia"],
  ["efectivo", "Efectivo"],
  ["debito", "Débito"],
  ["credito", "Crédito"],
  ["otro", "Otro"],
];

function hoyChile(): string {
  return new Intl.DateTimeFormat("en-CA", { timeZone: "America/Santiago" }).format(new Date());
}

const entero = (s: string): number | null => (/^\d{1,9}$/.test(s.trim()) ? Number(s.trim()) : null);

type Paso = EstadoPedido | null;

function Mensajes({ id, canal, demo }: { id: number; canal: string; demo: boolean }) {
  const { data, error } = useApi<{ numero: string; telefono_cliente: string; mensajes: MensajeCliente[] }>(
    `/api/admin/pedidos/${id}/mensajes`,
  );
  if (error) return <Aviso>{error}</Aviso>;
  if (!data) return <div className="h-24 animate-pulse rounded-[12px] bg-blanco" aria-hidden />;
  const lista = [...data.mensajes].sort((a, b) => Number(b.sugerido) - Number(a.sugerido));
  return (
    <section id="mensajes" aria-labelledby="mensajes-titulo" className="space-y-3">
      <h2 id="mensajes-titulo" className="titular text-[22px]">
        Mensajes para el cliente
      </h2>
      {demo ? (
        <p className="text-[13px] text-gris">Pedido de prueba (demo). Los mensajes se abren sin número, para que elijas el contacto.</p>
      ) : null}
      <ul className="space-y-3">
        {lista.map((m) => (
          <li key={m.clave} className={`rounded-[12px] bg-blanco p-4 ${m.sugerido ? "ring-2 ring-oro" : "ring-1 ring-negro/10"}`}>
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h3 className="text-[16px] font-extrabold">{m.titulo}</h3>
              {m.sugerido ? <span className="rounded-full bg-oro px-2 py-0.5 text-[12px] font-extrabold text-negro">Sugerido</span> : null}
            </div>
            <p className="mt-1 text-[13px] text-gris">Envíalo desde el WhatsApp de {canal === "regiones" ? "Regiones" : "Santiago"}.</p>
            <details className="mt-2" open={m.sugerido}>
              <summary className="min-h-11 py-2 text-[14px] font-bold underline underline-offset-2">Así le llega:</summary>
              <pre id={`msj-${m.clave}`} className="whitespace-pre-wrap break-words rounded-[10px] bg-panel p-3 font-sans text-[14px] leading-relaxed">
                {m.texto}
              </pre>
            </details>
            <div className="mt-3 flex flex-wrap gap-2">
              <a
                href={m.whatsapp_url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex min-h-12 flex-1 items-center justify-center gap-2 rounded-[12px] bg-whatsapp px-4 text-[15px] font-extrabold text-negro"
              >
                <IconoWhatsApp size={20} /> Abrir en WhatsApp
              </a>
              <CopiarTexto texto={m.texto} etiqueta="Copiar mensaje" idObjetivo={`msj-${m.clave}`} className="border-negro" />
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}

export function FichaPedido({ id }: { id: string }) {
  const { data: p, error: errorCarga, reemplazar, recargar } = useApi<Pedido>(`/api/admin/pedidos/${encodeURIComponent(id)}`);
  const opciones = useApi<OpcionesEntrega>("/api/opciones-entrega");
  const nuevo = useSearchParams().get("nuevo") === "1";
  const manejar = useManejarError();
  const [paso, setPaso] = useState<Paso>(null);
  const [error, setError] = useState("");
  const [aviso, setAviso] = useState("");
  const [guardando, setGuardando] = useState(false);
  const [vueltaMensajes, setVueltaMensajes] = useState(0);
  // Campos de los diálogos.
  const [nota, setNota] = useState("");
  const [medio, setMedio] = useState<MedioPago>("transferencia");
  const [monto, setMonto] = useState("");
  const [fecha, setFecha] = useState(hoyChile());
  const [confirmo, setConfirmo] = useState(false);
  const [agencia, setAgencia] = useState("");
  const [nSeg, setNSeg] = useState("");
  const [comuna, setComuna] = useState("");
  const [motivo, setMotivo] = useState("");
  const [errDialogo, setErrDialogo] = useState("");

  if (!p) {
    return errorCarga ? <Aviso>{errorCarga}</Aviso> : <div className="h-60 animate-pulse rounded-[12px] bg-blanco" aria-hidden />;
  }
  const pedido = p;

  const pideePago = (a: EstadoPedido) => a === "pagado" || (a === "entregado" && pedido.pago === null && pedido.entrega === "retiro");

  function abrir(a: EstadoPedido) {
    setErrDialogo("");
    setNota("");
    setMonto(String(pedido.total));
    setMedio(pedido.forma_pago === "transferencia" ? "transferencia" : "efectivo");
    setFecha(hoyChile());
    setConfirmo(false);
    setAgencia(pedido.agencia ?? "");
    setNSeg(pedido.n_seguimiento ?? "");
    setComuna(pedido.comuna ?? "");
    setMotivo("");
    setPaso(a);
  }

  async function cambiar() {
    if (!paso) return;
    const cuerpo: Record<string, unknown> = { estado: paso };
    if (nota.trim()) cuerpo.nota = nota.trim();
    if (pideePago(paso)) {
      const m = entero(monto);
      if (!m) {
        setErrDialogo("Escribe el monto recibido (solo números).");
        return;
      }
      cuerpo.pago = { medio, monto: m, fecha };
    }
    if (paso === "pagado") {
      if (!confirmo) {
        setErrDialogo("Marca que revisaste que la plata ya está en la cuenta.");
        return;
      }
      cuerpo.confirmo_pago_recibido = true;
    }
    if (paso === "despachado") {
      if (!agencia || !nSeg.trim() || !comuna.trim()) {
        setErrDialogo("Para despachar necesitas la agencia, el n.º de seguimiento y la comuna.");
        return;
      }
      cuerpo.agencia = agencia;
      cuerpo.n_seguimiento = nSeg.trim();
      cuerpo.comuna = comuna.trim();
    }
    if (paso === "cancelado" && motivo.trim()) cuerpo.motivo = motivo.trim();

    setGuardando(true);
    setErrDialogo("");
    try {
      const act = await api<Pedido>(`/api/admin/pedidos/${pedido.id}/estado`, { method: "PATCH", body: cuerpo });
      reemplazar(act);
      setVueltaMensajes((v) => v + 1);
      setPaso(null);
      setError("");
      setAviso(`Pedido marcado como "${ETIQUETA[act.estado]}". Mira abajo el mensaje sugerido.`);
    } catch (err) {
      if (err instanceof ErrorCliente && err.codigo === "stock_insuficiente") {
        setErrDialogo("No hay stock suficiente de algún producto del catálogo. Revisa el pedido o el stock.");
      } else if (err instanceof ErrorCliente && err.campos) {
        setErrDialogo(Object.values(err.campos).join(" · "));
      } else {
        setErrDialogo(manejar(err));
      }
    } finally {
      setGuardando(false);
    }
  }

  const acciones = pedido.transiciones_permitidas.filter((e) => e !== "cancelado");
  const puedeCancelar = pedido.transiciones_permitidas.includes("cancelado");
  const editable = pedido.estado === "cotizado" || pedido.estado === "esperando_pago";

  return (
    <div className="space-y-4">
      <Link href="/admin" className="inline-flex min-h-11 items-center gap-1 text-[14px] font-bold">
        <IconoFlechaIzq size={18} /> Hoy
      </Link>

      {nuevo ? <Aviso tipo="ok">Pedido {pedido.numero} guardado. Envíale la cotización por WhatsApp (abajo).</Aviso> : null}
      {aviso ? <Aviso tipo="ok">{aviso}</Aviso> : null}
      {error ? <Aviso>{error}</Aviso> : null}

      <section className="rounded-[12px] bg-blanco p-4">
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="font-display text-[32px] leading-none">{pedido.numero}</h1>
          <ChipEstado estado={pedido.estado} grande />
          <EtiquetaCanal canal={pedido.canal} />
        </div>
        {pedido.es_demo ? <p className="mt-1 text-[12px] font-semibold text-gris">Pedido de prueba (demo)</p> : null}
        <p className="mt-1 text-[13px] text-gris">
          Registrado por {pedido.registrado_por?.nombre ?? "—"} · {fechaHora(pedido.creado_en)}
        </p>
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <span className="text-[14px] font-bold">Link de seguimiento</span>
          <CopiarTexto texto={pedido.seguimiento_url} etiqueta="Copiar link" className="border-negro" />
        </div>

        <dl className="mt-4 grid grid-cols-[auto_1fr] gap-x-4 gap-y-1 text-[15px]">
          <dt className="text-gris">Cliente</dt>
          <dd className="font-semibold">{pedido.nombre || "Sin nombre"}</dd>
          <dt className="text-gris">Teléfono</dt>
          <dd className="font-semibold tabular-nums">{telefonoLegible(pedido.telefono)}</dd>
          <dt className="text-gris">Entrega</dt>
          <dd className="font-semibold">
            {pedido.entrega === "retiro"
              ? "Retiro en tienda"
              : `Despacho${pedido.comuna ? ` a ${pedido.comuna}` : ""}${pedido.region_nombre ? `, ${pedido.region_nombre}` : ""}`}
          </dd>
          {pedido.entrega === "envio" ? (
            <>
              <dt className="text-gris">Agencia</dt>
              <dd className="font-semibold">{pedido.agencia_nombre ?? "Por definir"}</dd>
              {pedido.n_seguimiento ? (
                <>
                  <dt className="text-gris">N.º seguimiento</dt>
                  <dd className="font-semibold tabular-nums">{pedido.n_seguimiento}</dd>
                </>
              ) : null}
            </>
          ) : null}
          <dt className="text-gris">Pago</dt>
          <dd className="font-semibold">
            {pedido.forma_pago === "transferencia" ? "Transferencia" : "En tienda"}
            {pedido.pago ? ` · recibido ${clp(pedido.pago.monto)} (${pedido.pago.medio})${pedido.pago.verificado_por ? ` · revisó ${pedido.pago.verificado_por.nombre}` : ""}` : ""}
          </dd>
        </dl>
        {pedido.nota_interna ? (
          <p className="mt-3 whitespace-pre-line rounded-[10px] bg-panel p-3 text-[14px]">
            <span className="font-bold">Nota interna: </span>
            {pedido.nota_interna}
          </p>
        ) : null}
        {pedido.motivo_cancelacion ? (
          <p className="mt-3 rounded-[10px] bg-st-cancel-bg p-3 text-[14px] text-st-cancel-fg">Motivo: {pedido.motivo_cancelacion}</p>
        ) : null}
      </section>

      <section className="rounded-[12px] bg-blanco p-4">
        <h2 className="text-[17px] font-extrabold">Productos</h2>
        <ul className="mt-2 divide-y divide-negro/10">
          {pedido.items.map((it, i) => (
            <li key={i} className="flex justify-between gap-3 py-2 text-[15px]">
              <span className="min-w-0">
                <span className="font-semibold">
                  {it.cantidad} x {it.descripcion}
                </span>
                <span className="block text-[13px] text-gris">
                  {it.codigo ? `Cód. ${it.codigo} · ` : "Línea libre · "}
                  {it.precio_unitario === null ? "sin precio" : `${clp(it.precio_unitario)} c/u`}
                </span>
              </span>
              <span className="shrink-0 font-semibold tabular-nums">{it.subtotal === null ? "—" : clp(it.subtotal)}</span>
            </li>
          ))}
        </ul>
        <dl className="mt-2 space-y-1 border-t-2 border-negro pt-2 text-[15px]">
          {pedido.descuento > 0 ? (
            <div className="flex justify-between">
              <dt>Descuento</dt>
              <dd>-{clp(pedido.descuento)}</dd>
            </div>
          ) : null}
          {pedido.costo_envio > 0 ? (
            <div className="flex justify-between">
              <dt>Envío</dt>
              <dd>{clp(pedido.costo_envio)}</dd>
            </div>
          ) : null}
          {pedido.envio_por_pagar ? <p className="text-[13px] text-gris">Envío por pagar en la agencia.</p> : null}
          <div className="flex justify-between text-[18px] font-extrabold">
            <dt>Total{pedido.total_manual ? " (escrito a mano)" : ""}</dt>
            <dd>{clp(pedido.total)}</dd>
          </div>
        </dl>
        {editable ? (
          <Link href={`/admin/pedidos/${pedido.id}/editar`} className="mt-3 inline-flex min-h-11 items-center text-[15px] font-bold underline underline-offset-2">
            Editar pedido
          </Link>
        ) : null}
      </section>

      {acciones.length > 0 || puedeCancelar ? (
        <section className="space-y-2 rounded-[12px] bg-blanco p-4">
          <h2 className="text-[17px] font-extrabold">Siguiente paso</h2>
          {acciones.map((a, i) => (
            <button key={a} type="button" onClick={() => abrir(a)} className={i === 0 ? BTN_PRIMARIO : BTN_SECUNDARIO}>
              {ACCION[a]}
            </button>
          ))}
          {puedeCancelar ? (
            <button type="button" onClick={() => abrir("cancelado")} className="min-h-11 w-full text-[15px] font-bold text-rojo underline underline-offset-2">
              Cancelar pedido
            </button>
          ) : null}
        </section>
      ) : null}

      <Mensajes key={`${pedido.estado}-${vueltaMensajes}`} id={pedido.id} canal={pedido.canal} demo={pedido.es_demo} />

      {pedido.historial && pedido.historial.length > 0 ? (
        <section className="rounded-[12px] bg-blanco p-4">
          <h2 className="text-[17px] font-extrabold">Historial</h2>
          <ul className="mt-2 space-y-1 text-[14px]">
            {pedido.historial.map((h, i) => (
              <li key={i}>
                {fechaHora(h.creado_en)} · {h.estado_anterior ? ETIQUETA[h.estado_anterior] : "—"} → {ETIQUETA[h.estado_nuevo]} ·{" "}
                {h.usuario?.nombre ?? "sistema"}
                {h.nota ? <span className="block pl-3 text-gris">{h.nota}</span> : null}
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      <button type="button" onClick={recargar} className="min-h-11 text-[14px] font-bold underline underline-offset-2">
        Recargar pedido
      </button>

      <Dialogo
        abierto={paso !== null}
        onCerrar={() => setPaso(null)}
        titulo={
          paso === "pagado"
            ? "¿Revisaste que la plata ya está en la cuenta?"
            : paso === "cancelado"
              ? `¿Cancelar el pedido ${pedido.numero}?`
              : paso
                ? `${ACCION[paso]} · ${pedido.numero}`
                : ""
        }
      >
        {paso === "pagado" ? <p className="text-[15px]">No te fíes solo del comprobante.</p> : null}

        {paso && pideePago(paso) ? (
          <div className="space-y-3">
            <div>
              <label htmlFor="d-medio" className="text-[14px] font-bold">
                Medio de pago
              </label>
              <select id="d-medio" value={medio} onChange={(e) => setMedio(e.target.value as MedioPago)} className={INPUT}>
                {MEDIOS.map(([v, t]) => (
                  <option key={v} value={v}>
                    {t}
                  </option>
                ))}
              </select>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label htmlFor="d-monto" className="text-[14px] font-bold">
                  Monto recibido
                </label>
                <input id="d-monto" inputMode="numeric" value={monto} onChange={(e) => setMonto(e.target.value)} className={INPUT} />
              </div>
              <div>
                <label htmlFor="d-fecha" className="text-[14px] font-bold">
                  Fecha
                </label>
                <input id="d-fecha" type="date" value={fecha} onChange={(e) => setFecha(e.target.value)} className={INPUT} />
              </div>
            </div>
          </div>
        ) : null}

        {paso === "pagado" ? (
          <label className="flex min-h-12 items-start gap-3 rounded-[10px] bg-st-espera-bg p-3 text-[15px] font-semibold text-st-espera-fg">
            <input type="checkbox" checked={confirmo} onChange={(e) => setConfirmo(e.target.checked)} className="mt-0.5 size-6 shrink-0 accent-negro" />
            Sí, revisé que la plata ya está en la cuenta.
          </label>
        ) : null}

        {paso === "despachado" ? (
          <div className="space-y-3">
            <div>
              <label htmlFor="d-agencia" className="text-[14px] font-bold">
                Agencia
              </label>
              <select id="d-agencia" value={agencia} onChange={(e) => setAgencia(e.target.value)} className={INPUT}>
                <option value="">Elige la agencia</option>
                {opciones.data?.agencias.map((a) => (
                  <option key={a.codigo} value={a.codigo}>
                    {a.nombre}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label htmlFor="d-nseg" className="text-[14px] font-bold">
                N.º de seguimiento
              </label>
              <input id="d-nseg" value={nSeg} onChange={(e) => setNSeg(e.target.value)} maxLength={60} className={INPUT} />
            </div>
            <div>
              <label htmlFor="d-comuna" className="text-[14px] font-bold">
                Comuna
              </label>
              <input id="d-comuna" value={comuna} onChange={(e) => setComuna(e.target.value)} maxLength={60} className={INPUT} />
            </div>
          </div>
        ) : null}

        {paso === "cancelado" ? (
          <div>
            <label htmlFor="d-motivo" className="text-[14px] font-bold">
              Motivo (opcional)
            </label>
            <input
              id="d-motivo"
              value={motivo}
              onChange={(e) => setMotivo(e.target.value)}
              maxLength={200}
              placeholder="Ej: se agotó, el cliente desistió"
              className={INPUT}
            />
          </div>
        ) : (
          <div>
            <label htmlFor="d-nota" className="text-[14px] font-bold">
              Nota (opcional)
            </label>
            <input id="d-nota" value={nota} onChange={(e) => setNota(e.target.value)} maxLength={200} className={INPUT} />
          </div>
        )}

        {errDialogo ? <Aviso>{errDialogo}</Aviso> : null}

        <div className="space-y-2">
          <button type="button" onClick={cambiar} disabled={guardando} className={BTN_PRIMARIO}>
            {guardando
              ? "Guardando…"
              : paso === "pagado"
                ? "Sí, está pagado"
                : paso === "despachado"
                  ? "Guardar y marcar despachado"
                  : paso === "cancelado"
                    ? "Sí, cancelar"
                    : paso
                      ? ACCION[paso]
                      : ""}
          </button>
          <button type="button" onClick={() => setPaso(null)} className={BTN_SECUNDARIO}>
            {paso === "pagado" ? "Todavía no" : paso === "cancelado" ? "No" : "Volver"}
          </button>
        </div>
      </Dialogo>
    </div>
  );
}
