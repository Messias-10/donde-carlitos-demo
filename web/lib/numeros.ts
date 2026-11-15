// Números de WhatsApp en el servidor: GET /api/config-publica y, si la API no
// responde, NEXT_PUBLIC_WHATSAPP_SANTIAGO / _REGIONES. Los botones nunca dependen
// de que la API esté arriba (arquitectura-landing.md §A).
import { cache } from "react";
import { configPublica } from "./api-servidor";
import { NUMEROS_ENV, type Numeros } from "./negocio";

const ok = (n: string | null | undefined) => (n ? /^\d{8,15}$/.test(n) : false);

export const numerosWhatsApp = cache(async (): Promise<Numeros> => {
  try {
    const c = await configPublica();
    return {
      santiago: ok(c.whatsapp_santiago) ? c.whatsapp_santiago : NUMEROS_ENV.santiago,
      regiones: ok(c.whatsapp_regiones) ? c.whatsapp_regiones : null,
    };
  } catch {
    return NUMEROS_ENV;
  }
});
