FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (disable CGO for static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /booking-backend ./cmd/main.go

# Runtime image
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /booking-backend .

# Expose application port (default 8080, can be overridden)
EXPOSE 8080

# Run the server
CMD ["./booking-backend"]
