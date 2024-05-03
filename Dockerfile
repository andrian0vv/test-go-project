FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o service ./cmd/service

FROM alpine:3.19.1

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/service .

EXPOSE 80

CMD ["./service"]
