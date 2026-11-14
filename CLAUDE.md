# Proyecto: Donde Carlitos Emporio — sitio web + pedidos online

## Qué es este repo
Soy Luckas, desarrollador freelance en Santiago. Este repo es una **propuesta comercial** para el dueño de Donde Carlitos Emporio. Incluye:
- un prototipo funcionando: landing + API en Go + panel de gestión de pedidos. **Desde el 2026-11-14 (decisión de Luckas) no hay carrito ni pedido desde la web:** el negocio vende en el local y por WhatsApp (un número para Santiago/retiro y otro para regiones); el panel ordena esos pedidos (canal, pago, despacho, seguimiento) y la landing manda a cada cliente al WhatsApp correcto;
- un documento de propuesta.

La meta: que el dueño mire el demo en su celular y diga "esta es mi tienda, pero más ordenada".

**El plan de trabajo completo está en `PLAN_MAESTRO.md`.** Si la sesión se corta, retoma desde `docs/BITACORA.md`.

## Datos del negocio (sacados de fuentes públicas en nov 2026; reconfirmar)
- **Nombre:** Donde Carlitos Emporio (en TikTok se escribe "🏆Donde Carlitos Emporio🏆")
- **Rubro:** importaciones y bazar para el hogar, sector barrio Meiggs / Santiago Centro
- **Dirección:** Alameda 2432, Santiago Centro, a 5 min del Metro Unión Latinoamericana y a 7 min del Metro República
- **Horario publicado:** lunes a viernes 09:00–17:00, sábado 09:00–15:00
- **WhatsApp de pedidos:** +56 9 2114 5368 (formato wa.me: 56921145368)
- **Envíos:** a todo Chile (dato de una página de Facebook que Luckas dice que no es del negocio; **reconfirmar** en TikTok/Instagram)
- **Lema en redes:** "El surtido más grande y los precios más económicos del mercado"
- **Confianza:** en su bio avisan por estafas y dicen que nunca le hablarán al cliente por interno

## Redes (confirmadas por Luckas el 2026-11-14, Fase 0)
- **TikTok:** https://www.tiktok.com/@mayte_emporio
- **Instagram:** https://www.instagram.com/mayte_emporioo/
- **Facebook:** no tiene (según Luckas).
- **Nota de Luckas:** aunque el usuario dice "mayte emporio", es la misma empresa Donde Carlitos Emporio.
- Encontradas antes en búsquedas, **no confirmadas como oficiales**: TikTok @_dondecarlitosemporio7 y Facebook donde.carlitos.33. Sirven solo como referencia secundaria; no citarlas como cuentas del negocio sin confirmar con el dueño.

### Por verificar (no afirmar sin confirmar)
- Un video de inauguración menciona un local en **Abate Molina 77, local 5** (Santiago Centro). ¿Sigue activo? ¿Es un segundo local?
- Ese mismo video menciona pago en **efectivo, transferencia y débito**, y "detalle y embalaje". ¿Venden también por mayor?
- **"Mayte Emporio"**: Luckas confirma que es la misma empresa (sus redes oficiales usan ese usuario). **"Andres Emporio"**: relación no confirmada; no afirmar nada.
- Antes se pensaba que la venta era solo presencial; una página de Facebook (no oficial según Luckas) dice "envíos a todo Chile". Confirmar en las redes oficiales.

## Contexto barrio Meiggs (punto de partida; analista-meiggs lo verifica y profundiza)
- Barrio comercial de Santiago Centro, junto a la Alameda y a pasos de Estación Central.
- Conocido por productos importados a bajo precio, venta al detalle y por mayor, y por sus peaks de temporada (sobre todo Navidad).
- Según una crónica de Ex-Ante (2025), concentra la comunidad china más grande de Chile.
- Hipótesis: una parte importante de los compradores son emprendedores que compran para revender, también desde regiones.
- Varios "emporios" del sector se hicieron conocidos vendiendo por TikTok (videos de recorrido, precios en pantalla, lives).

## GitHub
- **Cuenta:** Messias-10 (https://github.com/Messias-10)
- **Correo para commits:** luckasreyes945@gmail.com
- **Repo:** `Messias-10/donde-carlitos-demo`, **público** (decisión de Luckas del 2026-11-14, para portafolio). Por eso las capturas con datos de terceros (`research/raw/tt-*.png`, `secundaria-*.png`) quedan solo en local, fuera de git.

## Stack (decidido)
- **`web/`:** Next.js (App Router) + TypeScript + Tailwind. Mobile-first.
- **`api/`:**
  - Go 1.22 o superior, con `net/http` de la librería estándar.
  - SQLite (`modernc.org/sqlite`, sin CGO) detrás de una interfaz, para pasar a PostgreSQL en producción sin reescribir.
- **Ejecución:** todo corre con **Docker Compose** (Luckas tiene Docker instalado). Si Node o Go no están instalados en la máquina, se usan sus imágenes oficiales de Docker. No se instala software en el sistema.
- **Versiones:** las estables actuales, verificadas al empezar y fijadas (`package-lock.json`, `go.sum`, tags de imagen).
- **Comandos:** se usa `docker compose` directamente. No se depende de `make`, porque en Windows puede no existir.

## Reglas de contenido (todos los agentes)
1. **Nada inventado sobre el negocio.** Todo dato lleva su fuente (URL o captura en `research/raw/`). Si es inferencia, va marcada como *(inferencia)*.
2. **Precios y productos:** si no salen de un post o video real, se marcan como `[EJEMPLO]`.
3. **Prohibido inventar** testimonios, reseñas, cifras de ventas o de seguidores.
4. **Fotos, videos y logo de sus redes:** solo se usan en el prototipo que se le muestra al dueño. Anotar el origen de cada uno en `web/public/CREDITOS.md`. El logo no se rediseña.
5. **Mobile-first:** la mayoría de las visitas llegará desde TikTok o Instagram, en el navegador interno de la app.
6. **Español de Chile**, cercano y directo. Nada corporativo.
7. **Sin estereotipos** sobre el barrio, sus comerciantes ni sus comunidades. Se describe lo que se observa, con fuente.
8. Cada agente escribe solo en su carpeta y no pisa el trabajo de otro.

## Seguridad (no negociable)

### Secretos
- **Nunca** van en el código, los commits, los `.md`, el chat ni los briefs a subagentes: API keys, tokens, contraseñas, hashes reales, cookies ni archivos de sesión.
- Viven solo en `api/.env` (ignorado por git). `api/.env.example` lleva solo los nombres de las variables, con el valor vacío.
- Se generan con comandos que los escriben directo al archivo, sin imprimirlos.
- Prohibido leer, imprimir (`cat`, `type`, `echo`, `Get-Content`) o copiar `.env` y `*.local.txt`.
- Un hook pre-commit con gitleaks (vía Docker) bloquea los commits con secretos. Antes de cada push se escanea el historial completo.
- Si alguna vez se sube un secreto: parar, avisar a Luckas, darlo por comprometido y rotarlo.

### Código
- SQL siempre con parámetros. Validar toda entrada. Límite de tamaño del body. Timeouts en el servidor HTTP.
- **Admin:**
  - contraseña con bcrypt;
  - cookie de sesión httpOnly + SameSite=Strict (+ Secure en producción);
  - verificación del header Origin en las peticiones que cambian datos;
  - límite de intentos de login.
- **Límites y CORS:** límite de pedidos por IP. CORS cerrado: el front llama a la API por el mismo origen, con rewrites de Next.js.
- **Cabeceras:** CSP solo con los dominios necesarios (TikTok, mapas), X-Content-Type-Options, Referrer-Policy y frame-ancestors. HSTS en producción.
- **Errores y logs:** los errores no muestran detalles internos al cliente. Los logs no guardan datos personales completos (el teléfono va enmascarado).
- Nada secreto en variables `NEXT_PUBLIC_*`.
- **Dependencias:** las mínimas. `npm audit` y `govulncheck` sin hallazgos altos ni críticos. Dependabot configurado (`.github/dependabot.yml`).

### Docker
- Imágenes multi-stage, mínimas, con usuario no root.
- Sin secretos dentro de la imagen: el `.env` se pasa en tiempo de ejecución. `.dockerignore` en cada servicio.
- Puertos publicados solo en `127.0.0.1`.
- La base SQLite vive en un volumen de Docker, nunca en el repo.

### Git y GitHub
- Repo público por decisión de Luckas (2026-11-14). Verificar con `gh repo view --json visibility` y no cambiar la visibilidad sin que él lo pida. Nunca subir datos personales de terceros ni capturas con usuarios.
- Commits pequeños, en español, al cerrar cada fase.
- Prohibido el force push, cambiar la visibilidad del repo o borrarlo.
- Datos de clientes reales nunca van al repo. Los pedidos del demo son ficticios.

## Autonomía
Trabaja de punta a punta sin pedir permiso entre fases. **Solo** se interrumpe a Luckas en estos casos:
1. Al inicio, para confirmar los links de TikTok, Instagram y Facebook (una sola pregunta).
2. Si `gh` no tiene sesión iniciada con la cuenta Messias-10: correr `gh auth login --web -h github.com` y pedirle que lo complete en el navegador.
3. Si Docker no está corriendo: pedirle que abra Docker Desktop.
4. Si Instagram exige iniciar sesión en la ventana del navegador: pedírselo una sola vez. Si no se puede, seguir con lo demás y marcar lo que faltó.

**Nunca:** cambiar la visibilidad del repo sin pedido de Luckas, desplegar a internet, borrar el repo, hacer force push ni mostrar secretos.

## Estructura
```
research/            investigador-marca, analista-meiggs, scout-skills (raw/ = capturas)
propuesta/           estratega-propuesta
docs/                modelo de datos, openapi.yaml, BITACORA.md
api/                 ingeniero-go
web/                 constructor-web
docker-compose.yml   ingeniero-go
seguridad/           auditor-seguridad
qa/                  qa-revisor
entregables/         orquestador
.claude/skills/      skills instaladas + estilo-carlitos
```
