"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useRef, useState } from "react";
import { cotizar, MAX_CANTIDAD, quitarVarios, useCarrito, vaciar, type LineaCarrito } from "@/lib/carrito";
import { clp } from "@/lib/formato";
import type { ErrorApi, OpcionesEntrega, PedidoCreado } from "@/lib/tipos";
import { Ejemplo } from "../Marcas";
import { PasosMini, useHidratado } from "./Carrito";
import { guardarConfirmacion } from "./confirmacionGuardada";

type Campo = "nombre" | "telefono" | "entrega" | "region" | "comuna" | "agencia" | "comentario";
type Errores = Partial<Record<Campo, string>>;

const RE_NOMBRE = /^[\p{L}\p{M} '’.-]+$/u;
const RE_LUGAR = /^[\p{L}\p{M}\d '’.\-(),]+$/u;
const TEL = `+56 9 2114 5368`;

/** Valida igual que la API (pedidos.Solicitud.Validar), con los textos de copy-landing.md §15.5. */
function validar(d: Datos): Errores {
  const e: Errores = {};
  const nombre = d.nombre.trim().replace(/\s+/g, " ");
  if (nombre.length < 2) e.nombre = "Escribe tu nombre.";
  else if (!RE_NOMBRE.test(nombre)) e.nombre = "Usa solo letras en tu nombre.";
  else if (nombre.length > 80) e.nombre = "Tu nombre puede tener hasta 80 letras.";
  if (!telefonoValido(d.telefono)) e.telefono = "Revisa tu número: tiene que ser un celular chileno de 9 dígitos.";
  if (d.entrega !== "retiro" && d.entrega !== "envio") e.entrega = "Elige si retiras o te lo enviamos.";
  if (d.entrega === "envio") {
    const comuna = d.comuna.trim();
    if (!d.region || !d.agencia || comuna.length < 2 || comuna.length > 60 || !RE_LUGAR.test(comuna)) {
      const msg = "Elige tu región, tu comuna y la agencia.";
      if (!d.region) e.region = msg;
      if (comuna.length < 2 || comuna.length > 60 || !RE_LUGAR.test(comuna)) e.comuna = msg;
      if (!d.agencia) e.agencia = msg;
    }
  }
  if (d.comentario.length > 500) e.comentario = "Tu comentario puede tener hasta 500 caracteres.";
  return e;
}

function telefonoValido(t: string): boolean {
  const s = t.trim();
  if (!/^\+?[\d\s().-]{8,20}$/.test(s)) return false;
  let dig = s.replace(/\D/g, "");
  if (dig.length === 11 && dig.startsWith("56")) dig = dig.slice(2);
  return dig.length === 9 && dig.startsWith("9");
}

interface Datos {
  nombre: string;
  telefono: string;
  entrega: "" | "retiro" | "envio";
  region: string;
  comuna: string;
  agencia: string;
  comentario: string;
}

const INICIAL: Datos = { nombre: "", telefono: "", entrega: "retiro", region: "", comuna: "", agencia: "", comentario: "" };

function mapearCampoApi(k: string): Campo | null {
  return (["nombre", "telefono", "entrega", "region", "comuna", "agencia", "comentario"] as const).find((c) => c === k) ?? null;
}

export function FormPedido({ opciones }: { opciones: OpcionesEntrega | null }) {
  const items = useCarrito();
  const hidratado = useHidratado();
  const [d, setD] = useState<Datos>(INICIAL);
  const [errores, setErrores] = useState<Errores>({});
  const [errorGeneral, setErrorGeneral] = useState<string[]>([]);
  const [enviando, setEnviando] = useState(false);
  const router = useRouter();
  const [saliendo, setSaliendo] = useState(false);
  const resumenError = useRef<HTMLDivElement>(null);

  if (saliendo) return <p className="rounded-[16px] bg-blanco p-6 text-center text-[16px] font-bold" role="status">Pedido guardado. Abriendo tu confirmación…</p>;

  if (!hidratado) return <div className="h-96 animate-pulse rounded-[12px] bg-blanco/60" aria-hidden />;

  const { lineas, total, unidades } = cotizar(items);

  if (lineas.length === 0) {
    return (
      <div className="rounded-[16px] bg-blanco p-6 text-center">
        <h1 className="titular mb-3 text-[30px]">Tus datos</h1>
        <p className="text-[17px] font-bold">Tu pedido está vacío. ¡Mira lo que llegó!</p>
        <Link href="/#ofertas" className="mt-4 inline-flex h-12 items-center rounded-[12px] bg-negro px-5 font-extrabold text-oro">
          Ver ofertas
        </Link>
      </div>
    );
  }

  const set = <K extends keyof Datos>(k: K, v: Datos[K]) => {
    setD((x) => ({ ...x, [k]: v }));
    if (errores[k as Campo]) setErrores((e) => ({ ...e, [k]: undefined }));
  };

  function mostrarErrores(e: Errores, generales: string[] = []) {
    setErrores(e);
    setErrorGeneral(generales);
    requestAnimationFrame(() => {
      const primero = Object.keys(e)[0];
      if (primero) document.getElementById(`f-${primero}`)?.focus();
      else resumenError.current?.focus();
    });
  }

  async function enviar(ev: React.FormEvent) {
    ev.preventDefault();
    if (enviando) return;
    const e = validar(d);
    if (Object.keys(e).length > 0) {
      mostrarErrores(e);
      return;
    }
    const cuerpo: Record<string, unknown> = {
      nombre: d.nombre.trim().replace(/\s+/g, " "),
      telefono: d.telefono.trim(),
      entrega: d.entrega,
      items: items.map((l) => ({ producto_id: l.id, cantidad: l.cantidad })),
    };
    if (d.entrega === "envio") {
      cuerpo.region = d.region;
      cuerpo.comuna = d.comuna.trim();
      cuerpo.agencia = d.agencia;
    }
    if (d.comentario.trim()) cuerpo.comentario = d.comentario.trim();

    setEnviando(true);
    setErrorGeneral([]);
    try {
      const res = await fetch("/api/pedidos", {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify(cuerpo),
        credentials: "same-origin",
      });
      if (res.status === 201) {
        const pedido = (await res.json()) as PedidoCreado;
        guardarConfirmacion(pedido);
        setSaliendo(true);
        vaciar();
        // La confirmación vive en /pedido/listo para sobrevivir a recargas y a la vuelta desde WhatsApp (QA I-6).
        router.replace(`/pedido/listo?n=${encodeURIComponent(pedido.numero)}`);
        return;
      }
      const err = ((await res.json().catch(() => null)) as ErrorApi | null)?.error;
      if (res.status === 429) {
        mostrarErrores({}, [`Hiciste varios pedidos seguidos. Espera unos minutos o escríbenos al ${TEL}.`]);
      } else if (res.status === 409 && err?.codigo === "productos_no_disponibles") {
        mostrarErrores({}, mensajesStock(err.detalles, items));
      } else if (res.status === 409 && err?.codigo === "pedidos_abiertos") {
        mostrarErrores({}, [
          `Ya tienes pedidos esperando confirmación con este número. Escríbenos al ${TEL} para confirmarlos antes de hacer otro.`,
        ]);
      } else if (res.status === 400 && err?.codigo === "validacion" && err.campos) {
        const e2: Errores = {};
        const otros: string[] = [];
        for (const [k, v] of Object.entries(err.campos)) {
          const c = mapearCampoApi(k);
          if (c) e2[c] = v.charAt(0).toUpperCase() + v.slice(1) + ".";
          else otros.push(mensajeItem(k, v, items));
        }
        mostrarErrores(e2, otros);
      } else {
        mostrarErrores({}, [`No pudimos guardar tu pedido. Inténtalo de nuevo en un rato o escríbenos al ${TEL}.`]);
      }
    } catch {
      mostrarErrores({}, [
        navigator.onLine === false
          ? "Parece que no tienes internet. Tu carrito sigue guardado."
          : `No pudimos guardar tu pedido. Inténtalo de nuevo en un rato o escríbenos al ${TEL}.`,
      ]);
    } finally {
      setEnviando(false);
    }
  }

  const err = (c: Campo) =>
    errores[c] ? (
      <p id={`e-${c}`} className="mt-1 text-[14px] font-semibold text-rojo">
        {errores[c]}
      </p>
    ) : null;
  const aria = (c: Campo) => ({
    "aria-invalid": errores[c] ? true : undefined,
    "aria-describedby": errores[c] ? `e-${c}` : undefined,
  });
  const input =
    "mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px] focus:border-negro aria-[invalid=true]:border-rojo";

  return (
    <>
    <h1 className="titular mb-5 text-[36px]">Tus datos</h1>
    <div className="grid gap-6 lg:grid-cols-[1fr_360px] lg:items-start">
      <form onSubmit={enviar} noValidate className="space-y-5 rounded-[16px] bg-blanco p-5">
        {errorGeneral.length > 0 ? (
          <div ref={resumenError} tabIndex={-1} role="alert" className="space-y-1 rounded-[10px] border-2 border-rojo bg-rojo/5 p-3 text-[15px] font-semibold text-rojo">
            {errorGeneral.map((m) => (
              <p key={m}>{m}</p>
            ))}
          </div>
        ) : null}

        <div>
          <label htmlFor="f-nombre" className="text-[15px] font-bold">
            Tu nombre
          </label>
          <input
            id="f-nombre"
            className={input}
            value={d.nombre}
            onChange={(e) => set("nombre", e.target.value)}
            placeholder="Ej: Camila Rojas"
            autoComplete="name"
            maxLength={80}
            {...aria("nombre")}
          />
          {err("nombre")}
        </div>

        <div>
          <label htmlFor="f-telefono" className="text-[15px] font-bold">
            Tu WhatsApp
          </label>
          <input
            id="f-telefono"
            className={input}
            value={d.telefono}
            onChange={(e) => set("telefono", e.target.value)}
            placeholder="9 1234 5678"
            type="tel"
            inputMode="tel"
            autoComplete="tel-national"
            maxLength={20}
            {...aria("telefono")}
            aria-describedby={errores.telefono ? "e-telefono" : "a-telefono"}
          />
          <p id="a-telefono" className="mt-1 text-[13px] text-gris">
            Te confirmamos el pedido a este número.
          </p>
          {err("telefono")}
        </div>

        <fieldset>
          <legend className="text-[15px] font-bold">¿Cómo lo quieres?</legend>
          <div className="mt-2 grid gap-2">
            {(
              [
                ["retiro", "Retiro en tienda · Alameda 2432"],
                ["envio", "Envío por agencia (regiones)"],
              ] as const
            ).map(([v, t]) => (
              <label
                key={v}
                className={`flex min-h-12 cursor-pointer items-center gap-3 rounded-[10px] border-2 px-3 text-[15px] font-semibold ${
                  d.entrega === v ? "border-negro bg-crema" : "border-negro/15"
                }`}
              >
                <input
                  type="radio"
                  name="entrega"
                  value={v}
                  checked={d.entrega === v}
                  onChange={() => set("entrega", v)}
                  className="size-5 accent-negro"
                  id={v === "retiro" ? "f-entrega" : undefined}
                />
                {t}
              </label>
            ))}
          </div>
          {d.entrega === "envio" ? (
            <p className="mt-2 text-[13px] text-gris">
              Envíos a regiones: agencias, quién paga y monto mínimo están por confirmar con la tienda.
            </p>
          ) : null}
          {err("entrega")}
        </fieldset>

        {d.entrega === "envio" ? (
          <div className="space-y-4 rounded-[12px] bg-crema p-4">
            <div>
              <label htmlFor="f-region" className="text-[15px] font-bold">
                Región
              </label>
              <select id="f-region" className={input} value={d.region} onChange={(e) => set("region", e.target.value)} {...aria("region")}>
                <option value="">Elige tu región</option>
                {opciones?.regiones.map((r) => (
                  <option key={r.codigo} value={r.codigo}>
                    {r.nombre}
                  </option>
                ))}
              </select>
              {err("region")}
            </div>
            <div>
              <label htmlFor="f-comuna" className="text-[15px] font-bold">
                Comuna
              </label>
              <input
                id="f-comuna"
                className={input}
                value={d.comuna}
                onChange={(e) => set("comuna", e.target.value)}
                placeholder="Ej: Chillán"
                autoComplete="address-level2"
                maxLength={60}
                {...aria("comuna")}
              />
              {err("comuna")}
            </div>
            <div>
              <label htmlFor="f-agencia" className="text-[15px] font-bold">
                Agencia {opciones?.agencias_es_ejemplo ? <Ejemplo texto="lista de ejemplo" /> : null}
              </label>
              <select
                id="f-agencia"
                className={input}
                value={d.agencia}
                onChange={(e) => set("agencia", e.target.value)}
                {...aria("agencia")}
                aria-describedby={errores.agencia ? "e-agencia" : "a-agencia"}
              >
                <option value="">Elige tu agencia</option>
                {opciones?.agencias.map((a) => (
                  <option key={a.codigo} value={a.codigo}>
                    {a.nombre}
                  </option>
                ))}
              </select>
              <p id="a-agencia" className="mt-1 text-[13px] text-gris">
                El costo de envío te lo confirmamos por WhatsApp.
              </p>
              {err("agencia")}
            </div>
          </div>
        ) : null}

        <div>
          <label htmlFor="f-comentario" className="text-[15px] font-bold">
            ¿Algo más? (opcional)
          </label>
          <textarea
            id="f-comentario"
            className={`${input} h-28 py-2`}
            value={d.comentario}
            onChange={(e) => set("comentario", e.target.value)}
            placeholder="Ej: talla 10, paso el sábado en la mañana, es para mi negocio…"
            maxLength={500}
            {...aria("comentario")}
          />
          <p className="mt-1 text-right text-[12px] text-gris">{d.comentario.length}/500</p>
          {err("comentario")}
        </div>

        <p className="text-[13px] text-gris">Usamos tus datos solo para coordinar este pedido.</p>

        <button
          type="submit"
          disabled={enviando}
          className="flex h-[52px] w-full items-center justify-center rounded-[14px] bg-negro text-[17px] font-extrabold text-oro disabled:opacity-60"
        >
          {enviando ? "Guardando tu pedido…" : "Hacer pedido"}
        </button>
      </form>

      <aside className="space-y-3 rounded-[16px] bg-negro p-5 text-blanco lg:sticky lg:top-24">
        <h2 className="titular text-[22px] text-oro">Resumen</h2>
        <ul className="space-y-2 text-[14px]">
          {lineas.map((l) => (
            <li key={l.id} className="flex justify-between gap-3">
              <span className="min-w-0">
                {l.cantidad} x {l.nombre}
              </span>
              <span className="shrink-0 tabular-nums">{clp(l.subtotal)}</span>
            </li>
          ))}
        </ul>
        <div className="flex items-end justify-between border-t border-blanco/15 pt-3">
          <p className="text-[14px]">
            {unidades} {unidades === 1 ? "unidad" : "unidades"}
            <span className="block text-[12px] text-gris-claro">Total estimado, sin envío</span>
          </p>
          <p className="placa-precio text-[28px]">{clp(total)}</p>
        </div>
        <p className="text-[12px] text-gris-claro">El total final lo calcula la tienda al guardar el pedido.</p>
        <Link href="/carrito" className="block text-[14px] font-bold text-oro underline underline-offset-2">
          Cambiar productos
        </Link>
        <PasosMini />
      </aside>
    </div>
    </>
  );
}

/** Errores de "items[i].cantidad" y similares, con el nombre del producto. */
function mensajeItem(campo: string, motivo: string, items: LineaCarrito[]): string {
  const m = /^items\[(\d+)\]\.(\w+)$/.exec(campo);
  const l = m ? items[Number(m[1])] : undefined;
  if (m && m[2] === "cantidad") {
    return `Puedes pedir hasta ${MAX_CANTIDAD} unidades por producto${l ? ` (${l.nombre})` : ""}. Si necesitas más, escríbelo en el comentario y te cotizamos.`;
  }
  if (campo === "items") return `Revisa tu pedido: ${motivo}.`;
  return `Revisa los productos de tu pedido${l ? ` (${l.nombre})` : ""}: ${motivo}.`;
}

function mensajesStock(detalles: unknown, items: LineaCarrito[]): string[] {
  if (!Array.isArray(detalles)) {
    return ["Algunos productos ya no están disponibles en esa cantidad. Revisa tu pedido."];
  }
  const out: string[] = [];
  const sacar: number[] = [];
  for (const det of detalles as { producto_id?: number; motivo?: string }[]) {
    const l = items.find((x) => x.id === det.producto_id);
    const nombre = l?.nombre ?? "Un producto";
    if (det.motivo === "stock_insuficiente") {
      out.push(`No nos queda stock suficiente de ${nombre}. Baja la cantidad o sácalo del pedido.`);
    } else {
      out.push(`${nombre} ya no está disponible. Lo sacamos de tu pedido.`);
      if (l) sacar.push(l.id);
    }
  }
  if (sacar.length) quitarVarios(sacar);
  return out.length ? out : ["Algunos productos ya no están disponibles en esa cantidad. Revisa tu pedido."];
}

