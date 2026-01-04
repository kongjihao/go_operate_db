package main

import (
	"fmt"

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

// doDemo rdb.Do 方法使用示例:rdb.Do 是一个通用方法，可以执行任意 Redis 命令。
func doDemo() {
	// 直接执行命令获取错误
	err := rdb.Do("SET", "key3", 100, "EX", 3600).Err() // 设置key3，过期时间3600秒，.Err() 方法用于检查命令执行是否出错。
	if err != nil {
		fmt.Printf("set key3 err: %v\n", err)
		panic(err)
	}

	// 执行命令获取结果
	val, err := rdb.Do("GET", "key3").Result()
	if err != nil {
		fmt.Printf("get key3 err: %v\n", err)
		panic(err)
	}
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
