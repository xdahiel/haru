// Package jieba
// 分词模块，参考 jieba(https://github.com/fxsjy/jieba) 的实现
package jieba

import (
	"haru/jieba/config"
	"haru/jieba/types"
	"sync"
)

type Segmenter struct {
	querySeg  types.QuerySegment
	filter    types.StopWordFilter
	extractor Extractor
}

func newSplitter(dictPath, hmmPath, userDictPath, stopWordPath string) *Segmenter {
	return &Segmenter{
		querySeg: types.NewQuerySegment(dictPath, hmmPath, userDictPath),
		filter:   types.NewStopWordFilter(stopWordPath),
	}
}

func (s *Segmenter) Cut(sentence string) []types.SegmentResp {
	cutRes := s.querySeg.Cut(types.Rune(sentence), true)
	res := make([]types.SegmentResp, 0)

	for _, v := range cutRes {
		if s.filter.Has(v.Text) {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (s *Segmenter) Close() {
}

var (
	splitterOnce sync.Once
	segmenter    *Segmenter
)

func GetSplitter() *Segmenter {
	splitterOnce.Do(func() {
		path := config.GetDictPaths()
		segmenter = newSplitter(path[0], path[1], path[2], path[4])
	})

	return segmenter
}
