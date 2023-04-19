package types

import (
	"fmt"
	"haru/jieba/config"
	"testing"
)

func TestQuerySegment_Cut(t *testing.T) {
	path := config.GetDictPaths()
	q := NewQuerySegment(path[0], path[1], path[2])
	except := "我\\在\\北京\\大学\\北京大学\\读\\手扶\\拖拉\\拖拉机\\手扶拖拉机\\专业\\"
	res := q.Cut(Rune("我在北京大学读手扶拖拉机专业"), true)
	get := ""
	for _, v := range res {
		get += v.Text + "\\"
		fmt.Println(v.Text, v.Start)
	}
	if except != get {
		t.Errorf("except: %v \nget: %v", except, get)
	}

	except = "小明\\硕士\\毕业\\于\\中国\\科学\\学院\\科学院\\中国科学院\\计算\\计算所\\，\\后\\在\\日本\\京都\\大学\\日本京都大学\\深造\\"
	res = q.Cut(Rune("小明硕士毕业于中国科学院计算所，后在日本京都大学深造"), true)
	get = ""
	for _, v := range res {
		get += v.Text + "\\"
	}
	if except != get {
		t.Errorf("except: %v \nget: %v", except, get)
	}
}
