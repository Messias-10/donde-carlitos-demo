import { clp, miles } from "@/lib/formato";

type Tam = "sm" | "md" | "lg";

const TAM: Record<Tam, string> = {
  sm: "text-[22px]",
  md: "text-[30px]",
  lg: "text-[42px]",
};

/** Placa flúor con el precio que se cobra hoy y, si hay oferta real, el anterior tachado. */
export function Precio({
  precio,
  precioFinal,
  tam = "md",
  etiqueta = "detalle",
  oscuro = false,
}: {
  precio: number;
  precioFinal: number;
  tam?: Tam;
  etiqueta?: string;
  oscuro?: boolean;
}) {
  const conOferta = precioFinal < precio;
  return (
    <div className="flex flex-wrap items-end gap-x-2 gap-y-1">
      <span className={`placa-precio ${TAM[tam]}`} data-oferta={conOferta}>
        <span className="signo">$</span>
        {miles(precioFinal)}
        <span className="sr-only"> pesos</span>
      </span>
      <span className={`pb-0.5 text-[13px] leading-tight ${oscuro ? "text-gris-claro" : "text-gris"}`}>
        {conOferta ? (
          <>
            <s className={`block text-[14px] font-semibold ${oscuro ? "text-blanco/80" : "text-negro/70"}`}>
              <span className="sr-only">Antes </span>
              {clp(precio)}
            </s>
            oferta
          </>
        ) : (
          etiqueta
        )}
      </span>
    </div>
  );
}

/**
 * La unidad no aporta si repite el mínimo ("desde 2 u.") o solo dice "por mayor".
 * Así se evita "desde 2 u. · desde 2 u." (QA M-13).
 */
export function unidadUtil(unidad: string | null, minimo: number | null): string | null {
  if (!unidad) return null;
  const u = unidad.toLowerCase().replace(/\s+/g, " ").trim();
  if (/^(por|x|al) mayor$/.test(u)) return null;
  if (minimo && (u === `desde ${minimo} u.` || u === `desde ${minimo} u` || u === `desde ${minimo} unidades`)) return null;
  return unidad;
}

/** Línea de precio por mayor: "$3.500 x mayor desde 12 u. · caja de 12". */
export function LineaMayor({
  precioMayor,
  minimo,
  unidad,
  confirmado = true,
  oscuro = false,
  formato = "completo",
  className = "",
}: {
  precioMayor: number | null;
  minimo: number | null;
  unidad: string | null;
  confirmado?: boolean;
  oscuro?: boolean;
  /** vitrina: "$3.500 c/u · caja de 12" (copy-landing.md §3). */
  formato?: "completo" | "vitrina";
  className?: string;
}) {
  if (!precioMayor || !minimo) return null;
  if (formato === "vitrina") {
    const u = unidadUtil(unidad, minimo);
    return (
      <p className={`text-[14px] font-semibold leading-snug ${oscuro ? "text-blanco" : "text-negro"} ${className}`}>
        {clp(precioMayor)} c/u{" "}
        <span className={`font-normal ${oscuro ? "text-gris-claro" : "text-gris"}`}>· {u ?? `desde ${minimo} u.`}</span>
        {confirmado ? null : <MinimoPorConfirmar oscuro={oscuro} />}
      </p>
    );
  }
  const u = unidadUtil(unidad, minimo);
  const gris = oscuro ? "text-gris-claro" : "text-gris";
  return (
    <p className={`text-[14px] font-semibold leading-snug ${oscuro ? "text-blanco" : "text-negro"} ${className}`}>
      {clp(precioMayor)} <span className={`font-normal ${gris}`}>x mayor desde {minimo} u.</span>
      {u ? <span className={`font-normal ${gris}`}> · {u}</span> : null}
      {confirmado ? null : <MinimoPorConfirmar oscuro={oscuro} />}
    </p>
  );
}

/** I-2: el mínimo por mayor de la semilla es [EJEMPLO]. */
export function MinimoPorConfirmar({ oscuro = false }: { oscuro?: boolean }) {
  return (
    <span className={`ml-1 inline-block whitespace-nowrap rounded-full border px-2 py-px align-middle text-[11px] font-semibold leading-4 ${oscuro ? "border-gris-claro/60 bg-carbon text-gris-claro" : "border-gris/50 bg-blanco text-gris"}`}>
      mínimo por confirmar
    </span>
  );
}
