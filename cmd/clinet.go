package cmd

import "github.com/spf13/cobra"

// 初始化
func init() {
	rootCmd.AddCommand(clientCmd)
}

// 设置命令
var clientCmd = &cobra.Command{
	Use: "clinet",
	Run: ClientHandle,
}

func ClientHandle(cmd *cobra.Command, args []string) {
	// TODO client.Run
}
