# API — Donde Carlitos Emporio

API en Go (`net/http` de la librería estándar) + SQLite (`modernc.org/sqlite`, sin CGO). Contrato: `docs/openapi.yaml`. Modelo: `docs/modelo-datos.md`.

## Estructura
```
cmd/server         servidor HTTP (+ flag -healthcheck para Docker)
cmd/setup-admin    genera la clave del panel y escribe api/.env y api/ADMIN_PASSWORD.local.txt
internal/http      rutas, middleware (seguridad, límites, origen), handlers, DTOs  [paquete httpapi]
internal/productos precio por cantidad, disponibilidad, badges, búsqueda, validación del admin
internal/pedidos   validación del pedido, cotización, estados, teléfono, mensaje de WhatsApp
internal/auth      bcrypt, tokens de sesión, bloqueo de login
internal/store     interfaz Store + tipos; store/sqlite = implementación
internal/seed      semilla embebida (semilla.json) y cargador
internal/texto     limpieza y validación de texto libre
migrations         SQL versionado, embebido con embed
```

## Levantar y probar (sin Go instalado)
```sh
docker compose up -d --build api              # solo la API (web/ la construye constructor-web)
curl -s http://127.0.0.1:8080/healthz
curl -s "http://127.0.0.1:8080/api/productos?categoria=navidad&oferta=true"
curl -s -X POST http://127.0.0.1:8080/api/pedidos \
  -H "Origin: http://localhost:3000" -H "Content-Type: application/json" \
  -d '{"nombre":"Prueba Demo","telefono":"+56 9 0000 0009","entrega":"retiro","items":[{"producto_id":1,"cantidad":2}]}'
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
| `TRUSTED_PROXIES` | vacío | CIDRs desde donde se cree `X-Forwarded-For` (compose: redes privadas = contenedor web) |
| `WHATSAPP_NUMERO` | `56921145368` | destino del link wa.me (público, no secreto) |
| `TIENDA_DIRECCION` | `Alameda 2432, Santiago Centro` | dirección de retiro en `/api/opciones-entrega` |
| `PEDIDOS_DEMO` | `true` | marca los pedidos web como DEMO: el mensaje parte con `*[PEDIDO DE PRUEBA · DEMO]*` |
| `SEED_DEMO` / `SEED_PEDIDOS_DEMO` | `true` / `true` | siembra catálogo y 6 pedidos DEMO si la base está vacía |
| `DIAS_RECIEN_LLEGADO` | `14` | ventana del badge `recien_llegado` |
| `PEDIDOS_POR_IP` / `PEDIDOS_VENTANA_MIN` | `10` / `10` | límite de `POST /api/pedidos` por IP |
| `LOG_LEVEL` | `info` | |

## Decisiones
- **Precio recalculado en el servidor.** `POST /api/pedidos` acepta `total` y `precio_unitario` del navegador pero los ignora (se loguea si difieren). La cotización corre dentro de la misma transacción que descuenta stock e inserta el pedido.
- **Precio por mayor por unidad y con mínimo.** Solo aplica si `cantidad ≥ mayor_minimo`. Se aplica el menor entre detalle, oferta y mayor. Un solo escalón por producto (los videos a veces dicen dos; el segundo va en `descripcion`). Ver `docs/modelo-datos.md`.
- **Correlativo atómico sin huecos.** Contador en tabla, `UPDATE … RETURNING` dentro de una transacción `BEGIN IMMEDIATE` (DSN `_txlock=immediate`, `busy_timeout`). Probado con 60 pedidos concurrentes.
- **Stock se reserva al crear el pedido** (con `UPDATE … WHERE stock >= ?`, sin sobreventa; probado con 20 pedidos concurrentes por 5 unidades) y se devuelve al cancelar. Alternativa descartada: descontar al confirmar (permitiría confirmar sin stock).
- **Estados**: `listo` solo para retiro y `despachado` solo para envío; `despachado` no se cancela. El pedido trae `transiciones_permitidas`.
- **Mensaje de WhatsApp** generado por la API, idéntico a `propuesta/mensaje-pedido-whatsapp.md` (tests con los ejemplos A, B y C byte a byte). La respuesta también trae todos los campos para que el front lo arme igual si prefiere.
- **Sesiones**: token aleatorio de 256 bits en cookie HttpOnly + SameSite=Strict (+ Secure y prefijo `__Host-` en producción). En la base solo se guarda HMAC-SHA256(SESSION_SECRET, token). Vencen a las 12 h o con 2 h sin uso.
- **CSRF**: SameSite=Strict + verificación de `Origin` (y rechazo de `Sec-Fetch-Site: cross-site`) en todo método que cambia datos, incluidos los públicos.
- **Login**: bcrypt costo 12; siempre se corre bcrypt (hash señuelo si el usuario no existe); 5 fallos en 15 min por IP → bloqueo de 15 min (también con la clave correcta). IPv6 se agrupa por /64.
- **IP del cliente**: `X-Forwarded-For` solo se cree si la conexión viene de `TRUSTED_PROXIES`, recorriéndolo de derecha a izquierda. Así el límite por IP funciona detrás del rewrite de Next.js sin que un cliente lo falsifique.
- **Búsqueda**: columna `busqueda` normalizada (minúsculas, sin tildes) + `LIKE ? ESCAPE '\'` con `%` `_` `\` escapados. Máx. 60 caracteres y 5 palabras.
- **Validación**: JSON estricto (`DisallowUnknownFields`, un solo objeto, `Content-Type` JSON), body máx. 64 KB global y 16 KB en pedidos, texto sin caracteres de control ni de reordenamiento bidireccional, URLs de imagen solo `/…` o https, TikTok solo https a tiktok.com.
- **Servidor**: `ReadHeaderTimeout` 5 s, `ReadTimeout` 10 s, `WriteTimeout` 15 s, `IdleTimeout` 60 s, cabeceras máx. 16 KB, apagado ordenado con SIGTERM.
- **Cabeceras** en toda respuesta: `X-Content-Type-Options`, `Referrer-Policy: no-referrer`, `X-Frame-Options: DENY`, CSP `default-src 'none'; frame-ancestors 'none'`, CORP/COOP, `Permissions-Policy`, `Cache-Control: no-store`; HSTS en producción. Sin CORS (mismo origen).
- **Logs** JSON (`log/slog`): sin cuerpos, cookies ni nombre del cliente; teléfono enmascarado (`+569****5678`) e IP truncada a /24.
- **Errores**: genéricos hacia el cliente (`{"error":{"codigo","mensaje"}}`), detalle solo en el log.
- **Docker**: build `golang:1.27.4-alpine` y runtime `distroless/static-debian13:nonroot`, ambos fijados por digest; usuario 65532, sistema de archivos de solo lectura, `cap_drop: ALL`, `no-new-privileges`; healthcheck con el propio binario (distroless no tiene curl).
- **Imágenes de productos**: `imagen_url = /productos/tt-<videoId>.webp`; los archivos los pone constructor-web en `web/public/productos/` a partir de `research/raw/tt-<videoId>.png`. 9 productos no tienen captura y van con `imagen_url: null`.

## Semilla
`internal/seed/semilla.json` se generó el 2026-11-14 desde `research/productos-semilla.json` (37 productos con precio y video reales; 3 omitidos por no tener precio). Lo que no tiene fuente va con `[EJEMPLO]` en la descripción: stock/disponibilidad y los mínimos por mayor no dichos (se usó 6). Las ofertas del video ("antes $X") son `precio_oferta`. `fecha_llegada` = fecha del video *(inferencia)*. Los códigos `P0xx` son internos del demo. Pedidos DEMO: "Cliente Demo Uno…Seis", `+56 9 0000 000N`.

## Pendiente / siguiente paso
- Escalones múltiples de precio por mayor (tabla `precios_mayor`), si el dueño los confirma.
- Temporada activa configurable desde el panel (hoy el front filtra con `temporada=navidad`).
- Retiro en el local de Temuco, si el dueño lo quiere (hoy solo Alameda 2432).
- Límites por IP en memoria: con más de una instancia habría que moverlos a la base o a Redis.
- PostgreSQL: nueva implementación de `store.Store`.
