package console

import (
	"github.com/spf13/cobra"
	"integration.service/business/hydrus"
)

type RepairWechatBind struct {
}

func (rmd *RepairWechatBind) Init() {

}

func (rmd *RepairWechatBind) GetShort() string {
	return "修复company_user企微绑定企微昵称"
}

func (rmd *RepairWechatBind) GetUse() string {
	return "repairWechatBind"
}

func (rmd *RepairWechatBind) GetArgs() cobra.PositionalArgs {
	return cobra.ExactArgs(0)
}

func (rmd *RepairWechatBind) SetFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

/**
 * @note: 修复语音视频通话消息 将广州的数据库中数据 下载下来并且转移到 我们自己的图片服务器
 * @auth: tongwz
 * @date  2022年5月19日15:16:40
**/
func (rmd *RepairWechatBind) Run(cmd *cobra.Command, args []string) {
	var repairBindBusiness = new(hydrus.RepairWechatBindBusiness)
	repairBindBusiness.DoIt()
}
