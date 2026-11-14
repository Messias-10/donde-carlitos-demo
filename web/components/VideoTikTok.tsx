"use client";

import Image from "next/image";
import { useState } from "react";
import { idTikTok } from "@/lib/tiktok";
import { IconoPlay } from "./Iconos";


/**
 * Póster propio y, solo al tocar, el reproductor oficial de TikTok en un
 * iframe. Sin embed.js y sin autoplay al cargar la página.
 */
export function VideoTikTok({
  url,
  titulo,
  poster,
  avisoDireccion = false,
}: {
  url: string;
  titulo: string;
  poster?: string | null;
  avisoDireccion?: boolean;
}) {
  const [activo, setActivo] = useState(false);
  const id = idTikTok(url);

  return (
    <div className="flex flex-col gap-2">
      <div className="relative aspect-[9/16] w-full overflow-hidden rounded-[12px] bg-negro ring-1 ring-oro/40">
        {activo && id ? (
          <iframe
            src={`https://www.tiktok.com/player/v1/${id}?autoplay=1&rel=0&description=0&music_info=0`}
            title={`Video de TikTok: ${titulo}`}
            className="absolute inset-0 size-full"
            allow="autoplay; encrypted-media; fullscreen; picture-in-picture"
            referrerPolicy="strict-origin-when-cross-origin"
            sandbox="allow-scripts allow-same-origin allow-presentation allow-popups allow-popups-to-escape-sandbox"
            loading="lazy"
          />
        ) : (
          <button
            type="button"
            onClick={() => (id ? setActivo(true) : window.open(url, "_blank", "noopener,noreferrer"))}
            className="group absolute inset-0 flex flex-col items-center justify-center text-blanco"
            aria-label={`Ver el video de ${titulo}`}
          >
            {poster ? (
              <Image
                src={poster}
                alt=""
                fill
                sizes="(min-width: 1024px) 200px, 44vw"
                unoptimized={poster.startsWith("/api/imagenes/")}
                className="object-cover opacity-80"
              />
            ) : null}
            <span className="absolute inset-0 bg-gradient-to-t from-negro via-negro/30 to-transparent" aria-hidden />
            <span className="relative grid size-16 place-items-center rounded-full bg-negro/70 ring-2 ring-oro transition-transform group-active:scale-90">
              <IconoPlay size={30} className="translate-x-0.5 text-oro" />
            </span>
            <span className="relative mt-3 px-3 text-center text-[13px] font-bold">▶ Míralo en TikTok</span>
            {avisoDireccion ? (
              <span className="absolute inset-x-2 top-2 rounded-[8px] bg-fluor px-2 py-1 text-center text-[12px] font-extrabold text-negro">
                Hoy nos encuentras en Alameda 2432
              </span>
            ) : null}
          </button>
        )}
      </div>
      <a
        href={url}
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex min-h-11 items-center text-[13px] font-semibold underline underline-offset-2"
      >
        Abrir en TikTok
      </a>
    </div>
  );
}
