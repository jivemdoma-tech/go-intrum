package gointrum

import (
	"fmt"
	"strconv"
	"time"
)

const (
	datetimeLayout string = "2006-01-02 15:04:05-07:00" // Формат даты и времени Intrum
	dateLayout     string = "2006-01-02-07:00"          // Формат даты Intrum
	timeLayout     string = "15:04:05-07:00"            // Формат времени Intrum

	DatetimeLayout string = "2006-01-02 15:04:05-07:00" // Формат даты и времени Intrum
	DateLayout     string = "2006-01-02-07:00"          // Формат даты Intrum
	TimeLayout     string = "15:04:05-07:00"            // Формат времени Intrum

	TypeStock    string = "stock"    // Тип сущности "Объект"
	TypeCustomer string = "customer" // Тип сущности "Контакт"
	TypeSale     string = "sale"     // Тип сущности "Сделка"
	TypeRequest  string = "request"  // Тип сущности "Заявка"
)

func addSliceToParams[T string | uint64 | uint32](fieldName string, params map[string]string, slice []T) {
	if len(slice) == 0 {
		return
	}

	for i, v := range slice {
		switch v := any(v).(type) {
		case string:
			params[fmt.Sprintf("params[%s][%d]", fieldName, i)] = v
		case uint32:
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
