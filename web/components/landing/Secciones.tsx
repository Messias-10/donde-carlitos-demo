import Link from "next/link";
import { BotonAgregarRapido } from "@/components/BotonAgregarRapido";
import { IconoCaja, IconoCategoria, IconoTienda } from "@/components/Iconos";
import { PorConfirmar } from "@/components/Marcas";
import { Precio } from "@/components/Precio";
import { FilaProductos, ItemFila, Seccion } from "@/components/Seccion";
import { paraCarrito, TarjetaProducto } from "@/components/TarjetaProducto";
import { VideoTikTok } from "@/components/VideoTikTok";
import { clp } from "@/lib/formato";
import type { Categoria, ProductoVista } from "@/lib/tipos";

// ---------------------------------------------------------------- TikTok

export function LoVisteEnTikTok({ productos, conDireccionVieja }: { productos: ProductoVista[]; conDireccionVieja: Set<string> }) {
  if (productos.length === 0) return null;
  return (
    <Seccion
      id="tiktok"
      fondo="negro"
      titulo="¿Lo viste en TikTok?"
      bajada="Acá lo tienes con precio y stock. Toca el video para verlo de nuevo o agrégalo al tiro a tu pedido."
    >
      <div className="scroll-fila -mx-4 flex gap-4 overflow-x-auto px-4 pb-2 lg:mx-0 lg:grid lg:grid-cols-4 lg:overflow-visible lg:px-0" role="list" aria-label="Productos de nuestros videos">
        {productos.map((p) => (
          <div key={p.id} role="listitem" className="flex w-[44vw] max-w-[200px] shrink-0 flex-col gap-3 lg:w-auto lg:max-w-none">
            <VideoTikTok
              url={p.tiktok_url!}
              titulo={p.nombre}
              poster={p.imagen}
              avisoDireccion={conDireccionVieja.has(p.tiktok_url!)}
            />
            <div className="flex flex-1 flex-col gap-2 rounded-[12px] bg-carbon p-3">
              <h3 className="line-clamp-2 text-[14px] font-semibold leading-snug text-blanco">
                <Link href={`/producto/${p.slug}`} className="hover:underline">
                  {p.nombre}
                </Link>
              </h3>
              <Precio precio={p.precio} precioFinal={p.precio_final} tam="sm" oscuro />
              <div className="mt-auto">
                <BotonAgregarRapido producto={paraCarrito(p)} agotado={p.disponibilidad === "agotado"} promo={false} />
              </div>
            </div>
          </div>
        ))}
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Categorías

export function Categorias({ categorias }: { categorias: Categoria[] }) {
  const visibles = categorias.filter((c) => c.total_productos > 0);
  if (visibles.length === 0) return null;
  return (
    <Seccion id="categorias" titulo="¿Qué andas buscando?">
      <ul className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        {visibles.map((c) => (
          <li key={c.id}>
            <Link
              href={`/catalogo?categoria=${encodeURIComponent(c.slug)}`}
              className="flex h-full items-center gap-3 rounded-[12px] bg-blanco p-3 shadow-[0_1px_0_rgba(11,11,11,0.08)] ring-1 ring-negro/5 hover:ring-oro"
            >
              <span className="grid size-12 shrink-0 place-items-center rounded-full bg-negro text-oro">
                <IconoCategoria slug={c.slug} size={26} />
              </span>
              <span className="min-w-0">
                <span className="block text-[15px] font-bold leading-tight">{c.nombre}</span>
                <span className="text-[13px] text-gris">
                  {c.total_productos} {c.total_productos === 1 ? "producto" : "productos"}
                </span>
              </span>
            </Link>
          </li>
        ))}
      </ul>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Recién llegado

export function RecienLlegado({ productos }: { productos: ProductoVista[] }) {
  if (productos.length === 0) return null;
  return (
    <Seccion
      id="nuevo"
      fondo="blanco"
      titulo={
        <>
          ¡Nos llegó! <span aria-hidden>✨</span>
        </>
      }
      bajada="Lo último que bajó del contenedor. Cuando se acaba, se acaba."
    >
      <FilaProductos etiqueta="Recién llegados">
        {productos.map((p) => (
          <ItemFila key={p.id}>
            <TarjetaProducto p={p} mostrarLlegada sizes="(min-width: 1024px) 270px, 60vw" />
          </ItemFila>
        ))}
      </FilaProductos>
      <Link
        href="/catalogo?nuevo=1"
        className="mt-5 inline-flex h-12 items-center rounded-[12px] border-2 border-negro px-5 text-[15px] font-extrabold"
      >
        Ver todo lo nuevo
      </Link>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Pa negocio

export function PaNegocio({ ejemplo }: { ejemplo: ProductoVista | null }) {
  const detalle = ejemplo?.precio_final ?? 4000;
  const mayor = ejemplo?.precio_mayor ?? 3500;
  const unidad = ejemplo?.unidad_mayor ?? "caja de 12";
  const nombre = ejemplo?.nombre ?? "Squishy viral que cambia de color con el sol";
  return (
    <Seccion
      id="pa-negocio"
      titulo="¿Pa negocio?"
      bajada="Cada producto trae su precio al detalle y su precio por mayor, con el mínimo a la vista."
    >
      <div className="grid gap-6 lg:grid-cols-[1.1fr_1fr] lg:items-center">
        {/* Escalera de precio: detalle -> por mayor. */}
        <div className="rounded-[16px] bg-blanco p-5 ring-1 ring-negro/10">
          <p className="text-[14px] font-semibold text-gris">Ejemplo real</p>
          <p className="mt-1 text-[17px] font-bold leading-snug">
            {ejemplo ? (
              <Link href={`/producto/${ejemplo.slug}`} className="underline-offset-2 hover:underline">
                {nombre}
              </Link>
            ) : (
              nombre
            )}
          </p>
          <div className="mt-4 grid grid-cols-2 items-end gap-3">
            <div className="rounded-[12px] bg-crema p-3">
              <p className="text-[13px] font-semibold text-gris">Al detalle</p>
              <p className="mt-1 font-display text-[30px] leading-none">{clp(detalle)}</p>
            </div>
            <div className="rounded-[12px] bg-negro p-3 pt-6 text-blanco">
              <p className="text-[13px] font-semibold text-gris-claro">Por {unidad}</p>
              <p className="mt-1">
                <span className="placa-precio text-[30px]">{clp(mayor)}</span>
              </p>
              <p className="mt-1 text-[12px] text-gris-claro">c/u</p>
            </div>
          </div>
        </div>
        <div>
          <ul className="space-y-3 text-[16px]">
            <li className="flex gap-3">
              <span className="mt-2 size-2 shrink-0 rounded-full bg-oro" aria-hidden />
              En el carrito, el precio por mayor se aplica solo cuando llegas al mínimo.
            </li>
            <li className="flex gap-3">
              <span className="mt-2 size-2 shrink-0 rounded-full bg-oro" aria-hidden />
              ¿Necesitas más cantidad? Escríbelo en el comentario del pedido y te cotizamos.
            </li>
          </ul>
          <Link
            href="/catalogo"
            className="mt-5 inline-flex h-12 items-center rounded-[12px] bg-negro px-5 text-[15px] font-extrabold text-oro"
          >
            Ver precios por mayor
          </Link>
        </div>
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Pasos

export function PasosPedido() {
  const pasos = [
    {
      t: "Elige",
      d: "Agrega lo que quieres al carrito. Si llevas cantidad, el precio por mayor se aplica solo.",
    },
    {
      t: "Envía",
      d: "Deja tu nombre, tu WhatsApp y si retiras o te lo enviamos. Tu pedido queda con número y se abre WhatsApp con todo escrito: solo aprietas enviar.",
    },
    {
      t: "Te confirmamos",
      d: "Te respondemos desde nuestro número oficial con el total y el pago. Después retiras en Alameda 2432 o lo despachamos.",
    },
  ];
  return (
    <Seccion
      id="como-pedir"
      fondo="blanco"
      titulo={
        <>
          Pedir es así de fácil <span aria-hidden>👇</span>
        </>
      }
    >
      <ol className="grid gap-4 md:grid-cols-3">
        {pasos.map((p, i) => (
          <li key={p.t} className="flex gap-4 rounded-[12px] bg-crema p-4">
            <span className="font-display text-[44px] leading-none text-oro-profundo" aria-hidden>
              {i + 1}
            </span>
            <div>
              <h3 className="titular text-[22px]">{p.t}</h3>
              <p className="mt-1 text-[15px] text-negro/85">{p.d}</p>
            </div>
          </li>
        ))}
      </ol>
      <p className="mt-4 text-[14px] text-gris">
        La web no cobra nada. El pago lo coordinamos contigo por WhatsApp o en el local.
        <PorConfirmar texto="despacho por confirmar" />
      </p>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Retiro y envíos

export function RetiroEnvios() {
  return (
    <Seccion id="envios" titulo="Retira o te lo enviamos">
      <div className="grid gap-4 md:grid-cols-2">
        <div className="rounded-[12px] bg-blanco p-5 ring-1 ring-negro/10">
          <IconoTienda size={32} className="text-oro-profundo" />
          <h3 className="titular mt-3 text-[24px]">Retiro en tienda</h3>
          <p className="mt-2 text-[15px]">
            Haz tu pedido, espera nuestra confirmación por WhatsApp y pasa a buscarlo a Alameda 2432 con tu número de pedido.
          </p>
          <p className="mt-2 text-[15px] font-semibold">Llegas sabiendo qué llevas y cuánto pagas.</p>
        </div>
        <div className="rounded-[12px] bg-blanco p-5 ring-1 ring-negro/10">
          <IconoCaja size={32} className="text-oro-profundo" />
          <h3 className="titular mt-3 text-[24px]">Envío a regiones</h3>
          <p className="mt-2 text-[15px]">
            ¿Eres de región? Arma tu pedido, elige tu agencia y te confirmamos por WhatsApp el total y el despacho.
          </p>
          <p className="mt-3">
            <PorConfirmar texto="agencias, quién paga el envío y monto mínimo por confirmar" />
          </p>
        </div>
      </div>
    </Seccion>
  );
}
