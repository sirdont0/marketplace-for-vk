FROM golang:1.23 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o marketplace ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/marketplace .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./marketplace"]