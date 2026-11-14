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
| docker / compose | 29.1.3 / v5.0.1 | CLI OK, **daemon apagado** (se pidió abrir Docker Desktop) |
| node / npm | 24.20.0 / 11.19.0 | OK (opcional) |
| go | — | No instalado → se usa imagen `golang` |
| MCP playwright | disponible | OK |

- `git init -b main`, identidad local configurada, hook `.githooks/pre-commit` con gitleaks vía Docker creado y activado (`core.hooksPath`). Falla cerrado si Docker no corre. **Pendiente:** probar con token falso y primer commit (requiere Docker).
- **Pendiente:** crear repo privado `Messias-10/donde-carlitos-demo` (requiere `gh`).
- Fase 1 lanzada en paralelo mientras se resuelve Docker/gh (no los necesita).
