FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["/app/main"]
