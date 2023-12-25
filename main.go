package main

import (
	"cocoIM/service"
	"context"
	"flag"
	"github.com/spf13/cobra"
)

const version = "v0.1"

func main() {
	flag.Parse()
	cmdRoot := &cobra.Command{
		Use:     "im-server",
		Version: version,
		Short:   "IM chat demo",
	}
	ctx := context.Background()
	cmdRoot.AddCommand(service.NewServerStartCmd(ctx, version))
	if err := cmdRoot.Execute(); err != nil {
		panic(err)
	}
}
