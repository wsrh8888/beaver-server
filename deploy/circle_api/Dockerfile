FROM beaver_server AS builder

FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/circle_api/circle .

CMD ["./circle", "-f", "/etc/circle.yaml"]
