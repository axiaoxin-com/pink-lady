package webserver

import (
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/response"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
)

// 替换 gin 默认的 validator，更加友好的错误信息
func init() {
	binding.Validator = &goutils.GinStructValidator{}
}

// NewGinEngine 根据参数创建 gin 的 router engine
// middlewares 需要使用到的中间件列表，默认不为 engine 添加任何中间件
func NewGinEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	// set gin mode
	gin.SetMode(viper.GetString("server.mode"))

	engine := gin.New()

	// use middlewares
	for _, middleware := range middlewares {
		engine.Use(middleware)
	}

	return engine
}

// DefaultGinMiddlewares 默认的 gin server 使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		logging.GinLoggerWithConfig(logging.GinLoggerConfig{
			SkipPaths:           viper.GetStringSlice("logging.access_logger.skip_paths"),
			SkipPathRegexps:     viper.GetStringSlice("logging.access_logger.skip_path_regexps"),
			EnableDetails:       viper.GetBool("logging.access_logger.enable_details"),
			EnableContextKeys:   viper.GetBool("logging.access_logger.enable_context_keys"),
			EnableRequestHeader: viper.GetBool("logging.access_logger.enable_request_header"),
			EnableRequestForm:   viper.GetBool("logging.access_logger.enable_request_form"),
			EnableRequestBody:   viper.GetBool("logging.access_logger.enable_request_body"),
			EnableResponseBody:  viper.GetBool("logging.access_logger.enable_response_body"),
		}),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 500 JSON
		GinRecovery(response.Respond),
		// handler 超过指定时间没有处理完成直接返回 503
		GinTimeout(),
	}

	// 是否开启请求限频
	if viper.GetBool("server.ratelimiter") {
		fillInterval := viper.GetInt("ratelimiter.bucket_fill_every_microsecond")
		bucketCapacity := viper.GetInt("ratelimiter.bucket_capacity")
		// 添加限频中间件到中间件列表
		limiterType := strings.ToLower(viper.GetString("ratelimiter.type"))
		logging.Debugf(nil, "%s ratelimiter with bucket_fill_every_microsecond=%d bucket_capacity=%d", limiterType, fillInterval, bucketCapacity)
		if strings.HasPrefix(limiterType, "redis.") {
			which := strings.TrimPrefix(limiterType, "redis.")
			if rdb, err := goutils.RedisClient(which); err != nil {
				logging.Error(nil, "redis ratelimiter does not work. get redis client error:"+err.Error())
			} else {
				m = append(m, ratelimiter.GinRedisRatelimiter(rdb, fillInterval, bucketCapacity))
			}
		} else {
			m = append(m, ratelimiter.GinMemRatelimiter(fillInterval, bucketCapacity))
		}
	}
	return m
}
