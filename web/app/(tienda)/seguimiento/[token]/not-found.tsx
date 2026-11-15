import { BotonCanal } from "@/components/landing/Canales";
import { ES_DEMO } from "@/lib/negocio";
import { numerosWhatsApp } from "@/lib/numeros";

// Token inválido o vencido: mismo mensaje en los dos casos, sin revelar si el pedido existe.
export default async function SeguimientoNoEncontrado() {
  const numeros = await numerosWhatsApp();
  return (
    <div className="bg-crema px-4 pb-16 pt-8">
      <div className="mx-auto max-w-xl space-y-4 rounded-[16px] bg-blanco p-5">
        <h1 className="titular text-[32px]">No encontramos este pedido</h1>
        <p className="text-[16px]">Revisa que el link esté completo o escríbenos por WhatsApp con tu número de pedido.</p>
        <div className="grid gap-3">
          <BotonCanal canal="santiago" numeros={numeros} demo={ES_DEMO} texto="Santiago / retiro" />
          <BotonCanal canal="regiones" numeros={numeros} demo={ES_DEMO} texto="Regiones" />
        </div>
      </div>
    </div>
  );
}
