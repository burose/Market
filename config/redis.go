package config

import (
	"github.com/go-redis/redis"
	"log"
	"market/global"
)

func initRedisConfig() {
	addr := Appconfig.Redis.Addr
	db := Appconfig.Redis.DB
	password := Appconfig.Redis.Password

	redisclint := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := redisclint.Ping().Result()
	if err != nil {
		log.Fatalf("Fail to connected to redis,got error: %v", err)
	}
	global.RedisDB = redisclint
}
