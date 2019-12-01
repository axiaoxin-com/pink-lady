package demomod

import "pink-lady/app/models"

const (
	// AlertPolicyStatusEnable 告警策略状态开启
	AlertPolicyStatusEnable = 1
	// AlertPolicyStatusDisable 告警策略状态关闭
	AlertPolicyStatusDisable = 2
)

// AlertFilterRule 告警筛选规则表
type AlertFilterRule struct {
	models.BaseModel
	AlertPolicyID int64  `gorm:"column:alert_policy_id" json:"alert_policy_id" binding:"-" example:"-"` // 关联的告警策略ID
	Relation      int    `gorm:"column:relation" json:"relation" binding:"required" example:"1"`        // 与或关系（1=与，2=或）
	Field         string `gorm:"column:field" json:"field" binding:"required" example:"field"`          // 日志原始字段名
	Operating     string `gorm:"column:operating" json:"operating" binding:"required" example:"="`      // 操作符
	Value         string `gorm:"column:value" json:"value" binding:"required" example:"value"`          // 筛选值（多个值逗号分隔）
}

// TableName define tabel name
func (*AlertFilterRule) TableName() string {
	return "alert_filter_rule"
}

// AlertTriggerRule 告警策略触发条件表
type AlertTriggerRule struct {
	models.BaseModel
	AlertPolicyID        int64  `gorm:"column:alert_policy_id" json:"alert_policy_id" binding:"-" example:"-"`                      // 关联的告警策略主键id
	Relation             int    `gorm:"column:relation" json:"relation" binding:"required" example:"1"`                             // 与或关系（1=与，2=或）
	MetricID             int    `gorm:"column:metric_id" json:"metric_id" binding:"required" example:"-1"`                          // 指标id
	MetricType           int    `gorm:"column:metric_type" json:"metric_type" binding:"required" example:"1"`                       // 指标类型 1=普通指标 2=复合指标
	Operating            string `gorm:"column:operating" json:"operating" binding:"required" example:"="`                           // 操作符
	Value                string `gorm:"column:value" json:"value" binding:"required" example:"1"`                                   // 指标阈值
	ContinuousCycleCount int    `gorm:"column:continuous_cycle_count" json:"continuous_cycle_count" binding:"required" example:"1"` // 持续周期个数
}

// TableName define tabel name
func (*AlertTriggerRule) TableName() string {
	return "alert_trigger_rule"
}

// AlertPolicy 告警策略表
type AlertPolicy struct {
	models.BaseModel
	AppID              int                 `gorm:"column:appid" json:"appid" example:"1" binding:"-"`                                       // AppID
	Uin                string              `gorm:"column:uin" json:"uin" example:"axiaoxin" binding:"-"`                                    // Uin
	Name               string              `gorm:"column:name" json:"name" binding:"required" example:"swag-test-name"`                     // 策略名称
	MetricSetID        int64               `gorm:"column:metric_set_id" json:"metric_set_id" binding:"required" example:"-1"`               // 指标集ID
	NoticeFrequencySec int                 `gorm:"column:notice_frequency_sec" json:"notice_frequency_sec" binding:"required" example:"60"` // 通知频率（通知间隔秒数）
	Status             int                 `gorm:"column:status" json:"status" binding:"required" example:"1"`                              // 状态 1=已开启 2=未开启 3=已失效
	AlertGroupID       string              `gorm:"column:alert_group_id" json:"alert_group_id" binding:"required" example:"-1,0"`           // 告警接收组 逗号分隔
	AlertChannel       string              `gorm:"column:alert_channel" json:"alert_channel" binding:"required" example:"weixin,sms"`       // 告警接收渠道 1=邮件 2=短信 3=微信
	NoticePeriodBegin  int                 `gorm:"column:notice_period_begin" json:"notice_period_begin" binding:"-" example:"0"`           // 通知时段开始时间（从00:00:00开始计算的秒数）
	NoticePeriodEnd    int                 `gorm:"column:notice_period_end" json:"notice_period_end" binding:"required" example:"86399"`    // 通知时段结束时间（从00:00:00开始计算的秒数）
	URLScheme          string              `gorm:"column:url_scheme" json:"url_scheme" binding:"required" example:"http"`                   // 回调url的scheme
	CallbackURL        string              `gorm:"column:callback_url" json:"callback_url" binding:"required" example:"axiaoxin.com"`       // 回调url 不包含scheme部分
	LatestAlertTime    string              `gorm:"column:latest_alert_time" json:"latest_alert_time" binding:"-" example:""`                // 最后告警时间（产生告警后更改该字段）
	AlertFilterRules   []*AlertFilterRule  `json:"alert_filter_rules" binding:"dive"`                                                       // 告警过滤条件
	AlertTriggerRules  []*AlertTriggerRule `json:"alert_trigger_rules" binding:"required,dive"`                                             // 告警触发条件
}

// TableName define tabel name
func (*AlertPolicy) TableName() string {
	return "alert_policy"
}
