FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/platform_rpc/platformrpc .

CMD ["./platformrpc", "-f", "/etc/platformrpc.yaml"]
