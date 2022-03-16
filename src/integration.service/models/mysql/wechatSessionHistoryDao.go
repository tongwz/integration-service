package mysql

import (
	"github.com/jinzhu/gorm"
	"integration.service/pkg/logging"
)

type WechatSessionHistoryDao struct {
}

/**
 * @note: 查询表格的全部内容，筛选要查的日志
 * @auth: tongwz
 * @date  2022年2月18日09:18:56
**/
func (sh *WechatSessionHistoryDao) SessionList(tx *gorm.DB, suffix string, msgCreateAt string, pageSize int, page int) ([]WechatSessionHistory, error) {
	tableName := WechatSessionHistory.TableName(WechatSessionHistory{}, suffix)
	var wechatSessionList []WechatSessionHistory
	offset := pageSize * (page - 1)
	err := tx.Select("id,company_id,msg_type,content,msg_create_at").
		Table(tableName).
		Where("msg_create_at = ?", msgCreateAt).
		Offset(offset).
		Limit(pageSize).
		Find(&wechatSessionList).
		Error
	if err != nil {
		logging.Error("查询单独聊天记录失败：原因是", err.Error())
	}
	return wechatSessionList, err

}
