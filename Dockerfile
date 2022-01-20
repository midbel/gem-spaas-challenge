# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /spaas-gem

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN mkdir bin
RUN go build -o ./bin/gem cmd/gem/main.go

EXPOSE 8888

CMD ["./bin/gem"]
