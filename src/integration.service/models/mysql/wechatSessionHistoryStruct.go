package mysql

import "integration.service/pkg/setting"

type WechatSessionHistory struct {
	ID                int    `gorm:"primary_key" json:"id"`
	WechatUserId      int    `gorm:"wechat_user_id" json:"wechat_user_id"`
	ExternalContactId int    `gorm:"external_contact_id" json:"external_contact_id"`
	MsgId             string `gorm:"msg_id" json:"msg_id"`
	MsgCreateAt       string `gorm:"msg_create_at" json:"msg_create_at"`
	MsgCreateTime     int64  `gorm:"msg_create_time" json:"msg_create_time"`
	CompanyId         int64  `gorm:"company_id" json:"company_id"`
	MsgType           string `gorm:"msg_type" json:"msg_type"`
	Content           string `gorm:"content" json:"content"`
	JsonContent       string `gorm:"json_content" json:"json_content"`
	AddTime           string `gorm:"add_time" json:"add_time"`
	Name              string `gorm:"name" json:"name"`
	RevokeMsgId       int64  `gorm:"revoke_msg_id" json:"revoke_msg_id"`
	RevokeTime        int64  `gorm:"revoke_time" json:"revoke_time"`
}

func (ws WechatSessionHistory) TableName(suffix string) string {
	config, _ := setting.Cfg.GetSection("database")
	return config.Key("table_prefix").MustString("cm_") + "company_wechat_session_history_" + suffix
}
