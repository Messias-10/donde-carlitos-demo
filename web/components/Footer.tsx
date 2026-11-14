import Image from "next/image";
import Link from "next/link";
import { HORARIO, NEGOCIO } from "@/lib/negocio";
import { linkWhatsApp, MSJ_CONSULTA } from "@/lib/whatsapp";
import { PorConfirmar } from "./Marcas";

export function Footer() {
  return (
    <footer className="bg-negro pb-32 pt-12 text-blanco">
      <div className="mx-auto grid max-w-6xl gap-8 px-4 md:grid-cols-[auto_1fr_1fr]">
        <div className="flex items-center gap-3 md:block">
          <Image
            src="/marca/logo-donde-carlitos.png"
            alt="Donde Carlitos Emporio"
            width={64}
            height={64}
            className="size-16 rounded-full"
            sizes="64px"
            loading="lazy"
          />
          <p className="md:mt-3">
            <span className="block font-display text-[20px] uppercase leading-none text-oro">Donde Carlitos Emporio</span>
            <span className="font-script text-[18px] text-blanco">{NEGOCIO.firma}</span>
          </p>
        </div>

        <ul className="space-y-2 text-[15px]">
          <li>
            <span aria-hidden>📍 </span>
            {NEGOCIO.direccion}, {NEGOCIO.comuna}
          </li>
          <li>
            <span aria-hidden>🕘 </span>L a V {HORARIO.semana.abre}–{HORARIO.semana.cierra} · Sáb, dom y feriados{" "}
            {HORARIO.finde.abre}–{HORARIO.finde.cierra} (temporada)
            <PorConfirmar oscuro />
          </li>
          <li className="font-bold">
            Pedidos solo al{" "}
            <a href={linkWhatsApp(MSJ_CONSULTA)} target="_blank" rel="noopener noreferrer" className="inline-flex min-h-11 items-center text-oro underline underline-offset-2">
              {NEGOCIO.telefonoVisible}
            </a>
          </li>
          <li>
            TikTok{" "}
            <a href={NEGOCIO.tiktok.url} target="_blank" rel="noopener noreferrer" className="inline-flex min-h-11 items-center text-oro underline underline-offset-2">
              {NEGOCIO.tiktok.usuario}
            </a>{" "}
            · Instagram{" "}
            <a href={NEGOCIO.instagram.url} target="_blank" rel="noopener noreferrer" className="inline-flex min-h-11 items-center text-oro underline underline-offset-2">
              {NEGOCIO.instagram.usuario}
            </a>
          </li>
        </ul>

        <div className="space-y-3 text-[14px] text-gris-claro">
          <p className="text-blanco">Si te escriben desde otro número a nombre nuestro, no transfieras.</p>
          <p className="text-[15px] font-semibold text-oro">Los esperamos en Alameda 2432 con la mejor variedad de barrio Meiggs.</p>
          <p id="precios-referenciales">Los precios pueden cambiar. Tu pedido guarda el precio del momento en que lo hiciste.</p>
          <p className="rounded-[10px] border border-dashed border-gris-claro/60 p-3">
            Demo de prueba: los pedidos no son reales. Precios sacados de nuestros videos de TikTok (nov 2026).
          </p>
          <p>
            <Link href="/catalogo" className="inline-flex min-h-11 items-center underline underline-offset-2">
              Catálogo
            </Link>{" "}
            ·{" "}
            <Link href="/#preguntas" className="inline-flex min-h-11 items-center underline underline-offset-2">
              Preguntas frecuentes
            </Link>{" "}
            ·{" "}
            <Link href="/admin" className="inline-flex min-h-11 items-center underline underline-offset-2" prefetch={false}>
              Panel
            </Link>
          </p>
        </div>
      </div>
    </footer>
  );
}
