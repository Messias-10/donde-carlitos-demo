import type { Metadata } from "next";
import { Carrito } from "@/components/compra/Carrito";

export const metadata: Metadata = {
  title: "Tu pedido",
  robots: { index: false },
};

export default function PaginaCarrito() {
  return (
    <div className="bg-crema px-4 pb-20 pt-6">
      <div className="mx-auto max-w-6xl">
        <h1 className="titular mb-5 text-[36px]">Tu pedido</h1>
        <Carrito />
      </div>
    </div>
  );
}
