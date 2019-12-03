package response

import (
	"errors"
	"testing"
)

func TestNewRetCode(t *testing.T) {
	if rc := NewRetCode(1, "msg"); rc == nil {
		t.Fatal("NewRetCode return nil")
	}
}

func TestDecode(t *testing.T) {
	rc := NewRetCode(1, "msg", errors.New("err"))
	code, msg, err := rc.Decode()
	if code.(int) != 1 || msg != "msg" || err == nil {
		t.Fatal("Decode error")
	}
}

func TestError(t *testing.T) {
	rc := NewRetCode(1, "msg", errors.New("err"), errors.New("xx"))
	err := rc.Error()
	if err != "err; xx" {
		t.Fatal("Error error")
	}
}

func TestCode(t *testing.T) {
	rc := NewRetCode(1, "msg")
	code := rc.Code()
	if code.(int) != 1 {
		t.Fatal("Code error")
	}
}

func TestMsg(t *testing.T) {
	rc := NewRetCode(1, "msg")
	msg := rc.Msg()
	if msg != "msg" {
		t.Fatal("Msg error", msg)
	}
}

func TestErrs(t *testing.T) {
	rc := NewRetCode(1, "msg")
	errs := rc.Errs()
	if len(errs) != 0 {
		t.Fatal("Errs error", errs)
	}
	rc = NewRetCode(1, "msg", errors.New("x"))
	errs = rc.Errs()
	if len(errs) != 1 {
		t.Fatal("Errs error")
	}
}

func TestSetMsg(t *testing.T) {
	rc := NewRetCode(1, "msg")
	nmsg := rc.SetMsg("asdf").Msg()
	if rc.Msg() != "msg" || nmsg != "asdf" {
		t.Fatal("SetMsg error")
	}
}

func TestAppendError(t *testing.T) {
	rc := NewRetCode(1, "msg")
	nerrs := rc.AppendError(errors.New("err")).Errs()
	if len(rc.Errs()) != 0 || len(nerrs) != 1 {
		t.Fatal("AppendError error")
	}
}
