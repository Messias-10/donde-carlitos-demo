"use client";

import type { Ajustes as DatosAjustes } from "@/lib/tipos";
import { SoloAdmin } from "./Sesion";
import { Aviso } from "./ui";
import { useApi } from "./useApi";

/** "+56921145368" o "56921145368" -> "+56 9 2114 5368". */
function tel(n: string | null): string {
  if (!n) return "";
  const d = n.replace(/\D/g, "");
  const m = /^56(9)(\d{4})(\d{4})$/.exec(d);
  return m ? `+56 ${m[1]} ${m[2]} ${m[3]}` : `+${d}`;
}

function Fila({ etiqueta, children }: { etiqueta: string; children: React.ReactNode }) {
  return (
    <div className="border-t border-negro/10 py-3 first:border-t-0">
      <dt className="text-[13px] font-semibold text-gris">{etiqueta}</dt>
      <dd className="mt-0.5 break-words text-[16px] font-semibold">{children}</dd>
    </div>
  );
}

function Contenido() {
  const { data, error } = useApi<DatosAjustes>("/api/admin/ajustes");
  return (
    <div className="space-y-4">
      <h1 className="titular text-[28px]">Ajustes</h1>
      <p className="text-[15px] text-gris">Solo lectura. Estos datos se configuran en la instalación.</p>
      {error ? <Aviso>{error}</Aviso> : null}
      {data ? (
        <dl className="rounded-[12px] bg-blanco px-4">
          <Fila etiqueta="WhatsApp Santiago">{tel(data.whatsapp_santiago)}</Fila>
          <Fila etiqueta="WhatsApp Regiones">
            {data.whatsapp_regiones ? tel(data.whatsapp_regiones) : "Sin número: se usa el de Santiago"}
          </Fila>
          <Fila etiqueta="Dirección de retiro">{data.direccion_retiro}</Fila>
          <Fila etiqueta="Horario de retiro">{data.horario_retiro || "Sin horario"}</Fila>
          <Fila etiqueta="Dirección del sitio (links de seguimiento)">{data.url_sitio}</Fila>
          <Fila etiqueta="Datos de transferencia">
            {data.datos_transferencia_configurados ? (
              <span className="text-st-pagado-fg">Configurados ✔</span>
            ) : (
              <>Sin configurar — te los configura Luckas en la instalación</>
            )}
            <span className="mt-1 block text-[13px] font-normal text-gris">Se mandan solo por WhatsApp dentro de la cotización. Nunca se publican en la web.</span>
          </Fila>
          <Fila etiqueta="Modo demo">
            {data.modo_demo ? "Activo: los pedidos son de prueba y los mensajes llevan *[PRUEBA · DEMO]*" : "Apagado"}
          </Fila>
        </dl>
      ) : !error ? (
        <div className="h-60 animate-pulse rounded-[12px] bg-blanco" aria-hidden />
      ) : null}
    </div>
  );
}

export function Ajustes() {
  return (
    <SoloAdmin>
      <Contenido />
    </SoloAdmin>
  );
}
