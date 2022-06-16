package mysql

import (
	"github.com/jinzhu/gorm"
	"integration.service/http/req"
	"integration.service/pkg/logging"
)

type CompanyUserDao struct {

}

/**
 * @note: 查看付费订单用户企微绑定列表
 * @auth: tongWz
 * @date: 2022年6月14日10:02:14
**/
func (cud *CompanyUserDao) UserList(tx *gorm.DB, where req.CompanyUserWhereReq) ([]CompanyUser, error) {
	var companyUserList []CompanyUser
	err := tx.Select("id,customer_id,customer_cid,customer_wechat_nick,is_bind").
		Table(CompanyUser{}.TableName()).
		Where("is_bind = ?", where.IsBind).
		Where("customer_wechat_nick = ?", where.CustomerWechatNick).
		Find(&companyUserList).Error
	if err != nil {
		logging.Error("付费订单用户企微绑定列表失败", err.Error())
	}
	return companyUserList, err
}

/**
 * @note: 更新企微昵称字段
 * @auth: tongWz
 * @date: 2022年6月14日11:23:12
**/
func (cud *CompanyUserDao) UpdateWechatNick(tx *gorm.DB, id int, upInfo CompanyUser) error {
	return tx.Table(CompanyUser{}.TableName()).Where("id = ?", id).Update(&upInfo).Error
}
