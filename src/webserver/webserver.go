package webserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	// GinPprofURLPath 设置 gin 中的 pprof url 注册路径，可以通过外部修改
	GinPprofURLPath        = "/x/pprof"
	configFile      string = ""
)

// InitViperConfig 加载 server 配置文件到 viper
func InitViperConfig(configPath, configName, configType string) {
	configFile = configPath + "/" + configName + "." + configType

	if err := goutils.InitViper(configPath, configName, configType, func(e fsnotify.Event) {
		logging.Warn(nil, "Config file changed:"+e.Name)
	}); err != nil {
		logging.Error(nil, "Init viper error:"+err.Error())
	}

	// 设置配置默认值
	viper.SetDefault("env", "dev")

	viper.SetDefault("server.addr", ":4869")
	viper.SetDefault("server.mode", gin.ReleaseMode)
	viper.SetDefault("server.pprof", true)
	viper.SetDefault("server.read_timeout", 5)  // 服务器从 accept 到读取 body 的超时时间（秒）
	viper.SetDefault("server.write_timeout", 5) // 服务器从 accept 到写 response 的超时时间（秒）

	viper.SetDefault("apidocs.title", "pink-lady swagger apidocs")
	viper.SetDefault("apidocs.desc", "Using pink-lady to develop gin app on fly.")
	viper.SetDefault("apidocs.host", "localhost:4869")
	viper.SetDefault("apidocs.basepath", "/")
	viper.SetDefault("apidocs.schemes", []string{"http"})

	viper.SetDefault("basic_auth.username", "admin")
	viper.SetDefault("basic_auth.password", "admin")
}

// NewGinEngine 根据参数创建 gin 的 router engine
// middlewares 需要使用到的中间件列表，默认不为 engine 添加任何中间件
func NewGinEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	// set gin mode
	mode := viper.GetString("server.mode")
	if mode == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}
	gin.SetMode(mode)

	engine := gin.New()

	// use middlewares
	for _, middleware := range middlewares {
		engine.Use(middleware)
	}

	if viper.GetBool("server.pprof") {
		pprof.Register(engine, GinPprofURLPath)
	}
	return engine
}

// Run 以 viper 加载的 app 配置启动运行 http.Handler 的 app
// 注意：这里依赖 viper ，必须在外部先对 viper 配置进行加载
func Run(app http.Handler, routesRegister func(http.Handler)) {
	// 结束时关闭 db 连接
	defer goutils.CloseGormInstances()

	// 判断是否加载 viper 配置
	if !goutils.IsInitedViper() {
		panic("Running server must init viper by config file first!")
	}

	// 注册 api 路由
	routesRegister(app)

	// 创建 server
	addr := viper.GetString("server.addr")
	readTimeout := viper.GetInt("server.read_timeout")
	writeTimeout := viper.GetInt("server.write_timeout")
	srv := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}

	// 启动 http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatal(nil, "Server start error:"+err.Error())
		}
	}()
	logging.Infof(nil, "Server is running on %s with config file %s", srv.Addr, configFile)

	// 监听中断信号， WriteTimeout 时间后优雅关闭服务
	// syscall.SIGTERM 不带参数的 kill 命令
	// syscall.SIGINT ctrl-c kill -2
	// syscall.SIGKILL 是 kill -9 无法捕获这个信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Infof(nil, "Server will shutdown after %d seconds", writeTimeout)

	// 创建一个 context 用于通知 server 有 writeTimeout 秒的时间结束当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(writeTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Fatal(nil, "Server shutdown with error: "+err.Error())
	}
	logging.Info(nil, "Server shutdown")
}

// GinBasicAuth 加到 gin app 的路由中可以对该路由添加 basic auth 登录验证
// 传入 username 和 password 对可以替换默认的 username 和 password
func GinBasicAuth(args ...string) gin.HandlerFunc {
	username := viper.GetString("basic_auth.username")
	password := viper.GetString("basic_auth.password")
	switch len(args) {
	case 2:
		username = args[0]
		password = args[1]
	case 0:
		logging.Info(nil, "Set basic auth using the username and password in the configuration file.")
	default:
		logging.Error(nil, "Wrong number of username and password pair.")
	}
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}

// DefaultGinMiddlewares 默认的 gin 中间件
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		logging.GinTraceID(logging.GetGinTraceIDFromHeader, logging.GetGinTraceIDFromQueryString, logging.GetGinTraceIDFromPostForm),
		gin.Logger(),
	}
	return m
}
