# Kit multiagente: Donde Carlitos Emporio

Claude Code hace todo solo:
- revisa lo que tienes instalado;
- prepara git, GitHub y Docker;
- investiga la marca y el barrio;
- arma la propuesta;
- construye la web (Next.js) y la API (Go);
- audita la seguridad;
- deja el demo corriendo en Docker.

Al inicio solo te pregunta los links de TikTok, Instagram y Facebook.

## Qué necesitas antes
- **Claude Code** con tu sesión iniciada.
- **Docker Desktop abierto.**
- **Git.**
- **GitHub CLI (`gh`):** si no tiene sesión con Messias-10, Claude Code te pedirá completar `gh auth login` en el navegador.
- **Node.js (recomendado):** lo usa el navegador que mira Instagram y TikTok (`.mcp.json`). Sin Node, los agentes igual trabajan, pero ven menos de las redes.

No necesitas Go ni Node para construir: si faltan, Claude Code usa sus imágenes de Docker.

## Cómo arrancar
1. Descomprime el zip y entra a la carpeta `kit-donde-carlitos/`. En Mac, `.claude` queda oculta en Finder (Cmd + Shift + . para verla).
2. Abre una terminal ahí y ejecuta `claude`.
3. Cuando pregunte si apruebas el servidor MCP `playwright` del proyecto, acepta.
4. Pega `PROMPT_INICIAL.md`, con los links al final.
5. Si Instagram pide iniciar sesión en la ventana del navegador que se abre, entra tú una vez.

## Qué trae
```
CLAUDE.md                 contexto, stack, seguridad y reglas de autonomía
PLAN_MAESTRO.md           las fases que ejecuta el orquestador (0 a 5)
PROMPT_INICIAL.md         lo que pegas en Claude Code
.mcp.json                 navegador Playwright para ver las redes
.claude/settings.json     permisos: autoriza lo necesario y bloquea lo peligroso
.claude/agents/           8 agentes:
  investigador-marca      redes: logo, forma de vender, ofertas, personalidad → skill estilo-carlitos
  analista-meiggs         barrio Meiggs, comprador, competencia
  scout-skills            busca, audita e instala skills seguras
  estratega-propuesta     propuesta de valor, copy, mensaje de WhatsApp, pitch
  ingeniero-go            API en Go + SQLite + docker-compose
  constructor-web         front Next.js conectado a la API + panel admin
  auditor-seguridad       secretos, dependencias, API, front, Docker, GitHub
  qa-revisor              pruebas de punta a punta en el celular
.github/dependabot.yml    alertas de dependencias
.gitignore                deja fuera .env, claves, bases de datos y sesiones
```

## Seguridad, en corto
- **Repo:** privado en `Messias-10/donde-carlitos-demo`.
- **Secretos:** solo en `api/.env`, que está ignorado. Claude Code tiene bloqueado leerlo.
- **Commits:** un hook con gitleaks bloquea los que traigan secretos, y antes de cada push se escanea todo el historial.
- **Bloqueado en `.claude/settings.json`:** force push, hacer público, editar o borrar el repo.
- **Contraseña del admin:** queda en `api/ADMIN_PASSWORD.local.txt`. Ábrela tú; Claude Code no la ve.
- **Permisos:** `.claude/settings.json` deja que Claude Code trabaje sin pedirte permiso a cada rato (Docker, npm, go, git). Si quieres más control, borra líneas de `allow`.
