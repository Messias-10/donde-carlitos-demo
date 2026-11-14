# Fases y alcance — Donde Carlitos Emporio

*estratega-propuesta · 2026-11-14 · enfoque "Tu WhatsApp sigue igual, pero cada pedido queda ordenado"*

> **Precios y plazos:** todos van como `[COMPLETAR LUCKAS]`. No se estiman acá.
> **Costos de terceros** (dominio, servidor, Meta): no se inventan montos; se cotizan al momento.

## Resumen

| Fase | Qué es | Estado | Precio | Plazo |
|---|---|---|---|---|
| **1. Demo** | Landing con dos botones de WhatsApp, vitrina y seguimiento + panel de pedidos del WhatsApp, corriendo en el equipo de Luckas | En construcción para la presentación | Sin costo para el dueño `[COMPLETAR LUCKAS: confirmar]` | — |
| **2. Salida a producción** | Lo mismo, en internet, con sus números, sus usuarios, su dominio y respaldos | Propuesta | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **2B. WhatsApp oficial (opcional)** | Avisos automáticos y chats dentro del panel, con la API oficial de WhatsApp Business | Opcional; depende de Meta | `[COMPLETAR LUCKAS]` + costo de Meta por mensaje | `[COMPLETAR LUCKAS]` |
| **3. Extras (opcional)** | Pagos en línea, seguimiento automático con las agencias, catálogo completo y más | Se cotiza por ítem | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Costos recurrentes** | Dominio, servidor, base de datos y respaldos | Los paga el dueño directo o vía Luckas `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` | — |
| **Mantención / soporte** | Actualizaciones de seguridad, ayuda y cambios chicos | Opcional | `[COMPLETAR LUCKAS]` | — |

---

## Fase 1 — El demo

### Qué hace
**Landing pública (mobile-first, negro y oro, voz "¡Llegó…!"):**
- Dos botones grandes: **"Pedidos Santiago / retiro en tienda"** y **"Pedidos a regiones"**. Cada uno abre su WhatsApp con un mensaje ya escrito (`mensajes-whatsapp.md` a, b).
- **"Lo nuevo de los videos":** de 6 a 8 productos con precio, fecha del precio y link al video. Cada uno tiene "Pedir este por WhatsApp", que primero pregunta Santiago o Regiones y después abre el chat con el producto en el mensaje (c).
- **"Cómo comprar desde regiones"** en 5 pasos, con los datos dudosos marcados "por confirmar".
- **Retiro en Santiago, cómo llegar, horario y medios de pago en el local.**
- **Canal oficial (antiestafa):** números y cuentas oficiales, y "no publicamos datos bancarios".
- **Preguntas frecuentes** con las preguntas reales de los clientes.

**Página pública de seguimiento** (`/seguimiento/{token}`):
- muestra el estado, la línea de tiempo, la agencia + n.º de seguimiento y los productos;
- no muestra datos personales.

**Panel privado:**
- **Usuarios con nombre** y rol (atención / administrador). El historial dice quién registró, quién cambió el estado y quién despachó.
- **Registro rápido** (meta de diseño: menos de 30 s):
  - canal;
  - teléfono pegado del chat, que reconoce a los clientes conocidos;
  - nombre;
  - productos del catálogo o línea libre;
  - total;
  - comuna;
  - agencia.
- **Estados:** Cotizado → Esperando pago → Pagado → Preparando → Despachado (agencia + n.º obligatorios) o Listo para retiro → Entregado; más Cancelado.
- **Pantalla "Hoy":**
  - esperando pago > 24 h;
  - pagados sin despachar;
  - listos para retiro;
  - cotizados del día.
- **Mensajes listos con un toque:** cotización, datos de transferencia, pago recibido, despachado, listo para retiro, cancelado, recordatorio de pago y abrir chat. Abren WhatsApp con el texto puesto; **la persona aprieta enviar**.
- **Buscar** por teléfono, número de pedido o nombre.
- **Resumen** del día y de la semana, y **exportar a Excel (CSV)**.
- **Vitrina y catálogo interno** editables (precio, foto, link del video, orden).
- **Ajustes:**
  - datos de transferencia (vacíos por defecto; los escribe el dueño);
  - dirección y horario de retiro;
  - agencias.

**Seguridad base:**
- contraseñas con bcrypt;
- cookie de sesión httpOnly + SameSite=Strict;
- verificación de Origin;
- límite de intentos;
- SQL con parámetros;
- cabeceras;
- logs con teléfono enmascarado.
- Detalle en `seguridad/`.

### Qué NO hace (y se dice así en la reunión)
- **No envía mensajes solo** ni lee los chats. Solo abre WhatsApp con el texto escrito. Lo automático es la fase 2B.
- **No verifica pagos con el banco.** Recuerda revisar la cuenta antes de marcar "Pagado".
- **No cobra en línea** ni tiene carrito ni catálogo completo. La venta sigue en el local y por WhatsApp.
- **No se conecta con las agencias:** el n.º de seguimiento lo escribe quien despacha.
- **No está en internet:** corre en el equipo de Luckas con Docker.
- **Datos de demo:**
  - los pedidos y clientes son ficticios;
  - los precios son los de los videos de nov 2026, con fecha;
  - las fotos y videos son de sus redes, con permiso pendiente;
  - el número de regiones y los datos de transferencia están vacíos hasta que el dueño los dé.

---

## Fase 2 — Salida a producción

**Objetivo:** que el link de la bio lleve a la landing real, y que quien atiende los WhatsApp use el panel todos los días.

**Requisito para partir:** que el dueño confirme los datos de §2.1 y que seguridad y QA estén sin hallazgos altos.

### 2.1 Datos que confirma el dueño (sin esto no se publica)
1. **Números:** el de Santiago/retiro y el de regiones (¿son distintos?). ¿Quién atiende cada uno?
2. **Envíos a regiones:**
   - agencias;
   - si es por pagar o con costo;
   - compra mínima;
   - plazo de despacho.
3. **Pagos:** medio de pago para regiones. **Datos de transferencia: los escribe el dueño en Ajustes**; Luckas no los recibe por chat.
4. **Retiro:**
   - ¿se dejan pedidos listos?
   - ¿cuántos días se guardan?
   - ¿en cuál local?
5. **Horario** vigente y de temporada.
6. **Cancelaciones y devoluciones**, sobre todo de pedidos pagados.
7. **Nombre de la web** y qué **cuentas son oficiales**.
8. **Promo** "yo quiero mi promo": cómo se canjea.
9. **Permiso** para logo, fotos y videos, y el logo original en alta resolución.
10. **Usuarios:** quién entra al panel y con qué rol.
11. **Las respuestas de quien atiende regiones** a las 5 preguntas (`pitch.md` §d), para ajustar el panel antes de publicar.

### 2.2 Qué incluye

| Ítem | Detalle | Precio | Plazo |
|---|---|---|---|
| **Ajustes con quien atiende** | Una sesión de 30–60 min con quien atiende regiones usando el demo con pedidos parecidos a los reales. Se corrige lo que le complique antes de publicar | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Dominio** | Un `.cl` en NIC Chile **a nombre del dueño**, con renovación automática. Así los links de seguimiento y la landing son claramente oficiales | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Hosting de la web y la API** | Las dos bajo el mismo dominio (mismo origen, CORS cerrado), con HTTPS, HSTS, cookie `Secure` y la CSP final. Opciones en §2.3 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **PostgreSQL** | Se cambia SQLite por PostgreSQL con otra implementación de la interfaz `Store`, sin reescribir handlers (CLAUDE.md, Stack) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Respaldos** | Diarios, con copia fuera del servidor. Retención `[COMPLETAR LUCKAS: N días]` y prueba de restauración antes de publicar. Los pedidos son la libreta del negocio: perderlos es peor que no tenerlos | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Usuarios reales y capacitación** | Un usuario por persona. Capacitación en el local con quienes atienden los WhatsApp: registrar, cambiar estados, mandar mensajes desde el WhatsApp correcto (ojo con los celulares que tienen dos WhatsApp), revisar "Hoy" y exportar. Guía de una página con capturas | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Catálogo interno inicial** | Carga de `[COMPLETAR LUCKAS: N]` productos frecuentes, para que el registro sea rápido. Se parte con lo de los videos recientes | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Fotos reales para la vitrina** | Guía de una página para sacarlas con el celular (4:5, producto en mano), o una sesión de fotos | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Privacidad** | El panel guarda nombre, teléfono y comuna de los clientes. Se necesita un aviso de privacidad corto (qué se guarda, para qué y cuánto tiempo), un plazo de borrado y el acceso por roles. La **Ley 21.719** de protección de datos tiene fecha de vigencia el **1 de diciembre de 2026** según el Diario Oficial, pero hay un proyecto del Ejecutivo para postergarla a diciembre de 2027 `[VERIFICAR estado antes de publicar]` | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Monitoreo** | Aviso si la web o la API se caen; Dependabot, `npm audit` y `govulncheck` | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Lanzamiento** | Link en las bios de TikTok e IG. QR impreso en caja: "¿Eres de región? Escanea y escríbenos al WhatsApp correcto". Un video de ellos presentándolo, si quieren | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Medición** | Clics por botón (Santiago/Regiones/"Pedir este") con analítica sin datos personales. Métricas en `propuesta-valor.md` §6 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |

### 2.3 Hosting: VPS o PaaS

| | **VPS** (Docker Compose) | **PaaS** (plataforma + PostgreSQL administrado) |
|---|---|---|
| **Ventaja** | Usa el mismo `docker-compose.yml` del demo. Costo fijo | Menos mantención: sistema operativo, certificados y respaldos de la base a cargo del proveedor |
| **Desventaja** | Luckas mantiene el servidor (actualizaciones, firewall, respaldos) | Costo variable; más dependencia del proveedor |
| **Conviene si** | Hay contrato de mantención | El dueño quiere depender lo menos posible de una persona |
| **Costo mensual** | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |

**Recomendación:** `[COMPLETAR LUCKAS]`. Criterio: en diciembre el panel no se puede caer, porque es la libreta de pedidos. Hay que elegir lo que Luckas pueda vigilar mejor.

---

## Fase 2B — WhatsApp oficial (opcional, no está en el demo)

**Qué agrega:**
- **Avisos automáticos** al cliente cuando cambia el estado ("pago recibido", "despachado con n.º…"), sin que nadie apriete enviar.
- **Los chats dentro del panel:** varias personas atienden el mismo número y cada conversación queda ligada a su pedido.
- **Respuestas automáticas** a lo repetido (dirección, horario, "¿hacen envíos?"), y derivación a una persona.

**Qué requiere** (verificar todo al cotizar; las reglas de Meta cambian seguido):
- **Cuenta de Meta Business verificada** y aprobación del nombre visible.
- **Un número registrado en la API.** `[VERIFICAR]` si ese número puede seguir usándose al mismo tiempo en la app WhatsApp Business del celular (la llamada "coexistencia"), y si eso está disponible en Chile y con qué proveedor. Si no, conviene usar un número nuevo para lo automático.
- **Plantillas aprobadas por Meta** para los mensajes que el negocio inicia (los de §6 de `mensajes-whatsapp.md` son la base).
- **Un proveedor** (BSP) o la API directa de Meta (Cloud API).

**Costo** (sin montos; se cotiza al momento):
- **Meta cobra por mensaje de plantilla entregado** desde el 1 de julio de 2025. Antes cobraba por conversación.
- El precio depende de la **categoría** (marketing, utilidad o autenticación) y del **país** del cliente.
- Según la documentación y varios proveedores, los mensajes de **utilidad** enviados dentro de la ventana de 24 h de atención no se cobran. Hay versiones distintas sobre cambios en los mensajes de servicio desde octubre de 2026: **`[VERIFICAR en la página oficial de precios de Meta antes de cotizar]`**.
- Si se usa un proveedor (BSP), se suma su tarifa.

**Precio y plazo:** `[COMPLETAR LUCKAS]` + costo de Meta y del proveedor.

**Recomendación:** hacerla **después** de 1–2 meses de fase 2. Así se automatiza lo que de verdad se repite, con datos del panel, y no a ciegas.

---

## Fase 3 — Extras opcionales (se cotizan por separado)

| Ítem | Qué resuelve | Evidencia | Precio | Plazo |
|---|---|---|---|---|
| **Pagos en línea** (link de pago) | Pagar con tarjeta o débito desde regiones; responde a "¿se puede pagar por link de crédito?" | `faq-clientes.md` §9. Nadie del barrio en la muestra lo usa (`oportunidades.md`, "Qué NO") | `[COMPLETAR LUCKAS]` + comisión de la pasarela | `[COMPLETAR LUCKAS]` |
| **Seguimiento automático con agencias** | Que el estado "Entregado" se marque solo con la información de la agencia, si la agencia lo permite | Reclamos de cajas (TT-7706913997378948370) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Catálogo completo con pedido web** | Volver a lo que se sacó (carrito → WhatsApp), solo si los datos de la fase 2 muestran que la vitrina se queda corta | Competencia con carrito (C-1) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Canal mayorista** | Lista de precios por volumen y mínimo para revendedores | `oportunidades.md` P3-10 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Preventa con reserva** | "Llega el 16": reservar con número de pedido, en vez de "en el comentario" | tr. TT-7706993060130950408 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **"Avísame cuando llegue"** | El cliente deja su número en un producto agotado de la vitrina | `faq-clientes.md` §3 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Multi-local (Temuco)** | Pedidos y retiro por local | `adn-marca.md`, línea de tiempo | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Promo digital "Yo quiero mi promo"** | Código con vigencia que se canjea en caja | `ofertas.md`, regla 4 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |

---

## Riesgos y cómo se manejan

| Riesgo | Manejo |
|---|---|
| Quien atiende no registra los pedidos ("no hay tiempo") | Se valida con ella antes de publicar (§2.1 punto 11). Registro de 30 s, mensajes listos que le ahorran tiempo y "Hoy" como su lista de tareas. Se revisa la adopción a las 2 semanas |
| Mensaje enviado desde el WhatsApp equivocado | La etiqueta "Envíalo desde el WhatsApp de {canal}" y la capacitación |
| Comprobante falso | Confirmación obligatoria al marcar "Pagado". Solo ciertos roles pueden marcarlo `[CONFIRMAR CON EL DUEÑO]` |
| Datos de transferencia expuestos | Solo en Ajustes (rol admin) y en los mensajes; nunca en la web, en variables `NEXT_PUBLIC_*` ni en el repo |
| Fuga de datos de clientes | Roles, export solo para admin con aviso, link de seguimiento con token y sin datos personales, respaldos cifrados `[COMPLETAR LUCKAS]`, Ley 21.719 (§2.2) |
| Suplantación (alguien copia la landing) | Dominio a nombre del dueño, el link oficial en las bios y la sección "Canal oficial" |
| El panel se cae en diciembre | Monitoreo, respaldos y prueba de carga simple antes de la temporada |
| Dependencia de Luckas | Dominio y cuentas a nombre del dueño, documentación y mantención opcional |
