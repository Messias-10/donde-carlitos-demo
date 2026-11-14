import Image from "next/image";
import Link from "next/link";
import { HeaderAcciones } from "./HeaderAcciones";

export function Header() {
  return (
    <header className="sticky top-0 z-40 bg-negro text-blanco shadow-[0_1px_0_rgba(226,179,60,0.35)]">
      <div className="mx-auto flex h-16 max-w-6xl items-center gap-3 px-4 lg:h-[76px]">
        <Link href="/" className="flex shrink-0 items-center gap-2.5" aria-label="Donde Carlitos Emporio, inicio">
          <Image
            src="/marca/logo-donde-carlitos.png"
            alt="Donde Carlitos Emporio"
            width={44}
            height={44}
            className="size-11 rounded-full lg:size-14"
            sizes="(min-width: 1024px) 56px, 44px"
            loading="eager"
          />
          <span className="hidden font-display text-[20px] uppercase leading-none tracking-wide text-oro sm:block">
            Donde Carlitos Emporio
          </span>
        </Link>
        <HeaderAcciones />
      </div>
    </header>
  );
}
