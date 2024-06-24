# Build stage
FROM golang:1.22-alpine3.20 AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache --update build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ./chat-app ./cmd/server/main.go

# Run stage
FROM alpine:3.20

ENV CGO_ENABLED=1

RUN apk add --no-cache --update build-base

WORKDIR /app

COPY --from=builder /app/chat-app .
COPY .env .

EXPOSE 8080

CMD ["./chat-app"]