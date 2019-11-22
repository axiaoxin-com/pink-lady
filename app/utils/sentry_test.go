package utils

import "testing"

func TestInitSentry(t *testing.T) {
	err := InitSentry()
	if err != nil {
		t.Error(err)
	}
}
