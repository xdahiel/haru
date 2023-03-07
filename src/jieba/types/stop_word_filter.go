package types

import "haru/common/fileReader"

type StopWordFilter map[string]struct{}

func NewStopWordFilter(stopWordDict string) StopWordFilter {
	scanner := fileReader.NewScanner(stopWordDict)
	f := make(StopWordFilter)
	var (
		line string
		ok   bool
	)
	for {
		if line, ok = scanner.Next(); ok {
			f.insert(line)
		} else {
			break
		}
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
