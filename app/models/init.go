// Package models save database models
// define your database models in the package
package models

import (
	"github.com/axiaoxin/pink-lady/app/utils"
)

// BaseModel you should define you model with BaseModel
type BaseModel struct {
	ID        uint            `gorm:"primary_key" json:"id"`
	CreatedAt utils.JSONTime  `json:"createdAt"`
	UpdatedAt utils.JSONTime  `json:"updatedAt"`
	DeletedAt *utils.JSONTime `sql:"index" json:"-"`
}
