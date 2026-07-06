FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/open_rpc/openrpc .

CMD ["./openrpc", "-f", "/etc/openrpc.yaml"]
