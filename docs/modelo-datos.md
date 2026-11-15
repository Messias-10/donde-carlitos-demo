# Modelo de datos — Donde Carlitos Emporio (v2)

*ingeniero-go · 2026-11-15 · SQLite (`modernc.org/sqlite`), migraciones en `api/migrations/` (0001 a 0005). Contrato HTTP: `docs/openapi.yaml` 2.1.0.*

## Qué cambió en 2.1.0 (migración 0005)
Se aplica sobre la base existente, sin borrar nada:
- **`productos`:** `precio_es_promo` + `nota_promo` (precio de promo por comentario, no de detalle); `en_vitrina` + `orden_vitrina` (vitrina de la landing, máx. 8).
- **`usuarios`:** `canal_por_defecto` (`santiago` | `regiones` | NULL) y `debe_cambiar_clave` (clave temporal que hay que cambiar).
- **`pedidos`:** `token_vence_en` (el link de seguimiento vence 30 días después de entregado o cancelado).
- **`items_pedido`:** `tipo_precio` admite `manual` (línea del catálogo con precio editado). La tabla se reconstruye, porque SQLite no permite cambiar un CHECK.
- **Datos de la semilla ya cargada:**
  - promos P002, P007, P008, P011 y P013;
  - vitrina P012, P010, P016, P001, P005, P020, P017 y P027;
  - unidades sin paréntesis;
  - stock 300 en los productos del guion;
  - canal de los usuarios demo.
- **Estados:** `pagado` pasa directo a `despachado` o `listo_retiro` (transición directa, sin "preparando" implícito). Un `despachado` ya no se cancela.
- **Stock no estricto por defecto** (`STOCK_ESTRICTO=false`): si no alcanza al pagar o preparar, se acepta, queda en 0 y se avisa.

## Qué cambió en v2 (cambio de enfoque)
- **No hay pedido público.** Se eliminaron `POST /api/pedidos`, sus límites por IP y por teléfono y el mensaje de pedido web. El equipo registra los pedidos en el panel mientras atiende los dos WhatsApp (Santiago/retiro y Regiones).
- `pedidos`, `items_pedido` e `historial_estados` se **recrean** con el esquema nuevo (migración 0004). Los pedidos de la v1 eran todos ficticios (demo/QA) y se borran: antes se devuelve al catálogo el stock que tenían descontado, y la semilla vuelve a cargar pedidos DEMO.
- `admins` pasa a **`usuarios`**, con nombre visible, rol (`admin` | `vendedor`) y activo.
- Nuevo: canal, forma de pago acordada, pago registrado, n.º de seguimiento, token de seguimiento público, descuento, costo de envío / envío por pagar, total escrito a mano si hay líneas sin precio, motivo de cancelación, nota interna, líneas libres, `registrado_por`. Campos alineados con `propuesta/mensajes-whatsapp.md` §4.
- El catálogo (`categorias`, `productos`) y las imágenes subidas no cambian.

**Para el demo** no hace falta `docker compose down -v`: la 0004 deja la base consistente. Si se quiere partir de cero (catálogo con el stock original), `down -v` también sirve; la semilla se recrea sola.

## Principios
- Montos en **pesos chilenos enteros**.
- **El precio se congela** en `items_pedido` al registrar o editar el pedido.
- **El servidor calcula** subtotal y total cuando todas las líneas tienen precio (`Σ − descuento + costo_envio`; si el panel manda un total, debe coincidir). Si alguna línea libre va sin precio, el total lo escribe quien atiende (`total_manual = 1`).
- Fechas-hora en `TEXT` ISO-8601 **UTC** de ancho fijo (`2026-11-14T15:00:00.000Z`). Fechas simples `AAAA-MM-DD` = día de Chile.
- **Portabilidad a PostgreSQL:** parámetros `?`, `RETURNING`, índices parciales y `ON CONFLICT`. Al migrar: `INTEGER PRIMARY KEY` → `BIGSERIAL`, `TEXT` de fechas → `TIMESTAMPTZ`/`DATE`, `?` → `$n`. La aplicación solo ve `store.Store`.

## Diagrama
```
categorias 1──< productos 1──< items_pedido >──1 pedidos 1──< historial_estados >──0..1 usuarios
                                (producto_id NULL = línea libre)   │                   usuarios 1──< sesiones
                                                                   ├── registrado_por ──> usuarios
                                                                   └── pago_verificado_por ──> usuarios
contadores (correlativo DC-xxxx)    schema_migraciones
```

## Tablas

### `categorias` y `productos` (sin cambios desde v1.1)
Ver el detalle de campos en el historial de este archivo y en `docs/openapi.yaml` (`ProductoPublico`, `ProductoAdmin`). Resumen:
- precio al detalle, oferta con `oferta_hasta`, precio por mayor por unidad con `mayor_minimo` (+ `mayor_minimo_confirmado`) y `unidad_mayor`;
- stock interno (en público solo `disponibilidad`), `umbral_pocos`, activo, temporada, `fecha_llegada`;
- `imagen_url`, `tiktok_url`, `fuente_url`/`fuente_fecha`, `es_ejemplo`, `precio_actualizado_en`;
- borrado lógico con `eliminado_en`;
- (2.1.0) `precio_es_promo` INTEGER 0/1 + `nota_promo` TEXT NULL: el precio es de promo por comentario del video, no de detalle. En la semilla: P002, P007, P008, P011 y P013, según el campo `nota` de `research/productos-semilla.json`;
- (2.1.0) `en_vitrina` INTEGER 0/1 + `orden_vitrina` INTEGER NULL 1–99: vitrina de la landing, máx. 8. Si no hay ninguno, `GET /api/vitrina` muestra los 8 más recientes;
- `unidad_mayor` sin paréntesis (la web la muestra entre paréntesis).

**Precio de una línea de catálogo:** el menor entre detalle, oferta vigente y, **solo si `cantidad ≥ mayor_minimo`**, precio por mayor. Se guarda el tipo aplicado.

### `usuarios` (antes `admins`)
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| usuario | TEXT NOT NULL UNIQUE | 3–32: minúsculas, números, `. _ -` |
| nombre | TEXT NOT NULL | nombre visible en el panel y en el historial ("Vane (regiones)") |
| password_hash | TEXT NOT NULL | bcrypt costo 12. **Nunca** se guarda la clave en claro |
| rol | TEXT NOT NULL CHECK IN ('admin','vendedor') | `vendedor`: pedidos y clientes; ve productos sin editarlos; no ve usuarios, resumen ni CSV |
| activo | INTEGER 0/1 | desactivado = no puede entrar (sus sesiones se borran) |
| canal_por_defecto | TEXT NULL | `santiago` \| `regiones`; el panel lo preselecciona al registrar |
| debe_cambiar_clave | INTEGER 0/1 | 1 para claves temporales (creadas o regeneradas en el panel); mientras sea 1, solo puede ver su sesión, cambiar la clave o salir |
| creado_en, actualizado_en | TEXT | |

- El primer admin lo crea `setup-admin` (`api/.env`: `ADMIN_USUARIO`, `ADMIN_PASSWORD_HASH_B64`). Al arrancar, la API lo sincroniza: hash del `.env`, rol `admin`, activo. Así siempre hay una vía de recuperación.
- Los demás usuarios se crean en el panel (`POST /api/admin/usuarios`). La API genera la clave y la muestra una sola vez.
- La semilla **no** crea usuarios con clave.

### `sesiones`
`token_hash` (PK, HMAC-SHA256(SESSION_SECRET, token)), `usuario_id` → usuarios ON DELETE CASCADE, `creado_en`, `expira_en` (12 h), `ultimo_uso` (vence con 2 h sin uso). Solo valen si el usuario está activo.

### `contadores`
`('pedidos', valor)`: correlativo atómico. `UPDATE … RETURNING` dentro de la transacción que crea el pedido (`BEGIN IMMEDIATE`). Sin huecos ni duplicados.

### `pedidos`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| secuencia | INTEGER UNIQUE | |
| numero | TEXT UNIQUE | `DC-0001` |
| token_seguimiento | TEXT NOT NULL UNIQUE | 16 bytes de `crypto/rand` en base64url (22 caracteres); link público `URL_SITIO/seguimiento/<token>`. Nunca el `DC-xxxx` |
| canal | TEXT CHECK IN ('santiago','regiones') | WhatsApp por el que llegó |
| cliente_nombre | TEXT NOT NULL | |
| cliente_telefono | TEXT NOT NULL | E.164 `+569XXXXXXXX`; en logs `+569****1234` |
| busqueda | TEXT NOT NULL | número + nombre normalizados (minúsculas, sin tildes) + dígitos del teléfono |
| entrega | TEXT CHECK IN ('retiro','envio') | |
| region, comuna | TEXT NULL | opcionales al registrar; comuna obligatoria para `despachado` (CHECK) |
| agencia | TEXT NULL | código de agencia; obligatoria para `despachado` (CHECK) |
| n_seguimiento | TEXT NULL | n.º de la agencia; obligatorio para `despachado` (CHECK) |
| forma_pago | TEXT CHECK IN ('transferencia','en_tienda') | lo acordado; por defecto regiones → transferencia, santiago → en_tienda; con envío, solo transferencia |
| subtotal | INTEGER | suma de las líneas con precio |
| descuento | INTEGER ≥ 0 | |
| costo_envio | INTEGER ≥ 0 | solo con envío; se suma al total |
| envio_por_pagar | INTEGER 0/1 | envío pagado por el cliente en la agencia (solo con envío y sin costo) |
| total | INTEGER CHECK ≥ 0 | calculado, o escrito a mano si hay líneas sin precio |
| total_manual | INTEGER 0/1 | 1 = total escrito por quien atiende |
| motivo_cancelacion | TEXT NOT NULL DEFAULT '' | sale en el mensaje d6 |
| nota_interna | TEXT NOT NULL DEFAULT '' | solo equipo; no sale en mensajes ni seguimiento |
| estado | TEXT CHECK IN (cotizado, esperando_pago, pagado, preparando, despachado, listo_retiro, entregado, cancelado) | |
| estado_desde | TEXT NOT NULL | cuándo entró al estado actual (pantalla Hoy: "esperando pago > 24 h") |
| stock_descontado | INTEGER 0/1 | 1 desde que pasó a `pagado` o `preparando`; 0 de nuevo si se cancela |
| pago_medio | TEXT NULL CHECK IN (transferencia, efectivo, debito, credito, otro) | |
| pago_monto | INTEGER NULL | |
| pago_fecha | TEXT NULL (`AAAA-MM-DD`) | |
| pago_verificado_por | INTEGER NULL → usuarios | quien registró el pago (al marcar `pagado`, o al entregar si pagó al retirar) |
| registrado_por | INTEGER NULL → usuarios | NULL en pedidos DEMO de la semilla |
| es_demo | INTEGER 0/1 | con `MODO_DEMO=true`, todo pedido nuevo nace demo |
| token_vence_en | TEXT NULL (UTC) | se fija al pasar a entregado o cancelado (+30 días); NULL = no vence mientras siga abierto |
| creado_en, actualizado_en | TEXT | UTC |

Índices: `ix_pedidos_estado_desde (estado, estado_desde)`, `ix_pedidos_creado (creado_en)`, `ix_pedidos_canal_creado (canal, creado_en)`, `ix_pedidos_telefono (cliente_telefono)`.

### `items_pedido`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| pedido_id | → pedidos ON DELETE CASCADE | |
| producto_id | INTEGER NULL → productos | NULL = línea libre |
| codigo | TEXT NULL | congelado (catálogo) |
| descripcion | TEXT NOT NULL | nombre congelado o texto libre |
| cantidad | INTEGER CHECK > 0 | 1–10.000 |
| precio_unitario | INTEGER NULL CHECK ≥ 0 | congelado; NULL = línea libre sin precio (las de catálogo siempre tienen) |
| tipo_precio | TEXT CHECK IN ('detalle','oferta','mayor','manual','libre') | `manual` = catálogo con precio editado |
| subtotal | INTEGER NULL | NULL si no hay precio |

Editar un pedido (`cotizado` o `esperando_pago`) reemplaza todas sus líneas.

### `historial_estados`
`pedido_id`, `estado_anterior` (NULL al crear), `estado_nuevo`, `usuario_id` (NULL = semilla), `nota`, `creado_en`. La línea de tiempo pública sale de aquí (solo estado y fecha).

## Máquina de estados
```
cotizado ──► esperando_pago ──► pagado ──► preparando ──► despachado (envío) ───► entregado
    │              │                 ▲          └──► listo_retiro (retiro) ─► entregado
    └──────────────┴── (en_tienda) ──┴──► preparando     (paga al retirar)
cualquiera no final ──► cancelado          finales: entregado, cancelado
```
- **A `pagado`:** exige `pago` (medio, monto, fecha opcional) y `confirmo_pago_recibido` (la plata está en la cuenta: en el barrio hay comprobantes falsos). Guarda quién lo verificó.
- **A `preparando` sin pagar:** solo si `forma_pago = en_tienda` (retiro que paga en el local).
- **Stock:** se descuenta al pasar a `pagado` o a `preparando`, lo que ocurra primero, con `UPDATE … WHERE stock >= ?`. Si alguna línea no alcanza, no cambia nada y responde `409 stock_insuficiente`. Una cotización no compromete mercadería: en Meiggs el cliente a veces no paga, y así se evitan "agotados" falsos.
- **A `despachado`:** exige entrega `envio`, agencia, `n_seguimiento` y comuna.
- **A `listo_retiro`:** exige entrega `retiro`.
- **A `entregado`:** si no hay pago registrado (pagó al retirar), exige `pago`.
- **A `cancelado`:** `motivo` opcional. Si `stock_descontado = 1`, devuelve el stock y lo pone en 0.
- Cada cambio queda en `historial_estados` con usuario, fecha y nota.

## Seguimiento público
`GET /api/seguimiento/{token}` lee el pedido por `token_seguimiento` y devuelve solo número, estado, línea de tiempo, tipo de entrega, agencia y n.º de seguimiento. Nunca datos personales ni montos.

## Semilla del demo
- Catálogo: 37 productos de `research/productos-semilla.json` (ver `api/README.md`). Se carga si no hay categorías.
- Pedidos DEMO (`api/internal/seed/pedidos_demo.json`): 20 pedidos **ficticios** (`es_demo = 1`). DC-0101 a DC-0109 son los del guion de `propuesta/pitch.md` §c y los ejemplos de `propuesta/mensajes-whatsapp.md`. DC-0090 a DC-0100 son pedidos anteriores ya cerrados, para que el resumen tenga historia sin cambiar la pantalla Hoy del guion. Fechas relativas al arranque. Nombres "… Demo" y teléfonos `+56 9 0000 0xxx`. Se cargan si la tabla `pedidos` está vacía (`SEED_PEDIDOS_DEMO=true`); el correlativo parte en 89 para que el siguiente pedido real del demo sea DC-0110, como dice el guion.
- Usuarios demo `regiones1` ("Regiones 1") y `santiago1` ("Santiago 1"), vendedores, **sin clave utilizable** (hash de una clave aleatoria descartada). Para entrar con ellos, un admin les genera una con `POST /api/admin/usuarios/{id}/clave`.
