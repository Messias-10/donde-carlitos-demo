# Pitch — Donde Carlitos Emporio

*estratega-propuesta · 2026-11-14, ajustado el 2026-11-15 tras QA ronda 2 (B-1, M-6, M-7) · enfoque "Tu WhatsApp sigue igual, pero cada pedido queda ordenado"*

**A quién le habla:**
1. **El dueño:** decide y paga.
2. **Quien atiende los pedidos de regiones:** la que lo va a usar todos los días. En un comentario la llaman "la niña de los pedidos de regiones" (TT-7706913997378948370). **Nosotros no usamos esa expresión:** decimos "quien atiende regiones" o su nombre, si nos lo dicen.

**Tono:** cercano y directo, de tú. Si el dueño trata de usted, Luckas cambia a usted.

**Regla de oro:** solo se muestra lo que el demo hace. Todo lo demás va como "eso es la fase 2".

---

## (a) Mensaje de primer contacto (máximo 6 líneas)

**Canal:** DM de Instagram a @mayte_emporioo, o ir directo al local.
- **No** se escribe a los números de pedidos: ya están saturados (`faq-clientes.md`, quejas).
- Se adjunta un **video vertical de 60 s** del demo, en vez de un link. Un link de un desconocido a una marca que avisa por estafas puede parecer phishing, y además el demo no está en internet (CLAUDE.md, Autonomía).

```
Hola, ¿cómo están? Soy Luckas, desarrollador de Santiago, y sigo sus videos.
Vi que ahora atienden pedidos a regiones por WhatsApp, y en los comentarios preguntan harto "¿a qué número?" y "¿y mi caja?".
Les armé una muestra: su WhatsApp sigue igual, pero cada pedido queda anotado, con su estado, el aviso al cliente y el n.º de seguimiento.
Acá la ven en un minuto desde el celular: [VIDEO DE 60 S / LINK DEMO]
Si les tinca, paso por Alameda 2432 y se la muestro en 5 minutos, sin compromiso.
No les voy a pedir datos, claves ni pagos por este medio.
```

> La última línea va a propósito: le habla a una marca que vive alerta con las estafas.

---

## (b) Guion de 2 minutos en el local

### Checklist antes de la reunión

Se hace la noche anterior y se repasa 15 minutos antes de entrar al local.

- [ ] **Abrir el demo en una ventana privada** (incógnito) del navegador. Así no salen imágenes viejas guardadas en la caché (QA I-6).
- [ ] **Revisar que las ofertas de la vitrina sigan vigentes.** Por ejemplo, el tren navideño a $17.000 era "hasta el domingo" 15-11-2026: si ya venció, se saca el precio tachado o el producto desde Productos. La Vaca Lola llegaba aprox. el 16/11: actualizar o sacar la placa "PREVENTA".
- [ ] **Si el login quedó bloqueado** por claves mal tipeadas ("Demasiados intentos"), correr `docker compose restart api` y volver a entrar (QA M-14).
- [ ] **Tener listo el CSV de ejemplo:** exportarlo una vez desde Resumen → "Exportar a Excel (CSV)" y abrirlo en Excel para comprobar que se ve bien. Si falla algo en vivo, se muestra ese archivo.
- [ ] **Pedidos ficticios en su estado inicial** (ver §c). Si se ensayó el guion, volver a cargar la semilla.
- [ ] **Teléfono ficticio a mano** para registrar DC-0110: `+56 9 0000 0110`. Nunca un número real.
- [ ] **Pestañas abiertas:** "Hoy" en una, la landing en otra y el seguimiento de DC-0103 en una tercera.
- [ ] **Video de 60 s** descargado en el celular, como plan B.

### Preparación

- **Dónde corre:** el demo corre en el notebook de Luckas, en vista móvil (390 px). Los puertos solo escuchan en `127.0.0.1`, así que **no se abre desde un celular**. Plan B: el video de 60 s.
- **Modo demo:** en el demo **todos los pedidos son de prueba**, también los que se registran en vivo (MODO_DEMO). Los botones de mensajes abren WhatsApp **sin número** (tú eliges el contacto) y el texto parte con la línea `*[PRUEBA · DEMO]*`. Así es imposible que a un cliente real le llegue un mensaje del demo.
- **Landing:** **no se aprieta "enviar"**: se muestra el mensaje armado y se cierra.
- **Momento:** temprano, antes de la hora punta. No un sábado de temporada.

| Tiempo | Qué se dice | Qué se muestra |
|---|---|---|
| 0:00–0:15 | "Hola, soy Luckas, el que les escribió por Instagram. Dos minutos y me voy." | — |
| 0:15–0:35 | "Ustedes venden en el local y por WhatsApp, y ahora también a regiones. En los comentarios la gente pregunta a qué número escribir, si despachan y qué pasó con su caja. Con tanto pedido, no hay cómo llevar todo de memoria." | 2 o 3 comentarios reales, sin nombre de usuario |
| 0:35–0:55 | "**Tu WhatsApp sigue igual.** Lo único que cambia es que el cliente llega al número correcto con el mensaje ya escrito. Mira: 'Pedidos a regiones'… y se abre esto." | Landing → botón Regiones → mensaje (b) armado. **No se envía** |
| 0:55–1:30 | "Y esto es para quien atiende. Mientras chatea, anota el pedido en 30 segundos: teléfono, productos y comuna. Le sale el número DC-0110 y la cotización ya escrita, con un toque." | Panel → "+ Nuevo pedido" → pegar `+56 9 0000 0110` (ficticio) → guardar → "Enviar cotización" (vista previa, con la línea `*[PRUEBA · DEMO]*`) |
| 1:30–1:50 | "Y esta es la pantalla 'Hoy': quién lleva un día sin pagar, qué caja está pagada y todavía no sale, y qué está listo para retiro. Aprietas 'Despachar', pones la agencia, el número y la comuna, y queda despachado. Al cliente le mandas el aviso con un link para seguir su pedido." | "Hoy" → "Pagados sin despachar" → DC-0102 → **"Despachar"** → agencia, n.º y comuna → queda en Despachado → "Avisar despacho" (d4) → pestaña con el seguimiento |
| 1:50–2:00 | "Esto ya funciona. Para dejarlo andando necesito conversar 10 minutos con quien atiende regiones, para ajustarlo a cómo trabajan de verdad. ¿Me la presentas?" | — |

---

## (c) Cómo mostrar el demo: "así se vería tu día de pedidos a regiones"

Si el dueño o quien atiende regiones piden ver más, se recorre **un día completo** con pedidos ficticios. Son 3 a 4 minutos.

### Pedidos ficticios sugeridos (para la semilla de ingeniero-go)

Todos con `es_demo = true`, y con MODO_DEMO activo también los que se registran en vivo. Los nombres son inventados y los teléfonos, ficticios (formato `+56 9 0000 01xx`): el panel abre WhatsApp sin número y con la línea `*[PRUEBA · DEMO]*`.
- Productos y precios: de `research/productos-semilla.json`.
- Agencias, costos de envío, descuentos y n.º de seguimiento: `[EJEMPLO]`.
- Usuarios del demo: `Regiones 1`, `Santiago 1` y `Dueño` (no se inventan nombres del personal). `Regiones 1` con canal por defecto Regiones y `Santiago 1` con Santiago.
- Los ejemplos de `mensajes-whatsapp.md` usan estos mismos pedidos.

| Pedido | Canal | Cliente | Comuna / agencia | Detalle | Total | Estado al abrir el demo |
|---|---|---|---|---|---|---|
| DC-0101 | Regiones | Pedro Soto | Chillán / Starken, por pagar | 12 × Squishy $3.500 + 4 × Dinosaurio $14.000 | $98.000 | **Esperando pago hace 30 h** → alerta en "Hoy" |
| DC-0102 | Regiones | Marcela Fuentes | Concepción / Varmontt, envío $8.000 | 15 × Caja sorpresa premium $12.500, descuento $7.500 | $188.000 | **Pagado hace 20 h, sin despachar** → alerta |
| DC-0103 | Regiones | Jorge Muñoz | La Serena / Starken, n.º 000123456789 | 2 × Tren navideño grande $17.000 | $34.000 | Despachado ayer |
| DC-0104 | Santiago | Camila Rojas | Retiro, paga en tienda | 3 × Pulpo $4.000 + 1 × Telón $5.000 + 1 × "Bolsas de regalo surtidas" (línea libre) | $20.000 | **Listo para retiro hace 1 día** |
| DC-0105 | Regiones | Valentina Díaz | Puerto Montt / sin agencia aún | 2 × Juego magnético 139 piezas $20.000 | $40.000 | Cotizado hace 1 h |
| DC-0106 | Regiones | Luis Araya | Panguipulli / Pullman Cargo | 2 × Vaca Lola (preventa) $10.000 | $20.000 | Esperando pago hace 3 h (todavía sin alerta) |
| DC-0107 | Santiago | Andrea Pérez | Retiro, transferencia | 20 × Bolsa de regalo mediana $600 | $12.000 | Entregado (historial completo) |
| DC-0108 | Regiones | Rodrigo Vera | Valdivia / Starken | 10 × Pulpo viral $3.500 | $35.000 | Cancelado: "Se nos agotó el producto." |
| DC-0109 | Regiones | (sin nombre) | Rancagua / — | 2 × Tren navideño grande $17.000 | $34.000 | **Esperando pago hace 26 h** → alerta |

> Panguipulli se eligió porque es la comuna de un comentario real ("soy de Panguipulli", TT-7706669405505129735), pero el pedido es ficticio.

### Recorrido (lo que se dice y se muestra)

1. **"Así empieza tu día."** Abrir "Hoy":
   - "Hay dos clientes que llevan más de un día sin pagar (DC-0101 y DC-0109);
   - una caja pagada que no ha salido (DC-0102);
   - y un pedido listo para retiro (DC-0104)."
2. **"Un toque y le recuerdas el pago."** DC-0101 → "Recordar pago" → se ve el mensaje (d7) en WhatsApp, sin número y con `*[PRUEBA · DEMO]*`. **No se envía.**
3. **"Llega un cliente nuevo por el chat."** "+ Nuevo pedido", **con cronómetro a la vista**:
   - pegar el teléfono ficticio `+56 9 0000 0110` (nunca uno real);
   - Regiones;
   - buscar "squishy", 12;
   - comuna Talca;
   - Guardar → queda como DC-0110.
   - Luego "Enviar cotización" (d1).
   - Si el dueño pregunta por sus datos de transferencia, se muestra **Ajustes** ("Datos de transferencia: sin configurar") y se dice: "Aquí van tus datos de transferencia. **No se escriben en el panel ni me los mandas por chat: te los configuro yo en la instalación, en persona.** Salen solo dentro de la cotización y nunca se publican en la web."
4. **"Se pagó."** DC-0106 → "Marcar pagado":
   - el aviso "¿Revisaste que la plata está en la cuenta?";
   - "Avisar pago recibido" (d3).
5. **"Sale la caja."** Volver a "Hoy" → "Pagados sin despachar" → DC-0102 → **"Despachar"**:
   - Varmontt, el n.º y la comuna (Concepción);
   - el pedido pasa **directo a Despachado**, sin pasos intermedios;
   - "Avisar despacho" (d4).
6. **"Esto ve el cliente."** Abrir el link de seguimiento de DC-0103: estado, línea de tiempo, agencia y n.º de seguimiento. **No muestra nombre, teléfono, comuna, productos ni total**: si el cliente reenvía el link, nadie más ve qué compró ni cuánto pagó.
7. **"La vitrina la eliges tú."** Productos → un producto → marcar "Mostrar en la vitrina" y su orden (máximo 8) → la landing.
   - Mostrar una tarjeta con placa **"PROMO DEL VIDEO"** (por ejemplo, el telón a $5.000) y la nota `comenta "yo quiero mi promoción"`: "Los precios de promo por comentario salen marcados como promo, no como precio de detalle."
8. **"Y al final de la semana."** Resumen:
   - pedidos por canal y por estado;
   - **vendido**, que cuenta solo lo pagado o más avanzado (una cotización no es venta);
   - por agencia y por comuna;
   - por persona: cuántos registró y cuántos despachó;
   - "Exportar a Excel".
9. **Cierre:** "Todo esto es con pedidos inventados. Con los de verdad, en una semana tienes tus propios números."

---

## (d) Las 5 preguntas para quien atiende regiones

Se hacen **en persona**, con el dueño presente o con su permiso, y en 10 minutos. Se anotan las respuestas tal cual: son la línea base para medir (`propuesta-valor.md` §6).

1. **¿Cuántos pedidos de regiones atiendes en un día normal? ¿Y en un día fuerte de Navidad?**
   - *Para qué:* dimensionar el panel y la pantalla "Hoy". Saber si hace falta más de un usuario por canal.
2. **¿Dónde los anotas hoy?** (cuaderno, Excel, etiquetas de WhatsApp Business, notas del celular, en ninguna parte)
   - *Para qué:* ver de qué se migra y qué columnas ya usan. Si hay un Excel, el CSV debería calzar con él.
3. **¿Cómo confirmas que un cliente pagó?** (miras la cuenta del banco, te mandan el comprobante, le preguntas al dueño)
   - *Para qué:* ajustar el paso "Pagado" y quién puede marcarlo. Saber si han tenido comprobantes falsos.
4. **¿Por qué courier o agencia despachan? ¿Quién paga el envío? ¿Hay compra mínima?**
   - *Para qué:* llenar la lista de agencias y los `[CONFIRMAR CON EL DUEÑO]` de "Cómo comprar desde regiones".
5. **¿Qué es lo que más se te complica?** (dejar que responda; después, si hace falta: ¿los pagos?, ¿armar las cajas?, ¿los reclamos?, ¿que te escriban por todos lados?)
   - *Para qué:* decidir qué se pule primero. Si lo que más le pesa no está en el demo, va primero en la fase 2.

**Preguntas de repuesto:**
- ¿Es un número aparte para regiones o el mismo +56 9 2114 5368?
- ¿Cuántas personas atienden cada número?
- ¿Usan WhatsApp Business en el celular o en el computador?

---

## (e) Otras preguntas para el dueño (10 minutos, después del demo)

Salen de `research/pendientes.md`:
1. **Números:** ¿cuál es el de regiones? ¿Quién atiende cada uno?
2. **Datos de transferencia:** ¿a qué cuenta? **No se escriben en el panel ni se mandan por chat:** Luckas los configura en la instalación, en persona y con el dueño al lado. **Luckas no los pide por chat ni los anota en ningún documento.** En el panel, Ajustes solo muestra "configurados" o "sin configurar".
3. **Horario** vigente y de temporada.
4. **Retiro:** ¿se pueden dejar pedidos listos? ¿Cuántos días se guardan? ¿En cuál de los dos locales se retira?
5. **Promo:** ¿cómo se canjea "yo quiero mi promoción"? ¿Ese precio vale también por WhatsApp o solo comentando en el video?
6. **Nombre y cuentas:** ¿"Donde Carlitos Emporio", "Mayte Emporio" o "Familia Emporio"? ¿Qué cuentas son oficiales?
7. **Permisos:** fotos, videos y logo (archivo original). ¿Se incluye Temuco?
8. **Cambios y devoluciones**, y qué hacer si se cancela un pedido pagado.
9. **Equipo:** ¿cuántas personas van a usar el panel y en qué canal atiende cada una? (para crear sus usuarios con su canal por defecto).

---

## (f) Objeciones probables

**1. "Ya vendo bien por TikTok y WhatsApp, no necesito nada."**
> "Justamente: esto no cambia cómo vendes. Tu WhatsApp sigue igual. Lo que cambia es que cada pedido queda anotado, con su estado y quién lo atendió, y el cliente sabe en qué va su caja sin tener que reclamar en los comentarios."
- Evidencia: TT-7706913997378948370; `faq-clientes.md`, quejas.

**2. "No tenemos tiempo para anotar los pedidos."**
> "Son 30 segundos mientras chateas: pegas el teléfono, buscas el producto y pones la comuna. Y te devuelve el tiempo, porque la cotización, el aviso de pago, el de despacho y el de retiro ya están escritos. Y para despachar, desde 'Hoy' aprietas 'Despachar', pones agencia, número y comuna, y listo. Probemos un día con regiones y vemos si le sirve a quien atiende."
- Ojo: los 30 s son la meta de diseño (en QA tomó 13,6 s con datos de prueba). Se demuestra con cronómetro; no se promete como cifra.

**3. "¿Cuánto cuesta?"**
> "Lo dividí en partes. La salida a producción, que es lo básico para que funcione con tus números, tu dominio y tus usuarios, cuesta [COMPLETAR LUCKAS] y queda lista en [COMPLETAR LUCKAS]. Aparte van el dominio y el servidor, que son costos chicos y te los detallo por escrito. Lo automático con WhatsApp oficial es opcional y tiene un costo de Meta por mensaje, que también te detallo."
- Detalle en `fases.md`.

**4. "WhatsApp Business ya tiene etiquetas y respuestas rápidas." / "¿Y un bot?"**
> "Las etiquetas sirven y las puedes seguir usando. Pero no te dicen cuánto vendiste en la semana ni quién despachó qué, no le dan al cliente un link para seguir su caja, y si hay dos números o dos personas, cada una ve solo lo suyo. Lo de responder solo, el bot, es la fase 2: se hace con la API oficial de WhatsApp, que pide verificar el negocio con Meta y cobra por mensaje. Primero ordenamos; después automatizamos."

**5. "¿Y si se filtran los datos de mis clientes o mi cuenta del banco?"**
> "Cada persona entra con su usuario y su clave, y la clave se guarda cifrada. Cuando sumas a alguien, entra con una clave temporal y **la tiene que cambiar la primera vez que entra**, así nadie más la conoce. Los datos de tu cuenta **no se escriben en el panel ni me los mandas por chat: te los configuro yo en la instalación, en persona.** En el panel solo se ve si están configurados o no; no se publican en ninguna parte y solo salen dentro del mensaje de cotización que manda tu equipo. La página de seguimiento no muestra nombre, teléfono, comuna, productos ni total. Y para la salida a producción reviso la nueva ley de datos personales, que el Diario Oficial fija para diciembre de 2026."
- Evidencia: CLAUDE.md, Seguridad; `seguridad/reporte-2.md` (B9, cambio de clave obligatorio). Ley 21.719 `[VERIFICAR si se posterga]` (ver `fases.md`).

**6. "Mi gente prefiere venir al local."**
> "Perfecto, el local sigue siendo lo principal. Esto es solo para lo que ya entra por WhatsApp. Y si quieren dejar pedidos listos para retiro, el panel les avisa cuáles están esperando."

**7. "¿Y si la persona que atiende no lo usa?"**
> "Por eso quiero hablar con ella antes de cerrar nada. Si algo le complica, se cambia. Cada persona tiene su usuario y su canal por defecto, así que el panel ya le abre en Santiago o en Regiones, según lo que atienda. La capacitación es en el local, con sus pedidos reales, y queda una guía de una página."

---

## (g) Qué NO decir

- **"Vas a responder más rápido" o "vas a vender más".** El demo ordena, pero no agrega gente. No se dan cifras.
- **"Los clientes no te contestan" o "tu equipo no responde".** Se dice "no dan abasto con tanto pedido". Las quejas se muestran sin culpar a nadie.
- **"Envíos a todo Chile", "pagos en línea" o "respuestas automáticas"** como si existieran.
- **"Escribe tus datos de transferencia en Ajustes".** Ajustes es de solo lectura: los datos de transferencia los configura Luckas en la instalación.
- **"Puedes cambiar tus números, horario o dirección desde el panel".** En el demo, Ajustes solo los muestra; editarlos desde el panel es fase 2.
- **"La niña de los pedidos".** Se dice "quien atiende regiones".
- **Nada sobre la competencia, el ambiente del local ni la historia personal del dueño.**
- **No pedir datos bancarios, claves ni contraseñas** por ningún medio. Los datos de transferencia se configuran en persona, en la instalación, y no quedan anotados en ningún documento.
- **No registrar teléfonos reales en el demo.** Siempre `+56 9 0000 01xx`.
