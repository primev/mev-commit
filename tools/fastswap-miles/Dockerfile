FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev

WORKDIR /src
COPY . .

WORKDIR /src/tools
RUN go mod download

WORKDIR /src/tools/fastswap-miles
RUN CGO_ENABLED=1 GOOS=linux \
    go build -a -installsuffix cgo \
    -o /app/fastswap-miles .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/fastswap-miles .

RUN chmod +x fastswap-miles \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

USER appuser
ENTRYPOINT ["./fastswap-miles"]
