package utils

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// redis connection modes
const (
	RedisSingleInstanceMode = "single-instance"
	RedisSentinelMode       = "sentinel"
	RedisClusterMode        = "cluster"
	RedisAddrsSeparator     = ","
)

// redis clients
var (
	// Redis single instance and sentinel client
	Redis *redis.Client
	// RedisCluster redis cluster client
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
	if mode == RedisSingleInstanceMode {
		err = InitRedisClient(addr, password, db)
	} else if mode == RedisSentinelMode {
		addrs := strings.Split(addr, RedisAddrsSeparator)
		err = InitRedisSentinel(master, addrs, password, db)
	} else if mode == RedisClusterMode {
		addrs := strings.Split(addr, RedisAddrsSeparator)
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
	return errors.Wrap(err, "init redis client error")
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
	return errors.Wrap(err, "init redis sentinel error")
}

// InitRedisCluster init redis cluster client named `RedisCluster`
func InitRedisCluster(addrs []string, password string) error {
	RedisCluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})
	_, err := RedisCluster.Ping().Result()
	return errors.Wrap(err, "init redis cluster error")
}
