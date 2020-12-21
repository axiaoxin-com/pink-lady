package webserver

import (
	"html/template"
	"os"
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go/extra"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
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
	tmplPath := viper.GetString("static.tmpl_statik_path")
	if tmplPath != "" {
		t, err := GinLoadHTMLTemplate(tmplPath)
		if err != nil {
			panic(err)
		}
		engine.SetHTMLTemplate(t)
	}

	// register static
	staticURL := viper.GetString("static.url")
	if staticURL != "" {
		logging.Debugf(nil, "Static url: %s", staticURL)
		engine.StaticFS(staticURL, StatikFS)
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

// GinLoadHTMLTemplate 获取 gin 的 template
// tmplPath 是以 statik 的目录为根路径，以绝对路径表示
func GinLoadHTMLTemplate(tmplPath string) (*template.Template, error) {
	t := template.New("")
	// type WalkFunc func(path string, info os.FileInfo, err error) error
	err := fs.Walk(StatikFS, tmplPath, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "statik fs Walk error")
		}
		if info.IsDir() {
			return nil
		}
		filename := info.Name()
		if strings.HasSuffix(filename, ".tmpl") || strings.HasSuffix(filename, ".html") || strings.HasSuffix(filename, ".gohtml") || strings.HasSuffix(filename, ".gotmpl") {
			content, err := fs.ReadFile(StatikFS, filepath)
			if err != nil {
				return errors.Wrap(err, "statik fs ReadFile error")
			}
			t, err = t.New(filepath).Parse(string(content))
			if err != nil {
				return errors.Wrap(err, "template parse error")
			}
		}
		return nil
	})
	return t, err
}
