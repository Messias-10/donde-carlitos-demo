import Image from "next/image";
import Link from "next/link";
import { IconoWhatsApp } from "@/components/Iconos";
import { FilaProductos, ItemFila, Seccion } from "@/components/Seccion";
import { TarjetaProducto } from "@/components/TarjetaProducto";
import { miles } from "@/lib/formato";
import type { ProductoVista } from "@/lib/tipos";
import { linkWhatsApp, MSJ_CONSULTA } from "@/lib/whatsapp";

/** Banner único de temporada (estilo-carlitos §8.3) + fila de ofertas. */
export function Ofertas({ productos, esfera }: { productos: ProductoVista[]; esfera: ProductoVista | null }) {
  const precioEsfera = esfera?.precio_final ?? 500;
  return (
    <Seccion
      id="ofertas"
      titulo={
        <>
          Ofertas que vuelan <span aria-hidden>🔥</span>
        </>
      }
      bajada="Remates y promos de la semana. Cuando se acaban, se acaban."
    >
      <div className="relative mb-8 overflow-hidden rounded-[16px] bg-negro text-blanco ring-1 ring-oro lg:grid lg:grid-cols-[1.3fr_1fr]">
        <div className="relative z-10 p-5 lg:p-8">
          <p className="font-script text-[20px] text-blanco">Donde</p>
          <p className="titular texto-oro text-[40px] lg:text-[56px]">¡Remate de Navidad!</p>
          <p className="mt-3 max-w-[40ch] text-[16px] text-blanco/90">
            Esferas con luces, telones, trenes, lazos y bolsas de regalo.
          </p>
          <p className="mt-4 flex flex-wrap items-center gap-3 text-[16px]">
            Esferas con luces a
            <span className="placa-precio text-[34px]" data-oferta="true">
              <span className="signo">$</span>
              {miles(precioEsfera)}
            </span>
          </p>
          <p className="mt-4 text-[14px] text-gris-claro">
            <span aria-hidden>📍 </span>Alameda 2432
          </p>
          <div className="mt-5 flex flex-wrap gap-3">
            <Link
              href="/catalogo?temporada=navidad"
              className="inline-flex h-12 items-center rounded-[12px] bg-oro px-5 text-[15px] font-extrabold text-negro"
            >
              Ver ofertas
            </Link>
            <a
              href={linkWhatsApp(MSJ_CONSULTA)}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex h-12 items-center gap-2 rounded-[12px] bg-whatsapp px-5 text-[15px] font-extrabold text-negro"
            >
              <IconoWhatsApp size={20} />
              Escríbenos por WhatsApp
            </a>
          </div>
        </div>
        {esfera?.imagen ? (
          <div className="relative aspect-[16/10] lg:aspect-auto">
            <Image
              src={esfera.imagen}
              alt="Esferas navideñas con luces en su caja"
              fill
              sizes="(min-width: 1024px) 420px, 100vw"
              className="object-cover"
            />
            <span className="absolute inset-0 bg-gradient-to-t from-negro/60 to-transparent lg:bg-gradient-to-r" aria-hidden />
          </div>
        ) : null}
      </div>

      {productos.length > 0 ? (
        <FilaProductos etiqueta="Ofertas y remates">
          {productos.map((p) => (
            <ItemFila key={p.id}>
              <TarjetaProducto p={p} sizes="(min-width: 1024px) 270px, 60vw" />
            </ItemFila>
          ))}
        </FilaProductos>
      ) : (
        <p className="text-gris">Pronto subimos las ofertas de la semana.</p>
      )}
    </Seccion>
  );
}
