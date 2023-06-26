package utils

import (
	"gameServers/config"
	"github.com/go-redis/redis"
	"log"
	"sync"
)

var redisClient *redis.Client
var mtx sync.Mutex

var RedisServer = config.Config.GetString("RedisServer")
var Address = config.Config.GetString("Address")

func init() {
	doInit()
}

func doInit() {
	if redisClient != nil {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if redisClient != nil {
		return
	}

	if RedisServer != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     RedisServer,
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		})
	}
}

func GetId() (string, error) {
	doInit()
	id, err := redisClient.Get(Address).Result()
	// 判断查询是否出错
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("id的值：", id)
	return id, nil
}

func SetId(val string) error {
	doInit()
	log.Println("Address", Address, "redisServer:", RedisServer)
	err := redisClient.Set(Address, val, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
