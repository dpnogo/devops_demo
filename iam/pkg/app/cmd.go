package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

// RunCommandFunc 定义应用程序的命令启动回调函数
type RunCommandFunc func(args []string) error

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false // ?

	// 子节点绑定父节点
	if len(c.commands) > 0 {
		for _, subCmd := range c.commands {
			cmd.AddCommand(subCmd.cobraCommand())
		}
	}

	// 添加cmd对应运行函数
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	return cmd
}

func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}
