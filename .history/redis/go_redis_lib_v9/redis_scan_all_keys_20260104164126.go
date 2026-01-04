package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

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
// 但是在生产环境中不建议使用Keys命令，因为它会阻塞Redis服务器，影响性能。推荐使用Scan命令进行遍历查询。
// go-redis库中提供了Scan方法实现遍历查询key的功能，下面是一个使用Scan方法查询所有key的示例。
func scanAllKeys() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond) // 设置超时时间
	defer cancel()                                                                 // 释放资源

	var cursor uint64 // 游标位置，初始值为0，表示从头开始扫描，每次扫描后会更新为新的游标位置，当游标返回 0 时，表示扫描结束
	var keys []string // 存储扫描到的所有key
	for {
		var scanKeys []string // 存储每次扫描到的key
		var err error
		scanKeys, cursor, err = rdb.Scan(ctx, cursor, "*", 10).Result() // 每次扫描10个key，
		if err != nil {
			return nil, err
		}
		keys = append(keys, scanKeys...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// 针对这种需要遍历大量key的场景，go-redis中提供了一个简化方法——Iterator，其使用示例如下。
func scanAllKeysWithIterator() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond) // 设置超时时间
	defer cancel()                                                                 // 释放资源

	var keys []string
	iter := rdb.Scan(ctx, 0, "*", 10).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func main() {
	if err := initRedis(); err != nil {
		panic(err)
	}
	defer rdb.Close()

	keys, err := scanAllKeys()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		println(key)
	}

	keys2, err := scanAllKeysWithIterator()
	if err != nil {
		panic(err)
	}

	for _, key := range keys2 {
		println(key)
	}
}
