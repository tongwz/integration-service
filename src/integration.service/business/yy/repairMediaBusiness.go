package yy

import (
	"fmt"
	"integration.service/models/mysql"
	"integration.service/pkg/db"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"integration.service/utils"
	"strings"
)

type RepairMedia struct {
	WechatSessionHistoryDao mysql.WechatSessionHistoryDao
	TableNameMap            map[string]string
}

/**
 * @note: 修复媒体文件 typeChange 例如 msg_create_at msg_id  index 例如2022-05-19 16290973855815008786_1652189911210_external
 * @auth: tongWz
 * @date: 2022年5月24日11:18:52
**/
func (rm *RepairMedia) DoIt(typeChange string, index string) {
	dbYY := new(db.MysqlCommonDb).NewMysqlCommonDb("utf8mb4", "mysql_yy")
	dbHZ := new(db.MysqlCommonDb).NewMysqlCommonDb("utf8mb4", "mysql_hz")
	dbIndex := setting.Cfg.Section("repair_media").Key("need_repair_table_name_index").MustString("20225")
	yyTableName := mysql.WechatSessionHistory{}.TableName(rm.TableNameMap[dbIndex])
	hzTableName := mysql.WechatSessionHistory{}.TableName(dbIndex)
	var yySessionData, hzSessionData mysql.WechatSessionHistory
	if typeChange == "msg_id" {
		// 查询盈亚历史数据
		err := dbYY.Table(yyTableName).Where("`msg_id` = ? ", index).First(&yySessionData).Error
		if err != nil {
			logging.Error("通过msg_id查询盈亚mysql数据报错", err.Error())
			return
		}
		yyImageUrl := setting.Cfg.Section("image").Key("yy_image").MustString("https://image-server-api.gp622.com") + yySessionData.Content

		// 文件临时地址
		tmpFilePath := setting.Cfg.Section("repair_media").Key("tmp_file_path").MustString("runtime/temp") + "/" + yySessionData.MsgId + "." +
			setting.Cfg.Section("repair_media").Key("file_ext").MustString("mp4")
		// 测试生产环境下载
		// yyImageUrl = "https://image-server-api.gp622.com/api/attachment/show?id=627a6eeb3048ff6a7e1b688e"

		// 下载链接文件
		err = utils.DownFileUrl(tmpFilePath, yyImageUrl)
		if err != nil {
			logging.Error("读取盈亚的图片链接失败：", err.Error())
			return
		}

		// 文件上传到图片服务器
		fileUrl, _, err := utils.SyncToImageServer(tmpFilePath, "")
		if err != nil {
			return
		}

		if strings.Trim(fileUrl, " ") == "" {
			logging.Error("会话文件获取失败或上传失败fileUrl：", fileUrl)
			return
		}
		// 更新数据
		hzSessionData = mysql.WechatSessionHistory{
			Content: fileUrl,
		}
		// 将文件链接进行重新保存到杭州数据库
		err = dbHZ.Table(hzTableName).Where("`msg_id` = ? ", index).Update(&hzSessionData).Error
		if err != nil {
			logging.Error("语音会话链接更新失败：", err.Error())
			return
		}
		logging.Info("语音会话链接更新成功链接是：", setting.Cfg.Section("image").Key("hz_image").MustString("https://image-server-api.zq332.com"), fileUrl)
	} else if typeChange == "msg_create_at" {
		var yySessionDataBatch []mysql.WechatSessionHistory
		// 查询盈亚历史数据
		err := dbYY.Table(yyTableName).Where("`msg_create_at` = ? ", index).
			Where("`msg_type`= ? ", "meeting_voice_call").
			Find(&yySessionDataBatch).Error
		if err != nil {
			logging.Error("通过msg_create_at查询盈亚mysql数据报错", err.Error())
			return
		}
		if len(yySessionDataBatch) == 0 {
			logging.Info("通过msg_create_at查询盈亚数据为空")
			return
		}
		for _, yySessionData = range yySessionDataBatch {
			yyImageUrl := setting.Cfg.Section("image").Key("yy_image").MustString("https://image-server-api.gp622.com") + yySessionData.Content

			// 文件临时地址
			tmpFilePath := setting.Cfg.Section("repair_media").Key("tmp_file_path").MustString("runtime/temp") + "/" + yySessionData.MsgId + "." +
				setting.Cfg.Section("repair_media").Key("file_ext").MustString("mp4")

			// 下载链接文件
			err = utils.DownFileUrl(tmpFilePath, yyImageUrl)
			if err != nil {
				logging.Error("读取盈亚的图片链接失败：", err.Error())
				return
			}

			// 文件上传到图片服务器
			fileUrl, _, err := utils.SyncToImageServer(tmpFilePath, "")
			if err != nil {
				return
			}

			if strings.Trim(fileUrl, " ") == "" {
				logging.Error("会话文件获取失败或上传失败fileUrl：", fileUrl)
				return
			}
			// 更新数据
			hzSessionData = mysql.WechatSessionHistory{
				Content: fileUrl,
			}
			// 将文件链接进行重新保存到杭州数据库
			err = dbHZ.Table(hzTableName).Where("`msg_id` = ? ", yySessionData.MsgId).Update(&hzSessionData).Error
			if err != nil {
				logging.Error("语音会话链接更新失败：", err.Error())
				return
			}
			logging.Info("语音会话链接更新成功链接是：", setting.Cfg.Section("image").Key("hz_image").MustString("https://image-server-api.zq332.com"), fileUrl)
		}

	} else {
		fmt.Printf("无法识别typeChange = %s, index = %s", typeChange, index)
	}
}
