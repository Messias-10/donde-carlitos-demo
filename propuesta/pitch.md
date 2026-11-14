# Pitch al dueño — Donde Carlitos Emporio

*estratega-propuesta · Fase 2 · 2026-11-14*

**Tono:** cercano y directo, como se habla en el barrio. Se usa "tú". Si el dueño trata de usted, Luckas cambia a usted.

**Regla de oro:** solo se muestra lo que el demo hace. Todo lo demás va como "eso lo vemos en la fase 2".

---

## (a) Mensaje de primer contacto

**Canal recomendado:** DM de Instagram a @mayte_emporioo, o ir directo al local.
- **No** escribir al +56 9 2114 5368: es el número de pedidos y ya está saturado (`faq-clientes.md`, quejas). Escribir ahí sería sumarle carga justo a lo que se quiere aliviar.
- Un desconocido que manda un link a una marca que avisa por estafas puede parecer phishing. Por eso, **mejor adjuntar un video vertical de 60 s** grabado de la pantalla del demo, en vez de un link. El demo no se despliega a internet (CLAUDE.md, Autonomía).

```
Hola, ¿cómo están? Soy Luckas, desarrollador web de Santiago, y sigo sus videos.
En los comentarios se repite lo mismo: "¿dónde están?", "¿hacen envíos?", "¿te queda?".
Les armé una muestra de cómo sería la tienda online de Donde Carlitos: catálogo con precio detalle y por mayor, y cada pedido les llega con número y armado a su WhatsApp.
Acá la pueden ver en el celular: [VIDEO DE 60 S / LINK DEMO]
Si les tinca, paso por Alameda 2432 y se la muestro en 5 minutos, sin compromiso.
No les voy a pedir datos, claves ni pagos por este medio.
```

> La última línea va a propósito: le habla a una marca que vive alerta con las estafas.

---

## (b) Guion de 2 minutos en el local

**Preparación:**
- **Dónde corre el demo:** en el notebook de Luckas, en vista móvil (las herramientas de desarrollador del navegador, a 390 px). Los puertos solo escuchan en `127.0.0.1` (CLAUDE.md, Docker), así que **un celular no puede abrirlo por la red**, y esa regla no se cambia. Plan B: el mismo video vertical del DM, en el celular.
- **Antes de llegar:**
  - catálogo cargado y panel abierto en otra pestaña;
  - el número de WhatsApp de destino apuntando al de Luckas, no al del negocio (`mensaje-pedido-whatsapp.md` §1). Si el dueño quiere ver cómo le llega a él, se cambia en el momento, con su permiso.
- **Mejor momento:** temprano, antes de la hora punta. Evitar los sábados de temporada.

| Tiempo | Qué se dice | Qué se muestra |
|---|---|---|
| 0:00–0:15 | "Hola, soy Luckas, el que les escribió por Instagram. Te muestro algo en dos minutos y me voy, no te quito más." | — |
| 0:15–0:35 | "Tus videos venden solos. Pero en los comentarios la gente pregunta siempre lo mismo: dónde están, si hacen envíos, si queda, cuántos trae el embalaje. Y algunos dicen que escribieron al WhatsApp y nadie les contestó. Con todo lo que publican, no hay cómo contestar todo a mano." | Puede mostrar 2 o 3 comentarios reales en pantalla, sin nombres de usuario |
| 0:35–1:15 | "Mira: alguien ve tu video del squishy y toca el link de la bio. Cae acá. Ve el precio al detalle, $4.000, y por caja de 12, $3.500. Le pone 12 y el precio por mayor se aplica solo. Pone su nombre, si retira o es de región, y listo: pedido DC-0007. Se le abre WhatsApp con todo escrito." | Hero → tarjeta del squishy → carrito con precio por mayor → datos → confirmación → mensaje armado |
| 1:15–1:40 | "Y esto lo ves tú, en el celular. Llegó el pedido, lo marcas confirmado, después listo para retiro. Si se acabó el tren navideño, lo pones agotado y la web deja de venderlo. Si cambias un precio, es un minuto." | Panel: lista → detalle → cambiar estado → editar el stock de un producto |
| 1:40–1:50 | "Y aquí dice clarito cuál es tu único número y tus cuentas oficiales. Si alguien usa tu nombre, la gente tiene dónde verificar." | Sección "Canal oficial" |
| 1:50–2:00 | "Esto ya funciona. Para dejarlo andando con tu dominio, tus fotos y tus precios reales faltan unos datos que solo tú sabes. ¿Me das 10 minutos esta semana para verlos? Te traigo la propuesta con precio." | — |

**Las preguntas para esos 10 minutos** (salen de `research/pendientes.md`):
1. **Envíos:** ¿despachan a regiones? ¿Por qué agencias? ¿El envío es "por pagar"? ¿Hay monto mínimo?
2. **Horario:** ¿cuál es el vigente? ¿Hay horario de temporada?
3. **Por mayor:** ¿hay un mínimo general? ¿El precio que se publica en la web se respeta en caja?
4. **Promo:** ¿cómo se canjea "yo quiero mi promo"? ¿Quieren que la web la aplique?
5. **Cuentas y nombre:** ¿qué cuentas son suyas? ¿La web se llama "Donde Carlitos Emporio", "Mayte Emporio" o "Familia Emporio"?
6. **Extras:**
   - medios de pago por local;
   - si se incluye el local de Temuco;
   - política de cambios;
   - si pueden preparar pedidos para retiro;
   - el logo original en alta resolución.

**Qué NO decir:**
- **Porcentajes ni promesas de ventas.** Ni "vas a vender el doble", ni ninguna cifra parecida.
- **Nada contra la competencia ni contra el local** (música, limpieza, guardias): son quejas de comentarios, no tema del pitch.
- **Nada sobre su historia personal** (prensa) ni el nombre "Carlitos", si él no lo menciona primero.
- **"Envíos a todo Chile", "pagos online" o "avisos automáticos"** como si el demo ya los tuviera.

---

## (c) Objeciones probables

**1. "Ya vendo bien por TikTok, no necesito web."**
> "Justamente por eso. La web no reemplaza a TikTok: es donde cae la gente que tus videos ya traen. Hoy esa gente pregunta en los comentarios y te escribe uno por uno. Con la web llega con el pedido hecho, con número y total. TikTok sigue siendo la vitrina; la web ordena la caja. Y entre las cuentas grandes del barrio, tú eres la más grande y no tienes catálogo propio."
- Evidencia: `competencia.md`, conclusión 4; `faq-clientes.md`.

**2. "No tengo tiempo de subir productos."**
> "No hay que subir todo. Partimos con lo que muestras en los videos de la semana. El demo ya tiene productos sacados de tus videos de noviembre, con el precio que dijiste en cada uno. Para la salida, yo cargo el catálogo inicial (fase 2). Después, cambiar un precio, marcar agotado o poner algo en oferta se hace desde el celular, en un minuto, y lo puede hacer la misma persona que hoy atiende el WhatsApp."
- Ojo: no prometer que el dueño puede **crear** productos con foto desde el panel si el demo no lo hace. Verificar el alcance con ingeniero-go. Si no está, es fase 2 (`fases.md`).

**3. "¿Cuánto cuesta?"**
> "Depende de cuánto quieras dejar andando. Lo dividí en tres partes: la salida a producción, que es lo básico para que funcione con tu dominio, y los extras opcionales, como pagos en línea o avisos automáticos. La primera parte cuesta [COMPLETAR LUCKAS] y queda lista en [COMPLETAR LUCKAS]. Aparte van el dominio y el servidor, que son costos mensuales o anuales chicos y te los detallo por escrito: [COMPLETAR LUCKAS]."
- Detalle en `fases.md`. Si cobra o no comisión por venta lo define Luckas: [COMPLETAR LUCKAS].

**4. "Mi gente prefiere venir al local" / "Nosotros vendemos presencial."**
> "Perfecto, el retiro en tienda es la opción principal. El cliente llega sabiendo qué lleva y cuánto paga, y tú sabes antes qué viene a buscar. Si no quieres despachar, se oculta la opción de envío y queda solo retiro. Y en Navidad, que el barrio se repleta, sirve que la gente llegue con el pedido listo, si ustedes pueden tenerlo preparado."
- Evidencia: `contexto-meiggs.md` [8b] y [15]. Que puedan tener los pedidos preparados se pregunta; no se promete.

**5. "Me van a llegar más mensajes y no doy abasto." / "Ya tengo el catálogo de WhatsApp."**
> "Te escriben los mismos clientes de siempre, pero con el pedido completo en un solo mensaje: 'DC-0012, tres productos, $34.600, retira'. Eso se confirma mucho más rápido que una conversación de veinte mensajes. Y si el cliente no manda el WhatsApp, el pedido igual queda guardado en tu panel. El catálogo de WhatsApp lo reemplaza esta web: se actualiza desde el panel y la ve cualquiera desde el link de tu bio."
- Evidencia: el comentario "no esta actualizado el catalogo…" (TT-7707471893363789063).

**6. "¿Y si la web se cae o nadie la mantiene?"**
> "Buena pregunta. De cuatro webs de otras tiendas del barrio que revisamos, tres ni siquiera cargan. Por eso la fase 2 incluye el dominio a tu nombre, respaldos automáticos y [COMPLETAR LUCKAS: soporte o mantención]. La web tiene que estar siempre viva; si no, resta en vez de sumar."
- Evidencia: `competencia.md`, conclusión 3.

**7. "¿Y la promo de 'yo quiero mi promo'?"**
> "Tú decides. La web puede mostrar el precio promo directo en el producto, con el botón '¡Yo quiero mi promo!', o dejarlo solo para los comentarios. Lo que me tienes que decir es cómo la canjean hoy en caja, para que la web diga lo mismo."
- Evidencia: `pendientes.md` §3.
