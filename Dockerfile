FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pr-service ./cmd/app


FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/pr-service /app/pr-service

ENV APP_PORT=8080
EXPOSE 8080

CMD ["/app/pr-service"]
