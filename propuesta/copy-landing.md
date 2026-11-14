# Copy: landing, seguimiento y panel — Donde Carlitos Emporio

*estratega-propuesta · 2026-11-14 · voz: skill `estilo-carlitos` · enfoque "Tu WhatsApp sigue igual, pero cada pedido queda ordenado"*

## 1. Convenciones

**Cómo leer este archivo:**
- **Texto normal** es texto visible y final. Se copia tal cual.
- **Las líneas `> nota:`** son para el equipo y **no se muestran**.
- **`[CONFIRMAR CON EL DUEÑO]`** marca un dato no confirmado. En el demo se muestra el **texto de demo** que se indica, con una etiqueta gris chica: `dato por confirmar`. En la reunión sirve de pregunta (`pitch.md`).
- **`[EJEMPLO]`** es un dato ilustrativo.
- **`{campo}`** es un valor que pone el front o la API.

**Formato:**
- **Titulares:** en MAYÚSCULAS con Anton; máximo 1 emoji por titular y nunca 🏆.
- **Script Kaushan:** solo donde dice *script*.
- **Precios:** `$4.000`, con punto de miles. En titulares se puede decir "luquitas".
- **Botones:** sin emojis, salvo el ícono de WhatsApp.
- **Mensajes de WhatsApp:** no están acá; van en `mensajes-whatsapp.md`.

**Ya no existen** (se borran del front): carrito, "Agregar", "Ver mi pedido (n)", checkout, "Tus datos", "Hacer pedido", confirmación "¡Listo! Tu pedido es el…" y el buscador del catálogo.

### 1.1 Metadatos
- **title:** `Donde Carlitos Emporio · Pedidos por WhatsApp · Alameda 2432`
- **description:** `Juguetes virales, Navidad, hogar y escolar, al detalle y por mayor. Pide por WhatsApp para retiro en Santiago o con envío a regiones. Alameda 2432, a pasos del Metro Unión Latinoamericana.`
- **og:title:** `¿Lo viste en TikTok? Pídelo por WhatsApp · Donde Carlitos Emporio`
- **og:description:** igual que *description*.

> nota: "con envío a regiones" se apoya en el "número para región" (tr. TT-7707471893363789063) y en la "niña de los pedidos de regiones" (TT-7706913997378948370). Si el dueño dice que no despacha, cambiar a "Pide por WhatsApp y retira en Alameda 2432".

---

## 2. Hero

### Header
- Logo (alt): `Donde Carlitos Emporio`
- Link: `Sigue tu pedido`

### Titular elegido
*script:* **¡Hola chiquillos!**

# ¿LO VISTE EN TIKTOK? PÍDELO POR WHATSAPP.

**Bajada:**
Juguetes virales, Navidad, hogar y escolar, al detalle y por mayor.
Dinos si retiras en tienda o eres de regiones, y te abrimos el chat correcto con el mensaje listo.

**Botones (uno debajo del otro):**
- `Pedidos Santiago / retiro en tienda`
- `Pedidos a regiones`

**Debajo, en chico:** `Estos son nuestros únicos números de pedidos.` · link `Ver canal oficial`

**Barra de datos:**
`📍 Alameda 2432, Santiago Centro · a 5 min del Metro Unión Latinoamericana`
`🕘 Hoy: 9:00 a 18:00`

> notas:
> - El horario es el de nov 2026 (tr. TT-7707471893363789063) **[CONFIRMAR CON EL DUEÑO]**. Va sin "Abierto/Cerrado" hasta que se confirme, con la etiqueta `dato por confirmar`.
> - Las categorías de la bajada salen de `adn-marca.md` §E.
> - Si `WHATSAPP_REGIONES` está vacío, el botón igual dice "Pedidos a regiones", abre el número de Santiago con el mensaje (b) y lleva la etiqueta `dato por confirmar` en el demo.

**Por qué este titular:**
1. **Le habla a quien llega:** la mayoría viene desde un video de TikTok o IG (CLAUDE.md, regla 5).
2. **Promete solo lo que pasa:** el botón abre WhatsApp. Nada de "compra online".
3. **Respeta cómo venden hoy:** WhatsApp es su canal y el titular lo dice.
4. **Es corto:** 6 palabras, legible a 40–48 px en 360 px.
5. **Da continuidad** al titular anterior, "¿Lo viste en TikTok? Acá lo pides", ajustado al nuevo enfoque.

### Alternativa 1: el canal primero
*script:* **Familia Emporio**
# ¿DE SANTIAGO O DE REGIÓN? ESCRÍBENOS AL TIRO.
**Bajada:** Cada pedido llega al WhatsApp que corresponde, con el mensaje ya escrito. Así te atendemos más ordenado.

> nota: es la más clara para la decisión de los dos botones, pero no dice qué venden. "Te atendemos más ordenado" no promete rapidez.

### Alternativa 2: novedad primero
*script:* **Donde Carlitos**
# ¡LLEGÓ LO NUEVO, CHIQUILLOS! 👇
**Bajada:** Mira lo último de nuestros videos y pídelo por WhatsApp: retiro en Alameda 2432 o envío a regiones.

> nota: es la voz más pura de los videos (`adn-marca.md` §B), pero explica menos cómo se compra.

---

## 3. Lo nuevo de los videos (vitrina)

## LO NUEVO DE LOS VIDEOS ✨
Lo que mostramos esta semana en TikTok, con su precio. Tócalo y nos escribes con el producto ya puesto.

**Tarjeta:**
- Badges posibles: `¡LLEGÓ!` · `REMATE` · `OFERTA` · `PREVENTA · llega aprox. {dd/mm}` · `PRECIO X MAYOR`
- Precio: `{precio}` + en chico `detalle`
- Con oferta vigente: ~~`{precio_anterior}`~~ `{precio}`
- Por mayor: `{precio_mayor} c/u · {unidad_mayor}`
- Fecha: `Precio publicado el {dd/mm}`
- Link: `▶ Míralo en TikTok`
- Botón: `Pedir este por WhatsApp`

**Hoja "¿Dónde lo quieres?":**
- Título: `¿DÓNDE LO QUIERES?`
- Botón 1: `Santiago / retiro en tienda`
- Botón 2: `Regiones (envío)`
- Pie: `Te abrimos WhatsApp con el producto ya escrito. Tú solo envías.`
- Cerrar: `Cancelar`

**Pie de la sección:** `Los precios son los que dijimos en el video y pueden cambiar. Pregúntanos por WhatsApp si queda.`
**Al final:** `¿Buscas otra cosa? Escríbenos y te decimos si la tenemos.` · botones `Santiago` / `Regiones`

> notas:
> - Los productos sugeridos, con su fuente, están en `arquitectura-landing.md` §A2.
> - La urgencia ("solo hasta el domingo") solo va si `oferta_hasta` está vigente.
> - No se muestra stock.

---

## 4. Cómo comprar desde regiones

## ¿ERES DE REGIÓN? ASÍ COMPRAS 📦

**1. ESCRÍBENOS**
Toca "Pedidos a regiones" y cuéntanos tu comuna y lo que quieres. Si lo viste en un video, mándanos el nombre o el link.

**2. TE COTIZAMOS**
Te mandamos por el chat el resumen con tu número de pedido y el total.

**3. TRANSFIERES**
Los datos para transferir te los damos **solo por nuestro chat oficial**. No están publicados en ninguna parte.

**4. NOS MANDAS EL COMPROBANTE**
Por el mismo chat. Revisamos que el pago llegó y te confirmamos.

**5. DESPACHAMOS**
Por agencia, y te mandamos el número de seguimiento y un link para ver en qué va tu pedido.

Botón: `Pedidos a regiones`

### Lo que tienes que saber
| Dato | Texto final (cuando el dueño confirme) | Texto de demo |
|---|---|---|
| Agencias | `Despachamos por {agencias}.` `[CONFIRMAR CON EL DUEÑO]` | `Agencias: por confirmar` |
| Costo del envío | `El envío lo pagas tú en la agencia al recibirlo (por pagar).` `[CONFIRMAR CON EL DUEÑO]` | `Costo de envío: te lo confirmamos en la cotización` |
| Compra mínima | `Compra mínima para envío: {monto}.` `[CONFIRMAR CON EL DUEÑO]` | `Compra mínima: por confirmar` |
| Plazo | `Despachamos dentro de {n} días hábiles desde que confirmamos tu pago.` `[CONFIRMAR CON EL DUEÑO]` | `Plazo de despacho: por confirmar` |
| Medio de pago | `Transferencia bancaria.` `[CONFIRMAR CON EL DUEÑO]` | `Pago por transferencia (por confirmar)` |

> notas:
> - Los 5 pasos son el flujo que aprobó Luckas. Que el pago sea por transferencia es *(inferencia)*: es el medio más aceptado en Chile [26] y no hay link de pago. "Por pagar" es lo habitual en el barrio [11], pero no está confirmado para Carlitos.
> - **No escribir "Envíos a todo Chile".**
> - Si el dueño dice que no despacha: se saca esta sección y el botón de regiones, y se deja "Por ahora vendemos solo en el local".

---

## 5. ¿Eres de Santiago?

## ¿ERES DE SANTIAGO? RETIRA EN TIENDA 📍
1. Escríbenos por WhatsApp o ven directo al local.
2. Te confirmamos el total y te avisamos cuando tu pedido esté listo.
3. Lo retiras en Alameda 2432 con tu número de pedido.

Botón: `Pedidos Santiago / retiro en tienda`

> nota: "te avisamos cuando esté listo" depende de que puedan preparar pedidos con anticipación **[CONFIRMAR CON EL DUEÑO]**. Si no, el paso 2 queda así: "Te confirmamos el total y si queda stock".

---

## 6. Cómo llegar

## CÓMO LLEGAR 📍

**Alameda 2432, Santiago Centro**
- **Metro Unión Latinoamericana** (Línea 1): a 5 minutos caminando.
- **Metro República** (Línea 1): a 7 minutos caminando.
- **¿Vienes en bus desde región?** Estás a una estación de Metro Estación Central (Terminal San Borja) y a dos de Universidad de Santiago (terminales Alameda y Sur).

Botón: `Abrir en Google Maps`

> nota: las fuentes son las bios oficiales y [6][7]. La distancia en estaciones se sacó del orden de la Línea 1 *(inferencia)*.

### Horario
🕘 **Horario de temporada**
- Lunes a viernes: 9:00 a 18:00
- Sábado, domingo y feriados: 9:00 a 16:00

*Publicado en noviembre de 2026.* `[CONFIRMAR CON EL DUEÑO]`

### Cómo pagar en el local
En Alameda 2432 tenemos dos locales:
- **Hogar, electro y juguetes:** efectivo y débito.
- **Disney y licencias:** débito, crédito, transferencia y efectivo.

Boleta o factura, como necesites. `[CONFIRMAR CON EL DUEÑO]`

> nota: sale de tr. TT-7707464302034160914 (nov 2026).

---

## 7. Canal oficial

## ESTE ES NUESTRO ÚNICO CANAL OFICIAL

**Pedidos solo a estos números:**
- **Santiago y retiro en tienda:** +56 9 2114 5368
- **Regiones:** `{WHATSAPP_REGIONES}` `[CONFIRMAR CON EL DUEÑO]`

> nota: si no hay número de regiones o es el mismo, se muestra una sola línea: **Pedidos solo a este número: +56 9 2114 5368**.

**En redes nos encuentras como Mayte Emporio:**
- TikTok **@mayte_emporio**
- Instagram **@mayte_emporioo** (cuenta verificada)

**Para comprar seguro:**
1. Pide solo a nuestros números oficiales. Ningún otro.
2. **No publicamos datos bancarios en la web.** Te los damos solo por el chat oficial, junto con tu número de pedido.
3. Nunca te vamos a pedir claves ni códigos que te lleguen por SMS.
4. Antes de transferir, revisa que el chat sea con uno de nuestros números.
5. ¿Te escribieron a nombre nuestro desde otra cuenta o número? No transfieras y avísanos al número oficial.

> nota: "pedidos solo a este número" es de la bio de IG (`adn-marca.md` §H). "Nunca te escribimos primero por interno" (CLAUDE.md) **no se usa** hasta que el dueño lo confirme. Con el panel, el equipo **sí** escribe al cliente, pero solo dentro de un pedido que el cliente inició. No se nombran ni se enlazan las cuentas no confirmadas.

**Versión corta (footer y página de seguimiento):**
`Pedidos solo a nuestros números oficiales. Los datos para transferir te los damos solo por el chat oficial.`

---

## 8. Sigue tu pedido

## ¿YA PEDISTE? SIGUE TU PEDIDO
Cuando hagas tu pedido por WhatsApp, te mandamos un link para ver en qué va: si ya lo pagaste, si está listo o el número de seguimiento de la agencia. Sin dar tus datos.

> nota: no hay campo para buscar por número de pedido (ver `arquitectura-landing.md` §A7).

---

## 9. Quiénes somos (opcional)

*script:* **Familia Emporio**
## DEL TOLDO AMARILLO AL CONTENEDOR
Partimos con un toldo amarillo. Hoy abrimos contenedores en plena calle y te dejamos lo nuevo a precio de remate.
Botón: `▶ Mira cómo abrimos un contenedor`

> nota: el "toldo amarillo" sale de sus hashtags (`adn-marca.md` §A). Va sin nombres ni cifras. **[CONFIRMAR CON EL DUEÑO que quiere contarlo]**.

---

## 10. Preguntas frecuentes

## PREGUNTAS QUE NOS HACEN SIEMPRE 💬

> nota: van en orden de frecuencia, según `faq-clientes.md`.

**¿Dónde están? ¿Qué metro queda cerca?**
En Alameda 2432, Santiago Centro, a 5 minutos del Metro Unión Latinoamericana y a 7 del Metro República.

**¿Hacen envíos a regiones?**
Sí. Toca "Pedidos a regiones", te cotizamos por WhatsApp y, cuando confirmamos tu pago, lo despachamos por agencia con número de seguimiento. `[CONFIRMAR CON EL DUEÑO: agencias, quién paga y mínimo]`

**¿A qué número les escribo?**
Si eres de Santiago o vas a retirar: +56 9 2114 5368. Si eres de región: `{WHATSAPP_REGIONES}` `[CONFIRMAR CON EL DUEÑO]`. Los dos botones de esta página te llevan al correcto.

**¿Les queda…?**
Pregúntanos por WhatsApp. Si lo viste en un video, toca "Pedir este por WhatsApp" y nos llega el producto ya escrito.

**¿Cuántos vienen por embalaje? ¿Desde cuántas es precio por mayor?**
En los productos de "Lo nuevo de los videos" está el precio por mayor cuando lo dijimos en el video. Para lo demás, pregúntanos y te cotizamos.

**¿Hasta qué hora atienden? ¿Abren domingo y feriados?**
En temporada: lunes a viernes de 9:00 a 18:00, y sábado, domingo y feriados de 9:00 a 16:00. `[CONFIRMAR CON EL DUEÑO]`

**¿Es seguro transferir?**
Transfiere solo a los datos que te damos en el chat de nuestros números oficiales, y después de recibir tu cotización con número de pedido. Si alguien te escribe desde otra cuenta a nombre nuestro, no transfieras.

**¿Cómo sé en qué va mi pedido?**
Te mandamos un link de seguimiento por WhatsApp. Cuando lo despachamos, te llega también la agencia y el número de seguimiento.

**¿Mi caja llegó con un problema, a quién le digo?**
Escríbenos por el mismo chat donde hiciste el pedido, con tu número de pedido. Así lo vemos al tiro.

**¿Qué medios de pago aceptan en el local?**
En hogar, electro y juguetes: efectivo y débito. En Disney y licencias: débito, crédito, transferencia y efectivo. `[CONFIRMAR CON EL DUEÑO]`

**¿Cómo funciona "Yo quiero mi promo"?**
`[COMPLETAR CON EL DUEÑO: cómo se canjea la promo del comentario]`
> nota: no está explicado en los videos (`pendientes.md` §3). Si no se completa, la pregunta no se publica.

**¿Dan boleta o factura?**
Sí, boleta o factura, como necesites. `[CONFIRMAR CON EL DUEÑO]`

**¿Hacen cambios?**
`[COMPLETAR CON EL DUEÑO]`
> nota: no hay dato. Si no se completa, la pregunta no se publica.

**¿Tienen local en regiones?**
Sí, en Temuco: Manuel Rodríguez 1232. `[CONFIRMAR CON EL DUEÑO si se incluye]`
> nota: la fuente es la descripción de TT-7592989161653685522 (ene 2026). Si no se aprueba, la pregunta se saca.

---

## 11. Footer y barra fija

### Footer
- **Marca:** Donde Carlitos Emporio · *Familia Emporio*
- `📍 Alameda 2432, Santiago Centro`
- `🕘 L a V 9:00–18:00 · Sáb, dom y feriados 9:00–16:00 (temporada)`
- **Pedidos:** `Santiago y retiro: +56 9 2114 5368` · `Regiones: {WHATSAPP_REGIONES}`
- **Redes:** TikTok @mayte_emporio · Instagram @mayte_emporioo
- **Aviso corto:** la versión corta de §7
- **Cierre:** *Los esperamos en Alameda 2432 con la mejor variedad de barrio Meiggs.*
- **Solo en el demo:** `Demo de prueba. Precios sacados de nuestros videos de TikTok (nov 2026).`

> nota: el cierre es su despedida en los videos (tr. TT-7707649267464293639).

### Barra fija (móvil)
- `Santiago / retiro`
- `Regiones`
- aria-label: `Escribir por WhatsApp a pedidos de Santiago` / `Escribir por WhatsApp a pedidos de regiones`

---

## 12. Página de seguimiento (`/seguimiento/{token}`)

**Título:** `PEDIDO {numero}`

**Estado (grande) y explicación:**

| Estado API | Título | Explicación |
|---|---|---|
| `cotizado` | `Cotizado` | `Te mandamos la cotización por WhatsApp.` |
| `esperando_pago` | `Esperando tu pago` | `Cuando transfieras, mándanos el comprobante por el mismo chat.` |
| `pagado` | `Pago recibido` | `Ya confirmamos tu pago. ¡Gracias!` |
| `preparando` | `Preparando tu pedido` | `Estamos armando tu pedido.` |
| `despachado` | `En camino` | `Tu pedido salió por {agencia}.` |
| `listo_retiro` | `Listo para retiro` | `Te esperamos en {DIRECCION_RETIRO}.` |
| `entregado` | `Entregado` | `¡Gracias por comprar en Donde Carlitos Emporio!` |
| `cancelado` | `Cancelado` | `Este pedido se canceló. Si tienes dudas, escríbenos.` |

**Bloques:**
- **Línea de tiempo:** `{estado} · {dd/mm hh:mm}`
- **Si está despachado:** `Agencia: {agencia}` · `N.º de seguimiento: {nro_seguimiento}` · botón `Copiar número`, que cambia a `¡Copiado!`. Debajo: `Con este número lo sigues en {agencia}.`
- **Si está listo para retiro:** `📍 {DIRECCION_RETIRO}` · `🕘 {HORARIO_RETIRO}` · `Muestra este número en el local: {numero}`
- **Productos:** `{cantidad} x {descripcion}` · `Total: {total}`
- **Botón:** `Escríbenos por WhatsApp`
- **Pie:** la versión corta del canal oficial.
- **Solo en el demo:** `Pedido de prueba (demo).`

**Error:**
## NO ENCONTRAMOS ESTE PEDIDO
Revisa que el link esté completo o escríbenos por WhatsApp con tu número de pedido.
Botones: `Santiago / retiro` · `Regiones`

---

# P. Panel privado (microcopy)

## P1. Login
- **Título:** `PEDIDOS DEL WHATSAPP`
- **Campos y botón:** `Usuario` · `Contraseña` · `Entrar`
- **Errores:**
  - `Usuario o contraseña incorrectos.`
  - `Demasiados intentos. Espera unos minutos y vuelve a probar.`
- **Menú:** `Hoy` · `Nuevo pedido` · `Pedidos` · `Buscar` · `Resumen` · `Vitrina` · `Ajustes` (solo admin) · `Cerrar sesión`
- **Saludo:** `Hola, {usuario} · Canal: {Santiago | Regiones}`

## P2. Hoy
- **Título:** `HOY`
- **Contador:** `{n} nuevos · {n} pagados · {n} despachados · {n} entregados`
- **Filtro:** `Mi canal` / `Todos`

**Bloques:**
1. `ESPERANDO PAGO HACE MÁS DE 24 H ({n})`
   - tarjeta: `{numero} · {nombre} · {total} · hace {n} h`
   - botones: `Recordar pago` · `Cancelar`
2. `PAGADOS SIN DESPACHAR ({n})`
   - tarjeta: `{numero} · {nombre} · {comuna} · {agencia} · pagado hace {n} h`
   - botón: `Marcar despachado`
3. `LISTOS PARA RETIRO ({n})`
   - tarjeta: `{numero} · {nombre} · listo hace {n} días`
   - botón: `Marcar entregado`
4. `COTIZADOS HOY ({n})`

**Vacío:** `Nada pendiente por aquí. ¡Bien!`
**Botón fijo:** `+ Nuevo pedido`

## P3. Nuevo pedido
- **Título:** `NUEVO PEDIDO`
- **Campos:**
  - **Canal:** `Santiago` / `Regiones`
  - **Teléfono:** `Teléfono del cliente`, placeholder `Pega el número del chat`
  - **Si el cliente ya existe:** `Cliente conocido: {nombre} · {n} pedidos · última comuna: {comuna}` · botón `Usar estos datos`
  - **Nombre:** `Nombre`, placeholder `Ej: Pedro Soto`
  - **Productos:** `Productos`, buscador `Busca por nombre o código`
    - botón `+ Línea libre`
    - campos de la línea: `Cantidad` · `Descripción` · `Precio c/u (opcional)` · `Quitar`
  - **Descuento:** `Descuento (opcional)`
  - **Entrega:** `Despacho` / `Retiro en tienda`
  - **Si es despacho:**
    - `Comuna` (autocompletado), placeholder `Ej: Chillán`;
    - `Agencia`: botones `{agencias configuradas}` + `Otra`;
    - `Costo de envío (opcional)`;
    - casilla `Envío por pagar`.
  - **Pago:** `Transferencia` / `En tienda`
  - **Total:** `Total` · ayuda `Se calcula solo. Si una línea va sin precio, escríbelo tú.`
  - **Nota:** `Nota interna (no la ve el cliente)`
- **Botón:** `Guardar como cotizado`
- **Al guardar (toast):** `Pedido {numero} guardado.` · `Enviar cotización`

**Errores:**
- `Revisa el teléfono: tiene que ser un celular chileno (9 dígitos).`
- `Agrega al menos un producto.`
- `La cantidad tiene que ser 1 o más.`
- `El total no calza con los productos. Revísalo o deja una línea sin precio.`
- `Elige la comuna para el despacho.` (solo al marcar "Despachado")

## P4. Ficha del pedido
- **Encabezado:** `{numero}` · `{estado}` · `{Santiago | Regiones}` · `Registrado por {usuario} · {dd/mm hh:mm}`
- **Link:** `Link de seguimiento` · botón `Copiar link`

**Estados (panel):**

| Estado API | Etiqueta | Botón de la acción siguiente |
|---|---|---|
| `cotizado` | Cotizado | `Marcar esperando pago` (transferencia) / `Marcar preparando` (en tienda) |
| `esperando_pago` | Esperando pago | `Marcar pagado` |
| `pagado` | Pagado | `Marcar preparando` |
| `preparando` | Preparando | `Marcar despachado` / `Marcar listo para retiro` |
| `despachado` | Despachado | `Marcar entregado` |
| `listo_retiro` | Listo para retiro | `Marcar entregado` |
| `entregado` | Entregado | — |
| `cancelado` | Cancelado | — |

**Confirmaciones:**
- **Pagado:** `¿Revisaste que la plata ya está en la cuenta? No te fíes solo del comprobante.` · `Sí, está pagado` · `Todavía no`
- **Despachado:** el formulario pide `Agencia` · `N.º de seguimiento` · `Comuna` · `Guardar y marcar despachado`
- **Cancelar:** `¿Cancelar el pedido {numero}?` · `Motivo (opcional)`, placeholder `Ej: se agotó, el cliente desistió` · `Sí, cancelar` · `No`

**Mensajes:**
- Título: `MENSAJES PARA EL CLIENTE`
- Nombres de los botones:
  - `Enviar cotización`
  - `Enviar datos de transferencia` (solo si están configurados)
  - `Avisar pago recibido`
  - `Avisar despacho`
  - `Avisar listo para retiro`
  - `Recordar pago`
  - `Avisar cancelación`
  - `Abrir chat`
- Bajo cada uno: `Envíalo desde el WhatsApp de {Santiago | Regiones}.`
- Vista previa: `Así le llega:` + el texto
- Respaldo: `Copiar mensaje`, que cambia a `¡Copiado!`
- Si falta `DATOS_TRANSFERENCIA` (solo admin): `Aún no hay datos de transferencia. Agrégalos en Ajustes para que salgan en la cotización.`

**Historial:**
- Título: `Historial`
- Línea: `{dd/mm hh:mm} · {estado anterior} → {estado nuevo} · {usuario}`
- Solo en el demo: `Pedido de prueba (demo). Los mensajes se abren sin número, para que elijas el contacto.`

## P5. Resumen y exportar
- **Título:** `RESUMEN`
- **Selector:** `Hoy` · `Esta semana` · `Elegir fechas`
- **Tarjetas:**
  - `Pedidos` · `{n}` (`Santiago {n} · Regiones {n}`)
  - `Vendido (pagados y entregados)` · `{monto}`
  - `Ticket promedio` · `{monto}`
  - `Cancelados` · `{n}`
- **Listas:**
  - `Comunas que más piden`
  - `Agencias más usadas`
  - `Por persona: {usuario} · {n} registrados · {n} despachados`
- **Botón:** `Exportar a Excel (CSV)`
- **Aviso al exportar:** `Este archivo trae nombres y teléfonos de clientes. Guárdalo bien y no lo compartas.`

**Columnas del CSV, en orden** (UTF-8 con BOM, separador `;`, fechas `dd-mm-aaaa hh:mm`, montos como entero sin `$` ni puntos, para que Excel los sume):
`numero;fecha_registro;canal;entrega;estado;nombre;telefono;comuna;agencia;nro_seguimiento;productos;descuento;costo_envio;total;pago;registrado_por;fecha_pagado;fecha_despachado;despachado_por;fecha_entregado;motivo_cancelacion`
- `productos`: `2 x Squishy viral… | 1 x Tren navideño…`

> nota para ingeniero-go:
> - Exportar solo con rol `admin` *(sugerencia)*.
> - Escapar las celdas que empiecen con `=`, `+`, `-`, `@`, tabulación o retorno de carro, anteponiendo `'`. Así se evita la inyección de fórmulas en Excel: el nombre y la descripción son texto libre.

## P6. Buscar
- **Campo:** `Teléfono, n.º de pedido o nombre`
- **Sin resultados:** `No encontramos pedidos con eso.`
- **Resultado por teléfono:** `{nombre} · {n} pedidos · último: {dd/mm}`

## P7. Vitrina y catálogo
- **Título:** `VITRINA Y PRODUCTOS`
- **Campos:**
  - `Nombre`
  - `Código`
  - `Precio detalle`
  - `Precio anterior (solo si es real)`
  - `Oferta válida hasta`
  - `Precio por mayor (opcional)`
  - `Unidad (ej: caja de 12)`
  - `Fecha del precio`
  - `Foto`
  - `Link al video de TikTok`
  - `Mostrar en la vitrina`
  - `Orden en la vitrina`
- **Botón:** `Guardar`, que cambia a `Guardado`
- **Errores:**
  - `La vitrina tiene un máximo de 8 productos. Saca uno antes de agregar otro.`
  - `El precio anterior tiene que ser mayor que el precio actual.`
  - `El precio por mayor tiene que ser menor que el precio detalle.`

## P8. Ajustes (solo admin)
- **Usuarios:** `Usuarios` · `+ Agregar usuario` · `Nombre` · `Usuario` · `Canal por defecto` · `Rol: Atención / Administrador` · `Desactivar`
- **Datos de transferencia:** `Datos de transferencia`
  - ayuda: `Banco, tipo de cuenta, número, RUT, nombre y correo. Se mandan solo por WhatsApp dentro de la cotización. Nunca se publican en la web.`
  - vacío: `Sin datos. La cotización dirá que se los mandas aparte.`
- **Retiro:** `Dirección de retiro` · `Horario de retiro`
- **Agencias:** `Agencias` · `+ Agregar agencia`
- **Números de WhatsApp (en el demo, solo lectura):**
  - `Santiago: +56 9 2114 5368`
  - `Regiones: {número o "sin número: se usa el de Santiago"}`
- **Botón:** `Guardar cambios`

## P9. Errores generales del panel
- `No pudimos guardar. Revisa tu conexión y vuelve a intentar.`
- `Alguien más cambió este pedido. Recarga para ver la última versión.`
- `Tu sesión se cerró. Vuelve a entrar.`
- `Parece que no tienes internet. Lo que escribiste sigue aquí.`
