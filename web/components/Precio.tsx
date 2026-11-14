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

/** Línea de precio por mayor: "$3.500 x mayor desde 12 u. · caja de 12". */
export function LineaMayor({
  precioMayor,
  minimo,
  unidad,
  className = "",
}: {
  precioMayor: number | null;
  minimo: number | null;
  unidad: string | null;
  className?: string;
}) {
  if (!precioMayor || !minimo) return null;
  return (
    <p className={`text-[14px] font-semibold leading-snug text-negro ${className}`}>
      {clp(precioMayor)} <span className="font-normal text-gris">x mayor desde {minimo} u.</span>
      {unidad ? <span className="font-normal text-gris"> · {unidad}</span> : null}
    </p>
  );
}
