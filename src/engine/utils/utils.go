package utils

import "github.com/bytedance/sonic"

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FormatOutput(v any) string {
	data, _ := sonic.MarshalString(v)
	return data
}
