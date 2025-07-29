package cmd

import (
	"github.com/lin-snow/ech0/internal/cli"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份数据",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoBackup()
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "恢复数据",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取待恢复的备份文件路径
		if len(args) < 1 {
			cmd.Help()
			return
		}

		cli.DoRestore(args[0])
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
}
