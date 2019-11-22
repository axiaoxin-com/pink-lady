package utils

import (
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	s, err := MockRedis()
	if err != nil {
		t.Errorf("mock redis failed %s", err)
	}
	defer s.Close()

	r, err := NewRedisClient(s.Addr(), "", 0)
	if r == nil || err != nil {
		t.Error("init redis single instance failed", err)
	}
	err = r.Set("k", "v", 0).Err()
	if err != nil {
		t.Error(err)
	}
	v, err := r.Get("k").Result()
	if v != "v" || err != nil {
		t.Errorf("redis get %s %s", v, err)
	}

	/*
		Redis = nil
		err = InitRedis("sentinel", s.Addr(), "", 0, "master")
		if Redis == nil || err != nil {
			t.Errorf("init redis sentinel failed %s", err)
		}

		Redis = nil
		err = InitRedis("cluster", s.Addr(), "", 0, "")
		if Redis == nil || err != nil {
			t.Errorf("init redis cluster failed %s", err)
		}
	*/
}
