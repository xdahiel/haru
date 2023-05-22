package types

type trieSegment struct {
	segmentTagged
	dictTrie *DictTrie
	tagger   posTagger
}

func newTrieSegment(dictPath, userDictPath string) trieSegment {
	return trieSegment{
		segmentTagged: newSegmentTagged(),
		dictTrie:      NewDictTrie(dictPath, userDictPath, midWordWeight),
		tagger:        posTagger{},
	}
}

func newTrieSegmentReady(dictTrie *DictTrie) trieSegment {
	return trieSegment{
		segmentTagged: newSegmentTagged(),
		dictTrie:      dictTrie,
		tagger:        posTagger{},
	}
}

func (ts *trieSegment) Cut(sentence Rune) []segmentResp {
	return ts.cut(sentence, maxWordLength)
}

func (ts *trieSegment) cut(sentence Rune, length int) []segmentResp {
	f := NewPreFilter(sentence, ts.symbols) // 获取过滤器
	res := make([]segmentResp, 0)

	var l, r int
	for f.HasNext() {
		l, r = f.Next()                                                       // 获取下一个可分词区间
		dags := ts.dictTrie.findAsDagWithMaxWordLength(sentence[l:r], length) // 在树上构建有向无环图
		ts.calcDP(dags)                                                       // 初始化动态规划参数

		offset := 0
		cutRes := ts.cutByDag(sentence[l:r], dags) // 使用动态规划计算出权重最高的分词序列
		for _, cut := range cutRes {
			res = append(res, segmentResp{
				Text:  cut,
				Start: l + offset,
			})
			offset += len(cut)
		}
	}

	return res
}

func (ts *trieSegment) calcDP(dags []Dag) {
	var (
		nextPos int
		info    *DictUnit
		weight  float64
	)

	for i := len(dags) - 1; i >= 0; i-- {
		dags[i].info = nil
		dags[i].weight = minFloat64
		if len(dags[i].nexts) == 0 {
			panic(dags[i])
		}

		for _, v := range dags[i].nexts {
			nextPos = int(v.offset)
			info = v.info
			weight = 0.0

			if nextPos+1 < len(dags) {
				weight += dags[nextPos+1].weight
			}

			if info != nil {
				weight += info.weight
			} else {
				weight += ts.dictTrie.minWeight
			}

			if weight > dags[i].weight {
				dags[i].info = info
				dags[i].weight = weight
			}
		}
	}
}

func (ts *trieSegment) cutByDag(dst Rune, dags []Dag) []Rune {
	res := make([]Rune, 0, len(dags))

	for i := 0; i < len(dags); {
		p := dags[i].info
		if p != nil {
			if len(p.word) == 0 {
				panic(p.word)
			}
			res = append(res, dst[i:i+len(p.word)])
			i += len(p.word)
		} else {
			res = append(res, dst[i:i+1])
			i++
		}
	}

	return res
}

func (ts *trieSegment) isUserDictSingleChineseWord(val rune) bool {
	return ts.dictTrie.userDictSingleChineseWord.Has(val)
}

func (ts *trieSegment) GetDictTrie() *DictTrie {
	return ts.dictTrie
}

func (ts *trieSegment) Tag(src Rune) [][2]Rune {
	return ts.tagger.Tag(src, ts)
}
