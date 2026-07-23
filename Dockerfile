# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /root/

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy compiled binary
COPY --from=builder /app/app .

# Copy Swagger docs (optional)
COPY --from=builder /app/docs ./docs

# Copy migrations (if needed)
COPY --from=builder /app/migrations ./migrations

# Expose application port
EXPOSE 8080

# Start the application
CMD ["./app"]
