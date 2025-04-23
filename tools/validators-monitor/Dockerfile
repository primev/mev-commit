FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev

WORKDIR /src
COPY . .

WORKDIR /src/tools
RUN go mod download

WORKDIR /src/tools/validators-monitor
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -o /app/validator-monitor .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/validator-monitor .

RUN chmod +x validator-monitor \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

USER appuser
EXPOSE 8080
ENTRYPOINT ["./validator-monitor"]