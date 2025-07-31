FROM beaver_server AS builder

FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 移除配置文件复制，改为外部挂载
# COPY dictionary.yaml .

COPY --from=builder /build/dictionary_api/dictionary .

# 修改启动命令，使用外部配置文件路径
CMD ["./dictionary", "-f", "/etc/dictionary.yaml"] 