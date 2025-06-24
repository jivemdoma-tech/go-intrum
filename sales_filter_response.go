package gointrum

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type SalesFilterResponse struct {
	*Response
	Data *SalesFilterData `json:"data,omitempty"`
}
type SalesFilterData struct {
	List []*Sale `json:"list"`
	// Count any `json:"count"` // TODO
}
type Sale struct {
	ID                   uint64                `json:"id,string,omitempty"`              // ID сделки
	CustomersID          uint64                `json:"customers_id,string,omitempty"`    // ID контакта
	EmployeeID           uint64                `json:"employee_id,string,omitempty"`     // ID ответственного
	AdditionalEmployeeID []uint64              `json:"additional_employee_id,omitempty"` // Массив ID доп. ответственных
	DateCreate           time.Time             `json:"date_create,omitempty"`            // Дата создания
	SalesTypeID          uint64                `json:"sales_type_id,string,omitempty"`   // ID типа активности
	SaleStageID          uint64                `json:"sale_stage_id,string,omitempty"`   // ID стадии
	SaleName             string                `json:"sale_name,omitempty"`              // Название сделки
	SaleActivityType     string                `json:"sale_activity_type,omitempty"`     // Тип последней активности
	SaleActivityDate     time.Time             `json:"sale_activity_date,omitempty"`     // Дата последней активности сделк
	Fields               map[string]*SaleField `json:"fields,omitempty"`                 // Данные полей
}

// Использовать метод GetField для получения значения поля // TODO
type SaleField struct {
	DataType string `json:"datatype,omitempty"`
	Value    any    `json:"value,omitempty"`
}

func (s *Sale) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias Sale

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		AdditionalEmployeeID []string `json:"additional_employee_id,omitempty"`
		DateCreate           string   `json:"date_create,omitempty"`
		SaleActivityDate     string   `json:"sale_activity_date,omitempty"`
	}{
		Alias: (*Alias)(s), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена

	parsedDate, err := time.Parse(DatetimeLayout, aux.DateCreate)
	if err != nil {
		return err
	}
	s.DateCreate = parsedDate

	parsedDate, err = time.Parse(DatetimeLayout, aux.SaleActivityDate)
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

// Вспомогательная функция получения структуры поля
func (s *Sale) getField(fieldID uint64) (*SaleField, bool) {
	f, exists := s.Fields[strconv.FormatUint(fieldID, 10)]
	return f, exists
}

func (s *Sale) getFieldMap(fieldID uint64) (map[string]string, bool) {
	f, exists := s.getField(fieldID)
	if !exists {
		return nil, false
	}
	m, ok := f.Value.(map[string]string)
	if !ok {
		return nil, false
	}
	return m, true
}

// text
func (s *Sale) GetFieldText(fieldID uint64) string {
	f, exists := s.getField(fieldID)
	if !exists {
		return ""
	}
	vStr, ok := f.Value.(string)
	if !ok {
		return ""
	}
	return vStr
}

// radio
func (s *Sale) GetFieldRadio(fieldID uint64) bool {
	vStr := s.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// select
func (s *Sale) GetFieldSelect(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// multiselect
func (s *Sale) GetFieldMultiselect(fieldID uint64) []string {
	return strings.Split(s.GetFieldText(fieldID), ",")
}

// date
func (s *Sale) GetFieldDate(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, DateLayout)
}

// datetime
func (s *Sale) GetFieldDatetime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, DatetimeLayout)
}

// time
func (s *Sale) GetFieldTime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, TimeLayout)
}

// integer
func (s *Sale) GetFieldInteger(fieldID uint64) int64 {
	vStr := s.GetFieldText(fieldID)
	return parseInt(vStr)
}

// decimal
func (s *Sale) GetFieldDecimal(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// price
func (s *Sale) GetFieldPrice(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// file
func (s *Sale) GetFieldFile(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// point
func (s *Sale) GetFieldPoint(fieldID uint64) [2]string {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// integer_range
func (s *Sale) GetFieldIntegerRange(fieldID uint64) [2]int64 {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// decimal_range
func (s *Sale) GetFieldDecimalRange(fieldID uint64) [2]float64 {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// date_range
func (s *Sale) GetFieldDateRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DateLayout)
	})
}

// time_range
func (s *Sale) GetFieldTimeRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DateLayout)
	})
}

// datetime_range
func (s *Sale) GetFieldDatetimeRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DateLayout)
	})
}

// attach
func (s *Sale) GetFieldAttach(fieldID uint64) []uint64 {
	f, exists := s.getField(fieldID)
	if !exists {
		return nil
	}
	arr, ok := f.Value.([]interface{})
	if !ok || len(arr) == 0 {
		// fmt.Println(f.Value)
		return nil
	}
	vIDs := make([]uint64, 0, len(arr))
	for _, v := range arr {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		idStr, ok := m["id"].(string)
		if !ok {
			// Если по какой-то причине id пришел не строкой, можно попробовать float64 (стандартное поведение encoding/json)
			if idFloat, ok := m["id"].(float64); ok {
				vIDs = append(vIDs, uint64(idFloat))
				continue
			}
			continue
		}
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			vIDs = append(vIDs, id)
		}
	}
	return vIDs
}
