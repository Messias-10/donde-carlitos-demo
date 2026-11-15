"use client";

import { useState } from "react";
import { IconoCopiar } from "./Iconos";

/** Botón "Copiar …" que cambia a "¡Copiado!". Sin portapapeles: selecciona el texto. */
export function CopiarTexto({
  texto,
  etiqueta = "Copiar",
  className = "",
  idObjetivo,
}: {
  texto: string;
  etiqueta?: string;
  className?: string;
  /** id de un elemento cuyo contenido se selecciona si el portapapeles no está disponible. */
  idObjetivo?: string;
}) {
  const [copiado, setCopiado] = useState(false);
  async function copiar() {
    let ok = false;
    try {
      await navigator.clipboard.writeText(texto);
      ok = true;
    } catch {
      const el = idObjetivo ? document.getElementById(idObjetivo) : null;
      if (el) {
        const r = document.createRange();
        r.selectNodeContents(el);
        const sel = window.getSelection();
        sel?.removeAllRanges();
        sel?.addRange(r);
        try {
          ok = document.execCommand("copy");
        } catch {
          ok = false;
        }
      }
    }
    setCopiado(ok);
    if (ok) window.setTimeout(() => setCopiado(false), 2500);
  }
  return (
    <button
      type="button"
      onClick={copiar}
      className={`inline-flex min-h-11 items-center gap-2 rounded-[10px] border-2 px-3 text-[14px] font-bold ${className}`}
    >
      <IconoCopiar size={16} />
      <span aria-live="polite">{copiado ? "¡Copiado!" : etiqueta}</span>
    </button>
  );
}
