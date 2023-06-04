package crawl

type Data struct {
	Id           uint64 `json:"id"`
	Timestamp    uint64 `json:"timestamp"`
	Username     string `json:"username"`
	RepostsCount uint64 `json:"reposts_count"`
	Text         string `json:"text"`
}

type Crawler interface {
	Get(string) string
	Parse(string) Data
}
