# Build stage
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /chat-app cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /chat-app .
COPY .env .

EXPOSE 8080

CMD ["./chat-app"]