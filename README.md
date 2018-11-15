# pink-lady

The typically [gin](https://github.com/gin-gonic/gin)-based web application's organizational structure -> pink-lady.

The name comes from the Pink Lady which is a national standard cocktail with Gin as Base.

### Skeleton code organization structure

    > tree -I "vendor"
    .
    ├── app                                          // source code directory
    │   ├── apis                                    // write your apis at this directory
    │   │   ├── init.go                            // skeleton default api
    │   │   ├── init_test.go
    │   │   ├── router                             // gin router
    │   │   │   ├── router.go
    │   │   │   └── router_test.go
    │   │   ├── routes.go                          // register your handler function on url in here
    │   │   └── routes_test.go
    │   ├── app.yaml                                // configuration file
    │   ├── docs                                    // api docs generate by swag
    │   │   ├── docs.go
    │   │   └── swagger
    │   │       ├── swagger.json
    │   │       └── swagger.yaml
    │   ├── main.go                                 // main run a endless api server
    │   ├── middleware                              // skeleton default middlewares
    │   │   ├── errorhandler.go                    // handle 404 500 to return JSON
    │   │   ├── errorhandler_test.go
    │   │   ├── ginlogrus.go                       // logs use logrus
    │   │   ├── ginlogrus_test.go
    │   │   ├── init.go
    │   │   ├── requestid.go                       // log request id
    │   │   └── requestid_test.go
    │   ├── models                                  // write your models at here
    │   │   ├── init.go                            // provide a base model
    │   │   └── init_test.go
    │   ├── services                                // write your business handler at here
    │   │   ├── init.go
    │   │   ├── init_test.go
    │   │   └── retcode                            // write your business return code at here
    │   │       ├── retcode.go
    │   │       └── retcode_test.go
    │   └── utils                                   // add common utils at here
    │       ├── endless.go                          // provide a graceful stop server
    │       ├── gorequest.go                        // provide a http client
    │       ├── gorm.go                             // provide gorm db client
    │       ├── gorm_test.go
    │       ├── init.go
    │       ├── jsontime.go                         // provide a custom format time field for json
    │       ├── jsontime_test.go
    │       ├── logrus.go                           // provide a log
    │       ├── logrus_test.go
    │       ├── pagination.go                       // provide a pagination function
    │       ├── pagination_test.go
    │       ├── redis.go                            // provide a redis client
    │       ├── redis_test.go
    │       ├── response                            // provide united json response functions
    │       │   ├── response.go
    │       │   └── response_test.go
    │       ├── testing.go                          // provide GET/POST request function for testing
    │       ├── viper.go                            // provide configuration parser
    │       └── viper_test.go
    ├── Gopkg.lock                                   // dep file
    ├── Gopkg.toml                                   // dep file
    ├── misc                                         // write your tool scripts at here
    │   ├── pre-push.githook                        // a git pe-push hook for run test when push
    │   ├── README.md
    │   ├── release.sh                              // a tool for release the binary app
    │   └── supervisor.conf                         // a supervisor configure file demo
    └── README.md


### Develop requirements

- **[[Go >= 1.9]](https://golang.org/doc/devel/release.html)**  dep is a dependency management tool for Go. It requires Go 1.9 or newer to compile.
- **[[golang/dep]](https://github.com/golang/dep)**  Go dependency management tool <https://golang.github.io/dep/>
- **[[swaggo/swag]](https://github.com/swaggo/swag)**  Automatically generate RESTful API documentation with Swagger 2.0 for Go.
- **[[pilu/fresh]](https://github.com/pilu/fresh)**  Build and (re)start go web apps after saving/creating/deleting source files.



### Skeleton dependencies

- **[[spf13/viper]](https://github.com/spf13/viper)**  Go configuration with fangs
- **[[sirupsen/logrus]](https://github.com/sirupsen/logrus)**  Structured, pluggable logging for Go.
- **[[fvbock/endless]](https://github.com/fvbock/endless)**  Zero downtime restarts for go servers (Drop in replacement for http.ListenAndServe)
- **[[jinzhu/gorm]](https://github.com/jinzhu/gorm)**  The fantastic ORM library for Golang, aims to be developer friendly <http://gorm.io>
- **[[satori/go.uuid]](https://github.com/satori/go.uuid)** UUID package for Go
- **[[go-redis/redis]](https://github.com/go-redis/redis)**  Type-safe Redis client for Golang <https://godoc.org/github.com/go-redis/redis>
- **[[swaggo/gin-swagger]](https://github.com/swaggo/gin-swagger)**  gin middleware to automatically generate RESTful API documentation with Swagger 2.0.
- **[[gin-contrib/cors]](https://github.com/gin-contrib/cors)**  Official CORS gin's middleware
- **[[gin-contrib/sentry]](https://github.com/gin-contrib/sentry)**  Middleware to integrate with sentry crash reporting.
- **[[alicebob/miniredis]](https://github.com/alicebob/miniredis)**  Pure Go Redis server for Go unittests
- **[[parnurzeal/gorequest]](https://github.com/parnurzeal/gorequest)**  GoRequest -- Simplified HTTP client ( inspired by nodejs SuperAgent ) http://parnurzeal.github.io/gorequest/


### How to build web API server with pink-lady

1. Clone the pink-lady into your gopath and install dependebcies:

        cd $(go env GOPATH)/src
        git clone git@github.com:axiaoxin/pink-lady.git
        cd pink-lady
        dep ensure

    if you want to change the project path, you must change the import path too:

        mv pink-lady GOPATH/src/XXX/YYY/ZZZ
        cd GOPATH/src/XXX/YYY/ZZZ/app
        sed -i  "s|pink-lady|XXX/YYY/ZZZ|g"  `grep "pink-lady" . -rl`

2. Write your API in `apis` directory, you can create a file or a subdirectory as a package to save your gin API handler function, then register the handlers on url in `apis/routes.go` like the default `ping` API

3. Run develop server in `app` directory: (a convinent way to dynamic reload the server when code changing, you can use `fresh` to run server)

        cd app
        go test ./...
        go run main.go

4. Release your API server, run `misc/release.sh` will build a binary app in `build` directory if your tests all pass, and bump version, update api docs and make a git commit and add a tag with the version

        cd misc
        ./release.sh

If you wrote swag style comments you can generate the API docs in `app` directory manually:

    cd app
    swag init -g apis/init.go

You can test your API by curl or swagger API docs:<http://localhost:8080/x/apidocs/index.html>

### Develop suggestions

- You should put your API handler functions in `app/apis` as a single file or a package, then register the handlers in `app/apis/routes.go`
- Add your middleware in `app/middleware`, use it by adding to `app/apis/router/router.go`
- Define the database model in `app/models` by embed `BaseModel` in `app/models/init.go`
- Write reusable business code in `app/services` as single file or a package
- Define the return codes in `app/services/retcode/retcode.go`
- Add generic business independent tool type code in `app/utils`
- Add tool scripts etc in `misc`
- Write unit tests and doc for functions!
