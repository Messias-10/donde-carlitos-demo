import { Footer } from "@/components/Footer";
import { Header } from "@/components/Header";
import { ES_DEMO } from "@/lib/negocio";
import { numerosWhatsApp } from "@/lib/numeros";

export default async function LayoutTienda({ children }: { children: React.ReactNode }) {
  const numeros = await numerosWhatsApp();
  return (
    <>
      <a
        href="#contenido"
        className="sr-only focus:not-sr-only focus:fixed focus:left-3 focus:top-3 focus:z-50 focus:rounded-[10px] focus:bg-oro focus:px-4 focus:py-2 focus:font-bold focus:text-negro"
      >
        Saltar al contenido
      </a>
      <Header />
      <main id="contenido" className="min-h-[60vh]">
        {children}
      </main>
      <Footer numeros={numeros} demo={ES_DEMO} />
    </>
  );
}
