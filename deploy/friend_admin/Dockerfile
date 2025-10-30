FROM beaver_server AS builder

FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 移除配置文件复制，改为外部挂载
# COPY friend.yaml .

COPY --from=builder /build/friend_admin/friend .

# 修改启动命令，使用外部配置文件路径
CMD ["./friend", "-f", "/etc/friend.yaml"] 