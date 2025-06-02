package webserver

import (
	"fmt"
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
