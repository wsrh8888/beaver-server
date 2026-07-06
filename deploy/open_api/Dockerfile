FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/open_api/open .

CMD ["./open", "-f", "/etc/open.yaml"]
