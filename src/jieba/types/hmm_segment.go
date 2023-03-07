package types

type HmmSegment struct {
	segmentBase
	model *HmmModel
}

func newHmmSegment(modelPath string) HmmSegment {
	return HmmSegment{
		segmentBase: newSegmentBase(),
		model:       NewHmmModel(modelPath),
	}
}

func newHmmSegmentReady(model *HmmModel) HmmSegment {
	return HmmSegment{
		segmentBase: newSegmentBase(),
		model:       model,
	}
}

func (m *HmmSegment) Cut(sentence Rune) []Rune {
	res := make([]Rune, 0)
	f := NewPreFilter(sentence, m.symbols)
	for f.HasNext() {
		l, r := f.Next()
		res = append(res, m.cut(sentence[l:r])...)
	}
	return res
}

func (m *HmmSegment) cut(sentence Rune) []Rune {
	l, r := 0, 0
	res := make([]Rune, 0)
	for r < len(sentence) {
		if sentence[r] < 0x80 {
			if l != r {
				res = append(res, m.internalCut(sentence[l:r])...)
			}
			l = r
			for {
				r = m.sequentialLetterRule(sentence[l:])
				if l != r {
					break
				}
				r = m.numbersRule(sentence[l:])
				if l != r {
					break
				}
				r++
				break
			}
			res = append(res, sentence[l:r])
			l = r
		} else {
			r++
		}
	}

	if l != r {
		res = append(res, m.internalCut(sentence[l:r])...)
	}
	return res
}

func (m *HmmSegment) internalCut(sentence Rune) []Rune {
	status := m.viterbi(sentence)
	res := make([]Rune, 0)
	l, r := 0, 0
	for i, v := range status {
		if (v & 1) == 1 {
			r = i + 1
			res = append(res, sentence[l:r])
			l = r
		}
	}
	return res
}

func (m *HmmSegment) viterbi(sentence Rune) []int {
	x := len(sentence)
	y := statusSum

	size := x * y
	var (
		now, old, stat  int
		tmp, endE, endS float64
	)

	path := make([]int, size)
	weight := make([]float64, size)

	for j := 0; j < y; j++ {
		weight[j*x] = m.model.startProb[j] +
			m.model.getEmitProb(m.model.emitProbVec[j], sentence[0], minFloat64)
		path[j*x] = -1
	}

	emitProb := 0.0
	for i := 1; i < x; i++ {
		for j := 0; j < y; j++ {
			now = i + j*x
			weight[now] = minFloat64
			path[now] = wordEnd // warning

			emitProb = m.model.getEmitProb(m.model.emitProbVec[j], sentence[i], minFloat64)
			for k := 0; k < y; k++ {
				old = i - 1 + k*x
				tmp = weight[old] + m.model.transProb[k][j] + emitProb
				if tmp > weight[now] {
					weight[now] = tmp
					path[now] = k
				}
			}
		}
	}

	endE = weight[x-1+x*wordEnd]
	endS = weight[x-1+x*wordSep]
	stat = 0
	if endE >= endS {
		stat = wordEnd
	} else {
		stat = wordSep
	}

	status := make([]int, x)
	for i := x - 1; i >= 0; i-- {
		status[i] = stat
		stat = path[i+stat*x]
	}

	return status
}

func (m *HmmSegment) sequentialLetterRule(sentence Rune) int {
	x := sentence[0]
	if !((x >= 'a' && x <= 'z') || (x >= 'A' && x <= 'Z')) {
		return 0
	}
	for i, v := range sentence[1:] {
		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9')) {
			return i + 1
		}
	}

	return len(sentence)
}

func (m *HmmSegment) numbersRule(sentence Rune) int {
	x := sentence[0]
	if !(x >= '0' && x <= '9') {
		return 0
	}
	for i, v := range sentence {
		if !((v >= '0' && v <= '9') || v == '.') {
			return i + 1
		}
	}
	return len(sentence)
}
