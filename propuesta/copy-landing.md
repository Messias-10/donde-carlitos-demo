# Copy de la landing — Donde Carlitos Emporio

*estratega-propuesta · Fase 2 · 2026-11-14 · voz: skill `estilo-carlitos`*

## 1. Convenciones (léelas antes de copiar)

**Cómo leer este archivo:**
- **Texto normal** es el texto visible, final. Se copia tal cual.
- **Las líneas `> nota:`** son para el equipo y **no se muestran**. Traen la fuente, la regla y el campo de la API.
- **`[CONFIRMAR CON EL DUEÑO]`** marca un dato contradictorio o no confirmado. En el demo se muestra la versión indicada, con una etiqueta chica gris "dato por confirmar". En la reunión sirve de pregunta (ver `pitch.md`).
- **`[EJEMPLO]`** es un dato ilustrativo, que no sale de un post real.
- **Llaves `{campo}`:** valores que pone el front desde la API.

**Formato:**
- **Titulares:** en MAYÚSCULAS con Anton; máximo 1 emoji por titular y nunca 🏆.
- **Script Kaushan:** solo donde dice *script* (2–3 palabras).
- **Precios:** `$4.000`, con punto de miles y sin decimales. En titulares se puede decir "luquitas".
- **Botones de acción crítica:** sin emojis (salvo el ícono de WhatsApp).

### 1.1 Metadatos (SEO y vista previa al compartir)
- **title:** `Donde Carlitos Emporio · Alameda 2432, Santiago Centro`
- **description:** `Juguetes virales, Navidad, hogar y escolar al detalle y por mayor. Arma tu pedido y mándalo a nuestro WhatsApp oficial. Te esperamos en Alameda 2432, a pasos del Metro Unión Latinoamericana.`
- **og:title:** `¿Lo viste en TikTok? Acá lo pides · Donde Carlitos Emporio`
- **og:description:** igual que *description*.

> nota: la dirección y el Metro salen de las bios oficiales (`pendientes.md` §3).

---

## 2. Hero

### Header
- Logo (alt): `Donde Carlitos Emporio`
- Buscador, placeholder: `¿Qué andas buscando? Ej: squishy, tren, telón`
- Carrito (aria-label): `Ver mi pedido, {n} productos`

### Titular elegido
*script:* **¡Hola chiquillos!**

# ¿LO VISTE EN TIKTOK? ACÁ LO PIDES.

**Bajada:** Precio al detalle y por mayor, a la vista. Arma tu pedido y mándalo directo a nuestro WhatsApp oficial.

**Botones:**
- principal: `Ver ofertas`
- secundario: `Ver catálogo`

**Barra de datos (debajo de los botones):**
`📍 Alameda 2432, Santiago Centro · a 5 min del Metro Unión Latinoamericana`
`🕘 Hoy: 9:00 a 18:00` · `Abierto ahora` / `Cerrado ahora`

> nota: el horario es el de nov 2026 (tr. TT-7707471893363789063) **[CONFIRMAR CON EL DUEÑO]**. Mientras no esté confirmado, se muestra sin "Abierto/Cerrado" y con la etiqueta "dato por confirmar". El domingo, según esa fuente, es de 9:00 a 16:00.

**Por qué este titular:**
1. **Le habla a quien llega:** la mayoría de las visitas viene del navegador interno de TikTok o IG (CLAUDE.md, regla 5).
2. **Promete solo lo que el demo hace:** ver y pedir.
3. **Sirve para los dos perfiles:** el que vio el viral y el revendedor.
4. **Es corto:** 6 palabras, legible a 40–48 px en 360 px de ancho.
5. **Usa la frase de la propia guía de estilo:** "¿Lo viste en TikTok? Acá lo tienes con precio y stock" (estilo-carlitos §4, frase 6).

### Alternativa 1: novedad primero
*script:* **Familia Emporio**
# ¡LLEGÓ LO NUEVO, CHIQUILLOS! 👇
**Bajada:** Navidad, juguetes virales y hogar, con precio al detalle y por mayor. Arma tu pedido y mándalo a nuestro WhatsApp oficial.

> nota: es la voz más pura de los videos ("Hola chiquillos, llegaron…", `adn-marca.md` §B). Es más débil para explicar cómo se compra.

### Alternativa 2: su promesa de surtido
*script:* **Donde Carlitos**
# LA MEJOR VARIEDAD DE MEIGGS, CON EL PRECIO A LA VISTA
**Bajada:** Lo que ves en nuestros videos, con precio y stock. Arma tu pedido y te lo confirmamos por WhatsApp.

> nota: "la mejor variedad de productos de barrio Meiggs" es su propia frase de cierre (tr. TT-7707649267464293639). Es una comparación, así que solo va si el dueño la aprueba **[CONFIRMAR CON EL DUEÑO]**. Ojo: el lema "El surtido más grande y los precios más económicos del mercado" del CLAUDE.md no se vio en las bios oficiales; no se usa hasta confirmarlo.

---

## 3. Ofertas / Remate de Navidad

### Banner (uno solo, sin carrusel)
*script:* **Donde**
## ¡REMATE DE NAVIDAD!
Esferas con luces, telones, trenes, lazos y bolsas de regalo.
Esferas con luces a **$500**.
`📍 Alameda 2432`
- botón: `Ver ofertas`
- botón: `Escríbenos por WhatsApp`

> nota: el remate sale de "Comenzamos con los remates de Navidad" (tr. TT-7707268821178141970). Las esferas a $500 salen de TT-7707647176440171783 (13-11-2026), y era el precio para quien comentaba. En el demo va como precio detalle, **sin tachado** porque no hay precio anterior real. Si el dueño fija una fecha de término, se agrega `🔥 Hasta el {día dd/mm}`. Si no, no va urgencia.

### Fila de ofertas
## OFERTAS QUE VUELAN 🔥
Remates y promos de la semana. Cuando se acaban, se acaban.

> nota: "OFERTAS QUE VUELAN" está en la descripción de TT-7592989161653685522. Para la fila se filtra con `SoloOferta`. Ejemplo con tachado real: **Tren navideño para árbol (grande)**, ~~$20.000~~ **$17.000**, $15.000 por embalaje (TT-7706656443843284232, 10-11-2026). Su "válido solo hasta el domingo" ya venció (domingo 15-11-2026), así que no se muestra salvo que el dueño ponga una fecha nueva.

---

## 4. Lo viste en TikTok

## ¿LO VISTE EN TIKTOK?
Acá lo tienes con precio y stock. Toca el video para verlo de nuevo o agrégalo al tiro a tu pedido.

- Etiqueta de la miniatura: `▶ Míralo en TikTok`
- Sobre los videos con dirección antigua: `Hoy nos encuentras en Alameda 2432`
- Debajo va la tarjeta del producto (§15.1).

> nota: los videos sugeridos están en `arquitectura-landing.md` §3. El link abre en una pestaña nueva. Créditos en `web/public/CREDITOS.md`.

---

## 5. Categorías

## ¿QUÉ ANDAS BUSCANDO?
- `Navidad`
- `Juguetes y virales`
- `Hogar y organización`
- `Escolar`
- `Textil`
- `Electro`
- `Cuidado personal`
- `Cotillón`

**Bajo cada nombre:** `{n} productos`

> nota: se ocultan las categorías con 0 productos.

---

## 6. ¡Llegó esta semana!

## ¡NOS LLEGÓ! ✨
Lo último que bajó del contenedor. Cuando se acaba, se acaba.

- Badge: `¡LLEGÓ!`
- Bajo el nombre: `Llegó el {dd/mm}`
- Botón al final de la fila: `Ver todo lo nuevo`

> nota: usa `fecha_llegada`. La frase "cuando se acaban, se acaban" es de estilo-carlitos §4.

---

## 7. ¿Pa negocio?

## ¿PA NEGOCIO?
Cada producto trae su precio al detalle y su precio por mayor, con el mínimo a la vista.

**Ejemplo:** Squishy viral que cambia de color con el sol: **$4.000** al detalle · **$3.500** por caja de 12.

- En el carrito, el precio por mayor se aplica solo cuando llegas al mínimo.
- ¿Necesitas más cantidad? Escríbelo en el comentario del pedido y te cotizamos.

Botón: `Ver precios por mayor`

> nota: el ejemplo sale de TT-7707649267464293639 (13-11-2026). "Se aplica solo" depende del cotizador de la API **[VERIFICAR con ingeniero-go]**. **No** se menciona un mínimo de compra general por mayor: no está confirmado (`pendientes.md` §3).

---

## 8. Cómo pedir en 3 pasos

## PEDIR ES ASÍ DE FÁCIL 👇

**1. ELIGE**
Agrega lo que quieres al carrito. Si llevas cantidad, el precio por mayor se aplica solo.

**2. ENVÍA**
Deja tu nombre, tu WhatsApp y si retiras o te lo enviamos. Tu pedido queda con número y se abre WhatsApp con todo escrito: solo aprietas enviar.

**3. TE CONFIRMAMOS**
Te respondemos desde nuestro número oficial con el total y el pago. Después retiras en Alameda 2432 o lo despachamos.

*Pie chico:* La web no cobra nada. El pago lo coordinamos contigo por WhatsApp o en el local.

> nota: "lo despachamos" depende de la política de envíos **[CONFIRMAR CON EL DUEÑO]**. Si el dueño no despacha, el paso 3 queda así: "Te respondemos desde nuestro número oficial y retiras en Alameda 2432."

---

## 9. Retiro en tienda y envíos

## RETIRA O TE LO ENVIAMOS

**RETIRO EN TIENDA**
Haz tu pedido, espera nuestra confirmación por WhatsApp y pasa a buscarlo a Alameda 2432 con tu número de pedido.
Llegas sabiendo qué llevas y cuánto pagas.

**ENVÍO A REGIONES**
¿Eres de región? Arma tu pedido, elige tu agencia y te confirmamos por WhatsApp el total y el despacho.
`[CONFIRMAR CON EL DUEÑO: agencias, quién paga el envío y monto mínimo]`

> nota: la fuente más reciente es nov 2026, "número para región" (tr. TT-7707471893363789063), más el comentario de "la niña de los pedidos de regiones" (TT-7706913997378948370). En 2025 y abril 2026 decían "solamente venta presencial". **No escribir "Envíos a todo Chile".**
> Si el dueño confirma "por pagar": agregar "El envío lo pagas tú directo a la agencia cuando llega."
> Si el dueño dice que no despacha: reemplazar la tarjeta por "Por ahora vendemos solo en el local. Haz tu pedido y retíralo en Alameda 2432." y ocultar la opción "Envío" en el checkout.

---

## 10. Cómo llegar

## CÓMO LLEGAR 📍

**Alameda 2432, Santiago Centro**

- **Metro Unión Latinoamericana** (Línea 1): a 5 minutos caminando.
- **Metro República** (Línea 1): a 7 minutos caminando.
- **¿Vienes en bus desde región?** Estás a una estación de Metro Estación Central (Terminal San Borja) y a dos de Universidad de Santiago (terminales Alameda y Sur).

Botón: `Abrir en Google Maps`

> nota: los minutos salen del CLAUDE.md y las bios (reconfirmar). Los terminales salen de [6][7] de `contexto-meiggs.md`. Las estaciones de distancia salen del orden de la Línea 1 *(inferencia)*. No se enlaza el reel "cómo llegar en Metro" de jun 2025 (TT-7512969200437120262): es anterior al cambio a Alameda 2432.

### Horario
🕘 **Horario de temporada**
- Lunes a viernes: 9:00 a 18:00
- Sábado, domingo y feriados: 9:00 a 16:00

*Publicado en noviembre de 2026.* `[CONFIRMAR CON EL DUEÑO]`

> nota: hay tres versiones (`adn-marca.md` §Horario). Se usa la más reciente (tr. TT-7707471893363789063 y TT-7706656443843284232).

### Cómo pagar en el local
En Alameda 2432 tenemos dos locales:
- **Hogar, electro y juguetes:** efectivo y débito.
- **Disney y licencias:** débito, crédito, transferencia y efectivo.

Boleta o factura, como necesites.
`[CONFIRMAR CON EL DUEÑO]`

> nota: sale de tr. TT-7707464302034160914 (nov 2026): "boleta, factura lo que usted quiera".

---

## 11. Canal oficial

## ESTE ES NUESTRO ÚNICO CANAL OFICIAL

**Pedidos solo a este número:** +56 9 2114 5368
Botón: `Escribir a +56 9 2114 5368`

**En redes nos encuentras como Mayte Emporio:**
- TikTok **@mayte_emporio**
- Instagram **@mayte_emporioo** (cuenta verificada)

**Para comprar seguro:**
1. Los pedidos se hacen en esta web o al +56 9 2114 5368. Ningún otro número.
2. Nunca te escribimos primero por interno para venderte. `[CONFIRMAR CON EL DUEÑO la redacción]`
3. Nunca te vamos a pedir claves ni códigos que te lleguen por SMS.
4. Antes de pagar, revisa que el chat sea con nuestro número oficial.
5. ¿Te escribieron a nombre de nosotros desde otra cuenta o número? No transfieras y avísanos al número oficial.

> nota: "pedidos solo a este número" sale de la bio de IG (`adn-marca.md` §H). La regla 2 viene del CLAUDE.md, pero no se vio en las bios oficiales. Las reglas 3 a 5 son consejos de seguridad, no afirmaciones sobre el negocio. **No se nombran ni enlazan** @_dondecarlitosemporio7 ni páginas de Facebook. El nombre "Mayte Emporio" frente a "Donde Carlitos Emporio" está pendiente (`pendientes.md` §4). En F2 se agrega: "Nuestra única web es {dominio}".

**Versión corta (footer y confirmación):**
`Pedidos solo al +56 9 2114 5368. Si te escriben desde otro número a nombre nuestro, no transfieras.`

---

## 12. Preguntas frecuentes

## PREGUNTAS QUE NOS HACEN SIEMPRE 💬

> nota: van en orden de frecuencia según `faq-clientes.md`. Cada pregunta está escrita como la hace la gente, pero ordenada.

**¿Dónde están? ¿Qué metro queda cerca?**
En Alameda 2432, Santiago Centro, a 5 minutos del Metro Unión Latinoamericana y a 7 del Metro República. Mira el mapa en "Cómo llegar".

**¿Hacen envíos a regiones?**
Sí: arma tu pedido, elige "Envío por agencia" y te confirmamos por WhatsApp el total y el despacho. `[CONFIRMAR CON EL DUEÑO: agencias, quién paga y mínimo]`
> nota: si no despachan, la respuesta queda así: "Por ahora vendemos solo en el local de Alameda 2432."

**¿Les queda…? ¿Cómo sé si hay stock?**
Cada producto dice si está **Disponible**, si **quedan pocas unidades** o si está **Agotado**. Si no lo encuentras en la web, pregúntanos por WhatsApp: muchas veces vuelve en el próximo contenedor.
> nota: "vuelve en el próximo contenedor" es *(inferencia)* de la rotación rápida [18]. Si el dueño no lo quiere decir, cortar después de "por WhatsApp".

**¿Cuántos vienen por embalaje? ¿Desde cuántas es precio por mayor?**
Está en cada producto: el precio al detalle, el precio por mayor y desde cuántas unidades aplica (por ejemplo, "caja de 12"). En el carrito se aplica solo.

**¿Cómo compro? ¿Dónde les escribo?**
En 3 pasos: eliges, envías tu pedido y te confirmamos por WhatsApp. Nuestro único número de pedidos es el +56 9 2114 5368.

**¿Hasta qué hora atienden? ¿Abren domingo y feriados?**
En temporada: lunes a viernes de 9:00 a 18:00, y sábado, domingo y feriados de 9:00 a 16:00. `[CONFIRMAR CON EL DUEÑO]`

**¿Cuánto vale?**
Todos los productos tienen su precio en la web, al detalle y por mayor cuando aplica. Tu pedido guarda el precio del momento en que lo hiciste. `[CONFIRMAR CON EL DUEÑO que ese precio se respeta]`

**¿Qué medios de pago aceptan? ¿Puedo pagar con crédito?**
Depende del local. En hogar, electro y juguetes: efectivo y débito. En Disney y licencias: débito, crédito, transferencia y efectivo. `[CONFIRMAR CON EL DUEÑO]`

**¿Cómo funciona "Yo quiero mi promo"?**
En la web, los productos en promo ya muestran su precio especial. No tienes que comentar nada: lo agregas y listo. `[CONFIRMAR CON EL DUEÑO cómo se canjea la promo del video en el local]`
> nota: la regla de canje no está explicada en los videos (`pendientes.md` §3). Si el dueño no quiere que la web aplique la promo del comentario, esta pregunta se saca.

**¿Tienen tallas?**
Las tallas disponibles van en la descripción de cada producto. Si no ves la tuya, escríbela en el comentario del pedido o pregúntanos por WhatsApp.
> nota: el demo no maneja variantes por talla. Las preguntas de talla vienen de uniformes escolares (TT-7592989161653685522).

**¿Es seguro comprar con envío?**
Haz tu pedido en esta web o al +56 9 2114 5368 y paga solo cuando te confirmemos desde ese número. Si alguien te escribe desde otra cuenta a nombre nuestro, no transfieras.

**¿Dan boleta o factura?**
Sí, boleta o factura, como necesites. `[CONFIRMAR CON EL DUEÑO]`

**¿Hacen cambios?**
`[COMPLETAR CON EL DUEÑO: política de cambios y devoluciones]`
> nota: no hay dato (`pendientes.md` §3). Si no se completa, la pregunta no se publica.

**¿Tienen local en regiones?**
Sí, en Temuco: Manuel Rodríguez 1232. `[CONFIRMAR CON EL DUEÑO si se incluye y su horario]`
> nota: la fuente es la descripción de TT-7592989161653685522 (ene 2026). Si no se aprueba, la pregunta se saca.

---

## 13. Familia Emporio (opcional)

*script:* **Familia Emporio**
## DEL TOLDO AMARILLO AL CONTENEDOR
Partimos con un toldo amarillo. Hoy abrimos contenedores en plena calle y te dejamos lo nuevo a precio de remate.
Botón: `▶ Mira cómo abrimos un contenedor`

> nota: el "toldo amarillo" y el "local amarillo" son de sus propios hashtags (`adn-marca.md` §A). El video sugerido es TT-7441724226857454904 o TT-7683668603270352136. Van sin nombres de personas, cifras de la empresa ni historia personal. **[CONFIRMAR CON EL DUEÑO que quiere contar esto]**.

---

## 14. Footer y botón flotante

### Footer
- **Marca:** Donde Carlitos Emporio · *Familia Emporio*
- `📍 Alameda 2432, Santiago Centro`
- `🕘 L a V 9:00–18:00 · Sáb, dom y feriados 9:00–16:00 (temporada)`
- **Pedidos solo al +56 9 2114 5368**
- **Redes:** TikTok @mayte_emporio · Instagram @mayte_emporioo
- **Canal oficial (corto):** `Si te escriben desde otro número a nombre nuestro, no transfieras.`
- **Cierre:** *Los esperamos en Alameda 2432 con la mejor variedad de barrio Meiggs.*
- **Precios:** `Los precios pueden cambiar. Tu pedido guarda el precio del momento en que lo hiciste.`
- **Solo en el demo:** `Demo de prueba: los pedidos no son reales. Precios sacados de nuestros videos de TikTok (nov 2026).`

> nota: el cierre es su despedida en los videos (tr. TT-7707649267464293639, TT-7706989825747995911). La línea de precios queda sujeta a la confirmación de la política (§12, "¿Cuánto vale?").

### Botón flotante de WhatsApp
- **Sin carrito:** `Escríbenos por WhatsApp`
- **Con carrito:** `Ver mi pedido ({n}) · {total}`
- **Debajo, en chico:** `Único número de pedidos`

> nota: con carrito lleva al carrito, no al chat (ver `arquitectura-landing.md` §14). El mensaje de consulta está en `mensaje-pedido-whatsapp.md` §5.

---

## 15. Microcopy de compra

### 15.1 Tarjeta de producto
**Badges:**
- `REMATE` (si `etiqueta_oferta = remate`)
- `OFERTA`
- `¡LLEGÓ!`
- `NAVIDAD` / `ESCOLAR` / `HALLOWEEN` / `FIESTAS PATRIAS` / `INVIERNO` / `VERANO` (según `temporada`)
- `PRECIO X MAYOR` (si tiene `precio_mayor`)

**Precios:**
- Precio principal: `{precio}` + debajo, en chico: `detalle`
- Con oferta: ~~`{precio_detalle}`~~ `{precio_oferta}`
- Por mayor: `{precio_mayor} x mayor desde {mayor_minimo} u.`, y si hay unidad: `· {unidad_mayor}`. Ejemplo: `$3.500 x mayor desde 12 u. · caja de 12`

**Disponibilidad:**
- `Disponible`
- `Quedan pocas unidades`
- `Agotado`

**Botones:**
- `Agregar`
- En oferta: `¡Yo quiero mi promo!` (hace lo mismo que Agregar) `[CONFIRMAR CON EL DUEÑO; si no, usar "Agregar"]`
- Agotado: `¿Vuelve? Pregúntanos`. Abre WhatsApp con el mensaje de consulta por producto.
- Link al video: `▶ Míralo en TikTok`
- Al agregar (toast): `¡Listo! Agregaste {nombre}.` · `Ver mi pedido`

> nota: en el demo no hay "Avísame cuando llegue" automático (es F3); por eso el botón de agotado lleva a WhatsApp.

### 15.2 Ficha de producto
- Cantidad: `Cantidad`
- Precio por mayor: `Llevando {mayor_minimo} o más, pagas {precio_mayor} c/u.`
- Si hay unidad: `Viene en {unidad_mayor}.`
- Si tiene precio de demo: `Precio referencial (video del {dd/mm/aaaa}).`
- Si es ejemplo: `Producto de ejemplo`

> nota: el campo `es_ejemplo` de la API marca los productos de ejemplo.

### 15.3 Carrito
- **Título:** `TU PEDIDO`
- **Carrito vacío:** `Tu pedido está vacío. ¡Mira lo que llegó!` · botón `Ver ofertas`
- **Línea:** `{cantidad} x {precio_unitario} ({detalle | x mayor | oferta})` · `{subtotal}`
- **Aviso de escalera:** `Te faltan {n} para el precio por mayor ({precio_mayor} c/u).`
- **Al llegar al mínimo:** `Precio por mayor aplicado.`
- **Total:** `Total productos` · `{total}` · `(no incluye envío)`
- **Pie:** `La web no cobra. El pago lo coordinamos contigo por WhatsApp o en el local.`
- **Botones:** `Continuar` · `Seguir comprando` · `Quitar`

### 15.4 Datos del pedido
- **Título:** `TUS DATOS`
- **Nombre:** label `Tu nombre`, placeholder `Ej: Camila Rojas`
- **Teléfono:**
  - label `Tu WhatsApp`, placeholder `9 1234 5678`;
  - ayuda `Te confirmamos el pedido a este número.`
- **Entrega:** label `¿Cómo lo quieres?`
  - `Retiro en tienda · Alameda 2432`
  - `Envío por agencia (regiones)` `[CONFIRMAR CON EL DUEÑO]`
- **Región:** `Región`, placeholder `Elige tu región`
- **Comuna:** label `Comuna`, placeholder `Ej: Chillán`
- **Agencia:**
  - label `Agencia`, placeholder `Elige tu agencia`;
  - ayuda `El costo de envío te lo confirmamos por WhatsApp.`
  - > nota: las agencias del selector son `[EJEMPLO]` hasta que el dueño confirme.
- **Comentario:**
  - label `¿Algo más? (opcional)`;
  - placeholder `Ej: talla 10, paso el sábado en la mañana, es para mi negocio…`
- **Privacidad:** `Usamos tus datos solo para coordinar este pedido.`
- **Botón:** `Hacer pedido`

### 15.5 Errores
- **Validación del formulario:**
  - `Escribe tu nombre.`
  - `Revisa tu número: tiene que ser un celular chileno de 9 dígitos.`
  - `Elige si retiras o te lo enviamos.`
  - `Elige tu región, tu comuna y la agencia.`
- **Stock:** `No nos queda stock suficiente de {nombre}. Baja la cantidad o sácalo del pedido.`
- **Producto que ya no está:** `{nombre} ya no está disponible. Lo sacamos de tu pedido.`
- **Demasiados pedidos (límite por IP):** `Hiciste varios pedidos seguidos. Espera unos minutos o escríbenos al +56 9 2114 5368.`
- **Error general:** `No pudimos guardar tu pedido. Inténtalo de nuevo en un rato o escríbenos al +56 9 2114 5368.`
- **Sin conexión:** `Parece que no tienes internet. Tu carrito sigue guardado.`

> nota: los mensajes no muestran detalles internos (CLAUDE.md, Seguridad).

### 15.6 Confirmación
## ¡LISTO! TU PEDIDO ES EL {numero}
Falta un paso: mándanos el pedido por WhatsApp para que te lo confirmemos.

- **Botón principal:** `Enviar pedido por WhatsApp`
- **Respaldo:**
  - texto `¿No se abrió WhatsApp? Copia el mensaje y pégalo en un chat con el +56 9 2114 5368.`;
  - botón `Copiar mensaje`, que cambia a `¡Copiado!`.
- **Recordatorio:** `Guarda tu número de pedido: {numero}`
- **Si es retiro:** `Espera nuestra confirmación antes de venir. Te esperamos en Alameda 2432 con tu número de pedido.` + horario.
- **Si es envío:** `Te confirmamos por WhatsApp el total, el pago y el despacho.`
- **Canal oficial:** `Te respondemos solo desde el +56 9 2114 5368. No transfieras a otro número ni a otra cuenta.`
- **Botón:** `Volver a la tienda`

---

## 16. Panel de administración (microcopy)

### Login
- **Título:** `PANEL DE PEDIDOS`
- **Campos y botón:** `Usuario` · `Contraseña` · `Entrar`
- **Errores:**
  - `Usuario o contraseña incorrectos.`
  - `Demasiados intentos. Espera unos minutos y vuelve a probar.`
- **Botón de salida:** `Cerrar sesión`

### Pedidos
**Pestañas:**
- `Nuevos ({n})`
- `Confirmados`
- `Listos para retiro`
- `Despachados`
- `Entregados`
- `Cancelados`
- `Todos`

**Tarjeta de pedido:** `{numero}` · `{nombre}` · `{Retiro | Envío a {comuna}, {region}}` · `{total}` · `hace {tiempo}`

**Estados (etiqueta y qué significa):**

| Estado API | Etiqueta | Qué significa |
|---|---|---|
| `nuevo` | Nuevo | Llegó desde la web y nadie lo ha revisado |
| `confirmado` | Confirmado | Revisamos stock y total con el cliente |
| `listo` | Listo para retiro | Está armado en el local |
| `despachado` | Despachado | Salió por la agencia |
| `entregado` | Entregado | El cliente lo retiró o lo recibió |
| `cancelado` | Cancelado | No va. El stock vuelve al catálogo |

**Acciones:**
- `Confirmar`
- `Marcar listo para retiro`
- `Marcar despachado`
- `Marcar entregado`
- `Cancelar pedido`
- **Nota del cambio:** `Nota (opcional)`, placeholder `Ej: pagó por transferencia, retira el lunes`
- **Al cancelar:** `¿Cancelar el pedido {numero}? Las unidades vuelven al stock.` · `Sí, cancelar` · `No`
- **Opcional (si constructor-web lo incluye; es solo un link `wa.me` al teléfono del cliente):** `Abrir chat con el cliente`, con el mensaje `Hola {nombre}, te escribimos de Donde Carlitos Emporio por tu pedido {numero}.`
- **Historial:** `Historial` · `{fecha hora} · {estado anterior} → {estado nuevo} · {usuario}`
- **En el demo:** `Pedido de prueba (demo)`

> nota: que el stock vuelva al cancelar depende del dominio de la API (`ValidadorCambio` / `restaurarStock`).

### Productos
- **Título y buscador:** `PRODUCTOS` · `Buscar producto`
- **Campos:**
  - `Precio detalle`
  - `Precio oferta (opcional)`
  - `¿Oferta o remate?`
  - `Precio por mayor (opcional)`
  - `Desde cuántas unidades`
  - `Unidad (ej: caja de 12)`
  - `Stock`
  - `Mostrar "quedan pocas" con menos de`
  - `Temporada`
  - `Fecha de llegada`
  - `Link al video de TikTok`
  - `Visible en la web`
- **Botón:** `Guardar`, que cambia a `Guardado`
- **Errores (calzan con las reglas de la base):**
  - `El precio de oferta tiene que ser menor que el precio detalle.`
  - `El precio por mayor tiene que ser menor que el precio detalle.`
  - `Si pones precio por mayor, indica desde cuántas unidades (2 o más).`
  - `Alguien más cambió este producto. Recarga y vuelve a intentar.`
