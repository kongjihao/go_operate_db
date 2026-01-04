package main

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接执行命令获取错误
	err := rdb.Do(ctx, "SET", "key3", 100, "EX", 3600).Err()
	fmt.Printf("set key3 err: %v\n", err)

	// 执行命令获取结果
	val, err := rdb.Do(ctx, "GET", "key3").Result()
	fmt.Printf("get key3 value: %v, err: %v\n", val, err)
}

func main() {
	if err := initRedis(); err != nil {
		fmt.Println("connect redis failed!")
		return
	}

	fmt.Println("connect redis success")

	doDemo()
}
