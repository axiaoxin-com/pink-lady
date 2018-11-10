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
	now, _ := time.Parse(timeFormat, "2018-11-10 18:52:25.123")
	s.DT = now
	s.JT = JSONTime{Time: now}
	j, _ := json.Marshal(s)
	err := json.Unmarshal(j, &s)
	if err != nil {
		t.Error(err)
	}
	if s.JT.String() != "2018-11-10 18:52:25" {
		t.Error("json time field format error")
	}
}
