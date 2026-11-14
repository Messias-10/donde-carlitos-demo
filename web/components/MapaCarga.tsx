"use client";

import { useState } from "react";
import { NEGOCIO } from "@/lib/negocio";
import { IconoPin } from "./Iconos";

/** El mapa se carga solo si el usuario lo pide (ahorra datos en el celular). */
export function MapaCarga() {
  const [ver, setVer] = useState(false);
  return (
    <div className="overflow-hidden rounded-[12px] bg-carbon ring-1 ring-oro/40">
      {ver ? (
        <iframe
          src={NEGOCIO.mapsEmbed}
          title="Mapa: Alameda 2432, Santiago Centro"
          className="aspect-[4/3] w-full"
          loading="lazy"
          referrerPolicy="no-referrer-when-downgrade"
        />
      ) : (
        <button
          type="button"
          onClick={() => setVer(true)}
          className="flex aspect-[4/3] w-full flex-col items-center justify-center gap-3 text-blanco"
        >
          <span className="grid size-16 place-items-center rounded-full bg-negro ring-2 ring-oro">
            <IconoPin size={30} className="text-oro" />
          </span>
          <span className="text-[15px] font-bold">Ver mapa aquí</span>
          <span className="text-[12px] text-gris-claro">Se carga desde Google Maps</span>
        </button>
      )}
    </div>
  );
}
