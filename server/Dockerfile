FROM golang:1.20.5-bullseye

RUN apk update

RUN mkdir /app

WORKDIR /app

ADD . /app

RUN go build -o server main.go

CMD ["/app/server"]
