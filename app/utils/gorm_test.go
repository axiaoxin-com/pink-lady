package utils

import (
	"os"
	"testing"
)

func TestInitGormDB(t *testing.T) {
	db := "/tmp/pink-lady-unit-test.db"
	err := InitGormDB("sqlite3", "", db, "", "", 0, 0, 0, true)
	if DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer DB.Close()
	defer os.Remove(db)
}
