package cmd

import (
	"os"

	"github.com/lin-snow/ech0/internal/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ech0",
	Short: "开源、自托管、专注思想流动的轻量级发布平台",
	Long:  `开源、自托管、专注思想流动的轻量级发布平台`,

	// 这个 Run 会在没有子命令时执行
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoTui()
	},
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "启动 Ech0 TUI",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoTui()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoVersion()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看当前信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoEch0Info()
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "输出 Ech0 Logo",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoHello()
	},
}

func init() {
	// 解决Windows下使用 Cobra 触发 mousetrap 提示
	cobra.MousetrapHelpText = ""
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(helloCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
