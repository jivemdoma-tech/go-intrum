package intrum

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TODO: Реализовать оставшиеся поля SalesFilterParams:
//  Order
//  OrderField
//  Date
//  DateField
//  CountTotal
//  OnlyCountField

const SalesFilterMaxLimit int64 = 500

// SalesFilter - поиск сделок в CRM. Документация: https://www.intrumnet.com/api/#sales-filter
func SalesFilter(ctx context.Context, subdomain, apiKey string, p *SalesFilterParams) (*SalesFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/filter", subdomain)

	// Валидация
	if err := validateRequestArgs(methodURL, subdomain, apiKey); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, newErrEmptyParams(methodURL)
	}

	// Запрос
	resp := &SalesFilterResponse{}
	if err := request(ctx, apiKey, methodURL, p.params(), resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// =====================================================================================================================
// Request
// =====================================================================================================================

// SalesFilterParams - параметры запроса SalesFilter.
type SalesFilterParams struct {
	Type          []int64 // Массив ID типов сделок.
	Stage         []int64 // Массив ID стадий сделок.
	ByIDs         []int64 // Массив ID сделок.
	Search        string  // Поисковая строка.
	Manager       []int64 // Массив ID ответственных.
	Groups        []int64 // Массив ID CRM-групп.
	SaleCreatorID int64   // ID создателя сделки.
	Customer      int64   // ID контакта, прикрепленного к сделке.
	// Массив ID полей и значений.
	//  {ID ПОЛЯ}: "{ЗНАЧЕНИЕ}"
	// Для типов (integer, decimal, price, time, date, datetime) возможно указывать границы:
	//	Больше или равно: ">= {ЗНАЧЕНИЕ}"
	//	Меньше или равно: "<= {ЗНАЧЕНИЕ}"
	//	Между двумя значениям: "{ЗНАЧЕНИЕ} & {ЗНАЧЕНИЕ}"
	Fields      map[int64]string
	SliceFields []int64 // Массив ID полей, значения которых будут в ответе. По умолчанию выводятся все.
	// Активность объектов.
	//  Активные: "1" (по умолчанию)
	//  Удаленные: "0"
	//  Все: "ignore"
	Publish string
	Limit   int64 // Кол-во объектов в ответе.
	Page    int64 // Номер страницы ответа. Начинается с 1. Игнорируется StockFilterAll.
}

// params возвращает параметры запроса в формате map[string]string (с эффективным выделением памяти).
func (p SalesFilterParams) params() map[string]string {
	// Выделение памяти
	size := 6 // Поля с простыми типами
	size += len(p.Type)
	size += len(p.Stage)
	size += len(p.ByIDs)
	size += len(p.Manager)
	size += len(p.Groups)
	size += len(p.SliceFields)
	size += len(p.Fields) * 2
	paramsMap := make(map[string]string, size)

	// type
	addSliceToSingularParams(paramsMap, "type", p.Type)
	// stage
	addSliceToSingularParams(paramsMap, "stage", p.Stage)
	// byid + by_ids
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
	// sale_creator_id
	if v := p.SaleCreatorID; v > 0 {
		addToSingularParams(paramsMap, "sale_creator_id", v)
	}
	// customer
	if v := p.Customer; v > 0 {
		addToSingularParams(paramsMap, "customer", v)
	}
	// fields
	fieldsCount := 0
	for k, v := range p.Fields {
		if k == 0 || v == "" {
			continue
		}
		paramsMap[fmt.Sprintf("params[fields][%d][id]", fieldsCount)] = strconv.FormatInt(k, 10)
		paramsMap[fmt.Sprintf("params[fields][%d][value]", fieldsCount)] = v
		fieldsCount++
	}
	// slice_fields
	addSliceToSingularParams(paramsMap, "slice_fields", p.SliceFields)
	// publish
	addBoolToSingularParams(paramsMap, "publish", p.Publish)
	// limit
	switch v := p.Limit; {
	case v == 0, v >= 500:
		addToSingularParams(paramsMap, "limit", "500")
	default:
		addToSingularParams(paramsMap, "limit", v)
	}
	// page
	if v := p.Page; v >= 1 {
		addToSingularParams(paramsMap, "page", v)
	}

	return paramsMap
}

// =====================================================================================================================
// Response
// =====================================================================================================================

type (
	SalesFilterResponse struct {
		Status  string          `json:"status,omitempty"`
		Message string          `json:"message,omitempty"`
		Data    SalesFilterData `json:"data,omitempty"`
	}
	SalesFilterData struct {
		List []Sale `json:"list"`
	}
	Sale struct {
		ID                   int64                `json:"id,string"`              // ID сделки
		CustomersID          int64                `json:"customers_id,string"`    // ID контакта
		EmployeeID           int64                `json:"employee_id,string"`     // ID ответственного
		AdditionalEmployeeID []int64              `json:"additional_employee_id"` // Массив ID доп. ответственных
		DateCreate           time.Time            `json:"date_create"`            // Дата создания
		SaleTypeID           int64                `json:"sale_type_id,string"`    // ID типа активности
		SaleStageID          int64                `json:"sale_stage_id,string"`   // ID стадии
		SaleName             string               `json:"sale_name"`              // Название сделки
		SaleActivityType     string               `json:"sale_activity_type"`     // Тип последней активности
		SaleActivityDate     time.Time            `json:"sale_activity_date"`     // Дата последней активности сделк
		Fields               map[string]SaleField `json:"fields"`                 // Данные полей
		Publish              bool                 `json:"publish"`                // Опубликован/Удален
	}
	SaleField struct {
		DataType string `json:"datatype"`
		Value    any    `json:"value"`
	}
)

func (r *SalesFilterResponse) GetErrorMessage() string {
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

func (s *Sale) UnmarshalJSON(data []byte) error {
	// alias - обертка над оригинальной структурой для предотвращения рекурсии
	type alias Sale

	// Вспомогательная структура (Приведение типа к alias)
	var aux = &struct {
		*alias

		// Дата + время
		DateCreate       string `json:"date_create"`
		SaleActivityDate string `json:"sale_activity_date"`
		// Bool
		Publish string `json:"publish"`
		// Массивы
		AdditionalEmployeeID []string `json:"additional_employee_id"`
	}{alias: (*alias)(s)}

	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Обработка полей и добавление в оригинальную структуру

	// (time.Time) Поля типа datetime
	s.DateCreate = ParseDatetime(aux.DateCreate)
	s.SaleActivityDate = ParseDatetime(aux.SaleActivityDate)

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

	return nil
}

// field возвращает поле по id.
func (s *Sale) field(id int64) (*SaleField, bool) {
	if s == nil {
		return nil, false
	}
	// Проверка: поле существует
	field, exists := s.Fields[strconv.FormatInt(id, 10)]
	if !exists {
		return nil, false
	}

	return &field, true
}

// fieldMap возвращает значение поля (map[string]string) по id.
func (s *Sale) fieldMap(id int64) (map[string]string, bool) {
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
func (s *Sale) FieldText(id int64) (string, bool) {
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
func (s *Sale) FieldTextOrZero(id int64) string {
	result, _ := s.FieldText(id)
	return result
}

// FieldRadio возвращает значение поля (radio) по id.
func (s *Sale) FieldRadio(id int64) (bool, bool) {
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
func (s *Sale) FieldRadioOrZero(id int64) bool {
	result, _ := s.FieldRadio(id)
	return result
}

// FieldSelect возвращает значение поля (select) по id.
func (s *Sale) FieldSelect(id int64) (string, bool) { return s.FieldText(id) }

// FieldSelectOrZero возвращает значение поля (select) по id.
func (s *Sale) FieldSelectOrZero(id int64) string { return s.FieldTextOrZero(id) }

// FieldMultiselect возвращает значение поля (multiselect) по id.
func (s *Sale) FieldMultiselect(id int64) ([]string, bool) {
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
func (s *Sale) FieldMultiselectOrZero(id int64) []string {
	result, _ := s.FieldMultiselect(id)
	return result
}

// FieldDate возвращает значение поля (date) по id.
func (s *Sale) FieldDate(id int64) (time.Time, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return time.Time{}, false
	}

	return ParseDate(valueStr), true
}

// FieldDateOrZero возвращает значение поля (date) по id.
func (s *Sale) FieldDateOrZero(id int64) time.Time {
	result, _ := s.FieldDate(id)
	return result
}

// FieldDatetime возвращает значение поля (datetime) по id.
func (s *Sale) FieldDatetime(id int64) (time.Time, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return time.Time{}, false
	}

	return ParseDatetime(valueStr), true
}

// FieldDatetimeOrZero возвращает значение поля (datetime) по id.
func (s *Sale) FieldDatetimeOrZero(id int64) time.Time {
	result, _ := s.FieldDatetime(id)
	return result
}

// FieldInteger возвращает значение поля (integer) по id.
func (s *Sale) FieldInteger(id int64) (int64, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return 0, false
	}

	return ParseInt(valueStr), true
}

// FieldIntegerOrZero возвращает значение поля (integer) по id.
func (s *Sale) FieldIntegerOrZero(id int64) int64 {
	result, _ := s.FieldInteger(id)
	return result
}

// FieldDecimal возвращает значение поля (decimal) по id.
func (s *Sale) FieldDecimal(id int64) (float64, bool) {
	// Проверка: поле существует
	valueStr, exists := s.FieldText(id)
	if !exists {
		return 0, false
	}

	return ParseFloat(valueStr), true
}

// FieldDecimalOrZero возвращает значение поля (decimal) по id.
func (s *Sale) FieldDecimalOrZero(id int64) float64 {
	result, _ := s.FieldDecimal(id)
	return result
}

// FieldPrice возвращает значение поля (price) по id.
func (s *Sale) FieldPrice(id int64) (float64, bool) { return s.FieldDecimal(id) }

// FieldPriceOrZero возвращает значение поля (price) по id.
func (s *Sale) FieldPriceOrZero(id int64) float64 { return s.FieldDecimalOrZero(id) }

// FieldFile возвращает значение поля (file) по id.
func (s *Sale) FieldFile(id int64) ([]string, bool) {
	v, exists := s.FieldText(id)
	if !exists {
		return nil, false
	}

	v = strings.TrimPrefix(v, "[")
	v = strings.TrimSuffix(v, "]")

	return strings.Fields(v), true
}

// FieldFileOrZero возвращает значение поля (files) по id.
func (s *Sale) FieldFileOrZero(id int64) []string {
	result, _ := s.FieldFile(id)
	return result
}

// FieldPoint возвращает значение поля (point) по id.
func (s *Sale) FieldPoint(id int64) (*Point, bool) {
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
func (s *Sale) FieldPointOrZero(id int64) *Point {
	result, _ := s.FieldPoint(id)
	return result
}

// FieldIntegerRange возвращает значение поля (integer_range) по id.
func (s *Sale) FieldIntegerRange(id int64) ([2]int64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]int64{}, false
	}

	return ParseRange(valueMap, ParseInt), true
}

// FieldDecimalRange возвращает значение поля (decimal_range) по id.
func (s *Sale) FieldDecimalRange(id int64) ([2]float64, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]float64{}, false
	}

	return ParseRange(valueMap, ParseFloat), true
}

// FieldDateRange возвращает значение поля (date_range) по id.
func (s *Sale) FieldDateRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return ParseRange(valueMap, ParseDate), true
}

// FieldDatetimeRange возвращает значение поля (datetime_range) по id.
func (s *Sale) FieldDatetimeRange(id int64) ([2]time.Time, bool) {
	// Проверка: поле существует
	valueMap, exists := s.fieldMap(id)
	if !exists {
		return [2]time.Time{}, false
	}

	return ParseRange(valueMap, ParseDatetime), true
}

// FieldAttach возвращает значение поля (attach) по id.
func (s *Sale) FieldAttach(id int64) ([]int64, bool) {
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
