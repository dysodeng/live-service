FROM golang:1.13

WORKDIR /go/src/live-service

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

ADD go.mod .
ADD go.sum .

RUN go mod download

EXPOSE 8080

CMD nohup sh -c "go build && ./live-service"
