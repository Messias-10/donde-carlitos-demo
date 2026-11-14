-- 0002_stock_al_confirmar.sql
-- Cambio de regla (auditoría de seguridad, ronda 1, A2): el stock ya no se
-- reserva al crear el pedido sino al confirmarlo. Los pedidos que hoy están
-- en "nuevo" se crearon con la regla anterior y tienen stock descontado:
-- se devuelve para que confirmar no lo descuente dos veces.

UPDATE productos
   SET stock = stock + (
         SELECT COALESCE(SUM(i.cantidad), 0)
           FROM items_pedido i
           JOIN pedidos p ON p.id = i.pedido_id
          WHERE p.estado = 'nuevo' AND i.producto_id = productos.id)
 WHERE id IN (
         SELECT i.producto_id
           FROM items_pedido i
           JOIN pedidos p ON p.id = i.pedido_id
          WHERE p.estado = 'nuevo');

-- Para contar pedidos "nuevo" abiertos por teléfono al crear uno.
CREATE INDEX ix_pedidos_telefono_estado ON pedidos(cliente_telefono, estado);
