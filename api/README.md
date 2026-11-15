# API — Donde Carlitos Emporio

API en Go (`net/http` de la librería estándar) + SQLite (`modernc.org/sqlite`, sin CGO). Contrato: `docs/openapi.yaml` (**v2**). Modelo: `docs/modelo-datos.md`.

**v2 (2026-11-14):** no hay carrito ni pedido público. La API sirve la landing (vitrina, config pública y seguimiento sin datos personales) y el **panel de gestión de pedidos**. El equipo registra ahí los pedidos mientras atiende los WhatsApp de Santiago y Regiones. El panel tiene varios usuarios, pantalla Hoy, búsqueda, clientes, resumen, CSV y mensajes listos según `propuesta/mensajes-whatsapp.md`.

## Estructura
```
cmd/server         servidor HTTP (+ flag -healthcheck para Docker)
cmd/setup-admin    genera la clave del panel y escribe api/.env y api/ADMIN_PASSWORD.local.txt
internal/http      rutas, middleware (seguridad, límites, origen), handlers, DTOs  [paquete httpapi]
internal/productos precio por cantidad, disponibilidad, badges, búsqueda, validación del admin
internal/pedidos   pedido del panel: validación, cotización (catálogo y líneas libres), estados, teléfono
internal/mensajes  plantillas (d1)–(d8) de WhatsApp al cliente (plantillas.go) y links wa.me
internal/limite    LRU, ventana y semáforo para límites en memoria
internal/auth      bcrypt, tokens de sesión, bloqueo de login
internal/store     interfaz Store + tipos; store/sqlite = implementación
internal/seed      catálogo (semilla.json), pedidos DEMO (pedidos_demo.json) y cargadores
internal/texto     limpieza y validación de texto libre
migrations         SQL versionado, embebido con embed
```

## Levantar y probar (sin Go instalado)
Por defecto la API **no se publica** en la máquina: solo la web (`127.0.0.1:${WEB_PORT}`) le habla por la red interna. Para depurarla directo, el override `docker-compose.debug.yml` la publica en `127.0.0.1:8080`:
```sh
WEB_PORT=3100 docker compose -f docker-compose.yml -f docker-compose.debug.yml up -d --build
curl -s http://127.0.0.1:8080/healthz
curl -s "http://127.0.0.1:8080/api/productos?categoria=navidad&oferta=true"
curl -s http://127.0.0.1:8080/api/config-publica
docker compose down -v                        # borra la base: al volver a levantar se re-siembra
```
Las escrituras exigen `Origin` permitido (en el navegador va solo; con curl hay que ponerlo).

Calidad, en un contenedor (Git Bash: anteponer `MSYS_NO_PATHCONV=1`):
```sh
docker run --rm -v "$PWD/api:/src" -w /src golang:1.27.4-alpine sh -c \
  'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./...'
```

## Clave del panel (setup-admin)
Se corre **una vez** (o con `-rotar` para cambiarla). Escribe en `api/.env` `ADMIN_USUARIO`, `ADMIN_PASSWORD_HASH_B64` y `SESSION_SECRET`, y la clave en texto plano en `api/ADMIN_PASSWORD.local.txt`. **No imprime ningún secreto.** Conserva otras variables que ya estén en el `.env`.
```sh
docker compose build api
docker compose run --rm --no-deps -v "./api:/work" --entrypoint /app/setup-admin api -dir /work
docker compose up -d api      # la API lee el .env al arrancar
```
- PowerShell: igual. En Git Bash: `MSYS_NO_PATHCONV=1 docker compose run ... -v "$PWD/api:/work" ...`.
- En Linux, si el contenedor (usuario 65532) no puede escribir en `api/`, agregar `--user "$(id -u):$(id -g)"`.
- Rotar: `... --entrypoint /app/setup-admin api -dir /work -rotar` y reiniciar la API. Al cambiar el hash se cierran las sesiones abiertas.

### Por qué el hash va en base64
Un hash bcrypt (`$2a$12$…`) tiene `$`, y docker compose interpola `$` en `env_file` (desde Compose v2 `env_file` también interpola), lo que corrompe el valor sin avisar. `setup-admin` guarda `ADMIN_PASSWORD_HASH_B64` = base64 estándar del hash (sin `$`), y el servidor lo decodifica y valida que sea bcrypt al arrancar. `SESSION_SECRET` es hex. Un test verifica que el `.env` generado no tenga `$`.

## Variables de entorno
| Variable | Defecto | Uso |
|---|---|---|
| `APP_ENV` | `development` | `production` → cookie `__Host-dc_sesion` con `Secure` y HSTS |
| `LISTEN_ADDR` | `:8080` | |
| `DB_PATH` | `/data/carlitos.db` | en el volumen `datos-api` |
| `ADMIN_USUARIO` | `admin` | secreto/archivo `.env` |
| `ADMIN_PASSWORD_HASH_B64` | — | sin esto (o sin `SESSION_SECRET`) el panel queda deshabilitado (`503 admin_no_configurado`) |
| `SESSION_SECRET` | — | hex ≥ 32 bytes; clave del HMAC de los tokens de sesión |
| `ALLOWED_ORIGINS` | `http://localhost:3000,http://127.0.0.1:3000` | orígenes aceptados en POST/PUT/PATCH/DELETE |
| `TRUSTED_PROXIES` | vacío | CIDRs desde donde se cree el **último** valor de `X-Forwarded-For`. Compose: solo `172.30.57.10/32` (IP fija del contenedor web). Vacío = se limita por la IP de la conexión |
| `IMAGENES_DIR` | `/imagenes` | carpeta (volumen `imagenes-api`) de las imágenes subidas desde el panel |
| `LOGIN_GLOBAL_MIN` | `30` | intentos de login por minuto (todas las IP); además máximo 2 bcrypt simultáneos |
| `SEGUIMIENTO_POR_IP_MIN` | `30` | consultas a `/api/seguimiento/{token}` por minuto por IP (y 600/min en total) |
| `WHATSAPP_SANTIAGO` | `56921145368` | WhatsApp de Santiago/retiro (público); sale en `/api/config-publica` |
| `WHATSAPP_REGIONES` | vacío | WhatsApp de regiones. Vacío = [CONFIRMAR]: `whatsapp_regiones: null` y la landing usa el de Santiago |
| `URL_SITIO` (o `SITE_URL`) | `http://localhost:3000` | base del link de seguimiento (`/seguimiento/<token>`) |
| `DIRECCION_RETIRO` | `Alameda 2432, Santiago Centro` | mensajes d1/d5 y config pública |
| `HORARIO_RETIRO` | `L a V 9:00 a 18:00 · Sáb, dom y feriados 9:00 a 16:00` | línea 🕘 de d5 (definida pero vacía = se omite). Por confirmar |
| `HORARIO_CONFIRMADO` | `false` | la landing lo muestra como "por confirmar" mientras sea false |
| `DATOS_TRANSFERENCIA` | vacío | datos para pagar que escribe el dueño. Van **solo** en `api/.env`, nunca en la web. Una línea, con `\n` literales para los saltos. Vacío = la cotización omite el bloque y no se ofrece el mensaje d2. **Nunca se inventa** |
| `SEED_DEMO` / `SEED_PEDIDOS_DEMO` | `true` / `true` | catálogo si no hay categorías; 20 pedidos DEMO y 2 usuarios demo sin clave si no hay pedidos |
| `DIAS_RECIEN_LLEGADO` | `14` | ventana del badge `recien_llegado` |
| `MODO_DEMO` | `false` (compose del demo: `true`) | todo pedido nuevo nace `es_demo`: mensajes con `*[PRUEBA · DEMO]*` y `wa.me` sin número. En producción, `false` |
| `STOCK_ESTRICTO` | `false` | `false`: pagar o preparar con stock insuficiente se acepta, el stock queda en 0 y la respuesta trae `aviso_stock`. `true`: `409 stock_insuficiente` |
| `LOG_LEVEL` | `info` | |

## Decisiones
- **Sin pedido público (v2).** Se eliminaron `POST /api/pedidos`, sus límites por IP y por teléfono, el tope de pedidos abiertos y el mensaje de pedido web. El pedido lo registra el equipo en el panel.
- **Precios y total en el servidor.** Las líneas del catálogo toman el precio de la base y lo congelan: detalle, oferta vigente o, solo desde `mayor_minimo`, precio por mayor. Las líneas libres llevan precio opcional. Si todas tienen precio, `total = Σ − descuento + costo_envio` y, si el panel manda `total`, debe coincidir. Si alguna va sin precio, el total lo escribe quien atiende (`total_manual`).
- **Precio por mayor por unidad y con mínimo.** Un solo escalón por producto. Los videos a veces dicen dos: el segundo va en `descripcion`, o el panel usa una línea libre.
- **Correlativo atómico sin huecos.** Contador en tabla, `UPDATE … RETURNING` dentro de una transacción `BEGIN IMMEDIATE`. Probado con 60 pedidos concurrentes. La semilla lo deja en 89, así el primer pedido nuevo del demo es DC-0110, como dice el guion.
- **Estados y stock.**
  - Flujo: cotizado → esperando_pago → pagado → preparando → despachado | listo_retiro → entregado; cancelado desde cualquier estado no final.
  - Pagado exige pago y `confirmo_pago_recibido`; despachado exige agencia, n.º y comuna.
  - El stock se descuenta al pagar o al empezar a preparar (retiro con pago en tienda), con `UPDATE … WHERE stock >= ?`. Probado con 20 pagos simultáneos por 5 unidades: pasan 5. Se devuelve al cancelar.
  - Solo se puede editar antes de pagar.
- **Mensajes al cliente** (`internal/mensajes/plantillas.go`): texto exacto de (d1)–(d8), con tests byte a byte de todos los ejemplos del documento. `CLP()` separa los miles a mano. En pedidos `es_demo`, el link es `https://wa.me/?text=…`, sin número.
- **Seguimiento público.**
  - Token de 16 bytes en base64url (no el `DC-xxxx`).
  - Muestra solo número, estado, línea de tiempo, tipo de entrega, agencia y n.º. Sin nombre, teléfono, comuna, montos ni ítems.
  - 404 genérico, límite por IP y global, `noindex`.
  - *Nota:* `propuesta/mensajes-whatsapp.md` §5 sugiere mostrar también productos y total. Se siguió la regla más estricta del orquestador (sin montos ni ítems); agregarlos es un cambio pequeño si se decide.
- **Usuarios y roles.**
  - `admin` y `vendedor`. El vendedor ve pedidos y clientes, y ve productos sin editarlos; no ve usuarios, resumen ni CSV (`403 sin_permiso`).
  - Las claves de usuarios nuevos las genera la API (16 caracteres) y se muestran una sola vez; solo se guarda el bcrypt. Desactivar o regenerar la clave cierra sus sesiones.
  - El admin de `api/.env` se re-sincroniza al arrancar (rol admin y activo): es la vía de recuperación.
- **CSV:** UTF-8 con BOM y separador `;`. Las celdas de texto que empiezan con `=`, `+`, `-`, `@`, tab o CR se anteponen con un apóstrofo (CSV injection). Solo admin. Máx. 10.000 filas.
- **Sesiones**: token aleatorio de 256 bits en cookie HttpOnly + SameSite=Strict (+ Secure y prefijo `__Host-` en producción). En la base solo se guarda HMAC-SHA256(SESSION_SECRET, token). Vencen a las 12 h o con 2 h sin uso.
- **CSRF**: SameSite=Strict + verificación de `Origin` (y rechazo de `Sec-Fetch-Site: cross-site`) en todo método que cambia datos, incluidos los públicos.
- **Login**: bcrypt costo 12; siempre se corre bcrypt (hash señuelo si el usuario no existe); 5 fallos en 15 min por IP → bloqueo de 15 min (también con la clave correcta). Tope global de 30 intentos/min y máximo 2 bcrypt simultáneos (`429` + `Retry-After`), para que un flood desde muchas IPs no acapare la CPU. IPv6 se agrupa por /64.
- **IP del cliente** (auditoría A1): la API confía solo en el contenedor web, con IP fija `172.30.57.10` en la subred `172.30.57.0/24`, y de `X-Forwarded-For` toma **solo el último valor**, que es el que agrega ese proxy. Si la conexión no viene de un proxy confiable, se usa la IP de la conexión y se ignora XFF. **En producción debe haber un proxy de borde (Caddy, Nginx, el balanceador) que sobrescriba `X-Forwarded-For` con la IP real**; si solo agrega, un cliente puede inventar valores a la izquierda, y por eso se toma solo el último.
- **Limitadores** (login, seguimiento): LRU de tamaño fijo. Un flood de claves nuevas desaloja las más antiguas en vez de dejar afuera a todos.
- **Imágenes** (`POST /api/admin/productos/{id}/imagen`): solo JPG/PNG/WebP por firma de bytes, máx. 5 MB y 5000 px por lado (evita bombas de descompresión), 2 subidas simultáneas. Se re-codifican (quita EXIF/GPS y contenido políglota) y se guardan con nombre aleatorio en el volumen `imagenes-api`. Se sirven en `GET /api/imagenes/{archivo}` (solo nombres `[a-f0-9]{32}.jpg|png`) con Content-Type fijo, nosniff y CSP `default-src 'none'`.
- **Ofertas con vencimiento**: `oferta_hasta` se compara con la fecha de hoy en Chile. Vencida → no cuenta en precio, badges, filtro ni pedidos; en la API pública `precio_oferta` va en null.
- **UTC**: el reloj del servidor y todas las fechas-hora de la API están en UTC; las fechas simples (`fecha_llegada`, `oferta_hasta`) son días de Chile.
- **Búsqueda**: columna `busqueda` normalizada (minúsculas, sin tildes) + `LIKE ? ESCAPE '\'` con `%` `_` `\` escapados. Máx. 60 caracteres y 5 palabras.
- **Validación**: JSON estricto (`DisallowUnknownFields`, un solo objeto, `Content-Type` JSON), body máx. 64 KB (5 MB en la subida de imagen), texto sin caracteres de control ni de reordenamiento bidireccional, URLs de imagen solo `/…` o https, TikTok solo https a tiktok.com.
- **Servidor**: `ReadHeaderTimeout` 5 s, `ReadTimeout` 10 s, `WriteTimeout` 15 s, `IdleTimeout` 60 s, cabeceras máx. 16 KB, apagado ordenado con SIGTERM.
- **Cabeceras** en toda respuesta: `X-Content-Type-Options`, `Referrer-Policy: no-referrer`, `X-Frame-Options: DENY`, CSP `default-src 'none'; frame-ancestors 'none'`, CORP/COOP, `Permissions-Policy`, `Cache-Control: no-store`; HSTS en producción. Sin CORS (mismo origen).
- **Logs** JSON (`log/slog`): sin cuerpos, cookies ni nombre del cliente; teléfono enmascarado (`+569****5678`) e IP truncada a /24. Las claves generadas nunca se registran.
- **Errores**: genéricos hacia el cliente (`{"error":{"codigo","mensaje"}}`), detalle solo en el log.
- **Docker**: build `golang:1.27.4-alpine` y runtime `distroless/static-debian13:nonroot`, ambos fijados por digest; usuario 65532, sistema de archivos de solo lectura (escribe solo en los volúmenes `/data` e `/imagenes` y en `/tmp` tmpfs), `cap_drop: ALL`, `no-new-privileges`; healthcheck con el propio binario (distroless no tiene curl). La web también corre `read_only` con tmpfs en `/tmp` y `/app/.next/cache`.
- **Imágenes de productos**: `imagen_url = /productos/tt-<videoId>.webp`; los archivos los pone constructor-web en `web/public/productos/` a partir de `research/raw/tt-<videoId>.png`. 9 productos no tienen captura y van con `imagen_url: null`.

## Semilla
`internal/seed/semilla.json` se generó el 2026-11-14 desde `research/productos-semilla.json` (37 productos con precio y video reales; 3 omitidos por no tener precio). Lo que no tiene fuente va con `[EJEMPLO]` en la descripción: stock/disponibilidad y los mínimos por mayor no dichos (se usó 6). Las ofertas del video ("antes $X") son `precio_oferta`. `fecha_llegada` = fecha del video *(inferencia)*. Los códigos `P0xx` son internos del demo.

`internal/seed/pedidos_demo.json` trae 20 pedidos ficticios:
- DC-0101 a DC-0109 son los del guion (`propuesta/pitch.md` §c), con los mismos clientes, productos y estados.
- DC-0090 a DC-0100 son pedidos anteriores ya cerrados.
- Nombres "… Demo" y teléfonos `+56 9 0000 0xxx`. Usuarios demo `regiones1` y `santiago1`, sin clave utilizable.

Diferencias con el guion:
- DC-0102 usa una línea libre a $12.500, porque el catálogo no modela ese segundo escalón del video.
- Si la semilla se carga después del 15/11, DC-0103 y DC-0109 salen a $20.000 (vence la oferta del tren).

## Ronda 3 (contrato 2.1.0)
- **Despacho directo:** `pagado` → `despachado` o `listo_retiro` es una **transición directa**. No se registra un "preparando" implícito: el historial muestra lo que pasó. `preparando` sigue disponible para quien lo use.
- **No se cancela un despachado** (`arquitectura-landing.md` C3: la mercadería ya salió).
- **Resumen:**
  - "vendido" = pedidos con pago registrado y no cancelados; `pedidos` va aparte;
  - `por_comuna` muestra las 20 comunas con más pedidos;
  - `por_usuario` = registrados y despachados por persona.
- **`GET /api/admin/ajustes`** (solo admin, solo lectura): números, dirección, horario, URL y si hay datos de transferencia (nunca el texto) y modo demo.
- **Promo por comentario:** `precio_es_promo` + `nota_promo`. La API no arma ningún texto con "detalle" para esos productos; la web muestra la placa "PROMO DEL VIDEO".
- **Vitrina:** `en_vitrina` + `orden_vitrina`, máximo 8 (`409 vitrina_llena`). `GET /api/vitrina`.
- **Precio editable en línea de catálogo** (`tipo_precio: manual`).
- **CSV:** columnas de `copy-landing.md` P5, hora de Chile `dd-mm-aaaa hh:mm`, nombres visibles y teléfono como `569…`.
- **Búsqueda de productos:** primero los que empiezan con la primera palabra buscada, después por posición.
- **Seguridad:**
  - **Login (M4):** espera progresiva por usuario + IP: 1 s, 2 s, 4 s… hasta 5 min. Se olvida a los 15 min sin intentos o con un login correcto. Desde 5 fallos seguidos queda una línea `ALERTA` en el log.
  - **Seguimiento (I-5):** solo cuentan las consultas que terminan en 404, así quien tiene un link válido nunca queda bloqueado. El límite por IP solo aplica si hay proxy confiable; sin proxy, solo el global (600 fallos/min).
  - **Admin principal protegido (B8):** el `ADMIN_USUARIO` de `api/.env` solo se edita a sí mismo; desde otro usuario → `403 admin_principal_protegido`.
  - **Claves temporales (B9):** los usuarios creados o con clave regenerada nacen con `debe_cambiar_clave`. Hasta cambiarla (`POST /api/admin/sesion/clave`, mínimo 12 caracteres, cierra sus otras sesiones), el resto del panel responde `403 debe_cambiar_clave`. El admin de `api/.env` no la tiene: su clave la genera `setup-admin`.
  - **Link de seguimiento (B10):** vence 30 días después de entregado o cancelado (404 genérico). `POST /api/admin/pedidos/{id}/token` lo regenera.
  - **Límites de contenedor (B11):** API 256 MB y 100 procesos; web 512 MB y 200.

### Checklist antes de mostrar el demo (M-14)
- Si alguien probó claves malas, el login de ese usuario puede pedir esperar unos minutos. Para limpiarlo al tiro: `WEB_PORT=3100 docker compose restart api`. Los límites viven en memoria; reiniciar no toca la base.
- Ensayar el guion dos veces no rompe nada: el stock no bloquea (`STOCK_ESTRICTO=false`) y los productos del guion tienen stock 300.
- Para volver al estado inicial del guion (DC-0101…DC-0109 como al principio): `docker compose down -v` y levantar de nuevo. **Ojo:** esto borra la base, incluido el admin; hay que volver a correr `setup-admin`.

### Producción: proxy de borde (M4)
En el demo local la API ve todo el tráfico con la IP del contenedor web, y por eso los límites por IP los comparte todo el mundo. Para producción:
1. Poner un proxy de borde delante de `web` (Caddy o nginx, imagen fijada por digest, usuario no root). Tiene que **sobrescribir** `X-Forwarded-For` con la IP real, no agregarle: en nginx, `proxy_set_header X-Forwarded-For $remote_addr;`. Caddy lo hace por defecto si no se le declaran `trusted_proxies`.
2. Publicar solo ese proxy (443). `web` y `api` quedan en la red interna.
3. En la API: `TRUSTED_PROXIES=172.30.57.10/32` (la IP fija de `web`) y en la web reenviar el XFF que recibe. Así, la API cree solo el último valor, que es el que puso el borde.
4. `APP_ENV=production` (cookie `__Host-` + `Secure`, HSTS), `MODO_DEMO=false` y `URL_SITIO` con el dominio real.

## Pendiente / siguiente paso
- Escalones múltiples de precio por mayor (tabla `precios_mayor`), si el dueño los confirma.
- Temporada activa configurable desde el panel (hoy el front filtra con `temporada=navidad`).
- Retiro en el local de Temuco, si el dueño lo quiere (hoy solo Alameda 2432).
- Límites por IP en memoria: con más de una instancia habría que moverlos a la base o a Redis.
- Ajustes editables desde el panel (hoy son de solo lectura y salen de variables de entorno).
- Plazo de vencimiento del link de seguimiento: 30 días por defecto; la plantilla deja `[COMPLETAR LUCKAS]`.
- PostgreSQL: nueva implementación de `store.Store`.
