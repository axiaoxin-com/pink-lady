package webserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/chai2010/gettext-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

// GinBasicAuth gin 的基础认证中间件
// 加到 gin app 的路由中可以对该路由添加 basic auth 登录验证
// 传入 username 和 password 对可以替换默认的 username 和 password
func GinBasicAuth(args ...string) gin.HandlerFunc {
	username := viper.GetString("basic_auth.username")
	password := viper.GetString("basic_auth.password")
	if len(args) == 2 {
		username = args[0]
		password = args[1]
	}
	logging.Debug(nil, "Basic auth username:"+username+" password:"+password)
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}

// GinRecovery gin recovery 中间件
// save err in context and abort with recoveryHandler
func GinRecovery(
	recoveryHandler ...func(c *gin.Context, status int, data interface{}, err error, extraMsgs ...interface{}),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			status := c.Writer.Status()
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				status = http.StatusInternalServerError
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if brokenPipe {
					status = 499
					// log warning
					logging.Warnf(c, "Broken pipe: %v\n%s", err, string(debug.Stack()))
					if len(recoveryHandler) > 0 {
						c.Abort()
						recoveryHandler[0](c, status, nil, errors.New(http.StatusText(status)))
						return
					}
					c.AbortWithStatus(status)
					return
				}

				// save err in context when panic
				c.Error(fmt.Errorf("Recovery from panic: %v\n%s", err, string(debug.Stack())))
			}

			// 状态码大于 400 的以 status handler 进行响应
			if status >= 400 {
				// 有 handler 时使用 handler 返回
				if len(recoveryHandler) > 0 {
					c.Abort()
					recoveryHandler[0](c, status, nil, errors.New(http.StatusText(status)))
					return
				}
				// 否则只返回状态码
				c.AbortWithStatus(status)
				return
			}
		}()
		c.Next()
	}
}

// GinLogMiddleware 日志中间件
// 可根据实际需求自行修改定制
func GinLogMiddleware() gin.HandlerFunc {
	logging.CtxLoggerName = logging.Ctxkey("ctx_logger")
	logging.TraceIDKeyname = logging.Ctxkey("trace_id")
	logging.TraceIDPrefix = "logging_"
	loggerMiddleware := logging.GinLoggerWithConfig(logging.GinLoggerConfig{
		SkipPaths:           viper.GetStringSlice("logging.access_logger.skip_paths"),
		SkipPathRegexps:     viper.GetStringSlice("logging.access_logger.skip_path_regexps"),
		EnableDetails:       viper.GetBool("logging.access_logger.enable_details"),
		EnableContextKeys:   viper.GetBool("logging.access_logger.enable_context_keys"),
		EnableRequestHeader: viper.GetBool("logging.access_logger.enable_request_header"),
		EnableRequestForm:   viper.GetBool("logging.access_logger.enable_request_form"),
		EnableRequestBody:   viper.GetBool("logging.access_logger.enable_request_body"),
		EnableResponseBody:  viper.GetBool("logging.access_logger.enable_response_body"),
		SlowThreshold:       viper.GetDuration("logging.access_logger.slow_threshold") * time.Millisecond,
		TraceIDFunc:         nil,
		InitFieldsFunc:      nil,
		Formatter:           nil,
	})
	return loggerMiddleware
}

// GinRatelimitMiddleware 限频中间件
// 需先实现对应的 TODO ，可根据实际需求自行修改定制
func GinRatelimitMiddleware() gin.HandlerFunc {
	limiterConf := ratelimiter.GinRatelimiterConfig{
		// TODO: you should implement this function by yourself
		LimitKey: ratelimiter.DefaultGinLimitKey,
		// TODO: you should implement this function by yourself
		LimitedHandler: ratelimiter.DefaultGinLimitedHandler,
		// TODO: you should implement this function by yourself
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			// 每1秒填充1个token，桶容量为100（1秒最多100次请求）
			return time.Second * 1, 100
		},
	}

	// 根据 viper 中的配置信息选择限流类型
	var limiterMiddleware gin.HandlerFunc
	limiterType := strings.ToLower(viper.GetString("ratelimiter.type"))
	logging.Info(nil, "enable ratelimiter with type: "+limiterType)
	if strings.HasPrefix(limiterType, "redis.") {
		which := strings.TrimPrefix(limiterType, "redis.")
		if rdb, err := goutils.RedisClient(which); err != nil {
			panic("redis ratelimiter does not work. get redis client error:" + err.Error())
		} else {
			limiterMiddleware = ratelimiter.GinRedisRatelimiter(rdb, limiterConf)
		}
	} else {
		limiterMiddleware = ratelimiter.GinMemRatelimiter(limiterConf)
	}
	return limiterMiddleware
}

// I18nLangTags i18n支持的语言列表
var I18nLangTags = []language.Tag{
	language.SimplifiedChinese,
	language.English,
	language.TraditionalChinese,
	language.German,
	language.Spanish,
	language.French,
	language.Italian,
	language.Japanese,
	language.Korean,
	language.Portuguese,
	language.Russian,
	language.Turkish,
	language.Vietnamese,
	language.Arabic,
	language.Hindi,
	language.Bengali,
	language.Indonesian,
	language.Thai,
}

// i18nGettexters 多语言 gettexter
var i18nGettexters = map[string]gettext.Gettexter{}

// gettexter对象池
var (
	GettexterPoolEnglish            *sync.Pool
	GettexterPoolTraditionalChinese *sync.Pool
	GettexterPoolGerman             *sync.Pool
	GettexterPoolSpanish            *sync.Pool
	GettexterPoolFrench             *sync.Pool
	GettexterPoolItalian            *sync.Pool
	GettexterPoolJapanese           *sync.Pool
	GettexterPoolKorean             *sync.Pool
	GettexterPoolPortuguese         *sync.Pool
	GettexterPoolRussian            *sync.Pool
	GettexterPoolTurkish            *sync.Pool
	GettexterPoolVietnamese         *sync.Pool
	GettexterPoolArabic             *sync.Pool
	GettexterPoolHindi              *sync.Pool
	GettexterPoolBengali            *sync.Pool
	GettexterPoolIndonesian         *sync.Pool
	GettexterPoolThai               *sync.Pool
)

// CtxI18n 从context获取多语言
func CtxI18n(c *gin.Context, msg interface{}) string {
	if c == nil {
		return fmt.Sprint(msg)
	}
	return LangI18n(c.GetString("lang"), msg)
}

// LangI18n 获取多语言（模板方法使用）
func LangI18n(lang string, msg interface{}) string {
	msgid := fmt.Sprint(msg)
	switch lang {
	case language.English.String():
		if GettexterPoolEnglish != nil {
			gettexter := GettexterPoolEnglish.Get().(gettext.Gettexter)
			defer GettexterPoolEnglish.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.TraditionalChinese.String():
		if GettexterPoolTraditionalChinese != nil {
			gettexter := GettexterPoolTraditionalChinese.Get().(gettext.Gettexter)
			defer GettexterPoolTraditionalChinese.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.German.String():
		if GettexterPoolGerman != nil {
			gettexter := GettexterPoolGerman.Get().(gettext.Gettexter)
			defer GettexterPoolGerman.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Spanish.String():
		if GettexterPoolSpanish != nil {
			gettexter := GettexterPoolSpanish.Get().(gettext.Gettexter)
			defer GettexterPoolSpanish.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.French.String():
		if GettexterPoolFrench != nil {
			gettexter := GettexterPoolFrench.Get().(gettext.Gettexter)
			defer GettexterPoolFrench.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Italian.String():
		if GettexterPoolItalian != nil {
			gettexter := GettexterPoolItalian.Get().(gettext.Gettexter)
			defer GettexterPoolItalian.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Japanese.String():
		if GettexterPoolJapanese != nil {
			gettexter := GettexterPoolJapanese.Get().(gettext.Gettexter)
			defer GettexterPoolJapanese.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Korean.String():
		if GettexterPoolKorean != nil {
			gettexter := GettexterPoolKorean.Get().(gettext.Gettexter)
			defer GettexterPoolKorean.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Portuguese.String():
		if GettexterPoolPortuguese != nil {
			gettexter := GettexterPoolPortuguese.Get().(gettext.Gettexter)
			defer GettexterPoolPortuguese.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Russian.String():
		if GettexterPoolRussian != nil {
			gettexter := GettexterPoolRussian.Get().(gettext.Gettexter)
			defer GettexterPoolRussian.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Turkish.String():
		if GettexterPoolTurkish != nil {
			gettexter := GettexterPoolTurkish.Get().(gettext.Gettexter)
			defer GettexterPoolTurkish.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Vietnamese.String():
		if GettexterPoolVietnamese != nil {
			gettexter := GettexterPoolVietnamese.Get().(gettext.Gettexter)
			defer GettexterPoolVietnamese.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Arabic.String():
		if GettexterPoolArabic != nil {
			gettexter := GettexterPoolArabic.Get().(gettext.Gettexter)
			defer GettexterPoolArabic.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Hindi.String():
		if GettexterPoolHindi != nil {
			gettexter := GettexterPoolHindi.Get().(gettext.Gettexter)
			defer GettexterPoolHindi.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Bengali.String():
		if GettexterPoolVietnamese != nil {
			gettexter := GettexterPoolVietnamese.Get().(gettext.Gettexter)
			defer GettexterPoolVietnamese.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Indonesian.String():
		if GettexterPoolIndonesian != nil {
			gettexter := GettexterPoolIndonesian.Get().(gettext.Gettexter)
			defer GettexterPoolIndonesian.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	case language.Thai.String():
		if GettexterPoolThai != nil {
			gettexter := GettexterPoolThai.Get().(gettext.Gettexter)
			defer GettexterPoolThai.Put(gettexter)
			return gettexter.Gettext(msgid)
		}
	}

	return msgid
}

// I18nString 需要国际化的string
type I18nString string

func (s I18nString) String() string {
	return string(s)
}

// GinSetLanguage 设置i18n语言
func GinSetLanguage(supportedLangTags ...language.Tag) gin.HandlerFunc {
	if len(supportedLangTags) == 0 {
		supportedLangTags = I18nLangTags
	}
	matcher := language.NewMatcher(supportedLangTags)

	go func() {
		// 初始化全局变量i18nGettexters
		gettexterMap := map[string]gettext.Gettexter{}
		for _, lt := range supportedLangTags {
			gettexter := gettext.New(viper.GetString("i18n.domain"), viper.GetString("i18n.path")).SetLanguage(lt.String())
			gettexterMap[lt.String()] = gettexter
		}
		i18nGettexters = gettexterMap
		logging.Infow(nil, "GinSetLanguage Init i18nGettexters", "i18nGettexters", i18nGettexters)

		// 初始化gettexter对象池，保存和复用临时对象，减少内存分配，降低 GC 压力。
		GettexterPoolEnglish = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.English.String()]
			},
		}

		GettexterPoolTraditionalChinese = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.TraditionalChinese.String()]
			},
		}

		GettexterPoolGerman = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.German.String()]
			},
		}

		GettexterPoolSpanish = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Spanish.String()]
			},
		}

		GettexterPoolFrench = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.French.String()]
			},
		}

		GettexterPoolItalian = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Italian.String()]
			},
		}

		GettexterPoolJapanese = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Japanese.String()]
			},
		}

		GettexterPoolKorean = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Korean.String()]
			},
		}

		GettexterPoolPortuguese = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Portuguese.String()]
			},
		}

		GettexterPoolRussian = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Russian.String()]
			},
		}

		GettexterPoolTurkish = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Turkish.String()]
			},
		}

		GettexterPoolVietnamese = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Vietnamese.String()]
			},
		}

		GettexterPoolArabic = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Arabic.String()]
			},
		}

		GettexterPoolHindi = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Hindi.String()]
			},
		}

		GettexterPoolBengali = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Bengali.String()]
			},
		}

		GettexterPoolIndonesian = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Indonesian.String()]
			},
		}

		GettexterPoolThai = &sync.Pool{
			New: func() interface{} {
				return i18nGettexters[language.Thai.String()]
			},
		}
	}()

	return func(c *gin.Context) {
		var err error
		var langTags []language.Tag
		saveLangInCookie := false
		cookieName := "pink-lady.lang"

		if !strings.HasPrefix(c.Request.RequestURI, viper.GetString("statics.url")) {
			// 设置指定语言
			// 尝试从url获取lang参数
			lang := c.Query("lang")
			if lang != "" {
				langTags, _, err = language.ParseAcceptLanguage(lang)
				if err != nil {
					logging.Warn(c, "GinSetLanguage ParseAcceptLanguage from query error:"+err.Error())
				} else {
					saveLangInCookie = true
					// logging.Debugf(c, "GinSetLanguage ParseAcceptLanguage from query langTags:%+v", langTags)
				}
			}
			// 尝试从cookie获取lang
			if langTags == nil {
				cookieLang, err := c.Cookie(cookieName)
				if err != nil {
					logging.Debug(c, "GinSetLanguage get cookieLang error:"+err.Error())
				} else {
					langTags, _, err = language.ParseAcceptLanguage(cookieLang)
					if err != nil {
						logging.Warn(c, "GinSetLanguage ParseAcceptLanguage from cookieLang error:"+err.Error())
					}
				}
			}

			// 从请求头获取accept-language寻找最佳匹配
			if langTags == nil {
				langTags, _, err = language.ParseAcceptLanguage(c.Request.Header.Get("Accept-Language"))
				if err != nil {
					logging.Warn(c, "GinSetLanguage ParseAcceptLanguage from header error:"+err.Error())
				} else {
					// logging.Debugf(c, "GinSetLanguage ParseAcceptLanguage from header langTags:%+v supportedLangTags:%+v", langTags, supportedLangTags)
				}
			}

			if langTags == nil {
				lang = language.English.String()
			} else {
				code, _, _ := matcher.Match(langTags...)
				lang = code.String()
			}
			// logging.Debug(c, "GinSetLanguage match lang="+lang)

			if saveLangInCookie {
				c.SetCookie(cookieName, lang, 3153600000, "", "", false, true) // 100年过期
				// logging.Debugf(c, "GinSetLanguage set lang:%v to cookie", lang)
			}

			// 设置gettexter
			c.Set("lang", lang)
		}

		c.Next()
	}
}
