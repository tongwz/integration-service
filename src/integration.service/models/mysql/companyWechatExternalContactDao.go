package mysql

import "github.com/jinzhu/gorm"

type CompanyWechatExternalContactDao struct {

}

/**
 * @note: 查询客户详情
 * @auth: tongWz
 * @date: 2022年6月14日11:07:20
**/
func (cwecd *CompanyWechatExternalContactDao) GetInfo(db *gorm.DB, uid string, companyId int) (res CompanyWechatExternalContact, err error) {
	err = db.Select("*").
		Where("unionid = ? ", uid).
		Where("company_id = ?", companyId).
		First(&res).Error
	return res, err
}
