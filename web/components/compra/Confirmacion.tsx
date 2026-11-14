"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { clp } from "@/lib/formato";
import { HORARIO, NEGOCIO } from "@/lib/negocio";
import type { PedidoCreado } from "@/lib/tipos";
import { linkWhatsApp } from "@/lib/whatsapp";
import { IconoCopiar, IconoEscudo, IconoWhatsApp } from "../Iconos";
import { PorConfirmar } from "../Marcas";

export type DatosConfirmacion = Pick<PedidoCreado, "numero" | "total" | "entrega" | "es_demo" | "mensaje_whatsapp" | "whatsapp_url">;

/** Solo se acepta un link wa.me; si no, se arma con el mensaje de la API. */
function urlSegura(p: DatosConfirmacion): string {
  if (/^https:\/\/wa\.me\/\d{8,15}\?text=/.test(p.whatsapp_url)) return p.whatsapp_url;
  return linkWhatsApp(p.mensaje_whatsapp);
}

export function Confirmacion({ pedido }: { pedido: DatosConfirmacion }) {
  const [copiado, setCopiado] = useState(false);
  const titulo = useRef<HTMLHeadingElement>(null);
  const url = urlSegura(pedido);

  useEffect(() => {
    titulo.current?.focus();
  }, []);

  async function copiar() {
    try {
      await navigator.clipboard.writeText(pedido.mensaje_whatsapp);
      setCopiado(true);
    } catch {
      // Navegadores internos sin API de portapapeles: se selecciona el texto visible.
      const pre = document.getElementById("mensaje-pedido");
      if (pre) {
        const r = document.createRange();
        r.selectNodeContents(pre);
        const sel = window.getSelection();
        sel?.removeAllRanges();
        sel?.addRange(r);
        try {
          setCopiado(document.execCommand("copy"));
        } catch {
          setCopiado(false);
        }
      }
    }
    window.setTimeout(() => setCopiado(false), 2500);
  }

  return (
    <div className="mx-auto max-w-2xl space-y-5">
      <div className="rounded-[16px] bg-negro p-5 text-blanco">
        {pedido.es_demo ? (
          <p className="mb-2 inline-block rounded-full border border-dashed border-gris-claro px-2 py-0.5 text-[12px] font-semibold text-gris-claro">
            Pedido de prueba (demo)
          </p>
        ) : null}
        <h1 ref={titulo} tabIndex={-1} className="titular text-[34px] text-oro outline-none">
          ¡Listo! Tu pedido es el <span className="whitespace-nowrap">{pedido.numero}</span>
        </h1>
        <p className="mt-3 text-[16px]">Falta un paso: mándanos el pedido por WhatsApp para que te lo confirmemos.</p>
        <a
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          className="mt-5 flex h-[56px] items-center justify-center gap-2 rounded-[14px] bg-whatsapp text-[17px] font-extrabold text-negro"
        >
          <IconoWhatsApp size={24} />
          Enviar pedido por WhatsApp
        </a>
        <p className="mt-4 text-[14px] text-gris-claro">
          ¿No se abrió WhatsApp? Copia el mensaje y pégalo en un chat con el {NEGOCIO.telefonoVisible}.
        </p>
        <button
          type="button"
          onClick={copiar}
          className="mt-2 inline-flex h-11 items-center gap-2 rounded-[10px] border-2 border-oro px-4 text-[15px] font-bold text-oro"
        >
          <IconoCopiar size={18} />
          <span aria-live="polite">{copiado ? "¡Copiado!" : "Copiar mensaje"}</span>
        </button>
      </div>

      <div className="rounded-[16px] bg-blanco p-5">
        <p className="text-[16px] font-bold">
          Guarda tu número de pedido: <span className="whitespace-nowrap">{pedido.numero}</span>
        </p>
        <p className="mt-1 text-[15px]">
          Total: <strong>{clp(pedido.total)}</strong>
          {pedido.entrega === "envio" ? " (no incluye envío)" : ""}
        </p>
        {pedido.entrega === "retiro" ? (
          <p className="mt-3 text-[15px]">
            Espera nuestra confirmación antes de venir. Te esperamos en Alameda 2432 con tu número de pedido.
            <span className="mt-1 block text-[14px] text-gris">
              L a V {HORARIO.semana.abre}–{HORARIO.semana.cierra} · Sáb, dom y feriados {HORARIO.finde.abre}–{HORARIO.finde.cierra}{" "}
              <PorConfirmar />
            </span>
          </p>
        ) : (
          <p className="mt-3 text-[15px]">Te confirmamos por WhatsApp el total, el pago y el despacho.</p>
        )}
        <details className="mt-4">
          <summary className="min-h-11 py-2 text-[14px] font-bold underline underline-offset-2">Ver el mensaje</summary>
          <pre id="mensaje-pedido" className="mt-2 whitespace-pre-wrap break-words rounded-[10px] bg-crema p-3 font-sans text-[13px] leading-relaxed">
            {pedido.mensaje_whatsapp}
          </pre>
        </details>
      </div>

      <div className="flex gap-3 rounded-[16px] border-2 border-negro bg-crema p-4">
        <IconoEscudo size={26} className="shrink-0 text-negro" />
        <p className="text-[15px] font-semibold">
          Te respondemos solo desde el {NEGOCIO.telefonoVisible}. No transfieras a otro número ni a otra cuenta.
        </p>
      </div>

      <Link href="/" className="flex h-12 items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-oro">
        Volver a la tienda
      </Link>
    </div>
  );
}
