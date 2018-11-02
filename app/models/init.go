// package models save database models
// define your database models in the package
package models

import (
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type BaseModel gorm.Model

// Migrate run AutoMigrate to create database tables by models
// add your models &MODEL{} in AutoMigrate()
// running after InitGormDB
func Migrate() {
	if err := utils.DB.AutoMigrate(
	// add your models which need to migrate at here

	).Error; err != nil {
		logrus.Error(err)
	}
}
