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

func (q *QuerySegment) Cut(sentence Rune, hmm bool) []SegmentResp {
	f := NewPreFilter(sentence, q.mixSeg.symbols)
	res := make([]SegmentResp, 0)
	for f.HasNext() {
		l, r := f.Next()
		res = append(res, q.cut(sentence[l:r], hmm)...)
	}
	return res
}

func (q *QuerySegment) cut(sentence Rune, hmm bool) []SegmentResp {
	mixRes := q.mixSeg.cut(sentence, hmm)

	res := make([]SegmentResp, 0)
	for _, v := range mixRes {
		if len(v.Text) > 2 {
			for j := 0; j+1 < len(v.Text); j++ {
				if q.trie.findWord(v.Text[j:j+2]) != nil {
					res = append(res, SegmentResp{
						Text:  string(v.Text[j : j+2]),
						Start: v.Start + j,
					})
				}
			}
		}
		if len(v.Text) > 3 {
			for j := 0; j+2 < len(v.Text); j++ {
				if q.trie.findWord(v.Text[j:j+3]) != nil {
					res = append(res, SegmentResp{
						Text:  string(v.Text[j : j+3]),
						Start: v.Start + j,
					})
				}
			}
		}
		res = append(res, SegmentResp{
			Text:  string(v.Text),
			Start: v.Start,
		})
	}
	return res
}
