package req

/**
 * @note: 付费用户绑定企微搜索条件
 * @auth: tongWz
 * @date: 2022年6月14日09:51:54
**/
type CompanyUserWhereReq struct {
	IsBind             int    `json:"is_bind"`
	CustomerWechatNick string `json:"customer_wechat_nick"`
}
