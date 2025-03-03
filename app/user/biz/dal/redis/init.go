package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/trashwbin/dymall/app/user/conf"
)

var (
	RedisClient *redis.Client
)

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: conf.GetConf().Redis.Address,
		//Username: conf.GetConf().Redis.Username,//TODO redis版本可能不一样 后续调整
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
