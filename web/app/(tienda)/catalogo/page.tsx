import type { Metadata } from "next";
import Link from "next/link";
import { AvisoSinApi } from "@/components/AvisoSinApi";
import { IconoLupa } from "@/components/Iconos";
import { TarjetaProducto } from "@/components/TarjetaProducto";
import { listarCategorias, listarProductos, type FiltroCatalogo } from "@/lib/api-servidor";
import type { Categoria, Pagina, ProductoVista } from "@/lib/tipos";

export const metadata: Metadata = {
  title: "Catálogo",
  description: "Todo lo que tenemos en Alameda 2432, con precio al detalle y por mayor.",
};

type SP = Record<string, string | string[] | undefined>;

const ORDENES = {
  relevancia: "Destacados",
  recientes: "Lo más nuevo",
  precio_asc: "Menor precio",
  precio_desc: "Mayor precio",
} as const;
type Orden = keyof typeof ORDENES;

function uno(v: string | string[] | undefined): string {
  return (Array.isArray(v) ? v[0] : v) ?? "";
}

function leerFiltro(sp: SP) {
  const q = uno(sp.q).trim().slice(0, 60);
  const categoria = /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(uno(sp.categoria)) ? uno(sp.categoria).slice(0, 80) : "";
  const oferta = uno(sp.oferta) === "1";
  const nuevo = uno(sp.nuevo) === "1";
  const temporada = uno(sp.temporada) === "navidad" ? "navidad" : "";
  const o = uno(sp.orden);
  const orden: Orden = o in ORDENES ? (o as Orden) : "relevancia";
  const page = Math.min(Math.max(1, Number.parseInt(uno(sp.page), 10) || 1), 500);
  return { q, categoria, oferta, nuevo, temporada, orden, page };
}

type Filtro = ReturnType<typeof leerFiltro>;

/** Arma la URL del catálogo cambiando solo lo indicado (y volviendo a la página 1). */
function url(f: Filtro, cambios: Partial<Filtro>): string {
  const n = { ...f, page: 1, ...cambios };
  const p = new URLSearchParams();
  if (n.q) p.set("q", n.q);
  if (n.categoria) p.set("categoria", n.categoria);
  if (n.oferta) p.set("oferta", "1");
  if (n.nuevo) p.set("nuevo", "1");
  if (n.temporada) p.set("temporada", n.temporada);
  if (n.orden !== "relevancia") p.set("orden", n.orden);
  if (n.page > 1) p.set("page", String(n.page));
  const s = p.toString();
  return s ? `/catalogo?${s}` : "/catalogo";
}

function Chip({ href, activo, children }: { href: string; activo: boolean; children: React.ReactNode }) {
  return (
    <Link
      href={href}
      aria-current={activo ? "true" : undefined}
      className={`inline-flex h-10 shrink-0 items-center rounded-full px-4 text-[14px] font-bold ${
        activo ? "bg-negro text-oro" : "bg-blanco text-negro ring-1 ring-negro/15"
      }`}
    >
      {children}
    </Link>
  );
}

export default async function Catalogo({ searchParams }: { searchParams: Promise<SP> }) {
  const f = leerFiltro(await searchParams);
  const consulta: FiltroCatalogo = {
    q: f.q || undefined,
    categoria: f.categoria || undefined,
    oferta: f.oferta,
    recien_llegado: f.nuevo,
    temporada: f.temporada || undefined,
    orden: f.orden,
    page: f.page,
    page_size: 24,
  };

  let pagina: Pagina<ProductoVista> | null = null;
  let categorias: Categoria[] = [];
  let error = false;
  try {
    [pagina, categorias] = await Promise.all([listarProductos(consulta), listarCategorias()]);
  } catch {
    error = true;
  }

  const catActual = categorias.find((c) => c.slug === f.categoria);
  const titulo = f.q
    ? `Resultados para "${f.q}"`
    : catActual
      ? catActual.nombre
      : f.oferta
        ? "Ofertas y remates"
        : f.nuevo
          ? "Recién llegado"
          : f.temporada
            ? "Navidad"
            : "Catálogo";

  return (
    <div className="bg-crema pb-16">
      <div className="bg-negro px-4 pb-5 pt-6 text-blanco">
        <div className="mx-auto max-w-6xl">
          <p className="font-script text-[20px]">Todo el local</p>
          <h1 className="titular text-[36px] text-oro lg:text-[48px]">{titulo}</h1>
          <form action="/catalogo" method="get" role="search" className="mt-4 flex gap-2">
            {f.categoria ? <input type="hidden" name="categoria" value={f.categoria} /> : null}
            <label htmlFor="q-catalogo" className="sr-only">
              Buscar productos
            </label>
            <input
              id="q-catalogo"
              name="q"
              type="search"
              defaultValue={f.q}
              maxLength={60}
              placeholder="¿Qué andas buscando? Ej: squishy, tren, telón"
              className="h-12 min-w-0 flex-1 rounded-[10px] bg-blanco px-3 text-[16px] text-negro placeholder:text-gris"
              enterKeyHint="search"
            />
            <button type="submit" className="grid h-12 w-12 place-items-center rounded-[10px] bg-oro text-negro" aria-label="Buscar">
              <IconoLupa />
            </button>
          </form>
        </div>
      </div>

      <div className="mx-auto max-w-6xl px-4">
        <nav aria-label="Filtros" className="scroll-fila -mx-4 flex gap-2 overflow-x-auto px-4 py-4">
          <Chip href={url(f, { categoria: "", oferta: false, nuevo: false, temporada: "" })} activo={!f.categoria && !f.oferta && !f.nuevo && !f.temporada}>
            Todo
          </Chip>
          <Chip href={url(f, { oferta: !f.oferta })} activo={f.oferta}>
            Ofertas
          </Chip>
          <Chip href={url(f, { nuevo: !f.nuevo })} activo={f.nuevo}>
            ¡Llegó!
          </Chip>
          <Chip href={url(f, { temporada: f.temporada ? "" : "navidad" })} activo={!!f.temporada}>
            Navidad
          </Chip>
          {categorias
            .filter((c) => c.total_productos > 0)
            .map((c) => (
              <Chip key={c.id} href={url(f, { categoria: f.categoria === c.slug ? "" : c.slug })} activo={f.categoria === c.slug}>
                {c.nombre}
              </Chip>
            ))}
        </nav>

        {error || !pagina ? (
          <div className="-mx-4">
            <AvisoSinApi />
          </div>
        ) : (
          <>
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <p className="text-[14px] text-gris" aria-live="polite">
                {pagina.total} {pagina.total === 1 ? "producto" : "productos"}
              </p>
              <form action="/catalogo" method="get" className="flex items-center gap-2">
                {f.q ? <input type="hidden" name="q" value={f.q} /> : null}
                {f.categoria ? <input type="hidden" name="categoria" value={f.categoria} /> : null}
                {f.oferta ? <input type="hidden" name="oferta" value="1" /> : null}
                {f.nuevo ? <input type="hidden" name="nuevo" value="1" /> : null}
                {f.temporada ? <input type="hidden" name="temporada" value={f.temporada} /> : null}
                <label htmlFor="orden" className="text-[14px] font-semibold">
                  Ordenar
                </label>
                <select
                  id="orden"
                  name="orden"
                  defaultValue={f.orden}
                  className="h-10 rounded-[10px] bg-blanco px-2 text-[14px] ring-1 ring-negro/15"
                >
                  {Object.entries(ORDENES).map(([k, v]) => (
                    <option key={k} value={k}>
                      {v}
                    </option>
                  ))}
                </select>
                <button type="submit" className="h-10 rounded-[10px] bg-negro px-3 text-[14px] font-bold text-oro">
                  Ver
                </button>
              </form>
            </div>

            {pagina.items.length === 0 ? (
              <div className="rounded-[12px] bg-blanco p-6 text-center">
                <p className="text-[16px] font-bold">No encontramos productos con esos filtros.</p>
                <p className="mt-1 text-[15px] text-gris">Prueba con otra palabra o pregúntanos por WhatsApp: muchas cosas no alcanzan a subir.</p>
                <Link href="/catalogo" className="mt-4 inline-flex h-11 items-center rounded-[10px] bg-negro px-4 font-bold text-oro">
                  Ver todo el catálogo
                </Link>
              </div>
            ) : (
              <ul className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
                {pagina.items.map((p, i) => (
                  <li key={p.id}>
                    <TarjetaProducto p={p} preload={i < 2} mostrarLlegada={f.nuevo} />
                  </li>
                ))}
              </ul>
            )}

            {pagina.total_pages > 1 ? (
              <nav aria-label="Páginas" className="mt-8 flex items-center justify-center gap-3">
                {f.page > 1 ? (
                  <Link href={url(f, { page: f.page - 1 })} className="inline-flex h-11 items-center rounded-[10px] bg-blanco px-4 font-bold ring-1 ring-negro/15">
                    Anterior
                  </Link>
                ) : null}
                <span className="text-[14px] text-gris">
                  Página {pagina.page} de {pagina.total_pages}
                </span>
                {f.page < pagina.total_pages ? (
                  <Link href={url(f, { page: f.page + 1 })} className="inline-flex h-11 items-center rounded-[10px] bg-negro px-4 font-bold text-oro">
                    Siguiente
                  </Link>
                ) : null}
              </nav>
            ) : null}
          </>
        )}
      </div>
    </div>
  );
}
