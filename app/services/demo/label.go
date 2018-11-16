package demo

import (
	"strings"

	demoModels "pink-lady/app/models/demo"
	"pink-lady/app/utils"
)

// AddLabel 新建标签
// 返回标签ID
// 如果标签名称已经存在，不会报错，直接返回已存在标签的ID
// 如果标签名称已存在，且传递了remark参数则会更新remark字段
func AddLabel(name, remark string) (uint, error) {
	label := demoModels.Label{}
	err := utils.DB.Where(demoModels.Label{Name: name}).Assign(demoModels.Label{Remark: remark}).FirstOrCreate(&label).Error
	return label.ID, err
}

// GetLabelByName 按标签名称查询标签
func GetLabelByName(name string) (demoModels.Label, error) {
	label := demoModels.Label{}
	err := utils.DB.Where(&demoModels.Label{Name: name}).First(&label).Error
	return label, err
}

// GetLabelByID 按标签ID查询标签
func GetLabelByID(labelID uint) (demoModels.Label, error) {
	label := demoModels.Label{}
	err := utils.DB.First(&label, labelID).Error
	return label, err
}

// GetLabelsByIDs 按标签ID列表批量查询标签
func GetLabelsByIDs(labelIDs []uint) ([]demoModels.Label, error) {
	labels := []demoModels.Label{}
	err := utils.DB.Where(labelIDs).Find(&labels).Error
	return labels, err
}

// QueryLabel 根据参数组合查询label
// 传labelID则只返回对应ID的label，其他参数不生效
// offset为分页页码默认为1
// limit为分页查询每页数量, -1可取消限制
// order为排序方式 "字段 desc/asc" 默认按id降序排列
// 多条记录的结果返回分页信息
func QueryLabel(labelID uint, name, remark string, pageNum, pageSize int, order string) (interface{}, error) {
	if labelID > 0 {
		return GetLabelByID(labelID)
	}
	if name != "" {
		return GetLabelByName(name)
	}

	if order == "" {
		order = "id asc"
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	offset := (pageNum - 1) * pageSize
	limit := pageSize
	if pageSize <= 0 {
		pageSize = -1
		offset = -1
		limit = -1
	}

	count := 0
	totalCount := true
	scopedb := utils.DB.Order(order).Offset(offset).Limit(limit)
	if remark != "" {
		scopedb = scopedb.Where("remark LIKE ?", "%"+strings.TrimSpace(remark)+"%")
		totalCount = false
	}
	items := []demoModels.Label{}
	scopedb = scopedb.Find(&items)
	if totalCount {
		utils.DB.Model(&demoModels.Label{}).Count(&count)
	} else {
		scopedb.Count(&count)
	}
	result := map[string]interface{}{
		"items":      items,
		"pagination": utils.Paginate(count, pageNum, pageSize),
	}
	return result, scopedb.Error
}
