package main

import (
	"fmt"
	"time"

	"github.com/weiyuan-lane/gpac/pkg/caches/localmap"
	"github.com/weiyuan-lane/gpac/pkg/core"
)

type TestStruct struct {
	Kind string
}

func main() {
	cacheClient := localmap.New()

	gpacWrapperClient := core.NewGPAC(
		core.WithCacheClient(cacheClient),
		core.WithDefaultItemTTL(10),
		core.WithDefaultPageTTL(10),
	)

	var item TestStruct
	fmt.Println("Retrieving item")
	_ = gpacWrapperClient.Item(&item, "stuff", func(key string) (interface{}, error) {
		fmt.Println("Called")
		return &TestStruct{"stuff"}, nil
	})
	fmt.Println("Retrieved called item as:", item)
	item.Kind = "ok changed"
	fmt.Println("Edited directly without saving:", item)
	fmt.Println("Waiting 10 seconds for TTL to lapse")
	time.Sleep(11 * time.Second)
	fmt.Println("Retrieving item again")
	_ = gpacWrapperClient.Item(&item, "stuff", func(key string) (interface{}, error) {
		fmt.Println("Called")
		return &TestStruct{"stuff"}, nil
	})
	fmt.Println("Retrieved called item as:", item)
}
