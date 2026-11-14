# Modelo de datos — Donde Carlitos Emporio

*ingeniero-go · 2026-11-14 · implementado en `api/migrations/0001_inicial.sql`, `0002_stock_al_confirmar.sql` y `0003_campos_qa.sql` (SQLite, `modernc.org/sqlite`). Contrato HTTP en `docs/openapi.yaml` (v1.1).*

**Cambios v1.1 (auditoría de seguridad y QA, ronda 1):**
- El stock **ya no se descuenta al crear el pedido**: se descuenta al confirmarlo (ver "Máquina de estados").
- Nuevos campos de `productos`: `mayor_minimo_confirmado`, `oferta_hasta`, `precio_actualizado_en`.
- Imágenes subidas desde el panel: archivos en un volumen Docker, servidos en `/api/imagenes/{archivo}`.
- Todas las fechas-hora se guardan y se entregan en UTC.

## Principios
- **Precios en pesos chilenos enteros** (`INTEGER`). Nunca decimales.
- **El precio del pedido se congela** en `items_pedido` al momento de pedir. Si después cambia el precio del producto, el pedido no cambia.
- **El servidor recalcula todo.** El total del navegador se ignora.
- **Stock interno, disponibilidad pública.** La API pública solo expone `disponible | quedan_pocos | agotado`.
- **Fechas-hora** como `TEXT` ISO-8601 UTC de ancho fijo (`2026-11-14T15:00:00.000Z`), así el orden de texto es el orden temporal. Fechas simples como `AAAA-MM-DD`.
- **Borrado lógico** de productos (`eliminado_en`), para no romper pedidos antiguos.
- **Portabilidad a PostgreSQL:** el SQL usa parámetros `?`, `RETURNING`, índices parciales y `ON CONFLICT`, todo soportado por PostgreSQL. Al migrar: `INTEGER PRIMARY KEY` → `BIGSERIAL`, fechas `TEXT` → `TIMESTAMPTZ`/`DATE`, `?` → `$n`. La aplicación solo conoce la interfaz `store.Store` (`api/internal/store/store.go`).

## Diagrama

```
categorias 1 ──< productos 1 ──< items_pedido >── 1 pedidos 1 ──< historial_estados >── 0..1 admins
                                                                                     admins 1 ──< sesiones
contadores (correlativo de pedidos)      schema_migraciones (versiones aplicadas)
```

## Tablas

### `categorias`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| slug | TEXT NOT NULL UNIQUE | `navidad`, `juguetes-y-virales` |
| nombre | TEXT NOT NULL | |
| orden | INTEGER NOT NULL DEFAULT 0 | orden en menús |
| creado_en | TEXT NOT NULL | |

### `productos`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| slug | TEXT NOT NULL | único entre no eliminados (índice parcial) |
| codigo | TEXT NOT NULL | código interno (`P001`); único entre no eliminados; va en el mensaje de WhatsApp |
| nombre | TEXT NOT NULL | 2–120 caracteres |
| descripcion | TEXT NOT NULL DEFAULT '' | lo que no tiene fuente va con `[EJEMPLO]` |
| busqueda | TEXT NOT NULL | nombre + código + descripción + categoría en minúsculas y sin tildes; lo calcula Go al guardar |
| categoria_id | INTEGER NOT NULL FK → categorias | |
| precio_detalle | INTEGER NOT NULL CHECK > 0 | precio normal al detalle |
| precio_oferta | INTEGER NULL CHECK 0 < x < precio_detalle | precio rebajado |
| etiqueta_oferta | TEXT NOT NULL DEFAULT 'oferta' CHECK IN ('oferta','remate') | badge que muestra la oferta |
| oferta_hasta | TEXT NULL (`AAAA-MM-DD`) | último día (hora de Chile) en que vale `precio_oferta`; NULL = sin vencimiento. Vencida: no cuenta en `precio_final`, badges, filtro `oferta`, orden por precio ni pedidos. Semilla: solo P012 (tren grande) hasta 2026-11-15, porque su video dice "hasta el domingo" |
| precio_actualizado_en | TEXT NULL (UTC) | lo fija el servidor al crear un producto en el panel y cuando cambia `precio_detalle`, `precio_oferta` o `precio_mayor`; NULL en la semilla (precio de la fuente) |
| precio_mayor | INTEGER NULL CHECK 0 < x < precio_detalle | **por unidad**; aplica solo desde `mayor_minimo` |
| mayor_minimo | INTEGER NULL CHECK ≥ 2 | unidades mínimas para el precio por mayor; va junto con `precio_mayor` (CHECK) |
| unidad_mayor | TEXT NULL | etiqueta visible: "caja de 12", "embalaje de 16", "desde 10 u." |
| mayor_minimo_confirmado | INTEGER 0/1 NOT NULL DEFAULT 0 | 1 = el mínimo salió del video; 0 = mínimo `[EJEMPLO]` (en la semilla: P009, P012, P016, P017, P020, P028, P037, P038) |
| stock | INTEGER NOT NULL CHECK ≥ 0 | **interno**, nunca sale en la API pública |
| umbral_pocos | INTEGER NOT NULL DEFAULT 5 | stock ≤ umbral → `quedan_pocos` |
| activo | INTEGER 0/1 NOT NULL DEFAULT 1 | inactivo = oculto del catálogo público |
| temporada | TEXT NULL CHECK IN (navidad, escolar, fiestas_patrias, halloween, invierno, verano) | badge `temporada` |
| fecha_llegada | TEXT NULL (`AAAA-MM-DD`) | badge `recien_llegado` si está dentro de los últimos `DIAS_RECIEN_LLEGADO` días (14) |
| imagen_url | TEXT NULL | ruta del front (`/productos/tt-<id>.webp`), imagen subida (`/api/imagenes/<32 hex>.jpg\|png`) o https |
| tiktok_url | TEXT NULL | video de TikTok donde aparece (solo https a tiktok.com) |
| fuente_url | TEXT NULL | de dónde sale el precio (regla 1 y 2 de CLAUDE.md) |
| fuente_fecha | TEXT NULL (`AAAA-MM-DD`) | fecha de ese precio; si es antigua el front avisa "reconfirmar" |
| es_ejemplo | INTEGER 0/1 | 1 = producto o precio sin fuente real (`[EJEMPLO]`) |
| creado_en, actualizado_en | TEXT NOT NULL | |
| eliminado_en | TEXT NULL | borrado lógico |

Índices: `ux_productos_slug`, `ux_productos_codigo` (únicos parciales `WHERE eliminado_en IS NULL`), `ix_productos_categoria`, `ix_productos_visibles (activo, eliminado_en)`, `ix_productos_temporada`, `ix_productos_llegada`.

**Regla de precio por línea** (`productos.PrecioUnitario`): el menor entre `precio_detalle`, `precio_oferta` **si está vigente** (`oferta_hasta` NULL o ≥ hoy en Chile) y, **solo si `cantidad ≥ mayor_minimo`**, `precio_mayor`. Se guarda qué tipo aplicó (`detalle | oferta | mayor`). Máximo 100 unidades por línea en el pedido web.

**Imágenes subidas** (`POST /api/admin/productos/{id}/imagen`): no hay tabla. El archivo se re-codifica (sin EXIF), se guarda con nombre aleatorio en el volumen `imagenes-api` (`/imagenes`) y su ruta queda en `imagen_url`. Al reemplazarla se borra el archivo anterior.

**Por qué un solo escalón por mayor.** En los videos hay productos con dos escalones ("$9.000 desde 10 u., $8.500 por embalaje"). Para el demo se modela un escalón con mínimo conocido (`precio_mayor` + `mayor_minimo`) y el resto se cuenta en `unidad_mayor`/`descripcion`. Siguiente paso, si el dueño lo confirma: tabla `precios_mayor(producto_id, minimo, precio)`.

### `contadores`
| Campo | Tipo | Notas |
|---|---|---|
| nombre | TEXT PK | `'pedidos'` |
| valor | INTEGER NOT NULL | último correlativo usado |

El número se toma con `UPDATE contadores SET valor = valor + 1 WHERE nombre = 'pedidos' RETURNING valor` **dentro de la misma transacción** que inserta el pedido (SQLite en `BEGIN IMMEDIATE`). Si el pedido falla, se revierte y el número no se pierde: la serie no tiene huecos ni duplicados (probado con 60 pedidos concurrentes).

### `pedidos`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| secuencia | INTEGER NOT NULL UNIQUE | correlativo |
| numero | TEXT NOT NULL UNIQUE | `DC-0001` (`DC-%04d`; desde 10.000 crece a 5 dígitos) |
| cliente_nombre | TEXT NOT NULL | |
| cliente_telefono | TEXT NOT NULL | E.164 `+56912345678`; en logs va enmascarado `+569****5678` |
| entrega | TEXT NOT NULL CHECK IN ('retiro','envio') | |
| region | TEXT NULL | código (`nuble`); obligatorio si envío |
| comuna | TEXT NULL | obligatoria si envío; opcional si retiro |
| agencia | TEXT NULL | código (`starken`); obligatorio si envío. Lista `[EJEMPLO]` |
| comentario | TEXT NOT NULL DEFAULT '' | máx. 500 |
| total | INTEGER NOT NULL | suma de subtotales, calculada por el servidor |
| estado | TEXT NOT NULL CHECK IN (nuevo, confirmado, listo, despachado, entregado, cancelado) | |
| es_demo | INTEGER 0/1 | pedido ficticio (semilla o modo demo) |
| creado_en, actualizado_en | TEXT NOT NULL | |

CHECK: si `entrega = 'envio'`, `region`, `comuna` y `agencia` no pueden ser NULL. Índices: `ix_pedidos_estado_creado`, `ix_pedidos_creado`, `ix_pedidos_telefono_estado` (para el tope de 3 pedidos `nuevo` abiertos por teléfono, que se cuenta dentro de la misma transacción que crea el pedido).

No se guarda la IP ni ningún dato del cliente fuera de nombre, teléfono, comuna/región y comentario.

### `items_pedido`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| pedido_id | INTEGER FK → pedidos ON DELETE CASCADE | |
| producto_id | INTEGER FK → productos | |
| codigo, nombre | TEXT NOT NULL | congelados |
| cantidad | INTEGER CHECK > 0 | 1–500 (validado en Go) |
| precio_unitario | INTEGER CHECK > 0 | **congelado** |
| tipo_precio | TEXT CHECK IN ('detalle','oferta','mayor') | |
| subtotal | INTEGER | cantidad × precio_unitario |

### `historial_estados`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| pedido_id | INTEGER FK → pedidos | |
| estado_anterior | TEXT NULL | NULL en la creación |
| estado_nuevo | TEXT NOT NULL | |
| admin_id | INTEGER NULL FK → admins ON DELETE SET NULL | NULL = creación o semilla |
| nota | TEXT NOT NULL DEFAULT '' | máx. 200 |
| creado_en | TEXT NOT NULL | |

### `admins`
| Campo | Tipo | Notas |
|---|---|---|
| id | INTEGER PK | |
| usuario | TEXT NOT NULL UNIQUE | `ADMIN_USUARIO` (por defecto `admin`) |
| password_hash | TEXT NOT NULL | bcrypt costo 12. Llega desde `ADMIN_PASSWORD_HASH_B64` al arrancar; si cambió, se cierran las sesiones de ese admin |
| creado_en, actualizado_en | TEXT NOT NULL | |

### `sesiones`
| Campo | Tipo | Notas |
|---|---|---|
| token_hash | TEXT PK | **HMAC-SHA256(SESSION_SECRET, token)**; el token en claro solo vive en la cookie |
| admin_id | INTEGER FK → admins ON DELETE CASCADE | |
| creado_en | TEXT | |
| expira_en | TEXT | 12 h desde el login (absoluto) |
| ultimo_uso | TEXT | vence con 2 h de inactividad |

Las sesiones vencidas se purgan cada 30 minutos.

### `schema_migraciones`
`version TEXT PK, aplicada_en TEXT`. Las migraciones `api/migrations/*.sql` van embebidas en el binario y se aplican en orden al arrancar.

## Máquina de estados del pedido

```
nuevo ──► confirmado ──► listo (solo retiro) ──► entregado
  │            │    └──► despachado (solo envío) ──► entregado
  │            │                 
  └──► cancelado ◄──┴──────── listo
```
- `entregado` y `cancelado` son finales. `despachado` no se cancela (la mercadería ya salió).
- **Stock (v1.1):**
  - Crear un pedido (`nuevo`) **no descuenta**: solo verifica que el stock alcance. Así un pedido anónimo sin confirmar no puede dejar la tienda "agotada".
  - Pasar a `confirmado` **descuenta**, con un `UPDATE … WHERE stock >= ?` por línea. Si alguna no alcanza, se revierte todo y no se confirma (`409 stock_insuficiente`). Probado con 20 confirmaciones simultáneas para 5 unidades: 5 pasan, 15 no.
  - Cancelar un `confirmado` o `listo` **devuelve** el stock. Cancelar un `nuevo` no lo toca.
  - La migración 0002 devolvió el stock que la regla anterior había descontado a los pedidos que estaban en `nuevo`.

## Semilla del demo
- `api/internal/seed/semilla.json`, generada desde `research/productos-semilla.json` (investigador-marca, TikTok @mayte_emporio):
  - 37 productos con **nombre, precio, video (`tiktok_url` y `fuente_url`) y fecha reales**. 3 se omitieron por no tener precio dicho.
  - Ofertas reales del video ("antes $20.000, ahora $17.000") como `precio_oferta`.
  - `[EJEMPLO]` en la descripción para lo que completó ingeniero-go: stock/disponibilidad y los mínimos por mayor que el video no dijo (`mayor_minimo_confirmado = 0`).
  - Stock siempre sobre el umbral de "quedan pocos" (sin urgencia falsa), salvo P036 (repisa), agotado a propósito para mostrar ese estado.
  - Códigos `P0xx`: internos del demo.
- 6 pedidos `es_demo = 1` con nombres "Cliente Demo Uno…Seis", teléfonos `+56 9 0000 000N` y comentario `[DEMO]`, en distintos estados.
- Se carga solo si la base está vacía (`SEED_DEMO=true`). Para reiniciar el demo: `docker compose down -v`.
