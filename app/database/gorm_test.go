package database

import (
	"os"
	"pink-lady/app/logging"
	"pink-lady/app/utils"
	"testing"
)

func TestNewSQLite3Instance(t *testing.T) {
	// test sqlite3
	dbname := "/tmp/pink-lady-unit-test.db"
	db, err := NewSQLite3Instance(dbname, true, 10, 10, 10)
	if db == nil || err != nil {
		t.Fatal("init DB fail ", err)
	}
	_, err = os.Stat(dbname)
	if err != nil && os.IsNotExist(err) {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(dbname)

	// TODO: mock mysql, postgres, mssql
}

func TestInitGorm(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	utils.CopyFile("../config.toml.example", "./config.toml")
	// 清理测试用的配置文件
	defer func() { os.Remove("./config.toml") }()
	utils.InitViper("./", "config", "")
	logging.InitLogger()

	InitGorm()
	if InstanceMap["sqlite3"]["default"] == nil {
		t.Fatal("InitGorm failed")
	}
}
