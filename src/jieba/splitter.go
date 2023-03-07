// Package jieba
// 分词模块，参考 jieba(https://github.com/fxsjy/jieba) 的实现
package jieba

import (
	"haru/jieba/config"
	"haru/jieba/types"
	"sync"
)

type Splitter struct {
	querySeg  types.QuerySegment
	filter    types.StopWordFilter
	extractor Extractor
}

func newSplitter(dictPath, hmmPath, userDictPath, stopWordPath string) Splitter {
	return Splitter{
		querySeg: types.NewQuerySegment(dictPath, hmmPath, userDictPath),
		filter:   types.NewStopWordFilter(stopWordPath),
	}
}

func (s *Splitter) Cut(sentence string) []string {
	cutRes := s.querySeg.Cut(types.Rune(sentence), true)
	res := make([]string, 0)

	for _, v := range cutRes {
		if s.filter.Has(string(v)) {
			continue
		}
		res = append(res, string(v))
	}
	return res
}

var (
	splitterOnce sync.Once
	splitter     Splitter
)

func GetSplitter() Splitter {
	splitterOnce.Do(func() {
		path := config.GetDictPaths()
		splitter = newSplitter(path[0], path[1], path[2], path[4])
	})

	return splitter
}
