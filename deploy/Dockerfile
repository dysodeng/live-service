
FROM golang

RUN mkdir -p /go/src/golang.org/x \
 && cd /go/src/golang.org/x && git clone https://github.com/golang/sys \
 && git clone https://github.com/golang/crypto.git \
 && git clone https://github.com/golang/time \
 && mkdir -p /go/src/github.com/aliyun && cd /go/src/github.com/aliyun \
 && git clone https://github.com/aliyun/alibaba-cloud-sdk-go

RUN go get github.com/go-sql-driver/mysql && go get -u github.com/jinzhu/gorm \
 && go get github.com/gin-gonic/gin \
 && go get github.com/gomodule/redigo/redis && go get github.com/astaxie/beego/cache \
 && go get github.com/dgrijalva/jwt-go && go get github.com/aliyun/aliyun-oss-go-sdk/oss

RUN go get gopkg.in/yaml.v2

WORKDIR /go/src/live-service

EXPOSE 8080

CMD nohup sh -c "go build && ./live-service"
