package demo

import (
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/models"
	"github.com/axiaoxin/pink-lady/app/router"
	"github.com/axiaoxin/pink-lady/app/utils"

	jsoniter "github.com/json-iterator/go"
)

func TestAddObject(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)
	r := router.SetupRouter("test", "", false)

	r.POST("/", AddObject)
	w := utils.TestingPOSTRequest(r, "/", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id1"}`)
	body := jsoniter.Get(w.Body.Bytes())
	id := body.Get("data").ToInt()
	if id != 1 {
		t.Error("add object fail")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id2"}`)
	body = jsoniter.Get(w.Body.Bytes())
	id = body.Get("data").ToInt()
	if id != 2 {
		t.Error("add object error")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id1"}`)
	body = jsoniter.Get(w.Body.Bytes())
	id = body.Get("data").ToInt()
	if id != 1 {
		t.Error("object id should = 1")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"noappID": "appid1", "system": "sys1", "entity": "e1", "identity": "id1"}`)
	body = jsoniter.Get(w.Body.Bytes())
	code := body.Get("code").ToInt()
	if code != 3 {
		t.Error("code should be 3 invalidParams")
	}
	if w.Result().StatusCode != 400 {
		t.Error("http status code should be 400")
	}
}

func TestObject(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	models.Migrate()
	defer utils.DB.Close()
	defer os.Remove(db)
	r := router.SetupRouter("test", "", false)

	r.GET("/", Object)
	r.POST("/", AddObject)

	// add object
	w := utils.TestingPOSTRequest(r, "/", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id1"}`)
	body := jsoniter.Get(w.Body.Bytes())
	id := body.Get("data").ToInt()
	if id != 1 {
		t.Error("add object fail")
	}
	w = utils.TestingPOSTRequest(r, "/", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id2"}`)
	body = jsoniter.Get(w.Body.Bytes())
	id = body.Get("data").ToInt()
	if id != 2 {
		t.Error("add object fail")
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
	obj := body.Get("data")
	if obj.Get("identity").ToString() != "id1" {
		t.Error("should be id1")
	}

	// test field params
	w = utils.TestingGETRequest(r, "/?entity=e1")
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	items = body.Get("data", "items")
	if items.Size() != 2 {
		t.Error("items size should be 2")
	}
	// test field params
	w = utils.TestingGETRequest(r, "/?entity=e1&identity=id2")
	body = jsoniter.Get(w.Body.Bytes())
	code = body.Get("code").ToInt()
	if code != 0 {
		t.Error("json code is not 0")
	}
	items = body.Get("data", "items")
	if items.Size() != 1 {
		t.Error("items size should be 1")
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
