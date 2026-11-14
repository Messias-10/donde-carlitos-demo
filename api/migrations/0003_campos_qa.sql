-- 0003_campos_qa.sql — campos pedidos por QA ronda 1 (I-2, I-3, M-15).

-- I-2: el mínimo por mayor salió del video (1) o es [EJEMPLO] (0).
ALTER TABLE productos ADD COLUMN mayor_minimo_confirmado INTEGER NOT NULL DEFAULT 0
    CHECK (mayor_minimo_confirmado IN (0, 1));
-- I-3: cuándo el admin cambió por última vez un precio (NULL = precio de la fuente).
ALTER TABLE productos ADD COLUMN precio_actualizado_en TEXT;
-- M-15: último día (AAAA-MM-DD, hora de Chile) en que vale precio_oferta. NULL = sin vencimiento.
ALTER TABLE productos ADD COLUMN oferta_hasta TEXT;

-- Ajuste de la semilla ya cargada en bases existentes (en una base nueva no
-- hay filas y esto no hace nada; la semilla nueva ya trae estos valores).
-- Mínimos por mayor que sí salen del video: todos menos los 8 [EJEMPLO].
UPDATE productos SET mayor_minimo_confirmado = 1
 WHERE precio_mayor IS NOT NULL
   AND codigo IN ('P001','P003','P004','P005','P010','P021','P022','P025','P026','P039');
-- Tren navideño grande: el video (2026-11-10, martes) dice "válido hasta el domingo".
UPDATE productos SET oferta_hasta = '2026-11-15' WHERE codigo = 'P012' AND precio_oferta IS NOT NULL;
-- I-1: sin urgencia falsa; solo P036 queda agotado a propósito.
UPDATE productos SET stock = 30 WHERE codigo IN ('P005','P011','P027') AND stock <= umbral_pocos;
