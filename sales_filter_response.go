package gointrum

import (
	"encoding/json"
	"fmt"
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
	ID                   int64                 `json:"id,string,omitempty"`              // ID сделки
	CustomersID          int64                 `json:"customers_id,string,omitempty"`    // ID контакта
	EmployeeID           int64                 `json:"employee_id,string,omitempty"`     // ID ответственного
	AdditionalEmployeeID []int64               `json:"additional_employee_id,omitempty"` // Массив ID доп. ответственных
	DateCreate           time.Time             `json:"date_create,omitempty"`            // Дата создания
	SaleTypeID           int64                 `json:"sale_type_id,string,omitempty"`    // ID типа активности
	SaleStageID          int64                 `json:"sale_stage_id,string,omitempty"`   // ID стадии
	SaleName             string                `json:"sale_name,omitempty"`              // Название сделки
	SaleActivityType     string                `json:"sale_activity_type,omitempty"`     // Тип последней активности
	SaleActivityDate     time.Time             `json:"sale_activity_date,omitempty"`     // Дата последней активности сделк
	Fields               map[string]*SaleField `json:"fields,omitempty"`                 // Данные полей
	Publish              bool                  `json:"publish,omitempty"`                // Опубликован/Удален
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
		// Bool
		Publish string `json:"publish,omitempty"`
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

	newSlice := make([]int64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if newValue, err := strconv.ParseInt(v, 10, 64); err == nil {
			newSlice = append(newSlice, newValue)
		}
	}
	s.AdditionalEmployeeID = newSlice

	parsedBool, err := strconv.ParseBool(aux.Publish)
	switch err {
	case nil:
		s.Publish = parsedBool
	default:
		s.Publish = false
	}

	return nil
}

// Методы получения значений полей

// getField получает структуру поля по ID.
func (s *Sale) getField(fieldID int64) *SaleField {
	fieldIDStr := strconv.FormatInt(fieldID, 10)
	if f, exists := s.Fields[fieldIDStr]; exists {
		return f
	}
	return nil
}

func (s *Sale) getFieldMap(fieldID int64) map[string]string {
	f := s.getField(fieldID)
	if f == nil {
		return nil
	}
	switch m := f.Value.(type) {
	case map[string]string:
		return m
	case map[string]any:
		mStr := make(map[string]string, len(m))
		for k, v := range m {
			mStr[k] = fmt.Sprint(v)
		}
		return mStr
	}
	return nil
}

// Публичные методы

// Тип поля: "text".
func (s *Sale) GetFieldText(fieldID int64) string {
	f := s.getField(fieldID)
	if f == nil {
		return ""
	}
	vStr, ok := f.Value.(string)
	if !ok {
		return ""
	}
	return vStr
}

// Тип поля: "radio".
func (s *Sale) GetFieldRadio(fieldID int64) bool {
	vStr := s.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// Тип поля: "select".
func (s *Sale) GetFieldSelect(fieldID int64) string {
	return s.GetFieldText(fieldID)
}

// Тип поля: "multiselect".
func (s *Sale) GetFieldMultiselect(fieldID int64) []string {
	if vStr := s.GetFieldText(fieldID); vStr != "" {
		return strings.Split(vStr, ",")
	}
	return nil
}

// Тип поля: "date".
func (s *Sale) GetFieldDate(fieldID int64) time.Time {
	vStr := s.GetFieldText(fieldID)
	// Проверка на формат date
	if vDate := parseTime(vStr, DateLayout); !vDate.IsZero() {
		return vDate
	}
	// Проверка на формат datetime
	if vDatetime := parseTime(vStr, DatetimeLayout); !vDatetime.IsZero() {
		return time.Date(vDatetime.Year(), vDatetime.Month(), vDatetime.Day(), 0, 0, 0, 0, vDatetime.Location())
	}

	return time.Time{}
}

// Тип поля: "datetime".
func (s *Sale) GetFieldDatetime(fieldID int64) time.Time {
	vStr := s.GetFieldText(fieldID)
	// Проверка на формат datetime
	if vDatetime := parseTime(vStr, DatetimeLayout); !vDatetime.IsZero() {
		return vDatetime
	}
	// Проверка на формат date
	if vDate := parseTime(vStr, DateLayout); !vDate.IsZero() {
		return vDate
	}

	return time.Time{}
}

// Тип поля: "time".
func (s *Sale) GetFieldTime(fieldID int64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, TimeLayout)
}

// Тип поля: "integer".
func (s *Sale) GetFieldInteger(fieldID int64) int64 {
	vStr := s.GetFieldText(fieldID)
	return parseInt(vStr)
}

// Тип поля: "decimal".
func (s *Sale) GetFieldDecimal(fieldID int64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// Тип поля: "price".
func (s *Sale) GetFieldPrice(fieldID int64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// Тип поля: "file".
func (s *Sale) GetFieldFile(fieldID int64) string {
	return s.GetFieldText(fieldID)
}

// Тип поля: "point".
func (s *Sale) GetFieldPoint(fieldID int64) [2]string {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// Тип поля: "integer_range".
func (s *Sale) GetFieldIntegerRange(fieldID int64) [2]int64 {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// Тип поля: "decimal_range".
func (s *Sale) GetFieldDecimalRange(fieldID int64) [2]float64 {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// Тип поля: "date_range".
func (s *Sale) GetFieldDateRange(fieldID int64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DateLayout)
	})
}

// Тип поля: "time_range".
func (s *Sale) GetFieldTimeRange(fieldID int64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, TimeLayout)
	})
}

// Тип поля: "datetime_range".
func (s *Sale) GetFieldDatetimeRange(fieldID int64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DatetimeLayout)
	})
}

// Тип поля: "attach".
//
//	! ВНИМАНИЕ ! Возвращает ID только последней прикрепленной сущности.
func (s *Sale) GetFieldAttach(fieldID int64) []int64 {
	// TODO: Подружить метод с кривым API Интрума...
	f := s.getField(fieldID)
	if f == nil {
		return nil
	}
	m, ok := f.Value.(map[string]any)
	if !ok {
		return nil
	}
	idRaw, ok := m["id"]
	if !ok || idRaw == nil {
		return nil
	}
	switch id := idRaw.(type) {
	case string:
		if id == "" || id == "0" {
			return nil
		}
		val, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil
		}
		return []int64{val}
	}
	return nil
}

// Обертки методов с боле привычными названиями

func (s *Sale) GetFieldString(fieldID int64) string {
	return s.GetFieldText(fieldID)
}

func (s *Sale) GetFieldFloat(fieldID int64) float64 {
	return s.GetFieldDecimal(fieldID)
}

func (s *Sale) GetFieldFloatRange(fieldID int64) [2]float64 {
	return s.GetFieldDecimalRange(fieldID)
}

func (s *Sale) GetFieldBool(fieldID int64) bool {
	return s.GetFieldRadio(fieldID)
}
