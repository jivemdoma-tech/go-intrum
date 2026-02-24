package gointrum

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	StockFilterResponse struct {
		Status  string          `json:"status,omitempty"`
		Message string          `json:"message,omitempty"`
		Data    StockFilterData `json:"data,omitempty"`
	}
	StockFilterData struct {
		List []Stock `json:"list"`
		// Count bool `json:"count"` // TODO: Реализовать. Проблема: может быть int или bool
	}
	Stock struct {
		ID                   int64                `json:"id,string"`                // ID объекта
		Type                 int64                `json:"type,string"`              // ID типа объекта
		Category             int64                `json:"parent,string"`            // ID категории
		Name                 string               `json:"name"`                     // Название
		DateCreate           time.Time            `json:"date_add"`                 // Дата создания
		StockCreatorID       int64                `json:"stock_creator_id,string"`  // ID создателя
		EmployeeID           int64                `json:"employee_id,string"`       // ID гл. ответственного
		AdditionalEmployeeID []int64              `json:"additional_employee_id"`   // Массив ID доп. ответственных
		LastModify           time.Time            `json:"last_modify"`              // Дата последнего редактирования
		CustomerRelation     int64                `json:"customer_relation,string"` // ID прикрепленного контакта
		StockActivityType    string               `json:"stock_activity_type"`      // Тип последней активности
		StockActivityDate    time.Time            `json:"stock_activity_date"`      // Дата последней активности
		Publish              bool                 `json:"publish"`                  // Активен или удален
		Fields               map[int64]StockField `json:"fields"`                   // Поля

		// TODO: Оставшиеся поля.
		//  Count any `json:"count"`
		//  Log any `json:"log"`
		//  Copy int64 `json:"copy,string"`
		//  GroupID int64 `json:"group_id,string"`
	}
	StockField struct {
		ID    int64  `json:"id,string"`
		Type  string `json:"type"`
		Value any    `json:"value"`
	}
)

func (r *StockFilterResponse) GetErrorMessage() string {
	switch {
	case r == nil:
		return ""
	case r.Status == "" && r.Message == "":
		return ""
	default:
		return fmt.Sprintf("%s: %s", r.Status, r.Message)
	}
}

// UnmarshalJSON обрабатывает неподдерживаемые стандартным UnmarshalJSON форматы полей.
func (s *Stock) UnmarshalJSON(data []byte) error {
	// alias - обертка над оригинальной структурой для предотвращения рекурсии
	type alias Stock

	// Вспомогательная структура (Приведение типа к alias)
	aux := &struct {
		*alias

		// Нужные поля

		// Дата + время
		DateCreate        string `json:"date_add"`
		LastModify        string `json:"last_modify"`
		StockActivityDate string `json:"stock_activity_date"`
		// Bool
		Publish string `json:"publish"`
		// Массивы
		AdditionalEmployeeID []string     `json:"additional_employee_id"`
		Fields               []StockField `json:"fields"`
	}{alias: (*alias)(s)}

	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Обработка полей и добавление в оригинальную структуру

	// Дата + время

	newDateCreate, err := time.Parse(DatetimeLayout, aux.DateCreate)
	switch err {
	case nil:
		s.DateCreate = newDateCreate
	default:
		s.DateCreate = time.Time{}
	}

	newLastModify, err := time.Parse(DatetimeLayout, aux.LastModify)
	switch err {
	case nil:
		s.LastModify = newLastModify
	default:
		s.LastModify = time.Time{}
	}

	newStockActivityDate, err := time.Parse(DatetimeLayout, aux.StockActivityDate)
	switch err {
	case nil:
		s.StockActivityDate = newStockActivityDate
	default:
		s.StockActivityDate = time.Time{}
	}

	// Bool

	newPublish, err := strconv.ParseBool(aux.Publish)
	switch err {
	case nil:
		s.Publish = newPublish
	default:
		s.Publish = false
	}

	// Массивы

	newAdditionalEmployeeID := make([]int64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if value, err := strconv.ParseInt(v, 10, 64); err == nil {
			newAdditionalEmployeeID = append(newAdditionalEmployeeID, value)
		}
	}
	s.AdditionalEmployeeID = newAdditionalEmployeeID

	// Костыль: для некоторых типов полей Интрум передает значения несколько раз по идентичному ключу
	newFields, alreadyParsedFields := make(map[int64]StockField, len(aux.Fields)), make(map[int64]struct{})
	for _, f := range aux.Fields {
		// Реализация костыля
		switch f.Type {
		case "file", "attach":
			// Проверка, что поле уже обработано
			if _, ok := alreadyParsedFields[f.ID]; ok {
				continue
			}
			alreadyParsedFields[f.ID] = struct{}{}
			// Сбор значений по ключу во всех полях
			collectedValues := make([]string, 0)
			for _, ff := range aux.Fields {
				if f.ID == ff.ID {
					vStr, _ := ff.Value.(string)
					collectedValues = append(collectedValues, vStr)
				}
			}
			f.Value = strings.Join(collectedValues, ",")
		}
		newFields[f.ID] = f
	}
	s.Fields = newFields

	return nil
}

// getField получает поле по ID.
func (s *Stock) getField(fieldID int64) *StockField {
	if s == nil {
		return nil
	}

	switch f, ok := s.Fields[fieldID]; {
	case ok:
		return &f
	default:
		return nil
	}
}

func (s *Stock) getFieldMap(fieldID int64) map[string]string {
	if s == nil {
		return nil
	}

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
func (s *Stock) GetFieldText(fieldID int64) string {
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
func (s *Stock) GetFieldRadio(fieldID int64) bool {
	vStr := s.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// Тип поля: "select".
func (s *Stock) GetFieldSelect(fieldID int64) string {
	return s.GetFieldText(fieldID)
}

// Тип поля: "multiselect".
func (s *Stock) GetFieldMultiselect(fieldID int64) []string {
	if vStr := s.GetFieldText(fieldID); vStr != "" {
		return strings.Split(vStr, ",")
	}
	return nil
}

// Тип поля: "date".
func (s *Stock) GetFieldDate(fieldID int64) time.Time {
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
func (s *Stock) GetFieldDatetime(fieldID int64) time.Time {
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
func (s *Stock) GetFieldTime(fieldID int64) time.Time {
	vStr := s.GetFieldText(fieldID)
	return parseTime(vStr, TimeLayout)
}

// Тип поля: "integer".
func (s *Stock) GetFieldInteger(fieldID int64) int64 {
	vStr := s.GetFieldText(fieldID)
	return parseInt(vStr)
}

// Тип поля: "decimal".
func (s *Stock) GetFieldDecimal(fieldID int64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// Тип поля: "price".
func (s *Stock) GetFieldPrice(fieldID int64) float64 {
	vStr := s.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// Тип поля: "file".
func (s *Stock) GetFieldFile(fieldID int64) string {
	return s.GetFieldText(fieldID)
}

// Тип поля: "point".
func (s *Stock) GetFieldPoint(fieldID int64) [2]string {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// Тип поля: "integer_range".
func (s *Stock) GetFieldIntegerRange(fieldID int64) [2]int64 {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// Тип поля: "decimal_range".
func (s *Stock) GetFieldDecimalRange(fieldID int64) [2]float64 {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// Тип поля: "date_range".
func (s *Stock) GetFieldDateRange(fieldID int64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, DateLayout)
	})
}

// Тип поля: "time_range".
func (s *Stock) GetFieldTimeRange(fieldID int64) [2]time.Time {
	m := s.getFieldMap(fieldID)
	if m == nil {
		return [2]time.Time{}
	}
	return parseRange(m, func(s string) time.Time {
		return parseTime(s, TimeLayout)
	})
}

// Тип поля: "datetime_range".
func (s *Stock) GetFieldDatetimeRange(fieldID int64) [2]time.Time {
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
func (s *Stock) GetFieldAttach(fieldID int64) []int64 {
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

func (s *Stock) GetFieldString(fieldID int64) string {
	return s.GetFieldText(fieldID)
}

func (s *Stock) GetFieldFloat(fieldID int64) float64 {
	return s.GetFieldDecimal(fieldID)
}

func (s *Stock) GetFieldFloatRange(fieldID int64) [2]float64 {
	return s.GetFieldDecimalRange(fieldID)
}

func (s *Stock) GetFieldBool(fieldID int64) bool {
	return s.GetFieldRadio(fieldID)
}
