package mysql

import "integration.service/pkg/utils"

type CompanyWechatRelation struct {
	ID        int             `gorm:"primary_key" json:"id"`
	UserId    string          `gorm:"user_id"  json:"user_id"`
	Accid     string          `gorm:"accid" json:"accid"`
	Unionid   string          `gorm:"unionid" json:"unionid"`
	Cid       string          `gorm:"cid" json:"cid"`
	AccountId int64           `gorm:"accountId" json:"accountId"`
	CreatedAt utils.LocalTime `gorm:"created_at" json:"created_at"`
	IsDelete  int             `gorm:"is_delete" json:"is_delete"`
}

func (cu CompanyWechatRelation) TableName() string {
	return "company_wechat_relation"
}
