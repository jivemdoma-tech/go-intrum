package gointrum

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type ApplicationFilterResponse struct {
	Status string                 `json:"status"`
	Data   *ApplicationFilterData `json:"data"`
}

type ApplicationFilterData struct {
	List  []*Application `json:"list"`
	Count bool           `json:"count"`
}

type Application struct {
	ID                   uint64                       `json:"id,string"`
	Publish              bool                         `json:"publish,string"`
	EmployeeID           uint64                       `json:"employee_id,string"`
	CustomerID           uint64                       `json:"customer_id,string"`
	VisitID              uint64                       `json:"visit_id,string"`
	RequestTypeID        uint16                       `json:"request_type_id,string"`
	RequestTypeName      string                       `json:"request_type_name"`
	Source               string                       `json:"source"`
	DateCreate           time.Time                    `json:"date_create"`
	Comment              string                       `json:"comment"`
	RequestName          string                       `json:"request_name"`
	Status               string                       `json:"status"`
	RequestActivityType  string                       `json:"request_activity_type"`
	RequestActivityDate  time.Time                    `json:"request_activity_date"`
	RequestCreatorID     *uint64                      `json:"request_creator_id,string"`
	AdditionalEmployeeID []any                        `json:"additional_employee_id"`
	Fields               map[string]*ApplicationField `json:"fields,omitempty"`
}

type ApplicationField struct {
	Datatype string `json:"datatype"`
	Value    any    `json:"value"`
}

func (s *Application) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias Application

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		DateCreate          string `json:"date_create"`
		RequestActivityDate string `json:"request_activity_date"`
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

	parsedDate, err = time.Parse(datetimeLayout, aux.RequestActivityDate)
	if err != nil {
		return err
	}
	s.RequestActivityDate = parsedDate

	return nil
}

// Методы получения значений Application

// Вспомогательная функция получения структуры поля
func (s *Application) getField(fieldID uint64) (*ApplicationField, bool) {
	f, exists := s.Fields[strconv.FormatUint(fieldID, 10)]
	return f, exists
}

func (s *Application) getFieldMap(fieldID uint64) (map[string]string, bool) {
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
func (s *Application) GetFieldText(fieldID uint64) string {
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
func (s *Application) GetFieldRadio(fieldID uint64) bool {
	vStr := s.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// select
func (s *Application) GetFieldSelect(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// multiselect
func (s *Application) GetFieldMultiselect(fieldID uint64) []string {
	return strings.Split(s.GetFieldText(fieldID), ",")
}

// date
func (s *Application) GetFieldDate(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, dateLayout)
}

// datetime
func (s *Application) GetFieldDatetime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, datetimeLayout)
}

// time
func (s *Application) GetFieldTime(fieldID uint64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, timeLayout)
}

// integer
func (s *Application) GetFieldInteger(fieldID uint64) int64 {
	vStr := s.GetFieldText(fieldID)
	return parseInt(vStr)
}

// decimal
func (s *Application) GetFieldDecimal(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// price
func (s *Application) GetFieldPrice(fieldID uint64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// file
func (s *Application) GetFieldFile(fieldID uint64) string {
	return s.GetFieldText(fieldID)
}

// point
func (s *Application) GetFieldPoint(fieldID uint64) [2]string {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// integer_range
func (s *Application) GetFieldIntegerRange(fieldID uint64) [2]int64 {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// decimal_range
func (s *Application) GetFieldDecimalRange(fieldID uint64) [2]float64 {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// date_range
func (s *Application) GetFieldDateRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, dateLayout)
	})
}

// time_range
func (s *Application) GetFieldTimeRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, dateLayout)
	})
}

// datetime_range
func (s *Application) GetFieldDatetimeRange(fieldID uint64) [2]time.Time {
	m, ok := s.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, dateLayout)
	})
}

// attach
func (s *Application) GetFieldAttach(fieldID uint64) []uint64 {
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
