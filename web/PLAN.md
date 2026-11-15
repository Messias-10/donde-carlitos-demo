# Plan del front — Donde Carlitos Emporio

*constructor-web · Fase 3 · 2026-11-14*

**Insumos:**
- skill `estilo-carlitos`;
- `propuesta/arquitectura-landing.md`, `copy-landing.md` y `mensaje-pedido-whatsapp.md`;
- `research/oportunidades.md`;
- la API en `api/internal/http/*.go`. `docs/openapi.yaml` todavía no existía al empezar, así que el contrato se leyó del código.

**La idea:** el local de remate que el cliente ya conoce, pero ordenado. Negro y oro para la marca; flúor solo en los precios.

## 1. Stack y versiones (fijadas)

| Paquete | Versión | Nota |
|---|---|---|
| next | 16.4.6 | fijada en `docs/BITACORA.md` |
| react / react-dom | 19.3.1 | |
| typescript | 5.9.3 | No se usa TS 7: Next usa la API JS de `typescript` para chequear tipos |
| tailwindcss / @tailwindcss/postcss | 4.3.5 | config CSS-first (`@theme`) |
| eslint / eslint-config-next | 9.39.5 / 16.4.6 | No se usa ESLint 10: `eslint-plugin-react` admite hasta ^9.7. Next 16 ya no trae `next lint`; se corre `eslint .` |
| Node (imagen) | node:26.13.1-alpine | npm se corre dentro de esta imagen, con `web/` montado |

- **Sin librerías de UI, de estado ni de fetch:** carrito con `useSyncExternalStore` y localStorage; fetch nativo.
- **Fuentes** con `next/font/google` (se autoalojan en el build): Anton, Kaushan Script, Inter.

## 2. Páginas

| Ruta | Render | Qué hace |
|---|---|---|
| `/` | servidor, dinámico | Landing con las 14 secciones de `arquitectura-landing.md` |
| `/catalogo` | servidor, dinámico | Búsqueda `q`, chips de categoría, filtros "Ofertas", "Recién llegado" y "Navidad", orden y paginación por links. Funciona sin JS |
| `/producto/[slug]` | servidor, dinámico | Ficha: imagen 4:5, precios (detalle, oferta, mayor), disponibilidad, descripción, fuente del precio, cantidad + Agregar y video de TikTok (miniatura que carga el reproductor al tocar) |
| `/carrito` | cliente | Líneas, cantidades, aviso de escalera al precio por mayor, total estimado y 3 pasos en chico |
| `/pedido` | servidor + cliente | Formulario de datos con validación → `POST /api/pedidos` → confirmación con número, "Enviar pedido por WhatsApp", "Copiar mensaje" y canal oficial |
| `/admin` | cliente | Login, pedidos por estado (pestañas) y salida |
| `/admin/pedidos/[id]` | cliente | Detalle, cambio de estado con nota, cancelación con confirmación en línea, historial y "Abrir chat con el cliente" |
| `/admin/productos` | cliente | Lista con buscador |
| `/admin/productos/[id]` | cliente | Edición: precio, oferta, etiqueta, mayor, stock, umbral, temporada, llegada, TikTok y visible |
| `/healthz` | route handler | Healthcheck del contenedor |
| `robots.txt`, `sitemap.xml` | metadata routes | `/admin` va con `disallow` y `noindex` |

## 3. Componentes

- **Estructura:**
  - `Header`: logo de 44 px, buscador desplegable, menú de anclas y carrito con contador flúor;
  - `Footer`;
  - `BotonFlotante`: WhatsApp sin carrito y "Ver mi pedido (n) · $total" con carrito. Se oculta en `/carrito`, `/pedido` y `/admin`.
- **Producto:**
  - `TarjetaProducto`, `Precio` (placa flúor, tachado solo si `precio_oferta`), `Badges` (máximo 2 por tarjeta), `Disponibilidad` (punto + texto, nunca solo color);
  - `BotonAgregar` (microanimación del contador), `Toast`;
  - `VisualProducto`: foto si existe; si no, una placa propia por categoría con la etiqueta "Foto en el video". Ver §6.
- **Medios con carga al tocar:** `VideoTikTok` (póster → iframe `tiktok.com/player/v1/{id}`) y `MapaCarga` (botón → iframe de Google Maps).
- **Secciones de la landing:**
  - `Hero`, `BannerRemate`, `FilaProductos`, `LoVisteEnTikTok`, `Categorias`, `PaNegocio`, `PasosPedido`;
  - `RetiroEnvio`, `ComoLlegar`, `CanalOficial`, `FamiliaEmporio`, `Preguntas` (`<details>`).
- **Marcas de dato:** `PorConfirmar` ("dato por confirmar", gris chico) y `Ejemplo` ("ejemplo").
- **Compra:** `Carrito`, `FormPedido`, `Confirmacion`.
- **Admin:** `AdminShell` (sesión, salida), `ListaPedidos`, `DetallePedido`, `ListaProductos`, `FormProducto`.

## 4. Design tokens (Tailwind v4 `@theme`, de estilo-carlitos §2)

```
negro #0B0B0B · carbon #1A1A1A · oro #E2B33C · oro-claro #FCE48F · oro-profundo #9A6A12
blanco #FFFFFF · crema #FFF8E1 · fluor #F5E800 · magenta #D6125A · rojo #C8102E
whatsapp #25D366 · whatsapp-osc #0B6E4F · gris #6B6B6B · gris-claro #B8B8B8
```

**Tipografía:**
- `font-display` (Anton) para titulares en mayúsculas y precios;
- `font-script` (Kaushan) solo para acentos de 2–3 palabras sobre negro;
- `font-sans` (Inter) para el resto.

**Escala (móvil):**
- cuerpo 16 px;
- ficha 15 px;
- H1 hero 44 px, que sube a 64 px en desktop;
- H2 de sección 32 px;
- precio en la tarjeta 30 px.

**Forma:**
- tarjetas con radio de 12 px;
- badges en píldora;
- CTA de WhatsApp con radio de 14 px y 52 px de alto;
- botones de 44 px o más.

**Superficies:**
- header, hero, banner, canal oficial y footer en negro;
- catálogo y fichas en crema con tarjetas blancas.

**Movimiento** (solo como respuesta a una acción, y se respeta `prefers-reduced-motion`):
1. el contador del carrito hace "pop" al agregar;
2. el toast entra desde abajo;
3. la placa de precio de oferta se inclina levemente (-2°), como una etiqueta pegada.

**Revisión contra frontend-design:** el plan no cae en los clichés de la lista porque la paleta y las fuentes vienen del logo y del local, no de una plantilla.
- El elemento memorable es la **placa de precio flúor con forma de etiqueta**. Todo lo demás queda quieto.
- Los números (1-2-3) se usan solo en "Pedir es así de fácil", que sí es una secuencia.
- **Cambié** las "fade-up por sección", que había considerado, por cero animaciones de entrada.
- **Cambié** las etiquetas en mayúsculas sobre cada titular, que también había considerado, por nada: los titulares ya van en Anton mayúscula.

## 5. Datos

**Las llamadas pasan por el mismo origen:**
- **Desde el navegador:** el navegador llama a `/api/*`, y los `rewrites` de `next.config.ts` las mandan a `API_INTERNAL_URL` (por defecto `http://api:8080`; el valor se fija en el build). Así no hay CORS y la cookie `dc_sesion` (httpOnly, SameSite=Strict) funciona tal cual.
- **Desde el servidor (SSR):** el render en servidor hace `fetch(API_INTERNAL_URL + '/api/...')` con `cache: 'no-store'`, porque el stock cambia. Ese valor se lee en tiempo de ejecución.

**Qué endpoint usa cada parte:**

| Endpoint | Dónde |
|---|---|
| `GET /api/categorias` | landing (grilla) y catálogo (chips) |
| `GET /api/productos?oferta&recien_llegado&temporada&categoria&q&orden&page&page_size` | landing y catálogo |
| `GET /api/productos/{slug}` | ficha |
| `GET /api/opciones-entrega` | `/pedido` (regiones y agencias, con la marca `agencias_es_ejemplo`) |
| `POST /api/pedidos` | checkout. El servidor recalcula los precios. El front muestra el **total y el mensaje que devuelve la API** |
| `/api/admin/*` | panel: `login`, `logout`, `sesion`, `pedidos`, `pedidos/{id}`, `pedidos/{id}/estado`, `productos`, `productos/{id}` (PUT reemplaza todo el producto) |

**Carrito:**
- Se guarda en `localStorage["dc_carrito_v1"]`: `{v:1, items:[{id, slug, codigo, nombre, cantidad, precio, precio_final, precio_mayor, mayor_minimo, unidad_mayor, categoria}]}`.
- Nada sensible: ni nombre ni teléfono.
- El precio que muestra es un estimado con la misma regla de la API (el menor entre detalle, oferta y mayor si se llega al mínimo). El que vale es el que devuelve la API.

**Mensaje de WhatsApp:**
- Se usan `mensaje_whatsapp` y `whatsapp_url` de la respuesta de `POST /api/pedidos`, tal cual. No se arma con el carrito.
- El front solo acepta `whatsapp_url` si tiene la forma `https://wa.me/<dígitos>?text=`; si no, lo arma con `mensaje_whatsapp`.
- Los mensajes de consulta (§5 de la plantilla) los arma el front con `NEXT_PUBLIC_WHATSAPP_PEDIDOS` (por defecto 56921145368; es público, no secreto).

**Sin mocks:** la API ya está implementada, así que se integró directo.

## 6. Decisiones

1. **Fotos de producto.** La semilla apunta a `/productos/tt-<videoId>.webp`.
   - Por indicación del orquestador, se generaron **recortes limpios del cuadro del video** desde `research/raw/tt-*.png`: sin comentarios, sin usuarios y sin la interfaz de TikTok.
   - Cada recorte se revisó a mano. Uno se descartó porque aparece un menor de edad, y otro se recortó más abajo para dejar fuera un letrero con un teléfono antiguo. Detalle en `public/CREDITOS.md`.
   - El front revisa en el servidor si el archivo existe. Si no existe, muestra una placa propia por categoría.
2. **TikTok** sin script de embed (`embed.js` es pesado y pide `script-src` externo). El póster es propio. Al tocarlo se carga el iframe `https://www.tiktok.com/player/v1/{id}`. Siempre hay un link "Abrir en TikTok".
3. **Mapa:** botón "Ver mapa aquí" que carga el iframe de Google Maps solo al tocarlo, más el link "Abrir en Google Maps".
4. **CSP con nonce por solicitud**, en `proxy.ts` (Next 16 renombró `middleware` a `proxy`).
   - Sin `unsafe-inline` en `script-src`.
   - Las demás cabeceras van en `next.config.ts`: `X-Content-Type-Options`, `Referrer-Policy`, `X-Frame-Options`, `Permissions-Policy` y `Cross-Origin-Opener-Policy`. `frame-ancestors 'none'` va dentro de la CSP.
   - **Dominios:** `frame-src` permite www.tiktok.com, www.google.com y maps.google.com. `img-src` y `font-src` quedan en `'self'`, y `connect-src` en `'self'`.
   - `style-src 'self' 'unsafe-inline'`, porque `next/image` y React usan atributos `style`.
5. **Horario:**
   - Se muestra la versión de nov 2026 con la etiqueta "dato por confirmar" y **sin** "Abierto/Cerrado" (arquitectura §1).
   - El "Hoy" se calcula en el servidor con la zona America/Santiago.
6. **"Envíos a todo Chile" no se escribe.** La tarjeta de envíos lleva `[CONFIRMAR]` visible, como indica el copy.
7. **Preguntas sin dato** ("¿Hacen cambios?" y el local de Temuco) no se publican, como dice `copy-landing.md`.
8. **La confirmación del pedido** queda solo en memoria (no se guarda en el navegador). El botón abre `wa.me` en una pestaña nueva y la confirmación sigue a la vista.
9. **Admin:** solo cliente. Si la API responde 401, se vuelve al login. Sin tokens en el navegador: la cookie es httpOnly. La edición manda el producto completo (PUT), mapeando `categoria.id` → `categoria_id` y `fuente` → `fuente_url`/`fuente_fecha`.
10. **SEO:**
    - Metadata API con title, description y OG de `copy-landing.md` §1.1.
    - Imagen OG de 1200×630: el logo sin tocar, centrado sobre negro. Así se ve bien en la vista previa cuadrada de WhatsApp.
    - JSON-LD `Store` (LocalBusiness) con datos estáticos del negocio. No usa datos de usuarios ni de la API.
11. **Docker:**
    - `output: "standalone"`, imagen multi-stage `node:26.13.1-alpine` y usuario `node`.
    - Build args: `NEXT_PUBLIC_WHATSAPP_PEDIDOS`, `NEXT_PUBLIC_SITE_URL` y `API_INTERNAL_URL`.
    - `HEALTHCHECK` a `/healthz`.

## 7. Ronda 1 de Fase 4: correcciones y excepciones

**Excepción aceptada (seguridad B3): `style-src 'self' 'unsafe-inline'`.**
- **Por qué:** `next/image` y React escriben atributos `style` en línea (posición de `fill`, `aspect-ratio`), y un atributo `style` no se puede autorizar con nonce.
- **Riesgo:** solo inyección de CSS, no de scripts. `script-src` sigue con nonce y `strict-dynamic`, sin `unsafe-inline`.
- **Revisar** si Next agrega soporte de nonce para estilos.

**HSTS (B2):** solo se envía cuando el build es de producción y `NEXT_PUBLIC_SITE_URL` empieza con `https://`. En el demo local por http no se envía.

**Iframe de TikTok (B4):** lleva `sandbox="allow-scripts allow-same-origin allow-presentation allow-popups allow-popups-to-escape-sandbox"`. Se probó que el reproductor carga y reproduce.

**Imagen Docker (M2):**
- La etapa de ejecución actualiza `zlib` y borra npm, npx, corepack y yarn.
- El healthcheck usa `wget` de busybox.
- Trivy 0.69.3: 0 vulnerabilidades de cualquier severidad.

**Rendimiento (I-8):**
- Los botones "Agregar" de las tarjetas se pintan en el servidor (`BotonAgregarRapido`), con un solo listener (`AgregarDelegado`). Antes era un componente cliente por tarjeta.
- Las imágenes llevan `sizes` reales.
- Las secciones usan `content-visibility: auto`.
- La portada muestra 4 ofertas, 4 novedades y 4 videos.
- Lighthouse móvil del inicio: 85–87 (antes 72).

## 8. v2 (0.3.0): landing por WhatsApp + seguimiento + panel de pedidos

**Enfoque aprobado por Luckas:** no hay carrito ni pedido web. Se borraron:
- `/carrito`, `/pedido`, `/pedido/listo`, `/catalogo` y `/producto/[slug]`;
- el localStorage del carrito, los botones "Agregar", el botón flotante al carrito, el toast y el delegado de "Agregar";
- el iframe de TikTok y el del mapa (ahora son links). Con eso la CSP queda en `frame-src 'none'`.

**Landing `/`** (arquitectura-landing.md §A, 11 bloques):

| Bloque | Componente | Notas |
|---|---|---|
| Header | `Header` | logo + "Sigue tu pedido" |
| Hero | `landing/Hero` | titular "¿Lo viste en TikTok? Pídelo por WhatsApp.", dos `BotonCanal` (a)/(b), barra de datos |
| Vitrina | `landing/Vitrina` | `GET /api/productos?orden=recientes&page_size=8` (no hay `/api/vitrina` en el contrato). Botón "Pedir este por WhatsApp" → `HojaCanal` (`<dialog>`, plantilla c) |
| Regiones | `ComprarDesdeRegiones` | 5 pasos + "Lo que tienes que saber" con "por confirmar" |
| Santiago | `RetiroSantiago` | 3 pasos |
| Cómo llegar | `ComoLlegar` | link a Google Maps, horario, pagos |
| Canal oficial | `CanalOficial` | uno o dos números |
| Seguimiento | `SigueTuPedido` | sin buscador por número |
| Familia | `FamiliaEmporio` | link al video, sin miniatura |
| Preguntas | `Preguntas` | `<details>` |
| Barra fija | `BarraCanales` | aparece al salir del hero (IntersectionObserver) |

**Números de WhatsApp:**
- Se toman de `GET /api/config-publica` en el servidor (`lib/numeros.ts`). Si la API falla, salen de `NEXT_PUBLIC_WHATSAPP_SANTIAGO` y `NEXT_PUBLIC_WHATSAPP_REGIONES`.
- `whatsapp_regiones: null` → se usa el de Santiago y se marca "dato por confirmar".
- Los links `wa.me` se arman en el servidor: los botones funcionan sin JS.

**Mensajes y marca de demo:**
- Las plantillas (a), (b) y (c) están en `lib/whatsapp.ts`, con el texto exacto de `mensajes-whatsapp.md`.
- La línea `*[PRUEBA · DEMO]*` va mientras `NEXT_PUBLIC_DEMO` no sea `"false"`. Por defecto va, porque `config-publica` no trae un indicador de demo.

**Seguimiento `/seguimiento/[token]`:**
- Server component que consume `GET /api/seguimiento/{token}`.
- `noindex` y `referrer: no-referrer`.
- Si el token no existe, `not-found.tsx` responde con un 404 amable y los dos botones de canal.
- El contrato **no trae ítems, total ni canal**, así que no se muestran productos. El botón de WhatsApp usa Santiago si es retiro y Regiones si es envío.

**Panel `/admin`** (cliente; la sesión es la cookie httpOnly de la API):

| Ruta | Pantalla | Rol |
|---|---|---|
| `/admin` | Hoy (`GET /api/admin/hoy`), filtro por canal, alertas con borde rojo | todos |
| `/admin/pedidos/nuevo` | Registro: canal → teléfono (cliente conocido) → nombre → productos (buscador + línea libre) → entrega/pago → total | todos |
| `/admin/pedidos/[id]` | Ficha: transiciones como botones, diálogos (pagado con confirmación y datos de pago; despachado con agencia, n.º y comuna; cancelar con motivo), mensajes d1–d8 con "Abrir en WhatsApp" y "Copiar mensaje" (el sugerido, resaltado), historial | todos |
| `/admin/pedidos/[id]/editar` | PUT, solo en cotizado / esperando pago | todos |
| `/admin/buscar` | q, canal, estado, fechas; ficha de cliente si q es un teléfono | todos |
| `/admin/resumen` | Resumen + exportar CSV (con aviso de datos personales) | admin |
| `/admin/productos` | Lo de antes | admin |
| `/admin/usuarios` | Crear (clave visible una sola vez, con copiar), activar/desactivar, clave nueva | admin |

**Detalles del panel:**
- El 403 `sin_permiso` muestra "Esto lo puede hacer solo un administrador". El menú oculta lo que no corresponde al rol.
- La API no tiene un canal por usuario: el canal del registro se recuerda en `localStorage["dc_panel_canal"]`, que es una preferencia del puesto.

**Colores de estado:** tokens `--color-st-*` en `globals.css` (estilo-carlitos §8.6). Etiquetas en `lib/estados.ts`.
