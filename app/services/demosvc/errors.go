// 业务自身的错误码 (云API子错误码)

package demosvc

import "pink-lady/app/response"

var (
	ErrAlertPolicyCreateFailed        = response.RCFailure.SetMsg("创建告警策略失败")
	ErrAlertPolicyDeleteFailed        = response.RCFailure.SetMsg("删除告警策略失败")
	ErrAlertPolicyDescribeFailed      = response.RCFailure.SetMsg("查询告警策略失败")
	ErrAlertFilterRuleDescribeFailed  = response.RCFailure.SetMsg("查询告警过滤规则失败")
	ErrAlertTriggerRuleDescribeFailed = response.RCFailure.SetMsg("查询告警触发规则失败")
	ErrAlertPolicyModifyFailed        = response.RCFailure.SetMsg("修改告警策略失败")
	ErrAlertPolicyListFailed          = response.RCFailure.SetMsg("查询告警策略列表失败")
	ErrAlertFilterRuleDeleteFailed    = response.RCFailure.SetMsg("删除告警过滤条件失败")
	ErrAlertTriggerRuleDeleteFailed   = response.RCFailure.SetMsg("删除告警触发条件失败")
)
