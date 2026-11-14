// Marcas visibles para datos no confirmados o de ejemplo (copy-landing.md §1).

export function PorConfirmar({ oscuro = false, texto = "dato por confirmar" }: { oscuro?: boolean; texto?: string }) {
  return (
    <span
      className={`ml-1 inline-block whitespace-nowrap rounded-full border px-2 py-px align-middle text-[11px] font-semibold leading-4 ${
        oscuro ? "border-gris-claro/60 text-gris-claro" : "border-gris/50 text-gris"
      }`}
    >
      {texto}
    </span>
  );
}

export function Ejemplo({ oscuro = false, texto = "ejemplo" }: { oscuro?: boolean; texto?: string }) {
  return (
    <span
      className={`inline-block whitespace-nowrap rounded-full border border-dashed px-2 py-px align-middle text-[11px] font-semibold leading-4 ${
        oscuro ? "border-gris-claro text-gris-claro" : "border-gris text-gris"
      }`}
    >
      {texto}
    </span>
  );
}
