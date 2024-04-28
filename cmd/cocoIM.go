package cmd

import "github.com/spf13/cobra"

// 指定配置目录
var ConfigPath string

func init() {
	//cobra启动时候需要运行的方法
	cobra.OnInitialize(initConfig)

}

var rootCmd = cobra.Command{
	Use:   "coco",
	Short: "一个很装杯的IM系统",
}

func cocoIM() {}

// 初始化配置
func initConfig() {

}
