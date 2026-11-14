# Mensajes de WhatsApp: plantillas exactas

*estratega-propuesta · 2026-11-14 · enfoque "Tu WhatsApp sigue igual, pero cada pedido queda ordenado"*
*Reemplaza a `mensaje-pedido-whatsapp.md` (obsoleto: era del pedido desde el carrito, que ya no existe).*

**Quién lo usa:**
- **constructor-web:** plantillas (a), (b) y (c), que arma la landing.
- **ingeniero-go:** plantillas (d1)–(d8), que arma el panel. Se recomienda armarlas **en la API**, en un solo lugar y con tests byte a byte contra los ejemplos de este archivo, igual que se hizo con la plantilla anterior.
- **qa-revisor:** los ejemplos de §6 son la salida esperada.

**Marcas:**
- `[EJEMPLO]`: dato ilustrativo.
- `[CONFIRMAR CON EL DUEÑO]`: política del negocio que no está confirmada.
- Ninguna de las dos marcas aparece **dentro** de los mensajes. Van solo en las notas.

---

## 1. Cómo se abren

| Plantilla | Quién la dispara | Destino (`wa.me`) |
|---|---|---|
| (a) Botón "Pedidos Santiago / retiro en tienda" | Cliente, en la landing | `https://wa.me/{WHATSAPP_SANTIAGO}?text={encodeURIComponent(mensaje)}` |
| (b) Botón "Pedidos a regiones" | Cliente, en la landing | `https://wa.me/{WHATSAPP_REGIONES}?text=…` |
| (c) "Pedir este por WhatsApp" (vitrina) | Cliente. Primero elige Santiago o Regiones | El número del canal elegido |
| (d1)–(d8) Mensajes del panel | Quien atiende, desde la ficha del pedido | `https://wa.me/{telefono_cliente}?text=…` (teléfono normalizado, solo dígitos: `569XXXXXXXX`) |

### Configuración del negocio (no va fija en el código)

| Variable | Valor | Dónde vive | Nota |
|---|---|---|---|
| `WHATSAPP_SANTIAGO` | `56921145368` | web (`NEXT_PUBLIC_…`; es un número público) | Confirmado en las bios oficiales (`research/pendientes.md` §3) |
| `WHATSAPP_REGIONES` | **vacío** `[CONFIRMAR CON EL DUEÑO]` | web (`NEXT_PUBLIC_…`) | **Si está vacío, el botón de regiones usa `WHATSAPP_SANTIAGO`** y el mensaje igual dice "regiones", para que quien atiende lo reconozca. Ojo: según `adn-marca.md` §G, en nov 2026 el mismo +56 9 2114 5368 se presentó como "número para región". Puede que sea **un solo número** con dos personas atendiendo. Se pregunta en la reunión |
| `DATOS_TRANSFERENCIA` | **vacío por defecto** | API (solo panel; **no** va en `NEXT_PUBLIC_…` ni en la landing) | Texto de varias líneas que escribe el dueño. **Nunca se inventa.** Si está vacío, los mensajes lo omiten y el botón (d2) no aparece |
| `DIRECCION_RETIRO` | `Alameda 2432, Santiago Centro` | API | Si hay dos locales en la misma dirección, el dueño puede precisar cuál `[CONFIRMAR CON EL DUEÑO]` |
| `HORARIO_RETIRO` | `L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00` | API | Horario de nov 2026 (tr. TT-7707471893363789063) `[CONFIRMAR CON EL DUEÑO]`. Si está vacío, se omite la línea 🕘 |
| `URL_SITIO` | `http://localhost:{WEB_PORT}` en el demo; el dominio en producción | API | Base del link de seguimiento |

### Desde qué WhatsApp se envían los mensajes del panel
`wa.me` abre el WhatsApp que esté instalado en **el equipo de quien aprieta el botón**. Por eso:
- Al lado de cada botón de mensaje, el panel muestra: **"Envíalo desde el WhatsApp de {Santiago | Regiones}"**, según el canal del pedido.
- Si en un celular hay dos WhatsApp (personal y Business), Android pregunta con cuál abrir. Hay que elegir el del negocio. Esto va en la capacitación.

### Regla del demo (pedidos ficticios)
- **(a), (b) y (c)** en el demo: llevan la primera línea `*[PRUEBA · DEMO]*`, para que nadie del local confunda una prueba con un pedido real si alguien aprieta "enviar". El destino lo define el orquestador (hoy, el número del negocio; ver QA I-7). **En la presentación no se aprieta "enviar".**
- **(d1)–(d8)** con `es_demo = true`: el link es **`https://wa.me/?text=…`, sin número**. Así WhatsApp deja elegir el contacto y nunca se abre un chat con un desconocido que tenga, por azar, el teléfono ficticio del pedido.

---

## 2. Reglas de formato (iguales para todas)

1. **Saltos de línea:** `\n` (LF).
2. **Bloques:** cada mensaje es una lista de bloques. Entre bloques va **exactamente una línea en blanco**. Si una línea es condicional y no aplica, **no se escribe**. Si un bloque queda vacío, **desaparece entero**, con su línea en blanco.
3. **Sin espacios al final** de ninguna línea y sin salto de línea al final del mensaje. **Única excepción:** (a), (b) y (c) terminan en `: ` (dos puntos y un espacio), para que el cliente siga escribiendo justo ahí.
4. **Negrita:** `*texto*`.
5. **Emojis permitidos** (todos de la lista de `estilo-carlitos` §5): 🛍️ (U+1F6CD U+FE0F), 📦 (U+1F4E6), 📍 (U+1F4CD), 🕘 (U+1F558). Como máximo 1 por línea. **Nunca** en líneas de pago, montos ni avisos de seguridad.
6. **Primera línea del demo:** `*[PRUEBA · DEMO]*`, solo si `es_demo = true`.
7. **Montos:** ver §3.
8. **Texto libre** (nombre, descripción, comuna, agencia, n.º de seguimiento, motivo, código):
   1. `\r`, `\n` y `\t` pasan a ser un espacio;
   2. se borran `*`, `_`, `~` y `` ` `` (rompen el formato de WhatsApp);
   3. se juntan los espacios repetidos;
   4. se quitan los espacios de los extremos.
9. **`{nombre}` es el primer nombre** del cliente: la primera palabra del nombre ya limpio ("Pedro Soto" → "Pedro"). Si está vacío, el saludo queda `Hola, …`.
10. **`{datos_transferencia}`:** se respetan sus saltos de línea. Se quita `\r`, se recorta cada línea y se borran las líneas vacías. No se le aplica la regla 8.
11. **Más de 25 líneas de productos:** se muestran 25 y luego `• … y {n} productos más` (`producto` si `n = 1`). El total va completo.

## 3. Montos: `$12.500`

- Pesos chilenos, con `$` pegado, **punto de miles**, sin decimales y sin espacio: `$500`, `$4.000`, `$12.500`, `$187.500`, `$1.250.000`.
- Descuento: con signo menos antes del `$`: `-$7.500`.
- **No usar `Intl.NumberFormat` ni `toLocaleString`**, porque según el motor sale `4000`, `4,000` o `$ 4.000`. Usar:
  ```ts
  const clp = (n: number) => '$' + Math.round(n).toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
  ```
  En Go: separar los miles a mano con un punto (sin `golang.org/x/text`), con test para 0, 500, 4000, 12500, 1250000.

## 4. Variables

### Del pedido (panel)
| Variable | Tipo | Nota |
|---|---|---|
| `numero` | `DC-0101` | Correlativo |
| `canal` | `santiago` \| `regiones` | Qué WhatsApp lo atendió |
| `entrega` | `retiro` \| `despacho` | Por defecto: Santiago → `retiro`, Regiones → `despacho`. Se puede cambiar (ej.: alguien de región que pasa a buscar) |
| `nombre` | texto | Se usa el primer nombre (regla 9) |
| `lineas[]` | `{cantidad, descripcion, precio_unitario?}` | Producto del catálogo o **línea libre**. Si no tiene precio, la línea va sin precio |
| `descuento` | entero ≥ 0, opcional | Se muestra solo si es > 0 |
| `costo_envio` | entero ≥ 0, opcional | Si es > 0, se suma al total y se muestra |
| `envio_por_pagar` | bool | Solo si `entrega = despacho` y sin `costo_envio`. Política `[CONFIRMAR CON EL DUEÑO]` |
| `total` | entero | Si todas las líneas tienen precio: `Σ(cantidad × precio) − descuento + costo_envio`, y la API lo valida. Si alguna línea va sin precio, lo escribe quien atiende |
| `comuna` | texto | Opcional en la cotización; obligatorio para marcar "Despachado" |
| `agencia` | texto | Lista configurable `[EJEMPLO: Starken, Varmontt, Pullman Cargo…]` + "Otra". Obligatoria para "Despachado" |
| `nro_seguimiento` | texto | El número de la agencia. Obligatorio para "Despachado" |
| `pago` | `transferencia` \| `en_tienda` | Por defecto: Regiones → `transferencia`, Santiago → `en_tienda` |
| `pagado` | bool | `true` si el pedido pasó por el estado "Pagado" |
| `motivo` | texto, opcional | Motivo de la cancelación |
| `link_seguimiento` | URL | `{URL_SITIO}/seguimiento/{token}`; ver §5 |
| `es_demo` | bool | |

### De la vitrina (landing)
| Variable | Nota |
|---|---|
| `producto` | Nombre del producto |
| `precio` | Precio al detalle publicado |
| `precio_mayor`, `unidad_mayor` | Opcionales |
| `fecha_precio` | `dd/mm` de la última actualización: la fecha del video o del cambio en el panel |
| `codigo` | Código interno (los de los ejemplos son `[EJEMPLO]`) |

## 5. Link de seguimiento

- **Forma:** `{URL_SITIO}/seguimiento/{token}`.
- **El token** es aleatorio, de 16 bytes en base64url (22 caracteres). **No es el `DC-0101`**: un número correlativo se puede adivinar, y entonces cualquiera podría recorrer los pedidos de otros.
- **Qué muestra la página:** número de pedido, estado actual con su historial (solo estado y fecha), productos y cantidades, total, y, si está despachado, agencia + n.º de seguimiento.
- **Qué NO muestra:** nombre, teléfono, comuna, notas internas ni quién atendió.
- **Cabeceras:** `noindex` y `Referrer-Policy: no-referrer`.
- **Vencimiento:** `[COMPLETAR LUCKAS: p. ej. 90 días después de entregado]`.
- **Cuándo va en los mensajes:** la línea solo aparece si el pedido ya tiene token. Se recomienda generarlo al crear el pedido, y así sale siempre.

---

## 6. Plantillas y ejemplos renderizados

> Los nombres, teléfonos y pedidos son **ficticios**. Los precios de producto salen de `research/productos-semilla.json` (videos de nov 2026). Los códigos de producto, agencias, costos de envío, descuentos, n.º de seguimiento y tokens son `[EJEMPLO]`.
> El dominio de los ejemplos es `tienda.example` (dominio reservado; no existe).
> Los pedidos DC-0101 a DC-0109 son los mismos de la semilla sugerida del demo (`pitch.md` §c). DC-0120 existe solo como ejemplo.

### (a) Botón "Pedidos Santiago / retiro en tienda"

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola Donde Carlitos Emporio 🛍️
Vengo de la web. Quiero hacer un pedido para *retiro en tienda (Santiago)*.
Lo que busco es: 
```

**Ejemplo a1 (producción)**
```
Hola Donde Carlitos Emporio 🛍️
Vengo de la web. Quiero hacer un pedido para *retiro en tienda (Santiago)*.
Lo que busco es: 
```
URL resultante:
```
https://wa.me/56921145368?text=Hola%20Donde%20Carlitos%20Emporio%20%F0%9F%9B%8D%EF%B8%8F%0AVengo%20de%20la%20web.%20Quiero%20hacer%20un%20pedido%20para%20*retiro%20en%20tienda%20(Santiago)*.%0ALo%20que%20busco%20es%3A%20
```

**Ejemplo a2 (demo)**
```
*[PRUEBA · DEMO]*
Hola Donde Carlitos Emporio 🛍️
Vengo de la web. Quiero hacer un pedido para *retiro en tienda (Santiago)*.
Lo que busco es: 
```

> nota: no se ofrece despacho a domicilio en Santiago. Un cliente lo preguntó ("realiza envío a domicilio en santiago", TT-7706669405505129735), pero no hay dato `[CONFIRMAR CON EL DUEÑO]`.

### (b) Botón "Pedidos a regiones"

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola Donde Carlitos Emporio 📦
Vengo de la web. Soy de *regiones* y quiero cotizar un pedido con envío.
Mi comuna y lo que busco: 
```

**Ejemplo b1 (producción, con `WHATSAPP_REGIONES` vacío → va a `56921145368`)**
```
Hola Donde Carlitos Emporio 📦
Vengo de la web. Soy de *regiones* y quiero cotizar un pedido con envío.
Mi comuna y lo que busco: 
```

**Ejemplo b2 (demo)**
```
*[PRUEBA · DEMO]*
Hola Donde Carlitos Emporio 📦
Vengo de la web. Soy de *regiones* y quiero cotizar un pedido con envío.
Mi comuna y lo que busco: 
```

> nota: se pide la comuna desde el primer mensaje porque el envío depende de ella y la agencia también. Se deja en una sola línea abierta, porque en el celular el cursor queda al final: "Chillán, 12 squishy y 4 dinosaurios" se escribe de corrido.

### (c) "Pedir este por WhatsApp" (vitrina "Lo nuevo de los videos")

Al tocar el botón se abre una hoja con dos opciones: **"Santiago / retiro en tienda"** y **"Regiones (envío)"**. La elección se recuerda durante la visita.

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola Donde Carlitos Emporio 🛍️
Vengo de la web y me interesa este producto:

*{producto}*
Precio publicado el {fecha_precio}: {precio}
[[si precio_mayor]]Por mayor: {precio_mayor} c/u[[si unidad_mayor]] ({unidad_mayor})
Ref. {codigo}

[[santiago]]Es para *retiro en tienda (Santiago)*.
[[santiago]]¿Les queda? Quiero esta cantidad: 
[[regiones]]Es con *envío a regiones*.
[[regiones]]¿Les queda? Mi comuna y la cantidad: 
```

**Ejemplo c1 (Santiago, producción)**: squishy, TT-7707649267464293639.
```
Hola Donde Carlitos Emporio 🛍️
Vengo de la web y me interesa este producto:

*Squishy viral que cambia de color con el sol*
Precio publicado el 13/11: $4.000
Por mayor: $3.500 c/u (caja de 12)
Ref. P001

Es para *retiro en tienda (Santiago)*.
¿Les queda? Quiero esta cantidad: 
```

**Ejemplo c2 (Regiones, demo)**: tren grande, TT-7706656443843284232.
```
*[PRUEBA · DEMO]*
Hola Donde Carlitos Emporio 🛍️
Vengo de la web y me interesa este producto:

*Tren navideño para árbol (grande)*
Precio publicado el 10/11: $17.000
Por mayor: $15.000 c/u (embalaje)
Ref. P012

Es con *envío a regiones*.
¿Les queda? Mi comuna y la cantidad: 
```

> nota: "Precio publicado el dd/mm" va a propósito. Hay quejas de que el precio del video no se respetó en caja (`faq-clientes.md`, quejas), así que el mensaje deja claro **de cuándo** es ese precio, sin prometer que sigue vigente.

### (d1) Cotización / resumen del pedido

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *{numero}*:

[[por cada línea con precio]]• {cantidad} x {descripcion} · {precio_unitario} c/u = {subtotal}
[[por cada línea sin precio]]• {cantidad} x {descripcion}
[[si descuento > 0]]Descuento: -{descuento}
[[si costo_envio > 0]]Envío[[si agencia]] por {agencia}: {costo_envio}

*Total: {total}*
[[si despacho y envio_por_pagar y sin costo_envio]]El envío lo pagas tú en la agencia, al recibirlo (por pagar).

[[si despacho]]📦 Despacho[[si comuna]] a {comuna}[[si agencia]] por {agencia}.[[si no hay agencia]] ¿Por qué agencia te lo mandamos?
[[si retiro]]📍 Retiro en {DIRECCION_RETIRO}.

[[si pago = transferencia y hay datos]]*Para pagar, transfiere a:*
[[si pago = transferencia y hay datos]]{datos_transferencia}

[[si pago = en_tienda]]Lo pagas al retirar en el local.
[[si pago = transferencia y hay datos]]Después mándanos el comprobante por este chat y te confirmamos.
[[si pago = transferencia y no hay datos]]Si está todo bien, respóndenos *OK* y te mandamos los datos para pagar.

[[si link]]Sigue tu pedido aquí: {link_seguimiento}

[[si pago = transferencia]]Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d1-A (Regiones, por pagar, sin `DATOS_TRANSFERENCIA`, demo)**
```
*[PRUEBA · DEMO]*
Hola Pedro, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *DC-0101*:

• 12 x Squishy viral que cambia de color con el sol · $3.500 c/u = $42.000
• 4 x Dinosaurio gigante con ruedas · $14.000 c/u = $56.000

*Total: $98.000*
El envío lo pagas tú en la agencia, al recibirlo (por pagar).

📦 Despacho a Chillán por Starken.

Si está todo bien, respóndenos *OK* y te mandamos los datos para pagar.

Sigue tu pedido aquí: https://tienda.example/seguimiento/k7Qm2xVb9TfL4nRpA1c3Dw

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d1-B (Regiones, con descuento, costo de envío y `DATOS_TRANSFERENCIA` configurado, sin link)**
`DATOS_TRANSFERENCIA` = `[DATOS QUE CONFIGURA EL DUEÑO EN EL PANEL]` (marcador; **no** son datos reales).
Total: 15 × 12.500 = 187.500 − 7.500 + 8.000 = **188.000**.
```
Hola Marcela, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *DC-0102*:

• 15 x Caja sorpresa squishy premium · $12.500 c/u = $187.500
Descuento: -$7.500
Envío por Varmontt: $8.000

*Total: $188.000*

📦 Despacho a Concepción por Varmontt.

*Para pagar, transfiere a:*
[DATOS QUE CONFIGURA EL DUEÑO EN EL PANEL]

Después mándanos el comprobante por este chat y te confirmamos.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d1-C (Santiago, retiro, pago en tienda, con una línea libre sin precio)**
Como hay una línea sin precio, el total lo escribió quien atiende.
```
Hola Camila, te escribimos de *Donde Carlitos Emporio* 🛍️
Este es el resumen de tu pedido *DC-0104*:

• 3 x Pulpo viral (mueve ojos y tentáculos) · $4.000 c/u = $12.000
• 1 x Telón de Navidad 3 m x 2,5 m · $5.000 c/u = $5.000
• 1 x Bolsas de regalo surtidas

*Total: $20.000*

📍 Retiro en Alameda 2432, Santiago Centro.

Lo pagas al retirar en el local.
```

> notas:
> - Que el pago desde regiones sea por transferencia es *(inferencia)*: es el medio más aceptado en Chile [26] y no hay link de pago. `[CONFIRMAR CON EL DUEÑO]`.
> - "Por pagar" y los costos de envío: `[CONFIRMAR CON EL DUEÑO]`. Si el dueño cobra el envío, se usa `costo_envio`; si es por pagar, `envio_por_pagar`.
> - Enviar el mensaje **no** cambia el estado solo (la web no sabe si de verdad se envió). Al lado del botón va "Marcar como Esperando pago".

### (d2) Datos de transferencia (mensaje aparte)

Solo aparece si `DATOS_TRANSFERENCIA` tiene contenido. Sirve cuando el cliente responde "OK" o pide los datos de nuevo.

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, estos son los datos para pagar tu pedido *{numero}* por *{total}*:

{datos_transferencia}

Cuando transfieras, mándanos el comprobante por este chat.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d2-A (demo, datos en una línea)**
```
*[PRUEBA · DEMO]*
Hola Pedro, estos son los datos para pagar tu pedido *DC-0101* por *$98.000*:

[DATOS QUE CONFIGURA EL DUEÑO EN EL PANEL]

Cuando transfieras, mándanos el comprobante por este chat.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d2-B (producción, datos en varias líneas)**
`DATOS_TRANSFERENCIA` = `"  [Banco: lo escribe el dueño]\r\n\r\n[Cuenta: lo escribe el dueño]\n[RUT y nombre: lo escribe el dueño]  "`
```
Hola Luis, estos son los datos para pagar tu pedido *DC-0106* por *$20.000*:

[Banco: lo escribe el dueño]
[Cuenta: lo escribe el dueño]
[RUT y nombre: lo escribe el dueño]

Cuando transfieras, mándanos el comprobante por este chat.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

### (d3) Pago recibido

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, ¡recibimos tu pago! Muchas gracias.
Pedido *{numero}* · {total}

[[si despacho]]Ahora lo preparamos y te avisamos por aquí cuando salga, con el número de seguimiento de la agencia.
[[si retiro]]Ahora lo preparamos y te avisamos por aquí cuando esté listo para retiro.

[[si link]]Sigue tu pedido aquí: {link_seguimiento}
```

**Ejemplo d3-A (despacho, con link)**
```
Hola Marcela, ¡recibimos tu pago! Muchas gracias.
Pedido *DC-0102* · $188.000

Ahora lo preparamos y te avisamos por aquí cuando salga, con el número de seguimiento de la agencia.

Sigue tu pedido aquí: https://tienda.example/seguimiento/R2pX8vLq0WcN5sTe7YbH1g
```

**Ejemplo d3-B (retiro, sin link, demo)**
```
*[PRUEBA · DEMO]*
Hola Andrea, ¡recibimos tu pago! Muchas gracias.
Pedido *DC-0107* · $12.000

Ahora lo preparamos y te avisamos por aquí cuando esté listo para retiro.
```

> nota: antes de marcar "Pagado", quien atiende revisa que la plata **esté en la cuenta**, no solo el comprobante. En el barrio se han visto comprobantes falsos ([10] de `contexto-meiggs.md`). El panel lo recuerda al marcar "Pagado" (ver `copy-landing.md` §P).

### (d4) Pedido despachado

Requiere `agencia`, `nro_seguimiento` y `comuna` (el panel no deja marcar "Despachado" sin ellos).

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, ¡tu pedido *{numero}* ya va en camino! 📦

*Agencia:* {agencia}
*N.º de seguimiento:* {nro_seguimiento}
*Destino:* {comuna}

Con ese número lo puedes seguir en {agencia}.
[[si link]]Y aquí ves el estado de tu pedido: {link_seguimiento}

Si llega algún problema con tu caja, escríbenos por este mismo chat.
```

**Ejemplo d4-A (con link)**
```
Hola Jorge, ¡tu pedido *DC-0103* ya va en camino! 📦

*Agencia:* Starken
*N.º de seguimiento:* 000123456789
*Destino:* La Serena

Con ese número lo puedes seguir en Starken.
Y aquí ves el estado de tu pedido: https://tienda.example/seguimiento/Zt4uB9kQ2mVx7LcP0aRf3w

Si llega algún problema con tu caja, escríbenos por este mismo chat.
```

**Ejemplo d4-B (sin link, demo)**
```
*[PRUEBA · DEMO]*
Hola Marcela, ¡tu pedido *DC-0102* ya va en camino! 📦

*Agencia:* Varmontt
*N.º de seguimiento:* VT-998877
*Destino:* Concepción

Con ese número lo puedes seguir en Varmontt.

Si llega algún problema con tu caja, escríbenos por este mismo chat.
```

> nota: la última línea responde a un caso real: un cliente pidió en comentarios públicos que "conteste la niña de los pedidos de regiones" por un problema con su caja (TT-7706913997378948370). No se ponen links a las páginas de las agencias: no se inventan URLs.

### (d5) Listo para retiro

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, ¡tu pedido *{numero}* está listo para retiro! 🛍️

📍 {DIRECCION_RETIRO}
[[si HORARIO_RETIRO]]🕘 {HORARIO_RETIRO}

[[si pagado]]Ya está pagado.
[[si no pagado]]*A pagar al retirar:* {total}
Muestra tu número de pedido en el local: *{numero}*

[[si link]]Sigue tu pedido aquí: {link_seguimiento}
```

**Ejemplo d5-A (no pagado, sin link)**
```
Hola Camila, ¡tu pedido *DC-0104* está listo para retiro! 🛍️

📍 Alameda 2432, Santiago Centro
🕘 L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00

*A pagar al retirar:* $20.000
Muestra tu número de pedido en el local: *DC-0104*
```

**Ejemplo d5-B (pagado, con link, demo)**
```
*[PRUEBA · DEMO]*
Hola Andrea, ¡tu pedido *DC-0107* está listo para retiro! 🛍️

📍 Alameda 2432, Santiago Centro
🕘 L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00

Ya está pagado.
Muestra tu número de pedido en el local: *DC-0107*

Sigue tu pedido aquí: https://tienda.example/seguimiento/Hn6sD1qW8eYk3PzM5uJc0v
```

> nota: a quién se le muestra el número dentro del local, y cuántos días se guarda el pedido, se define con el dueño `[CONFIRMAR CON EL DUEÑO]`. Si fija un plazo, se agrega la línea `Te lo guardamos hasta el {dd/mm}.` después de "Muestra tu número…".

### (d6) Pedido cancelado

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, tu pedido *{numero}* quedó cancelado.
[[si motivo]]Motivo: {motivo}

[[si pagado]]Te escribimos por este chat para coordinar la devolución de tu pago.

Si quieres armar otro pedido, escríbenos por aquí cuando quieras.
```

**Ejemplo d6-A (con motivo, no pagado)**
```
Hola Rodrigo, tu pedido *DC-0108* quedó cancelado.
Motivo: Se nos agotó el producto.

Si quieres armar otro pedido, escríbenos por aquí cuando quieras.
```

**Ejemplo d6-B (pagado, sin motivo, demo)**
```
*[PRUEBA · DEMO]*
Hola Sofía, tu pedido *DC-0120* quedó cancelado.

Te escribimos por este chat para coordinar la devolución de tu pago.

Si quieres armar otro pedido, escríbenos por aquí cuando quieras.
```

> nota: cómo y en cuánto plazo se devuelve un pago es política del negocio `[CONFIRMAR CON EL DUEÑO]`. El mensaje solo compromete a escribir, no un plazo.

### (d7) Recordatorio de pago (extra, para la pantalla "Hoy")

Se ofrece en los pedidos con "Esperando pago" de más de 24 h.

**Plantilla**
```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, te escribimos de *Donde Carlitos Emporio* por tu pedido *{numero}* por *{total}*.

Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.
Si cambiaste de idea, avísanos y lo cancelamos, sin problema.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d7-A**
```
Hola Pedro, te escribimos de *Donde Carlitos Emporio* por tu pedido *DC-0101* por *$98.000*.

Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.
Si cambiaste de idea, avísanos y lo cancelamos, sin problema.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

**Ejemplo d7-B (demo, cliente sin nombre)**
```
*[PRUEBA · DEMO]*
Hola, te escribimos de *Donde Carlitos Emporio* por tu pedido *DC-0109* por *$34.000*.

Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.
Si cambiaste de idea, avísanos y lo cancelamos, sin problema.

Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.
```

### (d8) Abrir chat (sin plantilla de estado)

```
[[es_demo]]*[PRUEBA · DEMO]*
Hola {nombre}, te escribimos de *Donde Carlitos Emporio* por tu pedido *{numero}*.
```
Ejemplo: `Hola Jorge, te escribimos de *Donde Carlitos Emporio* por tu pedido *DC-0103*.`

---

## 7. Implementación de referencia (TypeScript)

Es la misma lógica que la API debe reproducir en Go. Lo que manda es el texto que produce, no el lenguaje.

```ts
type Canal = 'santiago' | 'regiones';
type Linea = { cantidad: number; descripcion: string; precio_unitario?: number | null };
type Pedido = {
  numero: string; canal: Canal; entrega: 'retiro' | 'despacho'; nombre?: string;
  lineas: Linea[]; descuento?: number | null; costo_envio?: number | null; envio_por_pagar?: boolean;
  total: number; comuna?: string; agencia?: string; nro_seguimiento?: string;
  pago: 'transferencia' | 'en_tienda'; pagado: boolean; motivo?: string;
  link_seguimiento?: string; es_demo: boolean;
};
type Config = { DATOS_TRANSFERENCIA?: string; DIRECCION_RETIRO: string; HORARIO_RETIRO?: string };
type Vitrina = { producto: string; precio: number; precio_mayor?: number | null;
  unidad_mayor?: string | null; fecha_precio: string; codigo: string };

const DEMO = '*[PRUEBA · DEMO]*';
const AVISO = 'Solo te escribimos desde este número. Si te piden transferir desde otro chat, no lo hagas.';
const MAX = 25;
const clp = (n: number) => '$' + Math.round(n).toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
const limpiar = (s = '') => s.replace(/[\r\n\t]+/g, ' ').replace(/[*_~`]/g, '').replace(/\s{2,}/g, ' ').trim();
const nombre1 = (s = '') => limpiar(s).split(' ')[0] ?? '';
const hola = (s?: string) => { const n = nombre1(s); return n ? `Hola ${n}` : 'Hola'; };
const datos = (s = '') => s.replace(/\r/g, '').split('\n').map(x => x.trim()).filter(Boolean);
type L = string | false | null | undefined;
const armar = (bloques: L[][]) => bloques
  .map(b => b.filter((l): l is string => typeof l === 'string' && l.length > 0))
  .filter(b => b.length > 0).map(b => b.join('\n')).join('\n\n');

// (a) y (b)
export const msgSantiago = (demo: boolean) => armar([[demo && DEMO,
  'Hola Donde Carlitos Emporio 🛍️',
  'Vengo de la web. Quiero hacer un pedido para *retiro en tienda (Santiago)*.',
  'Lo que busco es: ']]);
export const msgRegiones = (demo: boolean) => armar([[demo && DEMO,
  'Hola Donde Carlitos Emporio 📦',
  'Vengo de la web. Soy de *regiones* y quiero cotizar un pedido con envío.',
  'Mi comuna y lo que busco: ']]);

// (c)
export const msgVitrina = (p: Vitrina, canal: Canal, demo: boolean) => {
  const u = limpiar(p.unidad_mayor ?? '');
  return armar([
    [demo && DEMO, 'Hola Donde Carlitos Emporio 🛍️', 'Vengo de la web y me interesa este producto:'],
    [`*${limpiar(p.producto)}*`, `Precio publicado el ${p.fecha_precio}: ${clp(p.precio)}`,
      p.precio_mayor != null && `Por mayor: ${clp(p.precio_mayor)} c/u${u ? ` (${u})` : ''}`,
      `Ref. ${limpiar(p.codigo)}`],
    canal === 'santiago'
      ? ['Es para *retiro en tienda (Santiago)*.', '¿Les queda? Quiero esta cantidad: ']
      : ['Es con *envío a regiones*.', '¿Les queda? Mi comuna y la cantidad: '],
  ]);
};

// (d1)
export const msgCotizacion = (p: Pedido, c: Config) => {
  const ag = limpiar(p.agencia), co = limpiar(p.comuna);
  const d = p.pago === 'transferencia' ? datos(c.DATOS_TRANSFERENCIA) : [];
  const items: L[] = p.lineas.slice(0, MAX).map(l => {
    const base = `• ${l.cantidad} x ${limpiar(l.descripcion)}`;
    return l.precio_unitario != null
      ? `${base} · ${clp(l.precio_unitario)} c/u = ${clp(l.cantidad * l.precio_unitario)}` : base;
  });
  const resto = p.lineas.length - MAX;
  if (resto > 0) items.push(`• … y ${resto} producto${resto === 1 ? '' : 's'} más`);
  const envio = (p.costo_envio ?? 0) > 0;
  return armar([
    [p.es_demo && DEMO, `${hola(p.nombre)}, te escribimos de *Donde Carlitos Emporio* 🛍️`,
      `Este es el resumen de tu pedido *${p.numero}*:`],
    [...items, (p.descuento ?? 0) > 0 && `Descuento: -${clp(p.descuento!)}`,
      envio && `Envío${ag ? ` por ${ag}` : ''}: ${clp(p.costo_envio!)}`],
    [`*Total: ${clp(p.total)}*`,
      p.entrega === 'despacho' && p.envio_por_pagar && !envio && 'El envío lo pagas tú en la agencia, al recibirlo (por pagar).'],
    [p.entrega === 'despacho'
      ? `📦 Despacho${co ? ` a ${co}` : ''}${ag ? ` por ${ag}.` : '. ¿Por qué agencia te lo mandamos?'}`
      : `📍 Retiro en ${c.DIRECCION_RETIRO}.`],
    d.length > 0 ? ['*Para pagar, transfiere a:*', ...d] : [],
    [p.pago === 'en_tienda' ? 'Lo pagas al retirar en el local.'
      : d.length > 0 ? 'Después mándanos el comprobante por este chat y te confirmamos.'
      : 'Si está todo bien, respóndenos *OK* y te mandamos los datos para pagar.'],
    [p.link_seguimiento && `Sigue tu pedido aquí: ${p.link_seguimiento}`],
    [p.pago === 'transferencia' && AVISO],
  ]);
};

// (d2): devuelve null si no hay datos, y el panel oculta el botón
export const msgDatosTransferencia = (p: Pedido, c: Config) => {
  const d = datos(c.DATOS_TRANSFERENCIA);
  if (d.length === 0) return null;
  return armar([
    [p.es_demo && DEMO, `${hola(p.nombre)}, estos son los datos para pagar tu pedido *${p.numero}* por *${clp(p.total)}*:`],
    d, ['Cuando transfieras, mándanos el comprobante por este chat.'], [AVISO],
  ]);
};

// (d3)
export const msgPagoRecibido = (p: Pedido) => armar([
  [p.es_demo && DEMO, `${hola(p.nombre)}, ¡recibimos tu pago! Muchas gracias.`, `Pedido *${p.numero}* · ${clp(p.total)}`],
  [p.entrega === 'despacho'
    ? 'Ahora lo preparamos y te avisamos por aquí cuando salga, con el número de seguimiento de la agencia.'
    : 'Ahora lo preparamos y te avisamos por aquí cuando esté listo para retiro.'],
  [p.link_seguimiento && `Sigue tu pedido aquí: ${p.link_seguimiento}`],
]);

// (d4)
export const msgDespachado = (p: Pedido) => {
  const ag = limpiar(p.agencia);
  return armar([
    [p.es_demo && DEMO, `${hola(p.nombre)}, ¡tu pedido *${p.numero}* ya va en camino! 📦`],
    [`*Agencia:* ${ag}`, `*N.º de seguimiento:* ${limpiar(p.nro_seguimiento)}`, `*Destino:* ${limpiar(p.comuna)}`],
    [`Con ese número lo puedes seguir en ${ag}.`, p.link_seguimiento && `Y aquí ves el estado de tu pedido: ${p.link_seguimiento}`],
    ['Si llega algún problema con tu caja, escríbenos por este mismo chat.'],
  ]);
};

// (d5)
export const msgListoRetiro = (p: Pedido, c: Config) => armar([
  [p.es_demo && DEMO, `${hola(p.nombre)}, ¡tu pedido *${p.numero}* está listo para retiro! 🛍️`],
  [`📍 ${c.DIRECCION_RETIRO}`, c.HORARIO_RETIRO?.trim() && `🕘 ${c.HORARIO_RETIRO.trim()}`],
  [p.pagado ? 'Ya está pagado.' : `*A pagar al retirar:* ${clp(p.total)}`,
    `Muestra tu número de pedido en el local: *${p.numero}*`],
  [p.link_seguimiento && `Sigue tu pedido aquí: ${p.link_seguimiento}`],
]);

// (d6)
export const msgCancelado = (p: Pedido) => armar([
  [p.es_demo && DEMO, `${hola(p.nombre)}, tu pedido *${p.numero}* quedó cancelado.`, limpiar(p.motivo) && `Motivo: ${limpiar(p.motivo)}`],
  [p.pagado && 'Te escribimos por este chat para coordinar la devolución de tu pago.'],
  ['Si quieres armar otro pedido, escríbenos por aquí cuando quieras.'],
]);

// (d7)
export const msgRecordatorioPago = (p: Pedido) => armar([
  [p.es_demo && DEMO, `${hola(p.nombre)}, te escribimos de *Donde Carlitos Emporio* por tu pedido *${p.numero}* por *${clp(p.total)}*.`],
  ['Todavía no nos llega tu comprobante. Si ya transferiste, mándalo por este chat.',
    'Si cambiaste de idea, avísanos y lo cancelamos, sin problema.'],
  [AVISO],
]);

// (d8)
export const msgAbrirChat = (p: Pedido) => armar([[p.es_demo && DEMO,
  `${hola(p.nombre)}, te escribimos de *Donde Carlitos Emporio* por tu pedido *${p.numero}*.`]]);

// Link
export const linkWa = (texto: string, numero?: string) =>
  `https://wa.me/${numero ?? ''}?text=${encodeURIComponent(texto)}`;   // sin número = el usuario elige el contacto (demo)
```

**Teléfono del cliente (para `wa.me` y para buscar):** se normaliza a `569XXXXXXXX`. Se aceptan `+56 9 1234 5678`, `9 1234 5678`, `912345678` y `56912345678`. Cualquier otra cosa se rechaza con un mensaje claro.

## 8. Checklist de QA

- [ ] La salida es idéntica, byte a byte, a cada ejemplo de §6 con los mismos datos.
- [ ] Los montos van con punto de miles (`$12.500`), nunca `$12500`, `$12,500` ni `$ 12.500`.
- [ ] No hay dobles líneas en blanco. No hay espacios al final, salvo el `: ` final de (a), (b) y (c).
- [ ] Con `DATOS_TRANSFERENCIA` vacío: (d1) no muestra el bloque "Para pagar, transfiere a:" y (d2) no aparece en el panel.
- [ ] `DATOS_TRANSFERENCIA` no aparece en la landing, en el HTML público ni en ninguna variable `NEXT_PUBLIC_*`.
- [ ] Con `WHATSAPP_REGIONES` vacío, el botón de regiones abre `WHATSAPP_SANTIAGO` con el texto de (b).
- [ ] En pedidos `es_demo`, los botones del panel abren `https://wa.me/?text=…` (sin número).
- [ ] No se puede marcar "Despachado" sin agencia, n.º de seguimiento y comuna.
- [ ] El link de seguimiento usa el token, no el `DC-xxxx`, y la página no muestra nombre, teléfono ni comuna.
- [ ] Probado en Android e iOS, desde el navegador interno de TikTok e Instagram: saltos de línea, negritas y emojis llegan bien. Si algún cliente muestra `�`, se usa la variante sin emojis: se borra el emoji y el espacio que lo separa del texto.
