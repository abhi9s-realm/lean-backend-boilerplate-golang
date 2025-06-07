FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .

# Set default environment variables
ENV PORT=8080 \
    DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=lean_backend_boilerplate \
    LOG_LEVEL=info

EXPOSE 8080
CMD ["./api"]
