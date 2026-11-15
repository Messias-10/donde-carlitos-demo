import type { Metadata } from "next";
import { AdminShell } from "@/components/admin/Sesion";
import { getPanel } from "@/lib/api-panel";
import type { SesionPanel } from "@/lib/tipos";

export const metadata: Metadata = {
  title: "Pedidos del WhatsApp",
  robots: { index: false, follow: false },
};

export default async function LayoutAdmin({ children }: { children: React.ReactNode }) {
  // La sesión se resuelve en el servidor con la cookie httpOnly: sin fetch inicial en el navegador.
  const r = await getPanel<SesionPanel>("/api/admin/sesion");
  return <AdminShell sesionInicial={r.ok ? r.data : null}>{children}</AdminShell>;
}
