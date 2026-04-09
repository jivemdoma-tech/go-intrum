package intrum

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TODO: Реализовать оставшиеся поля StockFilterParams:
//  Nested
//  IndexFields
//  Order
//  OrderField
//  Date
//  DateField
//  GroupID
//  Copy
//  ObjectGroups
//  CountTotal
//  OnlyPrimaryID
//  OnlyCountField
//  SumField
//  Log
// TODO: Реализоваться оставшиеся поля в StockFilterResponse.StockFilterData:
//  Count
// TODO: Реализовать оставшиеся поля в StockFilterResponse.StockFilterData.Stock:
//  Count
//  Log
//  Copy
//  GroupID

const StockFilterMaxLimit int64 = 500

// StockFilter - поиск объектов в CRM. Документация: https://www.intrumnet.com/api/#stock-search
func StockFilter(ctx context.Context, subdomain, apiKey string, p *StockFilterParams) (*StockFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/filter", subdomain)

	// Валидация
	if err := validateRequestArgs(methodURL, subdomain, apiKey); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, newErrEmptyParams(methodURL)
	}

	// Обязательные поля
	if p.Type <= 0 && len(p.ByIDs) == 0 {
		return nil, newErrEmptyRequiredParams(methodURL)
	}

	// Запрос
	resp := &StockFilterResponse{}
	if err := request(ctx, apiKey, methodURL, p.params(), resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// StockFilterAll - поиск объектов в CRM по всем страницам. Документация: https://www.intrumnet.com/api/#stock-search
func StockFilterAll(ctx context.Context, subdomain, apiKey string, p *StockFilterParams) ([]Stock, error) {
	result := make([]Stock, 0, 500)
	for page := int64(1); ; page++ {
		// Shallow-копирование структуры для итерации
		pageParams := p.cloneWithPage(page)
		// Установка максимального кол-ва элементов в ответе
		if pageParams.Limit != StockFilterMaxLimit {
			pageParams.Limit = StockFilterMaxLimit
		}
		// Запрос
		resp, err := StockFilter(ctx, subdomain, apiKey, pageParams)
		if err != nil {
			return nil, err
		}
		stock := resp.Data.List

		if len(stock) == 0 {
			break
		}

		result = append(result, stock...)

		if len(stock) < int(pageParams.Limit) {
			break
		}
	}
	if len(result) == 0 {
		return nil, ErrNothingFound
	}

	return result, nil
}

// =====================================================================================================================
// Request
// =====================================================================================================================

// StockFilterParams - параметры запроса StockFilter.
//
// Обязательные поля:
//   - Type || ByIDs
//
// Основные параметры запроса:
//   - Type: ID типа объекта.
//   - Category: ID категории объекта.
//   - ByIDs: массив ID объектов. Все объекты в массиве должны быть одного типа.
//   - Publish: активность объектов. Активные (по умолчанию): "1". Удаленные: "0". Все: "ignore".
//   - Fields: массив ID полей и значений.
//     Для типов (integer, decimal, price, time, date, datetime) возможно указывать границы.
//     Больше или равно: ">= {ЗНАЧЕНИЕ}".
//     Меньше или равно: "<= {ЗНАЧЕНИЕ}".
//     Между двумя значениями: "{ЗНАЧЕНИЕ} & {ЗНАЧЕНИЕ}".
//
// Параметры ответа:
//   - SliceFields: массив ID полей, значения которых будут в ответе. По умолчанию выводятся все.
//   - Limit: кол-во объектов в ответе.
//   - Page: номер страницы ответа. Начинается с 1. Игнорируется StockFilterAll.
type StockFilterParams struct {
	Type           int64   // ID типа объекта.
	Category       int64   // ID категории объекта.
	ByIDs          []int64 // Массив ID объектов. Все объекты в массиве должны быть одного типа.
	Search         string  // Поисковая строка. Может содержать имся объекта или вхождения в поля с типами (text, select, multiselect).
	Manager        []int64 // Массив ID ответственных.
	Groups         []int64 // Массив ID CRM-групп.
	StockCreatorID int64   // ID создателя объекта.
	// Массив ID полей и значений.
	//  {ID ПОЛЯ}: "{ЗНАЧЕНИЕ}"
	// Для типов (integer, decimal, price, time, date, datetime) возможно указывать границы:
	//	Больше или равно: ">= {ЗНАЧЕНИЕ}"
	//	Меньше или равно: "<= {ЗНАЧЕНИЕ}"
	//	Между двумя значениям: "{ЗНАЧЕНИЕ} & {ЗНАЧЕНИЕ}"
	Fields              map[int64]string
	RelatedWithCustomer int64        // ID контакта, прикрепленного к объекту.
	Page                int64        // Номер страницы ответа. Начинается с 1. Игнорируется StockFilterAll.
	Date                [2]time.Time // {from: "2015-10-29", to: "2015-11-19"} выборка за определенный период
	// Активность объектов.
	//  Активные: "1" (по умолчанию)
	//  Удаленные: "0"
	//  Все: "ignore"
	Publish     string
	Limit       int64   // Кол-во объектов в ответе.
	SliceFields []int64 // Массив ID полей, значения которых будут в ответе. По умолчанию выводятся все.
}

// clone возвращает shallow-копию StockFilterParams.
func (p StockFilterParams) clone() *StockFilterParams {
	return new(p)
}

// cloneWithPage возвращает shallow-копию StockFilterParams с указанной страницей.
func (p StockFilterParams) cloneWithPage(page int64) *StockFilterParams {
	pageParams := p.clone()
	pageParams.Page = page
	return pageParams
}

// params возвращает параметры запроса в формате map[string]string (с эффективным выделением памяти).
func (p StockFilterParams) params() map[string]string {
	// Выделение памяти
	size := 8 // Поля с простыми типами
	size += len(p.ByIDs)
	size += len(p.Manager)
	size += len(p.Groups)
	size += len(p.SliceFields)
	size += len(p.Fields) * 2
	paramsMap := make(map[string]string, size)

	// type
	addToSingularParams(paramsMap, "type", p.Type)
	// category
	addToSingularParams(paramsMap, "category", p.Category)
	// byid | by_ids
	switch {
	case len(p.ByIDs) == 1:
		addToSingularParams(paramsMap, "byid", p.ByIDs[0])
	case len(p.ByIDs) >= 2:
		addSliceToSingularParams(paramsMap, "by_ids", p.ByIDs)
	}
	// search
	addToSingularParams(paramsMap, "search", p.Search)
	// manager
	addSliceToSingularParams(paramsMap, "manager", p.Manager)
	// groups
	addSliceToSingularParams(paramsMap, "groups", p.Groups)
	// stock_creator_id
	addToSingularParams(paramsMap, "stock_creator_id", p.StockCreatorID)
	// fields
	fieldsCount := 0
	for id, v := range p.Fields {
		if id == 0 || v == "" {
			continue
		}
		paramsMap[fmt.Sprintf("params[fields][%d][id]", fieldsCount)] = strconv.FormatInt(id, 10)
		paramsMap[fmt.Sprintf("params[fields][%d][value]", fieldsCount)] = v
		fieldsCount++
	}
	// related_with_customer
	addToSingularParams(paramsMap, "related_with_customer", p.RelatedWithCustomer)
	// page
	addToSingularParams(paramsMap, "page", p.Page)
	// date
	if !p.Date[0].IsZero() {
		paramsMap["params[date][from]"] = p.Date[0].Format(DatetimeLayout)
	}
	if !p.Date[1].IsZero() {
		paramsMap["params[date][to]"] = p.Date[1].Format(DatetimeLayout)
	}
	// publish
	addBoolToSingularParams(paramsMap, "publish", p.Publish)
	// limit
	switch v := p.Limit; {
	case v == 0, v >= 500:
		addToSingularParams(paramsMap, "limit", "500")
	default:
		addToSingularParams(paramsMap, "limit", v)
	}
	// slice_fields (SliceFields + Fields)
	if len(p.SliceFields) != 0 {
		sliceFields := make([]int64, 0, len(p.SliceFields)+len(p.Fields))
		sliceFields = append(sliceFields, p.SliceFields...)
		for id := range p.Fields {
			if id != 0 {
				sliceFields = append(sliceFields, id)
			}
		}
		addSliceToSingularParams(paramsMap, "slice_fields", sliceFields)
	}

	return paramsMap
}

// =====================================================================================================================
// Response
// =====================================================================================================================

type (
	StockFilterResponse struct {
		Status  string          `json:"status,omitempty"`
		Message string          `json:"message,omitempty"`
		Data    StockFilterData `json:"data,omitempty"`
	}
	StockFilterData struct {
		List []Stock `json:"list"`
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
	s.DateCreate = ParseDatetime(aux.DateCreate)
	s.LastModify = ParseDatetime(aux.LastModify)
	s.StockActivityDate = ParseDatetime(aux.StockActivityDate)

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

	return ParseDate(valueStr), true
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

	return ParseDatetime(valueStr), true
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

	return ParseInt(valueStr), true
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

	return ParseFloat(valueStr), true
}

// FieldDecimalOrZero возвращает значение поля (decimal) по id.
func (s *Stock) FieldDecimalOrZero(id int64) float64 {
	result, _ := s.FieldDecimal(id)
	return result
}

// FieldPrice возвращает значение поля (price) по id.
func (s *Stock) FieldPrice(id int64) (float64, bool) { return s.FieldDecimal(id) }

// FieldPriceOrZero возвращает значение поля (price) по id.
func (s *Stock) FieldPriceOrZero(id int64) float64 { return s.FieldDecimalOrZero(id) }

// FieldFile возвращает значение поля (file) по id.
func (s *Stock) FieldFile(id int64) ([]string, bool) { return s.FieldMultiselect(id) }

// FieldFileOrZero возвращает значение поля (files) по id.
func (s *Stock) FieldFileOrZero(id int64) []string {
	result, _ := s.FieldFile(id)
	return result
}

// FieldPoint возвращает значение поля (point) по id.
func (s *Stock) FieldPoint(id int64) (*Point, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return nil, false
	}

	var (
		x, _ = valueMap["x"]
		y, _ = valueMap["y"]
	)

	point, err := NewPointFromStrings(x, y)
	if err != nil {
		return nil, false
	}

	return point, true
}

// FieldPointOrZero возвращает значение поля (point) по id.
func (s *Stock) FieldPointOrZero(id int64) *Point {
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

	return ParseRange(valueMap, ParseInt), true
}

// FieldDecimalRange возвращает значение поля (decimal_range) по id.
func (s *Stock) FieldDecimalRange(id int64) ([2]float64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]float64{}, false
	}

	return ParseRange(valueMap, ParseFloat), true
}

// FieldDateRange возвращает значение поля (date_range) по id.
func (s *Stock) FieldDateRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return ParseRange(valueMap, ParseDate), true
}

// FieldDatetimeRange возвращает значение поля (datetime_range) по id.
func (s *Stock) FieldDatetimeRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return ParseRange(valueMap, ParseDatetime), true
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
		if vInt64 := ParseInt(vStr); vInt64 != 0 {
			valueInt64Slice = append(valueInt64Slice, vInt64)
		}
	}

	return valueInt64Slice, true
}
