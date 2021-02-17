package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/weiyuan-lane/gpac/pkg/caches/localmap"
	"github.com/weiyuan-lane/gpac/pkg/caches/rediswrapper"
	"github.com/weiyuan-lane/gpac/pkg/gpac"
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

	gpacWrapperClient := gpac.NewGPAC(
		gpac.WithCacheClient(cacheClient),
		gpac.WithDefaultItemTTL(10),
		gpac.WithDefaultPageTTL(10),
		gpac.WithUniqueNamespace("testingnamespace"),
	)

	// testSimpleItem(gpacWrapperClient)
	// testItem(gpacWrapperClient)
	testCollection(gpacWrapperClient)
}

func testLocalMap() {
	cacheClient := localmap.New()

	gpacWrapperClient := gpac.NewGPAC(
		gpac.WithCacheClient(cacheClient),
		gpac.WithDefaultItemTTL(10),
		gpac.WithDefaultPageTTL(10),
		gpac.WithUniqueNamespace("testingnamespace"),
	)

	testSimpleItem(gpacWrapperClient)
}

func testSimpleItem(gpacWrapperClient *gpac.PageAwareCache) {
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

func testItem(gpacWrapperClient *gpac.PageAwareCache) {
	var item TestStruct

	fmt.Println("Retrieving item")
	_ = gpacWrapperClient.Item(&item, func(k ...gpac.ArgReference) (interface{}, error) {
		if len(k) != 2 {
			panic("e1")
		}

		return &TestStruct{"teststructval-1"}, nil
	},
		gpac.NewArgReference("str", "teststructval"),
		gpac.NewArgReference("num", 1),
	)
	fmt.Println("Retrieved called item as:", item)

	_ = gpacWrapperClient.Item(&item, func(k ...gpac.ArgReference) (interface{}, error) {
		if len(k) != 2 {
			panic("e2")
		}

		return &TestStruct{"teststructval-2"}, nil
	},
		gpac.NewArgReference("str", "teststructval"),
		gpac.NewArgReference("num", 1),
	)
	fmt.Println("Retrieved called item (without calling direct but from cache) as:", item)

	item.Kind = "ok changed"
	fmt.Println("Edited directly without saving:", item)
	fmt.Println("Waiting 10 seconds for TTL to lapse")
	time.Sleep(11 * time.Second)
	fmt.Println("Retrieving item again")

	_ = gpacWrapperClient.Item(&item, func(k ...gpac.ArgReference) (interface{}, error) {
		if len(k) != 2 {
			panic("e3")
		}

		return &TestStruct{"teststructval-3"}, nil
	},
		gpac.NewArgReference("str", "teststructval"),
		gpac.NewArgReference("num", 1),
	)
	fmt.Println("Retrieved called item as:", item)
}

func testCollection(gpacWrapperClient *gpac.PageAwareCache) {
	itemMap := map[int]TestStruct{}
	fmt.Println("Retrieving collection")

	err := gpacWrapperClient.Items(&itemMap, func(keys []string) (interface{}, error) {
		return map[int]TestStruct{
			0: {"stuff 1"},
			1: {"stuff 2"},
		}, nil
	}, []string{"one", "two"})
	fmt.Println("Retrieved called item as:", itemMap, err)

	itemMap = map[int]TestStruct{}
	err = gpacWrapperClient.Items(&itemMap, func(keys []string) (interface{}, error) {
		return map[int]TestStruct{
			0: {"stuff 3"},
			1: {"stuff 4"},
		}, nil
	}, []string{"one", "two"})
	fmt.Println("Retrieved called item (without calling direct but from cache) as:", itemMap, err)

	fmt.Println("Waiting 10 seconds for TTL to lapse")
	time.Sleep(11 * time.Second)
	fmt.Println("Retrieving item again")

	itemMap = map[int]TestStruct{}
	err = gpacWrapperClient.Items(&itemMap, func(keys []string) (interface{}, error) {
		return map[int]TestStruct{
			0: {"stuff 5"},
			1: {"stuff 6"},
		}, nil
	}, []string{"one", "two"})
	fmt.Println("Retrieved called item as:", itemMap, err)
}
