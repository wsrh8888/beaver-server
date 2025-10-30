FROM beaver_server AS builder

FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 移除配置文件复制，改为外部挂载
# COPY track.yaml .

COPY --from=builder /build/track_api/track .

# 修改启动命令，使用外部配置文件路径
CMD ["./track", "-f", "/etc/track.yaml"]
