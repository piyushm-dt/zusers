# syntax=docker/dockerfile:1
FROM golang:1.12.0

RUN apt update && apt install ca-certificates libgnutls30 -y

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download
RUN go build -o main .

EXPOSE 8080

CMD [ "./main" ]