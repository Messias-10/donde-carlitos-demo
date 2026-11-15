// Bloques 3 a 9 de arquitectura-landing.md §A, con el texto de copy-landing.md §4–§10.
import type { ReactNode } from "react";
import { IconoEscudo, IconoMetro, IconoReloj } from "@/components/Iconos";
import { PorConfirmar } from "@/components/Marcas";
import { Seccion } from "@/components/Seccion";
import { HORARIO, NEGOCIO, telefonoVisible, type Numeros } from "@/lib/negocio";
import { BotonCanal } from "./Canales";

function Pasos({ pasos, oscuro = false }: { pasos: { t: string; d: ReactNode }[]; oscuro?: boolean }) {
  return (
    <ol className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
      {pasos.map((p, i) => (
        <li key={p.t} className={`flex gap-4 rounded-[12px] p-4 ${oscuro ? "bg-carbon" : "bg-blanco"}`}>
          <span className={`font-display text-[44px] leading-none ${oscuro ? "text-oro" : "text-oro-texto"}`} aria-hidden>
            {i + 1}
          </span>
          <div>
            <h3 className={`titular text-[22px] ${oscuro ? "text-blanco" : ""}`}>{p.t}</h3>
            <p className={`mt-1 text-[15px] ${oscuro ? "text-blanco/85" : "text-negro/85"}`}>{p.d}</p>
          </div>
        </li>
      ))}
    </ol>
  );
}

// ---------------------------------------------------------------- 3. Regiones

export function ComprarDesdeRegiones({ numeros, demo }: { numeros: Numeros; demo: boolean }) {
  const pasos = [
    { t: "Escríbenos", d: 'Toca "Pedidos a regiones" y cuéntanos tu comuna y lo que quieres. Si lo viste en un video, mándanos el nombre o el link.' },
    { t: "Te cotizamos", d: "Te mandamos por el chat el resumen con tu número de pedido y el total." },
    {
      t: "Transfieres",
      d: (
        <>
          Los datos para transferir te los damos <strong>solo por nuestro chat oficial</strong>. No están publicados en ninguna
          parte.
        </>
      ),
    },
    { t: "Nos mandas el comprobante", d: "Por el mismo chat. Revisamos que el pago llegó y te confirmamos." },
    { t: "Despachamos", d: "Por agencia, y te mandamos el número de seguimiento y un link para ver en qué va tu pedido." },
  ];
  const saber = [
    "Agencias: por confirmar",
    "Costo de envío: te lo confirmamos en la cotización",
    "Compra mínima: por confirmar",
    "Plazo de despacho: por confirmar",
    "Pago por transferencia (por confirmar)",
  ];
  return (
    <Seccion
      id="regiones"
      fondo="negro"
      titulo={
        <>
          ¿Eres de región? Así compras <span aria-hidden>📦</span>
        </>
      }
    >
      <Pasos pasos={pasos} oscuro />
      <div className="mt-5 rounded-[12px] border border-oro/40 p-4">
        <h3 className="text-[17px] font-extrabold text-blanco">Lo que tienes que saber</h3>
        <ul className="mt-2 space-y-1 text-[15px] text-blanco/90">
          {saber.map((s) => (
            <li key={s}>{s}</li>
          ))}
        </ul>
        <p className="mt-2">
          <PorConfirmar oscuro />
        </p>
      </div>
      <div className="mt-5 max-w-md">
        <BotonCanal canal="regiones" numeros={numeros} demo={demo} />
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- 4. Santiago

export function RetiroSantiago({ numeros, demo }: { numeros: Numeros; demo: boolean }) {
  const pasos = [
    { t: "Escríbenos", d: "Escríbenos por WhatsApp o ven directo al local." },
    {
      t: "Te confirmamos",
      d: (
        <>
          Te confirmamos el total y te avisamos cuando tu pedido esté listo. <PorConfirmar />
        </>
      ),
    },
    { t: "Retiras", d: "Lo retiras en Alameda 2432 con tu número de pedido." },
  ];
  return (
    <Seccion
      id="santiago"
      titulo={
        <>
          ¿Eres de Santiago? Retira en tienda <span aria-hidden>📍</span>
        </>
      }
    >
      <Pasos pasos={pasos} />
      <div className="mt-5 max-w-md">
        <BotonCanal canal="santiago" numeros={numeros} demo={demo} />
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- 5. Cómo llegar

export function ComoLlegar() {
  return (
    <Seccion
      id="como-llegar"
      fondo="blanco"
      titulo={
        <>
          Cómo llegar <span aria-hidden>📍</span>
        </>
      }
    >
      <div className="grid gap-6 lg:grid-cols-2">
        <div>
          <p className="font-display text-[26px] uppercase leading-tight">Alameda 2432, Santiago Centro</p>
          <ul className="mt-4 space-y-3 text-[15px]">
            <li className="flex gap-3">
              <IconoMetro size={22} className="mt-0.5 shrink-0 text-oro-texto" />
              <span>
                <strong>Metro Unión Latinoamericana</strong> (Línea 1): a 5 minutos caminando.
              </span>
            </li>
            <li className="flex gap-3">
              <IconoMetro size={22} className="mt-0.5 shrink-0 text-oro-texto" />
              <span>
                <strong>Metro República</strong> (Línea 1): a 7 minutos caminando.
              </span>
            </li>
            <li className="flex gap-3">
              <IconoMetro size={22} className="mt-0.5 shrink-0 text-oro-texto" />
              <span>
                <strong>¿Vienes en bus desde región?</strong> Estás a una estación de Metro Estación Central (Terminal San
                Borja) y a dos de Universidad de Santiago (terminales Alameda y Sur).
              </span>
            </li>
          </ul>
          <a
            href={NEGOCIO.mapsUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="mt-5 inline-flex h-12 items-center rounded-[12px] bg-negro px-5 text-[15px] font-extrabold text-oro"
          >
            Abrir en Google Maps
          </a>
        </div>
        <div className="space-y-4">
          <div className="rounded-[12px] bg-crema p-4">
            <h3 className="flex items-center gap-2 text-[17px] font-extrabold">
              <IconoReloj size={20} className="text-oro-texto" /> Horario de temporada
            </h3>
            <ul className="mt-2 space-y-1 text-[15px]">
              <li>
                Lunes a viernes: {HORARIO.semana.abre} a {HORARIO.semana.cierra}
              </li>
              <li>
                Sábado, domingo y feriados: {HORARIO.finde.abre} a {HORARIO.finde.cierra}
              </li>
            </ul>
            <p className="mt-2 text-[13px] text-gris">
              <em>{HORARIO.fuente}</em> <PorConfirmar />
            </p>
          </div>
          <div className="rounded-[12px] bg-crema p-4">
            <h3 className="text-[17px] font-extrabold">Cómo pagar en el local</h3>
            <p className="mt-2 text-[15px]">En Alameda 2432 tenemos dos locales:</p>
            <ul className="mt-1 space-y-1 text-[15px]">
              <li>
                <strong>Hogar, electro y juguetes:</strong> efectivo y débito.
              </li>
              <li>
                <strong>Disney y licencias:</strong> débito, crédito, transferencia y efectivo.
              </li>
            </ul>
            <p className="mt-2 text-[15px]">
              Boleta o factura, como necesites. <PorConfirmar />
            </p>
          </div>
        </div>
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- 6. Canal oficial

export function CanalOficial({ numeros }: { numeros: Numeros }) {
  const unSoloNumero = numeros.regiones === null || numeros.regiones === numeros.santiago;
  const reglas = [
    "Pide solo a nuestros números oficiales. Ningún otro.",
    <>
      <strong>No publicamos datos bancarios en la web.</strong> Te los damos solo por el chat oficial, junto con tu número de
      pedido.
    </>,
    "Nunca te vamos a pedir claves ni códigos que te lleguen por SMS.",
    "Antes de transferir, revisa que el chat sea con uno de nuestros números.",
    "¿Te escribieron a nombre nuestro desde otra cuenta o número? No transfieras y avísanos al número oficial.",
  ];
  return (
    <Seccion id="canal-oficial" fondo="carbon" titulo="Este es nuestro único canal oficial">
      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-[16px] border border-oro/50 p-5">
          {unSoloNumero ? (
            <>
              <p className="text-[15px] text-gris-claro">Pedidos solo a este número:</p>
              <p className="mt-1 font-display text-[34px] leading-none text-oro">{telefonoVisible(numeros.santiago)}</p>
              {numeros.regiones === null ? (
                <p className="mt-2 text-[13px] text-gris-claro">
                  Número de regiones: <PorConfirmar oscuro />
                </p>
              ) : null}
            </>
          ) : (
            <>
              <p className="text-[15px] text-gris-claro">Pedidos solo a estos números:</p>
              <p className="mt-2 text-[15px] text-blanco">Santiago y retiro en tienda:</p>
              <p className="font-display text-[30px] leading-none text-oro">{telefonoVisible(numeros.santiago)}</p>
              <p className="mt-3 text-[15px] text-blanco">Regiones:</p>
              <p className="font-display text-[30px] leading-none text-oro">{telefonoVisible(numeros.regiones!)}</p>
            </>
          )}
          <p className="mt-6 text-[15px] text-blanco">En redes nos encuentras como Mayte Emporio:</p>
          <ul className="mt-1 text-[15px]">
            <li>
              TikTok{" "}
              <a
                href={NEGOCIO.tiktok.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex min-h-11 items-center font-bold text-oro underline underline-offset-2"
              >
                {NEGOCIO.tiktok.usuario}
              </a>
            </li>
            <li>
              Instagram{" "}
              <a
                href={NEGOCIO.instagram.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex min-h-11 items-center font-bold text-oro underline underline-offset-2"
              >
                {NEGOCIO.instagram.usuario}
              </a>{" "}
              (cuenta verificada)
            </li>
          </ul>
        </div>
        <div>
          <h3 className="flex items-center gap-2 text-[18px] font-extrabold text-blanco">
            <IconoEscudo size={22} className="text-oro" /> Para comprar seguro
          </h3>
          <ol className="mt-3 space-y-3">
            {reglas.map((r, i) => (
              <li key={i} className="flex gap-3 text-[15px] text-blanco/90">
                <span className="font-display text-[22px] leading-none text-oro" aria-hidden>
                  {i + 1}
                </span>
                <span>{r}</span>
              </li>
            ))}
          </ol>
        </div>
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- 7. Sigue tu pedido

export function SigueTuPedido() {
  return (
    <Seccion id="seguimiento" fondo="blanco" titulo="¿Ya pediste? Sigue tu pedido">
      <p className="max-w-[60ch] text-[16px]">
        Cuando hagas tu pedido por WhatsApp, te mandamos un link para ver en qué va: si ya lo pagaste, si está listo o el
        número de seguimiento de la agencia. Sin dar tus datos.
      </p>
    </Seccion>
  );
}

// ---------------------------------------------------------------- 8. Familia Emporio

export function FamiliaEmporio() {
  return (
    <Seccion id="familia" titulo="Del toldo amarillo al contenedor" script="Familia Emporio">
      <p className="max-w-[52ch] text-[16px]">
        Partimos con un toldo amarillo. Hoy abrimos contenedores en plena calle y te dejamos lo nuevo a precio de remate.
      </p>
      <p className="mt-2">
        <PorConfirmar texto="historia por confirmar con el dueño" />
      </p>
      <a
        href={NEGOCIO.videoContenedor}
        target="_blank"
        rel="noopener noreferrer"
        className="mt-4 inline-flex h-12 items-center rounded-[12px] border-2 border-negro px-5 text-[15px] font-extrabold"
      >
        ▶ Mira cómo abrimos un contenedor
      </a>
    </Seccion>
  );
}

// ---------------------------------------------------------------- 9. Preguntas

export function Preguntas({ numeros }: { numeros: Numeros }) {
  const santiago = telefonoVisible(numeros.santiago);
  const regiones = numeros.regiones ? telefonoVisible(numeros.regiones) : null;
  const lista: { p: string; r: ReactNode }[] = [
    {
      p: "¿Dónde están? ¿Qué metro queda cerca?",
      r: "En Alameda 2432, Santiago Centro, a 5 minutos del Metro Unión Latinoamericana y a 7 del Metro República.",
    },
    {
      p: "¿Hacen envíos a regiones?",
      r: (
        <>
          Sí. Toca &quot;Pedidos a regiones&quot;, te cotizamos por WhatsApp y, cuando confirmamos tu pago, lo despachamos por
          agencia con número de seguimiento. <PorConfirmar texto="agencias, quién paga y mínimo por confirmar" />
        </>
      ),
    },
    {
      p: "¿A qué número les escribo?",
      r: (
        <>
          Si eres de Santiago o vas a retirar: {santiago}. Si eres de región:{" "}
          {regiones ?? (
            <>
              {santiago} <PorConfirmar />
            </>
          )}
          . Los dos botones de esta página te llevan al correcto.
        </>
      ),
    },
    {
      p: "¿Les queda…?",
      r: 'Pregúntanos por WhatsApp. Si lo viste en un video, toca "Pedir este por WhatsApp" y nos llega el producto ya escrito.',
    },
    {
      p: "¿Cuántos vienen por embalaje? ¿Desde cuántas es precio por mayor?",
      r: 'En los productos de "Lo nuevo de los videos" está el precio por mayor cuando lo dijimos en el video. Para lo demás, pregúntanos y te cotizamos.',
    },
    {
      p: "¿Hasta qué hora atienden? ¿Abren domingo y feriados?",
      r: (
        <>
          En temporada: lunes a viernes de 9:00 a 18:00, y sábado, domingo y feriados de 9:00 a 16:00. <PorConfirmar />
        </>
      ),
    },
    {
      p: "¿Es seguro transferir?",
      r: "Transfiere solo a los datos que te damos en el chat de nuestros números oficiales, y después de recibir tu cotización con número de pedido. Si alguien te escribe desde otra cuenta a nombre nuestro, no transfieras.",
    },
    {
      p: "¿Cómo sé en qué va mi pedido?",
      r: "Te mandamos un link de seguimiento por WhatsApp. Cuando lo despachamos, te llega también la agencia y el número de seguimiento.",
    },
    {
      p: "¿Mi caja llegó con un problema, a quién le digo?",
      r: "Escríbenos por el mismo chat donde hiciste el pedido, con tu número de pedido. Así lo vemos al tiro.",
    },
    {
      p: "¿Qué medios de pago aceptan en el local?",
      r: (
        <>
          En hogar, electro y juguetes: efectivo y débito. En Disney y licencias: débito, crédito, transferencia y efectivo.{" "}
          <PorConfirmar />
        </>
      ),
    },
    {
      p: "¿Dan boleta o factura?",
      r: (
        <>
          Sí, boleta o factura, como necesites. <PorConfirmar />
        </>
      ),
    },
  ];
  return (
    <Seccion
      id="preguntas"
      titulo={
        <>
          Preguntas que nos hacen siempre <span aria-hidden>💬</span>
        </>
      }
    >
      <div className="divide-y divide-negro/10 rounded-[12px] bg-blanco">
        {lista.map((q) => (
          <details key={q.p} className="group px-4">
            <summary className="flex min-h-14 items-center justify-between gap-3 py-3 text-[16px] font-bold">
              {q.p}
              <span className="text-[22px] leading-none text-oro-texto transition-transform group-open:rotate-45" aria-hidden>
                +
              </span>
            </summary>
            <p className="pb-4 text-[15px] text-negro/85">{q.r}</p>
          </details>
        ))}
      </div>
    </Seccion>
  );
}
