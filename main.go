package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"haru/common"
	"haru/crawl"
	"haru/crawl/weibo"
	"haru/engine"
	"haru/engine/types"
	"haru/logs"
	"haru/user"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.LstdFlags)

	common.Init()
	user.Init()
	flag.Parse()

	e := engine.GetEngine()
	go crawl.IndexToEngine(e)

	// 捕获ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Print("捕获Ctrl-c，退出服务器")
			e.Close()
			os.Exit(0)
		}
	}()

	r := gin.New()
	v1 := r.Group("/api/v1")

	user.InitRouter(v1)

	//r.Use(middlewares.Middlewares()...)

	r.SetFuncMap(template.FuncMap{
		"formatDate": Timestamp2String,
	})
	r.Static("/static", "web/public")
	r.LoadHTMLFiles("web/index.html", "web/login.html",
		"web/register.html", "web/result.html", "web/user.html")

	render(r)

	err := r.Run(":8080")
	if err != nil {
		logs.Error("run gin error: %v", err)
		return
	}
}

func render(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/register.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	r.GET("/login.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/result.html", func(ctx *gin.Context) {
		query := ctx.Query("query")
		logs.Debug("query: %v", query)
		e := engine.GetEngine()
		output := e.Search(types.SearchRequest{
			Text: query,
			RankOptions: &types.RankOptions{
				ScoringCriteria: &weibo.WeiboScoringCriteria{},
				OutputOffset:    0,
				MaxOutputs:      100,
			},
		})

		// 整理为输出格式
		var docs []weibo.Weibo
		for _, doc := range output.Docs {
			wb := weibo.Wbs[doc.DocId]
			docs = append(docs, wb)
		}
		ctx.HTML(http.StatusOK, "result.html", gin.H{
			"result": docs,
			"query":  query,
		})
	})
}
