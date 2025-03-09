package gointrum

import (
	"strconv"
	"strings"
	"time"
)

type SalesFilterResponse struct {
	Status string           `json:"status"`
	Data   *SalesFilterData `json:"data"`
}
type SalesFilterData struct {
	List  []*Sale `json:"list"`
	Count any     `json:"count"`
}
type Sale struct {
	ID                   string                `json:"id"`                     // ID сделки
	CustomersID          string                `json:"customers_id"`           // ID контакта
	EmployeeID           string                `json:"employee_id"`            // ID ответственного
	AdditionalEmployeeID []string              `json:"additional_employee_id"` // Массив ID дополнительных ответственных
	DateCreate           string                `json:"date_create"`            // Дата создания
	SalesTypeID          string                `json:"sales_type_id"`          // ID типа активности
	SaleStageID          string                `json:"sale_stage_id"`          // ID стадии
	SaleName             string                `json:"sale_name"`              // Название сделки
	SaleActivityType     string                `json:"sale_activity_type"`     // Тип последней активности
	SaleActivityDate     string                `json:"sale_activity_date"`     // Дата последней активности сделк
	Fields               map[string]*SaleField `json:"fields"`                 // Данные полей
}
type SaleField struct {
	DataType string `json:"datatype"`
	Value    any    `json:"value"`
}

// Методы получения значений Sale

func (s *Sale) GetSaleID() uint64 {
	r, err := strconv.ParseUint(s.ID, 10, 64)
	if err != nil {
		return 0
	}

	return r
}

func (s *Sale) GetCustomersID() uint64 {
	r, err := strconv.ParseUint(s.CustomersID, 10, 64)
	if err != nil {
		return 0
	}

	return r
}

func (s *Sale) GetEmployeeID() uint64 {
	r, err := strconv.ParseUint(s.EmployeeID, 10, 64)
	if err != nil {
		return 0
	}

	return r
}

func (s *Sale) GetAdditionalEmployeeID() []uint64 {
	r := make([]uint64, 0, len(s.AdditionalEmployeeID))
	for _, v := range s.AdditionalEmployeeID {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			continue
		}

		r = append(r, id)
	}

	return r
}

func (s *Sale) GetDateCreate() time.Time {
	r, err := time.Parse(datetimeLayout, s.DateCreate)
	if err != nil {
		return time.Time{}
	}

	return r
}

func (s *Sale) GetSalesTypeID() uint16 {
	r, err := strconv.ParseUint(s.SalesTypeID, 10, 64)
	if err != nil {
		return 0
	}

	return uint16(r)
}

func (s *Sale) GetSaleStageID() uint16 {
	r, err := strconv.ParseUint(s.SaleStageID, 10, 64)
	if err != nil {
		return 0
	}

	return uint16(r)
}

func (s *Sale) GetSaleName() string {
	return s.SaleName
}

func (s *Sale) GetSaleActivityType() string {
	return s.SaleActivityType
}

func (s *Sale) GetSaleActivityDate() time.Time {
	r, err := time.Parse(datetimeLayout, s.SaleActivityDate)
	if err != nil {
		return time.Time{}
	}

	return r
}

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
