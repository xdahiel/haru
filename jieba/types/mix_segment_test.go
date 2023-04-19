package types

import (
	"haru/jieba/config"
	"testing"
)

func TestMixSegment_Cut(t *testing.T) {
	path := config.GetDictPaths()
	x := NewMixSegment(path[0], path[1], path[2])
	except := "我\\在\\北京大学\\读\\手扶拖拉机\\专业\\"
	res := x.Cut(Rune("我在北京大学读手扶拖拉机专业"))
	get := ""
	for _, v := range res {
		get += string(v.Text) + "\\"
	}
	if except != get {
		t.Errorf("except: %v \n get: %v", except, get)
	}
}
