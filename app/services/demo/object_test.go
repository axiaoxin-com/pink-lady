package demo

import (
	"io"
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/db"
	demoModels "github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/utils"
)

func TestAddObject(t *testing.T) {
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

	id1, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	id2, err := AddObject(nil, "appid", "sys", "entity", "id")
	if id1 != id2 || err != nil {
		t.Error("same name should return same id")
	}
}

func TestGetObjectByID(t *testing.T) {
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

	id, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	object, err := GetObjectByID(nil, id)
	if err != nil || object.ID != id {
		t.Error("get object by id error ", err)
	}
}

func TestGetObjectsByIDs(t *testing.T) {
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

	id, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	objects, err := GetObjectsByIDs(nil, []uint{id})
	if err != nil || len(objects) != 1 {
		t.Error("get objects by ids error ", err)
	}
}

func TestQueryObject(t *testing.T) {
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

	// init data
	id, err := AddObject(nil, "appid", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	_, err = AddObject(nil, "appid1", "sys", "entity", "id")
	if err != nil {
		t.Error(err)
	}
	_, err = AddObject(nil, "appid", "sys1", "entity", "id")
	if err != nil {
		t.Error(err)
	}

	// test query by id
	objects, err := GetObjectsByIDs(nil, []uint{id})
	if err != nil || len(objects) != 1 {
		t.Error("get objects by ids error ", err)
	}

	data, err := QueryObject(nil, id, "", "", "", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	object, ok := data.(demoModels.Object)
	if !ok {
		t.Error("not return a object model ", ok, data)
	}
	if object.ID != id {
		t.Error("id not equal")
	}

	// test default query params
	data, err = QueryObject(nil, 0, "", "", "", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	result, ok := data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items := result["items"].([]demoModels.Object)
	pagi := result["pagination"].(utils.Pagination)
	if len(items) != 3 {
		t.Error("items should have 3 ", items)
	}
	if items[0].ID != 1 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 3 || pagi.PagesCount != 1 || pagi.PageNum != 1 || pagi.PageSize != -1 || pagi.HasPrev != false || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 1 {
		t.Error("pagination error ", pagi)
	}

	// test order
	data, err = QueryObject(nil, 0, "", "", "", "", 0, 0, "id desc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Object)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 3 {
		t.Error("items should have 3 ", items)
	}
	if items[0].ID != 3 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 3 || pagi.PagesCount != 1 || pagi.PageNum != 1 || pagi.PageSize != -1 || pagi.HasPrev != false || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 1 {
		t.Error("pagination error ", pagi)
	}

	// test page size
	data, err = QueryObject(nil, 0, "", "", "", "", 1, 1, "id asc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Object)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 1 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 3 || pagi.PagesCount != 3 || pagi.PageNum != 1 || pagi.PageSize != 1 || pagi.HasPrev != false || pagi.HasNext != true || pagi.PrevPageNum != 1 || pagi.NextPageNum != 2 {
		t.Error("pagination error ", pagi)
	}

	// test page size order
	data, err = QueryObject(nil, 0, "", "", "", "", 1, 1, "id desc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Object)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 3 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 3 || pagi.PagesCount != 3 || pagi.PageNum != 1 || pagi.PageSize != 1 || pagi.HasPrev != false || pagi.HasNext != true || pagi.PrevPageNum != 1 || pagi.NextPageNum != 2 {
		t.Error("pagination error ", pagi)
	}

	// test page num
	data, err = QueryObject(nil, 0, "", "", "", "", 2, 1, "id asc")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Object)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 2 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 3 || pagi.PagesCount != 3 || pagi.PageNum != 2 || pagi.PageSize != 1 || pagi.HasPrev != true || pagi.HasNext != true || pagi.PrevPageNum != 1 || pagi.NextPageNum != 3 {
		t.Error("pagination error ", pagi)
	}

	// test filter by appid
	data, err = QueryObject(nil, 0, "appid", "", "", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Object)
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

	// test filter by sys
	data, err = QueryObject(nil, 0, "", "sys1", "", "", 0, 0, "")
	if err != nil {
		t.Error(err)
	}
	result, ok = data.(map[string]interface{})
	if !ok {
		t.Error("not return a map data ", ok, data)
	}
	items = result["items"].([]demoModels.Object)
	pagi = result["pagination"].(utils.Pagination)
	if len(items) != 1 {
		t.Error("items should have 1 ", items)
	}
	if items[0].ID != 3 {
		t.Error("items order err")
	}
	if pagi.ItemsCount != 1 || pagi.PagesCount != 1 || pagi.PageNum != 1 || pagi.PageSize != -1 || pagi.HasPrev != false || pagi.HasNext != false || pagi.PrevPageNum != 1 || pagi.NextPageNum != 1 {
		t.Error("pagination error ", pagi)
	}
}
