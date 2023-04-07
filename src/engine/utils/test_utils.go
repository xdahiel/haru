package utils

import (
	"fmt"
	"runtime"
	"testing"
)

func Expect(t *testing.T, expect string, actual interface{}) {
	actualString := fmt.Sprint(actual)
	_, file, line, _ := runtime.Caller(1)
	if expect != actualString {
		t.Errorf("%s:%d 期待值=\"%s\", 实际=\"%s\"", file, line, expect, actualString)
	}
}
