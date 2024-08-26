FROM golang:1.22.5-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -o app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app /app
COPY --from=builder /app/.env /app

EXPOSE 8000
CMD ["./app"]