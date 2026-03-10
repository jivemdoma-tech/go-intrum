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

	// (time.Time) Поля типа datetime
	s.DateCreate = parseDatetime(aux.DateCreate)
	s.LastModify = parseDatetime(aux.LastModify)
	s.StockActivityDate = parseDatetime(aux.StockActivityDate)

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

// field возвращает поле по id.
func (s *Stock) field(id int64) (*StockField, bool) {
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

// fieldMap возвращает значение поля (map[string]string) по id.
func (s *Stock) fieldMap(id int64) (map[string]string, bool) {
	// Проверка: поле существует
	field, exists := s.field(id)
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

// FieldText возвращает значение поля (text) по id.
func (s *Stock) FieldText(id int64) (string, bool) {
	// Проверка: поле существует
	field, exists := s.field(id)
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

// FieldTextOrZero возвращает значение поля (text) по id.
func (s *Stock) FieldTextOrZero(id int64) string {
	result, _ := s.FieldText(id)
	return result
}

// FieldRadio возвращает значение поля (radio) по id.
func (s *Stock) FieldRadio(id int64) (bool, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return false, false
	}

	// Поле существует
	valueBool, _ := strconv.ParseBool(valueStr)
	return valueBool, true
}

// FieldRadioOrZero возвращает значение поля (radio) по id.
func (s *Stock) FieldRadioOrZero(id int64) bool {
	result, _ := s.FieldRadio(id)
	return result
}

// FieldSelect возвращает значение поля (select) по id.
func (s *Stock) FieldSelect(id int64) (string, bool) { return s.FieldText(id) }

// FieldSelectOrZero возвращает значение поля (select) по id.
func (s *Stock) FieldSelectOrZero(id int64) string { return s.FieldTextOrZero(id) }

// FieldMultiselect возвращает значение поля (multiselect) по id.
func (s *Stock) FieldMultiselect(id int64) ([]string, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
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

// FieldMultiselectOrZero возвращает значение поля (multiselect) по id.
func (s *Stock) FieldMultiselectOrZero(id int64) []string {
	result, _ := s.FieldMultiselect(id)
	return result
}

// FieldDate возвращает значение поля (date) по id.
func (s *Stock) FieldDate(id int64) (time.Time, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return time.Time{}, false
	}

	return parseDate(valueStr), true
}

// FieldDateOrZero возвращает значение поля (date) по id.
func (s *Stock) FieldDateOrZero(id int64) time.Time {
	result, _ := s.FieldDate(id)
	return result
}

// FieldDatetime возвращает значение поля (datetime) по id.
func (s *Stock) FieldDatetime(id int64) (time.Time, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return time.Time{}, false
	}

	return parseDatetime(valueStr), true
}

// FieldDatetimeOrZero возвращает значение поля (datetime) по id.
func (s *Stock) FieldDatetimeOrZero(id int64) time.Time {
	result, _ := s.FieldDatetime(id)
	return result
}

// FieldInteger возвращает значение поля (integer) по id.
func (s *Stock) FieldInteger(id int64) (int64, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return 0, false
	}

	return parseInt(valueStr), true
}

// FieldIntegerOrZero возвращает значение поля (integer) по id.
func (s *Stock) FieldIntegerOrZero(id int64) int64 {
	result, _ := s.FieldInteger(id)
	return result
}

// FieldDecimal возвращает значение поля (decimal) по id.
func (s *Stock) FieldDecimal(id int64) (float64, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return 0, false
	}

	return parseFloat(valueStr), true
}

// FieldDecimalOrZero возвращает значение поля (decimal) по id.
func (s *Stock) FieldDecimalOrZero(id int64) float64 {
	result, _ := s.FieldDecimal(id)
	return result
}

// FieldFilesOrZero возвращает значение поля (files) по id.
func (s *Stock) FieldFilesOrZero(id int64) []string {
	result, _ := s.FieldFile(id)
	return result
}

// FieldPrice возвращает значение поля (price) по id.
func (s *Stock) FieldPrice(id int64) (float64, bool) { return s.FieldDecimal(id) }

// FieldPriceOrZero возвращает значение поля (price) по id.
func (s *Stock) FieldPriceOrZero(id int64) float64 { return s.FieldDecimalOrZero(id) }

// FieldFile возвращает значение поля (file) по id.
func (s *Stock) FieldFile(id int64) ([]string, bool) { return s.FieldMultiselect(id) }

// FieldPoint возвращает значение поля (point) по id.
func (s *Stock) FieldPoint(id int64) ([2]string, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]string{}, false
	}

	var (
		x, _ = valueMap["x"]
		y, _ = valueMap["y"]
	)
	return [2]string{x, y}, true
}

// FieldPointOrZero возвращает значение поля (point) по id.
func (s *Stock) FieldPointOrZero(id int64) [2]string {
	result, _ := s.FieldPoint(id)
	return result
}

// FieldIntegerRange возвращает значение поля (integer_range) по id.
func (s *Stock) FieldIntegerRange(id int64) ([2]int64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]int64{}, false
	}

	return parseRange(valueMap, parseInt), true
}

// FieldDecimalRange возвращает значение поля (decimal_range) по id.
func (s *Stock) FieldDecimalRange(id int64) ([2]float64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]float64{}, false
	}

	return parseRange(valueMap, parseFloat), true
}

// FieldDateRange возвращает значение поля (date_range) по id.
func (s *Stock) FieldDateRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return parseRange(valueMap, parseDate), true
}

// FieldDatetimeRange возвращает значение поля (datetime_range) по id.
func (s *Stock) FieldDatetimeRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return parseRange(valueMap, parseDatetime), true
}

// FieldAttach возвращает значение поля (attach) по id.
func (s *Stock) FieldAttach(id int64) ([]int64, bool) {
	// Проверка: поле существует
	fieldStringSlice, exists := s.FieldMultiselect(id)
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
