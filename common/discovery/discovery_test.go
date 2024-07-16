package discovery

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestServiceDiscovery(t *testing.T) {
	ctx := context.Background()
	discovery := NewServiceDiscovery(&ctx)
	defer discovery.Close()
	discovery.WatchService("/web/", func(key, value string) {
		fmt.Println("this set_func")
	}, func(key, value string) {
		fmt.Println("this del_func")
	})
	discovery.WatchService("/gRPC/", func(key, value string) {
		fmt.Println("this set_func")
	}, func(key, value string) {
		fmt.Println("this del_func")
	})

	for {
		select {
		case <-time.Tick(10 * time.Second):
			fmt.Println("10s")
		}
	}

}
