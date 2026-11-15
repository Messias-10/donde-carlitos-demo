// "Lo nuevo de los videos" (arquitectura §A2, estilo-carlitos §8.1).
// Sin carrito, sin cantidades y sin stock. El botón abre la hoja de canal (HojaCanal).
import { BadgeVitrina } from "@/components/Badges";
import { IconoWhatsApp } from "@/components/Iconos";
import { Ejemplo } from "@/components/Marcas";
import { LineaMayor, Precio } from "@/components/Precio";
import { Seccion } from "@/components/Seccion";
import { VisualProducto } from "@/components/VisualProducto";
import { ddmm, fechaSantiago } from "@/lib/formato";
import type { Numeros } from "@/lib/negocio";
import type { ProductoVista } from "@/lib/tipos";
import type { DatosVitrina } from "@/lib/whatsapp";
import { BotonCanal } from "./Canales";

/** "dd/mm" del precio: el cambio hecho en el panel o, si no hay, la fecha del video. */
export function fechaPrecio(p: ProductoVista): string {
  if (p.precio_actualizado_en) return fechaSantiago(p.precio_actualizado_en).slice(0, 5);
  return ddmm(p.fuente?.fecha ?? null);
}

function datosMensaje(p: ProductoVista): DatosVitrina {
  return {
    nombre: p.nombre,
    codigo: p.codigo,
    precio: p.precio_final,
    precio_mayor: p.precio_mayor,
    unidad_mayor: p.unidad_mayor,
    fecha_precio: fechaPrecio(p),
    es_promo: p.precio_es_promo === true,
  };
}

function Tarjeta({ p, preload }: { p: ProductoVista; preload: boolean }) {
  const fecha = fechaPrecio(p);
  return (
    <article className="relative flex h-full flex-col overflow-hidden rounded-[12px] bg-blanco shadow-[0_1px_0_rgba(11,11,11,0.08),0_6px_18px_-10px_rgba(11,11,11,0.35)]">
      <div className="relative">
        <VisualProducto
          imagen={p.imagen}
          nombre={p.nombre}
          categoria={p.categoria.slug}
          sizes="(min-width: 1024px) 270px, (min-width: 640px) 31vw, 46vw"
          preload={preload}
          conVideo={!!p.tiktok_url}
        />
        <div className="absolute left-2 top-2">
          <BadgeVitrina producto={p} />
        </div>
      </div>
      <div className="flex flex-1 flex-col gap-2 p-3">
        <h3 className="line-clamp-2 text-[15px] font-semibold leading-snug">{p.nombre}</h3>
        {p.es_ejemplo ? <Ejemplo texto="producto de ejemplo" /> : null}
        <Precio precio={p.precio} precioFinal={p.precio_final} tam="sm" etiqueta={p.precio_es_promo ? "" : "detalle"} />
        {p.precio_es_promo && p.nota_promo ? <p className="text-[13px] font-semibold text-magenta">{p.nota_promo}</p> : null}
        <LineaMayor
          precioMayor={p.precio_mayor}
          minimo={p.mayor_minimo}
          unidad={p.unidad_mayor}
          confirmado={p.mayor_minimo_confirmado !== false}
          formato="vitrina"
        />
        {fecha ? (
          <p className="text-[13px] text-gris">
            {p.precio_es_promo ? "Promo publicada el" : "Precio publicado el"} {fecha}
          </p>
        ) : null}
        {p.tiktok_url ? (
          <a
            href={p.tiktok_url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex min-h-11 items-center text-[14px] font-semibold underline underline-offset-2"
          >
            ▶ Míralo en TikTok
          </a>
        ) : null}
        <button
          type="button"
          data-pedir={JSON.stringify(datosMensaje(p))}
          className="mt-auto flex min-h-12 w-full items-center justify-center gap-2 rounded-[12px] bg-negro px-2 text-center text-[14px] font-extrabold leading-tight text-oro"
        >
          <IconoWhatsApp size={18} className="hidden shrink-0 text-oro sm:inline-block" />
          Pedir este por WhatsApp
          <span className="sr-only">: {p.nombre}</span>
        </button>
      </div>
    </article>
  );
}

export function Vitrina({ productos, numeros, demo }: { productos: ProductoVista[]; numeros: Numeros; demo: boolean }) {
  return (
    <Seccion
      id="vitrina"
      titulo={
        <>
          Lo nuevo de los videos <span aria-hidden>✨</span>
        </>
      }
      bajada="Lo que mostramos esta semana en TikTok, con su precio. Tócalo y nos escribes con el producto ya puesto."
    >
      {productos.length > 0 ? (
        <ul className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
          {productos.map((p, i) => (
            <li key={p.id}>
              <Tarjeta p={p} preload={i < 2} />
            </li>
          ))}
        </ul>
      ) : (
        <p className="rounded-[12px] bg-blanco p-4 text-[15px]">Pronto subimos lo nuevo de la semana. Mientras, escríbenos por WhatsApp.</p>
      )}
      <p className="mt-5 text-[14px] text-gris">
        Los precios son los que dijimos en el video y pueden cambiar. Pregúntanos por WhatsApp si queda.
      </p>
      <div className="mt-6 rounded-[16px] bg-blanco p-4">
        <p className="text-[16px] font-bold">¿Buscas otra cosa? Escríbenos y te decimos si la tenemos.</p>
        <div className="mt-3 grid gap-2 sm:grid-cols-2">
          <BotonCanal canal="santiago" numeros={numeros} demo={demo} tam="medio" texto="Santiago" />
          <BotonCanal canal="regiones" numeros={numeros} demo={demo} tam="medio" texto="Regiones" />
        </div>
      </div>
    </Seccion>
  );
}
