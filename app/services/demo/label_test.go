package demo

import (
	"io"
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/db"
	"github.com/axiaoxin/pink-lady/app/logging"
	demoModels "github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/utils"
)

func TestAddLabel(t *testing.T) {
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
	logging.InitLogger()

	db.InitGorm()
	defer db.SQLite3("testing").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	id1, err := AddLabel(nil, "labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	id2, err := AddLabel(nil, "labelname", "labelremark")
	if id1 != id2 || err != nil {
		t.Error("same name should return same id ", err)
	}
	id3, err := AddLabel(nil, "labelname", "labelremark1")
	if id1 != id3 || err != nil {
		t.Error("same name should return same id ", err)
	}
	label, err := GetLabelByID(nil, id3)
	if label.Remark != "labelremark1" || err != nil {
		t.Error("label remark should update ", err)
	}
}

func TestGetLabelByID(t *testing.T) {
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

	db.InitGorm()
	defer db.SQLite3("testing").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	id, err := AddLabel(nil, "labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	label, err := GetLabelByID(nil, id)
	if err != nil || label.ID != id {
		t.Error("get label by id error ", err)
	}
}

func TestGetLabelsByIDs(t *testing.T) {
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

	db.InitGorm()
	defer db.SQLite3("testing").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	id, err := AddLabel(nil, "labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	labels, err := GetLabelsByIDs(nil, []uint{id})
	if err != nil || len(labels) != 1 {
		t.Error("get labels by ids error ", err)
	}
}

func TestQueryLabel(t *testing.T) {
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
	db.InitGorm()
	defer db.SQLite3("testing").Close()
	db.SQLite3("testing").AutoMigrate(&demoModels.Label{}, &demoModels.Object{})

	// init data
	id, err := AddLabel(nil, "labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	_, err = AddLabel(nil, "labelname1", "labelremark1")
	if err != nil {
		t.Error(err)
	}

	// test query by id
	labels, err := GetLabelsByIDs(nil, []uint{id})
	if err != nil || len(labels) != 1 {
		t.Error("get labels by ids error ", err)
	}

	data, err := QueryLabel(nil, id, "", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	label, ok := data.(demoModels.Label)
	if !ok {
		t.Error("not return a label model ", ok, data)
	}
	if label.ID != id {
		t.Error("id not equal")
	}

	// test get by name
	data, err = QueryLabel(nil, 0, "labelname", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	label, ok = data.(demoModels.Label)
	if !ok {
		t.Error("not return a label model ", ok, data)
	}
	if label.ID != id {
		t.Error("id not equal")
	}

	// test default query params
	data, err = QueryLabel(nil, 0, "", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	result, ok := data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items := result["items"].([]demoModels.Label)
	pagi := result["pagination"].(utils.Pagination)
	if len(items) != 2 {
		t.Error("items should have 2 ", items)
	}
	if items[0].ID != 1 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 2 || pagi.PagesCount != 1 || pagi.PageNum != 1 || pagi.PageSize != -1 || pagi.HasPrev != false || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 1 {
		t.Error("pagination error ", pagi)
	}

	// test order
	data, err = QueryLabel(nil, 0, "", "", 0, 0, "id desc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Label)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 2 {
		t.Error("items should have 2 ", items)
	}
	if items[0].ID != 2 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 2 || pagi.PagesCount != 1 || pagi.PageNum != 1 || pagi.PageSize != -1 || pagi.HasPrev != false || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 1 {
		t.Error("pagination error ", pagi)
	}

	// test page size
	data, err = QueryLabel(nil, 0, "", "", 1, 1, "id asc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Label)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 1 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 2 || pagi.PagesCount != 2 || pagi.PageNum != 1 || pagi.PageSize != 1 || pagi.HasPrev != false || pagi.HasNext != true || pagi.PrevPageNum != 1 || pagi.NextPageNum != 2 {
		t.Error("pagination error ", pagi)
	}

	// test page size order
	data, err = QueryLabel(nil, 0, "", "", 1, 1, "id desc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Label)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 2 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 2 || pagi.PagesCount != 2 || pagi.PageNum != 1 || pagi.PageSize != 1 || pagi.HasPrev != false || pagi.HasNext != true || pagi.PrevPageNum != 1 || pagi.NextPageNum != 2 {
		t.Error("pagination error ", pagi)
	}

	// test page num
	data, err = QueryLabel(nil, 0, "", "", 2, 1, "id asc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Label)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 2 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 2 || pagi.PagesCount != 2 || pagi.PageNum != 2 || pagi.PageSize != 1 || pagi.HasPrev != true || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 2 {
		t.Error("pagination error ", pagi)
	}

	// test filter by remark 模糊查询
	data, err = QueryLabel(nil, 0, "", "labelremark", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Label)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 2 {
		t.Error("items should have 2 ", items)
	}
	if items[0].ID != 1 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 2 || pagi.PagesCount != 1 || pagi.PageNum != 1 || pagi.PageSize != -1 || pagi.HasPrev != false || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 1 {
		t.Error("pagination error ", pagi)
	}
}
