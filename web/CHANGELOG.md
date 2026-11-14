# Changelog — web

## 0.1.0 — 2026-11-14 (Fase 3, primera iteración)

**Stack:**
- Next.js 16.4.6 (App Router, Turbopack), React 19.3.1 y TypeScript 5.9.3.
- Tailwind 4.3.5 y ESLint 9.39.5 con eslint-config-next 16.4.6.
- Versiones exactas (`save-exact`) y `package-lock.json` generado en `node:26.13.1-alpine`.
- `overrides` fuerza `minimatch` 11.0.8 dentro de `@typescript-eslint/typescript-estree` (GHSA-g3hv-75wg-m249, moderada, solo en desarrollo). Con eso, `npm audit` queda en 0.

**Páginas:**
- `/`: landing con las 14 secciones de `propuesta/arquitectura-landing.md` y el copy de `copy-landing.md`.
- `/catalogo`: búsqueda, filtros, orden y paginación. Funciona sin JS.
- `/producto/[slug]`: ficha con oferta, mayor, fuente del precio y TikTok al tocar.
- `/carrito`: carrito en localStorage, con estimado y aviso de escalera.
- `/pedido`: formulario validado → `POST /api/pedidos` → confirmación con `whatsapp_url` y `mensaje_whatsapp` de la API, más "Copiar mensaje".
- `/admin`: login, pedidos por estado, detalle con cambio de estado e historial, y edición de productos.
- Además: `/healthz`, `robots.txt`, `sitemap.xml`, 404 y página de error.

**Seguridad:**
- CSP con nonce por solicitud en `proxy.ts`, con `frame-src` limitado a TikTok y Google Maps.
- `X-Content-Type-Options`, `Referrer-Policy`, `X-Frame-Options`, `COOP` y `Permissions-Policy` en `next.config.ts`.
- Panel con `noindex` y `no-store`.
- Sin tokens en el navegador: la sesión es la cookie httpOnly de la API.
- El único `dangerouslySetInnerHTML` es el JSON-LD estático del negocio, con `<` escapado.

**SEO:**
- Metadata y Open Graph (`og.jpg` de 1200×630, el logo sobre negro).
- JSON-LD `Store` con dirección, horario de temporada y redes oficiales.

**Imágenes:**
- 26 recortes limpios del cuadro de video en `public/productos/`, con el logo y los íconos.
- Origen anotado en `public/CREDITOS.md`.

**Docker:**
- `Dockerfile` multi-stage con salida standalone y usuario `node`. Lleva `HEALTHCHECK` a `/healthz`.
- `.dockerignore` excluye `node_modules`, `.next` y `.env*`.

**Verificado:**
- `npm run lint` y `tsc --noEmit` limpios; `npm run build` OK; `npm audit` sin hallazgos.
- Pedido de punta a punta contra la API real en Docker: carrito → pedido DC-0008 → link `wa.me` decodificado igual a la plantilla. No se envió ningún mensaje.
- Admin: login, confirmar pedido y editar stock.
- A 375 px no hay desborde horizontal.
- Capturas en `qa/capturas/`.

**Pendiente:**
- Lighthouse móvil: no se pudo correr aquí (no hay Chrome en la imagen alpine); queda para qa-revisor.
- Fotos de productos sin recorte: muestran la placa por categoría.
