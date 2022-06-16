package mysql

import "github.com/jinzhu/gorm"

type CompanyWechatRelationDao struct {

}

/**
 * @note: 通过cid查询数据
 * @auth: tongWz
 * @date: 2022年6月14日10:27:20
**/
func (card *CompanyWechatRelationDao) GetInfo(db *gorm.DB, cid string) (CompanyWechatRelation, error){
	var res CompanyWechatRelation
	err := db.Select("*").
		Where("cid = ?", cid).
		First(&res).Error
	return res, err
}
