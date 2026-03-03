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
	default:
		return ""
	case r.Status != "" && r.Message != "":
		return r.Status + ": " + r.Message
	case r.Message != "":
		return r.Message
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

	// (bool) Publish
	newPublish, _ := strconv.ParseBool(aux.Publish)
	s.Publish = newPublish

	// ([]int64) AdditionalEmployeeID
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
			// Проверка: поле уже обработано
			switch _, ok := alreadyParsedFields[f.ID]; {
			case ok:
				continue
			default:
				alreadyParsedFields[f.ID] = struct{}{}
			}
			// Сбор значений
			collectedValues := make([]string, 0)
			for _, ff := range aux.Fields {
				if f.ID != ff.ID {
					continue
				}
				if ff.Value == nil {
					continue
				}
				// Строка
				if valueStr, ok := ff.Value.(string); ok && valueStr != "" {
					collectedValues = append(collectedValues, valueStr)
				}
				// Хэш-таблица
				if valueMap, ok := ff.Value.(map[string]any); ok && valueMap != nil {
					if valueMapID, ok := valueMap["id"]; ok && valueMapID != nil {
						collectedValues = append(collectedValues, fmt.Sprint(valueMapID))
					}
				}
			}
			f.Value = strings.Join(collectedValues, ",")
		}
		newFields[f.ID] = f
	}
	s.Fields = newFields

	return nil
}

// getField возвращает поле по id.
func (s *Stock) getField(id int64) (*StockField, bool) {
	if s == nil {
		return nil, false
	}
	// Проверка: поле существует
	field, exists := s.Fields[id]
	if !exists {
		return nil, false
	}

	return &field, true
}

// getFieldMap возвращает значение поля (map[string]string) по id.
func (s *Stock) getFieldMap(id int64) (map[string]string, bool) {
	// Проверка: поле существует
	field, exists := s.getField(id)
	if !exists {
		return nil, false
	}
	value := field.Value

	if resultMap, ok := value.(map[string]string); ok {
		return resultMap, true
	}
	if rawMap, ok := value.(map[string]any); ok {
		resultMap := make(map[string]string, len(rawMap))
		for k, v := range rawMap {
			switch v {
			default:
				resultMap[k] = fmt.Sprint(v)
			case nil:
				resultMap[k] = ""
			}
		}
		return resultMap, true
	}

	return nil, false
}

// GetFieldText возвращает значение поля (text) по id.
func (s *Stock) GetFieldText(id int64) (string, bool) {
	// Проверка: поле существует
	field, exists := s.getField(id)
	if !exists {
		return "", false
	}
	value := field.Value

	switch value {
	case nil:
		return "", true
	default:
		return fmt.Sprint(value), true
	}
}

// GetFieldRadio возвращает значение поля (radio) по id.
func (s *Stock) GetFieldRadio(id int64) (bool, bool) {
	// Проверка: поле существует
	valueStr, exists := s.GetFieldText(id)
	if !exists {
		return false, false
	}

	// Поле существует
	valueBool, _ := strconv.ParseBool(valueStr)
	return valueBool, true
}

// GetFieldSelect возвращает значение поля (select) по id.
func (s *Stock) GetFieldSelect(id int64) (string, bool) { return s.GetFieldText(id) }

// GetFieldMultiselect возвращает значение поля (multiselect) по id.
func (s *Stock) GetFieldMultiselect(id int64) ([]string, bool) {
	// Проверка: поле существует
	valueStr, exists := s.GetFieldText(id)
	if !exists {
		return nil, false
	}

	switch valueStr {
	case "":
		return nil, true
	default:
		return strings.Split(valueStr, ","), true
	}
}

// GetFieldDate возвращает значение поля (date) по id.
func (s *Stock) GetFieldDate(id int64) (time.Time, bool) {
	// Проверка: поле существует
	valueStr, exists := s.GetFieldText(id)
	if !exists {
		return time.Time{}, false
	}

	return parseDate(valueStr), true
}

// GetFieldDatetime возвращает значение поля (datetime) по id.
func (s *Stock) GetFieldDatetime(id int64) (time.Time, bool) {
	// Проверка: поле существует
	valueStr, exists := s.GetFieldText(id)
	if !exists {
		return time.Time{}, false
	}

	return parseDatetime(valueStr), true
}

// TODO: GetFieldTime
// // GetFieldTime возвращает значение поля (time) по id.
// func (s *Stock) GetFieldTime(id int64) (time.Time, bool) {}

// GetFieldInteger возвращает значение поля (integer) по id.
func (s *Stock) GetFieldInteger(id int64) (int64, bool) {
	// Проверка: поле существует
	valueStr, exists := s.GetFieldText(id)
	if !exists {
		return 0, false
	}

	return parseInt(valueStr), true
}

// GetFieldDecimal возвращает значение поля (decimal) по id.
func (s *Stock) GetFieldDecimal(id int64) (float64, bool) {
	// Проверка: поле существует
	valueStr, exists := s.GetFieldText(id)
	if !exists {
		return 0, false
	}

	return parseFloat(valueStr), true
}

// GetFieldPrice возвращает значение поля (price) по id.
func (s *Stock) GetFieldPrice(id int64) (float64, bool) { return s.GetFieldDecimal(id) }

// GetFieldFile возвращает значение поля (file) по id.
func (s *Stock) GetFieldFile(id int64) ([]string, bool) { return s.GetFieldMultiselect(id) }

// GetFieldPoint возвращает значение поля (point) по id.
func (s *Stock) GetFieldPoint(id int64) ([2]string, bool) {
	// Проверка: поле существует
	valueMap, exists := s.getFieldMap(id)
	if !exists {
		return [2]string{}, false
	}

	var (
		x, _ = valueMap["x"]
		y, _ = valueMap["y"]
	)
	return [2]string{x, y}, true
}

// GetFieldIntegerRange возвращает значение поля (integer_range) по id.
func (s *Stock) GetFieldIntegerRange(id int64) ([2]int64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.getFieldMap(id)
	if !exists {
		return [2]int64{}, false
	}

	return parseRange(valueMap, parseInt), true
}

// GetFieldDecimalRange возвращает значение поля (decimal_range) по id.
func (s *Stock) GetFieldDecimalRange(id int64) ([2]float64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.getFieldMap(id)
	if !exists {
		return [2]float64{}, false
	}

	return parseRange(valueMap, parseFloat), true
}

// GetFieldDateRange возвращает значение поля (date_range) по id.
func (s *Stock) GetFieldDateRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.getFieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return parseRange(valueMap, parseDate), true
}

// TODO: GetFieldTimeRange
// // GetFieldTimeRange возвращает значение поля (time_range) по id.
// func (s *Stock) GetFieldTimeRange(id int64) [2]time.Time {}

// GetFieldDatetimeRange возвращает значение поля (datetime_range) по id.
func (s *Stock) GetFieldDatetimeRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.getFieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return parseRange(valueMap, parseDatetime), true
}

// GetFieldAttach возвращает значение поля (attach) по id.
func (s *Stock) GetFieldAttach(id int64) ([]int64, bool) {
	// Проверка: поле существует
	fieldStringSlice, exists := s.GetFieldMultiselect(id)
	if !exists {
		return nil, false
	}

	valueInt64Slice := make([]int64, 0, len(fieldStringSlice))
	for _, vStr := range fieldStringSlice {
		if vInt64 := parseInt(vStr); vInt64 != 0 {
			valueInt64Slice = append(valueInt64Slice, vInt64)
		}
	}

	return valueInt64Slice, true
}
