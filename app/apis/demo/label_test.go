package demo

import (
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/apis/router"
	"github.com/axiaoxin/pink-lady/app/models"
	"github.com/axiaoxin/pink-lady/app/utils"

	jsoniter "github.com/json-iterator/go"
)

func TestAddLabel(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)
	r := router.SetupRouter("test", "", false)

	r.POST("/", AddLabel)
	w := utils.TestingPOSTRequest(r, "/", `{"name": "name", "remark": "remark"}`)
	body := jsoniter.Get(w.Body.Bytes())
	id := body.Get("data").ToInt()
	if id != 1 {
		t.Error("add label fail")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"name": "name1", "remark": "remark"}`)
	body = jsoniter.Get(w.Body.Bytes())
	id = body.Get("data").ToInt()
	if id != 2 {
		t.Error("add label fail")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"name": "name", "remark": "remark1"}`)
	body = jsoniter.Get(w.Body.Bytes())
	id = body.Get("data").ToInt()
	if id != 1 {
		t.Error("add label fail")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"noname": "name", "remark": "remark1"}`)
	body = jsoniter.Get(w.Body.Bytes())
	code := body.Get("code").ToInt()
	if code != 3 {
		t.Error("code should be 3 invalidParams")
	}
	if w.Result().StatusCode != 400 {
		t.Error("http status code should be 400")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"name": "name1", "noremark": "remark1"}`)
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("code should be 0 ok")
	}
	if w.Result().StatusCode != 200 {
		t.Error("http status code should be 200")
	}
}

func TestLabel(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)
	r := router.SetupRouter("test", "", false)

	r.GET("/", Label)
	r.POST("/", AddLabel)

	// add label
	w := utils.TestingPOSTRequest(r, "/", `{"name": "name", "remark": "remark"}`)
	body := jsoniter.Get(w.Body.Bytes())
	id := body.Get("data").ToInt()
	if id != 1 {
		t.Error("add label fail")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"name": "name1", "remark": "remark"}`)
	body = jsoniter.Get(w.Body.Bytes())
	id = body.Get("data").ToInt()
	if id != 2 {
		t.Error("add label fail")
	}

	// test no params
	w = utils.TestingGETRequest(r, "/")
	body = jsoniter.Get(w.Body.Bytes())
	code := body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	items := body.Get("data", "items")
	if items.Size() != 2 {
		t.Error("items size should be 2")
	}
	if items.Get(0).Get("id").ToInt() != 1 {
		t.Error("items 0 id should be 1")
	}
	pagi := body.Get("data", "pagination")
	if pagi.Get("itemsCount").ToInt() != 2 || pagi.Get("pagesCount").ToInt() != 1 || pagi.Get("pageNum").ToInt() != 1 || pagi.Get("pageSize").ToInt() != 10 || pagi.Get("hasPrev").ToBool() != false || pagi.Get("hasNext").ToBool() != false || pagi.Get("prevPageNum").ToInt() != 1 || pagi.Get("nextPageNum").ToInt() != 1 {
		t.Error("pagination error")
	}

	// test id params
	w = utils.TestingGETRequest(r, "/?id=1")
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	label := body.Get("data")
	if label.Get("name").ToString() != "name" {
		t.Error("label id 1 name should be name")
	}

	// test name params
	w = utils.TestingGETRequest(r, "/?name=name")
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	label = body.Get("data")
	if label.Get("id").ToInt() != 1 {
		t.Error("label name name id should be 1")
	}

	// test order params
	w = utils.TestingGETRequest(r, "/?order=id desc")
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	items = body.Get("data", "items")
	if items.Size() != 2 {
		t.Error("items size should be 2")
	}
	if items.Get(0).Get("id").ToInt() != 2 {
		t.Error("items 0 id should be 2")
	}
	pagi = body.Get("data", "pagination")
	if pagi.Get("itemsCount").ToInt() != 2 || pagi.Get("pagesCount").ToInt() != 1 || pagi.Get("pageNum").ToInt() != 1 || pagi.Get("pageSize").ToInt() != 10 || pagi.Get("hasPrev").ToBool() != false || pagi.Get("hasNext").ToBool() != false || pagi.Get("prevPageNum").ToInt() != 1 || pagi.Get("nextPageNum").ToInt() != 1 {
		t.Error("pagination error")
	}

	// test page params
	w = utils.TestingGETRequest(r, "/?pageNum=2&pageSize=1")
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	items = body.Get("data", "items")
	if items.Size() != 1 {
		t.Error("items size should be 1")
	}
	if items.Get(0).Get("id").ToInt() != 2 {
		t.Error("items 0 id should be 2")
	}
	pagi = body.Get("data", "pagination")
	if pagi.Get("itemsCount").ToInt() != 2 || pagi.Get("pagesCount").ToInt() != 2 || pagi.Get("pageNum").ToInt() != 2 || pagi.Get("pageSize").ToInt() != 1 || pagi.Get("hasPrev").ToBool() != true || pagi.Get("hasNext").ToBool() != false || pagi.Get("prevPageNum").ToInt() != 1 || pagi.Get("nextPageNum").ToInt() != 2 {
		t.Error("pagination error")
	}
}
