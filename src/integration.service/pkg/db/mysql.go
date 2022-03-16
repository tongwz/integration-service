package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"net/url"
	"time"
)

const cMysqlHeartbeat = 60

var Mysql *gorm.DB
var MysqlMb4 *gorm.DB

func init() {
	var err error
	Mysql, err = getMysqlClient("utf8mb4")
	if err != nil {
		panic("数据库utf8链接初始化打开失败:" + err.Error())
	}

	MysqlMb4, err = getMysqlClient("utf8mb4")
	if err != nil {
		panic("数据库utf8mb4链接初始化打开失败:" + err.Error())
	}

	_ = Mysql.DB().Ping()
	err = MysqlMb4.DB().Ping()
	if err != nil {
		logging.Error("连接mysql数据库失败：", err.Error())
	}

	// 心跳
	go mysqlHeartbeatRun()
}

/**
 * 获取 mysql 客户端
 */
func getMysqlClient(charset string) (*gorm.DB, error) {
	config, err := setting.Cfg.GetSection("database")
	if err != nil {
		panic("app.ini中数据库配置database未设置")
	}
	protocol := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		url.QueryEscape(config.Key("user").MustString("admin")),
		url.QueryEscape(config.Key("password").MustString("admin")),
		config.Key("host").String(),
		config.Key("name").String(),
		charset,
	)

	gDb, err := gorm.Open("mysql", protocol)
	if err != nil {
		return nil, err
	}

	// 加入sql写入日志
	if config.Key("is_log_mode").String() == "true" {
		gDb.LogMode(true)
	}

	return gDb, nil

}

/**
 * 重置 mysql 客户端
 */
func resetMysqlClient(charset string) {
	var err error

	if Mysql != nil && charset == "utf8" {
		_ = Mysql.Close()
		Mysql, err = getMysqlClient(charset)
	}

	if MysqlMb4 != nil && charset == "utf8mb4" {
		_ = MysqlMb4.Close()
		MysqlMb4, err = getMysqlClient(charset)
	}

	if err != nil {
		logging.Fatal("数据库初始化打开失败!%s", err.Error())
	}
}

/**
 * mysql 心跳
 */
func mysqlHeartbeatRun() {
	timer := time.NewTicker(cMysqlHeartbeat * time.Second)
	for {
		select {
		case <-timer.C: // 每一分钟 ping 一次
			if err := Mysql.DB().Ping(); err != nil {
				logging.Warn("mysql utf8 无法 ping 通!%s", err.Error())
				resetMysqlClient("utf8")
			}
			if err := MysqlMb4.DB().Ping(); err != nil {
				logging.Warn("mysql utf8mb4 无法 ping 通!%s", err.Error())
				resetMysqlClient("utf8mb4")
			}
		}
	}
}
