package utils

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// redis connection modes
const (
	RedisSingleInstanceMode = "normal"
	RedisSentinelMode       = "sentinel"
	RedisClusterMode        = "cluster"
	RedisAddrsSeparator     = ","
)

// redis clients map
var (
	// RedisClientMap redis normal and sentinel client map
	RedisClientMap = make(map[string]*redis.Client)
	// RedisClusterMap redis cluster client map
	RedisClusterMap = make(map[string]*redis.ClusterClient)
)

// InitRedis init redis client by different mode
// `mode` is the redis running mode, single-instance, sentinel or cluster
// single-instance and sentinel mode init the global var which named `Redis`, cluster mode init the gloabel var which named `RedisCluster`
// `addr` is redis address. when mode is sentinel or cluster, the addr is a mutilple address separated by comman
// `password` is redis auth password
// `dbindex` is redis db number, cluster mode don't use it
// `master` is redis sentinel master name, only need to be set on sentinel mode, others dont't use it
func InitRedis() error {
	var err error
	var cli *redis.Client
	var cluster *redis.ClusterClient

	redisMap := viper.GetStringMap("redis")
	for name, confItf := range redisMap {
		conf := confItf.(map[string]interface{})
		addr := conf["address"].(string)
		password := conf["password"].(string)
		dbindex := int(conf["dbindex"].(int64))
		mode := strings.ToLower(conf["mode"].(string))

		if mode == RedisSingleInstanceMode {
			cli, err = NewRedisClient(addr, password, dbindex)
			RedisClientMap[name] = cli
		} else if mode == RedisSentinelMode {
			addrs := strings.Split(addr, RedisAddrsSeparator)
			cli, err = NewRedisSentinel(viper.GetString("redis.master"), addrs, password, dbindex)
			RedisClientMap[name] = cli
		} else if mode == RedisClusterMode {
			addrs := strings.Split(addr, RedisAddrsSeparator)
			cluster, err = NewRedisCluster(addrs, password)
			RedisClusterMap[name] = cluster
		} else {
			err = errors.New("Invalid Redis config to InitRedis")
		}
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
