package demo

import (
	"strconv"

	demoService "github.com/axiaoxin/pink-lady/app/services/demo"
	"github.com/axiaoxin/pink-lady/app/services/retcode"
	"github.com/axiaoxin/pink-lady/app/utils/response"

	"github.com/gin-gonic/gin"
)

// AddLabelingBody bind for validator
type AddLabelingBody struct {
	// which object ids need to be labling with the label ids
	ObjectIDs []uint `json:"objectIDs" binding:"required"`
	// which label ids need to be labling with the object ids
	LabelIDs []uint `json:"labelIDs" binding:"required"`
}

// AddLabeling godoc
// @Summary Labeling for object
// @Description create association for a given object ID list and tag ID list.
// @Tags labeling
// @Accept json
// @Produce json
// @Param labeling body demo.AddLabelingBody true "labeling association info"
// @Router /demo/labeling [post]
// @Success 200 {object} response.Response
func AddLabeling(c *gin.Context) {
	body := AddLabelingBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
		return
	}
	results, err := demoService.AddLabeling(body.ObjectIDs, body.LabelIDs)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, results)
}

// ReplaceLabeling godoc
// @Summary Update labeling by replace way
// @Description replace association for a given object ID list and tag ID list.
// @Tags labeling
// @Accept json
// @Produce json
// @Param labeling body demo.AddLabelingBody true "labeling association info"
// @Router /demo/labeling [put]
// @Success 200 {object} response.Response
func ReplaceLabeling(c *gin.Context) {
	body := AddLabelingBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
		return
	}
	results, err := demoService.ReplaceLabeling(body.ObjectIDs, body.LabelIDs)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, results)
}

// DeleteLabeling godoc
// @Summary Delete labeling
// @Description delete association for a given object ID list and tag ID list.
// @Tags labeling
// @Accept json
// @Produce json
// @Param labeling body demo.AddLabelingBody true "labeling association info"
// @Router /demo/labeling [delete]
// @Success 200 {object} response.Response
func DeleteLabeling(c *gin.Context) {
	body := AddLabelingBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
		return
	}
	results, err := demoService.DeleteLabeling(body.ObjectIDs, body.LabelIDs)
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, results)
}

// GetLabelingByLabelID godoc
// @Summary Query labeling object list by label ID
// @Description return labeling object list
// @Tags labeling
// @Accept json
// @Produce json
// @Param id path uint true "label id"
// @Router /demo/labeling/label/{id} [get]
// @Success 200 {object} response.Response
func GetLabelingByLabelID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
		return
	}
	result, err := demoService.GetLabelingByLabelID(uint(id))
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, result)
}

// GetLabelingByObjectID godoc
// @Summary Query labeling label list by object ID
// @Description return labeling label list
// @Tags labeling
// @Accept json
// @Produce json
// @Param id path uint true "object id"
// @Router /demo/labeling/object/{id} [get]
// @Success 200 {object} response.Response
func GetLabelingByObjectID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.JSON400(c, retcode.InvalidParams, err.Error())
		return
	}
	result, err := demoService.GetLabelingByObjectID(uint(id))
	if err != nil {
		response.JSON(c, retcode.Failure, err.Error())
		return
	}
	response.JSON(c, retcode.Success, result)
}
