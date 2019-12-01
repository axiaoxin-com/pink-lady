package utils

import (
	"os"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	if err != nil {
		t.Fatal("GetLocalIP err:", err)
	}
	if ip == "" {
		t.Fatal("GetLocalIP ip is empty:")
	}
}

func TestStructToURLValues(t *testing.T) {
	v := StructToURLValues(&struct {
		I int `json:"int_i"`
		S string
	}{666, "testing"})
	if v.Get("int_i") != "666" || v.Get("S") != "testing" {
		t.Fatalf("convert failed: %+v", v)
	}
	if v.Encode() != "S=testing&int_i=666" {
		t.Fatal("encode error:", v.Encode())
	}
}

func TestCopyFile(t *testing.T) {
	if err := CopyFile("../config.toml.example", "/tmp/tst-cp-file.cfg"); err != nil {
		t.Fatal("copy file err:", err)
	}
	if _, err := os.Stat("/tmp/tst-cp-file.cfg"); os.IsNotExist(err) {
		t.Fatal("copy file is not exists")
	}
}

func TestRemoveAllWhiteSpace(t *testing.T) {
	rs := RemoveAllWhiteSpace(" a\tb\n \n\nc d   e ")
	if rs != "abcde" {
		t.Fatal("RemoveAllWhiteSpace error:", rs)
	}
}
