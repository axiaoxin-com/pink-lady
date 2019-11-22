package db

import (
	"os"
	"testing"
)

func TestSQLite3(t *testing.T) {
	i := SQLite3("none-exists")
	if i != nil {
		t.Error("should be nil")
	}
	i = SQLite3("testing")
	defer os.Remove("/tmp/pink-lady-testing.db")
	if i == nil {
		t.Error("instance is nil")
	}
	if _, err := os.Stat("/tmp/pink-lady-testing.db"); err != nil && os.IsNotExist(err) {
		t.Error(err)
	}
}
