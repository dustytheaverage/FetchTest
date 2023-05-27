FROM golang:latest

RUN mkdir /receipt
WORKDIR /receipt

COPY . .

RUN export GO111MODULE=on

RUN go mod init takehome/test
RUN go get .
RUN go run .

EXPOSE 8080