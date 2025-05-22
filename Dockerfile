FROM golang:1.24.3 AS builder

WORKDIR /app

# Copy all files including .env and bot code
COPY . .

# Build the binary
RUN cd bot && CGO_ENABLED=0 go build -o /app/bin/bot

# Final stage
FROM gcr.io/distroless/static

COPY --from=builder /app/bin/bot /app/bin/bot
COPY --from=builder /app/bot/.env .env

ENTRYPOINT ["/app/bin/bot"]
