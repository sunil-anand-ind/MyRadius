FROM golang:alpine

WORKDIR /app


COPY redisSubscriber.go ./

RUN go mod init redis-control-plane-logger

RUN go mod tidy

RUN go build .

CMD [./redis-control-plane-logger]