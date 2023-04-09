package engine

import (
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"haru/crawl/weibo"
	"haru/engine/types"
	"haru/logs"
	"net/http"
	"strings"
)

func InitRouter(v1 *gin.RouterGroup) {
	v1.POST("/search", func(ctx *gin.Context) {
		query := ctx.Query("search")
		logs.Debug("query: %v", query)
		e := GetEngine()
		output := e.Search(types.SearchRequest{
			Text: query,
			RankOptions: &types.RankOptions{
				ScoringCriteria: &weibo.WeiboScoringCriteria{},
				OutputOffset:    0,
				MaxOutputs:      100,
			},
		})

		// 整理为输出格式
		var docs []*weibo.Weibo
		for _, doc := range output.Docs {
			wb := weibo.Wbs[doc.DocId]
			for _, t := range output.Tokens {
				wb.Text = strings.Replace(wb.Text, t, "<font color=red>"+t+"</font>", -1)
			}
			docs = append(docs, &wb)
		}
		response, _ := sonic.Marshal(&weibo.JsonResponse{Docs: docs})
		ctx.JSON(http.StatusOK, gin.H{
			"code": "2000",
			"msg":  "success",
			"data": string(response),
		})
	})
}
