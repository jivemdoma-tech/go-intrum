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

func addSliceToParams[T uint16 | uint64](params map[string]string, fieldName string, values []T) {
	if len(values) == 0 {
		return
	}

	switch any(values[0]).(type) {
	case uint16, uint64:
		for i, v := range values {
			k := fmt.Sprintf("params[%s][%d]", fieldName, i)
			params[k] = strconv.FormatUint(uint64(v), 10)
		}
	}
}

func getParamsSize(data interface{}) int {
	v := reflect.ValueOf(data)

	if !v.IsValid() {
		return 0
	}

	// Если передан pointer - разыменовывается
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return 0
		}
		v = v.Elem()
	}

	// Обработка структуры, мапы, слайса и массива
	switch v.Kind() {
	case reflect.Struct:
		count := 0
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if f.IsZero() {
				continue
			}
			count += getParamsSize(f.Interface())
		}
		return count

	case reflect.Map:
		count := 0
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			count += getParamsSize(val.Interface())
		}
		return count + v.Len()

	case reflect.Slice, reflect.Array:
		count := 0
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			count += getParamsSize(elem.Interface())
		}
		return count + v.Len()

	default:
		// Для базовых типов просто считаем их как 1 элемент
		if v.IsZero() {
			return 0
		}
		return 1
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
