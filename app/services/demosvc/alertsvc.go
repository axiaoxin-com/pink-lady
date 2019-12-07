package demosvc

import (
	"pink-lady/app/database"
	"pink-lady/app/logging"
	"pink-lady/app/models"
	"pink-lady/app/models/demomod"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CreateAlertPolicy 新增告警策略
// 主体信息插入告警策略表，AlertTriggerRules 插入告警触发条件，AlertFilterRules 插入告警过滤条件
// 单个create是原子操作不需要使用事务
func CreateAlertPolicy(c *gin.Context, db *gorm.DB, policy *demomod.AlertPolicy) (int64, error) {
	// 创建告警策略
	if err := db.Create(policy).Error; err != nil {
		logging.CtxLogger(c).Error(err.Error())
		return 0, ErrAlertPolicyCreateFailed.AppendError(err)
	}
	return policy.ID, nil
}

// DeleteAlertPolicy 删除告警策略
// 根据appid和uin对指定id的数据进行删除，返回被删除对象和错误信息
// 关联的告警过滤条件和触发条件也会被删除
// 使用Association Clear不会删除记录，只是把关联的策略ID设置为空
// 这里采用真实删除，涉及多个delete操作需要使用事务
func DeleteAlertPolicy(c *gin.Context, db *gorm.DB, appID int, uin string, id int64) error {
	// 开启事务
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logging.CtxLogger(c).Error("DeleteAlertPolicy Begin tx error", zap.Error(err))
		return ErrAlertPolicyDeleteFailed.AppendError(err)
	}
	// 始终关闭事务
	defer tx.Rollback()

	// 查询要删除的记录,不存在返回错误
	policy := &demomod.AlertPolicy{}
	if err := db.Where("appid = ? AND uin = ? AND id = ?", appID, uin, id).Find(policy).Error; err != nil {
		logging.CtxLogger(c).Error("DeleteAlertPolicy Find policy error", zap.Error(err))
		return ErrAlertPolicyDeleteFailed.AppendError(err)
	}
	// 删除该记录
	if err := db.Delete(policy).Error; err != nil {
		logging.CtxLogger(c).Error("DeleteAlertPolicy Delete policy error", zap.Error(err))
		return ErrAlertPolicyDeleteFailed.AppendError(err)
	}

	// 删除告警过滤条件
	if err := tx.Where("alert_policy_id = ?", id).Delete(demomod.AlertFilterRule{}).Error; err != nil {
		logging.CtxLogger(c).Error("DeleteAlertPolicy Delete filter rule error", zap.Error(err))
		return ErrAlertPolicyDeleteFailed.AppendError(err)
	}

	// 删除告警触发条件
	if err := tx.Where("alert_policy_id = ?", id).Delete(demomod.AlertTriggerRule{}).Error; err != nil {
		logging.CtxLogger(c).Error("DeleteAlertPolicy Delete trigger rule error", zap.Error(err))
		return ErrAlertPolicyDeleteFailed.AppendError(err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logging.CtxLogger(c).Error("DeleteAlertPolicy Commit error", zap.Error(err))
		return ErrAlertPolicyDeleteFailed.AppendError(err)
	}
	return nil
}

// DescribeAlertPolicy 查询告警策略详情
// 及其关联的触发条件和筛选条件
func DescribeAlertPolicy(c *gin.Context, db *gorm.DB, appID int, uin string, id int64) (*demomod.AlertPolicy, error) {
	// 按appid uin id获取记录
	result := db.Where("appid = ? AND uin = ? AND id = ?", appID, uin, id)
	if result.Error != nil {
		logging.CtxLogger(c).Error("AlertPolicy Where error", zap.Error(result.Error))
		return nil, ErrAlertPolicyDescribeFailed.AppendError(result.Error)
	}
	// 加载关联的筛选条件
	result = result.Preload("AlertFilterRules", "alert_policy_id = ?", id)
	if result.Error != nil {
		logging.CtxLogger(c).Error("AlertPolicy Preload AlertFilterRules error", zap.Error(result.Error))
		return nil, ErrAlertPolicyDescribeFailed.AppendError(result.Error)
	}
	// 加载关联的触发条件
	result = result.Preload("AlertTriggerRules", "alert_policy_id = ?", id)
	if result.Error != nil {
		logging.CtxLogger(c).Error("AlertPolicy Preload AlertTriggerRules error", zap.Error(result.Error))
		return nil, ErrAlertPolicyDescribeFailed.AppendError(result.Error)
	}
	policy := &demomod.AlertPolicy{}
	if err := result.Find(policy).Error; err != nil {
		logging.CtxLogger(c).Error("AlertPolicy Find error", zap.Error(err))
		return nil, ErrAlertPolicyDescribeFailed.AppendError(err)
	}
	return policy, nil
}

// DescribeAlertFilterRules 获取告警策略筛选条件
func DescribeAlertFilterRules(c *gin.Context, db *gorm.DB, appID int, uin string, alertPolicyID int64) ([]demomod.AlertFilterRule, error) {
	ruleSlice := []demomod.AlertFilterRule{}
	policy := &demomod.AlertPolicy{
		BaseModel: models.BaseModel{
			ID: alertPolicyID,
		},
		AppID: appID,
		Uin:   uin,
	}
	if err := db.Model(policy).Related(&ruleSlice).Error; err != nil {
		logging.CtxLogger(c).Error("AlertFilterRules Related error", zap.Error(err))
		return nil, ErrAlertFilterRuleDescribeFailed.AppendError(err)
	}
	return ruleSlice, nil
}

// DescribeAlertTriggerRules 获取告警策略触发条件
func DescribeAlertTriggerRules(c *gin.Context, db *gorm.DB, policyID int64) ([]demomod.AlertTriggerRule, error) {
	ruleSlice := []demomod.AlertTriggerRule{}
	policy := &demomod.AlertPolicy{
		BaseModel: models.BaseModel{
			ID: policyID,
		},
	}
	if err := db.Model(policy).Related(&ruleSlice).Error; err != nil {
		logging.CtxLogger(c).Error("AlertTriggerRules Related error", zap.Error(err))
		return nil, ErrAlertTriggerRuleDescribeFailed.AppendError(err)
	}
	return ruleSlice, nil
}

// ModifyAlertPolicy 更新告警策略
// 使用Association Replace来更新关联条件如果传了主键ID则更新对应记录，没有主键ID的全部新增记录，剔除了关联并不会删除记录只会把关联的id置为0
// 这里需要真实删除记录 使用事务删除全部关联
func ModifyAlertPolicy(c *gin.Context, db *gorm.DB, policy *demomod.AlertPolicy) (*demomod.AlertPolicy, error) {
	// 查询被更新的记录是否存在
	rawPolicy := &demomod.AlertPolicy{}
	if err := db.Where("appid = ? AND uin = ? AND id = ?", policy.AppID, policy.Uin, policy.ID).Find(rawPolicy).Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Find error", zap.Error(err))
		return nil, ErrAlertPolicyModifyFailed.AppendError(err)
	}

	// 开启事务
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Begin tx error", zap.Error(err))
		return nil, ErrAlertPolicyModifyFailed.AppendError(err)
	}
	defer tx.Rollback()

	// 使用结构体更新只会更新非零值，这里使用map方式，只会更新其中有变化的属性
	if err := tx.Model(rawPolicy).Updates(map[string]interface{}{
		"name":                 policy.Name,
		"metric_id":            policy.MetricSetID,
		"notice_frequency_sec": policy.NoticeFrequencySec,
		"alert_group_id":       policy.AlertGroupID,
		"alert_channel":        policy.AlertChannel,
		"notice_period_begin":  policy.NoticePeriodBegin,
		"notice_period_end":    policy.NoticePeriodEnd,
		"url_scheme":           policy.URLScheme,
		"callback_url":         policy.CallbackURL,
	}).Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Updates policy error", zap.Error(err))
		return nil, ErrAlertPolicyModifyFailed.AppendError(err)
	}

	// 删除过滤条件
	if err := tx.Where("alert_policy_id = ?", policy.ID).Delete(demomod.AlertFilterRule{}).Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Delete filter rules error", zap.Error(err))
		return nil, ErrAlertFilterRuleDeleteFailed.AppendError(err)
	}
	// 删除触发条件
	if err := tx.Where("alert_policy_id = ?", policy.ID).Delete(demomod.AlertTriggerRule{}).Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Delete trigger rules error", zap.Error(err))
		return nil, ErrAlertTriggerRuleDeleteFailed.AppendError(err)
	}

	// 更新关联数据
	tx.Model(rawPolicy).Association("AlertFilterRules").Replace(policy.AlertFilterRules)
	tx.Model(rawPolicy).Association("AlertTriggerRules").Replace(policy.AlertTriggerRules)

	// 保存更新
	if err := tx.Save(rawPolicy).Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Save rules error", zap.Error(err))
		return nil, ErrAlertPolicyModifyFailed.AppendError(err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logging.CtxLogger(c).Error("ModifyAlertPolicy Commit error", zap.Error(err))
		return nil, ErrAlertPolicyModifyFailed.AppendError(err)
	}
	return rawPolicy, nil
}

// DescribeAlertPolicies 查询告警策略列表
// appid, uin必须要的过滤条件
// offset 指定开始返回记录前要跳过的记录数,-1表示取消offset条件
// limit 指定检索的最大记录数，-1表示取消limit条件
// order 字段名+desc/asc指定排序字段和方式: id desc
// id 按id搜索
// name 按名字模糊搜索
// 返回列表数据、最终数据总数、error
func DescribeAlertPolicies(c *gin.Context, db *gorm.DB, appID int, uin string, offset, limit int, order string, ID int64, name string) ([]*demomod.AlertPolicy, int, error) {
	// 策略列表
	policies := []*demomod.AlertPolicy{}
	// 记录总数
	totalCount := 0

	search := db.Where("appid = ? AND uin = ?", appID, uin)
	if search.Error != nil {
		logging.CtxLogger(c).Error("DescribeAlertPolicies base search error", zap.Error(search.Error))
		return nil, 0, ErrAlertPolicyListFailed.AppendError(search.Error)
	}
	if ID != 0 {
		// 按ID搜索单条记录
		search = search.Where("id = ?", ID)
		if search.Error != nil {
			logging.CtxLogger(c).Error("DescribeAlertPolicies search by id error", zap.Error(search.Error))
			return nil, 0, ErrAlertPolicyListFailed.AppendError(search.Error)
		}
	} else {
		// 按名称搜索多条记录
		if name != "" {
			search = search.Where("name LIKE ?", "%"+database.GormMySQLLikeFieldEscape(name)+"%")
			if search.Error != nil {
				logging.CtxLogger(c).Error("DescribeAlertPolicies search by name error", zap.Error(search.Error))
				return nil, 0, ErrAlertPolicyListFailed.AppendError(search.Error)
			}
		}
	}

	// 获取记录总数
	if err := search.Model(&demomod.AlertPolicy{}).Count(&totalCount).Error; err != nil {
		msg := "DescribeAlertPolicies count error"
		logging.CtxLogger(c).Error(msg, zap.Error(err))
		return nil, 0, errors.Wrap(err, msg)
	}

	// 发起最终查询
	if err := search.Preload("AlertFilterRules").Preload("AlertTriggerRules").Limit(limit).Offset(offset).Order(order).Find(&policies).Error; err != nil {
		msg := "DescribeAlertPolicies find error"
		logging.CtxLogger(c).Error(msg, zap.Error(err))
		return nil, 0, errors.Wrap(err, msg)
	}
	return policies, totalCount, nil
}
