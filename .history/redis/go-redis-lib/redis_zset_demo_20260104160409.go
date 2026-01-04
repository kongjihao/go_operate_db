package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func initRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis连接地址
		Password: "",               // 密码，本地redis默认没有密码
		DB:       0,                // use redis defalt db， 一共0~15
		PoolSize: 100,              // 连接池大小
	})

	_, err = rdb.Ping().Result()
	return err
}

// zsetDemo
func zsetDemo() {
	zsetKey := "language_rank"

	// ZAdd 向有序集合中添加元素
	rdb.ZAdd(zsetKey, &redis.Z{Score: 90, Member: "Golang"})
	rdb.ZAdd(zsetKey, &redis.Z{Score: 98, Member: "Java"})
	rdb.ZAdd(zsetKey, &redis.Z{Score: 95, Member: "Python"})
	rdb.ZAdd(zsetKey, &redis.Z{Score: 97, Member: "JavaScript"})
}