package main

import (
	"github.com/go-redis/redis"
	"golang.org/x/text/language"
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
// ZSET（Sorted Set，有序集合） 是其最强大、最具特色的数据结构之一。它完美地结合了 Set（集合） 和 Hash（哈希） 的特性，并提供了自动排序的功能。

// 核心特性:
// 你可以把它想象成一个排行榜或者一个带权重的唯一成员集合。
// 唯一性（Set 特性）： 每个成员（Member）都是唯一的，不能重复。
// 有序性（Sorted 特性）： 每个成员都关联一个浮点数类型的 分数（Score）。集合中的成员根据这个分数从小到大进行升序排列（分数小的在前，大的在后）。
// 高性能操作： 基于跳跃表（Skip List） 和哈希表的混合实现，保证了核心操作（增、删、查、按排名/分数范围查询）的高效性。
func zsetDemo() {
	// zset key
	zsetKey := "language_rank"

	// ZAdd 向有序集合中添加元素
	language := []redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}
	
	rdb.ZAdd(zsetKey, language...)
	
	// ZIncrBy 增加元素的分数
	rdb.ZIncrBy(zsetKey, 5.0, "Golang")

	// ZRangeByScore 根据分数范围查询元素
