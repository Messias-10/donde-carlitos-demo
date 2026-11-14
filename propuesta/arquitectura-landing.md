# Arquitectura de la landing — Donde Carlitos Emporio

*estratega-propuesta · Fase 2 · 2026-11-14 · para constructor-web*

**Insumos.** El texto final de cada sección está en `copy-landing.md`. El mensaje de pedido está en `mensaje-pedido-whatsapp.md`. El estilo visual viene de `.claude/skills/estilo-carlitos/SKILL.md`. Las referencias de evidencia son las mismas que en `propuesta-valor.md`.

**Principio:** la primera pantalla tiene que hacer lo que ya hace el caption de cualquier emporio del barrio: dirección, Metro, horario, WhatsApp y precio (C-conclusión 2, OP-9). Después viene lo que nadie tiene: catálogo con precio detalle/mayor, pedido ordenado y canal oficial.

**Supuestos de uso:**
- Móvil a 360 px.
- Navegador interno de TikTok o IG.
- Datos móviles: paginar el catálogo y comprimir imágenes. Nada de autoplay ni de embeds pesados (OP "Qué NO").

## Cambios respecto a la base sugerida (y por qué)

| Cambio | Por qué |
|---|---|
| La dirección, el Metro y el horario **suben al hero** como una barra de datos, además de tener su sección completa más abajo | "¿Dónde están?" es la pregunta n.º 1 (~45 comentarios) (`faq-clientes.md` §1) |
| **Ofertas / remate va en el 2.º lugar**, antes de "Lo viste en TikTok" | Se presenta en noviembre, con los remates de Navidad en curso (`ofertas.md`, OP-6). La oferta es la firma de la marca |
| Se agrega **"¿Pa negocio?"** (precio por mayor) como sección propia | Precio por embalaje o por mayor en casi todos los videos, y ~15 preguntas de "¿cuántos por embalaje?" (`ofertas.md`, `faq-clientes.md` §4) |
| Se agrega **"Familia Emporio"** (quiénes somos), corta y opcional | Los videos de contenedor generan confianza y simpatía (`top-contenido.md` n.º 5) |
| Retiro y envíos se separan de "Cómo llegar" | Son dos preguntas distintas: "¿hacen envíos?" (~30) y "¿dónde están?" (~45) |
| El **botón flotante** cambia según el carrito | Con productos en el carrito tiene que llevar al pedido (que se guarda en la base), no a un chat suelto. Si no, el pedido se salta el registro con número |

## Orden de las secciones

### 0. Header fijo
- **Objetivo:** marca siempre visible y carrito a un toque.
- **Contenido:**
  - logo de 44 px a la izquierda, sobre negro;
  - buscador (lupa que abre un campo);
  - carrito con contador flúor.
- **Datos:** cantidad del carrito (local, en el navegador).
- **Nota:** se mantiene liviano; sin menú de 8 ítems. Para las anclas a las secciones basta un menú desplegable.

### 1. Hero
- **Objetivo:** que quien llega desde un video entienda en 3 segundos qué es, dónde está y cómo se pide.
- **Contenido:**
  1. logo de 120–160 px;
  2. titular y bajada (`copy-landing.md` §2);
  3. CTA principal "Ver ofertas" (ancla a la sección 2) y CTA secundario "Ver catálogo";
  4. **barra de datos:** 📍 Alameda 2432 · Metro Unión Latinoamericana · 🕘 horario de hoy, con un estado "Abierto / Cerrado" calculado en el cliente con el horario vigente.
- **Evidencia:** OP-9; `faq-clientes.md` §1 y §6.
- **Nota:** el estado Abierto/Cerrado depende de qué horario sea el vigente **[CONFIRMAR CON EL DUEÑO]**. Mientras no esté confirmado, mostrar el horario con la fecha de su fuente y sin el estado, o con la etiqueta "por confirmar".

### 2. Ofertas / Remate de temporada
- **Objetivo:** mostrar la oferta como la cantan en los videos, con precio y escalera, y llevar al carrito.
- **Contenido:**
  - **un solo banner** de temporada (hoy "Remate de Navidad"; componente 8.3 de estilo-carlitos). Sin carrusel automático;
  - debajo, una fila deslizable de tarjetas con badge OFERTA o REMATE.
- **Datos (API):**
  - `precio_oferta` y `etiqueta_oferta` (oferta | remate);
  - `temporada = navidad`;
  - el precio tachado es `precio_detalle`.
- **Reglas:**
  - el tachado solo se carga si el precio anterior es real (ej.: tren de $20.000 a $17.000, TT-7706656443843284232);
  - si el video solo dijo el precio promo (ej.: esferas a $500), el producto va **sin oferta** y con ese precio como detalle, para no inventar un tachado;
  - la urgencia ("hasta el domingo") solo va si tiene fecha real y vigente.
- **Evidencia:** `ofertas.md` (fórmula Carlitos); estilo-carlitos §6.

### 3. Lo viste en TikTok
- **Objetivo:** unir el video con el producto comprable. Es el puente con el tráfico que llega desde redes.
- **Contenido:** de 4 a 6 miniaturas verticales 9:16 con ícono ▶. Cada una enlaza al video real en TikTok (pestaña nueva) y tiene debajo la tarjeta del producto con precio y "Agregar".
- **Datos:** `tiktok_url` e `imagen_url` del producto.
- **Videos sugeridos** (los de producto que hoy está en la semilla primero):

| Prioridad | Video | Producto | Nota |
|---|---|---|---|
| 1 | TT-7706656443843284232 | Tren navideño grande | Navidad, precio con rebaja real |
| 2 | TT-7707649267464293639 | Squishy viral que cambia de color | Nov 2026, detalle + caja de 12 |
| 3 | TT-7707296451398864136 | Caja sorpresa squishy | Escalera por volumen |
| 4 | TT-7707053973999471879 | Dinosaurio gigante con ruedas | Detalle, mayor de 4 y embalaje |
| 5 | TT-7671107953398664466 | Espada con luces | 2,0 M de vistas, dice Alameda 2432 |
| 6 | TT-7699270769963978002 | Calendario de adviento squishy | 1,4 M de vistas, temporada |

- **Reglas:**
  - si se usa un video con dirección vieja (⚠ en `top-contenido.md`, ej. las repisas TT-7476070344000933125), se superpone "Hoy nos encuentras en Alameda 2432";
  - **no usar** el sauna (TT-7627590547473059090) en esta sección, porque dice "solamente venta presencial" y contradice el pedido web;
  - todo material se anota en `web/public/CREDITOS.md`.
- **Evidencia:** `top-contenido.md`; OP-8; [25].

### 4. Categorías
- **Objetivo:** que el que no vino por un video encuentre su sección en un toque.
- **Contenido:** una grilla de 2 columnas en móvil, con nombre y cantidad de productos.
- **Datos:** `GET categorias` con `TotalProductos`.
- **Categorías sugeridas** (de `adn-marca.md` §E y `productos-semilla.json`):
  1. Navidad;
  2. Juguetes y virales;
  3. Hogar y organización;
  4. Escolar;
  5. Textil;
  6. Electro;
  7. Cuidado personal;
  8. Cotillón.
- **Nota:** las categorías sin productos en la semilla no se muestran vacías. Las licencias Disney se agregan solo si el dueño quiere separarlas **[CONFIRMAR CON EL DUEÑO]**.

### 5. ¡Llegó esta semana! (recién llegado)
- **Objetivo:** reflejar el "¡Llegaron…!" de los videos. La novedad es lo que más vende (`top-contenido.md`, patrones).
- **Contenido:** una fila deslizable con badge "¡LLEGÓ!" y la fecha de llegada.
- **Datos:** `fecha_llegada` reciente (filtro `LlegadaDesde`).
- **Nota:** los agotados se ven en gris como "Agotado" en vez de desaparecer (OP-7). Así el cliente puede preguntar por la reposición.

### 6. ¿Pa negocio? Precio por mayor
- **Objetivo:** hablarle directo al revendedor, que es cliente de primera (estilo-carlitos, regla 8).
- **Contenido:**
  - explicación en 3 líneas de detalle → por mayor desde N → por embalaje;
  - ejemplo con un producto real (squishy: $4.000 detalle, $3.500 por caja de 12);
  - CTA "Ver productos con precio por mayor".
- **Datos:** `precio_mayor`, `mayor_minimo` y `unidad_mayor`. El precio por mayor se aplica en el carrito al llegar al mínimo **[VERIFICAR con ingeniero-go que el cotizador lo hace así]**.
- **Nota:** el "mínimo de compra" general por mayor no está confirmado (`pendientes.md` §3). Solo se muestran los mínimos por producto.

### 7. Cómo pedir en 3 pasos
- **Objetivo:** bajar el miedo de "¿cómo lo hago?" (`faq-clientes.md` §5) y dejar claro que la web **no cobra**.
- **Contenido:** 3 pasos con ícono y una línea cada uno (Elige → Envía → Te confirmamos).
- **Nota:** se repite en miniatura dentro del carrito.

### 8. Retiro en tienda y envíos
- **Objetivo:** responder "¿hacen envíos?" en una línea visible.
- **Contenido:**
  - dos tarjetas: "Retiro en Alameda 2432" y "Envío por agencia";
  - en la tarjeta de envío: agencia a elección del cliente, quién paga y monto mínimo, todo **[CONFIRMAR CON EL DUEÑO]**.
- **Datos:** `entrega` (retiro | envio) y `agencia`. Las agencias que muestre el selector son `[EJEMPLO]` hasta que el dueño confirme.
- **Evidencia:** la versión más reciente es nov 2026, "número para región" (tr. TT-7707471893363789063). Antes decía "solamente venta presencial" (2025 y abril 2026). Ver `pendientes.md` §3.
- **Regla:** **no escribir "Envíos a todo Chile"** como hecho mientras no se confirme.

### 9. Cómo llegar
- **Objetivo:** cerrar la pregunta n.º 1 con todo el detalle.
- **Contenido:**
  - dirección;
  - Metro Unión Latinoamericana (5 min) y República (7 min);
  - terminales de buses cercanos;
  - mapa;
  - horario con fecha de vigencia;
  - medios de pago por local.
- **Mapa:** un botón "Abrir en Google Maps" o una imagen estática con link, mejor que un iframe pesado. Si se usa un iframe, su dominio va en la CSP.
- **Evidencia:** bios oficiales (`pendientes.md` §3); [4][6][7]; tr. TT-7707464302034160914 (pagos).
- **Nota:** el local de Temuco (Manuel Rodríguez 1232) va como bloque aparte **solo si el dueño lo aprueba** **[CONFIRMAR CON EL DUEÑO]**.

### 10. Canal oficial (antiestafa)
- **Objetivo:** que el cliente sepa cuál es la tienda real y qué hacer si le escriben a nombre de ella.
- **Contenido:**
  - el único número;
  - TikTok @mayte_emporio e Instagram @mayte_emporioo (verificada), con link;
  - "pedidos solo a este número";
  - 3 reglas simples;
  - en F2, el dominio oficial.
- **Dónde más aparece:** resumido en el footer y en la pantalla de confirmación del pedido (OP-1).
- **Reglas:**
  - **no se nombran ni se enlazan** las cuentas no confirmadas (@_dondecarlitosemporio7, páginas de Facebook);
  - no se publican datos bancarios (OP "Qué NO").

### 11. Familia Emporio (opcional)
- **Objetivo:** confianza y cercanía sin inventar nada.
- **Contenido:**
  - 2–3 líneas;
  - la miniatura del video de apertura de contenedor (TT-7441724226857454904 o TT-7683668603270352136).
- **Regla:** sin nombres de personas, cifras de la empresa (contenedores, trabajadores) ni historia personal del fundador hasta que el dueño lo apruebe (estilo-carlitos §7, NO 5).

### 12. Preguntas frecuentes
- **Objetivo:** contestar las preguntas reales, en el orden de su frecuencia.
- **Contenido:** un acordeón de 10–12 preguntas (`copy-landing.md` §12).
- **Evidencia:** `faq-clientes.md`.

### 13. Footer
- **Contenido:**
  - logo de 64 px;
  - dirección, horario y número único;
  - cuentas oficiales;
  - aviso corto de canal oficial;
  - link a "Precios referenciales";
  - en el demo: "Demo — pedidos de prueba".

### 14. Botón flotante de WhatsApp (sticky abajo, siempre)
- **Sin carrito:** "Escríbenos por WhatsApp". Abre `wa.me` con el mensaje de consulta de `mensaje-pedido-whatsapp.md` §5.
- **Con carrito:** "Ver mi pedido (N) · $total". Lleva al carrito y **no** abre WhatsApp directo, porque el pedido tiene que guardarse primero.
- **Estilo:** texto negro sobre el verde de WhatsApp, 52 px de alto. Abajo, en chico: "Único número de pedidos".
- **Nota:** no debe tapar el botón "Agregar" de la última tarjeta. Se deja un padding inferior en la página.

## Flujo fuera de la landing (pantallas que tocan el copy)

1. **Ficha de producto:**
   - foto 4:5;
   - precio detalle, precio por mayor con mínimo y unidad, y oferta si hay;
   - disponibilidad;
   - descripción, donde van las tallas si aplican (no hay variantes en F1);
   - "Míralo en TikTok";
   - selector de cantidad y "Agregar".
2. **Carrito:**
   - líneas con precio aplicado (detalle, x mayor u oferta);
   - aviso "Te faltan N para precio por mayor";
   - total;
   - los 3 pasos en miniatura.
3. **Datos del pedido:**
   - nombre, teléfono y retiro/envío;
   - si es envío: región, comuna y agencia;
   - comentario.
   - Nada más: sin RUT y sin cuenta (C-3, OP "Qué NO").
4. **Confirmación:**
   - número DC-xxxx;
   - botón "Enviar pedido por WhatsApp";
   - respaldo "Copiar mensaje" para cuando el navegador interno no abre WhatsApp;
   - bloque de canal oficial;
   - si es retiro: dirección y horario.
5. **Panel admin:** lista de pedidos por estado, detalle y cambio de estado, y edición de producto (precio, oferta, mayor, stock, activo). Su copy va en `copy-landing.md` §16.
