package types

import "haru/common/fileReader"

type StopWordFilter map[string]struct{}

func NewStopWordFilter(stopWordDict string) StopWordFilter {
	scanner := fileReader.NewScanner(stopWordDict)
	f := make(StopWordFilter)
	for scanner.HasNext() {
		f.insert(scanner.Next())
	}
	return f
}

func (f *StopWordFilter) insert(s string) {
	(*f)[s] = struct{}{}
}

func (f *StopWordFilter) Has(s string) bool {
	_, ok := (*f)[s]
	return ok
}
