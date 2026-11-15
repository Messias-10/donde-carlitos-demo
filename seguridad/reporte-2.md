# Auditoría de seguridad — Fase 4, ronda 2 (v2: landing + seguimiento + panel con roles)

- **Fecha:** 2026-11-15
- **Autor:** auditor-seguridad (guardado por el orquestador; el subagente no puede escribir reportes).
- **Alcance:** árbol e historial git (14 commits, hasta `b2e1ff1`); `api/` (rutas `/api/admin/*`, seguimiento, imágenes, usuarios, CSV, mensajes); `web/` (Next.js 16.4.6); imágenes `donde-carlitos-api:local` y `donde-carlitos-web:local`; `docker-compose.yml` y `docker-compose.debug.yml`; git y GitHub.
- **Entorno:** solo localhost. Stack principal (`WEB_PORT=3100`): solo GET de cabeceras y bundles. Instancia aislada `aud2-*` (127.0.0.1:18081 y 18100, base en tmpfs, admin de prueba fuera del repo, `DATOS_TRANSFERENCIA` con valor marcador), borrada al terminar. No se leyeron `api/.env` ni `*.local.txt`. No se envió ningún WhatsApp.
- **Salidas:** `seguridad/salidas/ronda-2/` (secretos, cookies y claves redactados).

## Veredicto

**Sin críticos ni altos.** M4 no bloquea el demo local, pero **sí bloquea ir a producción**.

| Severidad | Cantidad | Detalle |
|---|---|---|
| Crítica | 0 | |
| Alta | 0 | |
| Media | 2 | M4 (nuevo); M3 abierto (decisión de Luckas) |
| Baja | 7 | B3 y B7 siguen; B8–B12 nuevos |
| Informativo | 3 | I1–I3 |

## Estado de los hallazgos de la ronda 1

| # | Estado | Evidencia |
|---|---|---|
| A1 — XFF falsificable | **Resuelto** (efecto secundario: M4) | `TRUSTED_PROXIES: ""`; solo se cree el último XFF y solo desde proxy confiable. 6 logins con XFF rotado: `401×5, 429` (§7, §14). |
| A2 — pedido anónimo agota stock | **Resuelto por diseño** | No existe `POST /api/pedidos`; el stock se descuenta al pasar a `pagado`. |
| M1 — login sin tope global | **Resuelto** | 30/min global, semáforo de 2 bcrypt, LRU de 10.000. 12 logins simultáneos: `429×10`, `401×2` (§15). |
| M2 — npm en la imagen web | **Resuelto** | Trivy 0 vulnerabilidades. |
| M3 — repo público con fotos y propuesta | **Abierto** (decisión de Luckas) | `visibility: PUBLIC`; la v2 ya está publicada. |
| B1 — web sin `read_only` | **Resuelto** | `ReadonlyRootfs=true`, tmpfs, `cap_drop ALL`, `no-new-privileges`. |
| B2 — HSTS en web | **Resuelto en código** | No probado con https real. |
| B3 — `style-src 'unsafe-inline'` | **Abierto, aceptado** | Riesgo bajo: `script-src` con nonce y `strict-dynamic`. |
| B4 — iframe TikTok | **Resuelto** | Sin iframes: `frame-src 'none'`. |
| B5 — API publicada | **Resuelto** | Solo con `docker-compose.debug.yml` (127.0.0.1). |
| B6 — MCP `@latest` y deny | **Resuelto** | `@playwright/mcp@0.0.85`; deny ampliado; hook `pre-push`. |
| B7 — correo personal | **Abierto** | En `CLAUDE.md`, `PLAN_MAESTRO.md` y los 14 commits. |

## Hallazgos

### M4 — Media — Los límites por IP los comparte todo el tráfico (dueños: ingeniero-go y constructor-web)
- **Ubicación:** `docker-compose.yml:39`; `api/internal/http/sesion.go:62-68`; `api/internal/http/seguimiento.go:38-49`; `web/lib/api-servidor.ts` (`seguimiento()` no reenvía la IP del cliente); `web/app/(tienda)/seguimiento/[token]/page.tsx` (un 429 termina en 500).
- **Evidencia (§7, §14):** la API solo ve la IP del contenedor web. Login: 6 fallos anónimos (`401×5, 429`) y luego el admin con la clave correcta recibe `429 login_bloqueado` (~15 min). Seguimiento: 40 consultas a tokens inexistentes → 28 pasan, 12 `429`; después un cliente con token válido recibe `429` y la página responde **500**.
- **Impacto:** con 5 peticiones cada 15 min, sin claves, se bloquea el panel para todo el equipo; con 30/min se cae el seguimiento de todos. Bajo en local; explotable desde internet en producción.
- **Corrección:**
  1. (ingeniero-go) Proxy de borde (Caddy o nginx, imagen fijada, no root) delante de `web` que **sobrescriba** XFF; publicar solo ese proxy en 127.0.0.1 y `TRUSTED_PROXIES` con su /32. Si no se pone en local, documentarlo y aplicar el punto 3.
  2. (constructor-web) En `seguimiento()` reenviar el XFF recibido por la página.
  3. (ingeniero-go) Sin proxies confiables, no aplicar límite por IP al seguimiento (solo el global) o subirlo: el token de 128 bits ya impide enumerar.
  4. (ingeniero-go) Bloqueo de login por **usuario + IP** con espera progresiva (1 s → 30 s, tope de pocos minutos) y alerta en el log.
  5. (constructor-web) En seguimiento, un 429 se muestra como "Muchas consultas, intenta en un minuto" (429/503, no 500).
  6. Verificar repitiendo §7 y §14.

### M3 — Media — Repo público con fotos de personas y la propuesta (dueño: Luckas)
Sin cambios. Volver a privado o dejar por escrito que se acepta el riesgo.

### Bajos
| # | Dueño | Hallazgo | Corrección |
|---|---|---|---|
| B8 | ingeniero-go | Un admin puede regenerar la clave del admin principal (el de `api/.env`) y tomar su cuenta (§10) | 403 al actuar sobre `ADMIN_USUARIO` desde otro usuario, o pedir la clave actual de quien opera |
| B9 | ingeniero-go / constructor-web | Claves temporales para siempre: no hay cambio obligatorio ni ruta para cambiar la propia clave | `debe_cambiar_clave` + `POST /api/admin/sesion/clave`; forzar cambio al primer login |
| B10 | ingeniero-go | El token de seguimiento no vence ni se revoca | Vencer 30 días después de entregado/cancelado; endpoint para regenerarlo |
| B11 | ingeniero-go | Contenedores sin `mem_limit` ni `pids_limit` | API `256m`/`100`; web `512m`/`200`; probar subida 4000×4000 |
| B12 | constructor-web | `FormProducto.tsx:447`: `_blank` sin `rel` hacia `/producto/{slug}`, que ya no existe | Agregar `rel` y apuntar a una ruta existente o quitar el enlace |
| B3 | constructor-web | `style-src 'unsafe-inline'` | Aceptado en el demo; documentado |
| B7 | Luckas | Correo personal en repo público | Correo noreply desde ahora; sacarlo de los `.md` |

### Informativos
- **I1:** el vendedor ve todos los pedidos, cualquier cliente y `DATOS_TRANSFERENCIA` en mensajes. Es diseño (equipo compartido), no IDOR. Separar por canal requeriría filtrar en el store.
- **I2:** el teléfono en el CSV sale como `'+56…` por la protección contra fórmulas; exportarlo como dígitos sin `+`.
- **I3:** `api/.env.example` y `web/.env.example` no se pudieron leer (regla de denegación); gitleaks sin hallazgos.

## Lo que se revisó y salió bien
- **Autorización por rol en la API:** sin sesión → 401 en las 19 rutas; vendedor → 200 en panel y 403 `sin_permiso` en CSV, resumen, productos (escritura), imagen y usuarios.
- **Escalamiento:** el vendedor no puede darse rol admin (403); campos extra en login → 400; el rol se lee en cada request.
- **Usuarios:** desactivar borra sesiones y bloquea login; regenerar clave invalida la anterior; un admin no puede desactivarse ni quitarse el rol; clave temporal ~96 bits, bcrypt 12, `no-store`, nunca guardada en el front.
- **Sesión:** cookie `HttpOnly; SameSite=Strict` (+ `__Host-`/`Secure` en producción), token 256 bits como HMAC, 12 h / 2 h de inactividad, logout invalida.
- **CSRF/CORS:** Origin obligatorio en escrituras; `form-urlencoded` → 415; preflight ajeno sin CORS.
- **Seguimiento:** token de 16 bytes; JSON solo con número, estado, entrega, agencia, n.º de seguimiento y línea de tiempo; 404 idéntico para cualquier token inválido; tiempos indistinguibles; `noindex` y `no-referrer`.
- **CSV:** solo admin; celdas peligrosas escapadas (probado con `=HYPERLINK(...)`, `-2-3`, `=1+1`).
- **Imágenes:** formato por firma, re-codificación (políglota limpiado), límites de tamaño y píxeles, nombre aleatorio, sin path traversal, solo admin.
- **`DATOS_TRANSFERENCIA`:** no aparece en páginas públicas, API pública, bundles JS, `NEXT_PUBLIC_*` ni `docker history`; solo en mensajes del panel.
- **Validación:** SQL parametrizado, `LIKE` escapado, topes de paginación, cantidades y montos, teléfono E.164.
- **Logs:** teléfono enmascarado, IP /24, sin nombres de clientes, claves ni tokens.
- **Cabeceras:** CSP con nonce y `strict-dynamic`, `frame-ancestors 'none'`, nosniff, Referrer-Policy, X-Frame-Options, COOP, Permissions-Policy; `/admin` con `noindex` y `no-store`; API con CSP `default-src 'none'`.
- **Herramientas:** gitleaks historial 0; árbol solo `api/.env` (ignorado); govulncheck 0; npm audit 0; Trivy API 0, web 0, Dockerfiles 0.
- **Docker:** API distroless no root sin puertos; web `USER node`, sin npm, `read_only`, solo 127.0.0.1; datos en volúmenes.
- **Git:** pre-commit y pre-push con gitleaks (fallan cerrados), pre-push bloquea force push; Dependabot.

## No verificado
- Contenido completo de los `.env.example` (I3).
- HTTPS real (`Secure`, `__Host-`, HSTS).
- Secret scanning y push protection en GitHub (Settings → Code security).
- Tope global de login con más de 30 IPs distintas (revisado en código).
