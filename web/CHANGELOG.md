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

## 0.3.0 — 2026-11-14 (v2: landing por WhatsApp + seguimiento + panel de pedidos)

**Eliminado (enfoque aprobado por Luckas: sin carrito ni pedido web):**
- `/carrito`, `/pedido`, `/pedido/listo`, `/catalogo` y `/producto/[slug]`;
- el carrito en localStorage, los botones "Agregar", el botón flotante al carrito y el toast;
- los embeds de TikTok y el iframe del mapa. La CSP queda en `frame-src 'none'`.

**Landing:**
- **Titular:** "¿Lo viste en TikTok? Pídelo por WhatsApp.".
- **Canales:** botones Santiago (verde, texto negro) y Regiones (verde oscuro, texto blanco) con las plantillas (a) y (b), y barra fija de canales al pasar el hero.
- **Vitrina "Lo nuevo de los videos":**
  - "Pedir este por WhatsApp" abre la hoja "¿Dónde lo quieres?" con la plantilla (c);
  - el canal elegido se recuerda en sessionStorage durante la visita.
- **Secciones:** cómo comprar desde regiones (5 pasos), retiro en Santiago, cómo llegar, canal oficial, "sigue tu pedido", Familia Emporio y preguntas.
- **Números:** de `GET /api/config-publica`, con respaldo en `NEXT_PUBLIC_WHATSAPP_SANTIAGO` y `NEXT_PUBLIC_WHATSAPP_REGIONES`. Si Regiones es null, se usa Santiago y se marca "dato por confirmar".
- **Marca de demo:** `*[PRUEBA · DEMO]*` en los mensajes de la landing.

**Seguimiento:**
- `/seguimiento/[token]` muestra el estado con su color, la línea de tiempo, la agencia y el n.º de seguimiento con botón "Copiar número", y el botón de WhatsApp del canal.
- Lleva `noindex` y `no-referrer`, y responde un 404 amable.

**Panel:**
- Pantallas Hoy, Nuevo pedido (catálogo + línea libre, descuento, envío, total manual), Ficha (transiciones, confirmación de pago, datos de despacho, cancelación con motivo, mensajes con el sugerido resaltado, historial), Editar, Buscar con ficha de cliente, Resumen + CSV (admin), Productos (admin) y Usuarios (admin; la clave se muestra una vez).
- El 403 se maneja con un mensaje para el vendedor.

**Verificado en una instancia aislada (clave de prueba fuera del repo, ya borrada):**
- Login.
- Pedido DC-0110 registrado: cotizado → esperando pago → pagado (exige la confirmación) → preparando → despachado (agencia y n.º). El mensaje sugerido pasa a "Avisar despacho".
- La página de seguimiento muestra "En camino" con su n.º.
- El vendedor no ve Resumen, Productos ni Usuarios, y el CSV le responde 403.
- Los mensajes (a), (b) y (c) decodificados calzan con la plantilla.

**Lighthouse móvil del inicio:** rendimiento 90–100 en las corridas finales (78 en la más lenta) y accesibilidad 100 tras corregir el contraste de las etiquetas.

## 0.4.0 — 2026-11-15 (ronda 3: QA ronda 2 y seguridad ronda 2)

**Panel:**
- **B-1:**
  - "Despachar" desde Hoy (agencia, n.º y comuna, directo a despachado).
  - Página `/admin/ajustes` de solo lectura, que solo ve el admin.
  - Resumen con tablas por comuna y por persona; "Vendido" cuenta solo lo pagado o posterior.
- **I-8:** sesión resuelta en el servidor y "Hoy" como Server Component (sin fetch inicial ni polling).
- **B9:** pantalla "Cambia tu clave" obligatoria con clave temporal.
- **Ficha:**
  - **M-16:** un solo botón principal (según la forma de pago).
  - **Otras mejoras:**
    - "Generar link nuevo" del seguimiento (admin);
    - aviso de stock al pagar;
    - atajos `?accion=` desde Hoy;
    - `<dl>` válido (**M-13**).
- **Nuevo pedido:**
  - **M-5:** precio editable en líneas del catálogo.
  - **M-1:** el "cliente conocido" se limpia al cambiar el teléfono.
  - **M-10:** "¿Cuál?" para la agencia "Otra", que se anota en la nota interna.
  - **M-15:** canal por defecto del usuario.
- **Hoy:**
  - **M-9:** "Recordar pago" y "Cancelar" separados; bloque "Cotizados hoy".
  - **M-4:** plurales correctos ("hace 1 día").
  - **M-8:** espacio bajo el botón flotante.
- **Menú (M-8):** en dos filas en vez de un scroll escondido; saludo en su propia línea.
- **Usuarios (M-15):** canal por defecto editable.
- **Productos:**
  - **I-2:** "Mostrar en la vitrina" + orden, arriba; sin "quedan pocas"; stock en "Más opciones"; marca "precio de promo del video" + nota.
  - **B12/M-3:** sin link a `/producto/…`.

**Landing:**
- **B-2:** placa "PROMO DEL VIDEO" con la nota de cómo pedirla, sin "detalle"; el mensaje (c) dice "Precio promo del video".
- **I-2:** la vitrina sale de `GET /api/vitrina`, con respaldo a lo más reciente.
- **I-7:** en móvil el logo va al lado del saludo y hay una línea "📍 Alameda 2432 · Metro U. Latinoamericana" sobre los botones.
- **M-12:** el JSON-LD ya no publica el horario por confirmar.

**Seguimiento:**
- **I-5:** un 429 muestra "Hay muchas consultas" (sin 500).
- **Seguridad M4:** se reenvía la IP del cliente.
- Se quitó "Tu carrito sigue guardado" de `error.tsx`.

**Imágenes (I-6):** recortes renombrados con hash (`tt-<id>-<hash8>.webp`), con el mapeo en `lib/imagenes.ts`.

**Verificado (instancia aislada con la API 2.1; clave de prueba fuera del repo, ya borrada):**
- Landing a 375×812 y 360×640: la dirección y los dos botones de canal quedan en la primera pantalla.
- Vitrina desde `/api/vitrina`. Se marcó la esfera (promo) en la vitrina desde el panel y salió con "PROMO DEL VIDEO", su nota y "Precio promo del video, publicado el 13/11: $500" en el mensaje (c).
- "Despachar" desde Hoy (pagado → despachado directo; sugerido "Avisar despacho"), Ajustes, Resumen con comunas y personas, y precio editado en una línea del catálogo ($3.000 en vez de $4.000).
- Un vendedor nuevo entra y ve solo "Cambia tu clave"; la API responde 403 `debe_cambiar_clave` hasta que la cambia.
- Seguimiento con 40 consultas seguidas: 404 amable y luego "Hay muchas consultas" (sin 500).
- Lighthouse móvil:

  | Página | Rendimiento | Accesibilidad |
  |---|---|---|
  | Inicio | 89–96 | 100 |
  | Hoy (con sesión) | 81–83 | 100 |
  | Nuevo pedido | 82–98 | 100 |
