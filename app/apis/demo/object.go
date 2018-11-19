package demo

import (
	demoService "github.com/axiaoxin/pink-lady/app/services/demo"
	"github.com/axiaoxin/pink-lady/app/services/retcode"
	"github.com/axiaoxin/pink-lady/app/utils/response"

	"github.com/gin-gonic/gin"
)

// AddObjectBody use for bind json
type AddObjectBody struct {
	// APP ID
	AppID string `json:"appID" binding:"required,max=16" example:"1"`
	// system name
	System string `json:"system" binding:"required,max=64" example:"cmdb"`
	// entity name
	Entity string `json:"entity" binding:"required,max=64" example:"server"`
	// identity
	Identity string `json:"identity" binding:"required,max=64" example:"1"`
}

// ObjectQuery use for bind query
type ObjectQuery struct {
	// object ID
	ID uint `form:"id"`
	// APP ID
	AppID string `form:"appID"`
	// system name
	System string `form:"system"`
	// entity name
	Entity string `form:"entity"`
	// identity
	Identity string `form:"identity"`
	// page number
	PageNum int `form:"pageNum"`
	// page size
	PageSize int `form:"pageSize"`
	// order way
	Order string `form:"order"`
}

// AddObject godoc
// @Summary Add new object
// @Description add new object return object ID
// @Tags object
// @Accept json
// @Produce json
// @Param object body demo.AddObjectBody true "object info"
// @Router /demo/object [post]
// @Success 200 {object} response.Response
func AddObject(c *gin.Context) {
	body := AddObjectBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
	}
	objectID, err := demoService.AddObject(body.AppID, body.System, body.Entity, body.Identity)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, objectID)
}

// Object godoc
// @Summary Query object by params
// @Description query object by params
// @Tags object
// @Accept json
// @Produce json
// @Param id query uint false "query by object ID, other conditions do not enter into force."
// @Param appID query string false "filter result list by appid"
// @Param system query string false "filter result list by system"
// @Param entity query string false "filter result list by entity"
// @Param identity query string false "filter result list by identity"
// @Param pageNum query int false "page number" default(1)
// @Param pageSize query int false "page size" default(10)
// @Param order query string false "order field and way" default(id asc)
// @Router /demo/object [get]
// @Success 200 {object} response.Response
func Object(c *gin.Context) {
	query := ObjectQuery{
		PageNum:  1,
		PageSize: 10,
		Order:    "id asc",
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
	}
	result, err := demoService.QueryObject(query.ID, query.AppID, query.System, query.Entity, query.Identity, query.PageNum, query.PageSize, query.Order)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, result)
}
