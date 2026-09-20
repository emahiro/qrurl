FROM golang:1.27-bookworm

RUN mkdir /app

WORKDIR /app

ADD . /app

RUN go build -o server main.go

CMD ["/app/server"]
