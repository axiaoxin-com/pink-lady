package demo

import (
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/models"
	demoModels "github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/utils"
)

func init() {
	utils.InitLogger(os.Stdout, "debug", "text")
}

func TestAddLabel(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	models.Migrate()
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer utils.DB.Close()
	defer os.Remove(db)

	id1, err := AddLabel("labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	id2, err := AddLabel("labelname", "labelremark")
	if id1 != id2 || err != nil {
		t.Error("same name should return same id ", err)
	}
	id3, err := AddLabel("labelname", "labelremark1")
	if id1 != id3 || err != nil {
		t.Error("same name should return same id ", err)
	}
	label, err := GetLabelByID(id3)
	if label.Remark != "labelremark1" || err != nil {
		t.Error("label remark should update ", err)
	}
}

func TestGetLabelByID(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)

	id, err := AddLabel("labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	label, err := GetLabelByID(id)
	if err != nil || label.ID != id {
		t.Error("get label by id error ", err)
	}
}

func TestGetLabelsByIDs(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)

	id, err := AddLabel("labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	labels, err := GetLabelsByIDs([]uint{id})
	if err != nil || len(labels) != 1 {
		t.Error("get labels by ids error ", err)
	}
}

func TestQueryLabel(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)

	// init data
	id, err := AddLabel("labelname", "labelremark")
	if err != nil {
		t.Error(err)
	}
	_, err = AddLabel("labelname1", "labelremark1")
	if err != nil {
		t.Error(err)
	}

	// test query by id
	labels, err := GetLabelsByIDs([]uint{id})
	if err != nil || len(labels) != 1 {
		t.Error("get labels by ids error ", err)
	}

	data, err := QueryLabel(id, "", "", 0, 0, "")
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
	data, err = QueryLabel(0, "labelname", "", 0, 0, "")
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
	data, err = QueryLabel(0, "", "", 0, 0, "")
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
	data, err = QueryLabel(0, "", "", 0, 0, "id desc")
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
	data, err = QueryLabel(0, "", "", 1, 1, "id asc")
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
	data, err = QueryLabel(0, "", "", 1, 1, "id desc")
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
	data, err = QueryLabel(0, "", "", 2, 1, "id asc")
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
	data, err = QueryLabel(0, "", "labelremark", 0, 0, "")
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
