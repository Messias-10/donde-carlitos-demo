# Fases y alcance — Donde Carlitos Emporio

*estratega-propuesta · Fase 2 · 2026-11-14*

> **Precios y plazos:** todos van como `[COMPLETAR LUCKAS]`. No se estiman acá.

## Resumen

| Fase | Qué es | Estado | Precio | Plazo |
|---|---|---|---|---|
| **1. Demo** | Lo que el dueño ve hoy: web + pedidos + panel, corriendo en el equipo de Luckas | Construido para la presentación | Sin costo para el dueño `[COMPLETAR LUCKAS: confirmar]` | — |
| **2. Salida a producción** | La misma tienda, en internet, con su dominio, sus datos reales y respaldos | Propuesta | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **3. Extras (opcional)** | Pagos en línea, canal mayorista, notificaciones y más | Opcional, se cotiza por separado | `[COMPLETAR LUCKAS]` por ítem | `[COMPLETAR LUCKAS]` |
| **Costos recurrentes** | Dominio, servidor, base de datos y respaldos | Los paga el dueño directo o vía Luckas `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` mensual / anual | — |
| **Mantención / soporte** | Actualizaciones de seguridad, ayuda y cambios chicos | Opcional | `[COMPLETAR LUCKAS]` | — |

---

## Fase 1 — El demo (lo que ya existe)

**Qué hace:**
- **Web mobile-first** con la marca: negro y oro, la voz de "¡Llegó…!" y el logo tal cual.
- **Catálogo** por categorías, con buscador.
- **Precios:** detalle, oferta o remate (con tachado solo si el precio anterior es real) y por mayor con mínimo y unidad ("caja de 12").
- **Badges:** oferta/remate, recién llegado y temporada (Navidad activa).
- **Disponibilidad:** Disponible / Quedan pocas unidades / Agotado.
- **"Lo viste en TikTok":** miniaturas que llevan al video real.
- **Carrito sin registro**, que aplica el precio por mayor al llegar al mínimo.
- **Pedido** con nombre, teléfono y retiro o envío (región, comuna y agencia elegida por el cliente). Se **guarda en la base con número correlativo DC-0001** y abre WhatsApp con el mensaje armado (`mensaje-pedido-whatsapp.md`).
- **Panel con login:**
  - ver pedidos por estado y cambiarlos: nuevo → confirmado → listo / despachado → entregado, o cancelado;
  - el historial de cambios;
  - editar precio, oferta, precio por mayor, stock y visibilidad desde el celular.
- **Secciones informativas:** "Cómo llegar", "Canal oficial" (antiestafa) y preguntas frecuentes con las preguntas reales de los clientes.
- **Seguridad base:**
  - contraseña con bcrypt;
  - cookie de sesión protegida;
  - límite de intentos y de pedidos;
  - cabeceras de seguridad;
  - SQL con parámetros.
  - Detalle en `seguridad/`.

**Qué NO es (y se dice así en la reunión):**
- **No está en internet:** corre en el equipo de Luckas con Docker.
- **No cobra:** el pago se coordina por WhatsApp o en el local.
- **No avisa solo al cliente** cuando cambia el estado del pedido.
- **Los datos del demo no son definitivos:**
  - los precios son referenciales, sacados de videos con fecha;
  - las fotos y videos son de sus redes, con permiso pendiente;
  - los pedidos son ficticios.
- **Los datos marcados "por confirmar" no están confirmados** (horario, envíos, pagos, promo).
- **No maneja tallas ni colores como variantes:** van en la descripción.
- **Textos fijos:** el horario y los avisos no se editan desde el panel.
- **Alta de productos con foto desde el panel:** `[VERIFICAR con ingeniero-go si el demo la incluye; si no, es fase 2]`.

---

## Fase 2 — Salida a producción

**Objetivo:** que el link de la bio de @mayte_emporio y @mayte_emporioo lleve a la tienda real, segura y siempre viva.

**Requisito para partir:** que el dueño confirme los datos del §2.1 y que seguridad y QA estén sin hallazgos altos.

### 2.1 Datos que confirma el dueño (sin esto no se publica)
De `research/pendientes.md`:
1. **Envíos:**
   - ¿despachan?
   - ¿por qué agencias?
   - ¿"por pagar" o con costo?
   - ¿monto mínimo?
2. **Horario** vigente y de temporada.
3. **Medios de pago** por local.
4. **Mínimo** por mayor.
5. **Canje** de la promo.
6. **Política de cambios.**
7. **Nombre de la web** y qué **cuentas son oficiales**.
8. **Local de Temuco:** ¿se incluye?
9. **Respeto del precio publicado** en el pedido.
10. **Pedidos para retiro:** ¿se pueden preparar y con cuánta anticipación?
11. **Permiso** para usar logo, fotos y videos. Logo original en alta resolución.

### 2.2 Qué incluye

| Ítem | Detalle | Precio | Plazo |
|---|---|---|---|
| **Dominio** | Un `.cl` registrado en NIC Chile **a nombre del dueño**, con renovación automática y aviso de vencimiento. Así se evita la web que "no resuelve", como les pasa a 3 de 4 sitios de la competencia revisados (`competencia.md`, conclusión 3) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Hosting de la API y la web** | Hay dos opciones (ver la tabla 2.3). Web y API bajo el mismo dominio (mismo origen, CORS cerrado), con HTTPS, HSTS, cookie `Secure` y la CSP final con el dominio | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **PostgreSQL** | Se cambia SQLite por PostgreSQL con una nueva implementación de la interfaz `Store`, sin reescribir handlers ni reglas (CLAUDE.md, Stack). Se usan las mismas migraciones adaptadas | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Respaldos** | Respaldos automáticos diarios de la base, con copia fuera del servidor. Retención: `[COMPLETAR LUCKAS: N días]`. Una prueba de restauración antes de publicar y luego cada `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Fotos reales** | Fotos de los productos en formato 4:5, en mano o en contexto, como en sus videos (estilo-carlitos §7). Hay dos opciones: una sesión de fotos, o una guía de una página para que el equipo las saque con el celular | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Carga inicial del catálogo** | `[COMPLETAR LUCKAS: N productos]`. Se parte por la temporada en curso y los productos de los videos recientes, con los precios confirmados por el dueño | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Usuarios del panel** | Un usuario por persona que atienda los pedidos, para que el historial diga quién cambió qué. La contraseña de cada uno la define esa persona | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Capacitación** | Una sesión en el local con quienes atienden el WhatsApp: ver pedidos, cambiar estados, editar precio, stock y oferta, y qué hacer si algo falla. Se deja una guía de una página con capturas | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Textos editables** | Horario, aviso del día (ej.: "hoy cerramos a las 14:00") y banner de temporada, editables desde el panel | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Alta de productos con foto desde el panel** | Solo si no viene en el demo (ver fase 1) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Medición** | Analítica web sin datos personales y un resumen simple de pedidos (por semana, retiro o envío, regiones, precio por mayor). Son las métricas de `propuesta-valor.md` §5 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Privacidad** | Aviso de privacidad corto (qué datos se piden, para qué y por cuánto tiempo se guardan). Revisar la Ley 21.719 de protección de datos personales `[VERIFICAR fecha de entrada en vigencia]` | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Lanzamiento** | Link en la bio de TikTok e IG, QR impreso para caja y vitrina ("La próxima vez, pide acá"), y un video de ellos presentando la web, si el dueño quiere | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Monitoreo** | Aviso a Luckas si la web o la API se caen, y revisión de dependencias (Dependabot, `npm audit`, `govulncheck`) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |

### 2.3 Hosting: VPS o PaaS

| | **VPS** (servidor propio con Docker Compose) | **PaaS** (plataforma administrada + PostgreSQL administrado) |
|---|---|---|
| **Ventaja** | Usa el mismo `docker-compose.yml` del demo. Costo fijo y predecible | Menos mantención: el proveedor se encarga del sistema operativo, los certificados y, a menudo, de los respaldos de la base |
| **Desventaja** | Luckas tiene que mantener el servidor: actualizaciones, firewall y respaldos | El costo puede subir con el tráfico. Depende más del proveedor |
| **Conviene si** | Hay un contrato de mantención con Luckas | El dueño quiere depender lo menos posible de una persona |
| **Costo mensual** | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |

**Recomendación:** `[COMPLETAR LUCKAS]`. Criterio sugerido: en temporada alta, la web no se puede caer, así que hay que elegir la opción que Luckas pueda vigilar mejor. Antes de diciembre conviene una prueba de carga simple.

### 2.4 Lo que la fase 2 todavía no hace
Igual que la fase 1: no cobra en línea, no notifica sola al cliente y no maneja variantes. Todo eso es fase 3.

---

## Fase 3 — Extras opcionales (se cotizan por separado)

| Ítem | Qué resuelve | Evidencia | Precio | Plazo |
|---|---|---|---|---|
| **Pago "transferencia verificada"** (paso intermedio) | Marcar en el panel "pagado / verificado" antes de despachar | Comprobantes de transferencia falsos en el barrio ([10] de `contexto-meiggs.md`) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Pagos en línea** | Pago con tarjeta o débito a través de una pasarela chilena. Le sirve a regiones y a quien pregunta por crédito ("¿se puede pagar por link de crédito?") | `faq-clientes.md` §9. Ojo: nadie del barrio en la muestra la usa (`oportunidades.md`, "Qué NO") | `[COMPLETAR LUCKAS]` + comisión de la pasarela | `[COMPLETAR LUCKAS]` |
| **Canal mayorista** | Zona "Para tu emprendimiento" con lista de precios por volumen, mínimo de compra, precio de reventa sugerido (si el dueño quiere) y descarga de catálogo | `oportunidades.md` P3-10; referente Protrade (`competencia.md`) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Notificaciones** | Al cliente, cuando cambia el estado ("tu pedido está listo"), vía WhatsApp Business API (requiere aprobación de Meta y tiene costo por conversación), correo o SMS. Al equipo, una alerta de pedido nuevo | Hoy los reclamos terminan en comentarios públicos (`dolores-pedido.md` B7) | `[COMPLETAR LUCKAS]` + costo del proveedor | `[COMPLETAR LUCKAS]` |
| **"Avísame cuando llegue"** | El cliente deja su número en un producto agotado y le llega un aviso cuando vuelve | "¿Te queda…?", "y la avejita la trajeron de nuevo?" (`faq-clientes.md` §3) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Variantes (tallas y colores)** | Uniformes escolares y textil, con stock por talla | ~15 preguntas de talla en un solo video (`faq-clientes.md` §8) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` antes de feb (vuelta a clases) |
| **Preventa con reserva** | "Llega el día 16": reservar con un número, en vez de pedir "en el comentario" | tr. TT-7706993060130950408 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Promo digital "Yo quiero mi promo"** | Código o botón con vigencia, que se registra y se canjea en caja | `ofertas.md`, regla 4 | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Multi-local (Temuco)** | Stock y retiro por local | Local de Temuco (`adn-marca.md` §Línea de tiempo) | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |
| **Catálogo listo para marketplaces** | Exportar productos a Mercado Libre o Falabella, o a TikTok Shop si llega a Chile (sin fecha al 02-11-2026) | `oportunidades.md` P3-12; [24] | `[COMPLETAR LUCKAS]` | `[COMPLETAR LUCKAS]` |

---

## Riesgos y cómo se manejan

| Riesgo | Manejo |
|---|---|
| Precio de la web distinto al del local | Una sola fuente: el panel. El pedido guarda el precio del momento. Se define la política con el dueño (fase 2, §2.1) |
| Catálogo desactualizado (el mismo problema de hoy) | Capacitación, una persona responsable `[CONFIRMAR CON EL DUEÑO]`, y "Agotado" visible en vez de borrar |
| Más pedidos de los que se pueden confirmar en temporada | El panel ordena por estado y muestra los nuevos primero. Si hace falta, se puede pausar el envío a regiones o el pedido web desde la configuración `[COMPLETAR LUCKAS: incluir en fase 2]` |
| Suplantación de la propia web | Dominio a nombre del dueño, HTTPS y un solo link oficial en las bios. Se avisa en "Canal oficial" |
| La web se cae en diciembre | Monitoreo, respaldos y una prueba de carga antes de la temporada |
| Dependencia de Luckas | Dominio y cuentas a nombre del dueño, documentación de cómo correrla y un contrato de mantención opcional |
