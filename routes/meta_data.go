package routes

import (
	"context"
	"html/template"
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
	SiteName         string
	Slogan           string
	HostURL          string
	BuildID          string
	Env              string
	AppID            int
	AdsenseID        string
	SysNotice        string
	SysNoticeQRText  string
	Title            string
	IsCrawler        bool
	Lang             string
	Keywords         []string
	BaseDesc         string
	BootswatchTheme  string
	I18n             bool
	Beian            string
	AuthorName       string
	AuthorURL        string
	StaticsURL       string
	StaticsSelfhost  bool
	FlatpagesEnable  bool
	FlatpagesNavName string
	FlatpagesNavPath string
	BaiduTongJiID    string
	GtagID           string
	CanonicalURL     template.HTML
	CanonicalLinkTag template.HTML
	ShowAbout        bool
	SinceYear        string
	FriendLinkMap    map[string]string
}

// NewMetaData 返回页面元数据
func NewMetaData(c *gin.Context, title string) (m *MetaData) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	hostURL := GetHostURL(c)
	canonicalURL := hostURL + c.Request.RequestURI

	ua := c.GetHeader("User-Agent")
	isCrawler := crawlerdetect.IsCrawler(ua)

	m = &MetaData{
		SiteName:         webserver.CtxI18n(c, SiteName),
		Slogan:           webserver.CtxI18n(c, Slogan),
		HostURL:          hostURL,
		BuildID:          BuildID,
		Env:              viper.GetString("env"),
		AppID:            AppID,
		AdsenseID:        viper.GetString("server.adsense_id"),
		Title:            title,
		IsCrawler:        isCrawler,
		Lang:             c.GetString("lang"),
		BaseDesc:         webserver.CtxI18n(c, "pink-lady是一个golang gin的web开发模板"),
		BootswatchTheme:  "cosmo",
		I18n:             viper.GetBool("i18n.enable"),
		Beian:            viper.GetString("server.beian"),
		AuthorName:       viper.GetString("author.name"),
		AuthorURL:        viper.GetString("author.url"),
		StaticsURL:       hostURL + "/" + viper.GetString("statics.url"),
		StaticsSelfhost:  viper.GetBool("statics.selfhost"),
		FlatpagesEnable:  viper.GetBool("flatpages.enable"),
		FlatpagesNavName: webserver.CtxI18n(c, viper.GetString("flatpages.nav_name")),
		FlatpagesNavPath: viper.GetString("flatpages.nav_path"),
		BaiduTongJiID:    viper.GetString("server.baidu_tongji_id"),
		GtagID:           viper.GetString("server.gtag_id"),
		CanonicalURL:     template.HTML(canonicalURL),
		CanonicalLinkTag: template.HTML(`<link rel="canonical" href="` + canonicalURL + `">`),
		ShowAbout:        viper.GetBool("server.show_about"),
		SinceYear:        viper.GetString("server.since_year"),
		FriendLinkMap:    viper.GetStringMapString("friend_link"),
	}

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
