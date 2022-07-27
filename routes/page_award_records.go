package routes

import (
	"net/http"

	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

func PageAwardRecords(c *gin.Context) {
	meta := NewMetaData(c, webserver.CtxI18n(c, "赞赏记录"))
	data := gin.H{
		"meta": meta,
	}

	c.HTML(http.StatusOK, "award_records.html", data)
	return
}
