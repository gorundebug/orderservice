# syntax=docker/dockerfile:1.7

# Stage 1: Build
FROM servicelib-source AS servicelib-source

FROM golang:1.25-bookworm AS builder

ARG SERVICE_DIR=.
ARG TARGETARCH
ARG SERVICEGEN_RUNTIME_STRIP=ON
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}
ARG GOSUMDB=sum.golang.org
ENV GOSUMDB=${GOSUMDB}
COPY --from=servicelib-source / /tmp/servicelib-source
RUN source_dir=/tmp/servicelib-source; \
    if [ -f "$source_dir/context" ]; then \
      mkdir -p /tmp/servicelib-archive; \
      tar -xf "$source_dir/context" -C /tmp/servicelib-archive; \
      source_dir=/tmp/servicelib-archive; \
    fi; \
    if [ ! -f "$source_dir/go.mod" ]; then \
      source_dir=$(find "$source_dir" -mindepth 1 -maxdepth 1 -type d | head -n 1); \
    fi; \
    test -n "$source_dir" && test -f "$source_dir/go.mod"; \
    mkdir -p /servicelib; \
    cp -a "$source_dir/." /servicelib/
COPY . /workspace
WORKDIR /workspace/${SERVICE_DIR}
RUN if [ -f /workspace/go.work ]; then \
      cd /workspace \
      && go work edit \
        -replace github.com/gorundebug/servicelib=/servicelib; \
    else \
      go mod edit \
        -replace github.com/gorundebug/servicelib=/servicelib; \
    fi
RUN --mount=type=cache,id=servicegen-go-mod-v1-${TARGETARCH},target=/go/pkg/mod,sharing=locked \
    go mod download
RUN --mount=type=cache,id=servicegen-go-mod-v1-${TARGETARCH},target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=servicegen-go-build-v1-${TARGETARCH},target=/root/.cache/go-build,sharing=locked \
    if [ "${SERVICEGEN_RUNTIME_STRIP}" = "ON" ]; then \
      GO_LINKER_FLAGS="-s -w"; \
    else \
      GO_LINKER_FLAGS=""; \
    fi \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="${GO_LINKER_FLAGS}" -o /app/service ./cmd/service/main.go

# Stage 2: Runtime
FROM debian:bookworm-slim AS runtime

ARG SERVICE_DIR=.
ARG TARGETARCH
RUN rm -f /etc/apt/apt.conf.d/docker-clean
RUN --mount=type=cache,id=servicegen-go-apt-lists-v1-${TARGETARCH},target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,id=servicegen-go-apt-cache-v1-${TARGETARCH},target=/var/cache/apt,sharing=locked \
    apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates wget

WORKDIR /app

COPY --from=builder /app/service .
COPY ${SERVICE_DIR}/config/ ./config/

EXPOSE 9091
EXPOSE 9201

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO/dev/null http://localhost:9091status || exit 1

ENTRYPOINT ["./service", "-config", "./config/config.yaml", "-values", "./config/overrides.yaml"]