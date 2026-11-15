"use client";

import { useEffect, useState } from "react";
import { IconoWhatsApp } from "@/components/Iconos";
import type { Canal } from "@/lib/tipos";

/**
 * Barra fija abajo en móvil con los dos canales (arquitectura §A11). Aparece
 * después de pasar el hero para no duplicar los botones en la primera pantalla.
 * Los href vienen armados del servidor.
 */
export function BarraCanales({
  hrefs,
  textos,
}: {
  hrefs: Record<Canal, string>;
  textos: Record<Canal, { corto: string; aria: string }>;
}) {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const hero = document.getElementById("hero");
    if (!hero || !("IntersectionObserver" in window)) {
      // Sin observer (navegadores muy viejos): se muestra siempre.
      const t = window.setTimeout(() => setVisible(true), 0);
      return () => window.clearTimeout(t);
    }
    const io = new IntersectionObserver(([e]) => setVisible(!e.isIntersecting), { threshold: 0 });
    io.observe(hero);
    return () => io.disconnect();
  }, []);

  return (
    <nav
      aria-label="Pedidos por WhatsApp"
      aria-hidden={!visible}
      className={`fixed inset-x-0 bottom-0 z-40 border-t border-oro/30 bg-negro/95 px-3 pb-[max(8px,env(safe-area-inset-bottom))] pt-2 transition-transform duration-200 lg:hidden ${
        visible ? "translate-y-0" : "pointer-events-none translate-y-full"
      }`}
    >
      <div className="mx-auto grid max-w-md grid-cols-2 gap-2">
        {(["santiago", "regiones"] as const).map((c) => (
          <a
            key={c}
            href={hrefs[c]}
            target="_blank"
            rel="noopener noreferrer"
            tabIndex={visible ? 0 : -1}
            aria-label={textos[c].aria}
            className={`flex h-[52px] items-center justify-center gap-2 rounded-[14px] text-[15px] font-extrabold ${
              c === "santiago" ? "bg-whatsapp text-negro" : "bg-whatsapp-osc text-blanco"
            }`}
          >
            <IconoWhatsApp size={20} />
            {textos[c].corto}
          </a>
        ))}
      </div>
    </nav>
  );
}
