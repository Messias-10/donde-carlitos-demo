import { notFound } from "next/navigation";
import { DetallePedido } from "@/components/admin/DetallePedido";

export default async function AdminPedido({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  if (!/^\d{1,12}$/.test(id)) notFound();
  return <DetallePedido id={id} />;
}
