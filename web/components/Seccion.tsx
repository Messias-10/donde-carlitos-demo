import type { ReactNode } from "react";

type Fondo = "crema" | "blanco" | "negro" | "carbon";

const FONDO: Record<Fondo, string> = {
  crema: "bg-crema text-negro",
  blanco: "bg-blanco text-negro",
  negro: "bg-negro text-blanco",
  carbon: "bg-carbon text-blanco",
};

export function Seccion({
  id,
  fondo = "crema",
  titulo,
  script,
  bajada,
  children,
  className = "",
}: {
  id?: string;
  fondo?: Fondo;
  titulo: ReactNode;
  script?: string;
  bajada?: ReactNode;
  children: ReactNode;
  className?: string;
}) {
  const oscuro = fondo === "negro" || fondo === "carbon";
  const idTitulo = id ? `${id}-titulo` : undefined;
  return (
    <section id={id} aria-labelledby={idTitulo} className={`${FONDO[fondo]} py-12 [contain-intrinsic-size:auto_900px] [content-visibility:auto] lg:py-16 ${className}`}>
      <div className="mx-auto max-w-6xl px-4">
        {script ? <p className={`font-script text-[20px] ${oscuro ? "text-blanco" : "text-oro-texto"}`}>{script}</p> : null}
        <h2 id={idTitulo} className={`titular text-[32px] lg:text-[44px] ${oscuro ? "text-oro" : "text-negro"}`}>
          {titulo}
        </h2>
        {bajada ? (
          <p className={`mt-3 max-w-[62ch] text-[16px] ${oscuro ? "text-gris-claro" : "text-negro/80"}`}>{bajada}</p>
        ) : null}
        <div className="mt-6">{children}</div>
      </div>
    </section>
  );
}

/** Fila deslizable en móvil, grilla en escritorio. */
export function FilaProductos({ children, etiqueta }: { children: ReactNode; etiqueta: string }) {
  return (
    <div
      className="scroll-fila -mx-4 flex items-start gap-3 overflow-x-auto px-4 pb-2 lg:mx-0 lg:grid lg:grid-cols-4 lg:items-stretch lg:overflow-visible lg:px-0"
      role="list"
      aria-label={etiqueta}
    >
      {children}
    </div>
  );
}

export function ItemFila({ children }: { children: ReactNode }) {
  return (
    <div role="listitem" className="w-[60vw] max-w-[240px] shrink-0 lg:w-auto lg:max-w-none">
      {children}
    </div>
  );
}
