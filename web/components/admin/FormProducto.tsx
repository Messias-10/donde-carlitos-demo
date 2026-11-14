"use client";

import Link from "next/link";
import { useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import type { ProductoAdmin, Temporada } from "@/lib/tipos";
import { IconoFlechaIzq } from "../Iconos";
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
  precio: string;
  precio_oferta: string;
  etiqueta_oferta: "oferta" | "remate";
  precio_mayor: string;
  mayor_minimo: string;
  unidad_mayor: string;
  stock: string;
  umbral_pocos: string;
  temporada: Temporada | "";
  fecha_llegada: string;
  tiktok_url: string;
  activo: boolean;
}

function desde(p: ProductoAdmin): Campos {
  return {
    precio: String(p.precio),
    precio_oferta: p.precio_oferta ? String(p.precio_oferta) : "",
    etiqueta_oferta: p.etiqueta_oferta === "remate" ? "remate" : "oferta",
    precio_mayor: p.precio_mayor ? String(p.precio_mayor) : "",
    mayor_minimo: p.mayor_minimo ? String(p.mayor_minimo) : "",
    unidad_mayor: p.unidad_mayor ?? "",
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
function validar(c: Campos): Record<string, string> {
  const e: Record<string, string> = {};
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
  const s = entero(c.stock);
  if (s === null) e.stock = "Escribe el stock (0 o más).";
  const u = entero(c.umbral_pocos);
  if (u === null || u > 1000) e.umbral_pocos = "Escribe un número entre 0 y 1000.";
  if (c.tiktok_url.trim() && !/^https:\/\/(www\.|vm\.|m\.)?tiktok\.com\//.test(c.tiktok_url.trim())) {
    e.tiktok_url = "Pega un link de tiktok.com que empiece con https://";
  }
  return e;
}

export function FormProducto({ id }: { id: string }) {
  const { data, error, reemplazar } = useApi<ProductoAdmin>(`/api/admin/productos/${encodeURIComponent(id)}`);
  const [guardadoEn, setGuardadoEn] = useState("");
  if (!data) {
    return error ? (
      <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 font-semibold text-rojo">
        {error}
      </p>
    ) : (
      <div className="h-80 animate-pulse rounded-[12px] bg-blanco/70" aria-hidden />
    );
  }
  return (
    <Editor
      key={data.actualizado_en}
      p={data}
      guardado={guardadoEn === data.actualizado_en}
      onGuardado={(act) => {
        setGuardadoEn(act.actualizado_en);
        reemplazar(act);
      }}
    />
  );
}

function Editor({ p, guardado, onGuardado }: { p: ProductoAdmin; guardado: boolean; onGuardado: (p: ProductoAdmin) => void }) {
  const [c, setC] = useState<Campos>(() => desde(p));
  const [errores, setErrores] = useState<Record<string, string>>({});
  const [error, setError] = useState("");
  const [estado, setEstado] = useState<"" | "guardando" | "guardado">(guardado ? "guardado" : "");
  const manejar = useManejarError();

  const set = <K extends keyof Campos>(k: K, v: Campos[K]) => {
    setC((x) => ({ ...x, [k]: v }));
    setEstado("");
    if (errores[k]) setErrores((e) => ({ ...e, [k]: "" }));
  };

  async function guardar(ev: React.FormEvent) {
    ev.preventDefault();
    const e = validar(c);
    if (Object.keys(e).length) {
      setErrores(e);
      document.getElementById(`p-${Object.keys(e)[0]}`)?.focus();
      return;
    }
    const conMayor = c.precio_mayor.trim() !== "";
    // PUT reemplaza el producto completo: se reenvía todo lo que no se edita aquí.
    const cuerpo = {
      slug: p.slug,
      codigo: p.codigo,
      nombre: p.nombre,
      descripcion: p.descripcion,
      categoria_id: p.categoria.id,
      precio: entero(c.precio),
      precio_oferta: c.precio_oferta.trim() ? entero(c.precio_oferta) : null,
      etiqueta_oferta: c.etiqueta_oferta,
      precio_mayor: conMayor ? entero(c.precio_mayor) : null,
      mayor_minimo: conMayor ? entero(c.mayor_minimo) : null,
      unidad_mayor: conMayor ? c.unidad_mayor.trim() : "",
      stock: entero(c.stock),
      umbral_pocos: entero(c.umbral_pocos),
      activo: c.activo,
      temporada: c.temporada || null,
      fecha_llegada: c.fecha_llegada || null,
      imagen_url: p.imagen_url,
      tiktok_url: c.tiktok_url.trim() || null,
      fuente_url: p.fuente?.url ?? null,
      fuente_fecha: p.fuente?.fecha ?? null,
      es_ejemplo: p.es_ejemplo,
    };
    setEstado("guardando");
    setError("");
    try {
      const act = await api<ProductoAdmin>(`/api/admin/productos/${p.id}`, { method: "PUT", body: cuerpo });
      onGuardado(act);
    } catch (err) {
      setEstado("");
      if (err instanceof ErrorCliente && err.campos) setErrores(err.campos);
      setError(err instanceof ErrorCliente && err.status === 409 ? "Alguien más cambió este producto. Recarga y vuelve a intentar." : manejar(err));
    }
  }

  const campo = (k: keyof Campos, etiqueta: string, extra: React.InputHTMLAttributes<HTMLInputElement> = {}) => (
    <div>
      <label htmlFor={`p-${k}`} className="text-[14px] font-bold">
        {etiqueta}
      </label>
      <input
        id={`p-${k}`}
        value={String(c[k])}
        onChange={(e) => set(k, e.target.value as never)}
        aria-invalid={errores[k] ? true : undefined}
        aria-describedby={errores[k] ? `pe-${k}` : undefined}
        className="mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px] focus:border-negro aria-[invalid=true]:border-rojo"
        {...extra}
      />
      {errores[k] ? (
        <p id={`pe-${k}`} className="mt-1 text-[13px] font-semibold text-rojo">
          {errores[k]}
        </p>
      ) : null}
    </div>
  );
  const num = { inputMode: "numeric" as const, pattern: "[0-9]*", maxLength: 8 };

  return (
    <form onSubmit={guardar} noValidate className="space-y-4">
      <Link href="/admin/productos" className="inline-flex min-h-11 items-center gap-1 text-[14px] font-bold">
        <IconoFlechaIzq size={18} /> Productos
      </Link>
      <div className="rounded-[12px] bg-blanco p-4">
        <h1 className="text-[20px] font-extrabold leading-snug">{p.nombre}</h1>
        <p className="text-[13px] text-gris">
          {p.codigo} · {p.categoria.nombre}
          {p.es_ejemplo ? " · producto de ejemplo" : ""} ·{" "}
          <Link href={`/producto/${p.slug}`} className="underline" target="_blank">
            ver en la tienda
          </Link>
        </p>
      </div>

      <fieldset className="space-y-3 rounded-[12px] bg-blanco p-4">
        <legend className="sr-only">Precios</legend>
        {campo("precio", "Precio detalle", num)}
        {campo("precio_oferta", "Precio oferta (opcional)", num)}
        <div>
          <label htmlFor="p-etiqueta_oferta" className="text-[14px] font-bold">
            ¿Oferta o remate?
          </label>
          <select
            id="p-etiqueta_oferta"
            value={c.etiqueta_oferta}
            onChange={(e) => set("etiqueta_oferta", e.target.value === "remate" ? "remate" : "oferta")}
            className="mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px]"
          >
            <option value="oferta">Oferta</option>
            <option value="remate">Remate</option>
          </select>
        </div>
        {campo("precio_mayor", "Precio por mayor (opcional)", num)}
        {campo("mayor_minimo", "Desde cuántas unidades", num)}
        {campo("unidad_mayor", "Unidad (ej: caja de 12)", { maxLength: 60 })}
      </fieldset>

      <fieldset className="space-y-3 rounded-[12px] bg-blanco p-4">
        <legend className="sr-only">Stock y visibilidad</legend>
        {campo("stock", "Stock", num)}
        {campo("umbral_pocos", 'Mostrar "quedan pocas" con menos de', num)}
        <div>
          <label htmlFor="p-temporada" className="text-[14px] font-bold">
            Temporada
          </label>
          <select
            id="p-temporada"
            value={c.temporada}
            onChange={(e) => set("temporada", e.target.value as Temporada | "")}
            className="mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px]"
          >
            {TEMPORADAS.map((t) => (
              <option key={t.v} value={t.v}>
                {t.t}
              </option>
            ))}
          </select>
        </div>
        {campo("fecha_llegada", "Fecha de llegada", { type: "date" })}
        {campo("tiktok_url", "Link al video de TikTok", { type: "url", maxLength: 500, inputMode: "url" })}
        <label className="flex min-h-12 items-center gap-3 text-[15px] font-bold">
          <input type="checkbox" checked={c.activo} onChange={(e) => set("activo", e.target.checked)} className="size-6 accent-negro" />
          Visible en la web
        </label>
      </fieldset>

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
        {estado === "guardando" ? "Guardando…" : estado === "guardado" ? "Guardado" : "Guardar"}
      </button>
    </form>
  );
}
