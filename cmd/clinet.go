package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 初始化
func init() {
	rootCmd.AddCommand(clientCmd)
}

// 设置命令
var clientCmd = &cobra.Command{
	Use: "client",
	Run: ClientHandle,
}

func ClientHandle(cmd *cobra.Command, args []string) {
	// TODO client.Run
	fmt.Println("client 还没写好")
}
