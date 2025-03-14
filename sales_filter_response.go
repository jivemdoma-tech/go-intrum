package gointrum

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type SalesFilterResponse struct {
	Status string           `json:"status"`
	Data   *SalesFilterData `json:"data"`
}
type SalesFilterData struct {
	List []*Sale `json:"list"`
	// Count any `json:"count"` // TODO
}
type Sale struct {
	ID                   uint64                `json:"id,string"`              // ID сделки
	CustomersID          uint64                `json:"customers_id,string"`    // ID контакта
	EmployeeID           uint64                `json:"employee_id,string"`     // ID ответственного
	AdditionalEmployeeID []uint64              `json:"additional_employee_id"` // Массив ID доп. ответственных
	DateCreate           time.Time             `json:"date_create"`            // Дата создания
	SalesTypeID          uint16                `json:"sales_type_id,string"`   // ID типа активности
	SaleStageID          uint16                `json:"sale_stage_id,string"`   // ID стадии
	SaleName             string                `json:"sale_name"`              // Название сделки
	SaleActivityType     string                `json:"sale_activity_type"`     // Тип последней активности
	SaleActivityDate     time.Time             `json:"sale_activity_date"`     // Дата последней активности сделк
	Fields               map[string]*SaleField `json:"fields"`                 // Данные полей
}
// Использовать метод GetField для получения значения поля // TODO
type SaleField struct {
	DataType string `json:"datatype"`
	Value    any    `json:"value"`
}

func (s *Sale) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias Sale

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		AdditionalEmployeeID []string `json:"additional_employee_id"`
		DateCreate           string   `json:"date_create"`
		SaleActivityDate     string   `json:"sale_activity_date"`
	}{
		Alias: (*Alias)(s), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена

	parsedDate, err := time.Parse(datetimeLayout, aux.DateCreate)
	if err != nil {
		return err
	}
	s.DateCreate = parsedDate

	parsedDate, err = time.Parse(datetimeLayout, aux.SaleActivityDate)
	if err != nil {
		return err
	}
	s.SaleActivityDate = parsedDate

	newSlice := make([]uint64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if newValue, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, newValue)
		}
	}
	s.AdditionalEmployeeID = newSlice

	return nil
}

// Методы получения значений Sale

func (f *SaleField) getFieldStr() string {
	if v, ok := f.Value.(string); ok {
		return strings.Join(strings.Fields(v), " ")
	}
	return ""
}

func (s *Sale) GetField(fieldID uint64) any {
	f, exists := s.Fields[strconv.FormatUint(fieldID, 10)]
	if !exists {
		return ""
	}

	// В сделках не передаются типы "integer_range", "decimal_range", "datetime_range", "date_range", "time_range"
	// Вместо этого передается базовый тип + хэш-таблица со значениями "from", "to":
	/*
		"datatype": "integer",
		"value": {
			"from": "2",
			"to": "64"
		}
	*/
	// Поэтому для типов "integer", "decimal", "datetime", "date", "time" добавил дополнительные проверки
	switch f.DataType {
	// bool
	case "radio":
		if v, err := strconv.ParseBool(f.getFieldStr()); err == nil {
			return v
		}
		return false

	// string
	case "text", "select", "file":
		return f.getFieldStr()

	// []string
	case "multiselect":
		return strings.Split(f.getFieldStr(), ",")

	// [2]string
	case "point":
		if m, ok := f.Value.(map[string]string); ok && len(m) >= 2 {
			return [2]string{m["x"], m["y"]}
		}
		return [2]string{}

	// int64 | [2]int64
	case "integer":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, parseInt)
		}
		return parseInt(f.getFieldStr())

	// [2]int64
	case "integer_range":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, parseInt)
		}
		return [2]int64{}

	// float64 | [2]float64
	case "decimal":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, parseFloat)
		}
		return parseFloat(f.getFieldStr())

	// [2]float64
	case "decimal_range":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, parseFloat)
		}
		return [2]float64{}

	// float64
	case "price":
		return parseFloat(f.getFieldStr())

	// time.Time | [2]time.Time
	case "datetime":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, func(s string) time.Time {
				return parseTime(s, datetimeLayout)
			})
		}
		return parseTime(f.getFieldStr(), datetimeLayout)

	// [2]time.Time
	case "datetime_range":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, func(s string) time.Time {
				return parseTime(s, datetimeLayout)
			})
		}
		return [2]time.Time{}

	// time.Time | [2]time.Time
	case "date":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, func(s string) time.Time {
				return parseTime(s, dateLayout)
			})
		}
		return parseTime(f.getFieldStr(), dateLayout)

	// [2]time.Time
	case "date_range":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, func(s string) time.Time {
				return parseTime(s, dateLayout)
			})
		}
		return [2]time.Time{}

	// time.Time | [2]time.Time
	case "time":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, func(s string) time.Time {
				return parseTime(s, timeLayout)
			})
		}
		return parseTime(f.getFieldStr(), timeLayout)

	// [2]time.Time
	case "time_range":
		if m, ok := f.Value.(map[string]string); ok {
			return parseRange(m, func(s string) time.Time {
				return parseTime(s, timeLayout)
			})
		}
		return [2]time.Time{}

	// time.Duration
	case "duration":
		v := parseInt(f.getFieldStr())
		return v * int64(time.Minute)

	// []uint64
	case "attach":
		if vAttach, ok := f.Value.([]map[string]string); ok && len(vAttach) > 0 {
			vIDs := make([]uint64, 0, len(vAttach))
			for _, v := range vAttach {
				if id, err := strconv.ParseUint(v["id"], 10, 64); err == nil {
					vIDs = append(vIDs, id)
				}
			}
			return vIDs
		}
		return []uint64{}
	}

	return ""
}
