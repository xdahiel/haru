package weibo

import (
	"bufio"
	"haru/engine/engine"
	"haru/engine/types"
	"haru/logs"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	SecondsInADay     = 86400
	MaxTokenProximity = 2
)

var (
	weiboData = "/home/haru/Projects/haru/src/crawl/data/weibo_data.txt"
	Wbs       = map[uint64]Weibo{}
)

type Weibo struct {
	Id           uint64 `json:"id"`
	Timestamp    uint64 `json:"timestamp"`
	UserName     string `json:"user_name"`
	RepostsCount uint64 `json:"reposts_count"`
	Text         string `json:"text"`
}

/*
******************************************************************************

	索引

******************************************************************************
*/
func IndexWeibo(searcher *engine.Engine) {
	// 读入微博数据
	file, err := os.Open(weiboData)
	if err != nil {
		logs.Error("%v", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), "||||")
		if len(data) != 10 {
			continue
		}
		wb := Weibo{}
		wb.Id, _ = strconv.ParseUint(data[0], 10, 64)
		wb.Timestamp, _ = strconv.ParseUint(data[1], 10, 64)
		wb.UserName = data[3]
		wb.RepostsCount, _ = strconv.ParseUint(data[4], 10, 64)
		wb.Text = data[9]
		Wbs[wb.Id] = wb
	}

	log.Print("添加索引")
	for docId, weibo := range Wbs {
		searcher.IndexDocument(docId, types.DocumentIndexData{
			Content: weibo.Text,
			Fields: WeiboScoringFields{
				Timestamp:    weibo.Timestamp,
				RepostsCount: weibo.RepostsCount,
			},
		}, false)
	}

	searcher.FlushIndex()
	log.Printf("索引了%d条微博\n", len(Wbs))
}

/*
******************************************************************************

	评分

******************************************************************************
*/
type WeiboScoringFields struct {
	Timestamp    uint64
	RepostsCount uint64
}

type WeiboScoringCriteria struct {
}

func (criteria WeiboScoringCriteria) Score(
	doc types.IndexedDocument, fields interface{}) []float32 {
	if reflect.TypeOf(fields) != reflect.TypeOf(WeiboScoringFields{}) {
		return []float32{}
	}
	wsf := fields.(WeiboScoringFields)
	output := make([]float32, 3)
	if doc.TokenProximity > MaxTokenProximity {
		output[0] = 1.0 / float32(doc.TokenProximity)
	} else {
		output[0] = 1.0
	}
	output[1] = float32(wsf.Timestamp / (SecondsInADay * 3))
	output[2] = float32(doc.BM25 * (1 + float32(wsf.RepostsCount)/10000))
	return output
}

type JsonResponse struct {
	Docs []*Weibo `json:"docs"`
}
