package gointrum

import (
	"fmt"
	"strconv"
	"time"
)

const (
	datetimeLayout string = "2006-01-02 15:04:05" // Формат даты и времени Intrum
	dateLayout     string = "2006-01-02"          // Формат даты Intrum
	timeLayout     string = "15:04:05"            // Формат времени Intrum

	DatetimeLayout string = "2006-01-02 15:04:05" // Формат даты и времени Intrum
	DateLayout     string = "2006-01-02"          // Формат даты Intrum
	TimeLayout     string = "15:04:05"            // Формат времени Intrum

	TypeStock    string = "stock"    // Тип сущности "Объект"
	TypeCustomer string = "customer" // Тип сущности "Контакт"
	TypeSale     string = "sale"     // Тип сущности "Сделка"
	TypeRequest  string = "request"  // Тип сущности "Заявка"
)

func addToParams[T string | uint16 | uint64](params map[string]string, paramName string, v T) {
	k := fmt.Sprintf("params[%s]", paramName)
	switch v := any(v).(type) {
	case string:
		if v != "" {
			params[k] = v
		}
	case uint16:
		if v != 0 {
			params[k] = strconv.FormatUint(uint64(v), 10)
		}
	case uint64:
		if v != 0 {
			params[k] = strconv.FormatUint(v, 10)
		}
	}
}

func addSliceToParams[T string | uint64](params map[string]string, paramName string, paramSlice []T) {
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

func addToMultiParams[T string | uint16 | uint64](params map[string]string, paramIndex int, paramName string, v T) {
	k := fmt.Sprintf("params[%d][%s]", paramIndex, paramName)
	switch v := any(v).(type) {
	case string:
		if v != "" {
			params[k] = v
		}
	case uint16:
		if v != 0 {
			params[k] = strconv.FormatUint(uint64(v), 10)
		}
	case uint64:
		if v != 0 {
			params[k] = strconv.FormatUint(v, 10)
		}
	}
}

func addSliceToMultiParams[T string | uint64](params map[string]string, paramIndex int, paramName string, paramSlice []T) {
	if len(paramSlice) == 0 {
		return
	}

	for i, v := range paramSlice {
		k := fmt.Sprintf("params[%d][%s][%d]", paramIndex, paramName, i)
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
