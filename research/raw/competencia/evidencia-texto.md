# Evidencia en bruto — competencia y barrio (analista-meiggs)

Fecha de consulta: 2026-11-14. Método: WebSearch + WebFetch (sin navegador; el navegador lo usaba investigador-marca).
**Capturas de pantalla: PENDIENTES.** WebFetch no toma imágenes y TikTok no se deja leer sin navegador ("Please wait..."). Este archivo guarda los extractos de texto que respaldan cada ficha. Las capturas se pueden sacar en otra pasada con el navegador, con las URLs de abajo.

No se copian los teléfonos de la competencia: no hacen falta para el análisis.

---

## Importadora Sale Todo — https://importadorasaletodo.com (HTTP 200)
Descargado con curl y analizado localmente; no se guarda el HTML de ~9,5 MB.
- `<title>`: "Importadora Sale Todo SpA".
- Unas 1.900 fichas con el atributo `se_puede_pedir`. Cada producto tiene variantes de precio llamadas "Detalle", "Mayor" y, a veces, "Embalaje".
- Mínimos por variante: `"minimo":"3"` aparece unas 3.100 veces y `"minimo":"1"` unas 1.950; hay casos de 6, 10, 12 y 24.
- Categorías (cantidad de productos aproximada): NUEVOS PRODUCTOS 143, JUGUETES 126, CUIDADO PERSONAL 126, TECNOLOGÍA 116, AUTOS ACCESORIOS 92, ILUMINACIÓN 82, **NAVIDAD 81**, ESCOLAR 40, HOGAR 39, **INVIERNO 39**, **TEMPORADA 37**, LIMPIEZA 35, **FIESTAS PATRIAS 33**, SILLAS 27, **VERANO 22**, DECORACIÓN 13.
- Formulario de pedido: nombre y apellidos*, RUT*, teléfono*, correo, dirección (solo si es a domicilio), ciudad de envío*, aclaración del envío o del producto faltante, agencia de envío*.
- El pedido se manda armando un enlace `wa.me/<teléfono>?text=<pedido>` (`encodeURIComponent`).
- Envíos: "Todos los pedidos se envían por pagar". El cliente elige la agencia: Starken, Varmontt, Mena, Pullman, TVP, Chevalier, Cruz del Sur, Cacem, Dibama, JT, Transvalle, Transporte Barrios (Antofagasta/Calama), TransCargo (Punta Arenas), entre otras. También hay "Retiro en local (Abate Molina 686)".
- Avisa el horario: fuera de él muestra "HORARIO NOS ENCONTRAMOS CERRADOS" (lo vimos a la 01:44).
- Instagram @imp.sale.todo (https://www.instagram.com/imp.sale.todo/): 52,5 mil seguidores. Bio: Abate Molina 686, "Mayorista | Embalaje | Detalle". El link de la bio va a la web.
- También vende en Falabella como seller "IMPORTADORA Y COMERCIALIZADORA SALE TODO SPA": https://www.falabella.com/falabella-cl/seller/IMPORTADORA%20Y%20COMERCIALIZADORA%20SALE%20TODO%20SPA (no comprobé que sea la misma empresa, pero el nombre coincide).

## AliChile Importaciones
- Instagram https://www.instagram.com/alichileimportaciones/: 197 mil seguidores. Bio: "Importadora directa de juguetes, papelería y más", "Por mayor y al detalle", "Envíos a todo Chile". Destacados: Catálogo, Pago, Envíos, Ubicación.
- Linktree https://linktr.ee/alichileimportaciones: Catálogo (Whataform), WhatsApp, TikTok, Instagram, Facebook. **No tiene web propia.**
- Catálogo: https://s4f7cf41589304266.whataform.com/catalogoalichileimportaciones. El título dice "Catálogo ALICHILE"; el contenido se carga con JavaScript y no se pudo leer.
- Captions de TikTok (vistos en los resultados de búsqueda):
  - https://www.tiktok.com/@alichileimportaciones/video/7603750236724399380: chiste de que siempre responden "aunque el precio esté ahí mismo".
  - https://www.tiktok.com/@alichileimportaciones/video/7530672929575587079: "$6.500 x mayor $7.000 al detalle".
  - https://www.tiktok.com/@alichileimportaciones/video/7611589893550001429: les quisieron pagar con un comprobante de transferencia falso. El caption trae horario L–V 9:30–18:30 y sáb 9:30–16:30.
  - https://www.tiktok.com/@alichileimportaciones/video/7515544626904780038: menciona comentarios "en los videos y lives".
  - https://www.tiktok.com/@alichileimportaciones/video/7394461063053364485: productos recién llegados "se acaban rápido".
  - https://www.tiktok.com/@alichileimportaciones/video/7584428785374874901: "catálogo mayorista completo en el link de nuestra bio".
  - El caption se repite siempre igual: dirección (Bascuñán Guerrero 222, local 22) + metro Unión Latinoamericana + "x mayor y al detalle" + WhatsApp + "Envíos a todo Chile".
  - Según el resumen del buscador, el mínimo para envío es de $50.000 en los posts recientes. No lo verifiqué en un post puntual.

## DAKA Importaciones — https://www.instagram.com/daka_importaciones/
- 160 mil seguidores. La bio dice que es aliada en compras mayoristas y "la importadora preferida de los comerciantes en Meiggs"; también "Haz tu pedido por…" (texto cortado).
- Destacados: SAN ALFONSO 562, CONTÁCTANOS, EJECUTIVA/VENTA, ENVÍOS, HORARIO, PÁGINA WEB. No se ve la URL de la web.

## Emprende por mayor (Importadora AYJ) — https://www.instagram.com/importadora_ayj/
- 105 mil seguidores. Bio: WhatsApp, detalle y por mayor, Salvador Sanfuentes 3092 y Exposición 176, envíos.
- TikTok https://www.tiktok.com/@importadora_ayj/video/7589714221693209863: útiles escolares, "Envíos a todo Chile por compras desde $20.000", "Revisa nuestro catálogo y agenda tu pedido".

## Importadora Estrellita — https://www.instagram.com/estrellita_meiggs/
- WebFetch no devolvió el texto del perfil. Según el resultado de búsqueda: mayor y detalle, productos de temporada, Meiggs 75, lunes a sábado, envíos a todo Chile. **Por verificar.**

## Importadora Chinito Tao — https://www.tiktok.com/@chinitotao
- Solo tenemos el resumen del buscador: Alameda 2538, cerca del Metro ULA, compra por WhatsApp. **Por verificar.** Está a ~100 números de Alameda 2432.

## Tu Mercado Outlet — https://www.instagram.com/tumercadooutlet.cl/
- 1.189 seguidores. Tecnología, hogar y accesorios; Salvador Sanfuentes 2647, local 204; envíos a todo Chile.
- El link de la bio, www.tumercadooutlet.cl, **no resuelve DNS** (ENOTFOUND, 2026-11-14).

## Comercial Xing Xing II — https://www.instagram.com/comercialxingxingii/
- 319 seguidores. Venta mayorista y minorista de artículos para el hogar en Meiggs. El link de la bio, comercialxingxingii.com, **no resuelve DNS**.

## Cuentas que publican la dirección Alameda 2432 (no se afirma ninguna relación)
- Mayte emporio @mayte_emporioo (https://www.instagram.com/mayte_emporioo/): **cuenta verificada**, 659 mil seguidores. Bio: alameda 2432, número +56921145368, "pedidos solo a este número". El link de la bio va a Threads.
- ANDRES EMPORIO @andresemporio (https://www.instagram.com/andresemporio/): 60,4 mil seguidores, Alameda 2432, "Solo venta presencial". No figura como verificada.
- Emporio Premium @familia_emporio (https://www.instagram.com/familia_emporio/): 113 mil seguidores, "Productos para tu hogar y emprendimiento", Alameda 2432. El link de la bio, familiaemporio.com, **no resuelve DNS**.
- Facebook "Donde Carlitos Emporio" (https://www.facebook.com/donde.carlitos.33/) y Facebook "Mayte Emporio" (https://www.facebook.com/juan.jzjd/): aparecen en búsqueda. **Luckas dice que el negocio no tiene Facebook** (BITACORA, Fase 0). Hay que preguntarle al dueño si son suyas, antiguas o de terceros.
- TikTok @_dondecarlitosemporio7: aparece en búsqueda con videos que mencionan Alameda 2432. La cuenta oficial confirmada es @mayte_emporio. **Hay que confirmar con el dueño.**
