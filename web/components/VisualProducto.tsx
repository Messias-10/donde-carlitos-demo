import Image from "next/image";
import { IconoCategoria } from "./Iconos";

/**
 * Foto del producto si existe en public/. Si no, una placa propia por
 * categoría: las capturas de video con datos de terceros no se usan.
 */
export function VisualProducto({
  imagen,
  nombre,
  categoria,
  sizes,
  preload = false,
  agotado = false,
  conVideo = false,
}: {
  imagen: string | null;
  nombre: string;
  categoria: string;
  sizes?: string;
  preload?: boolean;
  agotado?: boolean;
  conVideo?: boolean;
}) {
  return (
    <div className={`relative aspect-[4/5] w-full overflow-hidden bg-negro ${agotado ? "grayscale" : ""}`}>
      {imagen ? (
        <Image
          src={imagen}
          alt={nombre}
          fill
          sizes={sizes ?? "(min-width: 1024px) 270px, (min-width: 640px) 31vw, 46vw"}
          className="object-cover"
          preload={preload}
          // Las fotos subidas desde el panel ya vienen achicadas y las sirve la API.
          unoptimized={imagen.startsWith("/api/imagenes/")}
          loading={preload ? undefined : "lazy"}
        />
      ) : (
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-3 p-4 text-center">
          <span
            className="absolute inset-3 rounded-[10px] border border-oro/35"
            aria-hidden
          />
          <IconoCategoria slug={categoria} size={56} className="text-oro" />
          <span className="text-[12px] font-semibold leading-tight text-gris-claro">
            {conVideo ? "Foto en el video" : "Foto pronto"}
          </span>
        </div>
      )}
    </div>
  );
}
