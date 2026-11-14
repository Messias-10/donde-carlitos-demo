import { CanalOficial, ComoLlegar, FamiliaEmporio, Preguntas } from "@/components/landing/Confianza";
import { Hero } from "@/components/landing/Hero";
import { Ofertas } from "@/components/landing/Ofertas";
import { Categorias, LoVisteEnTikTok, PaNegocio, PasosPedido, RecienLlegado, RetiroEnvios } from "@/components/landing/Secciones";
import { AvisoSinApi } from "@/components/AvisoSinApi";
import { listarCategorias, listarProductos, todosLosProductos } from "@/lib/api-servidor";
import { NEGOCIO } from "@/lib/negocio";
import type { Categoria, ProductoVista } from "@/lib/tipos";
import { idTikTok } from "@/lib/tiktok";

// Videos sugeridos para "Lo viste en TikTok" (arquitectura-landing.md §3), en orden.
const VIDEOS_TIKTOK = [
  "7706656443843284232", // tren navideño
  "7707649267464293639", // squishy que cambia de color
  "7707296451398864136", // caja sorpresa squishy
  "7707053973999471879", // dinosaurio gigante
  "7671107953398664466", // espada con luces
  "7699270769963978002", // calendario de adviento squishy
];
// Videos que muestran una dirección anterior (top-contenido.md ⚠).
const DIRECCION_VIEJA = new Set<string>([
  "https://www.tiktok.com/@mayte_emporio/video/7476070344000933125",
]);

async function seguro<T>(p: Promise<T>, def: T): Promise<{ v: T; ok: boolean }> {
  try {
    return { v: await p, ok: true };
  } catch {
    return { v: def, ok: false };
  }
}

export default async function Inicio() {
  const [cats, ofertas, nuevos, todos] = await Promise.all([
    seguro<Categoria[]>(listarCategorias(), []),
    seguro(listarProductos({ oferta: true, page_size: 12 }).then((r) => r.items), [] as ProductoVista[]),
    seguro(listarProductos({ recien_llegado: true, orden: "recientes", page_size: 8 }).then((r) => r.items), [] as ProductoVista[]),
    seguro(todosLosProductos(), [] as ProductoVista[]),
  ]);
  const sinApi = !cats.ok && !todos.ok;

  const porSlug = new Map(todos.v.map((p) => [p.slug, p]));

  // Ofertas reales primero; si hay pocas, se completa con lo del remate de Navidad.
  const filaOfertas = [...ofertas.v];
  if (filaOfertas.length < 6) {
    const ya = new Set(filaOfertas.map((p) => p.id));
    for (const p of todos.v) {
      if (filaOfertas.length >= 8) break;
      if (p.temporada === "navidad" && !ya.has(p.id) && p.disponibilidad !== "agotado") filaOfertas.push(p);
    }
  }

  const tiktok: ProductoVista[] = [];
  const vistos = new Set<string>();
  for (const id of VIDEOS_TIKTOK) {
    const p = todos.v.find((x) => idTikTok(x.tiktok_url) === id);
    if (p && !vistos.has(id)) {
      tiktok.push(p);
      vistos.add(id);
    }
  }

  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "Store",
    name: NEGOCIO.nombre,
    description: "Importaciones y bazar para el hogar: juguetes virales, Navidad, hogar y escolar, al detalle y por mayor.",
    image: "/marca/logo-donde-carlitos.png",
    telephone: NEGOCIO.telefonoE164,
    address: {
      "@type": "PostalAddress",
      streetAddress: NEGOCIO.direccion,
      addressLocality: NEGOCIO.ciudad,
      addressRegion: NEGOCIO.region,
      addressCountry: "CL",
    },
    openingHoursSpecification: [
      { "@type": "OpeningHoursSpecification", dayOfWeek: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"], opens: "09:00", closes: "18:00" },
      { "@type": "OpeningHoursSpecification", dayOfWeek: ["Saturday", "Sunday"], opens: "09:00", closes: "16:00" },
    ],
    sameAs: [NEGOCIO.tiktok.url, NEGOCIO.instagram.url],
  };

  return (
    <>
      <script
        type="application/ld+json"
        // Datos estáticos del negocio (no vienen de usuarios ni de la API). Se escapa "<".
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd).replace(/</g, "\\u003c") }}
      />
      <Hero />
      {sinApi ? <AvisoSinApi /> : null}
      <Ofertas productos={filaOfertas} esfera={porSlug.get("esfera-navidena-con-luces") ?? null} />
      <LoVisteEnTikTok productos={tiktok} conDireccionVieja={DIRECCION_VIEJA} />
      <Categorias categorias={cats.v} />
      <RecienLlegado productos={nuevos.v} />
      <PaNegocio ejemplo={porSlug.get("squishy-viral-que-cambia-de-color-con-el-sol") ?? null} />
      <PasosPedido />
      <RetiroEnvios />
      <ComoLlegar />
      <CanalOficial />
      <FamiliaEmporio />
      <Preguntas />
    </>
  );
}
