# gin-skeleton

Typically [gin](https://github.com/gin-gonic/gin)-based web application's organizational structure which in my mind.


### requirements

- **[[Go >= 1.9]](https://golang.org/doc/devel/release.html)**  dep is a dependency management tool for Go. It requires Go 1.9 or newer to compile.
- **[[golang/dep]](https://github.com/golang/dep)**  Go dependency management tool <https://golang.github.io/dep/>


### development

    cd gin-skeleton
    dep ensure
    cd app
    go run *.go


### dependencies

- **[[spf13/viper]](https://github.com/spf13/viper)**  Go configuration with fangs
- **[[sirupsen/logrus]](https://github.com/sirupsen/logrus)**  Structured, pluggable logging for Go.
- **[[fvbock/endless]](https://github.com/fvbock/endless)**  Zero downtime restarts for go servers (Drop in replacement for http.ListenAndServe)
- **[[jinzhu/gorm]](https://github.com/jinzhu/gorm)**  The fantastic ORM library for Golang, aims to be developer friendly <http://gorm.io>

### organizational structure

### todo

- [ ] CLI tool
- [ ] redis
- [ ] rabbitmq
- [ ] swagger
- [ ] validator
