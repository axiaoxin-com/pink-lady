package demo

import (
	demoModels "github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/gin-gonic/gin"
)

// AddLabeling add associations for the objectIDs and LabelIDs
func AddLabeling(c *gin.Context, objectIDs, labelIDs []uint) ([]map[string]interface{}, error) {
	objects, err := GetObjectsByIDs(c, objectIDs)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	labels, err := GetLabelsByIDs(c, labelIDs)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	var results []map[string]interface{}
	for _, object := range objects {
		// 记录结果，不使用批量操作
		for _, label := range labels {
			result := map[string]interface{}{
				"objectID": object.ID,
				"labelID":  label.ID,
				"result":   "ok",
			}
			err := utils.DB.Model(&object).Association("Labels").Append(label).Error
			if err != nil {
				utils.CtxLogger(c).Warn(err)
				result["result"] = err
			}
			results = append(results, result)
		}
	}
	return results, err
}

// ReplaceLabeling replace old associations with new given objectIDs and LabelIDs
func ReplaceLabeling(c *gin.Context, objectIDs, labelIDs []uint) ([]map[string]interface{}, error) {
	objects, err := GetObjectsByIDs(c, objectIDs)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	labels, err := GetLabelsByIDs(c, labelIDs)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	var results []map[string]interface{}
	for _, object := range objects {
		result := map[string]interface{}{
			"objectID": object.ID,
			"result":   "ok",
		}
		err := utils.DB.Model(&object).Association("Labels").Replace(labels).Error
		if err != nil {
			utils.CtxLogger(c).Warn(err)
			result["result"] = err
		}
		results = append(results, result)
	}
	return results, err
}

// DeleteLabeling delete associations for object ids and label ids
func DeleteLabeling(c *gin.Context, objectIDs, labelIDs []uint) ([]map[string]interface{}, error) {
	objects, err := GetObjectsByIDs(c, objectIDs)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	labels, err := GetLabelsByIDs(c, labelIDs)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	var results []map[string]interface{}
	for _, object := range objects {
		// 记录结果，不使用批量操作
		for _, label := range labels {
			result := map[string]interface{}{
				"objectID": object.ID,
				"labelID":  label.ID,
				"result":   "ok",
			}
			err := utils.DB.Model(&object).Association("Labels").Delete(label).Error
			if err != nil {
				utils.CtxLogger(c).Warn(err)
				result["result"] = err
			}
			results = append(results, result)
		}
	}
	return results, err
}

// GetLabelingByLabelID 根据标签ID查询已关联的对象列表
func GetLabelingByLabelID(c *gin.Context, labelID uint) ([]demoModels.Object, error) {
	objects := []demoModels.Object{}
	label, err := GetLabelByID(c, labelID)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	scopedb := utils.DB.Model(&label).Association("Objects")
	scopedb.Find(&objects)
	err = scopedb.Error
	return objects, err
}

//GetLabelingByObjectID 根据对象ID查询已关联的标签列表
func GetLabelingByObjectID(c *gin.Context, objectID uint) ([]demoModels.Label, error) {
	labels := []demoModels.Label{}
	object, err := GetObjectByID(c, objectID)
	if err != nil {
		utils.CtxLogger(c).Warn(err)
		return nil, err
	}
	scopedb := utils.DB.Model(&object).Association("Labels")
	scopedb.Find(&labels)
	err = scopedb.Error
	return labels, err
}
