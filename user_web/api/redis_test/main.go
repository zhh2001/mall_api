package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.120.21.77:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	value, err := rdb.Get(context.Background(), "13866660008").Result()
	if errors.Is(err, redis.Nil) {
		fmt.Println("Key 不存在")
	}
	fmt.Println(value)
}
