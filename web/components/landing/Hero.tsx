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
      <div className="relative mx-auto grid max-w-6xl gap-5 px-4 pb-7 pt-6 lg:grid-cols-[1fr_auto] lg:items-center lg:gap-12 lg:pb-14 lg:pt-14">
        <div className="lg:order-2">
          <Image
            src="/marca/logo-donde-carlitos.png"
            alt="Donde Carlitos Emporio"
            width={104}
            height={104}
            sizes="(min-width: 1024px) 260px, 104px"
            className="size-[104px] rounded-full shadow-[0_0_0_6px_rgba(226,179,60,0.12)] lg:size-[260px]"
            preload
          />
        </div>
        <div className="lg:order-1 lg:max-w-[620px]">
          <p className="font-script text-[22px] leading-none text-blanco lg:text-[30px]">¡Hola chiquillos!</p>
          <h1 id="hero-titulo" className="titular texto-oro mt-2 text-[42px] lg:text-[72px]">
            ¿Lo viste en TikTok? Pídelo por WhatsApp.
          </h1>
          <p className="mt-3 text-[18px] font-semibold leading-snug text-oro-claro">
            Juguetes virales, Navidad, hogar y escolar, al detalle y por mayor.
          </p>
          <p className="mt-2 max-w-[46ch] text-[16px] leading-snug text-blanco/90">
            Dinos si retiras en tienda o eres de regiones, y te abrimos el chat correcto con el mensaje listo.
          </p>
          <div className="mt-5 max-w-md">
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
