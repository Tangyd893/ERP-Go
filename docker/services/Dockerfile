# ── 通用 Go 微服务多阶段构建 ───────────────────────────
# 用法:
#   docker build -f docker/services/Dockerfile \
#     --build-arg SERVICE_PATH=backend/services/iam-service \
#     --build-arg SERVICE_NAME=iam-service \
#     --build-arg SERVICE_PORT=8081 \
#     -t erp-go/iam-service:latest .

FROM golang:1.25-alpine AS builder

ARG SERVICE_PATH
ARG SERVICE_NAME

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend/ ./backend/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/${SERVICE_NAME} \
    ./${SERVICE_PATH}/cmd/server/

FROM alpine:3.21

ARG SERVICE_NAME
ARG SERVICE_PORT

RUN apk add --no-cache ca-certificates tzdata curl && \
    adduser -D -u 10001 appuser

COPY --from=builder /app/bin/${SERVICE_NAME} /app/service

USER appuser
EXPOSE ${SERVICE_PORT}

HEALTHCHECK --interval=15s --timeout=3s --retries=3 \
    CMD curl -sf http://localhost:${SERVICE_PORT}/health || exit 1

ENTRYPOINT ["/app/service"]
