import type { Metadata } from "next";
import { PedidoListo } from "@/components/compra/PedidoListo";

export const metadata: Metadata = {
  title: "Pedido listo",
  robots: { index: false },
};

export default async function PaginaPedidoListo({ searchParams }: { searchParams: Promise<{ n?: string | string[] }> }) {
  const sp = await searchParams;
  const n = (Array.isArray(sp.n) ? sp.n[0] : sp.n) ?? "";
  return (
    <div className="bg-crema px-4 pb-20 pt-6">
      <PedidoListo numero={n.slice(0, 12)} />
    </div>
  );
}
