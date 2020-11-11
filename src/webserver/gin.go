package webserver

import (
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/response"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go/extra"
	"github.com/spf13/viper"
)

func init() {
	// 替换 gin 默认的 validator，更加友好的错误信息
	binding.Validator = &goutils.GinStructValidator{}
	// causes the json binding Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	binding.EnableDecoderUseNumber = true

	// jsoniter 启动模糊模式来支持 PHP 传递过来的 JSON。容忍字符串和数字互转
	extra.RegisterFuzzyDecoders()
	// jsoniter 设置支持 private 的 field
	extra.SupportPrivateFields()
}

// NewGinEngine 根据参数创建 gin 的 router engine
// middlewares 需要使用到的中间件列表，默认不为 engine 添加任何中间件
func NewGinEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	// set gin mode
	gin.SetMode(viper.GetString("server.mode"))

	engine := gin.New()
	// ///a///b -> /a/b
	engine.RemoveExtraSlash = true

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
	}

	// 是否开启请求限频
	if viper.GetBool("ratelimiter.enable") {
		limiterConf := ratelimiter.GinRatelimiterConfig{
			// TODO: 需自行实现限频 key 的生成函数
			LimitKey: nil,
			// TODO: 需自行实现限频的具体令牌桶配置信息
			TokenBucketConfig: nil,
		}
		// 添加限频中间件到中间件列表
		limiterType := strings.ToLower(viper.GetString("ratelimiter.type"))
		if strings.HasPrefix(limiterType, "redis.") {
			which := strings.TrimPrefix(limiterType, "redis.")
			if rdb, err := goutils.RedisClient(which); err != nil {
				logging.Error(nil, "redis ratelimiter does not work. get redis client error:"+err.Error())
			} else {
				m = append(m, ratelimiter.GinRedisRatelimiter(rdb, limiterConf))
			}
		} else {
			m = append(m, ratelimiter.GinMemRatelimiter(limiterConf))
		}
		logging.Info(nil, "enable ratelimiter with type: "+limiterType)
	}
	return m
}
