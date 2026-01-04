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

// doCommand go-redis基本使用示例
func doCommand() {
ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
defer cancel()

// 执行命令获取结果
val, err := rdb.Get(ctx, "key").Result()
fmt.Println(val, err)

// 先获取到命令对象
cmder := rdb.Get(ctx, "key")
fmt.Println(cmder.Val()) // 获取值
fmt.Println(cmder.Err()) // 获取错误

// 直接执行命令获取错误
err = rdb.Set(ctx, "key", 10, time.Hour).Err()

// 直接执行命令获取值
value := rdb.Get(ctx, "key").Val()
fmt.Println(value)
}


func main() {
	if err := initRedis(); err != nil {
		fmt.Printf("init Redis client failed, err: %v\n", err)
		return
	}

	fmt.Println("connect redis success")
}
