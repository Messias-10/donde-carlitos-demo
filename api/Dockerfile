# syntax=docker/dockerfile:1
# API Donde Carlitos — build multi-stage. Sin secretos en la imagen: el .env
# se pasa en tiempo de ejecución (docker compose env_file).

# ---- build ----
FROM golang:1.27.4-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS build
WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux GOFLAGS=-trimpath

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download && go mod verify

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w -buildid=" -o /out/server ./cmd/server && \
    go build -ldflags="-s -w -buildid=" -o /out/setup-admin ./cmd/setup-admin && \
    mkdir -p /out/data /out/imagenes

# ---- runtime ----
# distroless static: sin shell ni gestor de paquetes. Usuario nonroot (65532).
FROM gcr.io/distroless/static-debian13:nonroot@sha256:aafe993665bf6200db38a02c11a9351e782a84b6a45e7f33b6695638c88e3931
WORKDIR /app
COPY --from=build /out/server /app/server
COPY --from=build /out/setup-admin /app/setup-admin
# /data pertenece a nonroot: un volumen con nombre nuevo hereda este dueño.
COPY --from=build --chown=65532:65532 /out/data /data
COPY --from=build --chown=65532:65532 /out/imagenes /imagenes

USER 65532:65532
ENV DB_PATH=/data/carlitos.db IMAGENES_DIR=/imagenes LISTEN_ADDR=:8080
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=4s --start-period=10s --retries=3 CMD ["/app/server", "-healthcheck"]
ENTRYPOINT ["/app/server"]
