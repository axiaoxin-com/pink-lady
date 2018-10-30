package utils

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	REDIS_SINGLE_INSTANCE_MODE = "single-instance"
	REDIS_SENTINEL_MODE        = "sentinel"
	REDIS_CLUSTER_MODE         = "cluster"
	REDIS_ADDRS_SEPARATOR      = ","
)

var (
	Redis        *redis.Client
	RedisCluster *redis.ClusterClient
)

func InitRedis(mode string, addr string, password string, db int, master string) {
	mode = strings.ToLower(mode)
	if mode == REDIS_SINGLE_INSTANCE_MODE {
		InitRedisClient(addr, password, db)
	} else if mode == REDIS_SENTINEL_MODE {
		addrs := strings.Split(addr, REDIS_ADDRS_SEPARATOR)
		InitRedisSentinel(master, addrs, password, db)
	} else if mode == REDIS_CLUSTER_MODE {
		addrs := strings.Split(addr, REDIS_ADDRS_SEPARATOR)
		InitRedisCluster(addrs, password)
	}
}

func InitRedisClient(addr string, password string, db int) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err := Redis.Ping().Result(); err != nil {
		logrus.Error(err)
	}
}

func InitRedisSentinel(master string, addrs []string, password string, db int) {
	Redis = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    master,
		SentinelAddrs: addrs,
		Password:      password,
		DB:            db,
	})
}

func InitRedisCluster(addrs []string, password string) {
	RedisCluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})
	if _, err := RedisCluster.Ping().Result(); err != nil {
		logrus.Error(err)
	}
}
