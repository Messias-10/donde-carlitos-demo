import Image from "next/image";
import Link from "next/link";
import { IconoMetro, IconoPin, IconoReloj } from "@/components/Iconos";
import { PorConfirmar } from "@/components/Marcas";
import { horarioDeHoy, NEGOCIO } from "@/lib/negocio";

export function Hero() {
  const hoy = horarioDeHoy();
  return (
    <section aria-labelledby="hero-titulo" className="relative overflow-hidden bg-negro text-blanco">
      {/* Aro dorado fino que evoca el borde del logo (sin glitter ni animación). */}
      <span
        aria-hidden
        className="pointer-events-none absolute -right-40 -top-40 size-[420px] rounded-full border border-oro/25 lg:-right-24 lg:size-[620px]"
      />
      <span
        aria-hidden
        className="pointer-events-none absolute -right-28 -top-28 size-[340px] rounded-full border-2 border-oro/15 lg:size-[520px]"
      />
      <div className="relative mx-auto grid max-w-6xl gap-6 px-4 pb-8 pt-8 lg:grid-cols-[1fr_auto] lg:items-center lg:pb-14 lg:pt-14">
        <div className="lg:order-2">
          <Image
            src="/marca/logo-donde-carlitos.png"
            alt="Donde Carlitos Emporio"
            width={150}
            height={150}
            sizes="(min-width: 1024px) 280px, 150px"
            className="size-[150px] rounded-full shadow-[0_0_0_6px_rgba(226,179,60,0.12)] lg:size-[280px]"
            preload
          />
        </div>
        <div className="lg:order-1">
          <p className="font-script text-[24px] leading-none text-blanco lg:text-[30px]">¡Hola chiquillos!</p>
          <h1 id="hero-titulo" className="titular texto-oro mt-2 text-[46px] lg:text-[76px]">
            ¿Lo viste en TikTok? Acá lo pides.
          </h1>
          <p className="mt-4 max-w-[46ch] text-[17px] leading-snug text-blanco/90">
            Precio al detalle y por mayor, a la vista. Arma tu pedido y mándalo directo a nuestro WhatsApp oficial.
          </p>
          <div className="mt-6 flex flex-wrap gap-3">
            <a
              href="#ofertas"
              className="inline-flex h-12 items-center rounded-[12px] bg-oro px-6 text-[16px] font-extrabold text-negro"
            >
              Ver ofertas
            </a>
            <Link
              href="/catalogo"
              className="inline-flex h-12 items-center rounded-[12px] border-2 border-oro px-6 text-[16px] font-extrabold text-oro"
            >
              Ver catálogo
            </Link>
          </div>
        </div>
      </div>

      {/* Barra de datos: lo que ya dice cualquier caption del barrio. */}
      <div className="relative border-t border-oro/30 bg-carbon">
        <ul className="mx-auto grid max-w-6xl gap-2 px-4 py-3 text-[14px] sm:grid-cols-3 sm:gap-4">
          <li className="flex items-start gap-2">
            <IconoPin size={18} className="mt-0.5 shrink-0 text-oro" />
            <span>
              {NEGOCIO.direccion}, {NEGOCIO.comuna}
            </span>
          </li>
          <li className="flex items-start gap-2">
            <IconoMetro size={18} className="mt-0.5 shrink-0 text-oro" />
            <span>a 5 min del Metro Unión Latinoamericana</span>
          </li>
          <li className="flex items-start gap-2">
            <IconoReloj size={18} className="mt-0.5 shrink-0 text-oro" />
            <span>
              Hoy: {hoy.abre} a {hoy.cierra}
              <PorConfirmar oscuro />
            </span>
          </li>
        </ul>
      </div>
    </section>
  );
}
