# Auditoría de seguridad — Fase 4, ronda 1

- **Fecha:** 2026-11-14
- **Autor:** auditor-seguridad (guardado por el orquestador; el subagente no pudo escribir el archivo).
- **Alcance:** repo completo (árbol e historial git), `api/` (Go), `web/` (Next.js), imágenes Docker `donde-carlitos-api:local` y `donde-carlitos-web:local`, `docker-compose.yml`, configuración de git y GitHub.
- **Entorno de prueba:** solo localhost. Stack levantado (`WEB_PORT=3100`) y una API temporal aparte (`auditor-api-tmp`, 127.0.0.1:18080, clave de prueba propia fuera del repo, base en tmpfs), ya borrada. En el stack principal solo hubo logins fallidos con un usuario inexistente. No se leyó `api/.env` ni `*.local.txt`.
- **Salidas de herramientas:** `seguridad/salidas/` (secretos redactados).

## Veredicto

**Hay altos pendientes. Sin críticos.** Los 2 altos son de la API y se corrigen en `api/` y el compose.

| Severidad | Cantidad |
|---|---|
| Crítica | 0 |
| Alta | 2 |
| Media | 3 |
| Baja | 7 |

## Hallazgos

### A1 — Alta — IP del cliente falsificable (dueño: ingeniero-go)
- **Ubicación:** `api/internal/http/middleware.go:198-226` (`ipCliente`); `docker-compose.yml:31` (`TRUSTED_PROXIES: 172.16.0.0/12,10.0.0.0/8,192.168.0.0/16`).
- **Evidencia:** el rewrite de Next (`/api/*` → API) no agrega la IP real a `X-Forwarded-For`: solo la pone si el cliente no mandó la cabecera. La API confía en todas las redes privadas y toma el primer valor no confiable desde la derecha, que es el que escribió el cliente.
  - Por la web :3100: 7 logins fallidos con XFF fijo → `401×5, 429, 429`; con XFF rotado → `401×7`. El log registra la IP inventada.
  - Pedidos (API temporal): 12 con la misma IP → `201×10, 429, 429`; con XFF rotado → `201×12`.
  - DoS: 20.050 `POST /api/pedidos` con cuerpo `{}` y XFF únicos llenan el limitador (tope 20.000, falla cerrado en `limites.go:41`); después un cliente legítimo recibe `429`.
  - Ver `salidas/pruebas-dinamicas.txt` §3, §7, §9, §11.
- **Impacto:** el límite por IP y el bloqueo de login no sirven; se puede bloquear el checkout de todos por 10 minutos, repetible. Logs con IPs falsas. En producción, explotable desde internet.
- **Corrección:**
  1. Aceptar XFF solo desde la IP exacta del contenedor web (fijada con `ipv4_address`, `TRUSTED_PROXIES` con ese /32) y tomar solo el último valor.
  2. En producción, proxy de borde que **sobrescriba** XFF y confiar solo en ese salto. En local, sin proxy confiable el límite puede ser global.
  3. Validar el body antes de consumir cupo (`publico.go:164`).
  4. Tope global de pedidos por hora y límite por teléfono.
  5. Al llegar a `maxEntradas`, desalojar las más antiguas en vez de rechazar a todos.

### A2 — Alta — Un pedido anónimo deja la tienda agotada (dueño: ingeniero-go)
- **Ubicación:** `api/internal/store/sqlite/pedidos.go:84-95`; `api/internal/pedidos/pedidos.go:17-18` (`MaxItems=50`, `MaxCantidadPorItem=500`).
- **Evidencia:** un pedido `nuevo` descuenta stock al crearse y solo lo devuelve si el admin cancela. Un solo `POST /api/pedidos` de 1.262 bytes con los 36 productos activos a su stock completo → `201`; agotados: antes 2, después 38 de 38. Ver `salidas/pruebas-dinamicas.txt` §8 y §10.
- **Impacto:** con una petición, sin saltarse límites, la tienda queda agotada hasta que el dueño cancele a mano. Con A1 se repite sin fin y el panel se llena de pedidos falsos.
- **Corrección:** descontar stock solo al pasar a `confirmado` (en `nuevo` solo validar); o reserva con vencimiento automático. Además, topar cantidad por línea del pedido web y limitar pedidos `nuevo` abiertos por teléfono.

### M1 — Media — Login sin tope global, bcrypt costoso (dueño: ingeniero-go)
- **Ubicación:** `api/internal/http/admin.go:99`; `api/internal/auth/auth.go:129,153`.
- **Evidencia:** cada login corre bcrypt costo 12 (~250 ms CPU) exista o no el usuario; el único freno es el bloqueo por IP, que A1 anula. `BloqueoLogin` crece sin tope con IPs rotadas.
- **Impacto:** pocas peticiones por segundo copan la CPU de la API (la tienda cae con ella).
- **Corrección:** semáforo global a bcrypt (p. ej. 2 concurrentes; si está lleno, 429 con `Retry-After`), tope global de intentos por minuto y LRU en `BloqueoLogin`.

### M2 — Media — npm de la imagen base con vulnerabilidades (dueño: constructor-web)
- **Ubicación:** `web/Dockerfile:28-41` (etapa `run` sobre `node:26.13.1-alpine`).
- **Evidencia:** Trivy sobre `donde-carlitos-web:local`: 19 vulnerabilidades en el npm de la imagen base (6 HIGH: undici y brace-expansion; 12 MEDIUM; 1 LOW) y 1 MEDIUM en `zlib`. Los `node_modules` de la app tienen 0. Ver `salidas/trivy-donde-carlitos-web.txt`.
- **Impacto:** no alcanzables desde la app, pero la imagen no pasa un escaneo sin altos y npm sobra en ejecución.
- **Corrección:** en la etapa `run`, `apk upgrade --no-cache zlib` y borrar npm, npx, corepack y yarn; o usar distroless nodejs. Re-escanear con Trivy: 0 HIGH.

### M3 — Media — Repo público con fotos de personas y la propuesta comercial (dueño: decisión de Luckas)
- **Ubicación:** repo público; `web/public/productos/tt-*.webp` (26), capturas de perfil en `research/raw/`, `propuesta/*.md`.
- **Evidencia:** visibilidad `PUBLIC`. Los frames de producto y capturas de perfil muestran personas identificables (presentadores del negocio; transeúntes en `ig-perfil-02.png`). No se vieron menores reconocibles. La regla 4 de `CLAUDE.md` permite esas fotos solo en el prototipo para el dueño.
- **Impacto:** redistribución pública de imagen de personas y material del negocio sin su consentimiento; la propuesta queda pública antes del pitch. Al estar en el historial, quitarlos del árbol no los borra.
- **Corrección (decisión de Luckas):** volver a privado, o aceptar el riesgo por escrito y, desde ahora, no agregar más frames con personas; pedir el visto bueno del dueño antes de difundir.

### Bajos
| # | Dueño | Hallazgo | Corrección |
|---|---|---|---|
| B1 | ingeniero-go | Contenedor `web` sin `read_only` | `read_only: true` + tmpfs para `/tmp` y caché de `.next` |
| B2 | constructor-web | La web no manda HSTS en producción | HSTS cuando `NODE_ENV=production` y el sitio es https, o documentarlo para el proxy de borde |
| B3 | constructor-web | CSP con `style-src 'unsafe-inline'` | Aceptable en el demo; nonce en estilos si se puede, o documentar la excepción |
| B4 | constructor-web | Iframe de TikTok sin `sandbox` (`web/components/VideoTikTok.tsx:31`) | `sandbox="allow-scripts allow-same-origin allow-presentation allow-popups allow-popups-to-escape-sandbox"` y probar |
| B5 | ingeniero-go | API publicada en `127.0.0.1:8080`; desde el host el gateway cae en `TRUSTED_PROXIES` | Quitar el puerto o dejarlo en un perfil `debug` |
| B6 | orquestador | `.mcp.json` usa `@playwright/mcp@latest`; la lista de denegación no cubre `git push origin +main` ni `--force-with-lease` | Fijar versión; ampliar `deny`; considerar protección de rama `main` en GitHub |
| B7 | Luckas | Correo personal en archivos y commits de un repo público | Usar el correo noreply de GitHub desde ahora y sacarlo de los `.md` |

## Lo que salió bien (resumen)
- **Secretos:** gitleaks sobre el historial completo, 0 hallazgos; en el árbol solo `api/.env`, ignorado. `.gitignore` y `.dockerignore` cubren `.env*`, `*.local.txt`, bases y sesiones. `research/raw/tt-*.png` y `secundaria-*.png` nunca estuvieron en git. `docker history` sin secretos; el bundle no contiene secretos ni la URL interna. `NEXT_PUBLIC_*` solo con el número público y la URL del sitio. `setup-admin` no imprime secretos, escribe 0600 de forma atómica.
- **Autenticación:** todas las rutas admin exigen sesión; bcrypt costo 12 con hash señuelo; cookie `HttpOnly; SameSite=Strict` (+ `__Host-` y `Secure` en producción); token de 256 bits guardado como HMAC; vencimiento 12 h / 2 h de inactividad; logout invalida en servidor.
- **CSRF/CORS:** Origin obligatorio y en lista blanca en escrituras; `Content-Type` no JSON → 415; sin cabeceras CORS.
- **Inyección y validación:** SQL parametrizado; `ORDER BY` desde mapa cerrado; `LIKE` escapado; JSON estricto; límites de body; texto sin caracteres de control ni bidi; URLs del admin validadas.
- **Negocio:** total recalculado en servidor; stock con descuento condicional en transacción; la API pública no expone pedidos.
- **Disponibilidad:** timeouts del servidor; paginación con tope.
- **Front:** CSP con nonce y `strict-dynamic` efectiva; sin XSS; `rel="noopener noreferrer"` en enlaces externos; URL de WhatsApp validada.
- **Cabeceras:** nosniff, Referrer-Policy, X-Frame-Options, COOP, Permissions-Policy, CSP; `/admin` con `noindex` y `no-store`; API con CSP `default-src 'none'` y HSTS en producción.
- **Datos personales:** logs con teléfono e IP enmascarados; semilla con clientes ficticios; research sin usuarios de comentaristas.
- **Docker:** API distroless no root, `read_only`, `cap_drop: ALL`, `no-new-privileges`, fijada por digest; web multi-stage con `USER node`; puertos solo en 127.0.0.1; base en volumen. Trivy config 0; Trivy imagen API 0.
- **Dependencias:** govulncheck sin vulnerabilidades; `npm audit` 0; dependencias mínimas con versiones exactas.
- **Git:** hook pre-commit con gitleaks que falla cerrado; Dependabot para npm, gomod y Docker.

## No verificado en esta ronda
- Contenido de `api/.env.example` y `web/.env.example` (bloqueado por la regla de permisos; gitleaks sin hallazgos). Se pidió confirmación a ingeniero-go.
- Secret scanning y push protection en GitHub (`gh api` denegado). Luckas puede revisarlo en Settings → Code security.
- Comportamiento con HTTPS real (`Secure`, `__Host-`, HSTS): revisado solo en código.
