# Plan maestro: Donde Carlitos Emporio

Eres el **orquestador**. Coordinas a los subagentes de `.claude/agents/` y haces tú mismo solo la Fase 0 y la Fase 5.

A cada subagente pásale un brief con el contexto que necesita (no sabe lo que tú leíste) y el archivo que debe entregar. Nunca incluyas secretos en un brief.

**Al cerrar cada fase:**
- escribe 5 líneas de avance en el chat y sigue, sin esperar respuesta;
- registra el avance en `docs/BITACORA.md`;
- haz un commit.

Usa una lista de tareas y ve marcándola. Las únicas razones para detenerte están en la sección "Autonomía" de `CLAUDE.md`.

---

## Fase 0: preparación (la haces tú)

### 0.1 Links
Pídele a Luckas, en un solo mensaje, que confirme los links de TikTok, Instagram y Facebook. Muéstrale los que ya están en `CLAUDE.md` para que solo confirme o corrija. Si ya vienen en su mensaje, no preguntes. Actualiza la sección "Redes" de `CLAUDE.md`.

### 0.2 Entorno
Verifica y reporta en una tabla (herramienta | versión | estado):
- `git`
- `gh auth status`: la cuenta debe ser **Messias-10**
- `docker --version`, `docker compose version` y `docker info` (que el daemon esté corriendo)
- `node` / `npm` y `go` (opcionales)
- si el MCP de navegador `playwright` está disponible

**Si falta algo:**
- **Node o Go:** no los instales. Usa las imágenes oficiales `node` y `golang` de Docker para todo.
- **Docker apagado:** pide a Luckas que abra Docker Desktop.
- **`gh` sin sesión, o con otra cuenta:** corre `gh auth login --web -h github.com` y pídele que lo complete.
- **Sin navegador MCP:** sigue igual. Los agentes usarán WebSearch y WebFetch y marcarán lo que no pudieron ver.

**Versiones:** consulta las versiones estables actuales de Next.js (`npm view next version`, local o en contenedor), de Node LTS y de Go (tags oficiales de las imágenes). Anótalas en `docs/BITACORA.md`; los agentes las usarán fijas.

### 0.3 Git
1. `git init -b main`, si no es un repo.
2. Configura la identidad del repo (local, no global):
   - `git config user.name "Messias-10"`
   - `git config user.email "luckasreyes945@gmail.com"`
3. Crea `.githooks/pre-commit`, que corre **gitleaks vía Docker** sobre lo staged y bloquea el commit si encuentra secretos, y actívalo con `git config core.hooksPath .githooks`.
   - Revisa la sintaxis actual con `docker run --rm zricethezav/gitleaks:latest --help`.
   - En Windows/Git Bash, cuida la conversión de rutas al montar volúmenes (`MSYS_NO_PATHCONV=1`).
   - Pruébalo: stagea un archivo temporal con un token falso con forma de `ghp_…`, confirma que el commit se bloquea y borra el archivo.
4. Primer commit: `chore: estructura inicial del proyecto`.

### 0.4 GitHub
1. Si `Messias-10/donde-carlitos-demo` no existe: `gh repo create Messias-10/donde-carlitos-demo --private --source=. --remote=origin --push`.
2. Si existe, úsalo como remoto y confirma que es privado.
3. Verifica con `gh repo view Messias-10/donde-carlitos-demo --json visibility`. Tiene que decir `PRIVATE`; si no, detente y avisa.

---

## Fase 1: investigación (en paralelo)
Lanza a la vez:
- **investigador-marca:** redes de Carlitos. Logo, forma de vender, ofertas, personalidad. Crea la skill `estilo-carlitos`.
- **analista-meiggs:** barrio Meiggs, comprador y competencia.
- **scout-skills:** busca, audita e instala skills de riesgo bajo.

Si Instagram pide iniciar sesión, aplica el punto 4 de "Autonomía": se pide una vez y no se espera indefinidamente.

**Cierre:** commit `docs: investigación de marca, barrio y skills`.

## Fase 2: estrategia
- **estratega-propuesta:** propuesta de valor, arquitectura de la landing, copy, mensaje de WhatsApp, pitch y fases.

**Cierre:** commit `docs: propuesta de valor y copy`.

## Fase 3: construcción
1. **ingeniero-go** escribe primero el contrato: `docs/modelo-datos.md` y `docs/openapi.yaml`.
2. Luego, en paralelo:
   - **ingeniero-go:** implementa la API, el `api/Dockerfile` y el `docker-compose.yml`.
   - **constructor-web:** construye el front contra el contrato (con mocks mientras la API no esté) y el `web/Dockerfile`.
3. **Integración:**
   - `docker compose up -d --build`;
   - `/healthz` responde;
   - prueba de humo del pedido de punta a punta: web → API → BD → panel admin.

**Cierre:** commit `feat: web y API integradas en Docker`.

## Fase 4: seguridad y QA (en paralelo)
- **auditor-seguridad:** escribe `seguridad/reporte-N.md`.
- **qa-revisor:** escribe `qa/reporte-N.md`.

Las correcciones las hace el dueño de cada parte: ingeniero-go para `api/` y Docker, constructor-web para `web/`.

**Se repite hasta que:**
- seguridad no tenga hallazgos críticos ni altos, y
- QA diga "listo para mostrar al dueño".

El máximo es 3 rondas. Si no se logra, documenta qué falta. Haz un commit por ronda.

## Fase 5: entrega (la haces tú)
1. **Entregables en `entregables/`:**
   - `propuesta-donde-carlitos.md`: 1–2 páginas con diagnóstico, propuesta de valor, qué incluye el demo, fases, inversión `[COMPLETAR LUCKAS]` y siguiente paso.
   - `pitch-dueno.md`: mensaje de contacto, guion de 2 minutos y respuestas a objeciones.
   - `como-correr-el-demo.md`: levantar, detener y resetear con Docker; dónde está la clave admin; cómo se desplegaría después (no desplegar).
   - `resumen-seguridad.md`: qué está protegido y cómo, en lenguaje simple, para mostrarle al dueño.
   - `capturas/`: móvil y escritorio, incluido el panel admin.
2. **Clave admin:** corre el comando `setup-admin` de la API. Este comando escribe el hash en `api/.env` y la clave en `api/ADMIN_PASSWORD.local.txt`. No leas ni muestres ninguno de los dos.
3. **Dejar corriendo:** `docker compose up -d --build` y verifica que respondan `http://localhost:3000` y `/admin`.
4. **Cierre de repo:** escaneo completo con gitleaks de todo el historial. Si está limpio, commit final y `git push`. Vuelve a verificar que el repo siga privado.
5. **Resumen final a Luckas** (corto):
   - URLs locales;
   - dónde leer la clave admin (`api/ADMIN_PASSWORD.local.txt`, ábrela tú);
   - link del repo;
   - qué falta confirmar con el dueño;
   - qué le mostrarías primero;
   - pendientes.
