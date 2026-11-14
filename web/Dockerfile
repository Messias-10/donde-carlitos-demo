# Front Next.js de Donde Carlitos Emporio.
# Multi-stage, salida standalone, usuario no root. Sin secretos en la imagen.

ARG NODE_IMAGE=node:26.13.1-alpine

# ---------- 1. Dependencias (exactas, desde package-lock.json)
FROM ${NODE_IMAGE} AS deps
WORKDIR /app
COPY package.json package-lock.json .npmrc ./
RUN npm ci --no-audit --no-fund

# ---------- 2. Build
FROM ${NODE_IMAGE} AS build
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
# Valores públicos (no son secretos). Los rewrites hacia la API se fijan aquí.
ARG NEXT_PUBLIC_WHATSAPP_PEDIDOS=56921145368
ARG NEXT_PUBLIC_SITE_URL=http://localhost:3000
ARG API_INTERNAL_URL=http://api:8080
ENV NEXT_PUBLIC_WHATSAPP_PEDIDOS=${NEXT_PUBLIC_WHATSAPP_PEDIDOS} \
    NEXT_PUBLIC_SITE_URL=${NEXT_PUBLIC_SITE_URL} \
    API_INTERNAL_URL=${API_INTERNAL_URL}
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# ---------- 3. Ejecución
FROM ${NODE_IMAGE} AS run
WORKDIR /app
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000 \
    HOSTNAME=0.0.0.0
# Seguridad M2: parche de zlib y sin gestores de paquetes en ejecución
# (el npm de la imagen base trae dependencias con CVE altos y no se usa).
RUN apk upgrade --no-cache zlib \
 && rm -rf /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack \
           /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/corepack \
           /usr/local/bin/yarn /usr/local/bin/yarnpkg /opt/yarn-*
# Solo la salida standalone, los estáticos y public/. Dueño root y solo lectura;
# la caché de imágenes (.next/cache) es lo único escribible por el usuario node.
COPY --from=build /app/public ./public
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
RUN mkdir -p .next/cache && chown -R node:node .next/cache
USER node
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:3000/healthz || exit 1
CMD ["node", "server.js"]
