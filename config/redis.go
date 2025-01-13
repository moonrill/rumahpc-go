package config

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))

	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}

	Rdb = redis.NewClient(opt)
}
