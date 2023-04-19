package fileReader

import (
	"testing"
)

func TestNewScanner(t *testing.T) {
	path := "/home/chuns/code/haru/src/jieba/config/dict/hmm_model.utf8"
	s := NewScanner(path)

	cnt := 1
	for s.HasNext() {
		cnt++
	}

	if cnt != 10 {
		t.Errorf("except: %v, get: %v", 10, cnt)
	}
}

func d(err error) {
	if err != nil {
		panic(err)
	}
}
