package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"time"
)

const cHeart = 60

type MysqlCommonDb struct {

}

func (mcd *MysqlCommonDb) NewMysqlCommonDb(charset string, dbIndex string) *gorm.DB {
	var err error
	mysqlDb, err := mcd.getMysqlClient(charset, dbIndex)
	if err != nil {
		panic("数据库utf8链接初始化打开失败:" + err.Error())
	}

	err = mysqlDb.DB().Ping()
	if err != nil {
		logging.Error("连接mysql数据库失败：", err.Error())
	}

	// 心跳
	go mcd.mysqlHeartbeatRun(charset, dbIndex, mysqlDb)
	return mysqlDb
}

/**
 * 获取 mysql 客户端
 */
func (mcd *MysqlCommonDb) getMysqlClient(charset string, dbIndex string) (*gorm.DB, error) {
	config, err := setting.Cfg.GetSection(dbIndex)
	if err != nil {
		panic(fmt.Sprintf("app.ini中数据库配置%s未设置", dbIndex))
	}
	protocol := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.Key("user").MustString("admin"),
		config.Key("password").MustString("admin"),
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
func (mcd *MysqlCommonDb) resetMysqlClient(charset string, dbIndex string, mysqlDb *gorm.DB) {
	var err error

	if mysqlDb != nil {
		_ = mysqlDb.Close()
		mysqlDb, err = mcd.getMysqlClient(charset, dbIndex)
	}

	if err != nil {
		logging.Fatal("数据库初始化打开失败!%s", err.Error())
	}
}

/**
 * mysql 心跳
 */
func (mcd *MysqlCommonDb) mysqlHeartbeatRun(charset string, dbIndex string, mysqlDb *gorm.DB) {
	timer := time.NewTicker(cHeart * time.Second)
	for {
		select {
		case <-timer.C: // 每一分钟 ping 一次
			if err := mysqlDb.DB().Ping(); err != nil {
				logging.Warn("mysql 无法 ping 通!%s", err.Error())
				mcd.resetMysqlClient(charset, dbIndex, mysqlDb)
			}
		}
	}
}
