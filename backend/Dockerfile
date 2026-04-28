# syntax=docker/dockerfile:1.6
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary (关闭 CGO，剥离符号)
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /nosocial ./cmd/server

# Runtime stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata wget \
    && addgroup -S nosocial \
    && adduser -S -G nosocial -H -h /app nosocial
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /nosocial /app/nosocial
COPY config/config.yaml /app/config/config.yaml
RUN chown -R nosocial:nosocial /app

USER nosocial
EXPOSE 8080

# 健康检查：请确保后端实现 /api/v1/public/config 等无鉴权健康接口
# （当前有 /api/v1/public/config，返回 200 表示存活）
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/api/v1/public/config >/dev/null 2>&1 || exit 1

CMD ["/app/nosocial"]
