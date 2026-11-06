FROM beaver_server AS builder

FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 配置文件外部挂载（见 docker-compose.yaml volumes）

COPY --from=builder /build/doc_sync_ws/docsync .

CMD ["./docsync", "-f", "/etc/docsync.yaml"]
