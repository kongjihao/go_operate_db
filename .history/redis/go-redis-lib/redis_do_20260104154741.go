package main

import (
	"fmt"
	"time"
	
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func initRedis() (err error)