FROM golang:1.22.1 AS builder

WORKDIR /app
COPY . .
WORKDIR ./oracle
RUN CGO_ENABLED=0 GOOS=linux go build -o mev-commit-oracle ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/oracle/mev-commit-oracle /usr/local/bin/mev-commit-oracle
COPY --from=builder /app/oracle/integrationtest/key /key
COPY --from=builder /app/oracle/keystore /keystore
COPY --from=builder /app/oracle/integrationtest/config.yaml /config.yaml
COPY --from=builder /app/oracle/integrationtest/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
