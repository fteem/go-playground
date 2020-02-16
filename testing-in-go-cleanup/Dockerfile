FROM golang:1.14-rc-stretch AS builder

ENV GO111MODULE=on

WORKDIR /app

# Install the project dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

CMD ["./bin/sh"]
