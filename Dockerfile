FROM golang:1.13

WORKDIR /app

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /app
ADD ./go.mod /app
ADD ./go.sum /app

RUN export GOPROXY=https://goproxy.cn && go mod download

ADD . /app

RUN CGO_ENABLED=0 go build -o live-service

EXPOSE 8080

CMD /app/live-service
