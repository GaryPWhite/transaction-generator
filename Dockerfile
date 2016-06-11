FROM golang:1.7-alpine

COPY ./ $GOPATH

RUN go build speedtest.go

CMD $GOPATH/speedtest
