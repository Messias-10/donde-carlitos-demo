---
name: estilo-carlitos
description: Guía de estilo de Donde Carlitos Emporio (logo, paleta, tipografía, tono, emojis, ofertas, componentes, qué evitar). Cárgala antes de diseñar, escribir copy o revisar cualquier pantalla del sitio.
---

# Estilo Carlitos

**La idea en una línea:** es **el local de remate de Meiggs que tú conoces, pero ordenado**. Negro y oro del logo para la marca, amarillo flúor de etiqueta de precio para las ofertas, y la voz de "¡Hola chiquillos, llegaron…!" en cada titular.

**Fuentes:** todo sale de `research/adn-marca.md`, `research/ofertas.md` y `research/top-contenido.md`, sacados de las redes oficiales TikTok @mayte_emporio e IG @mayte_emporioo (2026-11-14). Las decisiones de diseño que no se observaron directamente van marcadas *(inferencia)*.

**Enfoque del sitio (actualizado 2026-11-14):** se resume en "Tu WhatsApp sigue igual, pero cada pedido queda ordenado". Detalle en `propuesta/arquitectura-landing.md`, `copy-landing.md` y `mensajes-whatsapp.md`.

**El sitio tiene cuatro piezas:**
1. una **landing** con dos botones de WhatsApp por canal (**Santiago / retiro** y **Regiones**), con el mensaje ya escrito;
2. una **vitrina liviana** "Lo nuevo de los videos", con el botón **"Pedir este por WhatsApp"**;
3. una **página pública de seguimiento** (`/seguimiento/{token}`);
4. un **panel privado** de gestión de pedidos para quienes atienden los WhatsApp.

**Ya no existen:** carrito, botón "Agregar", contador de carrito, checkout, formulario de pedido web ni stock en pantalla. Si un diseño los trae, está fuera de estilo.

**Nombre:** la marca visible es **Donde Carlitos Emporio**, así está en el logo. Las cuentas se llaman "Mayte emporio" y las descripciones usan "Familia Emporio".
- Hasta que el dueño confirme, en la web: "Donde Carlitos Emporio" como marca y "Familia Emporio" como firma cálida.
- Ver `research/pendientes.md` §4.

---

## 1. Logo

**Archivo:** `research/raw/logo-tiktok-avatar.png`. Para la web va una copia en `web/public/` con su crédito en `web/public/CREDITOS.md`.
- Es un círculo negro con un aro de oro con glitter, el Empire State, una corona, "Donde" y "Emporio" en script blanco y **"CARLITOS"** en mayúsculas doradas pesadas.
- **No tiene 🏆.**

**Reglas:**
- **No se rediseña, no se redibuja, no se recolorea, no se recorta y no se le saca el texto chico.** Se usa tal cual (regla 4 de CLAUDE.md).
- **Siempre como círculo completo.** Va sobre negro (`--c-negro`), o sobre cualquier fondo con un borde o sombra suave que separe el círculo. Su fondo negro es parte del logo.
- **Dónde va:**
  1. header, a la izquierda, 40–48 px en móvil y 56–64 px en desktop;
  2. hero, 120–160 px, centrado o a la izquierda del titular;
  3. footer, 64 px;
  4. favicon y ícono de app: el mismo círculo, con el entendido de que el texto chico no se lee.
- **Tamaño mínimo:** 40 px de diámetro en pantalla. Bajo eso, usar el **texto** "Donde Carlitos Emporio" en tipografía de marca, no el logo.
- **Espacio libre** alrededor: al menos el 15 % del diámetro.
- **Sobre fotos:** solo si la zona detrás es oscura, o poniendo el círculo encima de una placa negra.
- **Nunca:** sobre amarillo flúor sin borde, en monocromo, estirado, con filtros, ni junto a otro logo que compita.
- Si se necesita una versión limpia, va en `pendientes.md`: se le pide al dueño. Cualquier SVG propio se marca "PROPUESTA – no es el logo oficial" y no se usa en el demo.

## 1b. Fotos e imágenes: sin personas reconocibles

**Regla dura:** **no se usan fotos, miniaturas ni fotogramas donde se reconozca a una persona.** Esto incluye a vendedores, al dueño, al equipo, a clientes, a menores y a la gente de las filas. El material del repo y del sitio lo pueden ver terceros.

- **Sí:**
  - fotos **solo del producto** (sobre fondo neutro, en la góndola o en la caja);
  - fotos de **contenedores, cajas o fachada** sin caras;
  - **recortes** de un fotograma que dejen solo el producto o la mano;
  - ilustraciones o íconos.
- **No:**
  - miniaturas de TikTok tal cual: casi todas muestran la cara del vendedor;
  - fotos de inauguraciones con público;
  - capturas con comentarios y nombres de usuario (`research/raw/tt-*.png` es material interno, **nunca** va a `web/`);
  - avatares de cuentas de terceros.
- **Para enlazar un video,** no se muestra su miniatura: se pone la foto del producto y el link **"▶ Míralo en TikTok"**.
- **Si una imagen trae una cara de fondo:** se recorta, o se tapa con un desenfoque fuerte o un bloque sólido, antes de copiarla a `web/public/`. Ante la duda, no se usa.
- **Créditos:** cada imagen va en `web/public/CREDITOS.md` con su origen.
- **Excepción:** solo con autorización escrita del dueño y de la persona. Hasta entonces, no hay excepciones.

## 2. Paleta (design tokens)

```css
:root {
  /* Marca (del logo) */
  --c-negro:        #0B0B0B;  /* fondo del logo */
  --c-carbon:       #1A1A1A;  /* superficies oscuras secundarias */
  --c-oro:          #E2B33C;  /* oro medio: CARLITOS, aro */
  --c-oro-claro:    #FCE48F;  /* brillo del oro */
  --c-oro-profundo: #9A6A12;  /* sombra del oro; oro para texto sobre claro */
  --c-blanco:       #FFFFFF;  /* script "Donde/Emporio" */
  --c-crema:        #FFF8E1;  /* fondo claro cálido (inferencia: sirve para leer catálogo largo) */

  /* Oferta (del local y la historia del "toldo amarillo") */
  --c-fluor:        #F5E800;  /* etiqueta de precio flúor (#F6FC14 muestreado) */
  --c-magenta:      #D6125A;  /* carteles rosados/magenta del local (#DF125C muestreado) */
  --c-rojo:         #C8102E;  /* urgencia: "solo hasta el domingo", agotado */

  /* Funcionales */
  --c-whatsapp:     #25D366;
  --c-whatsapp-osc: #0B6E4F;  /* también botón "Regiones" con texto blanco */
  --c-gris:         #6B6B6B;  /* texto secundario sobre claro */
  --c-gris-claro:   #B8B8B8;  /* texto secundario sobre oscuro */

  /* Degradado oro (solo en titulares grandes y bordes, nunca en texto chico) */
  --g-oro: linear-gradient(180deg, #FCE48F 0%, #E2B33C 45%, #9A6A12 100%);
}
```

**Contraste (WCAG, calculado):**

| Texto / fondo | Ratio | Uso |
|---|---|---|
| blanco / negro | 19,7 | ✅ texto general en modo oscuro |
| oro / negro | 10,1 | ✅ titulares, precios y links sobre negro |
| oro / carbón | 8,9 | ✅ |
| negro / flúor | 15,4 | ✅ **precio y badge de oferta** |
| negro / oro | 10,1 | ✅ botones dorados con texto negro |
| blanco / magenta | 5,1 | ✅ badge "PROMO", texto de 14 px o más |
| blanco / rojo | 5,9 | ✅ urgencia y "agotado" |
| negro / WhatsApp | 9,9 | ✅ **CTA de WhatsApp: texto negro, no blanco** |
| blanco / WhatsApp oscuro | 6,3 | ✅ alternativa |
| oro profundo / blanco | 4,7 | ✅ solo para texto de 16 px o más |
| gris / blanco | 5,3 | ✅ secundario |
| **oro / blanco** | **1,95** | ❌ **prohibido para texto** |
| **blanco / oro** | **1,95** | ❌ **prohibido** |
| **blanco / WhatsApp** | ~1,9 | ❌ **prohibido** |
| magenta / negro | 3,8 | ⚠ solo para textos de 24 px o más, o elementos no textuales |

**Reglas:**
- **Superficies:** el header, el hero y el footer van en **negro y oro**. El catálogo puede ir en fondo **crema o blanco** para leer mejor, con tarjetas de borde suave *(inferencia)*.
- **El flúor es solo para precio y oferta.** Nunca se usa en fondos grandes: si todo grita, nada grita.
- **Proporción:** 60 % neutro (negro o crema), 25 % oro y blanco, 10 % flúor, 5 % magenta o rojo.

## 3. Tipografía (Google Fonts)

Equivalentes cercanos al logo y a su texto en pantalla:

| Rol | Fuente | Por qué |
|---|---|---|
| Display / precios / "CARLITOS" | **Anton** (400) | Sans pesada y condensada, como "CARLITOS"; los números grandes de precio se leen de lejos |
| Acento script ("Donde", "Emporio", "¡Llegó!") | **Kaushan Script** (400) | Manuscrito grueso y casual, como el script blanco del logo. Máximo 2–3 palabras |
| Texto, UI y fichas | **Inter** (400/600/800) | Legible en el navegador interno de TikTok/IG; buena con números (`font-variant-numeric: tabular-nums` para precios en tabla) |

**Reglas:**
- **Titulares en MAYÚSCULAS con Anton.** Inter va en tipo oración.
- **Precio:** Anton, 28–40 px en la tarjeta, con el "$" más chico. Formato chileno: `$4.000`, con punto de miles y sin decimales.
- **Kaushan Script nunca va en párrafos, botones ni precios.** Solo como acento sobre negro, en blanco u oro.
- **No usar** Bebas Neue y Anton juntas, ni más de 3 familias.
- **Tamaños base (móvil):** cuerpo 16 px, ficha 14–15 px, titular hero 40–48 px.

## 4. Tono de voz

**Quién habla:** el vendedor del video, cercano y de barrio, pero escrito, así que más ordenado. Español de Chile, sin modismos inventados y sin groserías (en los videos hay alguna; en la web no).

**Reglas:**
1. **Habla a "chiquillos / chicos"** o de tú. Nunca de "usted" formal en titulares. En la ficha puede ir neutro ("Precio por embalaje").
2. **Novedad primero:** "¡Llegó…!", "¡Nos llegaron…!", "Volvió…".
3. **Precio claro y redondo, en luquitas en el titular** y en `$` en la ficha.
4. **Siempre la escalera:** detalle → por embalaje / por mayor desde N unidades.
5. **Urgencia solo si es verdad:** fecha de término dicha en el video, o preventa con fecha de llegada. Sin stock en pantalla.
6. **Cierre con lugar:** "Te esperamos en Alameda 2432".
7. **Frases cortas.** Una idea por línea. Nada corporativo ("soluciones", "experiencia de compra", "estimado cliente").
8. **"Pa negocio" es bienvenido:** el revendedor es cliente de primera.
9. **Comparaciones de precio con otros lugares solo si el dueño las dijo** en un video ("en todos lados vale…"); nunca inventadas.
10. **Sin exagerar sobre la competencia, el barrio ni sus comunidades.**

**10 frases de ejemplo en su voz** (para web; las cifras van marcadas `[EJEMPLO]` salvo que vengan de un producto real del panel):
1. "¡Hola chiquillos! Llegó lo nuevo de la semana 👇"
2. "Squishy viral a 4 luquitas… y a $3.500 por caja de 12." *(dicho así en TT-7707649267464293639)*
3. "¿Pa negocio? Precio por embalaje desde 10 unidades."
4. "Ojo: esta oferta vale solo hasta el domingo." `[EJEMPLO de fecha]`
5. "Quedan pocas, chiquillos. Cuando se acaban, se acaban."
6. "¿Lo viste en TikTok? Acá está con su precio. Tócalo y nos escribes con el producto ya puesto."
7. "¿Eres de regiones? Escríbenos al WhatsApp de regiones y te cotizamos al tiro."
8. "Te esperamos en Alameda 2432, a pasos del Metro Unión Latinoamericana."
9. "Cuando hagas tu pedido te mandamos un link para ver en qué va. Sin dar tus datos."
10. "Somos Familia Emporio: traemos el contenedor y te lo dejamos a precio de remate."

## 5. Emojis

**Observado:** casi no hay emojis en pantalla. En las descripciones aparecen 🛍️ ✨ 📍 🕘 👉 🔥 💬 👇 💣 💯 ❤️ (TT-7592989161653685522). El 🏆 **no** está en el logo (solo en la cuenta secundaria no confirmada).

**Reglas:**
- **Máximo 1 emoji por titular** y 0 en botones de acción crítica (pagar, enviar pedido), salvo el ícono de WhatsApp.
- **Permitidos, con su uso:**
  - 📍 dirección;
  - 🕘 horario;
  - 🔥 oferta o tendencia;
  - 👇 "mira abajo";
  - 💬 comentar o pedir;
  - 🛍️ catálogo;
  - ✨ novedad;
  - 📦 envíos (solo si se confirma que despachan);
  - 💯 calidad u "original".
- **Evitar:** 🏆 hasta confirmar con el dueño, 💣 (en web puede leerse mal), 😱 en serie, y cualquier emoji como reemplazo de texto en precios.
- **Nunca** emojis dentro de un número de precio ni en textos legales o de confianza.

## 6. Cómo se comunica una oferta

**Estructura (la "fórmula Carlitos", en `research/ofertas.md`):**
```
¡LLEGÓ / REMATE!            (Anton, mayúsculas)
Nombre del producto
$4.000 detalle               (Anton grande, sobre flúor)
$3.500 por caja de 12        (Inter 600)
🔥 Solo hasta el domingo      (rojo; solo si es real)
Precio publicado el dd/mm    (Inter 13px gris)
[ Pedir este por WhatsApp ]  → hoja "¿Dónde lo quieres?" Santiago / Regiones
▶ Míralo en TikTok
📍 Alameda 2432
```

- **"Yo quiero mi promoción" en la web:** es la frase de los videos. Se puede citar en un titular o en la FAQ ("¿Comentaste 'yo quiero mi promoción'? Dilo en el chat").
- **No es un botón ni un mecanismo web.** La web no aplica precios promo ni guarda promociones; todo se resuelve por WhatsApp.
- Cómo se canjea en caja está pendiente de confirmar con el dueño.

**Palabras que sí:** llegó, llegaron, volvió, novedad, viral, remate, oferta, promo, "pa negocio", por embalaje, por mayor, desde N unidades, preventa, luquitas, "pídelo por WhatsApp".

**Palabras que no:** "descuento exclusivo online", "imperdible" en todo, "OFF" en inglés, "liquidación total" si no lo es, porcentajes inventados.

**Urgencia válida:**
- fecha de término dicha en el video (`oferta_hasta`; al vencer se quita el badge y el precio tachado);
- preventa con fecha de llegada;
- límite por persona, si el dueño lo define.

**No se usa** "quedan pocas", "hay stock" ni "agotado" en la web pública: nadie actualiza el stock en tiempo real. El mensaje de WhatsApp pregunta "¿Les queda?".

**Precio tachado:** solo si existe un precio anterior real (ej.: el tren de $20.000 que pasó a $17.000 en TT-7706656443843284232).

## 7. Sí y no visuales

**5 SÍ:**
1. ✅ **Producto en contexto real** (góndola, caja, contenedor), **sin personas reconocibles** (§1b). Foto 4:5.
2. ✅ **Negro + oro en la marca y flúor solo en precios.** El contraste es la firma.
3. ✅ **Precio grande y doble** (detalle + por mayor) en cada tarjeta, con "Precio publicado el dd/mm", legible sin sonido.
4. ✅ **Dos botones de WhatsApp por canal,** grandes y distinguibles por color (Santiago verde claro / Regiones verde oscuro), + bloque fijo "Cómo llegar" con 📍 Alameda 2432, Metro y horario con fecha.
5. ✅ **Prueba real y liviana:** link "▶ Míralo en TikTok", foto de contenedores sin caras, cuentas oficiales. Todo con fuente.

**5 NO:**
1. ❌ **Estética de retail corporativo** (blanco minimalista con tipografía fina, fotos de stock, modelos posando). No se parece a ellos.
2. ❌ **Caos de remate literal:** 5 colores flúor, todo parpadeando, banners rotando. "Ordenado" es la promesa.
3. ❌ **Oro sobre blanco para texto,** degradados dorados en textos chicos, glitter animado.
4. ❌ **Precios, reseñas, cifras de ventas, testimonios o "envíos a todo Chile"** sin fuente o sin confirmar.
5. ❌ **Fotos o miniaturas con personas reconocibles** (vendedores, clientes, menores), el 🏆, o nombres de personas (Carlitos, Mayte, Andrés) sin autorización. También **carrito, "Agregar", stock en pantalla** o datos bancarios en la web.

## 8. Componentes

### 8.1 Tarjeta de vitrina "Lo nuevo de los videos" (móvil primero)

La vitrina tiene de 6 a 8 tarjetas que vienen de `GET /api/vitrina`. Es liviana: **sin carrito, sin "Agregar", sin cantidades y sin stock.**

```
┌──────────────────────────────┐
│ [¡LLEGÓ! / OFERTA / PREVENTA]│  ← máx. 1 badge, arriba a la izquierda
│                              │
│   FOTO 4:5 del PRODUCTO      │  ← sin personas (§1b); object-fit: cover; lazy
│                              │
├──────────────────────────────┤
│ Squishy viral cambia color   │  Inter 600 15px, 2 líneas máx.
│ ▓ $4.000 ▓ detalle           │  Anton 32px negro sobre --c-fluor (placa)
│ ~~$20.000~~  (solo si hay oferta vigente con precio anterior real)
│ $3.500 c/u · caja de 12      │  Inter 600 14px (por mayor, si existe)
│ Precio publicado el 13/11    │  Inter 13px --c-gris
│ ▶ Míralo en TikTok           │  link de texto al video de origen
│ [  Pedir este por WhatsApp  ]│  botón ancho completo, 48 px
└──────────────────────────────┘
```

- **Fondo:** tarjeta blanca sobre sección crema, o carbón sobre negro. Radio de 12 px y sombra suave. En móvil va en grilla de 2 columnas a 360 px; o 1 columna con foto más grande *(inferencia; elegir según el peso)*.
- **Botón "Pedir este por WhatsApp":** fondo `--c-negro` y texto `--c-oro`, con el ícono de WhatsApp en oro.
  - No va en verde, para no confundirlo con los botones de canal. El canal se elige después.
  - Alto mínimo de 48 px.
- **Al tocar el botón** se abre la **hoja inferior "¿DÓNDE LO QUIERES?"**:
  - fondo `--c-negro`;
  - título en Anton blanco;
  - dos botones de canal (§8.4), apilados;
  - pie: "Te abrimos WhatsApp con el producto ya escrito. Tú solo envías.";
  - link "Cancelar".
  - La elección se recuerda durante la visita y se puede cambiar.
- **Pie de sección:** "Los precios son los que dijimos en el video y pueden cambiar. Pregúntanos por WhatsApp si queda."
- **Datos:** los precios salen del producto real con fecha (`fecha_precio`) o se marcan `[EJEMPLO]`. Al vencer `oferta_hasta` desaparecen el precio tachado y el badge de oferta.
- **Sin stock:** nada de "hay stock", "quedan pocas", "agotado" ni "Avísame".

### 8.2 Badges de la vitrina
| Variante | Fondo | Texto | Ejemplo |
|---|---|---|---|
| Precio (placa) | `--c-fluor` | negro, Anton | `$4.000` |
| Novedad | `--c-oro` | negro, Anton | `¡LLEGÓ!` |
| Oferta | `--c-magenta` | blanco, Inter 800, mayúsculas | `OFERTA` |
| Remate | `--c-negro` + borde 2 px `--c-oro` | `--c-oro`, Anton | `REMATE` |
| Preventa | `--c-negro` | `--c-fluor`, Inter 800 | `PREVENTA · llega aprox. 16/11` |
| Urgencia | `--c-rojo` | blanco | `SOLO HASTA EL DOMINGO` |
| Pa negocio | blanco + borde negro | negro | `PRECIO X MAYOR` |

- **Forma:** píldora (radio 999 px), padding 4×10 px, 12–13 px. **Máximo 1 badge** por tarjeta, más la placa de precio.
- **Opcional:** la placa de precio con forma de "estrella de etiqueta", como las del local (SVG de 12 puntas en flúor) *(inferencia; observado en `research/raw/tt-7592989161653685522.png`)*. Máximo una por tarjeta.
- **No hay** badge "YO QUIERO MI PROMO": no hay mecanismo web de promo (§6).

### 8.3 Banner de remate (opcional, uno a la vez)
```
┌───────────────────────────────────────┐ fondo --c-negro
│ Donde                                  │ Kaushan Script blanco 20px
│ ¡REMATE DE NAVIDAD!                    │ Anton 40px con --g-oro
│ Esferas, telones, trenes y más desde   │ Inter 16px blanco
│ ▓ $500 ▓                               │ Anton sobre flúor (precio real con fuente)
│ 🔥 Hasta el domingo · 📍 Alameda 2432    │ Inter 14px --c-gris-claro
│ [ Ver lo nuevo ↓ ]                     │ link a la vitrina (no compite con los botones de canal)
└───────────────────────────────────────┘
```
- **Opcional:** un aro dorado fino (1–2 px, `--c-oro`) que evoque el borde del logo. **Sin glitter animado ni carrusel automático.**
- **Imagen:** solo producto o contenedor **sin personas** (§1b).

### 8.4 Botones de canal de WhatsApp (el CTA principal)
Hay **dos canales**, cada uno con su número por configuración (`WHATSAPP_SANTIAGO`, `WHATSAPP_REGIONES`; si el de regiones está vacío se usa el de Santiago) y su mensaje pre-escrito (`propuesta/mensajes-whatsapp.md` (a), (b) y (c)).

| Botón | Fondo | Texto | Ícono | Uso |
|---|---|---|---|---|
| **Pedidos Santiago / retiro en tienda** | `--c-whatsapp` `#25D366` | **negro** `#0B0B0B` (9,9) | WhatsApp negro | hero, sección Santiago, barra fija, hoja de vitrina |
| **Pedidos a regiones** | `--c-whatsapp-osc` `#0B6E4F` | **blanco** (6,3) | WhatsApp blanco | hero, sección Regiones, barra fija, hoja de vitrina |

- **Hero:** los dos botones apilados, ancho completo, **56 px** de alto como mínimo, radio de 14 px, Inter 800 de 17 px.
  - Debajo, en chico: **"Estos son nuestros únicos números de pedidos"**, con link a la sección de canal oficial.
  - Es la evolución del "pedidos solo a este número" de la bio de IG.
- **Barra fija abajo (móvil):**
  - dos botones lado a lado, "Santiago / retiro" y "Regiones", de 52 px de alto y con los mismos colores;
  - aparece **después** de pasar el hero;
  - la página deja un padding inferior para que la barra no tape contenido.
- **Siempre** el ícono de WhatsApp + el texto del canal. **Nunca solo el color** para distinguirlos.
- **Enlace:** `https://wa.me/{numero}?text={mensaje codificado}`. El mensaje solo lleva el producto o el canal. **Nada de datos personales en la URL.**
- **Funcionan sin la API:** los números y mensajes no dependen del backend.
- **Prohibido:**
  - texto blanco sobre `#25D366` (1,9);
  - un tercer botón de WhatsApp "genérico";
  - "Enviar pedido (n)";
  - contadores;
  - "Agregar".

### 8.5 Página pública de seguimiento (`/seguimiento/{token}`)
```
┌──────────────────────────────┐ header negro con logo 44px
│ PEDIDO DC-0101               │ Anton 28px
│ ┌──────────────────────────┐ │
│ │ ● EN CAMINO              │ │ chip grande con el color del estado (§8.6)
│ │ Tu pedido salió por {ag} │ │ Inter 16px
│ └──────────────────────────┘ │
│ Línea de tiempo (vertical):  │
│ ● Cotizado        10/11 12:30│ pasos cumplidos: punto --c-oro-profundo
│ ● Pagado          10/11 15:02│
│ ◉ En camino       11/11 09:40│ paso actual: punto grande con el color del estado
│ Agencia · N.º seguimiento    │ número en Inter tabular + [Copiar número] → "¡Copiado!"
│ 2 x Squishy … · Total $8.000 │
│ [ Escríbenos por WhatsApp ]  │ botón del canal del pedido (colores §8.4)
│ Canal oficial (versión corta)│
└──────────────────────────────┘
```
- **Fondo:** crema con tarjeta blanca y header y footer negros. Es una página de **tranquilidad**: no lleva ofertas, banners ni vitrina.
- **No muestra** nombre, teléfono, comuna, notas ni quién atendió. Tampoco tiene un buscador por número.
- **Error o token vencido:** título "NO ENCONTRAMOS ESTE PEDIDO" en Anton, texto de ayuda y los dos botones de canal. Va en neutro, sin rojo alarmante.
- **Solo en el demo:** etiqueta gris "Pedido de prueba (demo)".

### 8.6 Estados de pedido: colores (panel y seguimiento)
Los usan el chip del panel (lista, "Hoy" y ficha) y el estado grande de la página de seguimiento. **Siempre** van con texto + punto o ícono; nunca el color solo. Los contrastes están calculados (WCAG AA, texto ≥ 4,5).

| Estado API | Panel | Seguimiento (título) | Fondo chip | Texto | Contraste | Idea |
|---|---|---|---|---|---|---|
| `cotizado` | Cotizado | Cotizado | `#E6E6E6` | `#0B0B0B` | 15,8 | neutro: aún nada |
| `esperando_pago` | Esperando pago | Esperando tu pago | `#FFE9A8` | `#6B4500` | 7,1 | ámbar: falta algo del cliente |
| `pagado` | Pagado | Pago recibido | `#D4F4E2` | `#0B6E4F` | 5,3 | verde claro: confirmado |
| `preparando` | Preparando | Preparando tu pedido | `#DCEAFF` | `#1D4E89` | 6,9 | azul: trabajo interno |
| `despachado` | Despachado | En camino | `#F3E1F7` | `#7A1E8C` | 7,1 | violeta: en tránsito |
| `listo_retiro` | Listo para retiro | Listo para retiro | `#E2B33C` (`--c-oro`) | `#0B0B0B` | 10,1 | oro: "ven al local" |
| `entregado` | Entregado | Entregado | `#0B6E4F` | `#FFFFFF` | 6,3 | verde sólido: cerrado |
| `cancelado` | Cancelado | Cancelado | `#FDE2E4` | `#A50E25` | 6,4 | rojo suave: cerrado sin venta |

```css
:root {
  --st-cotizado-bg:#E6E6E6;  --st-cotizado-fg:#0B0B0B;
  --st-espera-bg:#FFE9A8;    --st-espera-fg:#6B4500;
  --st-pagado-bg:#D4F4E2;    --st-pagado-fg:#0B6E4F;
  --st-prep-bg:#DCEAFF;      --st-prep-fg:#1D4E89;
  --st-desp-bg:#F3E1F7;      --st-desp-fg:#7A1E8C;
  --st-retiro-bg:#E2B33C;    --st-retiro-fg:#0B0B0B;
  --st-entregado-bg:#0B6E4F; --st-entregado-fg:#FFFFFF;
  --st-cancel-bg:#FDE2E4;    --st-cancel-fg:#A50E25;
}
```
- **No se usa** `--c-fluor` en los estados: el flúor es solo para precio y oferta.
- **Alerta de tiempo** (panel, pantalla "Hoy"):
  - "Esperando pago hace más de 24 h" y "Pagados sin despachar" llevan un **borde izquierdo de 4 px `--c-rojo`** y el texto "hace N h";
  - el chip del estado **no cambia de color**.
- **Línea de tiempo (seguimiento):**
  - pasos cumplidos con punto `--c-oro-profundo` y texto negro;
  - paso actual con el chip del estado;
  - Cancelado corta la línea y se muestra solo.

### 8.7 Panel de pedidos (`/admin`): estilo
- **Es una herramienta, no una vitrina:** fondo claro `#FFFFFF` / `#F7F7F7`, Inter 15–16 px y tablas o tarjetas densas pero legibles con el pulgar.
- **La marca se queda en el header:** negro con el logo de 32–40 px, el título "PEDIDOS DEL WHATSAPP" en Anton oro y el saludo "Hola, {usuario} · Canal: {Santiago | Regiones}".
- **Canal** del pedido en una mini etiqueta: Santiago (borde `#25D366` y texto negro) o Regiones (fondo `#0B6E4F` y texto blanco). Son los mismos colores de los botones de canal.
- **Acción siguiente:**
  - **un** botón grande por ficha (fondo `--c-negro` y texto blanco, 48 px de alto) con el verbo exacto de `copy-landing.md` P4 ("Marcar pagado", "Marcar despachado"…);
  - debajo, "Abrir en WhatsApp" para el mensaje sugerido (texto negro sobre `#25D366`).
- **"+ Nuevo pedido":**
  - botón fijo abajo a la derecha en móvil;
  - fondo `--c-oro` y texto negro;
  - campos ordenados para el pulgar; los botones de canal y de entrega van como segmentos grandes.
- **Confirmaciones sensibles** ("¿Revisaste que la plata está en la cuenta?") en un diálogo con el botón de confirmar en `--c-negro` y "Volver" secundario.
- **Sin flúor, sin glitter, sin emojis** en el panel (salvo dentro de la vista previa de los mensajes de WhatsApp).
- **Datos personales:** el teléfono va enmascarado en las listas (`+56 9 •••• 5368`) y completo solo en la ficha.

### 8.8 Otros patrones recurrentes
- **"▶ Míralo en TikTok":** link de texto bajo la tarjeta de vitrina. **Sin miniatura del video** (traen caras, §1b) y **sin embeds pesados** en la landing.
  - Si el video menciona una dirección antigua, no importa: la sección "Cómo llegar" deja clara la dirección vigente (📍 Alameda 2432).
- **Bloque "Canal oficial" (antiestafa):**
  - "Únicos números de pedidos" (Santiago y Regiones);
  - TikTok @mayte_emporio e Instagram @mayte_emporioo (verificada);
  - reglas para comprar seguro;
  - "No publicamos datos bancarios en la web: te los damos solo por el chat oficial".
  - No se nombran las cuentas no confirmadas.
  - Fondo `--c-carbon`, título en oro y texto blanco.
- **"Cómo comprar desde regiones":** 5 pasos numerados con números grandes en Anton `--c-oro` sobre negro, o `--c-oro-profundo` sobre crema. Los datos dudosos llevan la etiqueta gris chica "dato por confirmar".
- **Header de la landing:** negro, logo de 44 px y el link "Sigue tu pedido". **Sin buscador ni carrito.**

## 9. Checklist rápido antes de dar por buena una pantalla
- [ ] ¿El dueño reconocería el negro y oro, y el "¡Llegó…!"?
- [ ] ¿Está más ordenado que un video? (una oferta principal por pantalla)
- [ ] ¿En 3 segundos se entiende a qué botón apretar (Santiago / retiro o Regiones)?
- [ ] ¿Cero carrito, "Agregar", contadores o stock en pantalla?
- [ ] ¿Ninguna imagen muestra personas reconocibles? (§1b)
- [ ] ¿Cada precio tiene fuente y "Precio publicado el dd/mm", o dice `[EJEMPLO]`?
- [ ] ¿Se ven detalle y por mayor cuando existen?
- [ ] ¿La dirección y el horario están visibles y con fecha?
- [ ] ¿Pasa el contraste? (nada de oro o blanco sobre blanco u oro; nada de blanco sobre `#25D366`)
- [ ] ¿Los estados de pedido usan los tokens `--st-*` con texto, no solo color?
- [ ] ¿Máximo 1 emoji por titular, y sin 🏆?
- [ ] ¿Funciona a 360 px de ancho y en el navegador interno de TikTok, sin la API para los botones de WhatsApp?
