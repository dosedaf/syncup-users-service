FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

COPY .env .
COPY database/migration ./migrations

EXPOSE 3000

CMD ["/app/server"]