-- 0005_ronda3.sql — QA ronda 2 y seguridad ronda 2 (contrato 2.1.0).

-- B-2: precio de promo por comentario (no es precio de detalle). I-2: vitrina elegida por el dueño.
ALTER TABLE productos ADD COLUMN precio_es_promo INTEGER NOT NULL DEFAULT 0 CHECK (precio_es_promo IN (0, 1));
ALTER TABLE productos ADD COLUMN nota_promo TEXT;
ALTER TABLE productos ADD COLUMN en_vitrina INTEGER NOT NULL DEFAULT 0 CHECK (en_vitrina IN (0, 1));
ALTER TABLE productos ADD COLUMN orden_vitrina INTEGER CHECK (orden_vitrina IS NULL OR orden_vitrina BETWEEN 1 AND 99);
CREATE INDEX ix_productos_vitrina ON productos(en_vitrina, orden_vitrina);

-- M-15: canal por defecto del usuario. B9: clave temporal que hay que cambiar.
ALTER TABLE usuarios ADD COLUMN canal_por_defecto TEXT CHECK (canal_por_defecto IS NULL OR canal_por_defecto IN ('santiago', 'regiones'));
ALTER TABLE usuarios ADD COLUMN debe_cambiar_clave INTEGER NOT NULL DEFAULT 0 CHECK (debe_cambiar_clave IN (0, 1));

-- B10: el link de seguimiento vence 30 días después de entregado o cancelado.
ALTER TABLE pedidos ADD COLUMN token_vence_en TEXT;
UPDATE pedidos SET token_vence_en = strftime('%Y-%m-%dT%H:%M:%fZ', estado_desde, '+30 days')
 WHERE estado IN ('entregado', 'cancelado');

-- M-5: línea de catálogo con precio editado por quien atiende (tipo 'manual').
-- SQLite no permite cambiar un CHECK: se reconstruye la tabla.
CREATE TABLE items_pedido_nuevo (
    id              INTEGER PRIMARY KEY,
    pedido_id       INTEGER NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    producto_id     INTEGER REFERENCES productos(id),
    codigo          TEXT,
    descripcion     TEXT    NOT NULL,
    cantidad        INTEGER NOT NULL CHECK (cantidad > 0),
    precio_unitario INTEGER CHECK (precio_unitario IS NULL OR precio_unitario >= 0),
    tipo_precio     TEXT    NOT NULL CHECK (tipo_precio IN ('detalle', 'oferta', 'mayor', 'manual', 'libre')),
    subtotal        INTEGER CHECK (subtotal IS NULL OR subtotal >= 0),
    CHECK ((producto_id IS NULL) = (tipo_precio = 'libre')),
    CHECK (producto_id IS NULL OR precio_unitario IS NOT NULL),
    CHECK ((precio_unitario IS NULL) = (subtotal IS NULL))
);
INSERT INTO items_pedido_nuevo (id, pedido_id, producto_id, codigo, descripcion, cantidad, precio_unitario, tipo_precio, subtotal)
    SELECT id, pedido_id, producto_id, codigo, descripcion, cantidad, precio_unitario, tipo_precio, subtotal FROM items_pedido;
DROP TABLE items_pedido;
ALTER TABLE items_pedido_nuevo RENAME TO items_pedido;
CREATE INDEX ix_items_pedido_pedido   ON items_pedido(pedido_id);
CREATE INDEX ix_items_pedido_producto ON items_pedido(producto_id);

-- Ajuste de la semilla ya cargada (en una base nueva no hay filas; semilla.json ya trae estos valores).
-- B-2: precios de promo por comentario según research/productos-semilla.json (campo nota).
UPDATE productos SET precio_es_promo = 1, nota_promo = 'Comenta ''yo quiero mi promoción'' en el video' WHERE codigo = 'P002';
UPDATE productos SET precio_es_promo = 1, nota_promo = 'Precio de promo por comentario en el video' WHERE codigo IN ('P007', 'P011', 'P013');
UPDATE productos SET precio_es_promo = 1, nota_promo = 'Precio de promo del video (30 vasos)' WHERE codigo = 'P008';
-- I-2: vitrina con los productos del guion.
UPDATE productos SET en_vitrina = 1, orden_vitrina = CASE codigo
        WHEN 'P012' THEN 1 WHEN 'P010' THEN 2 WHEN 'P016' THEN 3 WHEN 'P001' THEN 4
        WHEN 'P005' THEN 5 WHEN 'P020' THEN 6 WHEN 'P017' THEN 7 WHEN 'P027' THEN 8 END
 WHERE codigo IN ('P012', 'P010', 'P016', 'P001', 'P005', 'P020', 'P017', 'P027');
-- M-11: sin paréntesis (la web pone la unidad entre paréntesis).
UPDATE productos SET unidad_mayor = 'embalaje surtido de 16' WHERE unidad_mayor = 'embalaje de 16 (surtido)';
UPDATE productos SET unidad_mayor = 'docena surtida a $6.000' WHERE unidad_mayor = 'docena surtida ($6.000)';
-- I-1: stock alto en los productos del guion (ensayar dos veces no lo agota).
UPDATE productos SET stock = 300 WHERE stock < 300
   AND codigo IN ('P001', 'P005', 'P007', 'P010', 'P012', 'P016', 'P017', 'P019', 'P020');
-- Usuarios demo: canal por defecto; si un admin les genera clave, deberán cambiarla.
UPDATE usuarios SET canal_por_defecto = 'regiones', debe_cambiar_clave = 1 WHERE usuario = 'regiones1';
UPDATE usuarios SET canal_por_defecto = 'santiago', debe_cambiar_clave = 1 WHERE usuario = 'santiago1';
