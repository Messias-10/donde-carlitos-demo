"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useRef, useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import type { Categoria, ProductoAdmin, Temporada } from "@/lib/tipos";
import { IconoFlechaIzq } from "../Iconos";
import { MAX_BYTES_ORIGINAL, prepararFoto, subirFoto } from "./foto";
import { useManejarError } from "./Sesion";
import { useApi } from "./useApi";

const TEMPORADAS: { v: Temporada | ""; t: string }[] = [
  { v: "", t: "Sin temporada" },
  { v: "navidad", t: "Navidad" },
  { v: "escolar", t: "Escolar" },
  { v: "fiestas_patrias", t: "Fiestas Patrias" },
  { v: "halloween", t: "Halloween" },
  { v: "invierno", t: "Invierno" },
  { v: "verano", t: "Verano" },
];

interface Campos {
  nombre: string;
  categoria_id: string;
  descripcion: string;
  precio: string;
  precio_oferta: string;
  etiqueta_oferta: "oferta" | "remate";
  precio_mayor: string;
  mayor_minimo: string;
  mayor_minimo_confirmado: boolean;
  unidad_mayor: string;
  oferta_hasta: string;
  stock: string;
  umbral_pocos: string;
  temporada: Temporada | "";
  fecha_llegada: string;
  tiktok_url: string;
  activo: boolean;
}

const VACIO: Campos = {
  nombre: "",
  categoria_id: "",
  descripcion: "",
  precio: "",
  precio_oferta: "",
  etiqueta_oferta: "oferta",
  precio_mayor: "",
  mayor_minimo: "",
  // Si el dueño escribe el mínimo, es un dato real.
  mayor_minimo_confirmado: true,
  unidad_mayor: "",
  oferta_hasta: "",
  stock: "",
  umbral_pocos: "5",
  temporada: "",
  fecha_llegada: "",
  tiktok_url: "",
  activo: true,
};

function desde(p: ProductoAdmin): Campos {
  return {
    nombre: p.nombre,
    categoria_id: String(p.categoria.id),
    descripcion: p.descripcion,
    precio: String(p.precio),
    precio_oferta: p.precio_oferta ? String(p.precio_oferta) : "",
    etiqueta_oferta: p.etiqueta_oferta === "remate" ? "remate" : "oferta",
    precio_mayor: p.precio_mayor ? String(p.precio_mayor) : "",
    mayor_minimo: p.mayor_minimo ? String(p.mayor_minimo) : "",
    mayor_minimo_confirmado: p.mayor_minimo_confirmado === true,
    unidad_mayor: p.unidad_mayor ?? "",
    oferta_hasta: p.oferta_hasta ?? "",
    stock: String(p.stock),
    umbral_pocos: String(p.umbral_pocos),
    temporada: p.temporada ?? "",
    fecha_llegada: p.fecha_llegada ?? "",
    tiktok_url: p.tiktok_url ?? "",
    activo: p.activo,
  };
}

const entero = (s: string): number | null => (/^\d{1,8}$/.test(s.trim()) ? Number(s.trim()) : null);

/** Mensajes de copy-landing.md §16 (calzan con las reglas de la API). */
function validar(c: Campos, nuevo: boolean): Record<string, string> {
  const e: Record<string, string> = {};
  if (nuevo) {
    const n = c.nombre.trim();
    if (n.length < 2 || n.length > 120) e.nombre = "Escribe el nombre del producto (2 a 120 letras).";
    if (!c.categoria_id) e.categoria_id = "Elige la categoría.";
  }
  const precio = entero(c.precio);
  if (!precio || precio < 1) e.precio = "Escribe el precio detalle en pesos (solo números).";
  if (c.precio_oferta.trim()) {
    const o = entero(c.precio_oferta);
    if (!o || !precio || o >= precio) e.precio_oferta = "El precio de oferta tiene que ser menor que el precio detalle.";
  }
  const pm = c.precio_mayor.trim();
  const mm = c.mayor_minimo.trim();
  if (pm || mm) {
    const m = entero(pm);
    if (!m || !precio || m >= precio) e.precio_mayor = "El precio por mayor tiene que ser menor que el precio detalle.";
    const n = entero(mm);
    if (!n || n < 2) e.mayor_minimo = "Si pones precio por mayor, indica desde cuántas unidades (2 o más).";
  } else if (c.unidad_mayor.trim()) {
    e.unidad_mayor = "La unidad solo se usa junto con el precio por mayor.";
  }
  if (entero(c.stock) === null) e.stock = "Escribe el stock (0 o más).";
  const u = entero(c.umbral_pocos);
  if (u === null || u > 1000) e.umbral_pocos = "Escribe un número entre 0 y 1000.";
  if (c.tiktok_url.trim() && !/^https:\/\/(www\.|vm\.|m\.)?tiktok\.com\//.test(c.tiktok_url.trim())) {
    e.tiktok_url = "Pega un link de tiktok.com que empiece con https://";
  }
  if (c.oferta_hasta && !c.precio_oferta.trim()) e.oferta_hasta = "La fecha de término solo se usa con precio de oferta.";
  if (c.descripcion.length > 2000) e.descripcion = "Máximo 2000 caracteres.";
  return e;
}

/** Cuerpo para POST/PUT (PUT reemplaza todo: se reenvía lo que no se edita). */
function cuerpo(c: Campos, base: ProductoAdmin | null) {
  const conMayor = c.precio_mayor.trim() !== "";
  return {
    slug: base?.slug ?? "",
    codigo: base?.codigo ?? codigoNuevo(),
    nombre: c.nombre.trim(),
    descripcion: c.descripcion.trim(),
    categoria_id: Number(c.categoria_id),
    precio: entero(c.precio),
    precio_oferta: c.precio_oferta.trim() ? entero(c.precio_oferta) : null,
    etiqueta_oferta: c.etiqueta_oferta,
    precio_mayor: conMayor ? entero(c.precio_mayor) : null,
    mayor_minimo: conMayor ? entero(c.mayor_minimo) : null,
    mayor_minimo_confirmado: conMayor ? c.mayor_minimo_confirmado : false,
    unidad_mayor: conMayor ? c.unidad_mayor.trim() : "",
    oferta_hasta: c.precio_oferta.trim() && c.oferta_hasta ? c.oferta_hasta : null,
    stock: entero(c.stock),
    umbral_pocos: entero(c.umbral_pocos),
    activo: c.activo,
    temporada: c.temporada || null,
    fecha_llegada: c.fecha_llegada || null,
    imagen_url: base?.imagen_url ?? null,
    tiktok_url: c.tiktok_url.trim() || null,
    fuente_url: base?.fuente?.url ?? null,
    fuente_fecha: base?.fuente?.fecha ?? null,
    es_ejemplo: base?.es_ejemplo ?? false,
  };
}

/** Código interno para productos nuevos: "N" + base 36 del reloj (el dueño no lo escribe). */
function codigoNuevo(): string {
  return ("N" + Date.now().toString(36)).toUpperCase().slice(0, 20);
}

function hoy(): string {
  return new Intl.DateTimeFormat("en-CA", { timeZone: "America/Santiago" }).format(new Date());
}

// ---------------------------------------------------------------- piezas

const claseInput =
  "mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px] focus:border-negro aria-[invalid=true]:border-rojo";

function Campo({
  k,
  etiqueta,
  c,
  set,
  errores,
  extra = {},
}: {
  k: keyof Campos;
  etiqueta: string;
  c: Campos;
  set: (k: keyof Campos, v: string) => void;
  errores: Record<string, string>;
  extra?: React.InputHTMLAttributes<HTMLInputElement>;
}) {
  return (
    <div>
      <label htmlFor={`p-${k}`} className="text-[14px] font-bold">
        {etiqueta}
      </label>
      <input
        id={`p-${k}`}
        value={String(c[k])}
        onChange={(e) => set(k, e.target.value)}
        aria-invalid={errores[k] ? true : undefined}
        aria-describedby={errores[k] ? `pe-${k}` : undefined}
        className={claseInput}
        {...extra}
      />
      {errores[k] ? (
        <p id={`pe-${k}`} className="mt-1 text-[13px] font-semibold text-rojo">
          {errores[k]}
        </p>
      ) : null}
    </div>
  );
}

const NUM = { inputMode: "numeric" as const, pattern: "[0-9]*", maxLength: 8 };

function SelectorFoto({ onFoto, actual }: { onFoto: (f: File | null) => void; actual: string | null }) {
  const [vista, setVista] = useState<string | null>(null);
  const [aviso, setAviso] = useState("");
  const input = useRef<HTMLInputElement>(null);
  return (
    <div>
      <p className="text-[14px] font-bold">Foto</p>
      <div className="mt-1 flex items-center gap-3">
        <div className="grid size-20 shrink-0 place-items-center overflow-hidden rounded-[10px] bg-negro text-[11px] text-gris-claro">
          {vista || actual ? (
            // Vista previa local (blob:) o la foto actual del sitio: no pasa por next/image.
            // eslint-disable-next-line @next/next/no-img-element
            <img src={vista ?? actual ?? ""} alt="Vista previa de la foto" className="size-full object-cover" />
          ) : (
            "Sin foto"
          )}
        </div>
        <div className="flex flex-col gap-2">
          <button
            type="button"
            onClick={() => input.current?.click()}
            className="h-11 rounded-[10px] bg-negro px-4 text-[15px] font-bold text-oro"
          >
            {vista || actual ? "Cambiar foto" : "Tomar o elegir foto"}
          </button>
          {vista ? (
            <button
              type="button"
              onClick={() => {
                setVista(null);
                onFoto(null);
                if (input.current) input.current.value = "";
              }}
              className="min-h-11 text-left text-[14px] font-bold text-rojo underline underline-offset-2"
            >
              Quitar foto nueva
            </button>
          ) : null}
        </div>
      </div>
      <input
        ref={input}
        type="file"
        accept="image/*"
        capture="environment"
        className="sr-only"
        tabIndex={-1}
        aria-label="Foto del producto"
        onChange={(e) => {
          const f = e.target.files?.[0] ?? null;
          if (f && f.size > MAX_BYTES_ORIGINAL) {
            setAviso("La foto es muy pesada. Prueba con otra.");
            return;
          }
          setAviso("");
          if (vista) URL.revokeObjectURL(vista);
          setVista(f ? URL.createObjectURL(f) : null);
          onFoto(f);
        }}
      />
      {aviso ? <p className="mt-1 text-[13px] font-semibold text-rojo">{aviso}</p> : null}
    </div>
  );
}

function MasOpciones({
  c,
  set,
  errores,
  categorias,
  conCategoria,
  abierto,
}: {
  c: Campos;
  set: (k: keyof Campos, v: string | boolean) => void;
  errores: Record<string, string>;
  categorias: Categoria[];
  conCategoria: boolean;
  abierto: boolean;
}) {
  const s = (k: keyof Campos, v: string) => set(k, v);
  return (
    <details className="group rounded-[12px] bg-blanco p-4" open={abierto}>
      <summary className="flex min-h-11 items-center justify-between text-[16px] font-extrabold">
        Más opciones
        <span className="text-[22px] leading-none transition-transform group-open:rotate-45" aria-hidden>
          +
        </span>
      </summary>
      <div className="mt-3 space-y-3">
        {conCategoria ? (
          <div>
            <label htmlFor="p-categoria_id" className="text-[14px] font-bold">
              Categoría
            </label>
            <select id="p-categoria_id" value={c.categoria_id} onChange={(e) => s("categoria_id", e.target.value)} className={claseInput}>
              {categorias.map((k) => (
                <option key={k.id} value={k.id}>
                  {k.nombre}
                </option>
              ))}
            </select>
          </div>
        ) : null}
        <Campo k="precio_oferta" etiqueta="Precio oferta (opcional)" c={c} set={s} errores={errores} extra={NUM} />
        <div>
          <label htmlFor="p-etiqueta_oferta" className="text-[14px] font-bold">
            ¿Oferta o remate?
          </label>
          <select
            id="p-etiqueta_oferta"
            value={c.etiqueta_oferta}
            onChange={(e) => set("etiqueta_oferta", e.target.value === "remate" ? "remate" : "oferta")}
            className={claseInput}
          >
            <option value="oferta">Oferta</option>
            <option value="remate">Remate</option>
          </select>
        </div>
        <Campo k="oferta_hasta" etiqueta="Oferta válida hasta (opcional)" c={c} set={s} errores={errores} extra={{ type: "date" }} />
        <Campo k="unidad_mayor" etiqueta="Unidad (ej: caja de 12)" c={c} set={s} errores={errores} extra={{ maxLength: 60 }} />
        <Campo k="umbral_pocos" etiqueta='Mostrar "quedan pocas" con esta cantidad o menos' c={c} set={s} errores={errores} extra={NUM} />
        <div>
          <label htmlFor="p-temporada" className="text-[14px] font-bold">
            Temporada
          </label>
          <select id="p-temporada" value={c.temporada} onChange={(e) => s("temporada", e.target.value)} className={claseInput}>
            {TEMPORADAS.map((t) => (
              <option key={t.v} value={t.v}>
                {t.t}
              </option>
            ))}
          </select>
        </div>
        <Campo k="fecha_llegada" etiqueta="Fecha de llegada" c={c} set={s} errores={errores} extra={{ type: "date" }} />
        <div>
          <label htmlFor="p-descripcion" className="text-[14px] font-bold">
            Descripción (tallas, colores, qué trae)
          </label>
          <textarea
            id="p-descripcion"
            value={c.descripcion}
            onChange={(e) => s("descripcion", e.target.value)}
            maxLength={2000}
            className={`${claseInput} h-28 py-2`}
          />
        </div>
        <label className="flex min-h-12 items-center gap-3 text-[15px] font-bold">
          <input type="checkbox" checked={c.activo} onChange={(e) => set("activo", e.target.checked)} className="size-6 accent-negro" />
          Visible en la web
        </label>
      </div>
    </details>
  );
}

// ---------------------------------------------------------------- formulario

function Formulario({
  base,
  categorias,
  guardado,
  aviso = "",
  onGuardado,
}: {
  base: ProductoAdmin | null;
  categorias: Categoria[];
  guardado: boolean;
  /** Mensaje que sobrevive al recargar el formulario (p. ej. "la foto no se subió"). */
  aviso?: string;
  onGuardado: (p: ProductoAdmin, aviso: string) => void;
}) {
  const nuevo = base === null;
  const [c, setC] = useState<Campos>(() =>
    base ? desde(base) : { ...VACIO, fecha_llegada: hoy(), categoria_id: categorias[0] ? String(categorias[0].id) : "" },
  );
  const [foto, setFoto] = useState<File | null>(null);
  const [errores, setErrores] = useState<Record<string, string>>({});
  const [error, setError] = useState(aviso);
  const [estado, setEstado] = useState<"" | "guardando" | "guardado">(guardado ? "guardado" : "");
  const manejar = useManejarError();

  const set = (k: keyof Campos, v: string | boolean) => {
    setC((x) => ({ ...x, [k]: v }));
    setEstado("");
    if (errores[k]) setErrores((e) => ({ ...e, [k]: "" }));
  };
  const s = (k: keyof Campos, v: string) => set(k, v);
  const opcionesConError = ["precio_oferta", "oferta_hasta", "unidad_mayor", "umbral_pocos", "descripcion"].some((k) => errores[k]);

  async function guardar(ev: React.FormEvent) {
    ev.preventDefault();
    const e = validar(c, nuevo);
    if (Object.keys(e).length) {
      setErrores(e);
      requestAnimationFrame(() => document.getElementById(`p-${Object.keys(e)[0]}`)?.focus());
      return;
    }
    setEstado("guardando");
    setError("");
    try {
      let p = nuevo
        ? await api<ProductoAdmin>("/api/admin/productos", { method: "POST", body: cuerpo(c, null) })
        : await api<ProductoAdmin>(`/api/admin/productos/${base.id}`, { method: "PUT", body: cuerpo(c, base) });
      let avisoFoto = "";
      if (foto) {
        try {
          const conFoto = await subirFoto(p.id, await prepararFoto(foto));
          if (conFoto) p = conFoto;
        } catch (err) {
          // El producto ya quedó guardado; solo falló la foto.
          avisoFoto = `El producto se guardó, pero la foto no se pudo subir (${manejar(err)}). Prueba de nuevo con "Cambiar foto".`;
        }
      }
      onGuardado(p, avisoFoto);
    } catch (err) {
      setEstado("");
      if (err instanceof ErrorCliente && err.campos) setErrores(err.campos);
      setError(
        err instanceof ErrorCliente && err.status === 409 && err.codigo === "conflicto"
          ? "Alguien más cambió este producto. Recarga y vuelve a intentar."
          : manejar(err),
      );
    }
  }

  return (
    <form onSubmit={guardar} noValidate className="space-y-4">
      <Link href="/admin/productos" className="inline-flex min-h-11 items-center gap-1 text-[14px] font-bold">
        <IconoFlechaIzq size={18} /> Productos
      </Link>

      {nuevo ? (
        <h1 className="titular text-[28px]">Agregar producto</h1>
      ) : (
        <div className="rounded-[12px] bg-blanco p-4">
          <h1 className="text-[20px] font-extrabold leading-snug">{base.nombre}</h1>
          <p className="text-[13px] text-gris">
            {base.codigo} · {base.categoria.nombre}
            {base.es_ejemplo ? " · producto de ejemplo" : ""} ·{" "}
            <Link href={`/producto/${base.slug}`} className="inline-flex min-h-11 items-center underline" target="_blank">
              ver en la tienda
            </Link>
          </p>
        </div>
      )}

      <fieldset className="space-y-3 rounded-[12px] bg-blanco p-4">
        <legend className="sr-only">Lo principal</legend>
        {nuevo ? (
          <>
            <Campo k="nombre" etiqueta="Nombre del producto" c={c} set={s} errores={errores} extra={{ maxLength: 120 }} />
            <div>
              <label htmlFor="p-categoria_id" className="text-[14px] font-bold">
                Categoría
              </label>
              <select
                id="p-categoria_id"
                value={c.categoria_id}
                onChange={(e) => s("categoria_id", e.target.value)}
                aria-invalid={errores.categoria_id ? true : undefined}
                className={claseInput}
              >
                <option value="">Elige la categoría</option>
                {categorias.map((k) => (
                  <option key={k.id} value={k.id}>
                    {k.nombre}
                  </option>
                ))}
              </select>
              {errores.categoria_id ? <p className="mt-1 text-[13px] font-semibold text-rojo">{errores.categoria_id}</p> : null}
            </div>
          </>
        ) : null}
        <div className="grid grid-cols-2 gap-3">
          <Campo k="precio" etiqueta="Precio detalle" c={c} set={s} errores={errores} extra={NUM} />
          <Campo k="stock" etiqueta="Stock" c={c} set={s} errores={errores} extra={NUM} />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <Campo k="precio_mayor" etiqueta="Precio por mayor" c={c} set={s} errores={errores} extra={NUM} />
          <Campo k="mayor_minimo" etiqueta="Desde cuántas" c={c} set={s} errores={errores} extra={NUM} />
        </div>
        {c.precio_mayor.trim() ? (
          <label className="flex min-h-11 items-center gap-3 text-[15px] font-semibold">
            <input
              type="checkbox"
              checked={c.mayor_minimo_confirmado}
              onChange={(e) => set("mayor_minimo_confirmado", e.target.checked)}
              className="size-6 accent-negro"
            />
            El mínimo por mayor está confirmado
          </label>
        ) : null}
        <Campo
          k="tiktok_url"
          etiqueta="Link al video de TikTok"
          c={c}
          set={s}
          errores={errores}
          extra={{ type: "url", maxLength: 500, inputMode: "url", placeholder: "https://www.tiktok.com/@mayte_emporio/video/…" }}
        />
        <SelectorFoto onFoto={setFoto} actual={base?.imagen_url ?? null} />
      </fieldset>

      <MasOpciones c={c} set={set} errores={errores} categorias={categorias} conCategoria={!nuevo} abierto={opcionesConError} />

      <div aria-live="polite">
        {error ? (
          <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 text-[15px] font-semibold text-rojo">
            {error}
          </p>
        ) : null}
      </div>

      <button
        type="submit"
        disabled={estado === "guardando"}
        className="sticky bottom-3 flex h-[52px] w-full items-center justify-center rounded-[14px] bg-negro text-[17px] font-extrabold text-oro shadow-lg disabled:opacity-60"
      >
        {estado === "guardando" ? "Guardando…" : estado === "guardado" ? "Guardado" : nuevo ? "Agregar producto" : "Guardar"}
      </button>
    </form>
  );
}

// ---------------------------------------------------------------- páginas

const CLAVE_AVISO = "dc_admin_aviso";

function avisoInicial(): string {
  try {
    const a = window.sessionStorage.getItem(CLAVE_AVISO) ?? "";
    window.sessionStorage.removeItem(CLAVE_AVISO);
    return a;
  } catch {
    return "";
  }
}

function Cargando({ error }: { error?: string }) {
  return error ? (
    <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 font-semibold text-rojo">
      {error}
    </p>
  ) : (
    <div className="h-80 animate-pulse rounded-[12px] bg-blanco/70" aria-hidden />
  );
}

export function FormProducto({ id }: { id: string }) {
  const prod = useApi<ProductoAdmin>(`/api/admin/productos/${encodeURIComponent(id)}`);
  const cats = useApi<{ items: Categoria[] }>("/api/categorias");
  const [guardadoEn, setGuardadoEn] = useState("");
  const [aviso, setAviso] = useState(() => avisoInicial());
  if (!prod.data || !cats.data) return <Cargando error={prod.error ?? cats.error} />;
  return (
    <Formulario
      key={`${prod.data.actualizado_en}|${prod.data.imagen_url}|${aviso}`}
      base={prod.data}
      categorias={cats.data.items}
      guardado={guardadoEn === prod.data.actualizado_en && !aviso}
      aviso={aviso}
      onGuardado={(act, a) => {
        setGuardadoEn(act.actualizado_en);
        setAviso(a);
        prod.reemplazar(act);
      }}
    />
  );
}

export function NuevoProducto() {
  const cats = useApi<{ items: Categoria[] }>("/api/categorias");
  const router = useRouter();
  if (!cats.data) return <Cargando error={cats.error} />;
  return (
    <Formulario
      base={null}
      categorias={cats.data.items}
      guardado={false}
      onGuardado={(p, aviso) => {
        // El aviso de la foto viaja en sessionStorage (no en la URL) hasta el formulario de edición.
        try {
          if (aviso) window.sessionStorage.setItem(CLAVE_AVISO, aviso);
        } catch {
          /* sin almacenamiento: se pierde el aviso, el producto igual quedó guardado */
        }
        router.replace(`/admin/productos/${p.id}`);
      }}
    />
  );
}
