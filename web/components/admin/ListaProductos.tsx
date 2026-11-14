"use client";

import Link from "next/link";
import { useState } from "react";
import { clp } from "@/lib/formato";
import type { Pagina, ProductoAdmin } from "@/lib/tipos";
import { useApi } from "./useApi";

export function ListaProductos() {
  const [q, setQ] = useState("");
  const [busqueda, setBusqueda] = useState("");
  const [page, setPage] = useState(1);
  const p = new URLSearchParams({ page: String(page), page_size: "30" });
  if (busqueda) p.set("q", busqueda);
  const { data: datos, error } = useApi<Pagina<ProductoAdmin>>(`/api/admin/productos?${p}`);

  return (
    <div>
      <h1 className="titular text-[28px]">Productos</h1>
      <form
        role="search"
        className="mt-3 flex gap-2"
        onSubmit={(e) => {
          e.preventDefault();
          setBusqueda(q.trim().slice(0, 60));
          setPage(1);
        }}
      >
        <label htmlFor="buscar-producto" className="sr-only">
          Buscar producto
        </label>
        <input
          id="buscar-producto"
          type="search"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Buscar producto"
          maxLength={60}
          className="h-12 min-w-0 flex-1 rounded-[10px] border-2 border-negro/20 bg-blanco px-3 text-[16px]"
        />
        <button type="submit" className="h-12 rounded-[10px] bg-negro px-4 font-bold text-oro">
          Buscar
        </button>
      </form>

      {error ? (
        <p role="alert" className="mt-4 rounded-[10px] border-2 border-rojo p-3 font-semibold text-rojo">
          {error}
        </p>
      ) : null}

      {datos ? (
        <>
          <p className="mt-3 text-[14px] text-gris">{datos.total} productos</p>
          <ul className="mt-2 space-y-2">
            {datos.items.map((p) => (
              <li key={p.id}>
                <Link href={`/admin/productos/${p.id}`} className="flex items-center justify-between gap-3 rounded-[12px] bg-blanco p-3 ring-1 ring-negro/5 hover:ring-oro">
                  <span className="min-w-0">
                    <span className="block truncate text-[15px] font-semibold">{p.nombre}</span>
                    <span className="text-[13px] text-gris">
                      {p.codigo} · stock {p.stock}
                      {!p.activo ? " · oculto" : ""}
                      {p.precio_oferta ? ` · ${p.etiqueta_oferta === "remate" ? "remate" : "oferta"}` : ""}
                    </span>
                  </span>
                  <span className="shrink-0 font-display text-[20px]">{clp(p.precio_final)}</span>
                </Link>
              </li>
            ))}
          </ul>
          {datos.total_pages > 1 ? (
            <div className="mt-5 flex items-center justify-center gap-3">
              <button type="button" disabled={page <= 1} onClick={() => setPage((x) => x - 1)} className="h-11 rounded-[10px] bg-blanco px-4 font-bold ring-1 ring-negro/15 disabled:opacity-40">
                Anterior
              </button>
              <span className="text-[14px] text-gris">
                {datos.page} de {datos.total_pages}
              </span>
              <button
                type="button"
                disabled={page >= datos.total_pages}
                onClick={() => setPage((x) => x + 1)}
                className="h-11 rounded-[10px] bg-negro px-4 font-bold text-oro disabled:opacity-40"
              >
                Siguiente
              </button>
            </div>
          ) : null}
        </>
      ) : (
        <div className="mt-4 h-40 animate-pulse rounded-[12px] bg-blanco/70" aria-hidden />
      )}
    </div>
  );
}
