package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// StockInsertParams - параметры запроса.
//
// Обязательные поля:
//   - Type
//
// Основные параметры запроса:
//   - Type: ID типа объекта.
//   - Name: Название объекта.
//   - Manager: ID главного ответственного ответственных.
//   - AdditionalManagers: Массив ID доп. ответственных.
//   - RelatedWithCustomer: ID контакта, прикрепленного к объекту
//   - Fields: массив ID полей и значений.
//     Для типа (multiselect) возможно указывать несколько вариантов: "{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ}".
type StockInsertParams struct {
	Type                int64   // ID типа объекта.
	Name                string  // Название объекта.
	Manager             int64   // ID главного ответственного ответственных.
	AdditionalManagers  []int64 // Массив ID доп. ответственных.
	RelatedWithCustomer int64   // ID контакта, прикрепленного к объекту
	// Fields: массив ID полей и значений.
	//
	// Для типа (multiselect) возможно указывать несколько вариантов:
	//  "{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ}".
	Fields      map[int64]string
	FieldsPoint map[int64]StockInsertPoint // Аналогично Fields для типа "point".
	FieldsFile  map[int64][]string         // Аналогично Fields для типа "file".

	// TODO: Оставшиеся поля. При реализации полей адаптируйте выделение памяти для paramsMap в методе params.
	//  GroupID
	//  Copy
}

type StockInsertPoint struct {
	Lat float64 // Широта
	Lon float64 // Долгота
}

// params возвращает параметры запроса в формате map[string]string (с эффективным выделением памяти).
func (p StockInsertParams) params() map[string]string {
	// Выделение памяти
	paramsMap := make(map[string]string,
		// Единичные поля
		4+
			// Слайсы
			len(p.AdditionalManagers)+
			// Мапы
			len(p.Fields)*2+
			len(p.FieldsPoint)*2+
			len(p.FieldsFile)*2,
	)

	// parent
	paramsMap["params[0][parent]"] = strconv.FormatInt(p.Type, 10)
	// name
	if p.Name != "" {
		paramsMap["params[0][name]"] = p.Name
	}
	// author
	if p.Manager != 0 {
		paramsMap["params[0][author]"] = strconv.FormatInt(p.Manager, 10)
	}
	// additional_author
	for i, v := range p.AdditionalManagers {
		paramsMap[fmt.Sprintf("params[0][additional_author][%d]", i)] = strconv.FormatInt(v, 10)
	}
	// related_with_customer
	if p.RelatedWithCustomer != 0 {
		paramsMap["params[0][related_with_customer]"] = strconv.FormatInt(p.RelatedWithCustomer, 10)
	}

	countFields := 0
	// fields
	for k, v := range p.Fields {
		paramsMap[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
		paramsMap[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}
	for k, v := range p.FieldsPoint {
		paramsMap[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
		paramsMap[fmt.Sprintf("params[0][fields][%d][value][lat]", countFields)] = strconv.FormatFloat(v.Lat, 'f', 10, 64)
		paramsMap[fmt.Sprintf("params[0][fields][%d][value][lon]", countFields)] = strconv.FormatFloat(v.Lon, 'f', 10, 64)
		countFields++
	}
	for k, fileNames := range p.FieldsFile {
		for _, fileName := range fileNames {
			paramsMap[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
			paramsMap[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = fileName
			countFields++
		}
	}

	return paramsMap
}

// StockInsert - поиск объектов в CRM. Документация: https://www.intrumnet.com/api/#stock-insert
//
// Ограничение: 1 запрос == 1 объект
func StockInsert(ctx context.Context, subdomain, apiKey string, p *StockInsertParams) (*StockInsertResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/insert", subdomain)

	// Валидация
	if p == nil {
		return nil, newErrNilParams(methodURL)
	}
	// Обязательные поля
	if p.Type <= 0 {
		return nil, newErrEmptyRequiredFields(methodURL)
	}

	// Запрос
	resp := &StockInsertResponse{}
	if err := request(ctx, apiKey, methodURL, p.params(), resp); err != nil {
		return nil, err
	}

	return resp, nil
}
