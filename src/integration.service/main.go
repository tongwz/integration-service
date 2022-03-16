package main

import (
	"github.com/spf13/cobra"
	"integration.service/command"
)

func main()  {
	var cmd = &cobra.Command{Use:"comm"}

	cmd.AddCommand(command.ConsoleCmd)

	_ = cmd.Execute()
}