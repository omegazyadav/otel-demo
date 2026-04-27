# --- Stage 1: Build ---
FROM golang:1.25-alpine AS builder

# Set workspace
WORKDIR /app

# Copy dependency files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary (statically linked for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o note-app .

# --- Stage 2: Run ---
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/note-app .
# Don't forget your HTML templates!
COPY --from=builder /app/templates ./templates

# Expose the app port
EXPOSE 8080

# Run the app
CMD ["./note-app"]
