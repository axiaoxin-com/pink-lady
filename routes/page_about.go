package routes

import (
	"net/http"

	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

func PageAbout(c *gin.Context) {
	meta := NewMetaData(c, webserver.CtxI18n(c, "关于网站"))
	data := gin.H{
		"meta": meta,
	}

	c.HTML(http.StatusOK, "about.html", data)
	return
}
