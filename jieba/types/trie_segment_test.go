package types

import (
	"haru/jieba/config"
	"log"
	"testing"
)

func TestTrieSegment_Cut(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	path := config.GetDictPaths()
	ts := newTrieSegment(path[0], path[2])

	except := "我\\在\\北京大学\\读\\手扶拖拉机\\专业\\"
	res := ts.Cut(Rune("我在北京大学读手扶拖拉机专业"))
	get := ""
	for _, v := range res {
		get += string(v.Text) + "\\"
	}
	if except != get {
		t.Errorf("except: %v\nget: %v", except, get)
	}

	except = "我\\是\\一个\\大学生\\，\\我\\有\\光明\\的\\前途\\"
	res = ts.Cut(Rune("我是一个大学生，我有光明的前途"))
	get = ""
	for _, v := range res {
		get += string(v.Text) + "\\"
	}
	if except != get {
		t.Errorf("except: %v\nget: %v", except, get)
	}
}
