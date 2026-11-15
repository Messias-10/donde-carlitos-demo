import { Hoy } from "@/components/admin/Hoy";

export default async function AdminHoy({ searchParams }: { searchParams: Promise<{ canal?: string | string[] }> }) {
  const sp = await searchParams;
  const c = Array.isArray(sp.canal) ? sp.canal[0] : sp.canal;
  return <Hoy canal={c === "santiago" || c === "regiones" ? c : "todos"} />;
}
