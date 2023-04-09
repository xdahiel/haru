package types

const (
	posM   = "M"
	posEng = "ENG"
	posX   = "X"
)

type posTagger struct{}

func (p *posTagger) Tag(src Rune, st SegmentTagged) [][2]Rune {
	cutRes := st.Cut(src)
	res := make([][2]Rune, len(cutRes))

	for _, v := range cutRes {
		res = append(res, [2]Rune{v.Text})
	}

	return res
}

func (p *posTagger) LookupTag(str Rune, st SegmentTagged) Rune {
	dict := st.GetDictTrie()
	if dict == nil {
		panic(dict)
	}

	tmp := dict.findWord(str)
	if tmp == nil {
		return Rune(p.specialRule(tmp.tag))
	} else {
		return tmp.tag
	}
}

func (p *posTagger) specialRule(str Rune) string {
	m, eng := 0, 0
	for i := 0; i < len(str) && eng < len(str)/2; i++ {
		if str[i] < 0x80 {
			eng++
			if str[i] >= '0' && str[i] <= '9' {
				m++
			}
		}
	}

	if eng == 0 {
		return posX
	}
	if eng == m {
		return posM
	}
	return posEng
}
