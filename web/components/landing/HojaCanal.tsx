"use client";

import { useEffect, useRef, useState } from "react";
import { IconoWhatsApp } from "@/components/Iconos";
import type { Canal } from "@/lib/tipos";
import { linkWhatsApp, mensajeProducto, type DatosVitrina } from "@/lib/whatsapp";

const CLAVE = "dc_canal";

function leerCanal(): Canal | null {
  try {
    const v = window.sessionStorage.getItem(CLAVE);
    return v === "santiago" || v === "regiones" ? v : null;
  } catch {
    return null;
  }
}

function guardarCanal(c: Canal) {
  try {
    window.sessionStorage.setItem(CLAVE, c);
  } catch {
    /* sin almacenamiento: no se recuerda, no pasa nada */
  }
}

function leerProducto(raw: string | undefined): DatosVitrina | null {
  if (!raw) return null;
  try {
    const p = JSON.parse(raw) as DatosVitrina;
    return typeof p.nombre === "string" && typeof p.precio === "number" ? p : null;
  } catch {
    return null;
  }
}

/**
 * Hoja "¿DÓNDE LO QUIERES?" (copy-landing.md §3). Un solo componente para toda la
 * vitrina: escucha los botones [data-pedir] (pintados en el servidor) y abre un
 * <dialog> con los dos canales. La elección se recuerda en sessionStorage durante
 * la visita y queda primera, pero siempre se puede cambiar.
 */
export function HojaCanal({ numeros, demo }: { numeros: Record<Canal, string>; demo: boolean }) {
  const ref = useRef<HTMLDialogElement>(null);
  const [producto, setProducto] = useState<DatosVitrina | null>(null);
  const [recordado, setRecordado] = useState<Canal | null>(null);

  useEffect(() => {
    function onClick(e: MouseEvent) {
      const b = (e.target as Element | null)?.closest<HTMLElement>("[data-pedir]");
      if (!b) return;
      const p = leerProducto(b.dataset.pedir);
      if (!p) return;
      setProducto(p);
      setRecordado(leerCanal());
      ref.current?.showModal();
    }
    document.addEventListener("click", onClick);
    return () => document.removeEventListener("click", onClick);
  }, []);

  const orden: Canal[] = recordado === "regiones" ? ["regiones", "santiago"] : ["santiago", "regiones"];
  const texto: Record<Canal, string> = { santiago: "Santiago / retiro en tienda", regiones: "Regiones (envío)" };

  return (
    <dialog
      ref={ref}
      aria-labelledby="hoja-titulo"
      className="m-0 mt-auto w-full max-w-none rounded-t-[20px] bg-negro p-0 text-blanco backdrop:bg-negro/60 sm:m-auto sm:max-w-md sm:rounded-[20px]"
      onClick={(e) => {
        if (e.target === ref.current) ref.current?.close();
      }}
    >
      <div className="px-5 pb-[max(20px,env(safe-area-inset-bottom))] pt-5">
        <span aria-hidden className="mx-auto mb-4 block h-1.5 w-12 rounded-full bg-blanco/30 sm:hidden" />
        <h2 id="hoja-titulo" className="titular text-[30px] text-blanco">
          ¿Dónde lo quieres?
        </h2>
        {producto ? <p className="mt-1 line-clamp-2 text-[15px] text-gris-claro">{producto.nombre}</p> : null}
        <div className="mt-5 flex flex-col gap-3">
          {producto
            ? orden.map((c) => (
                <a
                  key={c}
                  href={linkWhatsApp(numeros[c], mensajeProducto(producto, c, demo))}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={() => {
                    guardarCanal(c);
                    ref.current?.close();
                  }}
                  className={`flex min-h-14 items-center justify-center gap-2.5 rounded-[14px] px-4 text-[17px] font-extrabold ${
                    c === "santiago" ? "bg-whatsapp text-negro" : "bg-whatsapp-osc text-blanco"
                  }`}
                >
                  <IconoWhatsApp size={22} />
                  {texto[c]}
                  {recordado === c ? <span className="sr-only"> (tu elección anterior)</span> : null}
                </a>
              ))
            : null}
        </div>
        <p className="mt-4 text-[14px] text-gris-claro">Te abrimos WhatsApp con el producto ya escrito. Tú solo envías.</p>
        <button
          type="button"
          onClick={() => ref.current?.close()}
          className="mt-2 min-h-11 w-full text-[15px] font-bold text-oro underline underline-offset-2"
        >
          Cancelar
        </button>
      </div>
    </dialog>
  );
}
