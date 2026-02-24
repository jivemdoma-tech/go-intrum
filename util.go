package gointrum

import (
	"fmt"
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

func addToSingularParams[T string | int64 | time.Time](params map[string]string, paramName string, paramValue T) {
	k := fmt.Sprintf("params[%s]", paramName)
	switch v := any(paramValue).(type) {
	case string:
		if v != "" {
			params[k] = v
		}
	case int64:
		if v != 0 {
			params[k] = strconv.FormatInt(v, 10)
		}
	case time.Time:
		if !v.IsZero() {
			params[k] = v.Format(DatetimeLayout)
		}
	}
}

func addBoolToSingularParams(params map[string]string, paramName string, paramValue string) {
	switch lower := strings.ToLower(strings.TrimSpace(paramValue)); lower {
	case "1", "true", "да":
		addToSingularParams(params, paramName, "1")
	case "0", "false", "нет":
		addToSingularParams(params, paramName, "0")
	case "ignore":
		addToSingularParams(params, paramName, "ignore")
	}
}

func addSliceToSingularParams[T string | int64](params map[string]string, paramName string, paramValue []T) {
	if len(paramValue) == 0 {
		return
	}

	for index, value := range paramValue {
		k := fmt.Sprintf("params[%s][%d]", paramName, index)
		switch v := any(value).(type) {
		case string:
			if v != "" {
				params[k] = v
			}
		case int64:
			if v != 0 {
				params[k] = strconv.FormatInt(v, 10)
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
