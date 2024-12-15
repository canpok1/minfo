package internal

import "time"

type FormatTimeType int

const (
	YYYYMMDD FormatTimeType = iota
	YYYYMMDDHHMMSS
)

func ToJST(v time.Time) time.Time {
	jst := time.FixedZone("JST", 9*60*60)
	return v.In(jst)
}

func FormatTime(v time.Time, t FormatTimeType) string {
	switch t {
	case YYYYMMDD:
		return v.Format("2006/01/02")
	case YYYYMMDDHHMMSS:
		return v.Format("2006/01/02 15:04:05 MST")
	default:
		return ""
	}
}
