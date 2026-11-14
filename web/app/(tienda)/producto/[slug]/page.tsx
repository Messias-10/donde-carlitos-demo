import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Badges } from "@/components/Badges";
import { BotonAgregar } from "@/components/BotonAgregar";
import { Disponibilidad } from "@/components/Disponibilidad";
import { IconoFlechaIzq } from "@/components/Iconos";
import { Ejemplo } from "@/components/Marcas";
import { MinimoPorConfirmar, Precio, unidadUtil } from "@/components/Precio";
import { Descripcion } from "@/components/Descripcion";
import { paraCarrito } from "@/components/TarjetaProducto";
import { VideoTikTok } from "@/components/VideoTikTok";
import { VisualProducto } from "@/components/VisualProducto";
import { obtenerProducto } from "@/lib/api-servidor";
import { clp, ddmm, ddmmaaaa, fechaSantiago } from "@/lib/formato";
import { cache } from "react";

const producto = cache(obtenerProducto);

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
  const { slug } = await params;
  try {
    const p = await producto(slug);
    if (!p) return { title: "Producto no encontrado" };
    const desc = `${p.nombre}: ${clp(p.precio_final)} al detalle${
      p.precio_mayor && p.mayor_minimo ? `, ${clp(p.precio_mayor)} x mayor desde ${p.mayor_minimo} u.` : ""
    }. Pídelo por WhatsApp o retíralo en Alameda 2432.`;
    return {
      title: p.nombre,
      description: desc,
      openGraph: {
        title: `${p.nombre} · Donde Carlitos Emporio`,
        description: desc,
        images: p.imagen ? [{ url: p.imagen, width: 400, height: 500, alt: p.nombre }] : undefined,
      },
    };
  } catch {
    return { title: "Producto" };
  }
}

export default async function FichaProducto({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const p = await producto(slug);
  if (!p) notFound();

  const agotado = p.disponibilidad === "agotado";
  const enOferta = p.precio_final < p.precio;
  const unidad = unidadUtil(p.unidad_mayor, p.mayor_minimo);

  return (
    <div className="bg-crema pb-16">
      <div className="mx-auto max-w-6xl px-4 pt-4">
        <Link href="/catalogo" className="inline-flex min-h-11 items-center gap-1 text-[14px] font-bold">
          <IconoFlechaIzq size={18} /> Volver al catálogo
        </Link>
      </div>
      <article className="mx-auto grid max-w-6xl gap-6 px-4 pt-2 lg:grid-cols-[minmax(0,460px)_1fr] lg:gap-10">
        <div className="relative self-start overflow-hidden rounded-[16px]">
          <VisualProducto
            imagen={p.imagen}
            nombre={p.nombre}
            categoria={p.categoria.slug}
            sizes="(min-width: 1024px) 460px, 100vw"
            preload
            agotado={agotado}
            conVideo={!!p.tiktok_url}
          />
          <div className="absolute left-3 top-3">
            <Badges producto={p} max={3} />
          </div>
        </div>

        <div className="flex flex-col gap-4">
          <p className="text-[14px] font-semibold text-gris">
            <Link href={`/catalogo?categoria=${p.categoria.slug}`} className="underline underline-offset-2">
              {p.categoria.nombre}
            </Link>{" "}
            · Cód. {p.codigo}
          </p>
          <h1 className="font-display text-[34px] uppercase leading-[1.02] lg:text-[44px]">{p.nombre}</h1>
          {p.es_ejemplo ? (
            <p>
              <Ejemplo texto="Producto de ejemplo" />
            </p>
          ) : null}

          <Precio precio={p.precio} precioFinal={p.precio_final} tam="lg" />
          {enOferta && p.oferta_hasta ? (
            <p className="text-[15px] font-bold text-rojo">
              <span aria-hidden>🔥 </span>Oferta hasta el {ddmm(p.oferta_hasta)}
            </p>
          ) : null}

          {p.precio_mayor && p.mayor_minimo ? (
            <div className="rounded-[12px] bg-negro p-4 text-blanco">
              <p className="text-[16px]">
                Llevando <strong className="text-oro">{p.mayor_minimo} o más</strong>, pagas{" "}
                <strong className="text-fluor">{clp(p.precio_mayor)}</strong> c/u.
                {p.mayor_minimo_confirmado === false ? <MinimoPorConfirmar oscuro /> : null}
              </p>
              {unidad ? (
                <p className="mt-1 text-[15px] text-gris-claro">
                  {/^(por|desde|x) /i.test(unidad) ? `Se vende ${unidad}.` : `Viene en ${unidad}.`}
                </p>
              ) : null}
            </div>
          ) : null}

          <Disponibilidad valor={p.disponibilidad} />

          <div className="max-w-sm">
            <BotonAgregar producto={paraCarrito(p)} agotado={agotado} promo={enOferta} conCantidad />
          </div>

          {p.descripcion ? <Descripcion texto={p.descripcion} /> : null}

          {p.precio_actualizado_en ? (
            <p className="text-[13px] text-gris">Precio actualizado por la tienda el {fechaSantiago(p.precio_actualizado_en)}.</p>
          ) : p.fuente ? (
            <p className="text-[13px] text-gris">
              {p.fuente.fecha ? `Precio referencial (video del ${ddmmaaaa(p.fuente.fecha)}). ` : "Precio referencial. "}
              <a
                href={p.fuente.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex min-h-11 items-center underline underline-offset-2"
              >
                Ver fuente
              </a>
            </p>
          ) : null}

          {p.tiktok_url ? (
            <div className="mt-2 w-[60vw] max-w-[260px]">
              <h2 className="titular mb-3 text-[22px]">Míralo en TikTok</h2>
              <VideoTikTok url={p.tiktok_url} titulo={p.nombre} poster={p.imagen} />
            </div>
          ) : null}
        </div>
      </article>
    </div>
  );
}
