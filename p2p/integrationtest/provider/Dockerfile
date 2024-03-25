FROM golang:1.21.1 AS builder

WORKDIR /app
COPY . .

RUN cd integrationtest/provider && \
    CGO_ENABLED=0 GOOS=linux go build -o provider-emulator main.go client.go

FROM alpine:latest

COPY --from=builder /app/integrationtest/provider/provider-emulator /app/provider-emulator
COPY --from=builder /app/integrationtest/provider/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

