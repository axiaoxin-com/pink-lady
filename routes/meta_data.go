package routes

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/x-way/crawlerdetect"
)

const (
	// AppID common.app id
	AppID = 0

	// SiteName TODO: set your site name
	SiteName = webserver.I18nString("pink-lady")
	Slogan   = webserver.I18nString("")
)

// BuildID 编译ID -ldflags添加当前时间
var BuildID = ""

// TimeLocationBeijing 北京时区
var TimeLocationBeijing = time.FixedZone("BeijingTime", int((time.Hour * 8).Seconds()))

// MetaData 元数据
type MetaData struct {
	SiteName        string
	Slogan          string
	HostURL         string
	BuildID         string
	Env             string
	AppID           int
	AdsenseID       string
	SysNotice       string
	SysNoticeQRText string
	Title           string
	IsCrawler       bool
	Lang            string
	Keywords        []string
	BaseDesc        string
	BootswatchTheme string
	I18n            bool
	Beian           string
	AuthorName      string
	AuthorURL       string
	StaticsURL      string
	StaticsSelfhost bool
	FlatpagesConfig *FlatpagesConfig
	BaiduTongJiID   string
	GtagID          string
	ClarityID       string
	CloudflareToken string
	UmamiWebsiteID  string
	CanonicalURL    string
	ShowAbout       bool
	SinceYear       string
	FriendLinkMap   map[string]string
}

// NewMetaData 返回页面元数据
func NewMetaData(c *gin.Context, title string) (m *MetaData) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	ua := c.GetHeader("User-Agent")
	isCrawler := crawlerdetect.IsCrawler(ua)

	m = &MetaData{
		SiteName:        webserver.CtxI18n(c, SiteName),
		Slogan:          webserver.CtxI18n(c, Slogan),
		BuildID:         BuildID,
		Env:             viper.GetString("env"),
		AppID:           AppID,
		AdsenseID:       viper.GetString("server.adsense_id"),
		Title:           title,
		IsCrawler:       isCrawler,
		Lang:            c.GetString("lang"),
		BaseDesc:        webserver.CtxI18n(c, "pink-lady是一个golang gin的web开发模板"),
		BootswatchTheme: "cosmo",
		I18n:            viper.GetBool("i18n.enable"),
		Beian:           viper.GetString("server.beian"),
		AuthorName:      viper.GetString("author.name"),
		AuthorURL:       viper.GetString("author.url"),
		StaticsSelfhost: viper.GetBool("statics.selfhost"),
		FlatpagesConfig: flatpagesConfig,
		BaiduTongJiID:   viper.GetString("server.baidu_tongji_id"),
		GtagID:          viper.GetString("server.gtag_id"),
		ClarityID:       viper.GetString("server.clarity_id"),
		CloudflareToken: viper.GetString("server.cloudflare_token"),
		UmamiWebsiteID:  viper.GetString("server.umami_website_id"),
		ShowAbout:       viper.GetBool("server.show_about"),
		SinceYear:       viper.GetString("server.since_year"),
		FriendLinkMap:   viper.GetStringMapString("friend_link"),
	}

	m.SetCanonicalURL(c)
	m.SetSysNotice(c)

	logging.Debugf(ctx, "NewMetaData MetaData:%+v", *m)
	return m
}

// SetKeywords 设置Keywords字段
func (m *MetaData) SetKeywords(c *gin.Context, keywords []string) {
	for _, kw := range keywords {
		m.Keywords = append(m.Keywords, webserver.CtxI18n(c, kw))
	}
}

// SetSysNotice 设置系统公告
func (m *MetaData) SetSysNotice(ctx context.Context) {
	notice := viper.GetString("sys_notice.md")
	if notice == "" {
		return
	}

	startTime, err := time.ParseInLocation(time.DateTime, viper.GetString("sys_notice.start"), TimeLocationBeijing)
	if err != nil {
		logging.Errorw(ctx, "SetSysNotice ParseInLocation error", "error", err)
		return
	}
	endTime, err := time.ParseInLocation(time.DateTime, viper.GetString("sys_notice.end"), TimeLocationBeijing)
	if err != nil {
		logging.Errorw(ctx, "SetSysNotice ParseInLocation error", "error", err)
		return
	}

	nowTs := time.Now().Unix()
	if nowTs >= startTime.Unix() && nowTs <= endTime.Unix() {
		m.SysNotice = notice
		m.SysNoticeQRText = viper.GetString("sys_notice.qrtext")
	}
}

// SetCanonicalURL 设置典范URL
func (m *MetaData) SetCanonicalURL(c *gin.Context) {
	m.HostURL = GetHostURL(c)
	m.StaticsURL, _ = url.JoinPath(m.HostURL, viper.GetString("statics.url"))
	baseURL, err := url.Parse(m.HostURL)
	if err != nil {
		return
	}
	reqPath, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		return
	}
	fullReqURL := baseURL.ResolveReference(reqPath).String()
	parsedURL, err := url.Parse(fullReqURL)
	if err != nil {
		return
	}
	queryParams := parsedURL.Query()
	// 要删除的查询参数列表
	paramsToRemove := []string{"from", "hmsr", "utm_source", "utm_medium", "alert"}
	for _, param := range paramsToRemove {
		queryParams.Del(param)
	}
	parsedURL.RawQuery = queryParams.Encode()
	m.CanonicalURL = parsedURL.String()
}

func (m *MetaData) CanonicalLinkTag() string {
	if m.CanonicalURL == "" {
		return ""
	}
	return fmt.Sprintf(`<link rel="canonical" href="%s">`, m.CanonicalURL)
}

func (m *MetaData) HreflangLinkTags() string {
	// x-default:删除CanonicalURL中的lang参数
	canonicalURL, err := url.Parse(m.CanonicalURL)
	if err != nil {
		logging.Error(nil, "HreflangLinkTags parse CanonicalURL error:"+err.Error())
		return ""
	}
	queryParams := canonicalURL.Query()
	queryParams.Del("lang")
	canonicalURL.RawQuery = queryParams.Encode()
	defaultHreflangLink := canonicalURL.String()
	linkTags := []string{
		fmt.Sprintf(`<link rel="alternate" hreflang="x-default" href="%s" />`, defaultHreflangLink),
	}
	for _, langTag := range webserver.I18nLangTags {
		hreflangURL, err := url.Parse(defaultHreflangLink)
		if err != nil {
			logging.Error(nil, "HreflangLinkTags parse defaultHreflangLink error:"+err.Error())
			continue
		}
		queryParams := hreflangURL.Query()
		queryParams.Add("lang", langTag.String())
		hreflangURL.RawQuery = queryParams.Encode()
		linkTag := fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="%s">`, langTag, hreflangURL.String())
		linkTags = append(linkTags, linkTag)
	}
	return strings.Join(linkTags, "\n")
}

func GetHostURL(c *gin.Context) string {
	hostURL := "https://" + c.Request.Host
	host := strings.Split(c.Request.Host, ":")
	if len(host) == 2 && host[1] != "443" {
		hostURL = "http://" + c.Request.Host
	}

	if hosturl := viper.GetString("server.host_url"); hosturl != "" && viper.GetString("env") == "prod" {
		hostURL = hosturl
	}
	return hostURL
}

type Waline struct {
	ServerURL         string
	Type              WalineType // 评论|留言
	Path              string     // 当前页面完整url path
	WithCommentCount  bool
	WithPageviewCount bool
	ReactionTitle     string
	Reaction          interface{} // 为文章增加表情互动功能，通过设置表情地址数组来自定义表情图片，最大支持 8 个表情。
	Reaction0         string
	Reaction1         string
	Reaction2         string
	Reaction3         string
	Reaction4         string
	Reaction5         string
}

type WalineType string

const (
	WalineTypeMessageBoard WalineType = "留言"
	WalineTypeComment      WalineType = "评论"
)

func NewWaline(c *gin.Context, wtype WalineType, withCommentCount, withPageviewCount bool) *Waline {
	serverURL := viper.GetString("waline.server_url")
	if serverURL == "" {
		return nil
	}

	fullPath, err := url.JoinPath(GetHostURL(c), c.Request.URL.Path)
	if err != nil {
		logging.Error(c, "NewWaline url JoinPath error:"+err.Error())
		return nil
	}

	return &Waline{
		ServerURL:         serverURL,
		Type:              wtype,
		Path:              fullPath,
		WithCommentCount:  withCommentCount,
		WithPageviewCount: withPageviewCount,
		ReactionTitle:     viper.GetString("waline.reaction_title"),
		Reaction:          viper.Get("waline.reaction"),
		Reaction0:         viper.GetString("waline.reaction0"),
		Reaction1:         viper.GetString("waline.reaction1"),
		Reaction2:         viper.GetString("waline.reaction2"),
		Reaction3:         viper.GetString("waline.reaction3"),
		Reaction4:         viper.GetString("waline.reaction4"),
		Reaction5:         viper.GetString("waline.reaction5"),
	}
}
