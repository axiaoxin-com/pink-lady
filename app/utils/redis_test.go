package utils

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNewRedisClient(t *testing.T) {
	s, err := MockRedis()
	if err != nil {
		t.Fatalf("mock redis failed %s", err)
	}
	defer s.Close()

	r, err := NewRedisClient(s.Addr(), "", 0)
	if r == nil || err != nil {
		t.Fatal("init redis single instance failed", err)
	}
	err = r.Set("k", "v", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	v, err := r.Get("k").Result()
	if v != "v" || err != nil {
		t.Fatalf("redis get %s %s", v, err)
	}

	/*
		Redis = nil
		err = InitRedis("sentinel", s.Addr(), "", 0, "master")
		if Redis == nil || err != nil {
			t.Fatalf("init redis sentinel failed %s", err)
		}

		Redis = nil
		err = InitRedis("cluster", s.Addr(), "", 0, "")
		if Redis == nil || err != nil {
			t.Fatalf("init redis cluster failed %s", err)
		}
	*/
}

func TestInitRedis(t *testing.T) {
	err := InitRedis()
	if err == nil {
		t.Error("没有设置redis配置应该返回错误")
	}
	viper.SetDefault("redis.mode", "cluster")
	viper.SetDefault("redis.address", "127.0.0.1:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.dbindex", 0)
	InitRedis()
}
