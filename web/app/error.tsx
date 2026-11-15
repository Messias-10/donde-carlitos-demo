"use client";

import Link from "next/link";

// Sin detalles internos: solo una salida clara.
export default function ErrorPagina({ reset }: { error: Error; reset: () => void }) {
  return (
    <main className="grid min-h-dvh place-items-center bg-crema px-4 text-center">
      <div className="max-w-sm">
        <h1 className="titular text-[36px]">Algo falló al cargar</h1>
        <p className="mt-3 text-[16px] text-negro/80">Inténtalo de nuevo en un rato o escríbenos al +56 9 2114 5368.</p>
        <div className="mt-6 flex flex-col gap-3">
          <button type="button" onClick={reset} className="h-12 rounded-[12px] bg-negro font-extrabold text-oro">
            Reintentar
          </button>
          <Link href="/" className="flex h-12 items-center justify-center rounded-[12px] border-2 border-negro font-extrabold">
            Ir al inicio
          </Link>
        </div>
      </div>
    </main>
  );
}
