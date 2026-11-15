import { notFound } from "next/navigation";
import { EditarPedido } from "@/components/admin/EditarPedido";

export default async function AdminEditarPedido({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  if (!/^\d{1,12}$/.test(id)) notFound();
  return <EditarPedido id={id} />;
}
