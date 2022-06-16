package rsp

// 图片服务的响应
type ImageServerResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}