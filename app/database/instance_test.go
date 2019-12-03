package database

import (
	"testing"
)

func TestSQLite3(t *testing.T) {
	i := SQLite3("none-exists")
	if i != nil {
		t.Fatal("should be nil")
	}
}

func TestUTDB(t *testing.T) {
	db := UTDB()
	if db == nil {
		t.Fatal("UTDB实例为nil")
	}
}

func TestMySQL(t *testing.T) {
	i := MySQL("none-exists")
	if i != nil {
		t.Fatal("should be nil")
	}
}

func TestMsSQL(t *testing.T) {
	i := MsSQL("none-exists")
	if i != nil {
		t.Fatal("should be nil")
	}
}

func TestPostgres(t *testing.T) {
	i := Postgres("none-exists")
	if i != nil {
		t.Fatal("should be nil")
	}
}
