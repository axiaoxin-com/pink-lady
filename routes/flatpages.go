package routes

import (
	"net/http"
	"net/url"
	"path"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/statics"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Flatpages(app *gin.Engine) {
	navPath := viper.GetString("flatpages.nav_path")
	if navPath == "" {
		navPath = "fp"
	}
	fp := app.Group(navPath)
	fp.GET("/", func(c *gin.Context) {
		entries, err := statics.Files.ReadDir("flatpages")
		if err != nil {
			logging.Error(c, "unescape flatpage filename error:"+err.Error())
			c.Redirect(302, viper.GetString("server.host_url")+"/")
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
			c.Redirect(302, viper.GetString("server.host_url")+"/")
			return
		}
		filepath := path.Join("flatpages", filename)
		file, err := statics.Files.ReadFile(filepath)
		if err != nil {
			logging.Error(c, "read flatpages file error:"+err.Error())
			c.Redirect(302, viper.GetString("server.host_url")+"/")
			return
		}
		content := webserver.CtxI18n(c, string(file))

		meta := NewMetaData(c, filename)
		data := gin.H{
			"meta":    meta,
			"content": string(content),
		}
		c.HTML(http.StatusOK, "flatpages_single.html", data)
		return
	})
}
