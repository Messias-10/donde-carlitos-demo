// Botones de canal de WhatsApp (estilo-carlitos §8.4). Se pintan en el servidor:
// son links wa.me normales, funcionan sin JS y sin la API.
import { IconoWhatsApp } from "@/components/Iconos";
import { PorConfirmar } from "@/components/Marcas";
import type { Numeros } from "@/lib/negocio";
import type { Canal } from "@/lib/tipos";
import { linkWhatsApp, mensajeRegiones, mensajeSantiago } from "@/lib/whatsapp";

export const TEXTO_CANAL: Record<Canal, { largo: string; corto: string; aria: string }> = {
  santiago: {
    largo: "Pedidos Santiago / retiro en tienda",
    corto: "Santiago / retiro",
    aria: "Escribir por WhatsApp a pedidos de Santiago",
  },
  regiones: {
    largo: "Pedidos a regiones",
    corto: "Regiones",
    aria: "Escribir por WhatsApp a pedidos de regiones",
  },
};

export const CLASE_CANAL: Record<Canal, string> = {
  // Santiago: texto negro sobre #25D366 (9,9). Regiones: texto blanco sobre #0B6E4F (6,3).
  santiago: "bg-whatsapp text-negro",
  regiones: "bg-whatsapp-osc text-blanco",
};

export function numeroCanal(n: Numeros, canal: Canal): string {
  return canal === "regiones" ? (n.regiones ?? n.santiago) : n.santiago;
}

export function hrefCanal(n: Numeros, canal: Canal, demo: boolean): string {
  return linkWhatsApp(numeroCanal(n, canal), canal === "santiago" ? mensajeSantiago(demo) : mensajeRegiones(demo));
}

export function BotonCanal({
  canal,
  numeros,
  demo,
  tam = "grande",
  texto,
}: {
  canal: Canal;
  numeros: Numeros;
  demo: boolean;
  tam?: "grande" | "medio";
  texto?: string;
}) {
  const alto = tam === "grande" ? "min-h-14 text-[17px]" : "min-h-12 text-[16px]";
  return (
    <a
      href={hrefCanal(numeros, canal, demo)}
      target="_blank"
      rel="noopener noreferrer"
      className={`flex w-full items-center justify-center gap-2.5 rounded-[14px] px-4 py-2 text-center font-extrabold leading-tight ${alto} ${CLASE_CANAL[canal]}`}
    >
      <IconoWhatsApp size={24} className="shrink-0" />
      <span>{texto ?? TEXTO_CANAL[canal].largo}</span>
    </a>
  );
}

/** Los dos botones del hero, apilados, con la línea de canal oficial debajo. */
export function BotonesCanal({ numeros, demo, oscuro = true }: { numeros: Numeros; demo: boolean; oscuro?: boolean }) {
  return (
    <div className="flex flex-col gap-3">
      <BotonCanal canal="santiago" numeros={numeros} demo={demo} />
      <BotonCanal canal="regiones" numeros={numeros} demo={demo} />
      {numeros.regiones === null ? (
        <p className={`text-[12px] ${oscuro ? "text-gris-claro" : "text-gris"}`}>
          Regiones usa por ahora el mismo número de Santiago. <PorConfirmar oscuro={oscuro} />
        </p>
      ) : null}
    </div>
  );
}
