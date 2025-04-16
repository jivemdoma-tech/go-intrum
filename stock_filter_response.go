package gointrum

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type StockFilterResponse struct {
	*Response
	Data   *StockFilterData `json:"data,omitempty"`
}
type StockFilterData struct {
	List []*Stock `json:"list"`
	// Count bool               `json:"count"`
}
type Stock struct {
	ID                   uint64                 `json:"id,string"`
	StockType            uint16                 `json:"stock_type,string"`
	Type                 uint16                 `json:"type,string"`
	Parent               uint16                 `json:"parent,string"`
	Name                 string                 `json:"name"`
	DateAdd              time.Time              `json:"date_add"` // TODO
	Count                bool                   `json:"count"`
	Author               uint64                 `json:"author,string"`
	EmployeeID           uint64                 `json:"employee_id,string"`
	AdditionalAuthor     []uint64               `json:"additional_author"`
	AdditionalEmployeeID []uint64               `json:"additional_employee_id"`
	LastModify           time.Time              `json:"last_modify"` // TODO
	CustomerRelation     uint64                 `json:"customer_relation,string"`
	StockActivityType    string                 `json:"stock_activity_type"`
	StockActivityDate    time.Time              `json:"stock_activity_date"` // TODO
	Publish              bool                   `json:"publish"`
	Copy                 uint64                 `json:"copy,string"`
	GroupID              uint16                 `json:"group_id,string"`
	StockCreatorID       uint64                 `json:"stock_creator_id,string"`
	Fields               map[uint64]*StockField `json:"fields"`
	// Log                  interface{}       `json:"log"`
}
type StockField struct {
	ID    uint64 `json:"id,string"`
	Type  string `json:"type"`
	Value any    `json:"value"`
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
		DateAdd           string `json:"date_add"`
		LastModify        string `json:"last_modify"`
		StockActivityDate string `json:"stock_activity_date"`
		// Bool
		Count   string `json:"count"`
		Publish string `json:"publish"`
		// Массивы
		AdditionalAuthor     []string      `json:"additional_author"`
		AdditionalEmployeeID []string      `json:"additional_employee_id"`
		Fields               []*StockField `json:"fields"`
	}{
		Alias: (*Alias)(s), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена дата + время

	parsedDate, err := time.Parse(datetimeLayout, aux.DateAdd)
	switch err {
	case nil:
		s.DateAdd = parsedDate
	default:
		s.DateAdd = time.Time{}
	}

	parsedDate, err = time.Parse(datetimeLayout, aux.LastModify)
	switch err {
	case nil:
		s.LastModify = parsedDate
	default:
		s.LastModify = time.Time{}
	}

	parsedDate, err = time.Parse(datetimeLayout, aux.StockActivityDate)
	switch err {
	case nil:
		s.StockActivityDate = parsedDate
	default:
		s.StockActivityDate = time.Time{}
	}

	// Замена bool

	parsedBool, err := strconv.ParseBool(aux.Count)
	switch err {
	case nil:
		s.Count = parsedBool
	default:
		s.Count = false
	}

	parsedBool, err = strconv.ParseBool(aux.Publish)
	switch err {
	case nil:
		s.Publish = parsedBool
	default:
		s.Publish = false
	}

	// Замена массивов

	newSlice := make([]uint64, 0, len(aux.AdditionalAuthor))
	for _, v := range aux.AdditionalAuthor {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	s.AdditionalAuthor = newSlice

	newSlice = make([]uint64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	s.AdditionalEmployeeID = newSlice

	newMap := make(map[uint64]*StockField, len(aux.Fields))
	for _, v := range aux.Fields {
		newMap[v.ID] = v
	}
	s.Fields = newMap

	return nil
}

// Методы получения значений Stock

// Вспомогательная функция получения структуры поля
func (s *Stock) getField(fieldID uint64) (*StockField, bool) {
	f, exists := s.Fields[fieldID]
	return f, exists
}

func (s *Stock) getFieldMap(fieldID uint64) (map[string]string, bool) {
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
func (s *Stock) GetFieldText(fieldID uint64) string {
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
func (s *Stock) GetFieldRadio(fieldID uint64) bool {
	vStr := s.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// select
func (s *Stock) GetFieldSelect(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// multiselect
func (s *Stock) GetFieldMultiselect(fieldID uint64) []string {
	return strings.Split(s.GetFieldText(fieldID), ",")
}

// date
func (s *Stock) GetFieldDate(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, dateLayout)
}

// datetime
func (s *Stock) GetFieldDatetime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, datetimeLayout)
}

// time
func (s *Stock) GetFieldTime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, timeLayout)
}

// integer
func (s *Stock) GetFieldInteger(fieldID uint64) int64 {
	vStr := s.GetFieldText(fieldID)
	return parseInt(vStr)
}

// decimal
func (s *Stock) GetFieldDecimal(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// price
func (s *Stock) GetFieldPrice(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// file
func (s *Stock) GetFieldFile(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// point
func (s *Stock) GetFieldPoint(fieldID uint64) [2]string {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// integer_range
func (s *Stock) GetFieldIntegerRange(fieldID uint64) [2]int64 {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// decimal_range
func (s *Stock) GetFieldDecimalRange(fieldID uint64) [2]float64 {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// date_range
func (s *Stock) GetFieldDateRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, dateLayout)
	})
}

// time_range
func (s *Stock) GetFieldTimeRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, dateLayout)
	})
}

// datetime_range
func (s *Stock) GetFieldDatetimeRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, dateLayout)
	})
}

// attach
func (s *Stock) GetFieldAttach(fieldID uint64) []uint64 {
	f, exists := s.getField(fieldID)
	if !exists {
		return nil
	}
	vAttach, ok := f.Value.([]map[string]string)
	if !ok || len(vAttach) <= 0 {
		return nil
	}
	vIDs := make([]uint64, 0, len(vAttach))
	for _, v := range vAttach {
		if id, err := strconv.ParseUint(v["id"], 10, 64); err == nil {
			vIDs = append(vIDs, id)
		}
	}
	return vIDs
}
