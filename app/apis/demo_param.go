package apis

import (
	"pink-lady/app/models/demomod"
)

// BaseParam 公共参数
type BaseParam struct {
	// appid
	AppID int `json:"appid" form:"appid" uri:"appid" example:"1" binding:"required"`
	// 根帐号uin
	Uin string `json:"uin" form:"uin" uri:"uin" example:"axiaoxin" binding:"required"`
}

// CreateAlertPolicyParam 创建告警策略的参数
type CreateAlertPolicyParam struct {
	BaseParam
	// AlertPolicy不能使用匿名结构体，里面的AppID和Uin会和BaseParam里面的冲突，否则导致参数绑定失败
	AlertPolicy *demomod.AlertPolicy `json:"alert_policy" binding:"required"`
}

// DescribeAlertPolicyParam 获取告警策略详情的参数
type DescribeAlertPolicyParam struct {
	BaseParam
	ID int64 `json:"id" form:"id" uri:"id" example:"0" binding:"required"`
}

// DeleteAlertPolicyParam 删除告警策略的参数
type DeleteAlertPolicyParam struct {
	BaseParam
	ID int64 `json:"id" form:"id" uri:"id" example:"0" binding:"required"`
}

// ModifyAlertPolicyParam 更新告警策略的参数
type ModifyAlertPolicyParam struct {
	BaseParam
	// AlertPolicy不能使用匿名结构体，里面的AppID和Uin会和BaseParam里面的冲突，否则导致参数绑定失败
	AlertPolicy *demomod.AlertPolicy `json:"alert_policy" binding:"required"`
}

// DescribeAlertPoliciesParam 查询告警策略列表的参数
type DescribeAlertPoliciesParam struct {
	BaseParam
	Offset int    `json:"offset" form:"offset" example:"0" binding:"-"` // 偏移量，整型
	Limit  int    `json:"limit" form:"limit" example:"20" binding:"-"`  // 限制数目，整型
	Order  string `json:"order" form:"order" example:"id desc" binding:"-"`
	ID     int64  `json:"id" form:"id" example:"0" binding:"-"`
	Name   string `json:"name" form:"name" example:"" binding:"-"`
}
