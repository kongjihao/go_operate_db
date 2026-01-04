package main

import (
	"fmt"

	"github.com/go-redis/redis"	// go list -m all | grep go-redis
	// "github.com/redis/go-redis/v9"
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

	// zset values
	language := []redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}

	// zAdd 添加元素
	err := rdb.ZAdd(zsetKey, language...).Err()
	if err != nil {
		fmt.Println("Error adding elements to zset(language_rank): ", err)
		return
	}
	fmt.Println("Elements added to zset(language_rank) successfully.")

	// ZIncrBy 增加元素的分数
	newScore, err := rdb.ZIncrBy(zsetKey, 5.0, "Golang").Result()
	if err != nil {
		fmt.Println("Error incrementing score for Golang: ", err)
		return
	}
	fmt.Println("New score for Golang: ", newScore)

	// ZRevRangeWithScores 取分数最高的3个
	topLanguages, err := rdb.ZRevRangeWithScores(zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Println("Error retrieving top 3 languages: ", err)
		return
	}
	fmt.Println("Top 3 languages by score:")
	for _, lang := range topLanguages {
		fmt.Printf("Language: %s, Score: %.2f\n", lang.Member, lang.Score)
	}

	// ZRangeByScore 根据分数范围查询元素，取95~100分的
	min := "95"
	max := "100"
	zRangeBy := redis.ZRangeBy{
		Min: min,
		Max: max,
	}
	languageInRange, err := rdb.ZRangeByScore(zsetKey, zRangeBy).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	fmt.Printf("Languages with scores between %s and %s: %v\n", min, max, languageInRange)

	// ZRem 删除元素
	removed, err := rdb.ZRem(zsetKey, "JavaScript").Result()
	if err != nil {
		fmt.Println("Error removing JavaScript from zset: ", err)
		return
	}
	fmt.Printf("Number of elements removed (JavaScript): %d\n", removed)
}

func main() {
	err := initRedis()
	if err != nil {
		fmt.Println("Failed to connect to Redis: ", err)
		return
	}
	defer rdb.Close()

	zsetDemo()
}
