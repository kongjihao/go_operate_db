package main

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var rdb *redis.Client
var ctx = context.Background()

func initRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis连接地址