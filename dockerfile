#get a base image
FROM golang:1.12

RUN apt update && apt install ca-certificates libgnutls30 -y

WORKDIR E:/Go/zusers
COPY ./ .

RUN go get -d -v
RUN go build -v

EXPOSE 12345

CMD ["./zusers"]