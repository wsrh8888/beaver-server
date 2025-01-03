package core

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

func InitRedis(addr string, password string, db int) (client *redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 100,
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		logx.Error(err)
		panic(err)
	}
	fmt.Println("redis链接成功")

	return rdb
}
