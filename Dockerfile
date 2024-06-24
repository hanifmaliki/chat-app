# Build stage
FROM golang:1.22-alpine3.20 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /chat-app cmd/server/main.go

# Run stage
FROM alpine:3.20

ENV CGO_ENABLED=1

WORKDIR /root/

COPY --from=builder /chat-app .
COPY .env .

EXPOSE 8080

CMD ["./chat-app"]