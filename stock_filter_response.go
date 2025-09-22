package gointrum

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StockFilterResponse struct {
	*Response
	Data *StockFilterData `json:"data,omitempty"`
}
type StockFilterData struct {
	List []*Stock `json:"list"`
	// Count bool `json:"count"` // TODO Реализовать через кастомный UnmarshalJSON
}
type Stock struct {
	ID                   uint64                 `json:"id,string"`                          // ID объекта
	Type                 uint64                 `json:"type,string,omitempty"`              // ID типа объекта
	Category             uint64                 `json:"parent,string,omitempty"`            // ID категории
	Name                 string                 `json:"name,omitempty"`                     // Название
	DateCreate           time.Time              `json:"date_add,omitempty"`                 // Дата создания
	StockCreatorID       uint64                 `json:"stock_creator_id,string,omitempty"`  // ID создателя
	EmployeeID           uint64                 `json:"employee_id,string,omitempty"`       // ID гл. ответственного
	AdditionalEmployeeID []uint64               `json:"additional_employee_id,omitempty"`   // Массив ID доп. ответственных
	LastModify           time.Time              `json:"last_modify,omitempty"`              // Дата последнего редактирования
	CustomerRelation     uint64                 `json:"customer_relation,string,omitempty"` // ID прикрепленного контакта
	StockActivityType    string                 `json:"stock_activity_type,omitempty"`      // Тип последней активности
	StockActivityDate    time.Time              `json:"stock_activity_date,omitempty"`      // Дата последней активности
	Publish              bool                   `json:"publish,omitempty"`                  // Активен или удален
	Fields               map[uint64]*StockField `json:"fields,omitempty"`                   // Поля

	// TODO
	// Count any `json:"count,omitempty"`
	// Log any `json:"log,omitempty"`
	// Copy uint64 `json:"copy,string,omitempty"`
	// GroupID uint64 `json:"group_id,string,omitempty"`
}
type StockField struct {
	ID    uint64 `json:"id,string,omitempty"`
	Type  string `json:"type,omitempty"`
	Value any    `json:"value,omitempty"`
}

func (s *Stock) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type (
		Alias Stock
	)

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		// Дата + время
		DateCreate        string `json:"date_add,omitempty"`
		LastModify        string `json:"last_modify,omitempty"`
		StockActivityDate string `json:"stock_activity_date,omitempty"`
		// Bool
		Publish string `json:"publish,omitempty"`
		// Массивы
		AdditionalEmployeeID []string      `json:"additional_employee_id,omitempty"`
		Fields               []*StockField `json:"fields,omitempty"`
	}{
		Alias: (*Alias)(s), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена дата + время

	parsedDate, err := time.Parse(DatetimeLayout, aux.DateCreate)
	switch err {
	case nil:
		s.DateCreate = parsedDate
	default:
		s.DateCreate = time.Time{}
	}

	parsedDate, err = time.Parse(DatetimeLayout, aux.LastModify)
	switch err {
	case nil:
		s.LastModify = parsedDate
	default:
		s.LastModify = time.Time{}
	}

	parsedDate, err = time.Parse(DatetimeLayout, aux.StockActivityDate)
	switch err {
	case nil:
		s.StockActivityDate = parsedDate
	default:
		s.StockActivityDate = time.Time{}
	}

	// Замена bool

	parsedBool, err := strconv.ParseBool(aux.Publish)
	switch err {
	case nil:
		s.Publish = parsedBool
	default:
		s.Publish = false
	}

	// Замена массивов

	newSlice := make([]uint64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	s.AdditionalEmployeeID = newSlice

	var (
		newMap        = make(map[uint64]*StockField, len(aux.Fields))
		alreadyParsed = make(map[uint64]struct{}) // Костыль для сбора полей с дублирующимися ключами в 1 ключ
	)
	for _, f := range aux.Fields {
		// Реализация костыля
		switch f.Type {
		case "file", "attach":
			switch _, ok := alreadyParsed[f.ID]; {
			case ok:
				continue
			default:
				alreadyParsed[f.ID] = struct{}{}
			}
			// Прогон по всем полям с поиском значений под нашим ключом
			valuesCollected := make([]string, 0)
			for _, f2 := range aux.Fields {
				if f.ID != f2.ID {
					continue
				}
				vStr, _ := f2.Value.(string)
				valuesCollected = append(valuesCollected, vStr)
			}
			f.Value = strings.Join(valuesCollected, ",")
		}
		newMap[f.ID] = f
	}
	s.Fields = newMap

	return nil
}

// Методы получения значений полей

// getField получает структуру поля по ID.
func (s *Stock) getField(fieldID uint64) *StockField {
	if f, exists := s.Fields[fieldID]; exists {
		return f
	}
	return nil
}

func (s *Stock) getFieldMap(fieldID uint64) map[string]string {
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
func (s *Stock) GetFieldText(fieldID uint64) string {
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
func (s *Stock) GetFieldRadio(fieldID uint64) bool {
	vStr := s.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// Тип поля: "select".
func (s *Stock) GetFieldSelect(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// Тип поля: "multiselect".
func (s *Stock) GetFieldMultiselect(fieldID uint64) []string {
	if vStr := s.GetFieldText(fieldID); vStr != "" {
		return strings.Split(vStr, ",")
	}
	return nil
}

// Тип поля: "date".
func (s *Stock) GetFieldDate(fieldID uint64) time.Time {
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
func (s *Stock) GetFieldDatetime(fieldID uint64) time.Time {
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
func (s *Stock) GetFieldTime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, TimeLayout)
}

// Тип поля: "integer".
func (s *Stock) GetFieldInteger(fieldID uint64) int64 {
	vStr := s.GetFieldText(fieldID)
	return parseInt(vStr)
}

// Тип поля: "decimal".
func (s *Stock) GetFieldDecimal(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// Тип поля: "price".
func (s *Stock) GetFieldPrice(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// Тип поля: "file".
func (s *Stock) GetFieldFile(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// Тип поля: "point".
func (s *Stock) GetFieldPoint(fieldID uint64) [2]string {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// Тип поля: "integer_range".
func (s *Stock) GetFieldIntegerRange(fieldID uint64) [2]int64 {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// Тип поля: "decimal_range".
func (s *Stock) GetFieldDecimalRange(fieldID uint64) [2]float64 {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// Тип поля: "date_range".
func (s *Stock) GetFieldDateRange(fieldID uint64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DateLayout)
	})
}

// Тип поля: "time_range".
func (s *Stock) GetFieldTimeRange(fieldID uint64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, TimeLayout)
	})
}

// Тип поля: "datetime_range".
func (s *Stock) GetFieldDatetimeRange(fieldID uint64) [2]time.Time {
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
func (s *Stock) GetFieldAttach(fieldID uint64) []uint64 {
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
		val, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil
		}
		return []uint64{val}
	}
	return nil
}

// Обертки методов с боле привычными названиями

func (s *Stock) GetFieldString(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

func (s *Stock) GetFieldFloat(fieldID uint64) float64 {
	return s.GetFieldDecimal(fieldID)
}

func (s *Stock) GetFieldFloatRange(fieldID uint64) [2]float64 {
	return s.GetFieldDecimalRange(fieldID)
}

func (s *Stock) GetFieldBool(fieldID uint64) bool {
	return s.GetFieldRadio(fieldID)
}
