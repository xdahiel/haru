package types

import (
	"testing"
)

func TestUTrie_Seek(t *testing.T) {
	tr := GetTrie()
	tr.Insert("abc")
	tr.Insert("abcd")
	tr.Insert("a")
	tr.Insert("rtwrtg")
	tr.Insert("rtwrerfgtg")
	tr.Insert("成")
	tr.Insert("成功")
	tr.Insert("成都")
	tr.Insert("成都信息工程大学")
	tr.Insert("北京")
	tr.Insert("北漂")

	except := []string{"a", "abc", "abcd"}
	get := tr.Seek("a")
	if !isEqual(except, get) {
		t.Errorf("except: %v, get: %v", except, get)
	}

	except = []string{"rtwrtg", "rtwrerfgtg"}
	get = tr.Seek("r")
	if !isEqual(except, get) {
		t.Errorf("except: %v, get: %v", except, get)
	}

	except = []string{"成", "成功", "成都", "成都信息工程大学"}
	get = tr.Seek("成")
	if !isEqual(except, get) {
		t.Errorf("except: %v, get: %v", except, get)
	}

	except = []string{"北京", "北漂"}
	get = tr.Seek("北")
	if !isEqual(except, get) {
		t.Errorf("except: %v, get: %v", except, get)
	}
}

func isEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v1 := range s1 {
		if v1 != s2[i] {
			return false
		}
	}

	return true
}
