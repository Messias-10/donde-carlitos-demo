import { notFound } from "next/navigation";
import { FormProducto } from "@/components/admin/FormProducto";

export default async function AdminProducto({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  if (!/^\d{1,12}$/.test(id)) notFound();
  return <FormProducto id={id} />;
}
