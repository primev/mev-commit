FROM golang:1.22.1 AS builder

WORKDIR /app
COPY . .
WORKDIR ./p2p
RUN cd integrationtest/real-bidder && \
    go mod tidy && \
    go mod vendor && \
    CGO_ENABLED=0 GOOS=linux go build -o bidder-emulator main.go

FROM alpine:latest

COPY --from=builder /app/p2p/integrationtest/real-bidder/bidder-emulator /app/bidder-emulator
COPY --from=builder /app/p2p/integrationtest/real-bidder/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

