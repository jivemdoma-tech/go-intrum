package intrum

import (
	"context"
	"fmt"
	"strconv"
)

// TODO: Реализовать оставшиеся поля StockInsertParams.
//  GroupID
//  Copy
// TODO: Реализовать в StockInsertParams.Fields изменение полей типов: attach

// StockInsert - добавление объекта в CRM. Документация: https://www.intrumnet.com/api/#stock-insert
func StockInsert(ctx context.Context, subdomain, apiKey string, p *StockInsertParams) (*StockInsertResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/insert", subdomain)

	// Валидация
	if err := validateRequestArgs(methodURL, subdomain, apiKey); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, newErrEmptyParams(methodURL)
	}

	// Обязательные поля
	if p.Category <= 0 {
		return nil, newErrEmptyRequiredParams(methodURL)
	}

	// Запрос
	resp := &StockInsertResponse{}
	if err := request(ctx, apiKey, methodURL, p.params(), resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// =====================================================================================================================
// Request
// =====================================================================================================================

// StockInsertParams - параметры запроса StockInsert.
//
// Обязательные поля:
//   - Category
//
// Основные параметры запроса:
//   - Category:            Id категории объекта.
//   - Name:                Название объекта.
//   - Manager:             Id главного ответственного.
//   - AdditionalManagers:  Массив id доп. ответственных.
//   - RelatedWithCustomer: Id контакта, прикрепленного к объекту.
//   - Fields:              Массив ID полей и значений.
//     Для типа (multiselect) возможно указывать несколько вариантов: "{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ}".
type StockInsertParams struct {
	Category            int64   // Id категории объекта.
	Name                string  // Название объекта.
	Manager             int64   // Id главного ответственного.
	AdditionalManagers  []int64 // Массив id доп. ответственных.
	RelatedWithCustomer int64   // Id контакта, прикрепленного к объекту.
	// Fields: массив id полей и значений.
	//
	// Для типа (multiselect) возможно указывать несколько вариантов:
	//  "{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ}".
	Fields      map[int64]string
	FieldsPoint map[int64]*Point   // Аналогично Fields для типа "point".
	FieldsFile  map[int64][]string // Аналогично Fields для типа "file".
}

// params возвращает параметры запроса в формате map[string]string (с эффективным выделением памяти).
func (p StockInsertParams) params() map[string]string {
	// Выделение памяти
	size := 5 // Поля с простыми типами
	size += len(p.AdditionalManagers)
	size += len(p.Fields) * 2
	size += len(p.FieldsPoint) * 3
	for _, files := range p.FieldsFile {
		size += len(files) * 2
	}
	paramsMap := make(map[string]string, size)

	// parent
	if v := p.Category; v > 0 {
		paramsMap["params[0][parent]"] = strconv.FormatInt(v, 10)
	}
	// name
	if v := p.Name; v != "" {
		paramsMap["params[0][name]"] = v
	}
	// author
	if v := p.Manager; v > 0 {
		paramsMap["params[0][author]"] = strconv.FormatInt(v, 10)
	}
	// additional_author
	for i, id := range p.AdditionalManagers {
		if id > 0 {
			k, v := fmt.Sprintf("params[0][additional_author][%d]", i), strconv.FormatInt(id, 10)
			paramsMap[k] = v
		}
	}
	// related_with_customer
	if v := p.RelatedWithCustomer; v > 0 {
		paramsMap["params[0][related_with_customer]"] = strconv.FormatInt(p.RelatedWithCustomer, 10)
	}

	// fields
	fieldsCount := 0
	for id, v := range p.Fields {
		if id <= 0 || v == "" {
			continue
		}
		// ID
		paramsMap[fmt.Sprintf("params[0][fields][%d][id]", fieldsCount)] = strconv.FormatInt(id, 10)
		// Value
		paramsMap[fmt.Sprintf("params[0][fields][%d][value]", fieldsCount)] = v

		fieldsCount++
	}
	// fields (point)
	for id, point := range p.FieldsPoint {
		if id <= 0 || point == nil {
			continue
		}
		// Получение и валидация координат
		latStr, lonStr := point.StringLat(), point.StringLon()
		if latStr == "" || lonStr == "" {
			continue
		}
		// ID
		paramsMap[fmt.Sprintf("params[0][fields][%d][id]", fieldsCount)] = strconv.FormatInt(id, 10)
		// Value
		paramsMap[fmt.Sprintf("params[0][fields][%d][value][lat]", fieldsCount)] = latStr
		paramsMap[fmt.Sprintf("params[0][fields][%d][value][lon]", fieldsCount)] = lonStr

		fieldsCount++
	}
	// fields (file)
	for id, files := range p.FieldsFile {
		if id <= 0 || len(files) == 0 {
			continue
		}
		// Обработка слайсов имен файлов
		for _, f := range files {
			// ID
			paramsMap[fmt.Sprintf("params[0][fields][%d][id]", fieldsCount)] = strconv.FormatInt(id, 10)
			// Value
			paramsMap[fmt.Sprintf("params[0][fields][%d][value]", fieldsCount)] = f

			fieldsCount++
		}
	}

	return paramsMap
}

// =====================================================================================================================
// Response
// =====================================================================================================================

type StockInsertResponse struct {
	Status  string  `json:"status,omitempty"`
	Message string  `json:"message,omitempty"`
	Data    []int64 `json:"data,omitempty"`
}

func (r *StockInsertResponse) GetErrorMessage() string {
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
