package demo

import (
	"io"
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/db"
	demoModels "github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/utils"
)

func TestAddLabeling(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	confile, err := os.Open("../../config.toml.example")
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
	defer func() { os.Remove("/tmp/pink-lady-testing.db") }()
	workdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	utils.InitViper(workdir, "config", "envPrefix")

	err = db.InitGorm()
	if db.SQLite3("testing") == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer db.SQLite3("testing").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	// init test
	lid, err := AddLabel(nil, "label1", "remark1")
	if err != nil {
		t.Error(err)
	}
	lid2, err := AddLabel(nil, "label2", "remark2")
	if err != nil {
		t.Error(err)
	}
	oid, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	oid2, err := AddObject(nil, "appid", "sys", "entity", "id2")
	if err != nil {
		t.Error(err)
	}

	// test repeat labeling
	results, err := AddLabeling(nil, []uint{oid, oid2}, []uint{lid, lid2})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 4 {
		t.Error("should 4 items to labeling")
	}
	for _, result := range results {
		if result["result"] != "ok" {
			t.Error(result, "not ok")
		}
	}
	AddLabeling(nil, []uint{oid}, []uint{lid})
	objs, err := GetLabelingByLabelID(nil, lid)
	if err != nil {
		t.Error(err)
	}
	if len(objs) != 2 {
		t.Error("objs should 2")
	}
	labels, err := GetLabelingByObjectID(nil, oid)
	if err != nil {
		t.Error(err)
	}
	if len(labels) != 2 {
		t.Error("objs should 2")
	}
}

func TestReplaceLabeling(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	confile, err := os.Open("../../config.toml.example")
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
	defer func() { os.Remove("/tmp/pink-lady-testing.db") }()
	workdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	utils.InitViper(workdir, "config", "envPrefix")

	err = db.InitGorm()
	if db.SQLite3("testing") == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer db.SQLite3("testing").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	// init test
	lid, err := AddLabel(nil, "label1", "remark1")
	if err != nil {
		t.Error(err)
	}
	lid2, err := AddLabel(nil, "label2", "remark2")
	if err != nil {
		t.Error(err)
	}
	oid, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	oid2, err := AddObject(nil, "appid", "sys", "entity", "id2")
	if err != nil {
		t.Error(err)
	}

	results, err := AddLabeling(nil, []uint{oid, oid2}, []uint{lid, lid2})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 4 {
		t.Error("should 4 items to labeling")
	}
	for _, result := range results {
		if result["result"] != "ok" {
			t.Error(result, "not ok")
		}
	}
	objs, err := GetLabelingByLabelID(nil, lid)
	if err != nil {
		t.Error(err)
	}
	if len(objs) != 2 {
		t.Error("objs should 2")
	}
	labels, err := GetLabelingByObjectID(nil, oid)
	if err != nil {
		t.Error(err)
	}
	if len(labels) != 2 {
		t.Error("objs should 2")
	}

	results, err = ReplaceLabeling(nil, []uint{oid}, []uint{lid})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 1 {
		t.Error("should 1 items to labeling")
	}
	for _, result := range results {
		if result["result"] != "ok" {
			t.Error(result, "not ok")
		}
	}
	objs, err = GetLabelingByLabelID(nil, lid)
	if err != nil {
		t.Error(err)
	}
	if len(objs) != 2 {
		t.Error("objs should 2")
	}
	labels, err = GetLabelingByObjectID(nil, oid)
	if err != nil {
		t.Error(err)
	}
	if len(labels) != 1 {
		t.Error("objs should 1")
	}
}

func TestDeleteLabeling(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	confile, err := os.Open("../../config.toml.example")
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
	defer func() { os.Remove("/tmp/pink-lady-testing.db") }()
	workdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	utils.InitViper(workdir, "config", "envPrefix")

	err = db.InitGorm()
	if db.SQLite3("default") == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer db.SQLite3("default").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	// init test
	lid, err := AddLabel(nil, "label1", "remark1")
	if err != nil {
		t.Error(err)
	}
	lid2, err := AddLabel(nil, "label2", "remark2")
	if err != nil {
		t.Error(err)
	}
	oid, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	oid2, err := AddObject(nil, "appid", "sys", "entity", "id2")
	if err != nil {
		t.Error(err)
	}

	results, err := AddLabeling(nil, []uint{oid, oid2}, []uint{lid, lid2})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 4 {
		t.Error("should 4 items to labeling")
	}
	for _, result := range results {
		if result["result"] != "ok" {
			t.Error(result, "not ok")
		}
	}
	labels, _ := GetLabelingByObjectID(nil, oid)
	if len(labels) != 2 {
		t.Error("should 2")
	}

	results, err = DeleteLabeling(nil, []uint{oid}, []uint{lid})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 1 {
		t.Error("should 1 items to labeling")
	}
	for _, result := range results {
		if result["result"] != "ok" {
			t.Error(result, "not ok")
		}
	}
	labels, _ = GetLabelingByObjectID(nil, oid)
	if len(labels) != 1 {
		t.Error("should 1")
	}

	results, err = DeleteLabeling(nil, []uint{oid}, []uint{lid2})
	if err != nil {
		t.Error(err)
	}
	if len(results) != 1 {
		t.Error("should 1 items to labeling")
	}
	for _, result := range results {
		if result["result"] != "ok" {
			t.Error(result, "not ok")
		}
	}
	labels, _ = GetLabelingByObjectID(nil, oid)
	if len(labels) != 0 {
		t.Error("should 0")
	}
}
