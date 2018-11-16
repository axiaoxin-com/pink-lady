package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime custom format for json time field
type JSONTime struct {
	time.Time
}

// TimeFormat define the json time filed format
const (
	TimeFormat = "2006-01-02 15:04:05"
)

// UnmarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	t.Time = now
	return err
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

// String return  %Y-%m-%d %H:%M:%S
func (t JSONTime) String() string {
	return t.Time.Format(TimeFormat)
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
