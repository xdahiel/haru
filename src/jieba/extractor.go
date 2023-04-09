// 基于 IDF 的关键词提取
package jieba

import (
	"container/heap"
	"haru/common/fileReader"
	"haru/jieba/config"
	"haru/jieba/types"
	"log"
	"strconv"
	"strings"
	"sync"
)

type wordSet map[string]struct{}

func newWordSet() wordSet {
	return make(wordSet)
}

func (w wordSet) insert(s string) {
	w[s] = struct{}{}
}

func (w wordSet) has(s string) bool {
	_, ok := w[s]
	return ok
}

type Word struct {
	word    string
	offsets []uint
	weight  float64
}

func (w *Word) appendOffset(o uint) Word {
	w.offsets = append(w.offsets, o)
	return *w
}

type wordPq []*Word

func (q wordPq) Len() int           { return len(q) }
func (q wordPq) Less(i, j int) bool { return q[i].weight < q[j].weight }
func (q wordPq) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *wordPq) Push(w any) {
	item := w.(*Word)
	*q = append(*q, item)
}
func (q *wordPq) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*q = old[:n-1]
	return item
}

type Extractor struct {
	seg        types.MixSegment
	idfMap     map[string]float64
	idfAverage float64
	stopWords  wordSet
}

func newExtractor(dictPath, hmmPath, idfPath, stopWordPath, userDict string) Extractor {
	e := Extractor{
		seg:        types.NewMixSegment(dictPath, hmmPath, userDict),
		idfMap:     make(map[string]float64),
		idfAverage: 0,
	}
	e.loadIdfDict(idfPath)
	e.loadStopWordDict(stopWordPath)
	return e
}

func (e *Extractor) Extract(sentence string, topN int) []string {
	return e.extract(types.Rune(sentence), topN)
}

func (e *Extractor) extract(sentence types.Rune, topN int) []string {
	cutRes := e.seg.Cut(sentence)
	wordMap := make(map[string]*Word)
	offset := 0

	for _, v := range cutRes {
		t := uint(offset)
		offset += len(v.Text)

		sv := string(v.Text)
		if (len(v.Text) == 1 && len(string(v.Text)) == 3) || e.stopWords.has(sv) {
			continue
		}
		if _, ok := wordMap[sv]; !ok {
			wordMap[sv] = &Word{
				word:    "",
				offsets: make([]uint, 0, 1),
				weight:  0,
			}
		}
		wordMap[sv].appendOffset(t)
		wordMap[sv].weight += 1
	}

	if offset != len(sentence) {
		panic("word illegal!!")
	}

	res := make([]Word, 0)
	for k, word := range wordMap {
		if v, ok := e.idfMap[word.word]; ok {
			word.weight = word.weight * v
		} else {
			word.weight = word.weight * e.idfAverage
		}
		word.word = k
		res = append(res, *word)
	}
	if topN > len(res) {
		topN = len(res)
	}

	pq := make(wordPq, 0, topN)
	heap.Init(&pq)
	for i, v := range res {
		if len(pq) == topN {
			if pq[0].weight > v.weight {
				continue
			}
			heap.Pop(&pq)
		}
		heap.Push(&pq, &res[i])
	}
	for i := range res[:topN] {
		res[i] = *(heap.Pop(&pq).(*Word))
	}

	ret := make([]string, 0, topN)
	for _, v := range res[:topN] {
		ret = append(ret, v.word)
	}
	return ret
}

func (e *Extractor) loadIdfDict(idfPath string) {
	var (
		line   string
		idf    float64
		err    error
		idfSum = 0.0
		cnt    = 0
	)

	scanner := fileReader.NewScanner(idfPath)
	for scanner.HasNext() {
		line = scanner.Next()
		if len(line) == 0 {
			log.Println("skipped at ", cnt, ", line: ", line)
			continue
		}
		buf := strings.Split(line, " ")
		if len(buf) != 2 {
			log.Println("skipped at ", cnt, ", line: ", line)
			continue
		}
		idf, err = strconv.ParseFloat(buf[1], 64)
		if err != nil {
			panic(err)
		}
		e.idfMap[buf[1]] = idf
		idfSum += idf
		cnt++
	}

	if cnt == 0 {
		panic(cnt)
	}
	e.idfAverage = idfSum / float64(cnt)
	if e.idfAverage <= 0 {
		panic(e.idfAverage)
	}
}

func (e *Extractor) loadStopWordDict(filePath string) {
	var line string

	e.stopWords = make(wordSet)
	scanner := fileReader.NewScanner(filePath)
	for scanner.HasNext() {
		line = scanner.Next()
		e.stopWords.insert(line)
	}
}

var (
	extractorOnce sync.Once
	extractor     Extractor
)

func GetExtractor() Extractor {
	extractorOnce.Do(func() {
		path := config.GetDictPaths()
		extractor = newExtractor(path[0], path[1], path[3], path[4], path[2])
	})

	return extractor
}
