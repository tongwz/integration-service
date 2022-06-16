package mysql

type CompanyUser struct {
	ID                 int    `gorm:"primary_key" json:"id"`
	CustomerId         int    `gorm:"customer_id"  json:"customer_id"`
	CustomerCid        string `gorm:"customer_cid" json:"customer_cid"`
	CustomerWechatNick string `gorm:"customer_wechat_nick" json:"customer_wechat_nick"`
	IsBind             string `gorm:"is_bind" json:"is_bind"`
}

func (cu CompanyUser) TableName() string {
	return "company_user"
}