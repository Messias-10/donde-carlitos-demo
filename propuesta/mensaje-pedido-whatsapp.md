> **OBSOLETO (2026-11-14).** Este archivo describe el mensaje del pedido armado desde el carrito web. Desde el cambio de enfoque aprobado por Luckas, **ya no hay carrito, checkout ni pedido desde la web**.
> **No implementar.** La plantilla vigente está en [`mensajes-whatsapp.md`](mensajes-whatsapp.md).
> Se conserva solo como historial.

# Mensaje de pedido por WhatsApp — plantilla exacta (OBSOLETO)

*estratega-propuesta · Fase 2 · 2026-11-14 · lo implementa constructor-web y lo verifica qa-revisor*

## 1. Cómo se usa

1. **Primero se guarda el pedido.** El cliente aprieta "Hacer pedido", el front hace el POST a la API y la API responde con el pedido creado: número `DC-xxxx`, líneas con precio congelado, subtotales y total.
2. **El mensaje se arma con la respuesta de la API, no con el carrito.** Así el total del WhatsApp es el mismo que queda en el panel.
3. En la pantalla de confirmación, el botón "Enviar pedido por WhatsApp" abre:
   ```
   https://wa.me/{NUMERO}?text={encodeURIComponent(mensaje)}
   ```
4. El botón "Copiar mensaje" copia **el mismo texto** sin codificar.

### Número de destino
- **Producción:** `56921145368` (+56 9 2114 5368, el único número de pedidos según las bios oficiales).
- **Va por configuración**, no fijo en el código. Ejemplo: `NEXT_PUBLIC_WHATSAPP_PEDIDOS`. Es un dato público, no un secreto.
- **Regla para el demo y QA:** mientras se prueba, el número de destino tiene que ser **un número de prueba de Luckas**, no el del negocio. Así no le llegan pedidos ficticios al WhatsApp real del dueño (CLAUDE.md: "los pedidos del demo son ficticios"). Solo se apunta al número real para mostrárselo al dueño en persona, si él lo pide.

### Por qué no va el teléfono del cliente en el mensaje
El mensaje sale desde el WhatsApp del propio cliente, así que el negocio ya ve su número. Además, el teléfono queda guardado en el panel. Si se repitiera en el texto, solo se expondría más.

## 2. Plantilla

**Reglas de formato:**
- Saltos de línea con `\n` (LF).
- Exactamente **una línea en blanco** entre bloques.
- Sin espacios al final de las líneas ni salto de línea al final del mensaje.
- `*texto*` es negrita en WhatsApp.
- Lo que va entre `[[ ]]` es condicional y **no se escribe**.

```
[[si es_demo]]*[PEDIDO DE PRUEBA · DEMO]*
Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️

*Pedido {numero}*

*Productos*
[[por cada línea, en el orden de la API]]
• {cantidad} x {nombre_producto}
  Cód. {codigo} · {precio_unitario} c/u ({tipo}) = {subtotal}
[[fin]]
[[si hay más de 25 líneas]]• … y {n} productos más (detalle completo en el pedido {numero})

*Total: {total}*[[si envío]] (no incluye envío)

[[si retiro]]
📍 *Entrega:* Retiro en tienda, Alameda 2432
[[si retiro y hay comuna]]*Comuna:* {comuna}
[[si envío]]
📦 *Entrega:* Envío por agencia {agencia}
*Región:* {region}
*Comuna:* {comuna}
[[siempre]]
*Nombre:* {nombre_cliente}
[[si hay comentario]]💬 *Comentario:* {comentario}

Quedo esperando su confirmación por este chat.
```

### Campos

| Campo | Origen (respuesta de la API) | Formato |
|---|---|---|
| `numero` | `numero` | `DC-0007` tal cual |
| `cantidad` | `items[].cantidad` | entero, sin separador |
| `nombre_producto` | `items[].nombre` | limpiado (ver §3) |
| `codigo` | `items[].codigo` | limpiado |
| `precio_unitario`, `subtotal`, `total` | `items[].precio_unitario`, `items[].subtotal`, `total` | CLP: `$` + punto de miles, sin decimales (ver §3) |
| `tipo` | `items[].tipo_precio` | `detalle` → `detalle` · `mayor` → `x mayor` · `oferta` → `oferta` |
| `agencia` | `agencia` (código) | nombre visible de la agencia, el mismo del selector (ej. `Starken` `[EJEMPLO]`) |
| `region` | `region` (código) | nombre visible de la región, el mismo del selector (ej. `Ñuble`) |
| `comuna` | `comuna` | limpiado |
| `nombre_cliente` | `nombre` | limpiado |
| `comentario` | `comentario` | limpiado; máximo 300 caracteres en el mensaje (si es más largo: 299 + `…`). El panel guarda el completo |
| `es_demo` | `es_demo` | si es `true`, va la primera línea de prueba |

### Emojis (solo 4, todos de la lista de estilo-carlitos §5)
| Emoji | Código | Dónde |
|---|---|---|
| 🛍️ | U+1F6CD U+FE0F | fin del saludo |
| 📍 | U+1F4CD | entrega con retiro |
| 📦 | U+1F4E6 | entrega con envío |
| 💬 | U+1F4AC | comentario |

**Si QA ve que algún cliente de WhatsApp muestra `�` en lugar de un emoji:** usar la variante sin emojis. Se borra cada emoji junto con el espacio que lo separa del texto (el de antes en el saludo; el de después en las otras líneas). El resto queda idéntico.

## 3. Reglas de limpieza y formato

- **Moneda.** No usar `Intl` (según el motor, puede escribir `4000` sin punto). Usar:
  ```ts
  '$' + Math.round(n).toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  ```
  Resultado: `$500`, `$4.000`, `$17.000`, `$150.000`, `$1.250.000`.
- **Texto libre** (nombre del cliente, producto, comuna, comentario, código):
  1. reemplazar `\r`, `\n` y `\t` por un espacio;
  2. borrar `*`, `_`, `~` y `` ` ``, porque rompen el formato de WhatsApp;
  3. juntar los espacios repetidos;
  4. quitar los espacios de los extremos.
- **Orden de las líneas:** el que devuelve la API.
- **Más de 25 líneas:** se muestran las 25 primeras y la línea `• … y {n} productos más (detalle completo en el pedido {numero})`. Con `n = 1` va en singular: `producto`. El total va siempre completo.

## 4. Implementación de referencia (TypeScript)

```ts
type Linea = { codigo: string; nombre: string; cantidad: number; precio_unitario: number;
  tipo_precio: 'detalle' | 'oferta' | 'mayor'; subtotal: number };
type PedidoCreado = { numero: string; nombre: string; entrega: 'retiro' | 'envio';
  region_nombre?: string; comuna?: string; agencia_nombre?: string; comentario?: string;
  total: number; es_demo: boolean; items: Linea[] };

const clp = (n: number) => '$' + Math.round(n).toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
const limpiar = (s = '') => s.replace(/[\r\n\t]+/g, ' ').replace(/[*_~`]/g, '').replace(/\s{2,}/g, ' ').trim();
const TIPO = { detalle: 'detalle', oferta: 'oferta', mayor: 'x mayor' } as const;
const MAX_LINEAS = 25;
const MAX_COMENTARIO = 300;

export function mensajePedido(p: PedidoCreado): string {
  const out: string[] = [];
  if (p.es_demo) out.push('*[PEDIDO DE PRUEBA · DEMO]*');
  out.push('Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️', '');
  out.push(`*Pedido ${p.numero}*`, '');
  out.push('*Productos*');
  for (const l of p.items.slice(0, MAX_LINEAS)) {
    out.push(`• ${l.cantidad} x ${limpiar(l.nombre)}`);
    out.push(`  Cód. ${limpiar(l.codigo)} · ${clp(l.precio_unitario)} c/u (${TIPO[l.tipo_precio]}) = ${clp(l.subtotal)}`);
  }
  const resto = p.items.length - MAX_LINEAS;
  if (resto > 0) out.push(`• … y ${resto} producto${resto === 1 ? '' : 's'} más (detalle completo en el pedido ${p.numero})`);
  out.push('');
  out.push(`*Total: ${clp(p.total)}*` + (p.entrega === 'envio' ? ' (no incluye envío)' : ''), '');
  if (p.entrega === 'retiro') {
    out.push('📍 *Entrega:* Retiro en tienda, Alameda 2432');
    const comuna = limpiar(p.comuna);
    if (comuna) out.push(`*Comuna:* ${comuna}`);
  } else {
    out.push(`📦 *Entrega:* Envío por agencia ${limpiar(p.agencia_nombre)}`);
    out.push(`*Región:* ${limpiar(p.region_nombre)}`);
    out.push(`*Comuna:* ${limpiar(p.comuna)}`);
  }
  out.push(`*Nombre:* ${limpiar(p.nombre)}`);
  let c = limpiar(p.comentario);
  if (c) {
    if (c.length > MAX_COMENTARIO) c = c.slice(0, MAX_COMENTARIO - 1).trimEnd() + '…';
    out.push(`💬 *Comentario:* ${c}`);
  }
  out.push('', 'Quedo esperando su confirmación por este chat.');
  return out.join('\n');
}
```

> **Contrato:** `region_nombre` y `agencia_nombre` salen de mapear los códigos de la API con la misma lista del selector. Si la API ya devuelve el nombre visible, se usa ese. Los nombres de campo pueden ajustarse a `docs/openapi.yaml`, pero el texto que produce la función no cambia.

## 5. Mensajes de consulta (sin pedido)

No llevan emojis. Usan el mismo número de destino.

| Dónde | Texto exacto |
|---|---|
| Botón flotante sin carrito, banner y canal oficial | `Hola Donde Carlitos Emporio, vengo de la web y tengo una consulta: ` (con un espacio al final, para que el cliente siga escribiendo) |
| Producto agotado ("¿Vuelve? Pregúntanos") | `Hola Donde Carlitos Emporio, vengo de la web. ¿Les queda {nombre_producto} (cód. {codigo})?` |
| Panel admin, "Abrir chat con el cliente" (opcional; destino = teléfono del cliente sin `+`) | `Hola {nombre_cliente}, te escribimos de Donde Carlitos Emporio por tu pedido {numero}.` |

Consulta general ya codificada:
```
https://wa.me/56921145368?text=Hola%20Donde%20Carlitos%20Emporio%2C%20vengo%20de%20la%20web%20y%20tengo%20una%20consulta%3A%20
```

## 6. Ejemplos exactos (salida esperada para QA)

> Pedidos ficticios. Los precios salen de `research/productos-semilla.json`. Los códigos de producto son `[EJEMPLO]`: QA los compara con los que tenga la semilla real.

### Ejemplo A: retiro, sin comuna, con comentario, `es_demo = false`
Líneas:
- 2 × squishy a $4.000 detalle;
- 1 × tren navideño a $17.000 oferta;
- 12 × camino de mesa a $800 por mayor.

Total: 8.000 + 17.000 + 9.600 = **34.600**.

```
Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️

*Pedido DC-0007*

*Productos*
• 2 x Squishy viral que cambia de color con el sol
  Cód. JUG-0001 · $4.000 c/u (detalle) = $8.000
• 1 x Tren navideño para árbol (grande)
  Cód. NAV-0012 · $17.000 c/u (oferta) = $17.000
• 12 x Camino de mesa navideño
  Cód. NAV-0009 · $800 c/u (x mayor) = $9.600

*Total: $34.600*

📍 *Entrega:* Retiro en tienda, Alameda 2432
*Nombre:* Camila Rojas
💬 *Comentario:* Paso el sábado en la mañana.

Quedo esperando su confirmación por este chat.
```

### Ejemplo B: envío, con comentario, `es_demo = true`
Líneas:
- 12 × squishy a $3.500 por mayor (caja de 12);
- 4 × dinosaurio a $14.000 por mayor de 4.

Total: 42.000 + 56.000 = **98.000**.

```
*[PEDIDO DE PRUEBA · DEMO]*
Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️

*Pedido DC-0008*

*Productos*
• 12 x Squishy viral que cambia de color con el sol
  Cód. JUG-0001 · $3.500 c/u (x mayor) = $42.000
• 4 x Dinosaurio gigante con ruedas
  Cód. JUG-0010 · $14.000 c/u (x mayor) = $56.000

*Total: $98.000* (no incluye envío)

📦 *Entrega:* Envío por agencia Starken
*Región:* Ñuble
*Comuna:* Chillán
*Nombre:* Pedro Soto
💬 *Comentario:* Es para mi negocio. ¿Lo pueden despachar el lunes?

Quedo esperando su confirmación por este chat.
```

### Ejemplo C: retiro con comuna, sin comentario, con texto que hay que limpiar
- **Entrada:** nombre `"  Ana  *Pérez*\n"`, comuna `"Maipú"`, comentario vacío, 1 × esfera navideña con luces a $500 detalle.
- **Qué se verifica:** que el nombre sale limpio, que aparece `*Comuna:*` y que **no** aparece la línea de comentario.

```
Hola Donde Carlitos Emporio, les envío mi pedido de la web 🛍️

*Pedido DC-0009*

*Productos*
• 1 x Esfera navideña con luces
  Cód. NAV-0002 · $500 c/u (detalle) = $500

*Total: $500*

📍 *Entrega:* Retiro en tienda, Alameda 2432
*Comuna:* Maipú
*Nombre:* Ana Pérez

Quedo esperando su confirmación por este chat.
```

## 7. Checklist de QA

- [ ] El texto generado es idéntico, byte a byte, a los ejemplos A, B y C con los mismos datos.
- [ ] El total del mensaje coincide con el total del pedido en el panel.
- [ ] El total es la suma de los subtotales y cada subtotal es cantidad × precio unitario.
- [ ] Los montos van con punto de miles (`$4.000`, no `$4000` ni `$4,000`).
- [ ] No hay doble línea en blanco ni espacios al final de las líneas.
- [ ] Al abrir `wa.me` en Android e iOS, desde el navegador interno de TikTok e Instagram, el mensaje llega completo, con saltos de línea, negritas y emojis bien.
- [ ] Probado con 25 o más líneas: aparece la línea "… y N productos más".
- [ ] "Copiar mensaje" copia el mismo texto.
- [ ] En el demo, el número de destino **no** es el del negocio (ver §1).
- [ ] El mensaje no incluye el teléfono del cliente ni ningún dato que no esté en la plantilla.
