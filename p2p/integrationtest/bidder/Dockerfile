FROM golang:1.21.1 AS builder

WORKDIR /app
COPY . .

RUN cd integrationtest/bidder && \
    CGO_ENABLED=0 GOOS=linux go build -o bidder-emulator main.go

FROM alpine:latest

COPY --from=builder /app/integrationtest/bidder/bidder-emulator /app/bidder-emulator
COPY --from=builder /app/integrationtest/bidder/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

