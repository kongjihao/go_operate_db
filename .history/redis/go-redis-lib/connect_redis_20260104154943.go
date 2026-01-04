// package main

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

// doCommand go-redis基本使用示例
func doCommand() {
	// 执行命令获取结果
	val, err := rdb.Get("key1").Result()
	fmt.Printf("key1: %v, err: %v\n", val, err) // 获取值和错误，没有错误返回nil

	// 先获取到命令对象
	cmder := rdb.Get("key1")
	fmt.Printf("key1 value: %v\n", cmder.Val())   // 获取值
	fmt.Printf("get key1 err: %v\n", cmder.Err()) // 获取错误

	// 直接执行命令获取错误
	err = rdb.Set("key2", 2, time.Second*3).Err()
	fmt.Printf("set key2 err: %v\n", err)

	// 直接执行命令获取值
	value := rdb.Get("key2").Val()
	fmt.Printf("key2 value: %v\n", value)
}

func main() {
	if err := initRedis(); err != nil {
		fmt.Printf("init Redis client failed, err: %v\n", err)
		return
	}

	fmt.Println("connect redis success")

	doCommand()
}
