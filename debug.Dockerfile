FROM golang:1.17.2

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /src