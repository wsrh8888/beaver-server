FROM beaver_server AS builder

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/document_api/document .

CMD ["./document", "-f", "/etc/document.yaml"]
