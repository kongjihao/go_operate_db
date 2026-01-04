package main

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var rdb *redis.Client
var ctx = context.Background()

func initRedis() (err error)