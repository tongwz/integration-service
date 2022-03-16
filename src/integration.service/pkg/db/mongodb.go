package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"net/url"
	"time"
)
const cMongodbHeartbeat = 60
var MongodbClient *mongo.Client

func init() {
	MongodbClient, _ = getMongoClient()

	err := MongodbClient.Ping(context.TODO(), nil)
	if err != nil {
		logging.Error("连接mongodb数据库失败：", err.Error())
	}

	go mongodbHeartbeatRun()
}

/**
 * @note: 获取mongodb的连接
 * @auth: tongwz
 * @date 2022年2月18日17:49:02
**/
func getMongoClient() (*mongo.Client, error) {
	config, err := setting.Cfg.GetSection("mongodb")
	if err != nil {
		panic("app.ini中数据库配置database未设置")
	}
	var connStr string
	if config.HasKey("user") {
		connStr = "mongodb://%s:%s@%s:%s/?connect=direct"
		connStr = fmt.Sprintf(
			connStr,
			url.QueryEscape(config.Key("user").String()),
			url.QueryEscape(config.Key("password").String()),
			config.Key("host").String(),
			config.Key("port").String(),
		)
	} else {
		connStr = "mongodb://%s:%s/?connect=direct"
		connStr = fmt.Sprintf(
			connStr,
			config.Key("host").String(),
			config.Key("port").String(),
		)
	}
	// 连接选项
	clientOptions := options.Client().ApplyURI(connStr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	// 连接指针
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logging.Fatal("mongodb 创建连接~", err.Error())
		return nil, err
	}
	err = client.Connect(ctx)
	//client.Database(config.Key("name").String())
	if err != nil {
		logging.Fatal("mongodb 连接失败~", err.Error())
		return nil, err
	}
	return client, nil
}

// 重连mongodb
func resetMongodbClient() {
	var err error

	if MongodbClient != nil {
		_ = MongodbClient.Disconnect(context.TODO())
		MongodbClient, err = getMongoClient()
	}

	if err != nil {
		logging.Fatal("mongodb 初始化打开失败!%s", err.Error())
	}
}

// mongodb心跳包
func mongodbHeartbeatRun() {
	timer := time.NewTicker(cMongodbHeartbeat * time.Second)
	for {
		select {
		case <-timer.C: // 每一分钟 ping 一次
			if err := MongodbClient.Ping(context.TODO(), nil); err != nil {
				logging.Warn("mongodb无法 ping 通!%s", err.Error())
				resetMongodbClient()
			}
		}
	}
}
