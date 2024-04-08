FROM golang:1.22.1 AS builder

WORKDIR /app
COPY ./oracle ./oracle
COPY ./contracts-abi ./contracts-abi
WORKDIR ./oracle
RUN CGO_ENABLED=0 GOOS=linux go build -o mev-commit-oracle ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/oracle/mev-commit-oracle /usr/local/bin/mev-commit-oracle
COPY --from=builder /app/oracle/keystore /keystore
COPY --from=builder /app/oracle/config.yaml /config.yaml
COPY --from=builder /app/oracle/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
