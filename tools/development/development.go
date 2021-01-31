package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/weiyuan-lane/gpac/pkg/caches/localmap"
	"github.com/weiyuan-lane/gpac/pkg/caches/rediswrapper"
	"github.com/weiyuan-lane/gpac/pkg/core"
)

type TestStruct struct {
	Kind string
}

func main() {
	testRedis()
	// testLocalMap()
}

func testRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis.cache:6379",
		Password: "password",
		DB:       0,
	})

	cacheClient := rediswrapper.New(client)

	gpacWrapperClient := core.NewGPAC(
		core.WithCacheClient(cacheClient),
		core.WithDefaultItemTTL(10),
		core.WithDefaultPageTTL(10),
		core.WithUniqueNamespace("testingnamespace"),
	)

	testLogic(gpacWrapperClient)
}

func testLocalMap() {
	cacheClient := localmap.New()

	gpacWrapperClient := core.NewGPAC(
		core.WithCacheClient(cacheClient),
		core.WithDefaultItemTTL(10),
		core.WithDefaultPageTTL(10),
		core.WithUniqueNamespace("testingnamespace"),
	)

	testLogic(gpacWrapperClient)
}

func testLogic(gpacWrapperClient *core.PageAwareCache) {
	var item TestStruct
	fmt.Println("Retrieving item")
	_ = gpacWrapperClient.SimpleItem(&item, func(key string) (interface{}, error) {
		fmt.Println("Called")
		return &TestStruct{"stuff"}, nil
	}, "stuff")
	fmt.Println("Retrieved called item as:", item)
	_ = gpacWrapperClient.SimpleItem(&item, func(key string) (interface{}, error) {
		fmt.Println("Called")
		return &TestStruct{"stuff1"}, nil
	}, "stuff")
	fmt.Println("Retrieved called item (without calling direct but from cache) as:", item)
	item.Kind = "ok changed"
	fmt.Println("Edited directly without saving:", item)
	fmt.Println("Waiting 10 seconds for TTL to lapse")
	time.Sleep(11 * time.Second)
	fmt.Println("Retrieving item again")
	_ = gpacWrapperClient.SimpleItem(&item, func(key string) (interface{}, error) {
		fmt.Println("Called")
		return &TestStruct{"stuff2"}, nil
	}, "stuff")
	fmt.Println("Retrieved called item as:", item)
}
