package demo

import (
	"github.com/axiaoxin/pink-lady/app/models"
)

// Label model of label
type Label struct {
	models.BaseModel
	Name    string    `gorm:"type:varchar(32);not null;default:'';unique_index" json:"name" binding:"required,max=32" example:"server:linux"`
	Remark  string    `gorm:"type:varchar(64);not null;default:''" json:"remark" binding:"max=64" example:"Linux服务器"`
	Objects []*Object `gorm:"many2many:labeling;" json:"objects,omitempty"`
}

// TableName define the Label model's tabel name
func (Label) TableName() string {
	return "label"
}
