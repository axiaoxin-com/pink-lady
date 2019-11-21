package demo

import (
	"github.com/axiaoxin/pink-lady/app/models"
)

// Object model of object
type Object struct {
	models.BaseModel
	AppID    string   `gorm:"type:varchar(16);not null;default:'';unique_index:idx_appid_system_entity_identity" json:"appID"`
	System   string   `gorm:"type:varchar(64);not null;default:'';unique_index:idx_appid_system_entity_identity" json:"system"`
	Entity   string   `gorm:"type:varchar(64);not null;default:'';unique_index:idx_appid_system_entity_identity" json:"entity"`
	Identity string   `gorm:"type:varchar(64);not null;default:'';unique_index:idx_appid_system_entity_identity" json:"identity"`
	Labels   []*Label `gorm:"many2many:labeling;" json:"labels,omitempty"`
}

// TableName define the model's table name
func (Object) TableName() string {
	return "object"
}
