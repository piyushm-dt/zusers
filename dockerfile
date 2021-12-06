#get a base image
FROM golang:1.16-buster

WORKDIR E:/Go/zusers
COPY ./src .

RUN go get -d -v
RUN go build -v

CMD ["./zusers"]