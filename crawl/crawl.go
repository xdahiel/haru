package crawl

import (
	"github.com/robfig/cron/v3"
	"haru/common/fileReader"
	"haru/crawl/weibo"
	"haru/engine/engine"
)

var c cron.Cron

func InitCrawl() {
	c := cron.New()
	c.AddFunc("@day", crawl)
}

func crawl() {
	scanner := fileReader.NewScanner("crawl_list.txt")
	for scanner.HasNext() {
		go handle(scanner.Next())
	}
}

func handle(url string) {

}

func IndexToEngine(search *engine.Engine) {
	weibo.IndexWeibo(search)
}
