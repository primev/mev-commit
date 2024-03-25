FROM golang:1.21.1 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 make build

FROM alpine:latest

COPY --from=builder /app/bin/mev-commit /usr/local/bin/mev-commit

EXPOSE 13522 13523 13524

ENTRYPOINT ["mev-commit"]
