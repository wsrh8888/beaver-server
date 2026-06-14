FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/platform_api/platform .

CMD ["./platform", "-f", "/etc/platform.yaml"]
