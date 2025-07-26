package cmd

import (
	"github.com/lin-snow/ech0/internal/cli"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 Web 和 SSH 服务",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoSSH()
		cli.DoServeWithBlock()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
