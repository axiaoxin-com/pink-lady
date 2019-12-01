package apis

import (
	"pink-lady/app/database"
	"pink-lady/app/response"
	"pink-lady/app/services/demosvc"

	"github.com/gin-gonic/gin"
)

// CreateAlertPolicy godoc
// @Summary 新增告警策略
// @Description 新增告警策略返回对应ID
// @Tags demo
// @Accept json
// @Produce json
// @Param json body apis.CreateAlertPolicyParam true "参数包含BaseParam和demomod.AlertPolicy"
// @Router /demo/alert-policy [post]
// @Success 200 {object} response.Response
func CreateAlertPolicy(c *gin.Context) {
	p := &CreateAlertPolicyParam{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.ErrJSON400(c, response.RCInvalidParam, err)
		return
	}
	db := database.UTDB()
	// 使用外层appid和uin作为告警策略的字段
	p.AlertPolicy.AppID = p.AppID
	p.AlertPolicy.Uin = p.Uin
	result, err := demosvc.CreateAlertPolicy(c, db, p.AlertPolicy)
	if err != nil {
		response.ErrJSON(c, err)
		return
	}
	response.JSON(c, result)
}

// DescribeAlertPolicy godoc
// @Summary 查询指定告警策略
// @Description 根据appid、uin查询指定id的告警策略详情，包括其关联数据
// @Tags demo
// @Accept json
// @Produce json
// @Param appid path int true "param desc" default(1)
// @Param uin path string true "param desc" default(axiaoxin)
// @Param id path int true "param desc"
// @Router /demo/alert-policy/{appid}/{uin}/{id} [get]
// @Success 200 {object} response.Response
func DescribeAlertPolicy(c *gin.Context) {
	p := &DescribeAlertPolicyParam{}
	if err := c.ShouldBindUri(p); err != nil {
		response.ErrJSON400(c, response.RCInvalidParam, err)
		return
	}
	db := database.UTDB()
	result, err := demosvc.DescribeAlertPolicy(c, db, p.AppID, p.Uin, p.ID)
	if err != nil {
		response.ErrJSON(c, err)
		return
	}
	response.JSON(c, result)
}

// DeleteAlertPolicy godoc
// @Summary 删除告警策略
// @Description 根据appid、uin删除指定id的告警策略，包括其关联条件，成功Data为true
// @Tags demo
// @Accept json
// @Produce json
// @Param appid path int true "param desc" default(1)
// @Param uin path string true "param desc" default(axiaoxin)
// @Param id path int true "param desc"
// @Router /demo/alert-policy/{appid}/{uin}/{id} [delete]
// @Success 200 {object} response.Response
func DeleteAlertPolicy(c *gin.Context) {
	p := &DeleteAlertPolicyParam{}
	if err := c.ShouldBindUri(p); err != nil {
		response.ErrJSON400(c, response.RCInvalidParam, err)
		return
	}
	db := database.UTDB()
	err := demosvc.DeleteAlertPolicy(c, db, p.AppID, p.Uin, p.ID)
	if err != nil {
		response.ErrJSON(c, err)
		return
	}
	response.JSON(c, true)
}

// ModifyAlertPolicy godoc
// @Summary 更新告警策略
// @Description 更新告警策略返回对应ID
// @Tags demo
// @Accept json
// @Produce json
// @Param json body apis.ModifyAlertPolicyParam true "参数包含BaseParam和demomod.AlertPolicy"
// @Router /demo/alert-policy [put]
// @Success 200 {object} response.Response
func ModifyAlertPolicy(c *gin.Context) {
	p := &ModifyAlertPolicyParam{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.ErrJSON400(c, response.RCInvalidParam, err)
		return
	}
	db := database.UTDB()
	// 使用外层appid和uin作为告警策略的字段
	p.AlertPolicy.AppID = p.AppID
	p.AlertPolicy.Uin = p.Uin
	result, err := demosvc.ModifyAlertPolicy(c, db, p.AlertPolicy)
	if err != nil {
		response.ErrJSON(c, err)
		return
	}
	response.JSON(c, result)
}

// DescribeAlertPolicies godoc
// @Summary 查询告警策略列表
// @Description 根据appid、uin、搜索条件等查询告警策略列表，包括其关联数据
// @Tags demo
// @Accept json
// @Produce json
// @Param appid query int true "必填参数" default(1)
// @Param uin query string true "必填参数" default(axiaoxin)
// @Param offset query int false "分页offset" default(0)
// @Param limit query int false "分页limit" default(10)
// @Param order query string false "排序方式" default(id desc)
// @Param id query string false "按ID搜索"
// @Param name query string false "按名字模糊搜索"
// @Router /demo/alert-policy [get]
// @Success 200 {object} response.Response
func DescribeAlertPolicies(c *gin.Context) {
	p := &DescribeAlertPoliciesParam{
		Offset: 0,
		Limit:  10,
		Order:  "id desc",
	}
	if err := c.ShouldBindQuery(p); err != nil {
		response.ErrJSON400(c, response.RCInvalidParam, err)
		return
	}
	db := database.UTDB()
	result, count, err := demosvc.DescribeAlertPolicies(c, db, p.AppID, p.Uin, p.Offset, p.Limit, p.Order, p.ID, p.Name)
	if err != nil {
		response.ErrJSON(c, err)
		return
	}
	response.JSON(c, gin.H{
		"total_count":    count,
		"alert_policies": result,
	})
}
