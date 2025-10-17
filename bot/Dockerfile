FROM golang:1.25.3 AS builder

WORKDIR /app

# Copy all files including .env and bot code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o ./bin

ENTRYPOINT ["/build"]

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/.env .env

# Copy built binary
COPY --from=builder /app/bin bin

# Run binary
ENTRYPOINT ["./bin"]