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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	_, err = rdb.Ping(ctx).Result()
	return err
}

// Redis 是单线程执行命令的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行，比如不同客户端的指令在redis服务器端交叉执行。
// 有时我们需要确保一组命令作为一个整体来执行，中间不被其他客户端的命令插入，这时就需要使用事务（Transaction）。
// 在 go-redis 库中，可以使用 TxPipeline 或 TxPipelined 方法来实现事务操作，下面是一个使用事务的示例代码。
func redisTransactionDemo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "trans_key1", "value1", 0)
		pipe.Set(ctx, "trans_key2", "value2", 0)
		pipe.Incr(ctx, "trans_key1")
		return nil
	})
	return err
}

func main() {
	if err := initRedis(); err != nil {
		println("init Redis client failed, err:", err.Error())
		return
	}
	println("connect redis success")
	defer rdb.Close()

	if err := redisTransactionDemo(); err != nil {
		println("redis transaction demo failed, err:", err.Error())
		return
	}
	println("redis transaction demo success")
}
