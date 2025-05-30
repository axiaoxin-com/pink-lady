package routes

import (
	"github.com/axiaoxin-com/pink-lady/webserver"
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
	// AlertErrorSubmit 表单提交失败
	AlertErrorSubmit = "submit-error"
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
			Heading: webserver.CtxI18n(c, "温馨提示！"),
			Text:    text,
		}

	case AlertWarningCommon:
		return &AlertData{
			Color:   "warning",
			Heading: webserver.CtxI18n(c, "操作失败！"),
			Text:    text,
		}

	case AlertErrorCommon:
		return &AlertData{
			Color:   "danger",
			Heading: webserver.CtxI18n(c, "系统错误，请稍后重试！"),
			Text:    text,
		}
	case AlertErrorSubmit:
		return &AlertData{
			Color:   "danger",
			Heading: webserver.CtxI18n(c, "系统错误！"),
			Text:    webserver.CtxI18n(c, "提交失败，请稍后再试！") + text,
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
