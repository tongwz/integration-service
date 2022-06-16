package console

import (
	"fmt"
	"github.com/spf13/cobra"
	"integration.service/business/yy"
	"integration.service/models/mysql"
	"integration.service/pkg/setting"
)

type RepairMediaData struct {
}

func (rmd *RepairMediaData) Init() {

}

func (rmd *RepairMediaData) GetShort() string {
	return "修复语音通话数据"
}

func (rmd *RepairMediaData) GetUse() string {
	return "repairMediaData"
}

func (rmd *RepairMediaData) GetArgs() cobra.PositionalArgs {
	return cobra.ExactArgs(2)
}

func (rmd *RepairMediaData) SetFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

/**
 * @note: 修复语音视频通话消息 将广州的数据库中数据 下载下来并且转移到 我们自己的图片服务器
 * @auth: tongwz
 * @date  2022年5月19日15:16:40
**/
func (rmd *RepairMediaData) Run(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("需要传入两个参数,修改数据类型 和 目标值")
		return
	}
	// typeChange 有 msg_id , msg_create_at
	var typeChange = args[0]
	var cIndex = args[1]
	var tableNameMap = mysql.TableNameMap
	env := setting.Cfg.Section("app").Key("env").MustString("test")

	if env == "test" {
		tableNameMap = mysql.TableNameTestMap
	}

	var delImageFiles = yy.RepairMedia{TableNameMap:tableNameMap}
	delImageFiles.DoIt(typeChange, cIndex)
}