import type { Disponibilidad as D } from "@/lib/tipos";

const TXT: Record<D, { texto: string; punto: string; color: string }> = {
  disponible: { texto: "Disponible", punto: "bg-verde-stock", color: "text-verde-stock" },
  quedan_pocos: { texto: "Quedan pocas unidades", punto: "bg-naranjo-stock", color: "text-naranjo-stock" },
  agotado: { texto: "Agotado", punto: "bg-rojo", color: "text-rojo" },
};

/** Punto de color + texto (nunca solo color). */
export function Disponibilidad({ valor }: { valor: D }) {
  const t = TXT[valor] ?? TXT.disponible;
  return (
    <p className={`flex items-center gap-1.5 text-[13px] font-semibold ${t.color}`}>
      <span className={`inline-block size-2 rounded-full ${t.punto}`} aria-hidden />
      {t.texto}
    </p>
  );
}
