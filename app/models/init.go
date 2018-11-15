// package models save database models
// define your database models in the package
package models

import (
	"gin-skeleton/app/utils"
	"github.com/sirupsen/logrus"
)

type BaseModel struct {
	ID        uint            `gorm:"primary_key" json:"id"`
	CreatedAt utils.JSONTime  `json:"createdAt"`
	UpdatedAt utils.JSONTime  `json:"updatedAt"`
	DeletedAt *utils.JSONTime `sql:"index" json:"-"`
}

// Models save your models like &MODEL{} at there which will be auto migrate when server starting
var Models = []interface{}{}

// Migrate run AutoMigrate to create database tables which in Models array
// running after InitGormDB
func Migrate() {
	if err := utils.DB.AutoMigrate(Models...).Error; err != nil {
		logrus.Error(err)
	}
}
