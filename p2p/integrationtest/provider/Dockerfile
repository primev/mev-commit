FROM golang:1.23.2 AS builder

WORKDIR /app
COPY ./p2p ./p2p
COPY ./x ./x
WORKDIR ./p2p/integrationtest/provider
RUN CGO_ENABLED=0 GOOS=linux go build -o provider-emulator main.go client.go

FROM alpine:latest

COPY --from=builder /app/p2p/integrationtest/provider/provider-emulator /app/provider-emulator
COPY --from=builder /app/p2p/integrationtest/provider/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

