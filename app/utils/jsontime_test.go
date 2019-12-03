package utils

import (
	"encoding/json"
	"testing"
	"time"
)

type S struct {
	DT time.Time
	JT JSONTime
}

func TestJSONTime(t *testing.T) {
	s := S{}
	now, _ := time.Parse(TimeFormat, "2018-11-10 18:52:25.123")
	s.DT = now
	s.JT = JSONTime{Time: now}
	j, _ := json.Marshal(s)
	err := json.Unmarshal(j, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.JT.String() != "2018-11-10 18:52:25" {
		t.Fatal("json time field format error")
	}
	if _, err := s.JT.Value(); err != nil {
		t.Fatal(err)
	}
	if err := s.JT.Scan(""); err == nil {
		t.Fatal("scan convert to timestamp")
	}
	if err := s.JT.Scan(time.Now()); err != nil {
		t.Fatal("scan convert err", err)
	}
}
