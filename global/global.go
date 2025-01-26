package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB      //基于gorm的关系型数据库
	RedisDB *redis.Client //基于redis的键值数据库
)
