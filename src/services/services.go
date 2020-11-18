// Package services 加载或初始化外部依赖服务
package services

import "github.com/spf13/viper"

// Init 相关依赖服务的初始化或加载操作
func Init() error {
	env := viper.GetString("env")
	db, err := GormMySQL(env)
	if err != nil {
		panic(env + " get gorm mysql instance error:" + err.Error())
	}
	DB = db

	return nil
}
