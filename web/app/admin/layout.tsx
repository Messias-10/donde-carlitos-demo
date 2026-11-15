import type { Metadata } from "next";
import { AdminShell } from "@/components/admin/Sesion";

export const metadata: Metadata = {
  title: "Pedidos del WhatsApp",
  robots: { index: false, follow: false },
};

export default function LayoutAdmin({ children }: { children: React.ReactNode }) {
  return <AdminShell>{children}</AdminShell>;
}
