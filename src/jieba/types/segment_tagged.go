package types

type SegmentTagged interface {
	SegmentBase
	Tag(src Rune) [][2]Rune
	GetDictTrie() *DictTrie
}

type segmentTagged struct {
	segmentBase
}

func newSegmentTagged() segmentTagged {
	return segmentTagged{
		segmentBase: newSegmentBase(),
	}
}
