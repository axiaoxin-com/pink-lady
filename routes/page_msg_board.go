package routes

import (
	"net/http"

	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

func PageMsgBoard(c *gin.Context) {
	meta := NewMetaData(c, webserver.CtxI18n(c, "留言板"))
	data := gin.H{
		"meta": meta,
	}

	c.HTML(http.StatusOK, "msg_board.html", data)
	return
}
