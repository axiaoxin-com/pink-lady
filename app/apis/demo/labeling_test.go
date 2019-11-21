package demo

import (
	"io"
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/db"
	"github.com/axiaoxin/pink-lady/app/models/demo"
	"github.com/axiaoxin/pink-lady/app/router"
	"github.com/axiaoxin/pink-lady/app/utils"

	jsoniter "github.com/json-iterator/go"
)

func TestLabelingAPIs(t *testing.T) {
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
	db.SQLite3("testing").AutoMigrate(&demo.Label{}, &demo.Object{})

	r := router.SetupRouter("test", "", false)
	r.POST("/l", AddLabel)
	r.POST("/o", AddObject)
	r.POST("/", AddLabeling)
	r.POST("/r", ReplaceLabeling)
	r.POST("/d", DeleteLabeling)
	r.GET("/bl/:id", GetLabelingByLabelID)
	r.GET("/bo/:id", GetLabelingByObjectID)
	// init data
	w := utils.TestingPOSTRequest(r, "/l", `{"name": "name1", "remark": "remark1"}`)
	body := jsoniter.Get(w.Body.Bytes())
	l1id := body.Get("data").ToInt()
	if l1id != 1 {
		t.Error("id should = 1")
	}
	w = utils.TestingPOSTRequest(r, "/l", `{"name": "name2", "remark": "remark2"}`)
	body = jsoniter.Get(w.Body.Bytes())
	l2id := body.Get("data").ToInt()
	if l2id != 2 {
		t.Error("id should = 2")
	}
	w = utils.TestingPOSTRequest(r, "/o", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id1"}`)
	body = jsoniter.Get(w.Body.Bytes())
	o1id := body.Get("data").ToInt()
	if o1id != 1 {
		t.Error("id should = 1")
	}
	w = utils.TestingPOSTRequest(r, "/o", `{"appID": "appid1", "system": "sys1", "entity": "e1", "identity": "id2"}`)
	body = jsoniter.Get(w.Body.Bytes())
	o2id := body.Get("data").ToInt()
	if o2id != 2 {
		t.Error("id should = 2")
	}

	// test add labeling
	w = utils.TestingPOSTRequest(r, "/", `{"objectIDs": [1, 2], "labelIDs": [1, 2]}`)
	body = jsoniter.Get(w.Body.Bytes())
	data := body.Get("data")
	if data.Size() != 4 {
		t.Error("data should have 4 results ", data.ToString())
	}
	for i := 0; i < data.Size(); i++ {
		if data.Get(i).Get("result").ToString() != "ok" {
			t.Error("result not ok ", i)
		}
	}

	// test get by id
	w = utils.TestingGETRequest(r, "/bl/1")
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 2 {
		t.Error("data should have 2 results ", body.ToString())
	}
	w = utils.TestingGETRequest(r, "/bo/1")
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 2 {
		t.Error("data should have 2 results ", body.ToString())
	}

	// test replace
	w = utils.TestingPOSTRequest(r, "/r", `{"objectIDs": [1], "labelIDs": [1]}`)
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 1 {
		t.Error("data should have 1 results ", body.ToString())
	}
	w = utils.TestingGETRequest(r, "/bl/1")
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 2 {
		t.Error("data should have 2 results ", body.ToString())
	}
	w = utils.TestingGETRequest(r, "/bo/1")
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 1 {
		t.Error("data should have 1 results ", body.ToString())
	}

	// test delete
	w = utils.TestingPOSTRequest(r, "/d", `{"objectIDs": [1], "labelIDs": [1]}`)
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 1 {
		t.Error("data should have 1 results ", body.ToString())
	}
	w = utils.TestingGETRequest(r, "/bl/1")
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 1 {
		t.Error("data should have 1 results ", body.ToString())
	}
	w = utils.TestingGETRequest(r, "/bo/1")
	body = jsoniter.Get(w.Body.Bytes())
	data = body.Get("data")
	if data.Size() != 0 {
		t.Error("data should have 0 results ", body.ToString())
	}

}
