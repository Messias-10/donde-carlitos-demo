import Image from "next/image";
import { HORARIO, NEGOCIO, telefonoVisible, type Numeros } from "@/lib/negocio";
import { PorConfirmar } from "./Marcas";

const LINK = "inline-flex min-h-11 items-center text-oro underline underline-offset-2";

export function Footer({ numeros, demo }: { numeros: Numeros; demo: boolean }) {
  return (
    <footer className="bg-negro pb-28 pt-12 text-blanco lg:pb-12">
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
          <li className="font-bold">Santiago y retiro: {telefonoVisible(numeros.santiago)}</li>
          <li className="font-bold">
            Regiones:{" "}
            {numeros.regiones ? (
              telefonoVisible(numeros.regiones)
            ) : (
              <>
                {telefonoVisible(numeros.santiago)} <PorConfirmar oscuro />
              </>
            )}
          </li>
          <li>
            TikTok{" "}
            <a href={NEGOCIO.tiktok.url} target="_blank" rel="noopener noreferrer" className={LINK}>
              {NEGOCIO.tiktok.usuario}
            </a>{" "}
            · Instagram{" "}
            <a href={NEGOCIO.instagram.url} target="_blank" rel="noopener noreferrer" className={LINK}>
              {NEGOCIO.instagram.usuario}
            </a>
          </li>
        </ul>

        <div className="space-y-3 text-[14px] text-gris-claro">
          <p className="text-blanco">
            Pedidos solo a nuestros números oficiales. Los datos para transferir te los damos solo por el chat oficial.
          </p>
          <p className="text-[15px] font-semibold text-oro">Los esperamos en Alameda 2432 con la mejor variedad de barrio Meiggs.</p>
          {demo ? (
            <p className="rounded-[10px] border border-dashed border-gris-claro/60 p-3">
              Demo de prueba. Precios sacados de nuestros videos de TikTok (nov 2026).
            </p>
          ) : null}
        </div>
      </div>
    </footer>
  );
}
