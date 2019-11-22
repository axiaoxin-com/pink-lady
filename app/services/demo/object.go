package demo

import (
	"github.com/axiaoxin/pink-lady/app/db"
	demoModels "github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/gin-gonic/gin"
)

// AddObject insert an object to database
// return the object id
// if the object's fields value has existed, it will return the existed id
func AddObject(c *gin.Context, appID, system, entity, identity string) (uint, error) {
	obj := demoModels.Object{}
	err := db.SQLite3("testing").FirstOrCreate(&obj, demoModels.Object{AppID: appID, System: system, Entity: entity, Identity: identity}).Error
	return obj.ID, err
}

// GetObjectByID return a object model by id and error
func GetObjectByID(c *gin.Context, objectID uint) (demoModels.Object, error) {
	object := demoModels.Object{}
	err := db.SQLite3("testing").First(&object, objectID).Error
	return object, err
}

// GetObjectsByIDs return objects model list by ids and error
func GetObjectsByIDs(c *gin.Context, objectIDs []uint) ([]demoModels.Object, error) {
	objects := []demoModels.Object{}
	err := db.SQLite3("testing").Where(objectIDs).Find(&objects).Error
	return objects, err
}

// QueryObject 根据参数组合查询object
// 传objectID则只返回对应ID的object，其他参数不生效
// 同时传appID，system，entity，identity参数过滤条件叠加生效
// pageNum为分页查询页码 默认为1
// pageSize为分页查询每页数量, -1可取消限制
// order为排序方式 "字段 desc/asc" 默认按id降序排列
// 多条记录的结果返回分页信息
func QueryObject(c *gin.Context, objectID uint, appID, system, entity, identity string, pageNum, pageSize int, order string) (interface{}, error) {
	if objectID > 0 {
		return GetObjectByID(c, objectID)
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
	scopedb := db.SQLite3("testing").Order(order).Offset(offset).Limit(limit)
	if appID != "" {
		scopedb = scopedb.Where(&demoModels.Object{AppID: appID})
		totalCount = false
	}
	if system != "" {
		scopedb = scopedb.Where(&demoModels.Object{System: system})
		totalCount = false
	}
	if entity != "" {
		scopedb = scopedb.Where(&demoModels.Object{Entity: entity})
		totalCount = false
	}
	if identity != "" {
		scopedb = scopedb.Where(&demoModels.Object{Identity: identity})
		totalCount = false
	}
	items := []demoModels.Object{}
	scopedb = scopedb.Find(&items)
	if totalCount {
		db.SQLite3("testing").Model(&demoModels.Object{}).Count(&count)
	} else {
		scopedb.Count(&count)
	}
	result := map[string]interface{}{
		"items":      items,
		"pagination": utils.Paginate(count, pageNum, pageSize),
	}
	return result, scopedb.Error
}
