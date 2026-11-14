# Oportunidades — qué debe tener (y qué no) la web de Donde Carlitos

*analista-meiggs · 2026-11-14. Cada punto trae la decisión, para quién es, la evidencia y lo que falta confirmar. Las referencias [n] son las de `contexto-meiggs.md`; "C-" remite a `competencia.md`. Las deducciones van marcadas como *(inferencia)*.*

Prioridad: **P1** = tiene que estar en el demo · **P2** = conviene que esté · **P3** = mencionar en la propuesta como siguiente paso.

---

## P1 — Las cinco que más pesan

### 1. Sitio oficial como "fuente de verdad" contra las cuentas falsas
**Decisión:** una sección "Nuestras cuentas oficiales", visible en el footer y en la confirmación del pedido. Debe incluir:
- el IG verificado y el TikTok oficial, con link;
- **un solo número de WhatsApp**;
- la frase que ya usa la marca, "nunca te escribimos por interno" (reconfirmar el texto exacto con investigador-marca);
- qué hacer si alguien te escribe en nombre de la tienda.

El dominio propio, en HTTPS y siempre vivo, es parte del argumento.
**Para quién:** el cliente nuevo que llega desde TikTok y no sabe cuál de las cuentas "emporio" es la real.
**Evidencia:**
- La bio de Carlitos insiste en "pedidos solo a este número" [9].
- Aparecen cuentas con el nombre o la dirección del negocio que el dueño no reconoce como oficiales: dos páginas de Facebook y @_dondecarlitosemporio7 (ver `raw/competencia/evidencia-texto.md`).
- Otras dos marcas "emporio" publican la misma dirección (C-tabla).
- Según Norton 2026, el 47 % de las víctimas de estafas navideñas cayó por redes sociales, y el informe recomienda entrar directo a la página oficial del comercio [33].
- 3 de los 4 sitios de competidores enlazados en bios no cargan (C-conclusión 3).

**Confirmar con el dueño:** qué cuentas son suyas (FB donde.carlitos.33, FB juan.jzjd, TikTok @_dondecarlitosemporio7).

### 2. Precio siempre visible, al detalle y por mayor, con el mínimo claro
**Decisión:** cada producto muestra el precio al detalle y, si existe, el precio por mayor con su mínimo, por ejemplo "x mayor desde 3 u." `[EJEMPLO]`. Si se vende por caja, también el precio por embalaje. Además:
- un estado simple de disponibilidad: *Disponible / Quedan pocos / Agotado*;
- **no** mostrar un número exacto de stock mientras no haya inventario real.

**Para quién:** el cliente que quiere saber si el viaje vale la pena, y el emprendedor que calcula su margen.
**Evidencia:**
- AliChile reconoce en tono de chiste que le preguntan precios que ya están publicados [10].
- En el barrio se usa el doble precio "x mayor / al detalle" [10][11].
- Sale Todo usa los tres niveles con mínimo por variante [11].
- *(inferencia)* El precio visible baja las preguntas repetidas al WhatsApp.

**Confirmar:** si Carlitos vende por mayor y con qué mínimo (pendiente en CLAUDE.md). Todos los precios del demo van como `[EJEMPLO]` salvo los sacados de un post real.

### 3. Pedido armado en la web que llega como UN mensaje ordenado al WhatsApp, más el panel
**Decisión:**
- **Carrito sin registro** que genera un mensaje `wa.me` ya listo, con número de pedido, productos con código, cantidad, precio unitario y total, nombre, comuna o ciudad, y retiro o envío.
- El mismo pedido queda guardado en el **panel**, para que el equipo lo vea y le cambie el estado (nuevo, confirmado, listo para retiro, despachado).
- No reemplaza el WhatsApp: lo ordena.

**Para quién:** el cliente, que pide sin ir y venir de mensajes, y el equipo de 47 personas en Santiago [18], que hoy recibe todo por un solo número [9].
**Evidencia:**
- Sale Todo ya funciona así y AliChile lo hace con Whataform (C-3, C-1).
- WhatsApp tiene un carrito nativo, y ese es el estándar mínimo [40].
- El diferencial del demo es la marca propia, el panel y el estado del pedido *(inferencia)*.

### 4. "Pide antes, retira sin fila", sobre todo en temporada alta
**Decisión:**
- La opción de retiro en tienda es la de serie.
- Mostrar un mensaje del tipo "Te avisamos por WhatsApp cuando tu pedido esté listo".
- En la confirmación: dirección, metro, horario y que lleves el número de pedido.
- En diciembre, un banner de temporada que lo destaque ("Esta Navidad, evita la fila: pide online y retira").

**Para quién:** la familia que compra regalos y el comprador que no quiere andar con el celular en la mano entre la multitud.
**Evidencia:**
- 17 mil personas por hora antes de Navidad [8b].
- "Poco espacio" para los peatones y hurtos de celulares y billeteras [15][16].
- Accesos controlados y cierres puntuales de calles [14][28].
- Calor en las compras navideñas [29].
- Los malls ya usaron "filas virtuales" en Navidad (DF, vía buscador: https://www.df.cl/empresas/retail/centros-comerciales-extienden-horarios-de-atencion-de-cara-a-compras).
- *(inferencia)* Alameda 2432 queda en el borde, sobre la Alameda, así que retirar ahí evita meterse al centro del barrio. **Confirmar en terreno.**

**Confirmar:** si el local puede preparar pedidos para retiro y con cuánta anticipación.

### 5. Envíos a regiones pensados para el revendedor
**Decisión:** un bloque "Envíos a todo Chile" que explique en tres pasos:
1. pides en la web;
2. confirmas el pago;
3. se despacha por la agencia que elijas, "por pagar" o con el costo que defina el dueño.

En el checkout, un selector de agencias (Starken, Varmontt, Pullman, etc. `[EJEMPLO]`). Si hay monto mínimo, se dice antes del carrito.
**Para quién:** el emprendedor de regiones que hoy tendría que viajar.
**Evidencia:**
- Sale Todo ofrece más de 15 agencias regionales "por pagar" [11].
- AliChile pide $50.000 y AYJ $20.000 como mínimo de envío [10][30].
- Mayte Emporio tiene un local en Temuco [18].
- Los terminales de buses interregionales están al lado [6][7].

**Confirmar:** la política real de envíos y los montos (CLAUDE.md la marca como por confirmar).

---

## P2 — Conviene que estén

### 6. Sección de temporada que se activa sola (Navidad primero)
**Decisión:** una portada con un bloque de temporada configurable desde el panel. Calendario: Escolar (feb–mar), Fiestas Patrias (sep), Halloween (oct), Navidad (nov–dic), Invierno / Verano. **El demo debería lanzarse con "Navidad" activa**, porque se presenta en noviembre.
**Evidencia:**
- Las fiestas serían ~27 % de las ventas del barrio [12].
- La temporada escolar es otro peak [13].
- Sale Todo tiene categorías de Navidad, Fiestas Patrias, Invierno y Verano [11].
- Tiendanube recomienda planificar compras por temporada [5].

### 7. "Recién llegado": la novedad que se acaba
**Decisión:**
- Un carrusel "Llegó al local" con fecha de llegada.
- Una etiqueta "Nuevo" que se quita sola a los X días `[EJEMPLO]`.
- Si un producto se agota, queda visible como "Agotado" en vez de desaparecer, para que el cliente pregunte por la reposición.

**Evidencia:**
- Rotación rápida por política del fundador [18].
- "Se acaban rápido" [10].
- "Nuevos productos" es la categoría más grande de Sale Todo [11].

### 8. Ver el producto en uso: link al video de TikTok
**Decisión:** en la ficha, una miniatura con el botón "Míralo en TikTok", que lleva al video real del producto. **No** embeber reproductores con autoplay: pesan y el CSP solo debe permitir lo necesario (CLAUDE.md).
**Evidencia:** según IPSOS para TikTok, el 44 % busca ver el producto en uso y el 42 % considera decisivos los videos de las marcas [25]. Carlitos ya tiene ese contenido grabado *(inferencia; ver investigador-marca)*.

### 9. Primera pantalla igual al caption que ya funciona
**Decisión:** arriba de todo, en el celular: dirección, el metro a 5 minutos, el horario de hoy (abierto o cerrado) y los botones "Ver catálogo" y "WhatsApp". Más abajo, el mapa y cómo llegar desde los terminales de buses.
**Evidencia:**
- Es la fórmula de todos los captions del barrio [10].
- Sale Todo avisa "estamos cerrados" fuera de horario [11].
- Hubo cierres puntuales de calles [28], así que conviene poder publicar un aviso.

---

## P3 — Para la propuesta (siguiente paso)

### 10. Zona "Para tu emprendimiento"
Una sección para compra por mayor, con precio por volumen y, si el dueño quiere, un **"precio de reventa sugerido"**, como hace Protrade (C-referente 1).
**Evidencia:** las cuentas del barrio que le hablan al emprendedor ("Emprende por mayor", "hogar y emprendimiento") tienen entre 105 y 113 mil seguidores [30][31]. **Depende de confirmar que venden por mayor.**

### 11. Medios de pago y verificación, explicados
Una página corta: qué medios aceptan (**por confirmar**: efectivo, transferencia, débito), que el pedido se despacha solo con el pago verificado y a qué titular se transfiere (solo si el dueño lo aprueba).
**Evidencia:** la transferencia es el medio más aceptado en comercios chilenos [26], y AliChile recibió un comprobante falso [10].

### 12. Catálogo listo para marketplaces y TikTok Shop
Guardar los productos con datos estructurados (código, precios por nivel, categoría, temporada, foto, link al video), para reutilizarlos si llega TikTok Shop o si se publica en Falabella o Mercado Libre.
**Evidencia:** TikTok Shop no tiene fecha en Chile [24]; hay competidores que ya venden en Falabella o Mercado Libre [35][36].

---

## Qué NO debe tener el demo (y por qué)

| No | Por qué |
|---|---|
| **Pago online con tarjeta / pasarela** | Ningún competidor de la muestra la usa. El pago dominante es la transferencia o el pago en local [26][10]. Agrega fricción y riesgo para un demo. Se deja como fase 2 *(inferencia)*. |
| **Registro o login obligatorio para comprar** | Quien llega desde el navegador interno de TikTok o IG abandona si le piden cuenta *(inferencia)*. Sale Todo pide RUT incluso al detalle (C-3); hay que pedir lo mínimo. |
| **Stock numérico exacto** | Sin sistema de inventario sería inventado. Basta con Disponible / Quedan pocos / Agotado. |
| **Reseñas, testimonios o "X clientes felices"** | Prohibido inventarlos (CLAUDE.md). Si se quieren, que vengan de comentarios reales, con fuente. |
| **Cargar el catálogo completo de una vez** | Sale Todo pesa ~9,5 MB; en datos móviles dentro de TikTok se siente. Hay que paginar y comprimir imágenes. |
| **Videos con autoplay embebidos** | Peso y CSP. Mejor una miniatura con link. |
| **Chatbot o IA respondiendo por la tienda** | El valor de la marca es "un solo número oficial". Un bot nuevo puede parecer una suplantación *(inferencia)*. |
| **Prometer TikTok Shop o integraciones que no existen** | No tiene fecha en Chile [24]. |
| **Publicar datos bancarios en la web sin aprobación del dueño** | Es una decisión del negocio y puede usarse para estafas si cambian *(inferencia)*. |
| **Afirmar relación con Andres Emporio u otras cuentas** | No está confirmada (instrucción del orquestador). |

---

## Pendientes que condicionan estas decisiones
1. ¿Venden por mayor? ¿Desde cuántas unidades? (afecta 2, 10)
2. Política de envíos: agencias, "por pagar" o costo fijo, mínimo (afecta 5)
3. ¿Pueden preparar pedidos para retiro? ¿Con cuánta anticipación? (afecta 4)
4. Medios de pago reales (afecta 11)
5. ¿Qué cuentas son oficiales? (FB donde.carlitos.33, FB juan.jzjd, TikTok @_dondecarlitosemporio7) (afecta 1)
6. ¿El local de Abate Molina 77 sigue activo? (afecta 4, 9)
7. Capturas de la competencia con navegador (`research/raw/competencia/`)
