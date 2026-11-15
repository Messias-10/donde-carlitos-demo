import { notFound } from "next/navigation";
import { Suspense } from "react";
import { FichaPedido } from "@/components/admin/FichaPedido";

export default async function AdminPedido({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  if (!/^\d{1,12}$/.test(id)) notFound();
  return (
    <Suspense fallback={<div className="h-60 animate-pulse rounded-[12px] bg-blanco" aria-hidden />}>
      <FichaPedido id={id} />
    </Suspense>
  );
}
