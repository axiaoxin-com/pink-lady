package webserver

import (
	"fmt"
	"html"
	"strings"
	"sync"

	"github.com/axiaoxin-com/logging"
	"github.com/chai2010/gettext-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

// I18nString 需要国际化的string
type I18nString string

func (s I18nString) String() string {
	return string(s)
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

// gettexter对象池
// 全局变量：语言 -> 对应 sync.Pool
var GettexterPools map[string]*sync.Pool

// CtxI18n 从context获取多语言
func CtxI18n(c *gin.Context, msg any) string {
	if c == nil {
		return fmt.Sprint(msg)
	}
	return LangI18n(c.GetString("lang"), msg)
}

// LangI18n 获取多语言（模板方法使用）
func LangI18n(lang string, msg any) string {
	msgid := fmt.Sprint(msg)

	if pool, ok := GettexterPools[lang]; ok {
		gettexter := pool.Get().(gettext.Gettexter)
		defer pool.Put(gettexter)

		// 翻译结果
		text := gettexter.Gettext(msgid)

		// 处理机器翻译中带的 HTML 实体，例如 &#39; 转成 '
		return html.UnescapeString(text)
	}
	// 找不到语言就直接返回原文
	return msgid
}

// GinSetLanguage 设置i18n语言
func GinSetLanguage(supportedLangTags ...language.Tag) gin.HandlerFunc {
	if len(supportedLangTags) == 0 {
		supportedLangTags = I18nLangTags
	}
	matcher := language.NewMatcher(supportedLangTags)

	go func() {
		// 初始化全局变量GettexterPools
		GettexterPools = make(map[string]*sync.Pool)
		for _, lt := range supportedLangTags {
			gettexter := gettext.New(viper.GetString("i18n.domain"), viper.GetString("i18n.path")).SetLanguage(lt.String())
			GettexterPools[lt.String()] = &sync.Pool{
				New: func() any {
					return gettexter
				},
			}
			logging.Infof(nil, "GinSetLanguage Init GettexterPools: %v", lt.String())
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
