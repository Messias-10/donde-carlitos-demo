-- 0004_panel_pedidos.sql — v2: panel de gestión de pedidos (sin pedido público).
--
-- Los pedidos de la v1 eran todos ficticios (semilla DEMO y pruebas de QA). Se
-- borran y las tablas se recrean con el esquema nuevo. Antes, se devuelve al
-- catálogo el stock que esos pedidos tenían descontado (con la regla 0002, desde
-- "confirmado"). La semilla vuelve a cargar pedidos DEMO si la tabla queda vacía.

UPDATE productos
   SET stock = stock + (
         SELECT COALESCE(SUM(i.cantidad), 0)
           FROM items_pedido i
           JOIN pedidos p ON p.id = i.pedido_id
          WHERE p.estado IN ('confirmado', 'listo', 'despachado', 'entregado')
            AND i.producto_id = productos.id)
 WHERE id IN (
         SELECT i.producto_id
           FROM items_pedido i
           JOIN pedidos p ON p.id = i.pedido_id
          WHERE p.estado IN ('confirmado', 'listo', 'despachado', 'entregado'));

DROP TABLE historial_estados;
DROP TABLE items_pedido;
DROP TABLE pedidos;
UPDATE contadores SET valor = 0 WHERE nombre = 'pedidos';

-- admins -> usuarios, con nombre visible, rol y activo.
ALTER TABLE admins RENAME TO usuarios;
ALTER TABLE usuarios ADD COLUMN nombre TEXT NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN rol TEXT NOT NULL DEFAULT 'admin' CHECK (rol IN ('admin', 'vendedor'));
ALTER TABLE usuarios ADD COLUMN activo INTEGER NOT NULL DEFAULT 1 CHECK (activo IN (0, 1));
UPDATE usuarios SET nombre = 'Administrador' WHERE nombre = '';

ALTER TABLE sesiones RENAME COLUMN admin_id TO usuario_id;
DROP INDEX ix_sesiones_admin;
CREATE INDEX ix_sesiones_usuario ON sesiones(usuario_id);

CREATE TABLE pedidos (
    id                  INTEGER PRIMARY KEY,
    secuencia           INTEGER NOT NULL UNIQUE,
    numero              TEXT    NOT NULL UNIQUE,
    token_seguimiento   TEXT    NOT NULL UNIQUE,  -- 16 bytes aleatorios en base64url (22 caracteres)
    canal               TEXT    NOT NULL CHECK (canal IN ('santiago', 'regiones')),
    cliente_nombre      TEXT    NOT NULL,
    cliente_telefono    TEXT    NOT NULL,
    busqueda            TEXT    NOT NULL DEFAULT '',
    entrega             TEXT    NOT NULL CHECK (entrega IN ('retiro', 'envio')),
    region              TEXT,
    comuna              TEXT,
    agencia             TEXT,
    n_seguimiento       TEXT,
    forma_pago          TEXT    NOT NULL CHECK (forma_pago IN ('transferencia', 'en_tienda')),
    subtotal            INTEGER NOT NULL CHECK (subtotal >= 0),  -- suma de las líneas con precio
    descuento           INTEGER NOT NULL DEFAULT 0 CHECK (descuento >= 0),
    costo_envio         INTEGER NOT NULL DEFAULT 0 CHECK (costo_envio >= 0),
    envio_por_pagar     INTEGER NOT NULL DEFAULT 0 CHECK (envio_por_pagar IN (0, 1)),
    total               INTEGER NOT NULL CHECK (total >= 0),
    total_manual        INTEGER NOT NULL DEFAULT 0 CHECK (total_manual IN (0, 1)),  -- 1 = hay líneas sin precio
    motivo_cancelacion  TEXT    NOT NULL DEFAULT '',
    nota_interna        TEXT    NOT NULL DEFAULT '',
    estado              TEXT    NOT NULL CHECK (estado IN ('cotizado', 'esperando_pago', 'pagado', 'preparando',
                                                           'despachado', 'listo_retiro', 'entregado', 'cancelado')),
    estado_desde        TEXT    NOT NULL,
    stock_descontado    INTEGER NOT NULL DEFAULT 0 CHECK (stock_descontado IN (0, 1)),
    pago_medio          TEXT    CHECK (pago_medio IS NULL OR pago_medio IN ('transferencia', 'efectivo', 'debito', 'credito', 'otro')),
    pago_monto          INTEGER CHECK (pago_monto IS NULL OR pago_monto > 0),
    pago_fecha          TEXT,
    pago_verificado_por INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    registrado_por      INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    es_demo             INTEGER NOT NULL DEFAULT 0 CHECK (es_demo IN (0, 1)),
    creado_en           TEXT    NOT NULL,
    actualizado_en      TEXT    NOT NULL,
    CHECK (estado <> 'despachado' OR (agencia IS NOT NULL AND n_seguimiento IS NOT NULL AND comuna IS NOT NULL))
);
CREATE INDEX ix_pedidos_estado_desde ON pedidos(estado, estado_desde);
CREATE INDEX ix_pedidos_creado       ON pedidos(creado_en);
CREATE INDEX ix_pedidos_canal_creado ON pedidos(canal, creado_en);
CREATE INDEX ix_pedidos_telefono     ON pedidos(cliente_telefono);

CREATE TABLE items_pedido (
    id              INTEGER PRIMARY KEY,
    pedido_id       INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    producto_id     INTEGER REFERENCES productos(id),
    codigo          TEXT,
    descripcion     TEXT    NOT NULL,
    cantidad        INTEGER NOT NULL CHECK (cantidad > 0),
    precio_unitario INTEGER CHECK (precio_unitario IS NULL OR precio_unitario >= 0),  -- NULL = línea libre sin precio
    tipo_precio     TEXT    NOT NULL CHECK (tipo_precio IN ('detalle', 'oferta', 'mayor', 'libre')),
    subtotal        INTEGER CHECK (subtotal IS NULL OR subtotal >= 0),
    CHECK ((producto_id IS NULL) = (tipo_precio = 'libre')),
    CHECK (producto_id IS NULL OR precio_unitario IS NOT NULL),
    CHECK ((precio_unitario IS NULL) = (subtotal IS NULL))
);
CREATE INDEX ix_items_pedido_pedido   ON items_pedido(pedido_id);
CREATE INDEX ix_items_pedido_producto ON items_pedido(producto_id);

CREATE TABLE historial_estados (
    id              INTEGER PRIMARY KEY,
    pedido_id       INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    estado_anterior TEXT,
    estado_nuevo    TEXT    NOT NULL,
    usuario_id      INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    nota            TEXT    NOT NULL DEFAULT '',
    creado_en       TEXT    NOT NULL
);
CREATE INDEX ix_historial_pedido ON historial_estados(pedido_id, creado_en);
