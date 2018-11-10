package utils

import (
	"testing"
)

func TestInitGormDB(t *testing.T) {
	err := InitGormDB("sqlite3", "", "/tmp/gin-skeleton-unit-test.db", "", "", 0, 0, 0, true)
	if DB == nil || err != nil {
		t.Error("init DB fail ", err)
	}
	defer DB.Close()
}
