# Arquitectura: landing, seguimiento y panel — Donde Carlitos Emporio

*estratega-propuesta · 2026-11-14 · para constructor-web e ingeniero-go · enfoque "Tu WhatsApp sigue igual, pero cada pedido queda ordenado"*

**Insumos:**
- texto final en `copy-landing.md`;
- mensajes exactos en `mensajes-whatsapp.md`;
- estilo en `.claude/skills/estilo-carlitos/SKILL.md`;
- evidencia en `propuesta-valor.md`.

**Qué ya no existe:** catálogo navegable con carrito, carrito, checkout, formulario de pedido, confirmación de pedido web, botón "Agregar" ni contador del carrito. **Se eliminan del front y del contrato.** El estilo-carlitos §8.1 y §8.4 todavía hablan de "Agregar" y "Enviar pedido (3)": esa parte de la skill queda superada por este documento.

**Principio:** la landing tiene **un solo trabajo**, que es mandar a cada persona al WhatsApp correcto con el mensaje armado y quitarle el miedo a comprar desde lejos. Todo lo demás (vitrina, cómo llegar, preguntas) existe para que ese clic pase.

**Supuestos de uso:**
- móvil a 360 px, en el navegador interno de TikTok o IG;
- datos móviles: sin autoplay ni embeds pesados; imágenes comprimidas;
- la página tiene que funcionar aunque falle la API: los botones de WhatsApp no dependen de ella.

---

## A. Landing pública (`/`)

### 0. Header fijo
- **Objetivo:** marca siempre visible.
- **Contenido:** logo de 44 px sobre negro y un link "Sigue tu pedido", que baja a la sección 7.
- Sin carrito ni buscador.

### 1. Hero
- **Objetivo:** en 3 segundos, qué es, qué venden y **a qué botón apretar**.
- **Contenido:**
  1. logo de 120–160 px;
  2. titular, bajada y qué venden (`copy-landing.md` §2). Esto resuelve QA I-5, "el hero no dice qué venden";
  3. **dos botones grandes**, uno debajo del otro en móvil:
     - **"Pedidos Santiago / retiro en tienda"**: abre `WHATSAPP_SANTIAGO` con el mensaje (a);
     - **"Pedidos a regiones"**: abre `WHATSAPP_REGIONES` (o el de Santiago si está vacío) con el mensaje (b);
  4. bajo los botones, en chico: "Estos son nuestros únicos números de pedidos", con link a la sección 6;
  5. barra de datos: 📍 Alameda 2432 · Metro Unión Latinoamericana · 🕘 horario de hoy.
- **Estilo:**
  - Santiago: fondo `--c-whatsapp` con texto negro;
  - Regiones: fondo `--c-whatsapp-osc` con texto blanco (contraste 6,3);
  - los dos con el ícono de WhatsApp, 56 px de alto como mínimo y ancho completo. Así se distinguen sin leer.
- **Datos:** ninguno de la API. Los números van por configuración (`mensajes-whatsapp.md` §1).
- **Evidencia:** `faq-clientes.md` §1, §5; `dolores-pedido.md` B2.

### 2. Lo nuevo de los videos (vitrina liviana)
- **Objetivo:** que quien llega desde un video encuentre ese producto y lo pida con un toque.
- **Contenido:** de 6 a 8 tarjetas. Cada una lleva:
  - foto 4:5;
  - nombre;
  - precio detalle y, si hay, precio por mayor;
  - "Precio publicado el dd/mm";
  - link "▶ Míralo en TikTok" al video de origen;
  - botón **"Pedir este por WhatsApp"**.
- **Al tocar "Pedir este por WhatsApp"** se abre una hoja inferior: "¿Dónde lo quieres?", con dos botones (Santiago / retiro, o Regiones / envío). Luego se abre el número del canal con el mensaje (c). La elección se recuerda durante la visita (`sessionStorage`) y se puede cambiar.
- **Sin stock en pantalla.** La vitrina no dice "hay stock" ni "quedan pocas": el mensaje pregunta "¿Les queda?". Así no se promete stock que nadie actualiza (QA I-1).
- **Datos (API):** `GET /api/vitrina`, con los productos marcados `en_vitrina`, en orden y con un máximo de 8. Campos: `codigo`, `nombre`, `precio`, `precio_mayor`, `unidad_mayor`, `fecha_precio`, `imagen_url` y `tiktok_url`.
- **Productos sugeridos** (todos con precio hablado en videos de nov 2026; ver `research/productos-semilla.json`):

| # | Producto | Precio detalle | Por mayor | Fuente |
|---|---|---|---|---|
| 1 | Squishy viral que cambia de color con el sol | $4.000 | $3.500 c/u, caja de 12 | TT-7707649267464293639 (13/11) |
| 2 | Tren navideño para árbol (grande) | $17.000 (antes $20.000) | $15.000 c/u, embalaje | TT-7706656443843284232 (10/11) |
| 3 | Dinosaurio gigante con ruedas | $16.000 | $14.000 c/u desde 4; $12.500 por embalaje | TT-7707053973999471879 (11/11) |
| 4 | Caja sorpresa squishy premium | $15.000 | $13.500 desde 10; $12.500 desde 15 | TT-7707296451398864136 (12/11) |
| 5 | Pulpo viral (mueve ojos y tentáculos) | $4.000 | $3.500 c/u, caja surtida | TT-7706989825747995911 (11/11) |
| 6 | Juego magnético didáctico 139 piezas | $20.000 | $17.000 c/u, embalaje | TT-7706669405505129735 (10/11) |
| 7 | Telón de Navidad 3 m x 2,5 m | $5.000 | — | TT-7707411807194926343 (12/11). Era precio promo por comentario |
| 8 | Vaca Lola (preventa) | $10.000 | $7.000 c/u, embalaje | TT-7706993060130950408 (11/11). Llegada aprox. 16/11 |

- **Reglas:**
  - el precio tachado del tren solo va mientras su oferta esté vigente (`oferta_hasta`; según el video, "hasta el domingo" 15-11-2026);
  - la Vaca Lola va con el badge "PREVENTA · llega aprox. 16/11" y se saca o se actualiza cuando llegue;
  - el dueño elige y ordena la vitrina desde el panel;
  - créditos de fotos y videos en `web/public/CREDITOS.md`.

### 3. Cómo comprar desde regiones (5 pasos)
- **Objetivo:** responder "¿hacen envíos?" y "¿es seguro?" con un proceso claro.
- **Contenido:** 5 pasos numerados (Escribes → Te cotizamos → Transfieres → Mandas el comprobante → Despachamos con n.º de seguimiento). Debajo, un cuadro "Lo que tienes que saber", con agencias, quién paga el envío, mínimo de compra y plazo de despacho, **todo [CONFIRMAR CON EL DUEÑO]**, y el botón "Pedidos a regiones".
- **Regla:** **no escribir "Envíos a todo Chile"** mientras no se confirme. En el demo, los datos dudosos se muestran con una etiqueta gris chica: "dato por confirmar".
- **Evidencia:** tr. TT-7707471893363789063 ("número para región"); TT-7706913997378948370; `pendientes.md` §3.

### 4. ¿Eres de Santiago? Retiro en tienda
- **Objetivo:** dar al cliente de Santiago su camino corto.
- **Contenido:** 3 líneas (escríbenos o ven directo · te confirmamos el total · retiras con tu número de pedido) + el botón "Pedidos Santiago / retiro en tienda".
- **Nota:** es corto a propósito, porque el canal principal en Santiago es el local (`adn-marca.md` §G).

### 5. Cómo llegar y horario
- **Objetivo:** cerrar la pregunta n.º 1 (~45 comentarios).
- **Contenido:**
  - dirección;
  - Metro Unión Latinoamericana (5 min) y República (7 min);
  - terminales de buses cercanos;
  - botón "Abrir en Google Maps" (link, sin iframe);
  - horario con fecha de publicación;
  - cómo se paga en el local.
- **Evidencia:** bios oficiales; [4][6][7]; tr. TT-7707464302034160914.

### 6. Canal oficial (antiestafa)
- **Objetivo:** que el cliente sepa cuál es la tienda real **antes de transferir**.
- **Contenido:**
  - los dos números (el de regiones, si es distinto) con la etiqueta "Únicos números de pedidos";
  - TikTok @mayte_emporio e Instagram @mayte_emporioo (verificada);
  - 5 reglas para comprar seguro;
  - el aviso "No publicamos datos bancarios en la web: te los damos solo por el chat oficial".
- **Reglas:**
  - no se nombran ni se enlazan las cuentas no confirmadas (@_dondecarlitosemporio7, páginas de Facebook);
  - **no se publican datos bancarios** (OP "Qué NO").
- **Ubicación:** va en el 6.º lugar, pero el hero ya enlaza a esta sección y su versión corta se repite en el footer.

### 7. Sigue tu pedido
- **Objetivo:** explicar que el seguimiento existe y cómo llega.
- **Contenido:** "Cuando hagas tu pedido por WhatsApp te mandamos un link para ver en qué va, sin dar tus datos".
- **Importante:** **no hay buscador por número de pedido** en la landing. Los `DC-0101` son correlativos, así que un buscador permitiría ver pedidos ajenos. Solo se entra con el link que trae el token.

### 8. Quiénes somos (Familia Emporio), corta
- **Objetivo:** confianza y cercanía sin inventar.
- **Contenido:** 2–3 líneas y la miniatura del video de un contenedor (TT-7441724226857454904 o TT-7683668603270352136).
- **Regla:** sin nombres de personas, cifras de la empresa ni historia del fundador hasta que el dueño lo apruebe **[CONFIRMAR CON EL DUEÑO]**.

### 9. Preguntas frecuentes
- **Objetivo:** las preguntas reales, en orden de frecuencia (`faq-clientes.md`).
- **Contenido:** un acordeón de 10–12 preguntas (`copy-landing.md` §10).

### 10. Footer
- **Contenido:** logo de 64 px; dirección; horario; los números oficiales; las cuentas; el aviso corto antiestafa; y, en el demo, "Demo de prueba".

### 11. Barra fija abajo (móvil)
- **Contenido:** **dos botones** lado a lado, "Santiago / retiro" y "Regiones", con el ícono de WhatsApp y 52 px de alto. Los mismos colores del hero.
- **Comportamiento:** aparece después de pasar el hero, para no duplicar los botones en la primera pantalla.
- **Padding inferior** en la página, para que la barra no tape el último botón de la vitrina.

---

## B. Página pública de seguimiento (`/seguimiento/{token}`)

- **Objetivo:** que el cliente vea en qué va su pedido sin escribir "¿y mi caja?".
- **Contenido:**
  1. logo y "Pedido DC-0101";
  2. el estado actual, grande, con su explicación en una línea;
  3. una línea de tiempo con los estados por los que pasó y su fecha;
  4. si está despachado: agencia + n.º de seguimiento (con botón "Copiar");
  5. si está listo para retiro: dirección y horario;
  6. productos y cantidades, y el total;
  7. un botón "Escríbenos por WhatsApp", que abre el número del canal del pedido con el texto `Hola, les escribo por mi pedido {numero}.`;
  8. el bloque corto de canal oficial.
- **No muestra:** nombre, teléfono, comuna, notas internas ni quién atendió.
- **Token inválido o vencido:** "No encontramos este pedido. Revisa el link o escríbenos por WhatsApp". Mismo mensaje y mismo código HTTP en los dos casos, para no revelar si el pedido existe.
- **Técnico:**
  - `noindex`;
  - `Referrer-Policy: no-referrer`;
  - límite de peticiones por IP;
  - token de 16 bytes aleatorios en base64url (`mensajes-whatsapp.md` §5).

---

## C. Panel privado (`/admin`)

**Usuarios:** quienes atienden los WhatsApp y el dueño. Cada uno con su nombre, su contraseña (bcrypt) y su canal por defecto (Santiago o Regiones). Rol `admin` (el dueño: usuarios, ajustes, vitrina) o `atencion`.

### C1. Hoy (pantalla de inicio)
- **Objetivo:** la lista de tareas del día.
- **Bloques, en este orden:**
  1. **Esperando pago hace más de 24 h:** cada tarjeta trae el botón "Recordar pago" (d7) y "Cancelar".
  2. **Pagados sin despachar:** pedidos de despacho en estado Pagado o Preparando. Se ordenan del más antiguo al más nuevo y llevan la etiqueta "hace N h".
  3. **Listos para retiro:** con "hace N días".
  4. **Cotizados hoy, sin respuesta** *(sugerencia)*.
- **Arriba:** el contador del día (pedidos nuevos, pagados, despachados y entregados) y un filtro "Mi canal / Todos".
- **Botón fijo:** "+ Nuevo pedido".

### C2. Nuevo pedido (registro en menos de 30 s)
- **Objetivo:** anotar el pedido sin soltar el chat.
- **Orden de los campos (pensado para el pulgar):**
  1. **Canal:** dos botones, con el canal del usuario preseleccionado;
  2. **Teléfono:** campo para pegar. Se normaliza y, si ya existe, muestra "Cliente conocido: {nombre}, {n} pedidos" y rellena el nombre y la última comuna;
  3. **Nombre;**
  4. **Productos:** un buscador del catálogo interno por nombre o código, que trae el precio sugerido (editable), más "+ Línea libre" (descripción, cantidad y precio opcional);
  5. **Total:** se calcula solo; es editable si alguna línea va sin precio;
  6. **Si es despacho:** comuna (autocompletado con las comunas de Chile) y agencia (botones con las agencias configuradas `[EJEMPLO]`, más "Otra");
  7. **Nota interna** (opcional).
- **Guardar** deja el pedido en **Cotizado** y abre la ficha con el botón "Enviar cotización" (d1) a la vista.
- **Medición:** el demo debería medir el tiempo de registro en las pruebas de QA, con un cronómetro manual. La meta de 30 s es **de diseño**, no una métrica prometida.

### C3. Ficha del pedido
- **Arriba:** número, estado, canal y "Registrado por {usuario} · {fecha hora}".
- **Acción siguiente:** un solo botón grande con lo que toca según el estado (ver la tabla), y debajo el mensaje que corresponde.
- **Datos editables:** líneas, total, comuna, agencia, n.º de seguimiento, pago (transferencia / en tienda) y entrega (despacho / retiro).
- **Mensajes:** una lista con los mensajes de `mensajes-whatsapp.md` que aplican, cada uno con una vista previa y el botón "Abrir en WhatsApp". Al lado: "Envíalo desde el WhatsApp de {canal}".
- **Historial:** `{fecha hora} · {estado anterior} → {estado nuevo} · {usuario}`.

**Transiciones sugeridas (para ingeniero-go):**

| Desde | Hacia | Condición |
|---|---|---|
| Cotizado | Esperando pago | pago = transferencia |
| Cotizado | Preparando | pago = en_tienda (retiro en Santiago, paga al retirar) |
| Esperando pago | Pagado | Confirmación con aviso: "¿Revisaste que la plata está en la cuenta?" |
| Pagado | Preparando | — |
| Preparando | Despachado | entrega = despacho; **exige** agencia, n.º de seguimiento y comuna |
| Preparando | Listo para retiro | entrega = retiro |
| Despachado / Listo para retiro | Entregado | — |
| Cualquiera, salvo Despachado y Entregado | Cancelado | Motivo opcional |

**Mensaje sugerido por estado:**

| Estado al que se llega | Mensaje que el panel ofrece primero |
|---|---|
| Cotizado | (d1) Cotización |
| Esperando pago | (d2) Datos de transferencia, si existen; y si pasan 24 h, (d7) Recordatorio |
| Pagado | (d3) Pago recibido |
| Despachado | (d4) Despachado |
| Listo para retiro | (d5) Listo para retiro |
| Cancelado | (d6) Cancelado |

### C4. Buscar
- **Contenido:** un campo único que acepta teléfono (en cualquier formato), número de pedido o nombre.
- **Resultado por teléfono:** todos los pedidos de ese cliente, del más nuevo al más antiguo.

### C5. Resumen
- **Contenido:**
  - un selector Hoy / Esta semana / Rango;
  - pedidos por estado y por canal;
  - monto total y ticket promedio de los pedidos **pagados o entregados**;
  - las 5 comunas y las agencias más frecuentes;
  - cancelados con su motivo;
  - y, por usuario, cuántos registró y cuántos despachó.
- **Exportar CSV:** UTF-8 con BOM y separador `;`, para que Excel en español de Chile lo abra bien al hacer doble clic. Columnas en `copy-landing.md` §P5.

### C6. Vitrina y catálogo interno
- **Objetivo:** el catálogo sirve para registrar pedidos rápido. La vitrina es un subconjunto marcado.
- **Campos:** nombre, código, precio, precio por mayor y unidad, fecha del precio, foto (subida ya implementada en I-4), link TikTok y "Mostrar en la vitrina" (máximo 8, con orden).

### C7. Ajustes (solo admin)
- **Contenido:**
  - usuarios;
  - `DATOS_TRANSFERENCIA` (un textarea con el aviso "Se manda solo por WhatsApp, nunca se publica en la web");
  - horario de retiro;
  - dirección de retiro;
  - lista de agencias.
- **En el demo:** si los números de WhatsApp van por variables de entorno, en Ajustes solo se muestran, sin editarlos.
