# tools/Dockerfile
############################################
# 1) Builder
############################################
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

# set our module root
WORKDIR /src

# grab module metadata early for caching
COPY go.mod go.sum ./
RUN go mod download

# copy all code
COPY . .

# build only the validators‑monitor subpackage,
# emitting the binary into /app/validator-monitor
WORKDIR /src/validators-monitor
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -o /app/validator-monitor .

############################################
# 2) Final image
############################################
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# copy in our statically‑linked Go binary
COPY --from=builder /app/validator-monitor .

RUN chmod +x validator-monitor \
    && addgroup -S appgroup \
    && adduser  -S appuser -G appgroup

USER appuser

EXPOSE 8080

ENTRYPOINT ["./validator-monitor"]
