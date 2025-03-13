package gointrum

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	datetimeLayout string = "2006-01-02 15:04:05"
	dateLayout     string = "2006-01-02"
	timeLayout     string = "15:04:05"
)

func addSliceToParams[T string | uint64 | uint16](fieldName string, params map[string]string, slice []T) {
	if len(slice) == 0 {
		return
	}

	for i, v := range slice {
		switch v := any(v).(type) {
		case string:
			params[fmt.Sprintf("params[%s][%d]", fieldName, i)] = v
		case uint16:
			params[fmt.Sprintf("params[%s][%d]", fieldName, i)] = strconv.FormatUint(uint64(v), 10)
		case uint64:
			params[fmt.Sprintf("params[%s][%d]", fieldName, i)] = strconv.FormatUint(v, 10)
		}
	}
}

func parseInt(s string) int64 {
	if r, err := strconv.ParseInt(s, 10, 64); err == nil {
		return r
	}
	return 0
}

func parseFloat(s string) float64 {
	if r, err := strconv.ParseFloat(s, 64); err == nil {
		return r
	}
	return 0.0
}

func parseTime(s, layout string) time.Time {
	if t, err := time.Parse(layout, s); err == nil {
		return t
	}
	return time.Time{}
}

func parseRange[T any](m map[string]string, parseFunc func(string) T) [2]T {
	var r [2]T
	r[0] = parseFunc(m["from"])
	r[1] = parseFunc(m["to"])
	return r
}
