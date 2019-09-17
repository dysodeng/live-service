FROM golang:1.13

WORKDIR /go/src/live-service

ENV GOPROXY https://goproxy.io

ENV GO111MODULE on

EXPOSE 8080

CMD nohup sh -c "go mod download && go build && ./live-service"
