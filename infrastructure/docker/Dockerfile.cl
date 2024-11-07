FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY x/go.mod x/go.sum ./x/
COPY cl/go.mod cl/go.sum /app/cl/
COPY contracts-abi/go.mod contracts-abi/go.sum /app/contracts-abi/

RUN cd /app/x && go mod download
RUN cd /app/cl && go mod download
RUN cd /app/contracts-abi && go mod download

COPY . .

WORKDIR /app/cl

RUN go build -o /app/mev-commit-cl ./cmd/redisapp

FROM alpine:3.10

COPY --from=builder /app/mev-commit-cl /usr/local/bin/mev-commit-cl

ENTRYPOINT ["mev-commit-cl", "start"]
