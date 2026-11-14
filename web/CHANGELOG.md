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

## 0.2.0 — 2026-11-14 (Fase 4, correcciones de la ronda 1)

**Seguridad:**
- **M2:** la imagen de ejecución va sin npm, npx, corepack ni yarn, y con `zlib` actualizado. Trivy da 0 hallazgos.
- **B2:** HSTS condicionado a producción con https.
- **B4:** sandbox en el iframe de TikTok, probado.
- **B3:** documentado como excepción en `PLAN.md` §7.

**Veracidad:**
- **I-2:** etiqueta "mínimo por confirmar" cuando `mayor_minimo_confirmado` es false (tarjeta, ficha y carrito).
- **I-3:** "Precio actualizado por la tienda el {fecha}" cuando hay `precio_actualizado_en`.
- **`oferta_hasta`:** "Oferta hasta el dd/mm" en tarjeta y ficha.
- **M-7:** el texto "[EJEMPLO]" de las descripciones se muestra como la etiqueta gris "dato de ejemplo", y "[CONFIRMAR…]" como "por confirmar".

**Panel (I-4):**
- Nueva ruta `/admin/productos/nuevo` ("Agregar producto") con nombre, categoría, precio, mayor y mínimo, stock, link de TikTok y foto desde la cámara.
- La foto se achica en el celular a 1280 px y se sube a `POST /api/admin/productos/{id}/imagen`.
- En la edición, precio y stock van arriba y el resto en "Más opciones". Se agregan "mínimo confirmado" y "oferta válida hasta".
- **M-11:** teléfono con formato legible.

**Landing:**
- **I-5:** la bajada dice qué venden y el logo es de 110 px.
- **I-8:** rendimiento (ver `PLAN.md` §7).
- **M-3:** sin chip "Navidad" duplicado.
- **M-5:** táctiles de 44 px o más.
- **M-6:** oro de texto `#8A5E0F`.
- **M-9:** nombre completo en el header.
- **M-12:** recortes nuevos, centrados en el producto y sin caras. Cinco productos quedan con la placa genérica (ver `CREDITOS.md`).
- **M-13:** sin unidad redundante.
- **M-14:** sin hueco en las tarjetas del carrusel.
- **M-16:** sin preload duplicado del logo.

**Hecho antes del cambio de enfoque (queda en el código, sin seguir avanzando):**
- **I-6:** confirmación en `/pedido/listo` con sessionStorage.
- **M-4:** mensaje del nombre.
- **M-8:** número de pedido sin cortes.
- Mensajes para `pedidos_abiertos` y para el tope de 100 unidades por línea.
