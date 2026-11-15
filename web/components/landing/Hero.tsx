import Image from "next/image";
import { IconoMetro, IconoPin, IconoReloj } from "@/components/Iconos";
import { PorConfirmar } from "@/components/Marcas";
import { horarioDeHoy, NEGOCIO, type Numeros } from "@/lib/negocio";
import { BotonesCanal } from "./Canales";

export function Hero({ numeros, demo }: { numeros: Numeros; demo: boolean }) {
  const hoy = horarioDeHoy();
  return (
    <section id="hero" aria-labelledby="hero-titulo" className="relative overflow-hidden bg-negro text-blanco">
      {/* Aro dorado fino que evoca el borde del logo (sin glitter ni animación). */}
      <span
        aria-hidden
        className="pointer-events-none absolute -right-40 -top-40 size-[420px] rounded-full border border-oro/25 lg:-right-24 lg:size-[620px]"
      />
      <div className="relative mx-auto grid max-w-6xl gap-4 px-4 pb-7 pt-5 lg:grid-cols-[1fr_auto] lg:items-center lg:gap-12 lg:pb-14 lg:pt-14">
        <div className="hidden lg:order-2 lg:block">
          <Image
            src="/marca/logo-donde-carlitos.png"
            alt="Donde Carlitos Emporio"
            width={96}
            height={96}
            sizes="(min-width: 1024px) 260px, 96px"
            className="size-[96px] rounded-full shadow-[0_0_0_6px_rgba(226,179,60,0.12)] lg:size-[260px]"
            loading="lazy"
          />
        </div>
        <div className="lg:order-1 lg:max-w-[620px]">
          <div className="flex items-center gap-3">
            {/* En móvil el logo va al lado del saludo: deja la dirección y los botones en la primera pantalla (QA I-7). */}
            <Image
              src="/marca/logo-donde-carlitos.png"
              alt="Donde Carlitos Emporio"
              width={84}
              height={84}
              sizes="84px"
              className="size-[84px] shrink-0 rounded-full shadow-[0_0_0_5px_rgba(226,179,60,0.12)] lg:hidden"
              loading="eager"
            />
            <p className="font-script text-[24px] leading-none text-blanco lg:text-[30px]">¡Hola chiquillos!</p>
          </div>
          <h1 id="hero-titulo" className="titular texto-oro mt-3 text-[40px] lg:text-[72px]">
            ¿Lo viste en TikTok? Pídelo por WhatsApp.
          </h1>
          <p className="mt-3 text-[17px] font-semibold leading-snug text-oro-claro">
            Juguetes virales, Navidad, hogar y escolar, al detalle y por mayor.
          </p>
          <p className="mt-2 max-w-[46ch] text-[15px] leading-snug text-blanco/90">
            Dinos si retiras en tienda o eres de regiones, y te abrimos el chat correcto con el mensaje listo.
          </p>
          <p className="mt-4 flex items-center gap-2 text-[15px] font-semibold text-blanco">
            <IconoPin size={18} className="shrink-0 text-oro" />
            Alameda 2432 · Metro U. Latinoamericana
          </p>
          <div className="mt-3 max-w-md">
            <BotonesCanal numeros={numeros} demo={demo} />
            <p className="mt-3 text-[14px] text-blanco/85">
              Estos son nuestros únicos números de pedidos.{" "}
              <a href="#canal-oficial" className="inline-flex min-h-11 items-center font-bold text-oro underline underline-offset-2">
                Ver canal oficial
              </a>
            </p>
          </div>
        </div>
      </div>

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
