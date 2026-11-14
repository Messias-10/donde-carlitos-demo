import type { Metadata } from "next";
import { FormPedido } from "@/components/compra/FormPedido";
import { opcionesEntrega } from "@/lib/api-servidor";
import type { OpcionesEntrega } from "@/lib/tipos";

export const metadata: Metadata = {
  title: "Tus datos",
  robots: { index: false },
};

export default async function PaginaPedido() {
  let opciones: OpcionesEntrega | null = null;
  try {
    opciones = await opcionesEntrega();
  } catch {
    opciones = null;
  }
  return (
    <div className="bg-crema px-4 pb-20 pt-6">
      <div className="mx-auto max-w-6xl">
        <FormPedido opciones={opciones} />
      </div>
    </div>
  );
}
