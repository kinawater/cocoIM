package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// 指定配置目录
var ConfigPath string

func init() {
	//cobra启动时候需要运行的方法
	cobra.OnInitialize(initConfig)
	// 设置提示选项
	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "./config.yaml", "config file (default is ./config.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "coco",
	Short: "一个很装杯的IM系统",
	Run:   cocoIM,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cocoIM(cmd *cobra.Command, args []string) {
	fmt.Println("服务器启动……当然是假的")
}

// 初始化配置
func initConfig() {
}
