FROM golang:1.23.2 AS builder

WORKDIR /app
COPY ./p2p ./p2p
WORKDIR ./p2p/integrationtest/real-bidder
RUN CGO_ENABLED=0 GOOS=linux go build -o bidder-emulator main.go

FROM alpine:latest

COPY --from=builder /app/p2p/integrationtest/real-bidder/bidder-emulator /app/bidder-emulator
COPY --from=builder /app/p2p/integrationtest/real-bidder/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

