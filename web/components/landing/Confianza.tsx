import { IconoEscudo, IconoMetro, IconoReloj, IconoWhatsApp } from "@/components/Iconos";
import { MapaCarga } from "@/components/MapaCarga";
import { PorConfirmar } from "@/components/Marcas";
import { Seccion } from "@/components/Seccion";
import { VideoTikTok } from "@/components/VideoTikTok";
import { HORARIO, NEGOCIO } from "@/lib/negocio";
import { linkWhatsApp, MSJ_CONSULTA } from "@/lib/whatsapp";

// ---------------------------------------------------------------- Cómo llegar

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
              <IconoMetro size={22} className="mt-0.5 shrink-0 text-oro-profundo" />
              <span>
                <strong>Metro Unión Latinoamericana</strong> (Línea 1): a 5 minutos caminando.
              </span>
            </li>
            <li className="flex gap-3">
              <IconoMetro size={22} className="mt-0.5 shrink-0 text-oro-profundo" />
              <span>
                <strong>Metro República</strong> (Línea 1): a 7 minutos caminando.
              </span>
            </li>
            <li className="flex gap-3">
              <IconoMetro size={22} className="mt-0.5 shrink-0 text-oro-profundo" />
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

          <div className="mt-8 rounded-[12px] bg-crema p-4">
            <h3 className="flex items-center gap-2 text-[17px] font-extrabold">
              <IconoReloj size={20} className="text-oro-profundo" /> Horario de temporada
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

          <div className="mt-4 rounded-[12px] bg-crema p-4">
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
        <MapaCarga />
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Canal oficial

export function CanalOficial() {
  const reglas = [
    <>Los pedidos se hacen en esta web o al {NEGOCIO.telefonoVisible}. Ningún otro número.</>,
    <>
      Nunca te escribimos primero por interno para venderte. <PorConfirmar oscuro texto="redacción por confirmar" />
    </>,
    <>Nunca te vamos a pedir claves ni códigos que te lleguen por SMS.</>,
    <>Antes de pagar, revisa que el chat sea con nuestro número oficial.</>,
    <>¿Te escribieron a nombre de nosotros desde otra cuenta o número? No transfieras y avísanos al número oficial.</>,
  ];
  return (
    <Seccion id="canal-oficial" fondo="negro" titulo="Este es nuestro único canal oficial">
      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-[16px] border border-oro/50 p-5">
          <p className="text-[15px] text-gris-claro">Pedidos solo a este número:</p>
          <p className="mt-1 font-display text-[36px] leading-none text-oro">{NEGOCIO.telefonoVisible}</p>
          <a
            href={linkWhatsApp(MSJ_CONSULTA)}
            target="_blank"
            rel="noopener noreferrer"
            className="mt-4 inline-flex h-[52px] items-center gap-2 rounded-[14px] bg-whatsapp px-5 text-[16px] font-extrabold text-negro"
          >
            <IconoWhatsApp size={22} />
            Escribir a {NEGOCIO.telefonoVisible}
          </a>
          <p className="mt-6 text-[15px] text-blanco">En redes nos encuentras como Mayte Emporio:</p>
          <ul className="mt-2 space-y-1 text-[15px]">
            <li>
              TikTok{" "}
              <a href={NEGOCIO.tiktok.url} target="_blank" rel="noopener noreferrer" className="font-bold text-oro underline underline-offset-2">
                {NEGOCIO.tiktok.usuario}
              </a>
            </li>
            <li>
              Instagram{" "}
              <a href={NEGOCIO.instagram.url} target="_blank" rel="noopener noreferrer" className="font-bold text-oro underline underline-offset-2">
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

// ---------------------------------------------------------------- Familia Emporio

export function FamiliaEmporio() {
  return (
    <Seccion id="familia" titulo="Del toldo amarillo al contenedor" script="Familia Emporio">
      <div className="grid gap-6 sm:grid-cols-[1fr_200px] sm:items-start">
        <div>
          <p className="max-w-[52ch] text-[16px]">
            Partimos con un toldo amarillo. Hoy abrimos contenedores en plena calle y te dejamos lo nuevo a precio de
            remate.
          </p>
          <p className="mt-3">
            <PorConfirmar texto="historia por confirmar con el dueño" />
          </p>
        </div>
        <div className="w-[60vw] max-w-[200px]">
          <VideoTikTok url="https://www.tiktok.com/@mayte_emporio/video/7441724226857454904" titulo="Así abrimos un contenedor" />
        </div>
      </div>
    </Seccion>
  );
}

// ---------------------------------------------------------------- Preguntas

const PREGUNTAS: { p: string; r: React.ReactNode }[] = [
  {
    p: "¿Dónde están? ¿Qué metro queda cerca?",
    r: 'En Alameda 2432, Santiago Centro, a 5 minutos del Metro Unión Latinoamericana y a 7 del Metro República. Mira el mapa en "Cómo llegar".',
  },
  {
    p: "¿Hacen envíos a regiones?",
    r: (
      <>
        Sí: arma tu pedido, elige &quot;Envío por agencia&quot; y te confirmamos por WhatsApp el total y el despacho.{" "}
        <PorConfirmar texto="agencias, quién paga y mínimo por confirmar" />
      </>
    ),
  },
  {
    p: "¿Les queda…? ¿Cómo sé si hay stock?",
    r: "Cada producto dice si está Disponible, si quedan pocas unidades o si está Agotado. Si no lo encuentras en la web, pregúntanos por WhatsApp: muchas veces vuelve en el próximo contenedor.",
  },
  {
    p: "¿Cuántos vienen por embalaje? ¿Desde cuántas es precio por mayor?",
    r: 'Está en cada producto: el precio al detalle, el precio por mayor y desde cuántas unidades aplica (por ejemplo, "caja de 12"). En el carrito se aplica solo.',
  },
  {
    p: "¿Cómo compro? ¿Dónde les escribo?",
    r: `En 3 pasos: eliges, envías tu pedido y te confirmamos por WhatsApp. Nuestro único número de pedidos es el ${NEGOCIO.telefonoVisible}.`,
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
    p: "¿Cuánto vale?",
    r: (
      <>
        Todos los productos tienen su precio en la web, al detalle y por mayor cuando aplica. Tu pedido guarda el precio
        del momento en que lo hiciste. <PorConfirmar />
      </>
    ),
  },
  {
    p: "¿Qué medios de pago aceptan? ¿Puedo pagar con crédito?",
    r: (
      <>
        Depende del local. En hogar, electro y juguetes: efectivo y débito. En Disney y licencias: débito, crédito,
        transferencia y efectivo. <PorConfirmar />
      </>
    ),
  },
  {
    p: '¿Cómo funciona "Yo quiero mi promo"?',
    r: (
      <>
        En la web, los productos en promo ya muestran su precio especial. No tienes que comentar nada: lo agregas y
        listo. <PorConfirmar />
      </>
    ),
  },
  {
    p: "¿Tienen tallas?",
    r: "Las tallas disponibles van en la descripción de cada producto. Si no ves la tuya, escríbela en el comentario del pedido o pregúntanos por WhatsApp.",
  },
  {
    p: "¿Es seguro comprar con envío?",
    r: `Haz tu pedido en esta web o al ${NEGOCIO.telefonoVisible} y paga solo cuando te confirmemos desde ese número. Si alguien te escribe desde otra cuenta a nombre nuestro, no transfieras.`,
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

export function Preguntas() {
  return (
    <Seccion
      id="preguntas"
      fondo="blanco"
      titulo={
        <>
          Preguntas que nos hacen siempre <span aria-hidden>💬</span>
        </>
      }
    >
      <div className="divide-y divide-negro/10 rounded-[12px] bg-crema">
        {PREGUNTAS.map((q) => (
          <details key={q.p} className="group px-4">
            <summary className="flex min-h-14 items-center justify-between gap-3 py-3 text-[16px] font-bold">
              {q.p}
              <span className="text-[22px] leading-none text-oro-profundo transition-transform group-open:rotate-45" aria-hidden>
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
