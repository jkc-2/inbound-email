# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o email-forwarder .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/email-forwarder .
COPY .env.example .env

EXPOSE 25

CMD ["./email-forwarder"]
