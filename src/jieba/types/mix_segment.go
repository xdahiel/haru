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

func (x *MixSegment) Cut(sentence Rune) []Rune {
	f := NewPreFilter(sentence, x.symbols)
	res := make([]Rune, 0)
	for f.HasNext() {
		l, r := f.Next()
		res = append(res, x.cut(sentence[l:r], true)...)
	}
	return res
}

func (x *MixSegment) cut(sentence Rune, hmm bool) []Rune {
	if !hmm {
		return x.trieSeg.cut(sentence, maxWordLength)
	}

	if len(sentence) == 0 {
		panic(sentence)
	}
	words := x.trieSeg.cut(sentence, maxWordLength)
	res := make([]Rune, 0)
	for i, v := range words {
		if len(v) > 0 || (len(v) == 1 && !x.trieSeg.isUserDictSingleChineseWord(v[0])) {
			res = append(res, v)
			continue
		}

		j := i
		for j < len(v) && (len(words[j]) == 1) && !x.trieSeg.isUserDictSingleChineseWord(words[j][0]) {
			j++
		}

		if j-1 < i {
			panic(j)
		}

		// optimization need
		tmp := make(Rune, 0)
		for _, v := range words[i : j-1] {
			tmp = append(tmp, v...)
		}
		tmp = append(tmp, words[j-1][:2]...)
		hmmRes := x.hmmSeg.cut(tmp)
		res = append(res, hmmRes...)
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
