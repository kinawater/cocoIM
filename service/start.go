package service

import (
	"cocoIM/config"
	"context"
	"github.com/spf13/cobra"
	"strconv"
)

// 服务启动配置
type ServerStartOptions struct {
	id     string //当前服务器id
	listen string //监听地址
}

func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
	options := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "im-server",
		Short: "IM 服务启动……",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, options, version)
		},
	}
	defaultServerID := config.ServerConf.ID
	defaultServerListen := ":" + strconv.Itoa(config.ServerConf.HTTPPort)
	cmd.PersistentFlags().StringVarP(&options.id, "server-id", "i", defaultServerID, "current server local id")
	cmd.PersistentFlags().StringVarP(&options.listen, "listen", "l", defaultServerListen, "listen address")

	return cmd
}

func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	server := NewServer(opts.id, opts.listen)
	defer server.shutdown()
	return server.Start()
}
