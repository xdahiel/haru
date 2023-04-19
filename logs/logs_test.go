package logs

import "testing"

func TestDebug(t *testing.T) {
	Debug("1 + 1 = %v", 2)
	Error("1 + 1 = 2")
	Info("1 + 1 = %v", 2)
}
