FROM golang:1.21-alpine

RUN apk add --no-cache git bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 5002

RUN go build -o dashboard-api .

CMD ["./dashboard-api"]
