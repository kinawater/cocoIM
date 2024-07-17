package cmd

import (
	"github.com/spf13/cobra"
)

// 初始化
func init() {
	rootCmd.AddCommand(gatewayCmd)
}

// 设置命令
var gatewayCmd = &cobra.Command{
	Use: "gateway",
	Run: GatewayHandle,
}

func GatewayHandle(cmd *cobra.Command, args []string) {
	// TODO 尚未实现 RunMain
}
