FROM beaver_server AS builder

FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 移除配置文件复制，改为外部挂载
# COPY user.yaml .

COPY --from=builder /build/user_admin/user .

# 修改启动命令，使用外部配置文件路径
CMD ["./user", "-f", "/etc/user.yaml"] 