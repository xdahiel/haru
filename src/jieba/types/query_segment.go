package types

type QuerySegment struct {
	mixSeg MixSegment
	trie   *DictTrie
}

func NewQuerySegment(dictPath, modelPath, userDictPath string) QuerySegment {
	x := NewMixSegment(dictPath, modelPath, userDictPath)
	return QuerySegment{
		mixSeg: x,
		trie:   x.GetDictTrie(),
	}
}

func newQuerySegmentReady(trie *DictTrie, model *HmmModel) QuerySegment {
	return QuerySegment{
		mixSeg: NewMixSegmentReady(trie, model),
		trie:   trie,
	}
}

func (q *QuerySegment) Cut(sentence Rune, hmm bool) []Rune {
	f := NewPreFilter(sentence, q.mixSeg.symbols)
	res := make([]Rune, 0)
	for f.HasNext() {
		l, r := f.Next()
		res = append(res, q.cut(sentence[l:r], hmm)...)
	}
	return res
}

func (q *QuerySegment) cut(sentence Rune, hmm bool) []Rune {
	mixRes := q.mixSeg.cut(sentence, hmm)

	res := make([]Rune, 0)
	for _, v := range mixRes {
		if len(v) > 2 {
			for j := 0; j+1 < len(v); j++ {
				if q.trie.findWord(v[j:j+2]) != nil {
					res = append(res, v[j:j+2])
				}
			}
		}
		if len(v) > 3 {
			for j := 0; j+2 < len(v); j++ {
				if q.trie.findWord(v[j:j+3]) != nil {
					res = append(res, v[j:j+3])
				}
			}
		}
		res = append(res, v)
	}
	return res
}
