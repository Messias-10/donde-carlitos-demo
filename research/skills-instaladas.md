# Skills instaladas (scout-skills, Fase 1)

Fecha: 2026-11-14. Alcance: solo proyecto (`.claude/skills/`). Todas son archivos de texto (.md, .json, LICENSE). Ninguna tiene scripts, hooks, `allowed-tools` ni comandos de red. Descargadas con curl desde raw.githubusercontent.com fijadas a un commit, revisadas en una carpeta temporal y copiadas después.
Quedan disponibles al reiniciar la sesión. Antes, los agentes pueden leer el `SKILL.md` directamente.

## Instaladas (5 de 8 posibles), todas riesgo bajo

| Skill | Origen (commit) | Archivos | Para qué en este proyecto | Licencia |
|---|---|---|---|---|
| `frontend-design` | anthropics/skills, `skills/frontend-design` (2c5958e) | SKILL.md, LICENSE.txt | Dirección estética propia, evitar el look de plantilla IA, plan de tokens (color, tipografía, layout), movimiento con moderación. Es la original oficial. | Ver LICENSE.txt (el repo no declara licencia en la API) |
| `vercel-react-best-practices` | vercel-labs/agent-skills, `skills/react-best-practices` (1a2c79d) | SKILL.md, AGENTS.md, README.md, metadata.json, rules/ (70 reglas .md) | Rendimiento React/Next.js: cascadas, tamaño del bundle, server components. Clave para 4G y el navegador interno de TikTok/Instagram. | MIT (declarada en el frontmatter) |
| `owasp-security` | agamm/claude-code-owasp, `.claude/skills/owasp-security` (1ae29cb) | SKILL.md, LICENSE, reference/ (4 .md) | Revisión OWASP Top 10:2025 y ASVS, cabeceras, Dockerfile y CI. Para `auditor-seguridad`. Su flujo exige confirmar que el hallazgo es alcanzable antes de reportarlo, lo que reduce falsos positivos. | LICENSE en el repo (1 KB, tipo permisivo) |
| `golang` | saisudhir14/claude-skills, raíz del repo (edef07b) | SKILL.md, LICENSE, references/ (6 .md) | Buenas prácticas Go: errores, concurrencia, pruebas, rendimiento. Para `ingeniero-go`. | MIT |
| `go-security` | saisudhir14/claude-skills, `skills/go-security` (edef07b) | SKILL.md, LICENSE | Validación de entrada, SQL injection, path traversal, bcrypt, cabeceras HTTP, govulncheck. | MIT |

Desinstalar cualquiera: borrar su carpeta, por ejemplo `rm -r ".claude/skills/<nombre>"`. No hay nada instalado fuera de esas carpetas.

Reputación (consultada vía API de GitHub el 2026-11-14):
- anthropics/skills: 194k estrellas, activo (push 2026-11-06).
- vercel-labs/agent-skills: 34k estrellas, activo (push 2026-11-12).
- agamm/claude-code-owasp: 442 estrellas, 0 issues abiertos, activo (push 2026-11-03).
- saisudhir14/claude-skills: 8 estrellas, último push 2026-06-26. La reputación es baja. Lo compensa que el contenido es solo texto, lo revisé completo y es guía estándar (Google/Uber Go Style). Se puede desinstalar sin costo si molesta.

Notas de uso y conflictos:
- `golang` y `go-security` se solapan un poco (bcrypt, cabeceras). No se contradicen.
- `frontend-design` dice "confirmar con el cliente" si falta el brief. En este proyecto el brief lo da `CLAUDE.md` y el trabajo de `investigador-marca`; los agentes no deben detenerse a preguntar (regla de Autonomía).
- `owasp-security` incluye secciones de Kubernetes, Terraform, LLM, etc. que no aplican. Es solo texto cargado bajo demanda.
- `vercel-react-best-practices` asume a veces despliegue en Vercel. El proyecto NO se despliega a internet, se ignoran esas partes.
- `AGENTS.md` de esa skill pesa 108 KB (compilado de las reglas). Que los agentes lean SKILL.md y solo la regla que necesiten en `rules/`.

## Descartadas, con razón

| Skill | Razón |
|---|---|
| `web-design-guidelines` (vercel-labs) | En cada uso descarga las reglas desde una URL remota mutable (`.../web-interface-guidelines/main/command.md`) y las obedece como instrucciones. Instrucciones remotas sin versión: riesgo medio. Rechazada. |
| `skill-creator` (anthropics) | Tiene scripts Python (`eval-viewer/generate_review.py` y otros) y subagentes de evaluación que lanzan procesos. No cumple "sin scripts". Para `estilo-carlitos` basta escribir un SKILL.md simple. |
| `webapp-testing` (anthropics, también en awesome-copilot) | Trae `scripts/with_server.py`, que lanza procesos. Riesgo medio. Para pruebas de navegador se usa Playwright directo desde Docker. |
| `vercel-optimize`, `deploy-to-vercel`, `vercel-cli-with-tokens` | Piden tokens y despliegue a internet. Prohibido por `CLAUDE.md`. `vercel-optimize` además tiene muchos .mjs. |
| `react-view-transitions`, `composition-patterns` | Solo texto, pero son ruido para un prototipo y las transiciones de vista no son fiables en navegadores internos. Se puede reconsiderar. |
| `go-project-layout` (saisudhir14) | Propone Makefile. `CLAUDE.md` prohíbe depender de `make`. |
| `go-testing`, `go-code-review`, `go-concurrency`, etc. (saisudhir14) | Ya cubiertas por `golang` y sus `references/`. Evita duplicar contexto. |
| `multi-stage-dockerfile` (github/awesome-copilot) | Genérica y delgada. `owasp-security/reference/config-and-supply-chain.md` ya cubre Docker seguro. |
| `security-review` (github/awesome-copilot) | Duplica `owasp-security`. Se prefiere una sola. |
| Copias de `frontend-design` en agregadores | Duplicados. Se usó la original de anthropics. |
| Skills SEO Next.js (React SEO Skills, seo-technical, seo-metadata-generator, etc.) | Autores pequeños y sin verificar, vistos solo en listados de terceros y no auditados línea a línea. Para Open Graph y schema.org LocalBusiness basta la Metadata API de Next.js. |

## Huecos sin cubrir
Tailwind/design tokens (lo cubre en parte `frontend-design`), accesibilidad (idem, mínimos de calidad), SEO local (Open Graph y LocalBusiness) y pruebas Playwright. No se encontró candidata de riesgo bajo verificada. Se recomienda que `constructor-web` use la documentación oficial de Next.js, y que `qa-revisor` ejecute Playwright en Docker.

## Límites de esta auditoría
`gh` no estaba instalado, así que la búsqueda fue con WebSearch, WebFetch y la API pública de GitHub. No se revisaron issues abiertos uno a uno. No hice commits.
