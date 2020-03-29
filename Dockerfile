FROM golang:latest

RUN go get github.com/ajpatri/mirror-mirror
ENTRYPOINT /go/bin/mirror-mirror
