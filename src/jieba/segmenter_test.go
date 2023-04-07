package jieba

import (
	"testing"
)

func TestSplitter_Cut(t *testing.T) {
	s := GetSplitter()
	except := "小明\\硕士\\毕业\\中国\\科学\\学院\\科学院\\中国科学院\\计算\\计算所\\日本\\京都\\大学\\日本京都大学\\深造\\"
	res := s.Cut("小明硕士毕业于中国科学院计算所，后在日本京都大学深造")
	get := ""
	for _, v := range res {
		get += v + "\\"
	}
	if except != get {
		t.Errorf("except: %v \nget: %v", except, get)
	}
}
