import type { NextConfig } from "next";

// URL interna de la API en Go. Los rewrites se resuelven en el build:
// para cambiarla hay que reconstruir la imagen (build arg API_INTERNAL_URL).
const API = (process.env.API_INTERNAL_URL || "http://api:8080").replace(/\/+$/, "");

// Cabeceras fijas. La CSP va en proxy.ts porque lleva un nonce por solicitud.
const cabeceras = [
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
  { key: "X-Frame-Options", value: "DENY" },
  { key: "Cross-Origin-Opener-Policy", value: "same-origin" },
  {
    key: "Permissions-Policy",
    value: "camera=(), microphone=(), geolocation=(), payment=(), usb=()",
  },
];

// HSTS solo en producción real con https (QA/seguridad B2). En el demo local (http) no se manda.
const SITIO = process.env.NEXT_PUBLIC_SITE_URL || "";
if (process.env.NODE_ENV === "production" && SITIO.startsWith("https://")) {
  cabeceras.push({ key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains" });
}

const nextConfig: NextConfig = {
  output: "standalone",
  poweredByHeader: false,
  reactStrictMode: true,
  devIndicators: false,
  images: {
    formats: ["image/webp"],
    deviceSizes: [360, 414, 640, 768, 1080, 1280],
    imageSizes: [44, 64, 96, 128, 160, 256, 384],
    minimumCacheTTL: 86400,
  },
  async rewrites() {
    return [{ source: "/api/:path*", destination: `${API}/api/:path*` }];
  },
  async headers() {
    return [
      { source: "/:path*", headers: cabeceras },
      {
        source: "/admin/:path*",
        headers: [
          { key: "X-Robots-Tag", value: "noindex, nofollow" },
          { key: "Cache-Control", value: "no-store" },
        ],
      },
      {
        source: "/admin",
        headers: [
          { key: "X-Robots-Tag", value: "noindex, nofollow" },
          { key: "Cache-Control", value: "no-store" },
        ],
      },
    ];
  },
};

export default nextConfig;
