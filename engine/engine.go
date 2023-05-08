package engine

import (
	"haru/engine/engine"
	"haru/engine/types"
	"sync"
)

var (
	search     engine.Engine
	searchOnce sync.Once
)

func GetEngine() *engine.Engine {
	searchOnce.Do(func() {
		search = engine.Engine{}
		search.Init(types.EngineInitOptions{
			SegmenterDictionaries: "../jieba/config/dict/jieba.dict.utf8",
			StopTokenFile:         "../jieba/config/dict/stop_words.utf8",
			IndexerInitOptions: &types.IndexerInitOptions{
				IndexType: types.LocationsIndex,
			},
		})
	})

	return &search
}
