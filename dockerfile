FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main

CMD [ "/app/main" ]