package console

import (
	"fmt"
	"github.com/spf13/cobra"
	"integration.service/business/jz"
	"integration.service/models/mysql"
	"integration.service/pkg/setting"
)

type FilesRemove struct {
}

func (fr *FilesRemove) Init() {

}

func (fr *FilesRemove) GetShort() string {
	return "图片服务器数据删除"
}

func (fr *FilesRemove) GetUse() string {
	return "filesRemove"
}

func (fr *FilesRemove) GetArgs() cobra.PositionalArgs {
	return cobra.ExactArgs(2)
}

func (fr *FilesRemove) SetFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

/**
 * @note: 删除图片服务器的文件，第一个参数 开始年月 第二个参数 结束年月
 * @auth: tongwz
 * @date  2022年2月17日16:39:24
**/
func (fr *FilesRemove) Run(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("需要传入两个参数,起始年月和结束年月")
		return
	}
	var startMonth = args[0]
	var endMonth = args[1]
	var tableNameMap = mysql.TableNameMap
	exceptCompanyId := setting.Cfg.Section("").Key("except_company_id").MustInt64(3)
	fileBasePath := setting.Cfg.Section("").Key("file_base_path").MustString("/data2/wwwroot/image-server-api/public/")
	env := setting.Cfg.Section("app").Key("env").MustString("test")

	if env == "test" {
		tableNameMap = mysql.TableNameTestMap
	}

	var delImageFiles = jz.DelImageFiles{ExceptCompanyId: exceptCompanyId, FileBasePath: fileBasePath, TableNameMap: tableNameMap}
	delImageFiles.DoIt(startMonth, endMonth)
}
