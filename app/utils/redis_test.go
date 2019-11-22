package utils

import (
	"testing"
)

func TestInitRedis(t *testing.T) {
	s, err := MockRedis()
	if err != nil {
		t.Errorf("mock redis failed %s", err)
	}
	defer s.Close()

	err = InitRedis()
	if Redis == nil || err != nil {
		t.Errorf("init redis single instance failed %s", err)
	}
	err = Redis.Set("k", "v", 0).Err()
	if err != nil {
		t.Error(err)
	}
	v, err := Redis.Get("k").Result()
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
