package mysql

import "integration.service/pkg/utils"

type CompanyWechatExternalContact struct {
	ID        int             `gorm:"primary_key" json:"id"`
	Name      string          `gorm:"name"  json:"name"`
	Remark    string          `gorm:"remark" json:"remark"`
	Unionid   string          `gorm:"unionid" json:"unionid"`
	CompanyId int             `gorm:"company_id" json:"company_id"`
	CreatedAt utils.LocalTime `gorm:"created_at" json:"created_at"`
}

func (cwec CompanyWechatExternalContact) TableName() string {
	return "cm_company_wechat_external_contact"
}
