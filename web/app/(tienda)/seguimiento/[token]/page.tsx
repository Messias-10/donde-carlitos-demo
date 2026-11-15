import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { CopiarTexto } from "@/components/CopiarTexto";
import { IconoWhatsApp } from "@/components/Iconos";
import { Ejemplo, PorConfirmar } from "@/components/Marcas";
import { numeroCanal } from "@/components/landing/Canales";
import { headers } from "next/headers";
import { ErrorApiServidor, seguimiento } from "@/lib/api-servidor";
import { CHIP, PUNTO, SEGUIMIENTO } from "@/lib/estados";
import { fechaHora } from "@/lib/formato";
import { ES_DEMO, HORARIO } from "@/lib/negocio";
import { numerosWhatsApp } from "@/lib/numeros";
import { linkWhatsApp, mensajeSeguimiento } from "@/lib/whatsapp";

export const metadata: Metadata = {
  title: "Sigue tu pedido",
  robots: { index: false, follow: false },
  referrer: "no-referrer",
};

const DIRECCION = "Alameda 2432, Santiago Centro";

export default async function PaginaSeguimiento({ params }: { params: Promise<{ token: string }> }) {
  const { token } = await params;
  const h = await headers();
  // IP del cliente: el primer valor de X-Forwarded-For (si hay proxy) o X-Real-IP.
  const ip = (h.get("x-forwarded-for") ?? "").split(",")[0].trim() || h.get("x-real-ip") || undefined;
  const numeros = await numerosWhatsApp();
  let s: Awaited<ReturnType<typeof seguimiento>>;
  try {
    s = await seguimiento(token, ip);
  } catch (e) {
    if (e instanceof ErrorApiServidor && e.status === 429) return <MuchasConsultas />;
    throw e;
  }
  if (!s) notFound();

  const info = SEGUIMIENTO[s.estado];
  const canal = s.entrega === "envio" ? "regiones" : "santiago";
  const cancelado = s.estado === "cancelado";
  // Cancelado corta la línea y se muestra solo (estilo-carlitos §8.6).
  const linea = cancelado ? s.linea_tiempo.filter((p) => p.estado === "cancelado") : s.linea_tiempo;

  return (
    <div className="bg-crema px-4 pb-16 pt-6">
      <article className="mx-auto max-w-xl space-y-4">
        {ES_DEMO ? <Ejemplo texto="Pedido de prueba (demo)" /> : null}
        <h1 className="titular text-[34px]">
          Pedido <span className="whitespace-nowrap">{s.numero}</span>
        </h1>

        <div className={`rounded-[16px] p-5 ${CHIP[s.estado]}`}>
          <p className="flex items-center gap-2 font-display text-[28px] uppercase leading-none">
            <span className={`inline-block size-3 rounded-full ${PUNTO[s.estado]} ring-2 ring-current`} aria-hidden />
            {info.titulo}
          </p>
          <p className="mt-2 text-[16px]">{info.texto({ agencia: s.agencia_nombre, direccion: DIRECCION })}</p>
        </div>

        {s.estado === "despachado" && (s.agencia_nombre || s.n_seguimiento) ? (
          <div className="rounded-[16px] bg-blanco p-5">
            {s.agencia_nombre ? <p className="text-[15px]">Agencia: <strong>{s.agencia_nombre}</strong></p> : null}
            {s.n_seguimiento ? (
              <>
                <p className="mt-1 text-[15px]">
                  N.º de seguimiento:{" "}
                  <strong id="n-seguimiento" className="tabular-nums">
                    {s.n_seguimiento}
                  </strong>
                </p>
                <CopiarTexto texto={s.n_seguimiento} etiqueta="Copiar número" idObjetivo="n-seguimiento" className="mt-3 border-negro" />
                {s.agencia_nombre ? <p className="mt-2 text-[14px] text-gris">Con este número lo sigues en {s.agencia_nombre}.</p> : null}
              </>
            ) : null}
          </div>
        ) : null}

        {s.estado === "listo_retiro" ? (
          <div className="rounded-[16px] bg-blanco p-5 text-[15px]">
            <p>
              <span aria-hidden>📍 </span>
              {DIRECCION}
            </p>
            <p className="mt-1">
              <span aria-hidden>🕘 </span>L a V {HORARIO.semana.abre} a {HORARIO.semana.cierra} · Sáb, dom y feriados{" "}
              {HORARIO.finde.abre} a {HORARIO.finde.cierra} <PorConfirmar />
            </p>
            <p className="mt-2 font-bold">
              Muestra este número en el local: <span className="whitespace-nowrap">{s.numero}</span>
            </p>
          </div>
        ) : null}

        <section aria-labelledby="linea-titulo" className="rounded-[16px] bg-blanco p-5">
          <h2 id="linea-titulo" className="text-[17px] font-extrabold">
            En qué va
          </h2>
          <ol className="mt-3 space-y-3">
            {linea.map((p, i) => {
              const actual = i === linea.length - 1;
              return (
                <li key={`${p.estado}-${p.fecha}`} className="flex items-start gap-3">
                  <span
                    aria-hidden
                    className={`mt-1 inline-block shrink-0 rounded-full ${actual ? `size-4 ${PUNTO[p.estado]} ring-2 ring-negro/20` : "size-3 bg-oro-profundo"}`}
                  />
                  <span className="flex flex-1 flex-wrap justify-between gap-x-3 text-[15px]">
                    <span className={actual ? "font-bold" : ""}>
                      {SEGUIMIENTO[p.estado].titulo}
                      {actual ? <span className="sr-only"> (estado actual)</span> : null}
                    </span>
                    <span className="tabular-nums text-gris">{fechaHora(p.fecha)}</span>
                  </span>
                </li>
              );
            })}
          </ol>
        </section>

        <a
          href={linkWhatsApp(numeroCanal(numeros, canal), mensajeSeguimiento(s.numero, ES_DEMO))}
          target="_blank"
          rel="noopener noreferrer"
          className={`flex min-h-14 items-center justify-center gap-2.5 rounded-[14px] px-4 text-[17px] font-extrabold ${
            canal === "santiago" ? "bg-whatsapp text-negro" : "bg-whatsapp-osc text-blanco"
          }`}
        >
          <IconoWhatsApp size={22} />
          Escríbenos por WhatsApp
        </a>

        <p className="rounded-[12px] border-2 border-negro p-4 text-[15px] font-semibold">
          Pedidos solo a nuestros números oficiales. Los datos para transferir te los damos solo por el chat oficial.
        </p>
      </article>
    </div>
  );
}

/** QA I-5: un 429 de la API se explica, no termina en error 500. */
function MuchasConsultas() {
  return (
    <div className="bg-crema px-4 pb-16 pt-8">
      <div className="mx-auto max-w-xl space-y-3 rounded-[16px] bg-blanco p-5" role="status">
        <h1 className="titular text-[30px]">Hay muchas consultas</h1>
        <p className="text-[16px]">Intenta de nuevo en un minuto. Tu link sigue funcionando.</p>
      </div>
    </div>
  );
}
