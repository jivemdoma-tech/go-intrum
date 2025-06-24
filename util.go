package gointrum

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DatetimeLayout string = "2006-01-02 15:04:05" // Формат даты и времени Intrum
	DateLayout     string = "2006-01-02"          // Формат даты Intrum
	TimeLayout     string = "15:04:05"            // Формат времени Intrum

	TypeStock    string = "stock"    // Тип сущности "Объект"
	TypeCustomer string = "customer" // Тип сущности "Контакт"
	TypeSale     string = "sale"     // Тип сущности "Сделка"
	TypeRequest  string = "request"  // Тип сущности "Заявка"
)

func returnErrBadParams(methodURL string) error {
	u, _ := url.ParseRequestURI(methodURL)
	return fmt.Errorf("failed to create request for method %s: %s", u.Path, statusBadParams)
}

func addToParams[T string | int8 | uint8 | uint16 | uint64 | int16 | int64 | time.Time](params map[string]string, paramName string, v T) {
	k := fmt.Sprintf("params[%s]", paramName)
	switch v := any(v).(type) {
	case string:
		if v != "" {
			params[k] = v
		}
	case int8, uint8, int16, uint16:
		vInt := v.(int64)
		if vInt > 0 {
			params[k] = strconv.FormatInt(vInt, 10)
		}
	case int64:
		if v > 0 {
			params[k] = strconv.FormatInt(v, 10)
		}
	case uint64:
		if v > 0 {
			params[k] = strconv.FormatUint(v, 10)
		}
	case time.Time:
		if !v.IsZero() {
			params[k] = v.Format(DatetimeLayout)
		}
	}
}

func addBoolStringToParams(params map[string]string, paramName string, v string) {
	k := fmt.Sprintf("params[%s]", paramName)
	switch vLower := strings.ToLower(strings.TrimSpace(v)); {
	case vLower == "1", vLower == "true":
		params[k] = "1"
	case vLower == "0", vLower == "false":
		params[k] = "0"
	case vLower == "ignore":
		params[k] = "ignore"
	}
}

func addSliceToParams[T string | int64 | uint64](params map[string]string, paramName string, paramSlice []T) {
	if len(paramSlice) == 0 {
		return
	}

	for i, v := range paramSlice {
		k := fmt.Sprintf("params[%s][%d]", paramName, i)
		switch v := any(v).(type) {
		case string:
			if v != "" {
				params[k] = v
			}
		case uint64:
			if v != 0 {
				params[k] = strconv.FormatUint(v, 10)
			}
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
