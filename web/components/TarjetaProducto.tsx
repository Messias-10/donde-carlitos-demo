import Link from "next/link";
import type { ProductoVista } from "@/lib/tipos";
import type { ProductoParaCarrito } from "@/lib/carrito";
import { ddmm } from "@/lib/formato";
import { Badges } from "./Badges";
import { BotonAgregar } from "./BotonAgregar";
import { Disponibilidad } from "./Disponibilidad";
import { IconoPlay } from "./Iconos";
import { Ejemplo } from "./Marcas";
import { LineaMayor, Precio } from "./Precio";
import { VisualProducto } from "./VisualProducto";

export function paraCarrito(p: ProductoVista): ProductoParaCarrito {
  return {
    id: p.id,
    slug: p.slug,
    codigo: p.codigo,
    nombre: p.nombre,
    precio: p.precio,
    precio_final: p.precio_final,
    precio_mayor: p.precio_mayor,
    mayor_minimo: p.mayor_minimo,
    unidad_mayor: p.unidad_mayor,
    categoria: p.categoria.slug,
  };
}

/** Tarjeta de producto (estilo-carlitos §8.1). */
export function TarjetaProducto({
  p,
  mostrarLlegada = false,
  preload = false,
}: {
  p: ProductoVista;
  mostrarLlegada?: boolean;
  preload?: boolean;
}) {
  const agotado = p.disponibilidad === "agotado";
  const enOferta = p.precio_final < p.precio;
  return (
    <article
      className={`relative flex h-full flex-col overflow-hidden rounded-[12px] bg-blanco shadow-[0_1px_0_rgba(11,11,11,0.08),0_6px_18px_-10px_rgba(11,11,11,0.35)] ${
        agotado ? "opacity-80" : ""
      }`}
    >
      <Link href={`/producto/${p.slug}`} className="relative block" aria-label={p.nombre} tabIndex={-1}>
        <VisualProducto
          imagen={p.imagen}
          nombre={p.nombre}
          categoria={p.categoria.slug}
          agotado={agotado}
          preload={preload}
          conVideo={!!p.tiktok_url}
        />
        <div className="absolute inset-x-2 top-2 flex items-start justify-between gap-2">
          <Badges producto={p} />
          {p.tiktok_url ? (
            <span className="ml-auto inline-flex shrink-0 items-center gap-1 rounded-full bg-negro/80 px-2 py-1 text-[11px] font-bold text-blanco">
              <IconoPlay size={12} /> TikTok
            </span>
          ) : null}
        </div>
        {agotado ? (
          <span className="absolute inset-x-0 bottom-3 mx-auto w-fit rounded-full bg-rojo px-3 py-1 text-[13px] font-extrabold text-blanco">
            AGOTADO
          </span>
        ) : null}
      </Link>
      <div className="flex flex-1 flex-col gap-2 p-3">
        <h3 className="line-clamp-2 text-[15px] font-semibold leading-snug">
          <Link href={`/producto/${p.slug}`} className="hover:underline">
            {p.nombre}
          </Link>
        </h3>
        {p.es_ejemplo ? <Ejemplo texto="producto de ejemplo" /> : null}
        {mostrarLlegada && p.fecha_llegada ? (
          <p className="text-[13px] text-gris">Llegó el {ddmm(p.fecha_llegada)}</p>
        ) : null}
        <Precio precio={p.precio} precioFinal={p.precio_final} />
        <LineaMayor precioMayor={p.precio_mayor} minimo={p.mayor_minimo} unidad={p.unidad_mayor} />
        <Disponibilidad valor={p.disponibilidad} />
        <div className="mt-auto pt-1">
          <BotonAgregar producto={paraCarrito(p)} agotado={agotado} promo={enOferta} />
        </div>
      </div>
    </article>
  );
}
