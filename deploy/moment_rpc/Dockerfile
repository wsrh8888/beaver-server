FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/moment_rpc/momentrpc .

CMD ["./momentrpc", "-f", "/etc/momentrpc.yaml"]
