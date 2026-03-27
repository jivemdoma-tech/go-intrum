package intrum

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Типы основных сущностей

	EntityTypeStock    string = "stock"    // Тип сущности "Объект"
	EntityTypeCustomer string = "customer" // Тип сущности "Контакт"
	EntityTypeSale     string = "sale"     // Тип сущности "Сделка"
	EntityTypeRequest  string = "request"  // Тип сущности "Заявка"

	// Форматы даты и времени

	DatetimeLayout   string = "2006-01-02 15:04:05" // Формат даты и времени Intrum
	DatetimeLayoutUI string = "02.01.2006 15:04:05" // Формат даты и времени Intrum (UI)
	DateLayout       string = "2006-01-02"          // Формат даты Intrum
	DateLayoutUI     string = "02.01.2006"          // Формат даты Intrum (UI)
	TimeLayout       string = "15:04:05"            // Формат времени Intrum

	// Форматы bool

	BoolYes    string = "1"      // "Да"
	BoolNo     string = "0"      // "Нет"
	BoolIgnore string = "ignore" // "Да" + "Нет"
)

var (
	dateLayouts     = []string{DateLayout, DatetimeLayoutUI, DatetimeLayout, DateLayoutUI}
	datetimeLayouts = []string{DatetimeLayout, DateLayoutUI, DateLayout, DatetimeLayoutUI}
)

// LocalizeTime применяет локальный часовой пояс (без сдвига) к time.Time.
func LocalizeTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
}

// ParseDate возвращает дату из переданной строки.
func ParseDate(s string) time.Time {
	var result time.Time
	// Отказоустойчивый парсинг
	if s != "" {
		for _, layout := range dateLayouts {
			if parsed, err := time.Parse(layout, s); err == nil {
				result = parsed
			}
		}
	}
	if result.IsZero() {
		return time.Time{}
	}

	result = LocalizeTime(result)
	result = time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, result.Location())
	return result
}

// ParseDatetime возвращает дату и время из переданной строки.
func ParseDatetime(s string) time.Time {
	var result time.Time
	// Отказоустойчивый парсинг
	if s != "" {
		for _, layout := range datetimeLayouts {
			if parsed, err := time.Parse(layout, s); err == nil {
				result = parsed
			}
		}
	}
	if result.IsZero() {
		return time.Time{}
	}

	result = LocalizeTime(result)
	result = time.Date(result.Year(), result.Month(), result.Day(), result.Hour(), result.Day(), result.Minute(), 0, result.Location())
	return result
}

func ParseInt(s string) int64 {
	if r, err := strconv.ParseInt(s, 10, 64); err == nil {
		return r
	}
	return 0
}

func ParseFloat(s string) float64 {
	if r, err := strconv.ParseFloat(s, 64); err == nil {
		return r
	}
	return 0.0
}

func ParseRange[T any](m map[string]string, parseFunc func(string) T) [2]T {
	return [2]T{parseFunc(m["from"]), parseFunc(m["to"])}
}

// Координаты

// Point - координата на карте.
type Point struct {
	Lat float64 // Широта
	Lon float64 // Долгота
}

// NewPoint возвращает Point.
func NewPoint(lat, lon float64) *Point { return &Point{Lat: lat, Lon: lon} }

// NewPointFromStrings парсит строковые значения координат и возвращает Point.
func NewPointFromStrings(latStr, lonStr string) (*Point, error) {
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lat: %w", err)
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lon: %w", err)
	}

	return &Point{Lat: lat, Lon: lon}, nil
}

func (p *Point) StringLat() string {
	if p == nil {
		return ""
	}
	return strconv.FormatFloat(p.Lat, 'f', 10, 64)
}

func (p *Point) StringLon() string {
	if p == nil {
		return ""
	}
	return strconv.FormatFloat(p.Lon, 'f', 10, 64)
}

// Добавление в параметры запроса

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
