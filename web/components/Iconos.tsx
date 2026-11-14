// Íconos SVG en línea (sin librerías). Todos decorativos: aria-hidden.
import type { SVGProps } from "react";

type P = SVGProps<SVGSVGElement> & { size?: number };

function base({ size = 24, ...rest }: P) {
  return {
    width: size,
    height: size,
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 2,
    strokeLinecap: "round" as const,
    strokeLinejoin: "round" as const,
    "aria-hidden": true,
    focusable: false,
    ...rest,
  };
}

export function IconoWhatsApp({ size = 24, ...rest }: P) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden focusable={false} {...rest}>
      <path
        fill="currentColor"
        d="M12 2a9.9 9.9 0 0 0-8.5 15l-1.4 5 5.2-1.4A10 10 0 1 0 12 2Zm0 18.2a8.2 8.2 0 0 1-4.2-1.1l-.3-.2-3.1.8.8-3-.2-.3A8.2 8.2 0 1 1 12 20.2Zm4.5-6.1c-.2-.1-1.5-.7-1.7-.8-.2-.1-.4-.1-.6.1l-.8 1c-.1.2-.3.2-.5.1a6.7 6.7 0 0 1-3.3-2.9c-.3-.4.2-.4.7-1.4.1-.2 0-.3 0-.4l-.8-1.8c-.2-.5-.4-.4-.6-.4h-.5a1 1 0 0 0-.7.3 3 3 0 0 0-.9 2.2 5.2 5.2 0 0 0 1.1 2.7 11.8 11.8 0 0 0 4.5 4c1.7.7 2.4.8 3.2.6.5-.1 1.5-.6 1.7-1.2.2-.6.2-1.1.2-1.2-.1-.1-.3-.2-.5-.3Z"
      />
    </svg>
  );
}

export function IconoCarrito(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M3 4h2l2.2 10.5a1.5 1.5 0 0 0 1.5 1.2h8.6a1.5 1.5 0 0 0 1.5-1.1L21 8H6.2" />
      <circle cx="9.5" cy="19.5" r="1.4" />
      <circle cx="17" cy="19.5" r="1.4" />
    </svg>
  );
}

export function IconoLupa(p: P) {
  return (
    <svg {...base(p)}>
      <circle cx="11" cy="11" r="6.5" />
      <path d="m16 16 4.5 4.5" />
    </svg>
  );
}

export function IconoMenu(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M4 7h16M4 12h16M4 17h16" />
    </svg>
  );
}

export function IconoCerrar(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M6 6l12 12M18 6 6 18" />
    </svg>
  );
}

export function IconoPlay({ size = 24, ...rest }: P) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden focusable={false} {...rest}>
      <path fill="currentColor" d="M8 5.5v13a1 1 0 0 0 1.5.9l10.2-6.5a1 1 0 0 0 0-1.8L9.5 4.6A1 1 0 0 0 8 5.5Z" />
    </svg>
  );
}

export function IconoPin(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M12 21s-6.5-6.2-6.5-11.2a6.5 6.5 0 0 1 13 0C18.5 14.8 12 21 12 21Z" />
      <circle cx="12" cy="9.8" r="2.3" />
    </svg>
  );
}

export function IconoReloj(p: P) {
  return (
    <svg {...base(p)}>
      <circle cx="12" cy="12" r="8.5" />
      <path d="M12 7.5V12l3 2" />
    </svg>
  );
}

export function IconoMetro(p: P) {
  return (
    <svg {...base(p)}>
      <rect x="5" y="3.5" width="14" height="14" rx="3" />
      <path d="M5 11h14M8.5 14.5h.01M15.5 14.5h.01M8 20.5l2-3M16 20.5l-2-3" />
    </svg>
  );
}

export function IconoEscudo(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M12 3 5 6v5.5c0 4.4 3 8 7 9.5 4-1.5 7-5.1 7-9.5V6l-7-3Z" />
      <path d="m9 12 2 2 4-4" />
    </svg>
  );
}

export function IconoCaja(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M3.5 7.5 12 3l8.5 4.5v9L12 21l-8.5-4.5v-9Z" />
      <path d="M3.5 7.5 12 12l8.5-4.5M12 12v9" />
    </svg>
  );
}

export function IconoTienda(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M4 9.5V20h16V9.5M3 9.5 5 4h14l2 5.5c0 1.4-1.1 2.5-2.5 2.5S16 10.9 16 9.5c0 1.4-1.1 2.5-2.5 2.5h-3C9.1 12 8 10.9 8 9.5 8 10.9 6.9 12 5.5 12S3 10.9 3 9.5Z" />
      <path d="M10 20v-5h4v5" />
    </svg>
  );
}

export function IconoCopiar(p: P) {
  return (
    <svg {...base(p)}>
      <rect x="8.5" y="8.5" width="11" height="11" rx="2" />
      <path d="M15.5 8.5V6a1.5 1.5 0 0 0-1.5-1.5H6A1.5 1.5 0 0 0 4.5 6v8A1.5 1.5 0 0 0 6 15.5h2.5" />
    </svg>
  );
}

export function IconoMas(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M12 5v14M5 12h14" />
    </svg>
  );
}

export function IconoMenos(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M5 12h14" />
    </svg>
  );
}

export function IconoFlechaIzq(p: P) {
  return (
    <svg {...base(p)}>
      <path d="M15 5l-7 7 7 7" />
    </svg>
  );
}

// ---------------------------------------------------------------- categorías

export function IconoCategoria({ slug, ...p }: P & { slug: string }) {
  const b = base({ strokeWidth: 1.6, ...p });
  switch (slug) {
    case "navidad":
      return (
        <svg {...b}>
          <path d="M12 2.5 13 5l2.6.2-2 1.7.6 2.6L12 8.1l-2.2 1.4.6-2.6-2-1.7L11 5l1-2.5Z" />
          <path d="M12 10 7 16h3l-4 4.5h12L14 16h3l-5-6ZM12 20.5v1.5" />
        </svg>
      );
    case "juguetes-y-virales":
      return (
        <svg {...b}>
          <circle cx="12" cy="13" r="6.5" />
          <circle cx="7" cy="6.5" r="2.5" />
          <circle cx="17" cy="6.5" r="2.5" />
          <path d="M9.5 12h.01M14.5 12h.01M10 15.5c1.2 1 2.8 1 4 0" />
        </svg>
      );
    case "hogar":
    case "hogar-y-organizacion":
      return (
        <svg {...b}>
          <path d="M3.5 11 12 4l8.5 7M5.5 9.5V20h13V9.5" />
          <path d="M10 20v-5.5h4V20" />
        </svg>
      );
    case "escolar":
      return (
        <svg {...b}>
          <path d="M4 5.5h6a2 2 0 0 1 2 2V20a2 2 0 0 0-2-2H4V5.5ZM20 5.5h-6a2 2 0 0 0-2 2V20a2 2 0 0 1 2-2h6V5.5Z" />
        </svg>
      );
    case "textil":
      return (
        <svg {...b}>
          <path d="M8.5 3.5 4 6l1.5 4.5L7.5 10v10.5h9V10l2 .5L20 6l-4.5-2.5c-.5 1.5-1.9 2.5-3.5 2.5s-3-1-3.5-2.5Z" />
        </svg>
      );
    case "electro":
      return (
        <svg {...b}>
          <path d="M13 2.5 5.5 13.5H11l-1 8 7.5-11H12l1-8Z" />
        </svg>
      );
    case "cuidado-personal":
      return (
        <svg {...b}>
          <path d="M9 8.5h6v12a1 1 0 0 1-1 1h-4a1 1 0 0 1-1-1v-12ZM10.5 8.5V5.5h3v3M11 2.5h4" />
          <path d="M9 13h6" />
        </svg>
      );
    case "cotillon":
      return (
        <svg {...b}>
          <path d="M4 20.5 9 8l7 7-12 5.5ZM14 4.5l.5 2M18 6l-1.5 1.5M19.5 10.5l-2 .5M13 9.5c1-2 3-3 5-2.5" />
        </svg>
      );
    default:
      return (
        <svg {...b}>
          <path d="M3.5 7.5 12 3l8.5 4.5v9L12 21l-8.5-4.5v-9Z" />
          <path d="M3.5 7.5 12 12l8.5-4.5M12 12v9" />
        </svg>
      );
  }
}
