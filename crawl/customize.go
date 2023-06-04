package crawl

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"strings"
)

type CustomizeCrawler struct {
	Url    string `json:"url"`
	Cookie string `json:"cookie"`

	recursionDepth int
	maxCrawlCount  int
}

func (c *CustomizeCrawler) Get(suffix string) string {
	var ck map[string]string
	// 解析cookie
	cookies := strings.Split(c.Cookie, ";")
	for _, cookie := range cookies {
		kvs := strings.Split(cookie, ":")
		ck[kvs[0]] = kvs[1]
	}

	formatCookie := []*http.Cookie{}
	for k, v := range ck {
		formatCookie = append(formatCookie, &http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	// 发送请求
	cookiesJar, _ := cookiejar.New(nil)
	url, _ := url2.Parse(c.Url + suffix)
	cookiesJar.SetCookies(url, formatCookie)
	client := &http.Client{Jar: cookiesJar}

	req, _ := http.NewRequest("GET", c.Url+suffix, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(body)

}

func (c *CustomizeCrawler) Parse(content string) Data {
	panic("implement me!")
}
