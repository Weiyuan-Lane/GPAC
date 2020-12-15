package main

import (
	"github.com/weiyuan-lane/gpac/pkg/caches/localmap"
	"github.com/weiyuan-lane/gpac/pkg/core"
)

func main() {
	cacheClient := localmap.New()

	core.NewGPAC(
		core.WithCacheClient(cacheClient),
		core.WithDefaultItemTTL(100),
		core.WithDefaultPageTTL(100),
	)
}
