# QA · Ronda 1 — Donde Carlitos Emporio

*qa-revisor · Fase 4 · 2026-11-14 · stack local con `WEB_PORT=3100` (web http://localhost:3100, API http://127.0.0.1:8080). Guardado por el orquestador: el subagente no pudo escribir el archivo.*
*Capturas: `qa/capturas/ronda-1/` (prefijo `m` = móvil 375×812, `d` = escritorio 1366×900)*

## Veredicto: **listo para mostrar al dueño**, con correcciones importantes antes de la reunión

- **No hay bloqueantes.** Funcionan el arranque limpio, el pedido de punta a punta, el mensaje de WhatsApp idéntico a la plantilla, el panel con login y la resistencia a manipulación.
- **Hay que corregir I-1 a I-4 antes de ponérselo en la mano:** hoy el demo dice cosas que no son verdad (stock y "quedan pocas" inventados, mínimos `[EJEMPLO]` sin marcar, precio editado que sigue diciendo "precio del video") y al dueño le falta subir un producto nuevo.

## Resumen de pruebas

| # | Prueba | Resultado |
|---|---|---|
| 1 | Arranque limpio | ✅ 16 s; `api` y `web` `healthy`; `/healthz` 200 en 3100 y 8080 |
| 2 | Build y tests (en contenedor) | ✅ `npm ci`, lint, `tsc`, build, `npm audit` (0), `go vet`, `gofmt`, `go test ./...`, `govulncheck` v1.9.0 sin vulnerabilidades |
| 3 | Pedido E2E en 375×812 | ✅ "squishy" → 3 productos → escalera a 12 u. a $3.500 → **DC-0007 $58.000**; envío Ñuble/Chillán y retiro |
| 3b | WhatsApp decodificado contra la plantilla | ✅ idéntico byte a byte a la referencia §4 (tildes, ñ, emojis, LF, total) |
| 3c | Pedido en /admin | ✅ (instancia aislada) DC-0007, $58.000, líneas y total correctos |
| 4 | Panel | ✅ login incorrecto/correcto, Nuevo → Confirmado con nota e historial, edición de precio y stock reflejada, logout invalida la cookie, 429 tras 5 fallos, cookie HttpOnly + SameSite=Strict |
| 4b | Sin sesión | ✅ `/api/admin/*` → 401; PATCH → 403/401 |
| 5 | Manipulación | ✅ ver tabla |
| 6 | Lighthouse móvil | ⚠ inicio 72 |
| 7 | Accesibilidad | ⚠ footer y chips < 44 px; un contraste de 4,44 |
| 8 | Marca | ✅ en general |
| 9 | Veracidad | ⚠ I-1, I-2, I-3 |
| 10 | 5 segundos | ⚠ no se entiende qué venden (I-5) |
| 11 | Dueño | ⚠ no puede cargar productos (I-4) |

**Manipulación**

| Envío | Respuesta |
|---|---|
| `total:1`, `precio_unitario:1` | 201, total recalculado a $8.000 ✅ |
| `cantidad` -5 / 0 / 99999999 | 400 "entre 1 y 500" ✅ |
| `cantidad` 1.5 | 400 ✅ |
| Producto inactivo / agotado / inexistente / sin stock | 409 con motivo ✅ |
| Campos extra (`es_demo`, `tipo_precio`) | 400 ✅ |
| `<script>` en nombre, teléfono extranjero, `items:[]`, región inventada | 400 ✅ |
| Sin Origin / Origin ajeno | 403 ✅ |
| `text/plain` / 20 KB | 415 / 413 ✅ |
| 11.º pedido en 10 min | 429 ✅ |

**Lighthouse móvil 13.6.0**

| Página | Rend. | Acc. | BP* | SEO | LCP | TBT | CLS |
|---|---|---|---|---|---|---|---|
| `/` | **72** | 97 | 74 | 100 | 1,3 s | **2.270 ms** | 0 |
| `/catalogo` | 83 | 99 | 74 | 100 | 1,3 s | 670 ms | 0 |
| ficha squishy | 89 | 100 | 74 | 100 | 3,2 s | 240 ms | 0 |

\*BP baja por HTTP en local; no es hallazgo.

## Bloqueantes
Ninguno.

## Importantes

| # | Hallazgo | Sugerencia | Dueño |
|---|---|---|---|
| I-1 | La semilla da stock inventado y la tarjeta muestra "Quedan pocas unidades" (p. ej. calendario de adviento = 2): urgencia falsa. El `[EJEMPLO]` solo está en la descripción. | Stocks sobre el umbral en `semilla.json` (salvo la repisa agotada); que el dueño provoque "quedan pocas" en vivo desde el panel | ingeniero-go |
| I-2 | 8 mínimos por mayor inventados (6 u.; P009, P012, P016, P017, P020, P028, P037, P038) se muestran como reales | "Mínimo por confirmar" o etiqueta gris; idealmente un flag por campo en la API | constructor-web (+ ingeniero-go) |
| I-3 | Tras editar un precio en el panel, la ficha sigue diciendo "Precio referencial (video del …)" | API: marcar cuándo la tienda cambió el precio; web: "Precio actualizado por la tienda el {fecha}" | ingeniero-go + constructor-web |
| I-4 | El panel no permite agregar productos ni fotos (la API ya tiene POST) | "Agregar producto" con nombre, categoría, precio, mayor, stock, link de TikTok y foto; si no entra, declararlo Fase 2 | constructor-web (+ ingeniero-go para subida de imagen) |
| I-5 | El hero no dice qué venden; el logo de 150 px empuja todo bajo el pliegue | Bajada "Juguetes virales, Navidad, hogar y escolar, al detalle y por mayor." y logo de ~110 px | constructor-web |
| I-6 | La confirmación se pierde al recargar o al volver de WhatsApp | Guardar número, `whatsapp_url` y total en `sessionStorage`, o redirigir a `/pedido/listo?n=DC-xxxx` sin datos personales | constructor-web |
| I-7 | El demo apunta al WhatsApp real del negocio | Número de prueba por defecto, el real solo para la reunión | ingeniero-go |
| I-8 | Inicio pesado en móvil (TBT 2,3 s): miniaturas de 165 px pedidas en `w=1280` por falta de `sizes` | `sizes` correctos, menos `"use client"`, carrusel y embed diferidos | constructor-web |

## Menores

| # | Hallazgo | Sugerencia | Dueño |
|---|---|---|---|
| M-1 | `X-Forwarded-For` falso a través de la web salta el bloqueo de login y el límite de pedidos (= A1 de seguridad) | Acotar `TRUSTED_PROXIES` | ingeniero-go |
| M-2 | og:url, og:image, sitemap y robots apuntan a `localhost:3000` | Build arg `NEXT_PUBLIC_SITE_URL` con `${WEB_PORT}` | ingeniero-go |
| M-3 | Chip "Navidad" duplicado (temporada + categoría) | Ocultar o renombrar | constructor-web |
| M-4 | Nombre con `*` o `()` → "Escribe tu nombre." aunque esté lleno | "Usa solo letras en tu nombre." | constructor-web |
| M-5 | Táctiles < 44 px: footer, chips, select, "Ver fuente", "Abrir en TikTok" | `min-h-11` | constructor-web |
| M-6 | "Familia Emporio" con contraste 4,44:1 | Sobre negro o #8A5E0F | constructor-web |
| M-7 | Textos internos visibles: "[CONFIRMAR CON EL DUEÑO: …]" y "[EJEMPLO] Stock…" | Etiqueta gris "dato por confirmar" | constructor-web |
| M-8 | El número se corta en "DC-" / "0007" | `nowrap` | constructor-web |
| M-9 | En escritorio el header dice "DONDE CARLITOS" sin "Emporio" | Nombre completo | constructor-web |
| M-10 | `actualizado_en` 5 h atrás de `creado_en` | Guardar en UTC | ingeniero-go |
| M-11 | Teléfono sin formato en el panel | `+56 9 8765 4321` | constructor-web |
| M-12 | Fotos donde se ve más a la persona que al producto | Recorte más cerca u otro fotograma | constructor-web |
| M-13 | "desde 2 u. · desde 2 u." y "· por mayor" redundantes | Ocultar lo redundante | constructor-web |
| M-14 | Hueco en la tarjeta del carrusel de ofertas | `mt-auto` | constructor-web |
| M-15 | La promo del Tren grande vence el 15/11 | Revisar ofertas antes de la reunión; campo `oferta_hasta` | ingeniero-go |
| M-16 | Avisos de preload del logo sin usar en consola | Ajustar `priority` | constructor-web |

## Marca
Negro y oro y voz Carlitos ✅ · una oferta por pantalla ✅ · precios con fuente ✅ (salvo I-1, I-2) · detalle y mayor ✅ · horario con fecha ✅ · contraste ✅ salvo M-6 · sin 🏆 ✅ · 360 px sin scroll horizontal · logo sin alterar ✅ · ofertas con placa flúor y tachado solo con precio real ✅.

## Veracidad
Sin lorem, inglés ni testimonios. 78 links internos en 200. Horario, pagos, envíos, historia y aviso antiestafa van marcados. No se usan "envíos a todo Chile", el lema sin confirmar ni el 🏆. Problemas: I-1, I-2, I-3.

## Prueba del dueño
La reconoce como su tienda (logo, voz, productos de esta semana). Usa el panel de pedidos desde el celular. En productos edita pero no puede agregar (I-4); el formulario de 12 campos es largo: precio y stock arriba, el resto en "Más opciones".

## Enfoque "catálogo oficial de lo que sale en los videos + pedidos por WhatsApp" (no afecta el veredicto)
**Ya lo cumple:** titular "¿Lo viste en TikTok? Acá lo pides.", sección con embed, "Míralo en TikTok" y "Ver fuente" por producto, pedido que termina en WhatsApp con DC-xxxx, "La web no cobra", canal oficial.

**Cambiaría:**
1. "Carrito" → "Mi pedido"; "Hacer pedido" → "Armar mi pedido para WhatsApp".
2. "Salió en TikTok el dd/mm" en cada tarjeta en vez de "¡LLEGÓ!" genérico.
3. Botón "Pedir este por WhatsApp" por producto.
4. Ordenar por video reciente por defecto.
5. Alta de producto desde el link de TikTok en el panel (I-4).
6. "Pregúntanos si queda" en vez de un stock que nadie va a mantener.

## Notas de método
- Login probado en una instancia aislada (proyecto `dcqa`, API 127.0.0.1:18181, web 127.0.0.1:3200, clave de prueba fuera del repo, `WHATSAPP_NUMERO` de prueba); borrada al terminar.
- No se leyeron `api/.env` ni `api/ADMIN_PASSWORD.local.txt`. No se abrió `wa.me` ni se envió nada.
- En el stack principal queda un pedido de prueba, DC-0007.
