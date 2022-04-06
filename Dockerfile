# syntax=docker/dockerfile:1

FROM golang:1.18.0-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /docker-build

EXPOSE 3000

CMD [ "/docker-build" ]
