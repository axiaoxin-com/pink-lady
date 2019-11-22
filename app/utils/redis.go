package utils

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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
	// Redis single instance
	Redis *redis.Client
	// Redis sentinel instance
	RedisSentinel *redis.Client
	// RedisCluster redis cluster client
	RedisCluster *redis.ClusterClient
)

// InitRedis init redis client by different mode
// `mode` is the redis running mode, single-instance, sentinel or cluster
// single-instance and sentinel mode init the global var which named `Redis`, cluster mode init the gloabel var which named `RedisCluster`
// `addr` is redis address. when mode is sentinel or cluster, the addr is a mutilple address separated by comman
// `password` is redis auth password
// `dbindex` is redis db number, cluster mode don't use it
// `master` is redis sentinel master name, only need to be set on sentinel mode, others dont't use it
func InitRedis() error {
	addr := viper.GetString("redis.address")
	password := viper.GetString("redis.password")
	dbindex := viper.GetInt("redis.dbindex")
	mode := strings.ToLower(viper.GetString("redis.mode"))

	var err error
	if mode == RedisSingleInstanceMode {
		Redis, err = NewRedisClient(addr, password, dbindex)
	} else if mode == RedisSentinelMode {
		addrs := strings.Split(addr, RedisAddrsSeparator)
		RedisSentinel, err = NewRedisSentinel(viper.GetString("redis.master"), addrs, password, dbindex)
	} else if mode == RedisClusterMode {
		addrs := strings.Split(addr, RedisAddrsSeparator)
		RedisCluster, err = NewRedisCluster(addrs, password)
	}
	return err
}

// NewRedisClient return a single instance redis client
func NewRedisClient(addr string, password string, dbindex int) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbindex,
	})
	_, err := r.Ping().Result()
	return r, errors.Wrap(err, "init redis client error")
}

// NewRedisSentinel return redis sentinel client
func NewRedisSentinel(master string, addrs []string, password string, dbindex int) (*redis.Client, error) {
	r := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    master,
		SentinelAddrs: addrs,
		Password:      password,
		DB:            dbindex,
	})
	_, err := r.Ping().Result()
	return r, errors.Wrap(err, "init redis sentinel error")
}

// NewRedisCluster return redis cluster client
func NewRedisCluster(addrs []string, password string) (*redis.ClusterClient, error) {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})
	_, err := c.Ping().Result()
	return c, errors.Wrap(err, "init redis cluster error")
}
