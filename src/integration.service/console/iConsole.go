package console

import "github.com/spf13/cobra"

type IConsole interface {
	// 获取名称
	GetShort() string
	GetUse() string
	GetArgs() cobra.PositionalArgs
	// 设置命令參數
	SetFlags(cmd *cobra.Command) *cobra.Command
	// 运行命令
	Run(cmd *cobra.Command, args []string)

	// 初始化
	Init()
}