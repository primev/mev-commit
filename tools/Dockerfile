# tools/Dockerfile

FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /src

# pull in your module files (from tools/)
COPY tools/go.mod tools/go.sum ./
RUN go mod download

# now copy the *entire* repo so local replace‐paths exist
COPY . .

# build just your validator-monitor binary
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
    && adduser  -S appuser -G appgroup
USER appuser
EXPOSE 8080
ENTRYPOINT ["./validator-monitor"]
