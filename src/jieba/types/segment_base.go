package types

const SPECIAL_SEPARATORS = " \t\n\xEF\xBC\x8C\xE3\x80\x82"

type RSet map[rune]struct{}

func NewRSet() RSet {
	return make(map[rune]struct{})
}

func (r *RSet) Insert(s rune) {
	(*r)[s] = struct{}{}
}

func (r *RSet) Has(s rune) bool {
	_, ok := (*r)[s]
	return ok
}

func (r *RSet) clear() {
	for k, _ := range *r {
		delete(*r, k)
	}
}

type SegmentBase interface {
	Cut(sentence Rune) []Rune
}

type segmentBase struct {
	symbols RSet
}

func newSegmentBase() segmentBase {
	sb := segmentBase{symbols: NewRSet()}
	sb.resetSeparators(SPECIAL_SEPARATORS)
	return sb
}

func (sb *segmentBase) resetSeparators(s string) {
	sb.symbols.clear()
	for _, v := range s {
		sb.symbols[v] = struct{}{}
	}
}
