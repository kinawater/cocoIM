package cmd

import (
	"cocoIM/gateway"
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
	gateway.RunMain(ConfigPath)
}
