package cmd

import (
	"cocoIM/mytest"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testConnCMD)
	rootCmd.PersistentFlags().Int64Var(&mytest.TcpConnNum, "tcp_conn_num", 10000, "tcp 连接的数量，默认10000")
}

// 设置新命令
var testConnCMD = &cobra.Command{
	Use: "testconn",
	Run: testConnHandle,
}

func testConnHandle(command *cobra.Command, args []string) {
	mytest.TestConnRunMain()
}
