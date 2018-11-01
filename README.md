# gin-skeleton

Typically [gin](https://github.com/gin-gonic/gin)-based web application's organizational structure which in my mind.


### requirements

- **[[Go >= 1.9]](https://golang.org/doc/devel/release.html)**  dep is a dependency management tool for Go. It requires Go 1.9 or newer to compile.
- **[[golang/dep]](https://github.com/golang/dep)**  Go dependency management tool <https://golang.github.io/dep/>
- **[[swaggo/swag]](https://github.com/swaggo/swag)**  Automatically generate RESTful API documentation with Swagger 2.0 for Go.


### dependencies

- **[[spf13/viper]](https://github.com/spf13/viper)**  Go configuration with fangs
- **[[sirupsen/logrus]](https://github.com/sirupsen/logrus)**  Structured, pluggable logging for Go.
- **[[fvbock/endless]](https://github.com/fvbock/endless)**  Zero downtime restarts for go servers (Drop in replacement for http.ListenAndServe)
- **[[jinzhu/gorm]](https://github.com/jinzhu/gorm)**  The fantastic ORM library for Golang, aims to be developer friendly <http://gorm.io>
- **[[satori/go.uuid]](https://github.com/satori/go.uuid)** UUID package for Go
- **[[go-redis/redis]](https://github.com/go-redis/redis)**  Type-safe Redis client for Golang <https://godoc.org/github.com/go-redis/redis>
- **[[swaggo/gin-swagger]](https://github.com/swaggo/gin-swagger)**  gin middleware to automatically generate RESTful API documentation with Swagger 2.0.
- **[[gin-contrib/cors]](https://github.com/gin-contrib/cors)**  Official CORS gin's middleware
- **[[gin-contrib/sentry]](https://github.com/gin-contrib/sentry)**  Middleware to integrate with sentry crash reporting.

### organizational structure


### develop

    cd gin-skeleton
    dep ensure
    cd app
    go run *.go

    curl localhost:8080/x/ping

    # update API docs
    # localhost:8080/x/apidocs/index.html
    cd apis
    swag init -g *.go

### todo

- [ ] CLI tool
- [ ] sentry
- [ ] cron
- [ ] markdown flatpage
