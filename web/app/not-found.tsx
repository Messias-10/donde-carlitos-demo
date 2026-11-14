import Image from "next/image";
import Link from "next/link";

export default function NoEncontrado() {
  return (
    <main className="grid min-h-dvh place-items-center bg-negro px-4 text-center text-blanco">
      <div className="max-w-sm">
        <Image src="/marca/logo-donde-carlitos.png" alt="Donde Carlitos Emporio" width={120} height={120} className="mx-auto size-[120px] rounded-full" />
        <h1 className="titular texto-oro mt-6 text-[40px]">Esto no está en el local</h1>
        <p className="mt-3 text-[16px] text-blanco/85">Puede que el producto ya no esté o que el link venga cortado.</p>
        <div className="mt-6 flex flex-col gap-3">
          <Link href="/catalogo" className="flex h-12 items-center justify-center rounded-[12px] bg-oro font-extrabold text-negro">
            Ver el catálogo
          </Link>
          <Link href="/" className="flex h-12 items-center justify-center rounded-[12px] border-2 border-oro font-extrabold text-oro">
            Ir al inicio
          </Link>
        </div>
      </div>
    </main>
  );
}
