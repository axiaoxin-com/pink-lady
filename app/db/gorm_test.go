package db

import (
	"io"
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/axiaoxin/pink-lady/app/utils"
)

func TestNewSQLite3Instance(t *testing.T) {
	// test sqlite3
	dbname := "/tmp/pink-lady-unit-test.db"
	db, err := NewSQLite3Instance(dbname, true, 10, 10, 10)
	if db == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	_, err = os.Stat(dbname)
	if err != nil && os.IsNotExist(err) {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbname)

	// TODO: mock mysql, postgres, mssql
}

func TestInitGorm(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	confile, err := os.Open("../config.toml.example")
	if err != nil {
		t.Error(err)
	}
	defer confile.Close()
	tmpConfile, err := os.Create("./config.toml")
	if err != nil {
		t.Error(err)
	}
	defer tmpConfile.Close()
	io.Copy(tmpConfile, confile)
	// 清理测试用的配置文件
	defer func() { os.Remove("./config.toml") }()
	workdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	utils.InitViper(workdir, "config", "envPrefix")
	logging.InitLogger()

	InitGorm()
	if InstanceMap["sqlite3"]["default"] == nil {
		t.Error("InitGorm failed")
	}
	t.Logf("====> InstanceMap %+v", InstanceMap)
}
