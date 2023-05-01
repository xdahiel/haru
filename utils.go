package main

import "time"

func Timestamp2String(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format(time.DateTime)
}
