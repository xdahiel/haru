package common

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

func InitRDB() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetRDB() *redis.Client {
	return rdb
}
