// package models save database models
// define your database models in the package
package models

import (
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type BaseModel gorm.Model

// Models save your models like &MODEL{} at there which will be auto migrate when server starting
var Models = []interface{}{}

// Migrate run AutoMigrate to create database tables by models
// add your models
// running after InitGormDB
func Migrate() {
	if err := utils.DB.AutoMigrate(Models...).Error; err != nil {
		logrus.Error(err)
	}
}
