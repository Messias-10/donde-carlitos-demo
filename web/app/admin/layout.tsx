import type { Metadata } from "next";
import { AdminShell } from "@/components/admin/Sesion";

export const metadata: Metadata = {
  title: "Panel de pedidos",
  robots: { index: false, follow: false },
};

export default function LayoutAdmin({ children }: { children: React.ReactNode }) {
  return <AdminShell>{children}</AdminShell>;
}
