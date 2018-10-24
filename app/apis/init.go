// package apis save the web api code,
// routes.go register handle function on url
//
// package apis init in init.go
// gin-skeleton system default apis in init.go
package apis

import (
	"github.com/axiaoxin/gin-skeleton/app/apis/retcode"
	"github.com/axiaoxin/gin-skeleton/app/common"
	"github.com/gin-gonic/gin"
)

// package init function
func init() {

}

// response current api version for ping request
func Ping(c *gin.Context) {
	data := gin.H{"version": common.VERSION}
	Respond(c, retcode.SUCCESS, data)
}
