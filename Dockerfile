# Use the official Go image as the base image
FROM golang:1.22.5-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -o app
EXPOSE 8000
CMD ["./app"]