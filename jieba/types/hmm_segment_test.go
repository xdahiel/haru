package types

import (
	"haru/jieba/config"
	"testing"
)

func TestHmmSegment_Cut(t *testing.T) {
	path := config.GetDictPaths()
	m := newHmmSegment(path[1])

	except := "全世界\\都\\在\\学中国话\\"
	res := m.Cut(Rune("全世界都在学中国话"))
	get := ""
	for _, v := range res {
		get += string(v) + "\\"
	}
	if except != get {
		t.Errorf("except: %v \n get: %v", except, get)
	}

	except = "我\\在\\北京大学\\读手\\扶\\拖拉机\\专业\\"
	res = m.Cut(Rune("我在北京大学读手扶拖拉机专业"))
	get = ""
	for _, v := range res {
		get += string(v) + "\\"
	}
	if except != get {
		t.Errorf("except: %v \n get: %v", except, get)
	}
}
