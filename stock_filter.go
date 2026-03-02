package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

const (
	StockFilterMaxLimit int64 = 500
)

// StockFilterParams - параметры запроса.
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
	RelatedWithCustomer int64 // ID контакта, прикрепленного к объекту.
	Page                int64 // Номер страницы ответа. Начинается с 1. Игнорируется StockFilterAll.
	// Активность объектов.
	//  Активные: "1" (по умолчанию)
	//  Удаленные: "0"
	//  Все: "ignore"
	Publish     string
	Limit       int64   // Кол-во объектов в ответе.
	SliceFields []int64 // Массив ID полей, значения которых будут в ответе. По умолчанию выводятся все.

	// TODO: Оставшиеся поля. При реализации полей адаптируйте выделение памяти для resultParams в методе params.
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
}

// copy возвращает shallow-копию структуры.
func (p StockFilterParams) copy() *StockFilterParams {
	copyParams := p
	return &copyParams
}

// copyWithPage возвращает shallow-копию структуры с новой страницей.
func (p StockFilterParams) copyWithPage(page int64) *StockFilterParams {
	pageParams := p.copy()
	pageParams.Page = page
	return pageParams
}

// params возвращает параметры запроса в формате map[string]string (с эффективным выделением памяти).
func (p StockFilterParams) params() map[string]string {
	// Выделение памяти
	resultParams := make(map[string]string,
		// Единичные поля
		8+
			// Слайсы
			len(p.ByIDs)+
			len(p.Manager)+
			len(p.Groups)+
			len(p.SliceFields)+
			// Мапы
			len(p.Fields)*2,
	)

	// type
	addToSingularParams(resultParams, "type", p.Type)
	// category
	addToSingularParams(resultParams, "category", p.Category)
	// byid | by_ids
	switch {
	case len(p.ByIDs) == 1:
		addToSingularParams(resultParams, "byid", p.ByIDs[0])
	case len(p.ByIDs) >= 2:
		addSliceToSingularParams(resultParams, "by_ids", p.ByIDs)
	}
	// search
	addToSingularParams(resultParams, "search", p.Search)
	// manager
	addSliceToSingularParams(resultParams, "manager", p.Manager)
	// groups
	addSliceToSingularParams(resultParams, "groups", p.Groups)
	// stock_creator_id
	addToSingularParams(resultParams, "stock_creator_id", p.StockCreatorID)
	// fields
	fieldsCount := 0
	for k, v := range p.Fields {
		if k == 0 || v == "" {
			continue
		}
		resultParams[fmt.Sprintf("params[fields][%d][id]", fieldsCount)] = strconv.FormatInt(k, 10)
		resultParams[fmt.Sprintf("params[fields][%d][value]", fieldsCount)] = v
		fieldsCount++
	}
	// related_with_customer
	addToSingularParams(resultParams, "related_with_customer", p.RelatedWithCustomer)
	// page
	addToSingularParams(resultParams, "page", p.Page)
	// publish
	addBoolToSingularParams(resultParams, "publish", p.Publish)
	// limit
	switch v := p.Limit; {
	case v == 0, v >= 500:
		addToSingularParams(resultParams, "limit", "500")
	default:
		addToSingularParams(resultParams, "limit", v)
	}
	// slice_fields (SliceFields + Fields)
	sliceFields := make([]int64, 0, len(p.SliceFields)+len(p.Fields))
	sliceFields = append(sliceFields, p.SliceFields...)
	for id := range p.Fields {
		if id != 0 {
			sliceFields = append(sliceFields, id)
		}
	}
	addSliceToSingularParams(resultParams, "slice_fields", sliceFields)

	return resultParams
}

// StockFilter - поиск объектов в CRM. Документация: https://www.intrumnet.com/api/#stock-search
func StockFilter(ctx context.Context, subdomain, apiKey string, p *StockFilterParams) (*StockFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/filter", subdomain)

	// Валидация
	if p == nil {
		return nil, newErrNilParams(methodURL)
	}
	// Обязательные поля
	if p.Type <= 0 && len(p.ByIDs) == 0 {
		return nil, newErrEmptyRequiredFields(methodURL)
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
	resultStock := make([]Stock, 0, 500)
	for page := int64(1); ; page++ {
		// Shallow-копирование структуры для итерации
		pageParams := p.copyWithPage(page)
		// Установка максимального кол-ва элементов в ответе
		if pageParams.Limit != StockFilterMaxLimit {
			pageParams.Limit = StockFilterMaxLimit
		}
		// Запрос
		resp, err := StockFilter(ctx, subdomain, apiKey, pageParams)
		if err != nil {
			return nil, err
		}

		if len(resp.Data.List) == 0 {
			break
		}

		resultStock = append(resultStock, resp.Data.List...)

		if len(resp.Data.List) < int(pageParams.Limit) {
			break
		}
	}
	if len(resultStock) == 0 {
		return nil, nil
	}

	return resultStock, nil
}
