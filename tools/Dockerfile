############################################
# 1) Builder
############################################
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev

# Put *everything* into /src so that your ../bridge/standard code
# is present when go mod download runs.
WORKDIR /src
COPY . .

# Download all modules (external + local via replace)
WORKDIR /src/tools
RUN go mod download

# Build only the validators‑monitor binary
WORKDIR /src/tools/validators-monitor
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -o /app/validator-monitor .

############################################
# 2) Final image
############################################
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
