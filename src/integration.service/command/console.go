package command

import (
	"github.com/spf13/cobra"
	"integration.service/console"
)

func init()  {
	Register()
}

var ConsoleCmd = &cobra.Command{
	Use:   "console",
	Short: "执行命令",
	Args:  cobra.ExactArgs(1),
	PersistentPostRun: func(cmd *cobra.Command, args []string) {

	},
	//命令执行
	Run: func(cmd *cobra.Command, args []string) {
		//执行动作
	},
}

func Register()  {
	var consoleList []console.IConsole
	consoleList = append(consoleList, new(console.FilesRemove))
	consoleList = append(consoleList, new(console.RepairMediaData))
	consoleList = append(consoleList, new(console.RepairWechatBind))

	for _, c := range consoleList {
		cmd := &cobra.Command{
			Use:   c.GetUse(),
			Short: c.GetShort(),
			Args:  c.GetArgs(),
			Run:   c.Run,
		}
		ConsoleCmd.AddCommand(cmd)
		c.Init()
	}
}