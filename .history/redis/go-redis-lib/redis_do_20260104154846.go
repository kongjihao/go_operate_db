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

// doDemo rdb.Do 方法使用示例
func doDemo() {
	// 执行命令获取结果
	val, err := rdb.Do("GET", "key1").Result()
	fmt.Printf("key1: %v, err: %v\n", val, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接执行命令获取错误
	err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
	fmt.Println(err)

	// 执行命令获取结果
	val, err := rdb.Do(ctx, "get", "key").Result()
	fmt.Println(val, err)
}