package main

import (
	"fmt"
	"time"
	"context"
	
	"github.com/redis/go-redis/v9"
	
)

var rdb *redis.Client
var ctx = context.Background()

func initRedis() (err error)