package main

import (
	"context"
	"flag"
	"haru/user/model"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"

	"haru/associate"
	assoTypes "haru/associate/types"
	"haru/common"
	"haru/crawl"
	"haru/crawl/weibo"
	"haru/engine"
	"haru/engine/types"
	"haru/logs"
	"haru/middlewares"
	"haru/user"
	userController "haru/user/controller"
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
	v2 := r.Group("/api/v2")

	v2.Use(middlewares.Middlewares()...)

	user.InitRouter(v1)
	user.InitUserRouter(v2)

	associate.InitRouter(v1)

	//r.Use(middlewares.Middlewares()...)

	r.SetFuncMap(template.FuncMap{
		"formatDate": Timestamp2String,
	})
	r.Static("/static", "web/public")
	r.LoadHTMLFiles("web/index.html", "web/login.html", "web/company.html",
		"web/register.html", "web/result.html", "web/user.html", "web/image_result.html")

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
	r.GET("/company.html", userController.Company)
	r.GET("/result.html", func(ctx *gin.Context) {
		query := ctx.Query("query")
		assoTypes.GetTrie().Insert(query)
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

		advs := make([]*model.Advertise, 0)
		for _, token := range output.Tokens {
			tmp, err := model.FindAdvertiseByKeyword(token)
			if err != nil {
				continue
			}
			advs = append(advs, tmp...)
		}
		if len(advs) > 10 {
			advs = advs[:10]
		}

		// 整理为输出格式
		var docs []weibo.Weibo
		for _, doc := range output.Docs {
			wb := weibo.Wbs[doc.DocId]
			docs = append(docs, wb)
		}
		ctx.HTML(http.StatusOK, "result.html", gin.H{
			"result": docs,
			"query":  query,
			"advs":   advs,
		})
	})

	r.GET("/image_result.html", func(c *gin.Context) {
		logs.Info("visit image_result.html")

		rdb := common.GetRDB()
		ctx := context.Background()

		length, _ := rdb.LLen(ctx, "names").Result()

		//images, _ := rdb.LRange(ctx, "images", 0, length-1).Result()
		names, _ := rdb.LRange(ctx, "names", 0, length-1).Result()
		logs.Debug("get from redis: %v", names)

		//type Image struct {
		//	Name string `json:"name"`
		//	Data string `json:"data"`
		//}

		//var results []Image
		//for i, image := range images {
		//	results = append(results, Image{
		//		Name: names[i],
		//		Data: image,
		//	})
		//}

		c.HTML(http.StatusOK, "image_result.html", gin.H{
			"results": names,
		})
	})
}
