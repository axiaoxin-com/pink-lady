package webserver

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	configFile string = ""
)

// InitWithConfigFile 根据 webserver 配置文件初始化 webserver
func InitWithConfigFile(configPath, configName, configType string) {
	configFile = configPath + "/" + configName + "." + configType

	// 加载配置文件内容到 viper 中以便使用
	if err := goutils.InitViper(configPath, configName, configType, func(e fsnotify.Event) {
		logging.Warn(nil, "Config file changed:"+e.Name)
		logging.SetLevel(viper.GetString("logging.level"))
	}); err != nil {
		// 文件不存在时 1 使用默认配置，其他 err 直接 panic
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(err)
		}
		logging.Error(nil, "Init viper error:"+err.Error())
	}

	// 设置 viper 中 webserver 配置项默认值
	viper.SetDefault("env", "dev")

	viper.SetDefault("server.addr", ":4869")
	viper.SetDefault("server.mode", gin.ReleaseMode)
	viper.SetDefault("server.pprof", true)
	viper.SetDefault("server.handler_timeout", 5)

	viper.SetDefault("apidocs.title", "pink-lady swagger apidocs")
	viper.SetDefault("apidocs.desc", "Using pink-lady to develop gin app on fly.")
	viper.SetDefault("apidocs.host", "localhost:4869")
	viper.SetDefault("apidocs.basepath", "/")
	viper.SetDefault("apidocs.schemes", []string{"http"})

	viper.SetDefault("basic_auth.username", "admin")
	viper.SetDefault("basic_auth.password", "admin")

	// 初始化 sentry 并创建 sentry 客户端
	sentryDSN := viper.GetString("sentry.dsn")
	if sentryDSN == "" {
		sentryDSN = os.Getenv(logging.SentryDSNEnvKey)
	}
	sentryDebug := true
	if viper.GetString("server.mode") == "release" {
		sentryDebug = false
	}
	if viper.GetBool("sentry.debug") {
		sentryDebug = true
	}
	logging.Debug(nil, "Sentry use dns: "+sentryDSN)
	sentry, err := logging.NewSentryClient(sentryDSN, sentryDebug)
	if err != nil {
		logging.Error(nil, "Sentry client create error:"+err.Error())
	}

	// 根据配置创建 logging 的 logger 并将 logging 的默认 logger 替换为当前创建的 logger
	logger, err := logging.NewLogger(logging.Options{
		Level:             viper.GetString("logging.level"),
		Format:            viper.GetString("logging.format"),
		OutputPaths:       viper.GetStringSlice("logging.output_paths"),
		DisableCaller:     viper.GetBool("logging.disable_caller"),
		DisableStacktrace: viper.GetBool("logging.disable_stacktrace"),
		AtomicLevelServer: logging.AtomicLevelServerOption{
			Addr:     viper.GetString("logging.atomic_level_server.addr"),
			Path:     viper.GetString("logging.atomic_level_server.path"),
			Username: viper.GetString("basic_auth.username"),
			Password: viper.GetString("basic_auth.password"),
		},
		SentryClient: sentry,
	})
	if err != nil {
		logging.Error(nil, "Logger create error:"+err.Error())
	} else {
		logging.ReplaceLogger(logger)
	}
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
	handlerTimeout := viper.GetDuration("server.handler_timeout") * time.Second
	addr := viper.GetString("server.addr")
	srv := &http.Server{
		Addr:    addr,
		Handler: app,
		//	Handler: http.TimeoutHandler(app, handlerTimeout, ""),
	}
	// Shutdown 时关闭 db 和 redis 连接
	srv.RegisterOnShutdown(func() {
		goutils.CloseGormInstances()
		goutils.CloseRedisInstances()
	})

	// 启动 http server
	go func() {
		var ln net.Listener
		var err error
		if strings.ToLower(strings.Split(addr, ":")[0]) == "unix" {
			ln, err = net.Listen("unix", strings.Split(addr, ":")[1])
			if err != nil {
				panic(err)
			}
		} else {
			ln, err = net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
		}
		if err := srv.Serve(ln); err != nil {
			logging.Fatal(nil, err.Error())
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
	logging.Infof(nil, "Server is shutting down.")

	// 创建一个 context 用于通知 server 有 writeTimeout 秒的时间结束当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Fatal(nil, "Server shutdown with error: "+err.Error())
	}
	logging.Info(nil, "Server exit.")
}
