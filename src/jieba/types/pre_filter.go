package types

type PreFilter struct {
	cursor   int
	sentence Rune
	symbols  RSet
}

func NewPreFilter(sentence Rune, symbols RSet) PreFilter {
	return PreFilter{
		cursor:   0,
		sentence: sentence,
		symbols:  symbols,
	}
}

func (f *PreFilter) HasNext() bool {
	return f.cursor < len(f.sentence)
}

func (f *PreFilter) Next() (int, int) {
	l, r := f.cursor, 0
	for f.cursor < len(f.sentence) {
		if f.symbols.Has(f.sentence[f.cursor]) {
			if l == f.cursor {
				f.cursor++
			}
			r = f.cursor
			return l, r
		}
		f.cursor++
	}
	r = len(f.sentence)
	return l, r
}
