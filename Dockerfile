# Stage 1: Build
FROM golang:1.25-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/service ./cmd/service/main.go

# Stage 2: Runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates wget && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/service .
COPY config/ ./config/

EXPOSE 9091
EXPOSE 9201

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO/dev/null http://localhost:9091status || exit 1

ENTRYPOINT ["./service", "-config", "./config/config.yaml", "-values", "./config/overrides.yaml"]