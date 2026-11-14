---
name: estilo-carlitos
description: Guía de estilo de Donde Carlitos Emporio (logo, paleta, tipografía, tono, emojis, ofertas, componentes, qué evitar). Cárgala antes de diseñar, escribir copy o revisar cualquier pantalla del sitio.
---

# Estilo Carlitos

**La idea en una línea:** es **el local de remate de Meiggs que tú conoces, pero ordenado**. Negro y oro del logo para la marca, amarillo flúor de etiqueta de precio para las ofertas, y la voz de "¡Hola chiquillos, llegaron…!" en cada titular.

**Fuentes:** todo sale de `research/adn-marca.md`, `research/ofertas.md` y `research/top-contenido.md`, sacados de las redes oficiales TikTok @mayte_emporio e IG @mayte_emporioo (2026-11-14). Las decisiones de diseño que no se observaron directamente van marcadas *(inferencia)*.

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
  --c-whatsapp-osc: #0B6E4F;
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
5. **Urgencia solo si es verdad:** fecha, stock real o "pocas unidades" que venga del panel.
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
6. "¿Lo viste en TikTok? Acá lo tienes con precio y stock."
7. "Arma tu pedido y nos llega directo al WhatsApp. Así de fácil."
8. "Te esperamos en Alameda 2432, a pasos del Metro Unión Latinoamericana."
9. "¡Yo quiero mi promo! Tócalo y te guardamos el precio especial."
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
[¡Yo quiero mi promo!]  [Pedir por WhatsApp]
📍 Alameda 2432
```

**Palabras que sí:** llegó, llegaron, volvió, novedad, viral, remate, oferta, promo, "pa negocio", por embalaje, por mayor, desde N unidades, quedan pocas, preventa, luquitas.

**Palabras que no:** "descuento exclusivo online", "imperdible" en todo, "OFF" en inglés, "liquidación total" si no lo es, porcentajes inventados.

**Urgencia válida:**
- fecha de término;
- stock del panel (`quedan ≤ 5`);
- preventa con fecha de llegada;
- límite por persona, si el dueño lo define.

**Precio tachado:** solo si existe un precio anterior real (ej.: el tren de $20.000 que pasó a $17.000 en TT-7706656443843284232).

## 7. Sí y no visuales

**5 SÍ:**
1. ✅ **Producto en mano o en contexto real** (pasillo, caja, contenedor): fotos verticales tipo video.
2. ✅ **Negro + oro en la marca y flúor solo en precios.** El contraste es la firma.
3. ✅ **Precio grande y doble** (detalle + embalaje) en cada tarjeta, legible sin sonido.
4. ✅ **Bloque fijo "Cómo llegar"** con 📍 Alameda 2432, Metro y horario con fecha de vigencia.
5. ✅ **Prueba real:** video de TikTok embebido, contenedores, filas en inauguración. Todo con fuente.

**5 NO:**
1. ❌ **Estética de retail corporativo** (blanco minimalista con tipografía fina, fotos de stock, modelos posando). No se parece a ellos.
2. ❌ **Caos de remate literal:** 5 colores flúor, todo parpadeando, banners rotando. "Ordenado" es la promesa.
3. ❌ **Oro sobre blanco para texto,** degradados dorados en textos chicos, glitter animado.
4. ❌ **Precios, reseñas, cifras de ventas, testimonios o "envíos a todo Chile"** sin fuente o sin confirmar.
5. ❌ **Usar el 🏆, el nombre de personas (Carlitos, Mayte, Andrés) o fotos de clientes** sin autorización del dueño.

## 8. Componentes

### 8.1 Tarjeta de producto (móvil primero)
```
┌──────────────────────────┐
│ [badge PROMO]   [▶ TikTok]│  ← badge arriba a la izquierda; ícono si viene de un video
│                           │
│     FOTO 4:5 (producto    │  ← fondo neutro o foto real; object-fit: cover
│      en mano/contexto)    │
├──────────────────────────┤
│ Squishy viral cambia color│  Inter 600 15px, 2 líneas máx.
│ ▓ $4.000 ▓  detalle        │  Anton 32px negro sobre --c-fluor (placa)
│ $3.500 x caja de 12        │  Inter 600 14px --c-negro / --c-gris
│ ● Hay stock · quedan pocas │  verde/rojo + texto (no solo color)
│ [ + Agregar ] [Yo quiero mi promo] │
└──────────────────────────┘
```
- **Fondo de la tarjeta:** blanco sobre catálogo crema, o carbón sobre fondo negro. Radio de 12 px y sombra suave.
- **Botón principal "Agregar":** fondo `--c-negro` con texto `--c-oro`, o fondo `--c-oro` con texto `--c-negro`. Alto mínimo de 44 px.
- **Estados de stock:** "Hay stock", "Quedan pocas", "Agotado" (con la tarjeta en gris y el botón "Avísame") y "Preventa · llega el DD/MM".
- Los precios salen del producto real (`research/productos-semilla.json`) o se marcan `[EJEMPLO]`.

### 8.2 Badge de oferta
| Variante | Fondo | Texto | Ejemplo |
|---|---|---|---|
| Precio | `--c-fluor` | negro, Anton | `$4.000` |
| Promo comentario | `--c-magenta` | blanco, Inter 800, mayúsculas | `YO QUIERO MI PROMO` |
| Remate | `--c-negro` + borde 2 px `--c-oro` | `--c-oro`, Anton | `REMATE` |
| Urgencia | `--c-rojo` | blanco | `SOLO HASTA EL DOMINGO` |
| Novedad | `--c-oro` | negro | `¡LLEGÓ!` |
| Pa negocio | blanco + borde negro | negro | `PRECIO X MAYOR` |

- **Forma:** píldora (radio 999 px), padding 4×10 px, 12–13 px.
- **Opcional:** el badge de precio puede tener forma de "estrella de etiqueta" como las del local (SVG estrella de 12 puntas en flúor) *(inferencia; observado en `research/raw/tt-7592989161653685522.png`)*. Máximo una estrella por tarjeta.

### 8.3 Banner de remate
```
┌───────────────────────────────────────┐ fondo --c-negro
│ Donde                                  │ Kaushan Script blanco 20px
│ ¡REMATE DE NAVIDAD!                    │ Anton 40px con --g-oro
│ Esferas, telones, trenes y más desde   │ Inter 16px blanco
│ ▓ $500 ▓                               │ Anton sobre flúor
│ 🔥 Hasta el domingo · 📍 Alameda 2432    │ Inter 14px --c-gris-claro
│ [ Ver ofertas ]  [ 💬 WhatsApp ]        │
└───────────────────────────────────────┘
```
- **Opcional:** un aro dorado fino (1–2 px, `--c-oro`) que evoque el borde del logo. **Sin glitter animado.**
- **Imagen a la derecha** (desktop) o abajo (móvil): foto real de productos o un fotograma del video.
- Va **un solo banner** activo a la vez. No hay carrusel automático.

### 8.4 CTA de WhatsApp
- **Botón fijo abajo en móvil** (sticky). Fondo `--c-whatsapp`, **texto negro** `#0B0B0B`, ícono de WhatsApp, alto de 52 px y radio de 14 px. Texto: **"Pedir por WhatsApp"**, o con carrito: **"Enviar pedido (3) por WhatsApp"**.
- **Número único** (+56 9 2114 5368) y abajo, en chico: "Este es nuestro único número de pedidos" *(dicho en la bio de IG: "pedidos solo a este número")*.
- **Enlace:** `https://wa.me/56921145368?text=` con un mensaje prearmado (producto, cantidad, detalle o mayor, nombre, comuna/región).
  - No se guarda nada sensible en la URL además del pedido.
  - El texto final lo define estratega-propuesta.
- **No usar** texto blanco sobre el verde de WhatsApp (contraste 1,9).

### 8.5 Otros patrones recurrentes
- **"Lo viste en TikTok":** miniatura vertical 9:16, ícono ▶ y debajo la tarjeta del producto. Si el video menciona una dirección antigua, superponer "Hoy en Alameda 2432".
- **Bloque de confianza:** cuentas oficiales (@mayte_emporio, @mayte_emporioo), el número único, "Nunca te pediremos datos por interno" (*confirmar redacción con el dueño*), medios de pago por local y boleta o factura.
- **Header:** negro, logo de 44 px, buscador y carrito con contador flúor.

## 9. Checklist rápido antes de dar por buena una pantalla
- [ ] ¿El dueño reconocería el negro y oro, y el "¡Llegó…!"?
- [ ] ¿Está más ordenado que un video? (una oferta principal por pantalla)
- [ ] ¿Cada precio tiene fuente o dice `[EJEMPLO]`?
- [ ] ¿Se ven detalle y embalaje o mayor?
- [ ] ¿La dirección y el horario están visibles y con fecha?
- [ ] ¿Pasa el contraste? (nada de oro o blanco sobre blanco u oro)
- [ ] ¿Máximo 1 emoji por titular, y sin 🏆?
- [ ] ¿Funciona a 360 px de ancho y en el navegador interno de TikTok?
