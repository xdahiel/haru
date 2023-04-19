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
			// 如果你希望使用持久存储，启用下面的选项
			// 默认使用boltdb持久化，如果你希望修改数据库类型
			// 请修改 WUKONG_STORAGE_ENGINE 环境变量
			// UsePersistentStorage: true,
			// PersistentStorageFolder: "weibo_search",
		})
	})

	return &search
}
