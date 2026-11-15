"use client";

import type { Pedido } from "@/lib/tipos";
import { FormPedido } from "./FormPedido";
import { Aviso } from "./ui";
import { useApi } from "./useApi";

/** Editar solo en cotizado o esperando pago (la API responde 409 no_editable si no). */
export function EditarPedido({ id }: { id: string }) {
  const { data, error } = useApi<Pedido>(`/api/admin/pedidos/${encodeURIComponent(id)}`);
  if (error) return <Aviso>{error}</Aviso>;
  if (!data) return <div className="h-60 animate-pulse rounded-[12px] bg-blanco" aria-hidden />;
  if (data.estado !== "cotizado" && data.estado !== "esperando_pago") {
    return <Aviso tipo="info">Este pedido ya no se puede editar: ya está pagado o más adelante.</Aviso>;
  }
  return <FormPedido base={data} />;
}
