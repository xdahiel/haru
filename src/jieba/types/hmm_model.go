package types

import (
	"haru/common/fileReader"
	"strconv"
	"strings"
)

type emitProbMap map[rune]float64

const (
	wordBeg = iota
	wordEnd
	wordMid
	wordSep

	statusSum
)

type HmmModel struct {
	stateMap  [statusSum]rune
	startProb [statusSum]float64 // 初始状态向量
	transProb [statusSum][statusSum]float64

	emitProbMapBeg emitProbMap
	emitProbMapEnd emitProbMap
	emitProbMapMid emitProbMap
	emitProbMapSep emitProbMap

	emitProbVec []*emitProbMap
}

func NewHmmModel(modelPath string) *HmmModel {
	model := &HmmModel{
		stateMap:       [4]rune{'B', 'E', 'M', 'S'},
		startProb:      [4]float64{},
		transProb:      [4][4]float64{},
		emitProbMapBeg: make(map[rune]float64),
		emitProbMapEnd: make(map[rune]float64),
		emitProbMapMid: make(map[rune]float64),
		emitProbMapSep: make(map[rune]float64),
		emitProbVec:    make([]*emitProbMap, 0, 4),
	}
	model.emitProbVec = append(model.emitProbVec,
		&model.emitProbMapBeg,
		&model.emitProbMapEnd,
		&model.emitProbMapMid,
		&model.emitProbMapSep,
	)
	model.loadModel(modelPath)
	return model
}

func (m *HmmModel) loadModel(filePath string) {
	scanner := fileReader.NewScanner(filePath)
	var (
		line string
		err  error
	)
	cnt := 0
	for scanner.HasNext() {
		line = scanner.Next()
		if cnt == 0 {
			tmp := strings.Split(line, " ")
			for i, v := range tmp {
				m.startProb[i], err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(tmp[i])
				}
			}
		} else if cnt <= statusSum {
			tmp := strings.Split(line, " ")
			for i, v := range tmp {
				m.transProb[cnt-1][i], err = strconv.ParseFloat(v, 64)
			}
		} else if cnt-statusSum == 1 {
			m.loadEmitProb(line, m.emitProbMapBeg)
		} else if cnt-statusSum == 2 {
			m.loadEmitProb(line, m.emitProbMapEnd)
		} else if cnt-statusSum == 3 {
			m.loadEmitProb(line, m.emitProbMapMid)
		} else if cnt-statusSum == 4 {
			m.loadEmitProb(line, m.emitProbMapSep)
		} else {
			break
		}

		cnt++
	}
}

func (m *HmmModel) loadEmitProb(line string, probMap emitProbMap) {
	if len(line) == 0 {
		panic(line)
	}

	var err error
	tmp := strings.Split(line, ",")
	for _, v := range tmp {
		pos := strings.LastIndex(v, ":")
		if !(pos > 0 && pos < len(v)-1) {
			panic(v)
		}
		probMap[[]rune(v)[0]], err = strconv.ParseFloat(v[pos+1:], 64)
		if err != nil {
			panic(err)
		}
	}
}

func (m *HmmModel) getEmitProb(mp *emitProbMap, key rune, defaultVal float64) float64 {
	v, ok := (*mp)[key]
	if ok {
		return v
	}
	return defaultVal
}
