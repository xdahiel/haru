package jieba

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestExtractor_Extract(t *testing.T) {
	log.SetFlags(log.LstdFlags)
	e := GetExtractor()

	f, err := os.Open("/home/chuns/code/haru/src/jieba/config/dict/test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	text, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(e.Extract(string(text), 10))
}
