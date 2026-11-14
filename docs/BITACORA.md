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
