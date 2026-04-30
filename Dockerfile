FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o note-app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/note-app .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./note-app"]
