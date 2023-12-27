package routes

import (
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/statics"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Flatpages(app *gin.Engine) {
	navPath := viper.GetString("flatpages.nav_path")
	if navPath == "" {
		navPath = "/fp"
	}
	fp := app.Group(navPath)
	fp.GET("/", func(c *gin.Context) {
		entries, err := statics.Files.ReadDir("flatpages")
		if err != nil {
			logging.Error(c, "unescape flatpage filename error:"+err.Error())
			c.Redirect(302, "/")
			return
		}

		meta := NewMetaData(c, viper.GetString("flatpages.nav_name"))
		data := gin.H{
			"meta":    meta,
			"entries": entries,
		}
		c.HTML(http.StatusOK, "flatpages_index.html", data)
		return
	})
	fp.GET("/:filename/", func(c *gin.Context) {
		filename, err := url.PathUnescape(c.Param("filename"))
		if err != nil {
			logging.Error(c, "unescape flatpage filename error:"+err.Error())
			c.Redirect(302, "/")
			return
		}
		filepath := path.Join("flatpages", filename)
		file, err := statics.Files.ReadFile(filepath)
		if err != nil {
			logging.Error(c, "read flatpages file error:"+err.Error())
			c.Redirect(302, "/")
			return
		}
		content := webserver.CtxI18n(c, string(file))
		title := ""
		if lines := strings.Split(content, "\n"); len(lines) > 0 {
			title = strings.TrimPrefix(lines[0], "#")
		}

		meta := NewMetaData(c, title)
		data := gin.H{
			"meta":    meta,
			"content": string(content),
		}
		c.HTML(http.StatusOK, "flatpages_single.html", data)
		return
	})
}
