#get a base image
FROM golang:1.12

RUN apt update && apt install ca-certificates libgnutls30 -y

WORKDIR E:/Go/zusers
COPY ./src .

RUN go get -d -v
RUN go build -v

CMD ["./zusers"]