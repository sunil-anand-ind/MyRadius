FROM golang:alpine

WORKDIR /app

RUN mkdir datastore
COPY server.go ./
COPY datastore/* ./datastore/

RUN go mod init radius-control-plane

RUN go mod tidy

RUN go build .

EXPOSE 1813

CMD [./radius-control-plane]