FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/auth_rpc/authrpc .

CMD ["./authrpc", "-f", "/etc/authrpc.yaml"]
