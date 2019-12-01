// Package models save database models
// define your database models in the package
package models

import (
	"pink-lady/app/utils"
)

// BaseModel you should define you model with BaseModel
type BaseModel struct {
	ID        int64          `gorm:"primary_key,column:id" json:"id" example:"0"`     // 主键ID
	CreatedAt utils.JSONTime `gorm:"column:created_at" json:"created_at" example:"-"` // 创建时间
	UpdatedAt utils.JSONTime `gorm:"column:updated_at" json:"updated_at" example:"-"` // 更新时间
}
