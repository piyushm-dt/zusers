# syntax=docker/dockerfile:1
FROM golang:1.12.13

RUN apt update && apt install ca-certificates libgnutls30 -y

WORKDIR /app

COPY . ./

RUN go build -o /zusers

EXPOSE 8080

CMD [ "/zusers" ]