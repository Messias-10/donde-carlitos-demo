# QA · Ronda 2 — Donde Carlitos Emporio (v2: WhatsApp + seguimiento + panel)

*qa-revisor · Fase 4 · 2026-11-15 · guardado por el orquestador (el subagente no puede escribir reportes). Stack principal en http://localhost:3100 (solo lectura); pruebas con escritura en la instancia aislada `dcqa2` (127.0.0.1:3201, clave de prueba fuera del repo), ya borrada. Capturas: `qa/capturas/ronda-2/` (`m` = móvil 375×812, `d` = escritorio 1366×900).*

## Veredicto: **no listo para mostrar al dueño (por ahora)**

La base está bien: los mensajes de la landing (a, b, c) y del panel (d1–d8) calzan letra por letra con `propuesta/mensajes-whatsapp.md` (251 comparaciones, 0 diferencias); el panel recorre todos los estados; registrar un pedido tomó **13,6 s** (meta < 30 s); el seguimiento no muestra datos personales; vendedor y admin bien separados; build y tests en verde.

Hay **3 bloqueantes** que se notarían frente al dueño:
- **B-1.** El guion del pitch promete pantallas que el demo no tiene.
- **B-2.** 3 de los 8 precios de la vitrina son de promo por comentario y se muestran como precio de detalle.
- **B-3.** Los pedidos registrados en vivo no son demo: el botón abre WhatsApp al número pegado y sin `*[PRUEBA · DEMO]*`.

Con B-1 a B-3, I-1 e I-5 corregidos, queda listo; el resto puede ir como fase 2.

## Tabla de pruebas

| # | Prueba | Resultado |
|---|---|---|
| 1 | Arranque limpio (instancia aislada) | ✅ 74 s; `healthy`; `/healthz` 200. Stack principal intacto |
| 2 | Build y tests | ✅ `npm ci`, eslint, tsc, `next build` (sin `/carrito` ni `/catalogo`), `npm audit` 0; gofmt, `go vet`, `go test ./...`, govulncheck sin vulnerabilidades |
| 3 | Landing 375×812 y escritorio | ✅ con hallazgos B-2, I-6, I-7 |
| 3b | Links (a)(b)(c) | ✅ (a)=a2, (b)=b2; 16 variantes de (c) iguales a la referencia §7. Con `WHATSAPP_REGIONES` vacío, Regiones abre el 56921145368 con el texto (b) |
| 4 | Guion del pitch en el panel | ⚠ funciona entero pero no calza con el guion (B-1) |
| 4b | Registrar DC-0110 | ✅ 13,6 s en 7 acciones |
| 4c | Estados | ✅ regiones: cotizado → esperando pago → pagado (confirmación obligatoria) → preparando → despachado (agencia, n.º, comuna) → entregado; Santiago con pago en tienda: cotizado → preparando → listo → entregado (exige pago); cancelados con y sin motivo |
| 4d | d1–d8 decodificados | ✅ 0 diferencias; pedidos demo abren `wa.me/?text=` sin número |
| 4e | Seguimiento | ✅ sin nombre, teléfono ni comuna; `noindex`, `no-referrer`, `no-store`; token falso → 404 amable. Problema: I-5 |
| 4f | Vendedor y admin | ✅ vendedor 403 en CSV, resumen, usuarios y productos (escritura); puede registrar pedidos |
| 4g | Sin sesión | ✅ 12 rutas → 401; login sin Origin u Origin ajeno → 403 |
| 4h | Login/logout | ✅ 429 tras 5 fallos (XFF falso no lo salta); logout 204 |
| 4i | Editar precio → vitrina | ✅ "Precio publicado el 15/11: $3.200" |
| 5 | Manipulación | ✅ total falso, cantidades inválidas, precio negativo, descuento > total, producto inexistente, `<script>`, teléfono extranjero, campo extra → 400; `POST /api/pedidos` → 404 |
| 6 | CSV | ⚠ BOM, `;` y `no-store` ✅; columnas no son las de `copy-landing.md` P5 (I-4) |
| 7 | Lighthouse móvil | ver tabla |
| 8 | Accesibilidad | ✅ táctiles ≥ 44 px a 360 px, sin scroll horizontal, `alt`, `lang="es-CL"`; menor M-13 |
| 9 | Veracidad | ⚠ B-2 |
| 10 | 5 segundos | ⚠ qué venden y cómo pedir sí; dónde están no (I-7) |
| 11 | Prueba del dueño | ⚠ ver abajo |

**Lighthouse móvil 13.6.0** (BP baja por HTTP local; SEO 66 por `noindex` intencional)

| Página | Rend. | Acc. | LCP | TBT | CLS |
|---|---|---|---|---|---|
| `/` (stack principal) | 86 | 100 | 1,0 s | 550 ms | 0 |
| `/` (aislada, recién levantada) | 72 | 100 | 1,0 s | 1.900 ms | 0 |
| `/seguimiento/{token}` | 76 | 100 | 1,2 s | 1.160 ms | 0 |
| `/admin` (login) | 88 | 100 | 3,0 s | 280 ms | 0 |
| `/admin` (Hoy) | **63** | 100 | 3,9 s | **1.380 ms** | 0,034 |
| `/admin/pedidos/nuevo` | 86 | 100 | 3,3 s | 280 ms | 0,052 |
| `/admin/pedidos/{id}` | 74 | 96 | 3,7 s | 630 ms | 0,029 |

## Bloqueantes

| # | Hallazgo | Evidencia | Sugerencia | Dueño |
|---|---|---|---|---|
| B-1 | El guion de `pitch.md` muestra lo que el demo no hace: (a) "Hoy → DC-0102 → Marcar despachado", pero pagado solo pasa a preparando; (b) resumen "por comuna" y "quién despachó cuántos" no existe; (c) "los datos de transferencia los escribes en Ajustes", pero `/admin/ajustes` da 404 | `m34-despachar-dc0102.png`, `m62-resumen.png`, `m63-resumen-completo.png` | Botón de despacho desde Hoy; resumen por comuna y por usuario; Ajustes de solo lectura; o ajustar el guion | ingeniero-go + constructor-web; guion: orquestador/estratega |
| B-2 | Precios de promo por comentario presentados como precio de detalle (Esfera $500, Vasos $1.000, Telón $5.000) | Tarjetas P002, P008, P007; la nota solo está en `descripcion` | Campo `precio_es_promo` y placa "PROMO DEL VIDEO" con "comenta 'yo quiero mi promoción'", sin "detalle"; o sacarlos de la vitrina | ingeniero-go + constructor-web |
| B-3 | Pedidos registrados en vivo con `es_demo=false`: "Abrir en WhatsApp" va al número pegado y sin la línea demo | `m28-guardado-dc0110.png` | En modo demo, todo pedido nuevo `es_demo=true`; en el guion, teléfonos `+56 9 0000 01xx` | ingeniero-go |

## Importantes

| # | Hallazgo | Sugerencia | Dueño |
|---|---|---|---|
| I-1 | El stock inventado bloquea "Pagado" (`409 stock_insuficiente`); ensayar el guion dos veces lo rompe | Stock no bloqueante (aviso) por defecto; stock alto en los productos del guion | ingeniero-go |
| I-2 | El dueño no elige la vitrina; el formulario conserva "Stock" y "quedan pocas" de la v1 | `en_vitrina` + `orden_vitrina` (máx. 8), `GET /api/vitrina`; quitar "quedan pocas" y mover stock a "Más opciones" | ingeniero-go + constructor-web |
| I-3 | "Vendido" suma cotizaciones | Sumar solo pagado o posterior; "Pedidos" aparte | ingeniero-go |
| I-4 | CSV sin `productos`, `fecha_pagado`, `fecha_despachado`, `despachado_por`, `motivo_cancelacion`; fechas UTC ISO; códigos en vez de nombres; teléfono `'+569…` | Columnas de P5, hora de Chile `dd-mm-aaaa hh:mm`, nombres, teléfono como dígitos `569…` | ingeniero-go |
| I-5 | Seguimiento cae con 500 tras 16 tokens falsos (límite por IP compartido); `error.tsx` habla del carrito | Web: 429 → "muchas consultas, intenta en un minuto", borrar frase del carrito. API: no limitar por IP sin proxy confiable o contar solo los 404 | constructor-web + ingeniero-go |
| I-6 | Caché de `/_next/image` (24 h) puede mostrar los recortes viejos con caras | Renombrar recortes con hash; checklist del pitch: ventana privada | constructor-web |
| I-7 | La dirección queda bajo el pliegue; en 360×640 "Pedidos a regiones" se corta (logo de 150 px repetido) | Logo del hero ~96 px o línea "📍 Alameda 2432 · Metro U. Latinoamericana" sobre los botones | constructor-web |
| I-8 | "Hoy" lento en móvil (63, TBT 1,38 s) | Server Components en tarjetas; revisar hidratación y pollings | constructor-web |

## Menores

| # | Hallazgo | Dueño |
|---|---|---|
| M-1 | "Cliente conocido" sigue tras cambiar el teléfono | constructor-web |
| M-2 | Buscar "squishy" deja el squishy viral último | ingeniero-go |
| M-3 | "ver en la tienda" apunta a `/producto/{slug}` (404) | constructor-web |
| M-4 | "hace 1 días", "hace 0 días" | constructor-web |
| M-5 | Precio de línea del catálogo no editable | constructor-web + ingeniero-go |
| M-6 | Seguimiento sin productos ni total (decisión de privacidad; alinear el spec) | orquestador |
| M-7 | La API permite cancelar un despachado; `arquitectura-landing.md` C3 lo prohíbe | ingeniero-go |
| M-8 | "+ Nuevo pedido" tapa la última tarjeta; menú a 360 px corta "Resumen" sin señal de scroll; saludo truncado | constructor-web |
| M-9 | "Hoy": acción "Recordar pago o cancelar" en vez de dos botones; falta "Cotizados hoy" | constructor-web |
| M-10 | Agencia "Otra" sin campo "¿Cuál?" | constructor-web |
| M-11 | Paréntesis anidados en (c): "embalaje de 16 (surtido)" | ingeniero-go |
| M-12 | JSON-LD publica el horario por confirmar | constructor-web |
| M-13 | `<dl>` con `<p>` hijo en la ficha | constructor-web |
| M-14 | Bloqueo de login compartido: documentar `docker compose restart api` en la checklist | ingeniero-go |
| M-15 | Usuarios sin "canal por defecto" | ingeniero-go + constructor-web |
| M-16 | Desde cotizado se ofrecen dos botones principales | constructor-web |

## Marca
✅ negro y oro, voz, flúor solo en precios, sin carrito ni stock en la landing, botones de canal distinguibles, barra fija, hoja "¿DÓNDE LO QUIERES?", chips de estado, logo sin alterar, sin 🏆, sin caras con caché limpia. ❌ B-2. ⚠ I-2.

## Veracidad
✅ los 8 productos de la vitrina salen de `research/productos-semilla.json` con su video; "por confirmar" en lo dudoso; sin "envíos a todo Chile", testimonios, cifras, lorem ni inglés; datos de transferencia ausentes de lo público; teléfono enmascarado en logs. ❌ B-2. ⚠ M-12.

## Prueba del dueño y de quien atiende regiones
- La reconoce como su tienda.
- Quien atiende regiones lo usaría para lo diario: registrar ~14 s, "Hoy" muestra las alertas del guion, mensajes con un toque, confirmación de pago obligatoria.
- Se traba en: stock que bloquea el pago (I-1), despacho con paso extra (B-1a), precio de línea no editable (M-5), "Hoy" lento (I-8).
- El dueño no arma su vitrina (I-2), el resumen infla lo vendido (I-3) y el CSV no dice qué se vendió ni a qué hora local (I-4).

## Notas de método
- Login solo en `dcqa2` con clave de prueba fuera del repo, nunca en pantalla; todo borrado.
- No se leyeron `api/.env` ni `api/ADMIN_PASSWORD.local.txt`; no se abrió `wa.me` ni se envió nada; no se tocaron el stack principal ni `mouchett-*`.
- Comparaciones de mensajes con la referencia de `mensajes-whatsapp.md` §7 en `node:26.13.1-alpine`.
