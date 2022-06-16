package hydrus

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	req2 "integration.service/http/req"
	"integration.service/models/mysql"
	"integration.service/pkg/db"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
)

type RepairWechatBindBusiness struct {
	CompanyUserDao *mysql.CompanyUserDao
	CompanyWechatRelationDao *mysql.CompanyWechatRelationDao
	CompanyWechatExternalContactDao *mysql.CompanyWechatExternalContactDao
}

/**
 * @note: 修复企微绑定 互客 数据表
 * @auth: tongWz
 * @date: 2022年6月13日19:51:12
**/
func (robb *RepairWechatBindBusiness) DoIt() {
	/**
	1	查询出is_bind = 1 并且 customer_wechat_nick = ""的数据
	2	通过查到的 customer_cid 查询 company_wechat_relation 表数据
	3 	查询不到不进行处理，查询到了我们通过unionid 查询 company_module库中 的unionid的客户数据
	4 	将name 数据更新 customer_wechat_nick 字段
	 */
	hydrusDbClient := new(db.MysqlCommonDb).NewMysqlCommonDb("utf8mb4", "mysql_hydrus")
	// company_module库客户端
	cmDbClient := new(db.MysqlCommonDb).NewMysqlCommonDb("utf8mb4", "mysql_hz")

	// step 1
	var companyUserList []mysql.CompanyUser
	req := req2.CompanyUserWhereReq{
		IsBind:1,
		CustomerWechatNick:"",
	}
	companyUserList, err := robb.CompanyUserDao.UserList(hydrusDbClient, req)
	if err != nil {
		logging.Error("查询数据失败：", err.Error())
		return
	}

	var relationInfo mysql.CompanyWechatRelation
	var companyId = setting.Cfg.Section("repair_company_user").Key("company_id").MustInt(3)
	for _, info := range companyUserList {
		// step 2
		relationInfo, err = robb.CompanyWechatRelationDao.GetInfo(hydrusDbClient, info.CustomerCid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logging.Error("企微绑定关系cid查询为空：", err.Error(), "cid =", info.CustomerCid)
				continue
			}
			logging.Error("企微绑定关系查询失败：", err.Error())
			continue
		}
		// step 3 通过unionid获取客户姓名
		contactInfo, err := robb.CompanyWechatExternalContactDao.GetInfo(cmDbClient, relationInfo.Unionid, companyId)
		if err != nil {
			logging.Error("企微外部联系人查询失败：", err.Error())
			continue
		}
		upInfo := mysql.CompanyUser{
			CustomerWechatNick:contactInfo.Name,
		}
		err = robb.CompanyUserDao.UpdateWechatNick(hydrusDbClient, info.ID, upInfo)
		if err !=nil {
			logging.Error("更新company_user企微昵称失败：", err.Error())
		}
	}
	fmt.Println("执行完成！")
}