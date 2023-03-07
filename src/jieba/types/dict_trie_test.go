package types

import (
	"haru/jieba/config"
	"log"
	"testing"
)

func TestNewDictTrie(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	path := config.GetDictPaths()
	dt := NewDictTrie(path[0], path[2], midWordWeight)

	find := dt.has(Rune("成都"))
	if !find {
		t.Errorf("except: %v, get: %v", true, find)
	}

	find = dt.has(Rune("信息"))
	if !find {
		t.Errorf("except: %v, get: %v", true, find)
	}

	find = dt.has(Rune("工程"))
	if !find {
		t.Errorf("except: %v, get: %v", true, find)
	}

	find = dt.has(Rune("大学"))
	if !find {
		t.Errorf("except: %v, get: %v", true, find)
	}

	dags := dt.findAsDag(Rune("成都信息工程大学手扶拖拉机专业"))
	except := "成都\\信息\\信息工程\\工程\\大学\\手扶\\手扶拖拉机\\拖拉\\拖拉机\\专业\\"
	ret := travelDag(dags, "\\")
	if ret != except {
		t.Errorf("except: %v, get: %v", except, ret)
	}

	dags = dt.findAsDag(Rune("秦始皇派蒙恬攻打西域"))
	except = "秦始皇\\始皇\\派蒙\\蒙恬\\攻打\\西域\\"
	ret = travelDag(dags, "\\")
	if ret != except {
		t.Errorf("except: %v, get: %v", except, ret)
	}
}

func travelDag(dags []Dag, dim string) string {
	res := ""
	for i, v1 := range dags {
		for _, v2 := range v1.nexts {
			if v2.offset == uint(i) {
				continue
			}
			for _, v := range dags[i : v2.offset+1] {
				res += string(v.r)
			}
			res += dim
		}
	}

	return res
}
