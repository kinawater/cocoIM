package cmd

import (
	"cocoIM/mytest"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCMD)
}

// 设置新命令
var testCMD = &cobra.Command{
	Use: "mytest",
	Run: testHandle,
}

func testHandle(command *cobra.Command, args []string) {
	mytest.RunMain()
}
