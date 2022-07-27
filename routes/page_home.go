package routes

import (
	"net/http"

	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

func PageHome(c *gin.Context) {
	meta := NewMetaData(c, webserver.CtxI18n(c, "首页"))

	data := gin.H{
		"meta":  meta,
		"alert": Alert(c, "", ""),
	}

	c.HTML(http.StatusOK, "home.html", data)
	return
}
