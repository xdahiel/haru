package types

type MixSegment struct {
	segmentTagged
	trieSeg trieSegment
	hmmSeg  HmmSegment
	tagger  posTagger
}

func NewMixSegment(trieDict, hmmDict, userDict string) MixSegment {
	ts := newTrieSegment(trieDict, userDict)
	hs := newHmmSegment(hmmDict)
	return MixSegment{
		segmentTagged: newSegmentTagged(),
		trieSeg:       ts,
		hmmSeg:        hs,
		tagger:        struct{}{},
	}
}

func NewMixSegmentReady(trie *DictTrie, model *HmmModel) MixSegment {
	return MixSegment{
		segmentTagged: newSegmentTagged(),
		trieSeg:       newTrieSegmentReady(trie),
		hmmSeg:        newHmmSegmentReady(model),
		tagger:        struct{}{},
	}
}

func (x *MixSegment) Cut(sentence Rune) []segmentResp {
	f := NewPreFilter(sentence, x.symbols)
	res := make([]segmentResp, 0)
	for f.HasNext() {
		l, r := f.Next()
		res = append(res, x.cut(sentence[l:r], true)...)
	}
	return res
}

func (x *MixSegment) cut(sentence Rune, hmm bool) []segmentResp {
	if !hmm {
		return x.trieSeg.cut(sentence, maxWordLength)
	}

	if len(sentence) == 0 {
		panic(sentence)
	}
	words := x.trieSeg.cut(sentence, maxWordLength)
	res := make([]segmentResp, 0)
	for i, v := range words {
		if len(v.Text) > 0 || (len(v.Text) == 1 && !x.trieSeg.isUserDictSingleChineseWord(v.Text[0])) {
			res = append(res, v)
			continue
		}

		j := i
		for j < len(v.Text) && (len(words[j].Text) == 1) && !x.trieSeg.isUserDictSingleChineseWord(words[j].Text[0]) {
			j++
		}

		if j-1 < i {
			panic(j)
		}

		// TODO: 待优化
		tmp := make(Rune, 0)
		for _, v1 := range words[i : j-1] {
			tmp = append(tmp, v1.Text...)
		}
		tmp = append(tmp, words[j-1].Text[:2]...)
		hmmRes := x.hmmSeg.cut(tmp)
		offset := 0
		for _, v1 := range hmmRes {
			res = append(res, segmentResp{
				Text:  v1,
				Start: v.Start + offset,
			})
			offset += len(v1)
		}
		i = j - 1
	}
	return res
}

func (x *MixSegment) GetDictTrie() *DictTrie {
	return x.trieSeg.dictTrie
}

func (x *MixSegment) Tag(src Rune) [][2]Rune {
	return x.tagger.Tag(src, x)
}

func (x *MixSegment) LookupTag(str Rune) Rune {
	return x.tagger.LookupTag(str, x)
}
