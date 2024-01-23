package routes

import (
	"github.com/gin-gonic/gin"
)

// AlertData HTML提示框数据
type AlertData struct {
	Color   string `json:"color"`   // error, warning, info
	Heading string `json:"heading"` // 标题
	Text    string `json:"text"`    // 内容
}

const (
	// AlertWarningCommon 系统提示通用提示
	AlertInfoCommon = "sys-info"
	// AlertWarningCommon 系统警告通用提示
	AlertWarningCommon = "sys-warning"
	// AlertErrorCommon 系统错误通用提示
	AlertErrorCommon = "sys-error"
	// AlertOK OK
	AlertOK = "ok"
)

// Alert 页面提示
func Alert(c *gin.Context, alert, text string) *AlertData {
	if alert == "" {
		alert = c.Query("alert")
	}
	switch alert {
	case AlertInfoCommon:
		return &AlertData{
			Color:   "info",
			Heading: "温馨提示！",
			Text:    text,
		}
	case AlertWarningCommon:
		return &AlertData{
			Color:   "warning",
			Heading: "操作失败！",
			Text:    text,
		}
	case AlertErrorCommon:
		return &AlertData{
			Color:   "danger",
			Heading: "系统错误，请稍后重试！",
			Text:    text,
		}
	case AlertOK:
		return &AlertData{
			Color:   "success",
			Heading: "Success",
			Text:    text,
		}
	}
	return nil
}
