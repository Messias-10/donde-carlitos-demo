# Bitácora — Donde Carlitos Emporio

Registro de avance por fase. Si la sesión se corta, retomar desde la última entrada.

## Versiones fijadas (consultadas el 2026-11-14)
| Componente | Versión | Fuente |
|---|---|---|
| Next.js | 16.4.6 | `npm view next version` (dist-tag `latest`) |
| Node.js (LTS activa "Lithium") | 26.13.1 → imagen `node:26.13.1-alpine` | nodejs.org/dist/index.json |
| Go | 1.27.4 → imagen `golang:1.27.4-alpine` | go.dev/dl/?mode=json |
| gitleaks | v8.28.0 → imagen `zricethezav/gitleaks:v8.28.0` | tag fijo (verificar al levantar Docker) |

Los tags de imagen se verifican con `docker pull` cuando Docker esté corriendo. Si un tag no existe, se usa el más cercano y se anota aquí.

## Fase 0 — Preparación (2026-11-14)
- Links confirmados por Luckas: TikTok @mayte_emporio, Instagram @mayte_emporioo, Facebook: no tiene. "Mayte Emporio" = misma empresa Donde Carlitos. `CLAUDE.md` actualizado.
- Entorno:

| Herramienta | Versión | Estado |
|---|---|---|
| git | 2.56.1 | OK |
| gh | — | **No instalado** (se pidió a Luckas) |
| docker / compose | 29.1.3 / v5.0.1 | OK (Docker Desktop abierto por Luckas) |
| node / npm | 24.20.0 / 11.19.0 | OK (opcional) |
| go | — | No instalado → se usa imagen `golang` |
| MCP playwright | disponible | OK |

- `git init -b main`, identidad local configurada, hook `.githooks/pre-commit` con gitleaks vía Docker activado (`core.hooksPath`). Falla cerrado si Docker no corre. **Probado:** un token falso `ghp_…` bloqueó el commit ("leaks found: 1"); archivo borrado.
- Imágenes verificadas con `docker pull`: gitleaks v8.28.0, golang:1.27.4-alpine, node:26.13.1-alpine.
- Primer commit `chore: estructura inicial del proyecto`. `.gitattributes` con `eol=lf` para que scripts y Dockerfiles no se rompan en Windows.
- **Pendiente:** crear repo privado `Messias-10/donde-carlitos-demo` (requiere `gh`).
- Fase 1 lanzada en paralelo mientras se resuelve Docker/gh (no los necesita).
- Para ganar tiempo, ingeniero-go adelantó el contrato y la API (Fase 3, paso 1 y 2 backend); solo depende de `research/oportunidades.md`, ya listo.

## Fase 1 — Investigación (en curso)
- scout-skills: 5 skills de solo texto instaladas (ver `research/skills-instaladas.md`). Revisadas por el orquestador: sin scripts ni hooks.
- analista-meiggs: `research/contexto-meiggs.md`, `competencia.md`, `oportunidades.md`. Capturas de competencia pendientes (sin navegador).
- investigador-marca: en curso.
- investigador-marca: `adn-marca.md`, `faq-clientes.md`, `dolores-pedido.md`, `ofertas.md`, `top-contenido.md`, `pendientes.md`, `productos-semilla.json` (40 productos con precio real hablado en video) y skill `estilo-carlitos`. TikTok visto a fondo; Instagram bloqueado por login (sin comentarios ni destacadas). Capturas en `research/raw/` (39 MB; las `tt-*.png` muestran usuarios de comentaristas: solo uso interno).
- Contradicciones abiertas en `research/pendientes.md`: envíos a regiones, horario (3 versiones), medios de pago por tienda, mínimo por mayor, ícono de Facebook en el logo.
- Commit `docs: investigación de marca, barrio y skills`.

## Fase 2 — Estrategia (en curso)
- estratega-propuesta lanzado.
- Luckas pidió (2026-11-14) hacer el repo público y que el orquestador instale `gh`. `gh` se instala con winget a pedido suyo. Lo de repo público queda en consulta: el repo contiene logo, capturas y research de un negocio real sin su permiso.
- Luckas decidió repo **público** (2026-11-14). Antes del primer push: capturas con usuarios de comentaristas (`research/raw/tt-*.png`) y de la cuenta no oficial (`secundaria-*.png`) sacadas del historial (commit local reescrito, nunca publicado) y añadidas a `.gitignore`; quedan solo en local. Textos de research ya citan comentarios sin nombre de usuario. `CLAUDE.md` y `.claude/settings.json` actualizados con la decisión.
- `gh` instalado con winget (C:\Program Files\GitHub CLI). Login web en curso por Luckas.
- estratega-propuesta: `propuesta/propuesta-valor.md`, `arquitectura-landing.md`, `copy-landing.md`, `mensaje-pedido-whatsapp.md`, `pitch.md`, `fases.md`. Titular: "¿LO VISTE EN TIKTOK? ACÁ LO PIDES." Descartada la hipótesis "la competencia no tiene web" (Sale Todo sí tiene). Commit `docs: propuesta de valor y copy`.

## Fase 3 — Construcción (en curso)
- ingeniero-go: API implementada en `api/`; se le pidió priorizar `docs/openapi.yaml`, ajustar el mensaje de WhatsApp a la plantilla y hacer configurable el número de destino.
- constructor-web lanzado (front contra el contrato; repo será público, nada secreto en `NEXT_PUBLIC_*`).
- Login `gh` completado por Luckas (Messias-10). Escaneo gitleaks del historial completo: 5 commits, sin hallazgos. Repo creado: https://github.com/Messias-10/donde-carlitos-demo — `visibility: PUBLIC` (decisión de Luckas). Push de `main` OK.
- ingeniero-go (2026-11-14): contrato `docs/openapi.yaml` + `docs/modelo-datos.md` y API completa en `api/` (Go 1.27.4, `net/http`, SQLite modernc v1.62.0, x/crypto v0.59.0). `gofmt`/`go vet` limpios, `go test ./...` en verde (numeración concurrente, total recalculado, transiciones, mensaje WhatsApp = ejemplos A/B/C de la plantilla byte a byte, 401 admin, bloqueo de login, inyección SQL en búsqueda, setup-admin), `govulncheck` v1.9.0: sin vulnerabilidades. Imagen `distroless/static-debian13:nonroot` (27,9 MB), fijada por digest junto con `golang:1.27.4-alpine`. `docker-compose.yml` en la raíz (api en 127.0.0.1:8080, web en 127.0.0.1:3000, volumen `datos-api`, `env_file` opcional). Semilla: 37 productos de `research/productos-semilla.json` con su video de origen (3 omitidos sin precio) + 6 pedidos DEMO. setup-admin probado solo en un directorio temporal fuera del repo (borrado); el hash va en base64 (`ADMIN_PASSWORD_HASH_B64`) para que compose no interpole los `$`. **No se ejecutó contra `api/.env`.** Decisiones en `api/README.md`.
- constructor-web: front completo en `web/` (landing de 14 secciones, catálogo, ficha, carrito, pedido, admin, healthz, OG, JSON-LD). Lint, tsc y build en verde; `npm audit` sin hallazgos. CSP con nonce en `proxy.ts`. 26 imágenes de producto recortadas del cuadro del video; el orquestador las revisó en una hoja de contactos: sin comentarios, usuarios ni interfaz de TikTok. Capturas del demo en `qa/capturas/`.
- Puerto 3000 ocupado en la máquina de Luckas por otro proyecto (`mouchett-app`): no se tocó. El compose ahora acepta `WEB_PORT` (por defecto 3000) y ajusta `ALLOWED_ORIGINS`. Healthcheck de web pasa a `/healthz`.
- **Integración (2026-11-14):** `docker compose down -v` + `WEB_PORT=3100 docker compose up -d --build` → api y web `healthy`. Páginas `/`, `/catalogo`, `/carrito`, `/pedido`, `/admin` en 200. Pedido de humo DC-0008 (envío, 3 productos) por el mismo origen: total $20.500 = suma de líneas, mensaje WhatsApp igual a la plantilla y `wa.me` decodificado idéntico. Sin header Origin → 403; admin sin sesión → 401; logs con teléfono enmascarado. Cabeceras: CSP, nosniff, Referrer-Policy, X-Frame-Options DENY.
- `setup-admin` ejecutado: `api/.env` y `api/ADMIN_PASSWORD.local.txt` creados (no leídos, ignorados por git). Login incorrecto → 401 `credenciales_invalidas`.

## Fase 4 — Seguridad y QA, ronda 1 (2026-11-14)
- auditor-seguridad: 0 críticos, 2 altos (A1 IP falsificable vía X-Forwarded-For; A2 un pedido anónimo agota el stock), 3 medios, 7 bajos. Reporte guardado por el orquestador en `seguridad/reporte-1.md` (el subagente no pudo escribirlo); salidas redactadas en `seguridad/salidas/`.
- Correcciones asignadas: ingeniero-go → A1, A2, M1, B1, B5 y confirmar `.env.example` vacío. constructor-web → M2, B2, B4 (junto con los hallazgos de QA).
- Orquestador (B6): `@playwright/mcp` fijado en 0.0.85; `deny` ampliado a `--force-with-lease` y `+rama`; nuevo hook `.githooks/pre-push` que bloquea pushes que no avanzan en línea recta y escanea todo el historial con gitleaks (falla cerrado).
- M3 (repo público con fotos de personas del negocio y la propuesta): decisión de Luckas del 2026-11-14 de mantenerlo público; riesgo informado. Desde ahora no se agregan más frames con personas sin necesidad.
- qa-revisor: veredicto **"listo para mostrar al dueño"**, sin bloqueantes; 8 importantes (I-1 stock inventado/urgencia falsa, I-2 mínimos [EJEMPLO] sin marcar, I-3 precio editado sigue diciendo "del video", I-4 no se pueden agregar productos, I-5 hero no dice qué venden, I-6 confirmación se pierde al recargar, I-7 número real del negocio, I-8 inicio pesado: Lighthouse 72) y 16 menores. Reporte guardado por el orquestador en `qa/reporte-1.md`.
- I-7 se mantiene (número del negocio por defecto); se documenta para la reunión. Cambio de enfoque "catálogo oficial" pendiente de decisión de Luckas; no se aplica en esta ronda.
- Correcciones ronda 1 en curso: ingeniero-go (A1, A2, M1, B1, B5, I-1, I-2, I-3, I-4 subida de imagen, M-2, M-10, M-15) y constructor-web (M2, B2, B4, I-2..I-6, I-8 y menores).
- ingeniero-go, correcciones ronda 1 (seguridad + QA), contrato v1.1:
  - **Seguridad.**
    - A1/M1: la IP se toma del último valor de XFF y solo desde un proxy confiable. Comprobado en vivo que Next.js reenvía el XFF del cliente sin sobrescribirlo, así que en el demo `TRUSTED_PROXIES` va vacío y se limita por la IP de la conexión. Web tiene IP fija 172.30.57.10 para producción, detrás de un proxy de borde que sobrescriba XFF.
    - Limitadores LRU en `internal/limite`. Login con tope global de 30/min y semáforo de 2 bcrypt.
    - Pedidos (implementado antes del cambio de dirección; probablemente se elimina): stock al confirmar (migración 0002), tope de 100 por línea y de 3 `nuevo` por teléfono, límites por teléfono y global, y validación antes de consumir cupo.
    - B1: web `read_only` con tmpfs `/tmp` y `/app/.next/cache`; healthy en 3100.
    - B5: la API ya no se publica (override `docker-compose.debug.yml`). `.env.example` reescrito solo con nombres vacíos.
  - **QA.**
    - I-1: semilla sin "quedan pocos"; solo P036 agotado.
    - I-2: `mayor_minimo_confirmado`.
    - I-3: `precio_actualizado_en`.
    - I-4: `POST /api/admin/productos/{id}/imagen` + `GET /api/imagenes/{archivo}`, re-codificadas y en el volumen `imagenes-api`. Dependencia nueva `golang.org/x/image` v0.48.0, para webp.
    - M-2: `NEXT_PUBLIC_SITE_URL`.
    - M-10: UTC.
    - M-15: `oferta_hasta`; P012 hasta 2026-11-15, según su video.
    - La migración 0003 ajusta la semilla ya cargada sin `down -v`.
  - **Calidad.** `go vet`, `go test ./...` y `govulncheck` en verde. Stack recreado con `WEB_PORT=3100` (`down` sin `-v`).

## Cambio de enfoque (2026-11-14, decisión de Luckas)
- **Fuera:** carrito, checkout y pedido desde la web. Motivo: el negocio vende en el local y por WhatsApp, con un número para Santiago/retiro y otro para regiones ("número para región" y "la niña de los pedidos de regiones" en videos/comentarios de nov 2026, ver `research/dolores-pedido.md`).
- **Dentro:** landing pública (dos botones de WhatsApp Santiago/Regiones con mensaje pre-escrito, "cómo comprar desde regiones", antiestafa, vitrina liviana "Lo nuevo de los videos" con botón por producto, página pública de seguimiento sin datos personales) + panel de gestión de pedidos (registro rápido desde el chat, estados cotizado → esperando pago → pagado → preparando → despachado/listo retiro → entregado, pantalla "Hoy", mensajes listos al cliente, varios usuarios, búsqueda por teléfono, resumen y CSV).
- Fase 2 fuera del demo: API oficial de WhatsApp Business.
- Número de regiones desconocido: [CONFIRMAR]; datos de transferencia nunca se inventan (variable `DATOS_TRANSFERENCIA`, vacía por defecto).
- En curso: estratega-propuesta reescribe propuesta/copy/pitch + `propuesta/mensajes-whatsapp.md`; ingeniero-go rediseña contrato (openapi v2) y API. constructor-web entra cuando el contrato esté listo.
- estratega-propuesta: propuesta, arquitectura, copy, pitch y fases reescritos al nuevo enfoque; nuevo `propuesta/mensajes-whatsapp.md` (plantillas a, b, c, d1–d8 con ejemplos); `mensaje-pedido-whatsapp.md` marcado obsoleto. Titular: "¿LO VISTE EN TIKTOK? PÍDELO POR WHATSAPP." Ojo: según `research/adn-marca.md` §G el mismo +56 9 2114 5368 se presentó como "número para región" → puede ser un solo número atendido por dos personas (pregunta para la reunión). Datos [VERIFICAR] en `fases.md`: precios de la API de WhatsApp Business (cobro por mensaje de plantilla desde jul-2025) y vigencia de la Ley 21.719 de datos personales (1-dic-2026, con proyecto de postergación).
- ingeniero-go, **contrato v2 y API del panel** (2026-11-14/15):
  - **Contrato y migración.** `docs/openapi.yaml` 2.0.0 (25 rutas) y `docs/modelo-datos.md` reescritos. Migración 0004: pedidos/items/historial recreados, `admins` → `usuarios` (rol admin|vendedor, activo), stock de los pedidos v1 devuelto. Probada sobre una base v1.1 y aplicada a la base real sin `down -v`.
  - **Fuera.** `POST /api/pedidos`, sus límites por IP y teléfono y el mensaje de pedido web.
  - **Dentro.**
    - Panel: registrar, editar antes de pagar, estados con requisitos (pagado con `confirmo_pago_recibido`; despachado con agencia, n.º y comuna; en_tienda prepara sin pagar), Hoy, búsqueda, clientes, resumen y CSV (solo admin, anti CSV injection).
    - Usuarios con clave generada que se muestra una vez.
    - Mensajes d1–d8 exactos de `propuesta/mensajes-whatsapp.md`, con tests byte a byte de todos sus ejemplos (`wa.me` sin número en DEMO).
    - Seguimiento público con token base64url de 16 bytes, sin datos personales.
    - `/api/config-publica`.
  - **Semilla.** 20 pedidos DEMO: DC-0101 a DC-0109 del guion de `pitch.md` más DC-0090 a DC-0100 cerrados. Usuarios demo "Regiones 1" y "Santiago 1" sin clave utilizable.
  - **Variables nuevas.** `WHATSAPP_SANTIAGO`, `WHATSAPP_REGIONES`, `URL_SITIO`, `DIRECCION_RETIRO`, `HORARIO_RETIRO`, `DATOS_TRANSFERENCIA` (solo en `api/.env`, vacío por defecto).
  - **Calidad.** gofmt, vet, `go test ./...` y govulncheck en verde. API reconstruida y sana; la web no se reconstruyó (la está cambiando constructor-web).
