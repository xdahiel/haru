package types

import (
	"haru/common/fileReader"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	minFloat64    = -3.14e100
	maxFloat64    = 3.14e100
	dictColumnNum = 3
	unknownTag    = ""
)

type UserWordWeightOption uint8

const (
	minWordWeight UserWordWeightOption = iota
	midWordWeight
	maxWordWeight
)

type DictTrie struct {
	freqSum                   float64
	minWeight                 float64
	maxWeight                 float64
	midWeight                 float64
	userWordDefaultWeight     float64
	userDictSingleChineseWord RSet

	trie           *Trie
	staticNodeInfo []DictUnit
	activeNodeInfo []DictUnit
}

func NewDictTrie(dictPath, userDictPath string, option UserWordWeightOption) *DictTrie {
	dt := new(DictTrie)
	dt.userDictSingleChineseWord = NewRSet()
	dt.staticNodeInfo = make([]DictUnit, 0)
	dt.activeNodeInfo = make([]DictUnit, 0)

	dt.loadDict(dictPath)
	dt.freqSum = dt.calcFreqSum(dt.staticNodeInfo)
	dt.calcWeight(dt.staticNodeInfo, dt.freqSum)
	dt.setStaticWordWeight(option)

	if len(userDictPath) > 0 {
		dt.loadUserDict(userDictPath)
	}
	dt.createTrie(dt.staticNodeInfo)

	return dt
}

func (dt *DictTrie) findWord(word Rune) *DictUnit {
	return dt.trie.findWord(word)
}

func (dt *DictTrie) findAsDag(word Rune) []Dag {
	return dt.findAsDagWithMaxWordLength(word, maxWordLength)
}

func (dt *DictTrie) findAsDagWithMaxWordLength(word Rune, length int) []Dag {
	return dt.trie.findAsDagWithMaxWordLength(word, length)
}

func (dt *DictTrie) has(word Rune) bool {
	return dt.findWord(word) != nil
}

func (dt *DictTrie) loadDict(dictPath string) {
	s := fileReader.NewScanner(dictPath)
	var line string

	for s.HasNext() {
		line = s.Next()
		tmp := strings.Split(line, " ")
		if len(tmp) != dictColumnNum {
			panic(line)
		}

		weight, err := strconv.ParseFloat(tmp[1], 64)
		if err != nil {
			panic(tmp[1])
		}
		dt.staticNodeInfo = append(dt.staticNodeInfo, DictUnit{
			word:   Rune(tmp[0]),
			weight: weight,
			tag:    Rune(tmp[2]),
		})
	}
}

func (dt *DictTrie) loadUserDict(userDictPath string) {
	files := strings.Split(userDictPath, "|;")
	for _, v := range files {
		var line string
		s := fileReader.NewScanner(v)
		for s.HasNext() {
			line = s.Next()
			if len(line) == 0 {
				continue
			}
			dt.insertUserDictNode(line)
		}
	}
}

func (dt *DictTrie) insertUserDictNode(line string) {
	buf := strings.Split(line, " ")
	var info DictUnit
	if len(buf) == 1 {
		info = DictUnit{
			word:   Rune(buf[0]),
			weight: dt.userWordDefaultWeight,
			tag:    Rune(unknownTag),
		}
	} else if len(buf) == 2 {
		info = DictUnit{
			word:   Rune(buf[0]),
			weight: dt.userWordDefaultWeight,
			tag:    Rune(buf[1]),
		}
	} else if len(buf) == 3 {
		weight, err := strconv.ParseFloat(buf[1], 64)
		if err != nil {
			panic(err)
		}
		if dt.freqSum <= 0 {
			panic(dt.freqSum)
		}
		weight = weight / dt.freqSum

		info = DictUnit{
			word:   Rune(buf[0]),
			weight: weight,
			tag:    Rune(unknownTag),
		}
	}

	dt.staticNodeInfo = append(dt.staticNodeInfo, info)
	if len(info.word) == 1 {
		dt.userDictSingleChineseWord.Insert(info.word[0])
	}
}

func (dt *DictTrie) calcFreqSum(infos []DictUnit) float64 {
	sum := 0.0
	for _, v := range infos {
		sum += v.weight
	}
	return sum
}

func (dt *DictTrie) calcWeight(infos []DictUnit, sum float64) {
	if sum <= 0 {
		panic(sum)
	}

	for i, _ := range infos {
		if infos[i].weight <= 0 {
			panic(infos[i].weight)
		}
		infos[i].weight = math.Log(infos[i].weight / sum)
	}
}

func (dt *DictTrie) setStaticWordWeight(option UserWordWeightOption) {
	if len(dt.staticNodeInfo) == 0 {
		panic(dt.staticNodeInfo)
	}

	infos := dt.staticNodeInfo
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].weight < infos[j].weight
	})

	dt.minWeight = infos[0].weight
	dt.midWeight = infos[len(infos)/2].weight
	dt.maxWeight = infos[len(infos)-1].weight

	switch option {
	case minWordWeight:
		dt.userWordDefaultWeight = dt.minWeight
	case midWordWeight:
		dt.userWordDefaultWeight = dt.midWeight
	case maxWordWeight:
		dt.userWordDefaultWeight = dt.maxWeight
	}
}

func (dt *DictTrie) createTrie(infos []DictUnit) {
	if len(infos) == 0 {
		panic(infos)
	}

	words := make([]Rune, 0, len(infos))
	vals := make([]*DictUnit, 0, len(infos))

	for i, v := range infos {
		words = append(words, v.word)
		vals = append(vals, &infos[i])
	}

	dt.trie = NewTrie(words, vals)
}
