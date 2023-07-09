FROM golang:1.20.5-bullseye

RUN mkdir /app

WORKDIR /app

ADD . /app

RUN go build -o server main.go

CMD ["/app/server"]
