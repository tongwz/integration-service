package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"strings"
	"time"
)

const cRedisHeartbeat = 60
var Redis *redis.Client

func init() {
	Redis = getRedisClient()
	if _, err := Redis.Ping().Result(); err != nil {
		fmt.Println("redis连初始化错误：", err.Error())
		logging.Error("连接redis数据库失败：", err.Error())
	}
	go redisHeartbeatRun()
}

func getRedisClient() *redis.Client {
	config, err := setting.Cfg.GetSection("redis")
	if err != nil {
		panic("app.ini中数据库配置redis未设置")
	}

	return redis.NewClient(&redis.Options{
		Addr:     config.Key("host").String(),
		Password: config.Key("password").String(), // 没有password也得有个空值
		DB:       config.Key("db").MustInt(),
	})

}

/**
 * @note: 重启redis
 * @auth: tongwz
 * @date  2022年2月23日19:21:10
**/
func resetRedis() {
	if Redis != nil {
		_ = Redis.Close()
	}

	Redis = getRedisClient()

	if _, err := Redis.Ping().Result(); err != nil {
		fmt.Println("redis连初始化错误：", err.Error())
		logging.Error("redis初始化打开失败:%s", err.Error())
	}
}

/**
 * @note: redis的心跳包
 * @auth: tongwz
 * @date  2022年2月23日19:21:26
**/
func redisHeartbeatRun() {
	timer := time.NewTicker(cRedisHeartbeat * time.Second)
	for {
		select {
		case <-timer.C: // 每一分钟 ping 一次
			if pong, err := Redis.Ping().Result(); err != nil || strings.ToLower(pong) != "pong" {
				resetRedis()
				logging.Error("心跳检测 redis 客户端有误:%+v;pong:%s", err, pong)
			}
		}
	}
}
