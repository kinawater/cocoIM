package cmd

import (
	"cocoIM/ipconf"
	"github.com/spf13/cobra"
)

// ipconf的主要功能是给请求端分配ip地址，ip地址会根据一定的规则进行动态排序，给请求端分配一个当前较为空闲的机器
// 所谓空闲机器不能简单的靠剩余资源来决定，比如ABC三台机器，AB两台是老机器，MAX也就1000，而第三台新机器C的MAX是10000，如果
// A,B 各有900的空闲，而C有1000，那么AB的负载率应该是10%，但是C的负载率已经到了90%虽然看起来C有更多的空间，但是分配的时候应该优先
// 负载率低的AB，大致是这么个意思

func init() {
	rootCmd.AddCommand(ipConfCmd)
}

// 设置新命令
var ipConfCmd = &cobra.Command{
	Use: "ipconf",
	Run: IpConfHandle,
}

func IpConfHandle(command *cobra.Command, args []string) {
	ipconf.RunMain(ConfigPath)
}
