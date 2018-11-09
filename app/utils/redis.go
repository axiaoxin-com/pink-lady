package utils

import (
	"strings"

	"github.com/alicebob/miniredis"
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

// InitRedis init redis client by different mode
// `mode` is the redis running mode, single-instance, sentinel or cluster
// single-instance and sentinel mode init the global var which named `Redis`, cluster mode init the gloabel var which named `RedisCluster`
// `addr` is redis address. when mode is sentinel or cluster, the addr is a mutilple address separated by comman
// `password` is redis auth password
// `db` is redis db number, cluster mode don't use it
// `master` is redis sentinel master name, only need to be set on sentinel mode, others dont't use it
func InitRedis(mode string, addr string, password string, db int, master string) error {
	mode = strings.ToLower(mode)
	var err error
	if mode == REDIS_SINGLE_INSTANCE_MODE {
		err = InitRedisClient(addr, password, db)
	} else if mode == REDIS_SENTINEL_MODE {
		addrs := strings.Split(addr, REDIS_ADDRS_SEPARATOR)
		err = InitRedisSentinel(master, addrs, password, db)
	} else if mode == REDIS_CLUSTER_MODE {
		addrs := strings.Split(addr, REDIS_ADDRS_SEPARATOR)
		err = InitRedisCluster(addrs, password)
	}
	return err
}

// InitRedisClient init a single instance redis client named `Redis`
func InitRedisClient(addr string, password string, db int) error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Error(err)
	}
	return err
}

// InitRedisSentinel init redis sentinel client also named `Redis`
func InitRedisSentinel(master string, addrs []string, password string, db int) error {
	Redis = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    master,
		SentinelAddrs: addrs,
		Password:      password,
		DB:            db,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Error(err)
	}
	return err
}

// InitRedisCluster init redis cluster client named `RedisCluster`
func InitRedisCluster(addrs []string, password string) error {
	RedisCluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})
	_, err := RedisCluster.Ping().Result()
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func MockRedis() (*miniredis.Miniredis, error) {
	s, err := miniredis.Run()
	if err != nil {
		logrus.Error(err)
	}
	return s, err
}
