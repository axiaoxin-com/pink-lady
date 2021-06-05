package webserver

import (
	"html/template"
	"net/http"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/pink-lady/response"
	"github.com/axiaoxin-com/pink-lady/statics"
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

	// set template funcmap, must befor load templates
	engine.SetFuncMap(TemplFuncs)

	// load html template
	tmplPath := viper.GetString("statics.tmplpath")
	if tmplPath != "" {
		t, err := template.ParseFS(statics.Files, tmplPath)
		if err != nil {
			panic(err)
		}
		engine.SetHTMLTemplate(t)
	}

	// register statics
	staticsURL := viper.GetString("statics.url")
	if staticsURL != "" {
		engine.StaticFS(staticsURL, http.FS(statics.Files))
	}

	return engine
}

// DefaultGinMiddlewares 默认的 gin server 使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		GinLogMiddleware(),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 500 JSON
		GinRecovery(response.Respond),
	}

	// 配置开启请求限频则添加限频中间件
	if viper.GetBool("ratelimiter.enable") {
		m = append(m, GinRatelimitMiddleware())
	}
	return m
}
