"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { useCarrito } from "@/lib/carrito";
import { IconoCarrito, IconoCerrar, IconoLupa, IconoMenu } from "./Iconos";
import { useContadorAvisos } from "./Toast";

const ANCLAS = [
  { href: "/#ofertas", texto: "Ofertas y remates" },
  { href: "/catalogo", texto: "Catálogo completo" },
  { href: "/#pa-negocio", texto: "¿Pa negocio? Precio por mayor" },
  { href: "/#como-pedir", texto: "Cómo pedir" },
  { href: "/#envios", texto: "Retiro y envíos" },
  { href: "/#como-llegar", texto: "Cómo llegar y horario" },
  { href: "/#canal-oficial", texto: "Canal oficial" },
  { href: "/#preguntas", texto: "Preguntas frecuentes" },
];

export function HeaderAcciones() {
  const items = useCarrito();
  const n = items.reduce((s, l) => s + l.cantidad, 0);
  const [buscar, setBuscar] = useState(false);
  const [menu, setMenu] = useState(false);
  const [q, setQ] = useState("");
  const input = useRef<HTMLInputElement>(null);
  const router = useRouter();
  const ruta = usePathname();
  const [rutaVista, setRutaVista] = useState(ruta);

  // Al cambiar de página se cierran el buscador y el menú (estado derivado, sin efecto).
  if (ruta !== rutaVista) {
    setRutaVista(ruta);
    setMenu(false);
    setBuscar(false);
  }

  useEffect(() => {
    if (buscar) input.current?.focus();
  }, [buscar]);

  // Pop del contador solo cuando el usuario agrega algo.
  const pop = useContadorAvisos();

  return (
    <>
      <div className="ml-auto flex items-center gap-1">
        <button
          type="button"
          className="grid size-11 place-items-center rounded-full text-blanco hover:text-oro"
          aria-label={buscar ? "Cerrar buscador" : "Buscar productos"}
          aria-expanded={buscar}
          aria-controls="buscador"
          onClick={() => {
            setBuscar((b) => !b);
            setMenu(false);
          }}
        >
          {buscar ? <IconoCerrar /> : <IconoLupa />}
        </button>
        <Link
          href="/carrito"
          className="relative grid size-11 place-items-center rounded-full text-blanco hover:text-oro"
          aria-label={`Ver mi pedido, ${n} productos`}
        >
          <IconoCarrito />
          {n > 0 ? (
            <span
              key={pop}
              className={`absolute -right-0.5 -top-0.5 grid min-w-[22px] ${pop > 0 ? "animate-pop" : ""} place-items-center rounded-full bg-fluor px-1 text-[12px] font-extrabold leading-[22px] text-negro`}
              aria-hidden
            >
              {n > 99 ? "99+" : n}
            </span>
          ) : null}
        </Link>
        <button
          type="button"
          className="grid size-11 place-items-center rounded-full text-blanco hover:text-oro"
          aria-label={menu ? "Cerrar menú" : "Abrir menú"}
          aria-expanded={menu}
          aria-controls="menu-principal"
          onClick={() => {
            setMenu((m) => !m);
            setBuscar(false);
          }}
        >
          {menu ? <IconoCerrar /> : <IconoMenu />}
        </button>
      </div>

      {buscar ? (
        <form
          id="buscador"
          role="search"
          className="absolute inset-x-0 top-full border-t border-oro/30 bg-negro px-4 py-3"
          onSubmit={(e) => {
            e.preventDefault();
            const t = q.trim().slice(0, 60);
            router.push(t ? `/catalogo?q=${encodeURIComponent(t)}` : "/catalogo");
          }}
        >
          <div className="mx-auto flex max-w-6xl gap-2">
            <label htmlFor="q-header" className="sr-only">
              Buscar productos
            </label>
            <input
              ref={input}
              id="q-header"
              type="search"
              value={q}
              maxLength={60}
              onChange={(e) => setQ(e.target.value)}
              placeholder="¿Qué andas buscando? Ej: squishy, tren, telón"
              className="h-12 min-w-0 flex-1 rounded-[10px] bg-blanco px-3 text-[16px] text-negro placeholder:text-gris"
              enterKeyHint="search"
            />
            <button type="submit" className="h-12 rounded-[10px] bg-oro px-4 font-extrabold text-negro">
              Buscar
            </button>
          </div>
        </form>
      ) : null}

      {menu ? (
        <nav
          id="menu-principal"
          aria-label="Secciones"
          className="absolute inset-x-0 top-full border-t border-oro/30 bg-negro px-4 pb-4"
        >
          <ul className="mx-auto max-w-6xl divide-y divide-blanco/10">
            {ANCLAS.map((a) => (
              <li key={a.href}>
                <Link
                  href={a.href}
                  className="block py-3 text-[16px] font-semibold text-blanco hover:text-oro"
                  onClick={() => setMenu(false)}
                >
                  {a.texto}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      ) : null}
    </>
  );
}
