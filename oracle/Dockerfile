FROM golang:1.21.1 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mev-commit-oracle ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/mev-commit-oracle /usr/local/bin/mev-commit-oracle
COPY --from=builder /app/keystore /keystore
COPY --from=builder /app/config.yaml /config.yaml
COPY --from=builder /app/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
