package types

import (
	"haru/jieba/config"
	"testing"
)

func TestStopWordFilter(t *testing.T) {
	path := config.GetDictPaths()
	f := NewStopWordFilter(path[4])
	get := f.Has("不尽然")
	except := true
	if get != except {
		t.Errorf("except: %v, get: %v", except, get)
	}

	get = f.Has("北京大学")
	except = false
	if get != except {
		t.Errorf("except: %v, get: %v", except, get)
	}
}
