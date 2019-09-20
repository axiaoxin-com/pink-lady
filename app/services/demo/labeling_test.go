package demo

import (
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/models"
	"github.com/axiaoxin/pink-lady/app/utils"
)

func TestAddLabeling(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	models.Migrate()
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer utils.DB.Close()
	defer os.Remove(db)

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
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	models.Migrate()
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer utils.DB.Close()
	defer os.Remove(db)

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
	db := "/tmp/pink-lady-unit-test.db"
	err := utils.InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	models.Migrate()
	if utils.DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer utils.DB.Close()
	defer os.Remove(db)

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
