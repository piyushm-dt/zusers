# syntax=docker/dockerfile:1
FROM golang:1.12.13

RUN apt update && apt install ca-certificates libgnutls30 -y

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080
EXPOSE 8081

RUN go build

CMD ["./zusers"]