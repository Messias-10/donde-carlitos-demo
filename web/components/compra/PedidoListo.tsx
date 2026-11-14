"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import { NEGOCIO } from "@/lib/negocio";
import { linkWhatsApp, MSJ_CONSULTA } from "@/lib/whatsapp";
import { IconoWhatsApp } from "../Iconos";
import { Confirmacion } from "./Confirmacion";
import { leerConfirmacion, RE_NUMERO } from "./confirmacionGuardada";

const noop = () => () => {};

export function PedidoListo({ numero }: { numero: string }) {
  const valido = RE_NUMERO.test(numero);
  // sessionStorage solo existe en el navegador: el servidor pinta el estado "cargando".
  const raw = useSyncExternalStore(
    noop,
    () => (valido ? JSON.stringify(leerConfirmacion(numero)) : "null"),
    () => "cargando",
  );

  if (raw === "cargando") return <div className="h-80 animate-pulse rounded-[16px] bg-blanco/60" aria-hidden />;
  const datos = raw === "null" ? null : (JSON.parse(raw) as Parameters<typeof Confirmacion>[0]["pedido"]);
  if (datos) return <Confirmacion pedido={datos} />;

  // Otra pestaña o almacenamiento bloqueado: solo el número, sin datos personales.
  return (
    <div className="mx-auto max-w-2xl space-y-4 rounded-[16px] bg-negro p-5 text-blanco">
      {valido ? (
        <>
          <h1 className="titular text-[32px] text-oro">
            Tu pedido es el <span className="whitespace-nowrap">{numero}</span>
          </h1>
          <p className="text-[16px]">
            Si todavía no nos mandaste el pedido por WhatsApp, escríbenos al {NEGOCIO.telefonoVisible} con tu número de
            pedido y te lo confirmamos.
          </p>
        </>
      ) : (
        <h1 className="titular text-[32px] text-oro">No encontramos ese pedido</h1>
      )}
      <a
        href={linkWhatsApp(MSJ_CONSULTA)}
        target="_blank"
        rel="noopener noreferrer"
        className="flex h-[52px] items-center justify-center gap-2 rounded-[14px] bg-whatsapp text-[16px] font-extrabold text-negro"
      >
        <IconoWhatsApp size={22} /> Escríbenos por WhatsApp
      </a>
      <Link href="/" className="flex h-12 items-center justify-center rounded-[12px] border-2 border-oro font-extrabold text-oro">
        Volver a la tienda
      </Link>
    </div>
  );
}
