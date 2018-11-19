package demo

import (
	demoService "github.com/axiaoxin/pink-lady/app/services/demo"
	"github.com/axiaoxin/pink-lady/app/services/retcode"
	"github.com/axiaoxin/pink-lady/app/utils/response"

	"github.com/gin-gonic/gin"
)

// AddLabelBody use for bind json
type AddLabelBody struct {
	// label name
	Name string `json:"name" binding:"required,max=32" example:"label name"`
	// label remark
	Remark string `json:"remark" binding:"max=64" example:"label remark"`
}

// LabelQuery use for bind query
type LabelQuery struct {
	// label ID
	ID uint `form:"id"`
	// label name
	Name string `form:"name"`
	// label remark
	Remark string `form:"remark"`
	// page number
	PageNum int `form:"pageNum"`
	// page size
	PageSize int `form:"pageSize"`
	// how to order
	Order string `form:"order"`
}

// AddLabel godoc
// @Summary Add new label
// @Description give name and remark to add new label, return label ID, if label exist, update then remark field
// @Tags label
// @Accept json
// @Produce json
// @Param label body demo.AddLabelBody true "label info"
// @Router /demo/label [post]
// @Success 200 {object} response.Response
func AddLabel(c *gin.Context) {
	body := AddLabelBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
	}
	labelID, err := demoService.AddLabel(body.Name, body.Remark)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, labelID)
}

// Label godoc
// @Summary Query label by query params
// @Description query label by id or name, or paging query
// @Tags label
// @Accept json
// @Produce json
// @Param id query uint false "query by ID, other conditions do not enter into force."
// @Param name query string false "query by name, other conditions do not enter into force."
// @Param remark query string false "fuzzy query by remark"
// @Param pageNum query int false "page number" default(1)
// @Param pageSize query int false "page size" default(10)
// @Param order query string false "order field and order type" default(id asc)
// @Router /demo/label [get]
// @Success 200 {object} response.Response
func Label(c *gin.Context) {
	query := LabelQuery{
		PageNum:  1,
		PageSize: 10,
		Order:    "id asc",
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
	}

	result, err := demoService.QueryLabel(query.ID, query.Name, query.Remark, query.PageNum, query.PageSize, query.Order)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, result)
}
