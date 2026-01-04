package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
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

// Redis Pipeline 允许通过使用单个 client-server-client 往返执行多个命令来提高性能。区别于一个接一个地执行100个命令，你可以将这些命令放入 pipeline 中，然后使用1次读写操作像执行单个命令一样执行它们。这样做的好处是节省了执行命令的网络往返时间（RTT）。
// 在下面的示例代码中演示了使用 pipeline 通过一个 write + read 操作来执行多个命令。
func redisPipelineDemo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond) // 设置超时时间
	defer cancel()                                                                 // 释放资源
	pipe := rdb.Pipeline()

	incr1 := pipe.Incr(ctx, "pipeline_counter1")
	incr2 := pipe.Incr(ctx, "pipeline_counter2")
	pipe.Expire(ctx, "pipeline_counter1", time.Hour)
	pipe.Expire(ctx, "pipeline_counter2", time.Hour)

	_, err := pipe.Exec(ctx) // 将四个命令一次性发送到Redis服务器执行，比不使用pipeline减少了3次网络往返时间
	if err != nil {
		return err
	}

	println("pipeline_counter1:", incr1.Val())
	println("pipeline_counter2:", incr2.Val())
	return nil
}

// 或者，你也可以使用Pipelined 方法，它会在函数退出时调用 Exec。


func main() {
	if err := initRedis(); err != nil {
		println("connect redis failed:", err.Error())
		return
	}
	println("connect redis success")
	defer rdb.Close() // 关闭客户端连接

	if err := redisPipelineDemo(); err != nil {
		println("redis pipeline demo failed:", err.Error())
		return
	}
}
