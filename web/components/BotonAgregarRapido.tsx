// Botón "Agregar" que se pinta en el servidor (sin componente cliente por tarjeta).
// Un solo listener en AgregarDelegado.tsx lee data-agregar y suma al carrito.
// Así la portada hidrata mucho menos JS (QA I-8).
import type { ProductoParaCarrito } from "@/lib/carrito";
import { linkWhatsApp, msjAgotado } from "@/lib/whatsapp";
import { IconoWhatsApp } from "./Iconos";

export function BotonAgregarRapido({
  producto,
  agotado,
  promo,
}: {
  producto: ProductoParaCarrito;
  agotado: boolean;
  promo: boolean;
}) {
  if (agotado) {
    return (
      <a
        href={linkWhatsApp(msjAgotado(producto.nombre, producto.codigo))}
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex min-h-11 w-full items-center justify-center gap-2 rounded-[10px] border-2 border-negro bg-blanco px-3 text-[14px] font-extrabold text-negro"
      >
        <IconoWhatsApp size={18} className="text-whatsapp-osc" />
        ¿Vuelve? Pregúntanos
      </a>
    );
  }
  return (
    <button
      type="button"
      data-agregar={JSON.stringify(producto)}
      className={`inline-flex min-h-11 w-full items-center justify-center rounded-[10px] px-3 text-[15px] font-extrabold transition-transform active:scale-95 ${
        promo ? "bg-magenta text-blanco" : "bg-negro text-oro"
      }`}
    >
      {promo ? "¡Yo quiero mi promo!" : "Agregar"}
      <span className="sr-only"> {producto.nombre}</span>
    </button>
  );
}
