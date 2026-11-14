-- 0001_inicial.sql — esquema inicial de Donde Carlitos.
-- SQL portable a PostgreSQL salvo los tipos (INTEGER PRIMARY KEY -> BIGSERIAL,
-- TEXT de fechas -> TIMESTAMPTZ). Ver docs/modelo-datos.md.

CREATE TABLE categorias (
    id         INTEGER PRIMARY KEY,
    slug       TEXT    NOT NULL UNIQUE,
    nombre     TEXT    NOT NULL,
    orden      INTEGER NOT NULL DEFAULT 0,
    creado_en  TEXT    NOT NULL
);

CREATE TABLE productos (
    id              INTEGER PRIMARY KEY,
    slug            TEXT    NOT NULL,
    codigo          TEXT    NOT NULL,
    nombre          TEXT    NOT NULL,
    descripcion     TEXT    NOT NULL DEFAULT '',
    busqueda        TEXT    NOT NULL DEFAULT '',
    categoria_id    INTEGER NOT NULL REFERENCES categorias(id),
    precio_detalle  INTEGER NOT NULL CHECK (precio_detalle > 0),
    precio_oferta   INTEGER CHECK (precio_oferta IS NULL OR (precio_oferta > 0 AND precio_oferta < precio_detalle)),
    etiqueta_oferta TEXT    NOT NULL DEFAULT 'oferta' CHECK (etiqueta_oferta IN ('oferta', 'remate')),
    precio_mayor    INTEGER CHECK (precio_mayor IS NULL OR (precio_mayor > 0 AND precio_mayor < precio_detalle)),
    mayor_minimo    INTEGER CHECK (mayor_minimo IS NULL OR mayor_minimo >= 2),
    unidad_mayor    TEXT,   -- etiqueta visible del precio por mayor: "caja de 12", "embalaje de 16"
    stock           INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    umbral_pocos    INTEGER NOT NULL DEFAULT 5 CHECK (umbral_pocos >= 0),
    activo          INTEGER NOT NULL DEFAULT 1 CHECK (activo IN (0, 1)),
    temporada       TEXT    CHECK (temporada IS NULL OR temporada IN ('navidad', 'escolar', 'fiestas_patrias', 'halloween', 'invierno', 'verano')),
    fecha_llegada   TEXT,
    imagen_url      TEXT,
    tiktok_url      TEXT,
    fuente_url      TEXT,   -- de dónde sale el precio (video de TikTok, post)
    fuente_fecha    TEXT,   -- YYYY-MM-DD en que se publicó ese precio
    es_ejemplo      INTEGER NOT NULL DEFAULT 0 CHECK (es_ejemplo IN (0, 1)),
    creado_en       TEXT    NOT NULL,
    actualizado_en  TEXT    NOT NULL,
    eliminado_en    TEXT,
    CHECK ((precio_mayor IS NULL) = (mayor_minimo IS NULL))
);

-- Slug y código únicos solo entre productos no eliminados (borrado lógico).
CREATE UNIQUE INDEX ux_productos_slug   ON productos(slug)   WHERE eliminado_en IS NULL;
CREATE UNIQUE INDEX ux_productos_codigo ON productos(codigo) WHERE eliminado_en IS NULL;
CREATE INDEX ix_productos_categoria     ON productos(categoria_id);
CREATE INDEX ix_productos_visibles      ON productos(activo, eliminado_en);
CREATE INDEX ix_productos_temporada     ON productos(temporada);
CREATE INDEX ix_productos_llegada       ON productos(fecha_llegada);

-- Contador atómico para el número correlativo de pedido (DC-0001, ...).
CREATE TABLE contadores (
    nombre TEXT    PRIMARY KEY,
    valor  INTEGER NOT NULL
);
INSERT INTO contadores (nombre, valor) VALUES ('pedidos', 0);

CREATE TABLE admins (
    id             INTEGER PRIMARY KEY,
    usuario        TEXT NOT NULL UNIQUE,
    password_hash  TEXT NOT NULL,
    creado_en      TEXT NOT NULL,
    actualizado_en TEXT NOT NULL
);

CREATE TABLE pedidos (
    id               INTEGER PRIMARY KEY,
    secuencia        INTEGER NOT NULL UNIQUE,
    numero           TEXT    NOT NULL UNIQUE,
    cliente_nombre   TEXT    NOT NULL,
    cliente_telefono TEXT    NOT NULL,
    entrega          TEXT    NOT NULL CHECK (entrega IN ('retiro', 'envio')),
    region           TEXT,
    comuna           TEXT,
    agencia          TEXT,
    comentario       TEXT    NOT NULL DEFAULT '',
    total            INTEGER NOT NULL CHECK (total >= 0),
    estado           TEXT    NOT NULL DEFAULT 'nuevo'
                     CHECK (estado IN ('nuevo', 'confirmado', 'listo', 'despachado', 'entregado', 'cancelado')),
    es_demo          INTEGER NOT NULL DEFAULT 0 CHECK (es_demo IN (0, 1)),
    creado_en        TEXT    NOT NULL,
    actualizado_en   TEXT    NOT NULL,
    CHECK (entrega = 'retiro' OR (region IS NOT NULL AND comuna IS NOT NULL AND agencia IS NOT NULL))
);
CREATE INDEX ix_pedidos_estado_creado ON pedidos(estado, creado_en);
CREATE INDEX ix_pedidos_creado        ON pedidos(creado_en);

CREATE TABLE items_pedido (
    id              INTEGER PRIMARY KEY,
    pedido_id       INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    producto_id     INTEGER NOT NULL REFERENCES productos(id),
    codigo          TEXT    NOT NULL,
    nombre          TEXT    NOT NULL,
    cantidad        INTEGER NOT NULL CHECK (cantidad > 0),
    precio_unitario INTEGER NOT NULL CHECK (precio_unitario > 0),
    tipo_precio     TEXT    NOT NULL CHECK (tipo_precio IN ('detalle', 'oferta', 'mayor')),
    subtotal        INTEGER NOT NULL CHECK (subtotal > 0)
);
CREATE INDEX ix_items_pedido_pedido   ON items_pedido(pedido_id);
CREATE INDEX ix_items_pedido_producto ON items_pedido(producto_id);

CREATE TABLE historial_estados (
    id              INTEGER PRIMARY KEY,
    pedido_id       INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    estado_anterior TEXT,
    estado_nuevo    TEXT    NOT NULL,
    admin_id        INTEGER REFERENCES admins(id) ON DELETE SET NULL,
    nota            TEXT    NOT NULL DEFAULT '',
    creado_en       TEXT    NOT NULL
);
CREATE INDEX ix_historial_pedido ON historial_estados(pedido_id, creado_en);

CREATE TABLE sesiones (
    token_hash TEXT    PRIMARY KEY,
    admin_id   INTEGER NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    creado_en  TEXT    NOT NULL,
    expira_en  TEXT    NOT NULL,
    ultimo_uso TEXT    NOT NULL
);
CREATE INDEX ix_sesiones_admin  ON sesiones(admin_id);
CREATE INDEX ix_sesiones_expira ON sesiones(expira_en);
