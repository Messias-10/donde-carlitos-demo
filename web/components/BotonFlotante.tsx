"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cotizar, useCarrito } from "@/lib/carrito";
import { clp } from "@/lib/formato";
import { linkWhatsApp, MSJ_CONSULTA } from "@/lib/whatsapp";
import { IconoCarrito, IconoWhatsApp } from "./Iconos";

const OCULTO_EN = ["/carrito", "/pedido", "/admin"];

/**
 * Sin carrito: WhatsApp de consulta. Con carrito: lleva al pedido, para que
 * se guarde con número antes de ir a WhatsApp (arquitectura-landing.md §14).
 */
export function BotonFlotante() {
  const ruta = usePathname();
  const items = useCarrito();
  if (OCULTO_EN.some((r) => ruta === r || ruta.startsWith(r + "/"))) return null;

  const { total, unidades } = cotizar(items);
  const conCarrito = unidades > 0;

  return (
    <div className="pointer-events-none fixed inset-x-0 bottom-0 z-40 px-3 pb-[max(10px,env(safe-area-inset-bottom))] pt-2 lg:inset-x-auto lg:right-6 lg:pb-6">
      <div className="pointer-events-auto mx-auto max-w-md lg:w-[360px]">
        {conCarrito ? (
          <Link
            href="/carrito"
            className="flex h-[52px] items-center justify-center gap-2 rounded-[14px] bg-oro px-4 text-[16px] font-extrabold text-negro shadow-[0_8px_24px_-8px_rgba(0,0,0,0.6)] ring-2 ring-negro"
          >
            <IconoCarrito size={22} />
            Ver mi pedido ({unidades}) · {clp(total)}
          </Link>
        ) : (
          <a
            href={linkWhatsApp(MSJ_CONSULTA)}
            target="_blank"
            rel="noopener noreferrer"
            className="flex h-[52px] items-center justify-center gap-2 rounded-[14px] bg-whatsapp px-4 text-[16px] font-extrabold text-negro shadow-[0_8px_24px_-8px_rgba(0,0,0,0.6)]"
          >
            <IconoWhatsApp size={24} />
            Escríbenos por WhatsApp
          </a>
        )}
        <p className="mt-1 text-center text-[11px] font-semibold text-negro">
          <span className="rounded-full bg-crema/90 px-2 py-0.5">Único número de pedidos</span>
        </p>
      </div>
    </div>
  );
}
