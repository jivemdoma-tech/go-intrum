package gointrum

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DatetimeLayout   string = "2006-01-02 15:04:05" // Формат даты и времени Intrum
	DatetimeLayoutUI string = "02.01.2006 15:04:05" // Формат даты и времени Intrum (UI)
	DateLayout       string = "2006-01-02"          // Формат даты Intrum
	DateLayoutUI     string = "02.01.2006"          // Формат даты Intrum (UI)
	TimeLayout       string = "15:04:05"            // Формат времени Intrum

	EntityTypeStock    string = "stock"    // Тип сущности "Объект"
	EntityTypeCustomer string = "customer" // Тип сущности "Контакт"
	EntityTypeSale     string = "sale"     // Тип сущности "Сделка"
	EntityTypeRequest  string = "request"  // Тип сущности "Заявка"
)

var (
	dateLayouts     = []string{DateLayout, DatetimeLayoutUI, DatetimeLayout, DateLayoutUI}
	datetimeLayouts = []string{DatetimeLayout, DateLayoutUI, DateLayout, DatetimeLayoutUI}
)

func localizeTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
}

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

func parseDate(s string) time.Time {
	var result time.Time

	if s != "" {
		// Отказоустойчивый парсинг
		for _, layout := range dateLayouts {
			if parsed, err := time.Parse(layout, s); err == nil {
				result = parsed
			}
		}
		// Локализация часового пояса
		if !result.IsZero() {
			result = localizeTime(result)
		}
	}

	return result
}

func parseDatetime(s string) time.Time {
	var result time.Time

	if s != "" {
		// Отказоустойчивый парсинг
		for _, layout := range datetimeLayouts {
			if parsed, err := time.Parse(layout, s); err == nil {
				result = parsed
			}
		}
		// Локализация часового пояса
		if !result.IsZero() {
			result = localizeTime(result)
		}
	}

	return result
}

func parseRange[T any](m map[string]string, parseFunc func(string) T) [2]T {
	var r [2]T
	r[0] = parseFunc(m["from"])
	r[1] = parseFunc(m["to"])
	return r
}
