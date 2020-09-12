FROM golang:latest

WORKDIR $GOPATH/src/github.com/axiaoxin-com/pink-lady/src

COPY . $GOPATH/src/github.com/axiaoxin-com/pink-lady

ENV GOPROXY="https://goproxy.cn,direct"

RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init --dir ./ --generalInfo apis/apis.go --propertyStrategy snakecase --output ./apis/docs
RUN go build -o pink-lady-apiserver -tags=jsoniter

EXPOSE 4869 4870
ENTRYPOINT ["./pink-lady-apiserver", "-p", ".", "-c", "config.default", "-t", "toml"]
