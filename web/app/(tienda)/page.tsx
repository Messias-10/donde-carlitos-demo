import { BarraCanales } from "@/components/landing/BarraCanales";
import { hrefCanal, numeroCanal, TEXTO_CANAL } from "@/components/landing/Canales";
import {
  CanalOficial,
  ComoLlegar,
  ComprarDesdeRegiones,
  FamiliaEmporio,
  Preguntas,
  RetiroSantiago,
  SigueTuPedido,
} from "@/components/landing/Confianza";
import { Hero } from "@/components/landing/Hero";
import { HojaCanal } from "@/components/landing/HojaCanal";
import { Vitrina } from "@/components/landing/Vitrina";
import { listarProductos, vitrinaPublica } from "@/lib/api-servidor";
import { ES_DEMO, NEGOCIO } from "@/lib/negocio";
import { numerosWhatsApp } from "@/lib/numeros";
import type { ProductoVista } from "@/lib/tipos";

// Vitrina: GET /api/vitrina (los productos que el dueño marcó, en su orden).
// Si la API no lo tiene o falla, lo más reciente del catálogo.
async function vitrina(): Promise<ProductoVista[]> {
  try {
    return await vitrinaPublica();
  } catch {
    try {
      return (await listarProductos({ orden: "recientes", page_size: 8 })).items;
    } catch {
      return [];
    }
  }
}

export default async function Inicio() {
  const [numeros, productos] = await Promise.all([numerosWhatsApp(), vitrina()]);
  const demo = ES_DEMO;

  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "Store",
    name: NEGOCIO.nombre,
    description: "Importaciones y bazar para el hogar: juguetes virales, Navidad, hogar y escolar, al detalle y por mayor.",
    image: "/marca/logo-donde-carlitos.png",
    telephone: `+${numeros.santiago}`,
    address: {
      "@type": "PostalAddress",
      streetAddress: NEGOCIO.direccion,
      addressLocality: NEGOCIO.ciudad,
      addressRegion: NEGOCIO.region,
      addressCountry: "CL",
    },
    sameAs: [NEGOCIO.tiktok.url, NEGOCIO.instagram.url],
  };

  return (
    <div className="pb-20 lg:pb-0">
      <script
        type="application/ld+json"
        // Datos estáticos del negocio (no vienen de usuarios ni de la API). Se escapa "<".
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd).replace(/</g, "\\u003c") }}
      />
      <Hero numeros={numeros} demo={demo} />
      <Vitrina productos={productos} numeros={numeros} demo={demo} />
      <ComprarDesdeRegiones numeros={numeros} demo={demo} />
      <RetiroSantiago numeros={numeros} demo={demo} />
      <ComoLlegar />
      <CanalOficial numeros={numeros} />
      <SigueTuPedido />
      <FamiliaEmporio />
      <Preguntas numeros={numeros} />
      <HojaCanal numeros={{ santiago: numeroCanal(numeros, "santiago"), regiones: numeroCanal(numeros, "regiones") }} demo={demo} />
      <BarraCanales
        hrefs={{ santiago: hrefCanal(numeros, "santiago", demo), regiones: hrefCanal(numeros, "regiones", demo) }}
        textos={{
          santiago: { corto: TEXTO_CANAL.santiago.corto, aria: TEXTO_CANAL.santiago.aria },
          regiones: { corto: TEXTO_CANAL.regiones.corto, aria: TEXTO_CANAL.regiones.aria },
        }}
      />
    </div>
  );
}
