package main

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var rdb *redis.Client
var ctx = context.Background()

func initRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis连接地址
		Password: "",               // 密码，本地redis默认没有密码
		DB:       0,                // use redis defalt db， 一共0~15
		PoolSize: 100,              // 连接池大小
	})

	_, err = rdb.Ping(ctx).Result()
	return err
}

// 在Redis中可以使用KEYS prefix* 命令按前缀查询所有符合条件的 key，go-redis库中提供了Keys方法实现类似查询key的功能。