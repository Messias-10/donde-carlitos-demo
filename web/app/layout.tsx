import type { Metadata, Viewport } from "next";
import { Anton, Inter, Kaushan_Script } from "next/font/google";
import { headers } from "next/headers";
import "./globals.css";

const anton = Anton({ weight: "400", subsets: ["latin"], variable: "--font-anton", display: "swap" });
const kaushan = Kaushan_Script({
  weight: "400",
  subsets: ["latin"],
  variable: "--font-kaushan",
  display: "swap",
  preload: false,
});
const inter = Inter({ subsets: ["latin"], variable: "--font-inter", display: "swap" });

const SITIO = (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000").replace(/\/+$/, "");

const DESCRIPCION =
  "Juguetes virales, Navidad, hogar y escolar, al detalle y por mayor. Pide por WhatsApp para retiro en Santiago o con envío a regiones. Alameda 2432, a pasos del Metro Unión Latinoamericana.";

export const metadata: Metadata = {
  metadataBase: new URL(SITIO),
  title: {
    default: "Donde Carlitos Emporio · Pedidos por WhatsApp · Alameda 2432",
    template: "%s · Donde Carlitos Emporio",
  },
  description: DESCRIPCION,
  applicationName: "Donde Carlitos Emporio",
  openGraph: {
    type: "website",
    locale: "es_CL",
    siteName: "Donde Carlitos Emporio",
    title: "¿Lo viste en TikTok? Pídelo por WhatsApp · Donde Carlitos Emporio",
    description: DESCRIPCION,
    url: "/",
    images: [{ url: "/og.jpg", width: 1200, height: 630, alt: "Logo de Donde Carlitos Emporio" }],
  },
  twitter: {
    card: "summary_large_image",
    title: "¿Lo viste en TikTok? Pídelo por WhatsApp · Donde Carlitos Emporio",
    description: DESCRIPCION,
    images: ["/og.jpg"],
  },
  formatDetection: { telephone: false, address: false, email: false },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  themeColor: "#0b0b0b",
};

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  // Leer las cabeceras hace dinámico el render: necesario para el nonce de la CSP.
  await headers();
  return (
    <html lang="es-CL" className={`${anton.variable} ${kaushan.variable} ${inter.variable}`}>
      <body className="min-h-dvh antialiased">{children}</body>
    </html>
  );
}
