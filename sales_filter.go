package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter
type SalesFilterParams struct {
	Manager     []uint64 // Массив ID ответственных
	Type        []uint64 // Массив ID типов сделок
	Stage       []uint64 // Массив ID стадий сделок
	ByIDs       []uint64 // Получение сделок по массиву ID
	SliceFields []uint64 // Массив ID дополнительных полей, которые будут в ответе (по умолчанию выводятся все)
	Limit       uint64   // Число записей в выборке (Макс. 500)
	Search      string   // Поисковая строка
	// Массив условий поиска.
	//	Ключ - ID поля
	//	Значение - значение поля
	// Для полей с типом integer,decimal,price,time,date,datetime возможно указывать границы:
	//	Value: ">= {значение}" - больше или равно
	//	Value: "<= {значение}" - меньше или равно
	//	Value: "{значение_1} & {значение_2}" - между значением 1 и 2
	Fields map[uint64]string

	// TODO: Добавить больше параметров запроса
	// Customer       uint64       // ID контакта
	// Groups         []uint64     // Массив CRM групп
	// SaleCreatorID  uint64       // ID создателя
	// Page           uint16       // Номер страницы выборки
	// Publish        bool         // 1 - активные, 0 - удаленные, по умолчанию 1
	// ByID           uint64       // Получение сделки по ее id
	// Order          string       // Направление сортировки asc - по возрастанию, desc - по убыванию
	// OrderField     string       // Если в качестве значения указать sale_activity_date выборка будет сортироваться по дате активности, create_date - по дате создания, delete_date - по дате удаления, id - по id
	// Date           [2]time.Time // {from: "2015-10-29", to: "2015-11-19"} выборка за определенный период
	// DateField      string       // Если в качестве значения указать sale_activity_date выборка по параметру активности
	// CountTotal     bool         // Подсчет общего количества найденых записей, 1 - считать, 0 - нет (по умолчанию 0)
	// OnlyCountField bool         // 1 - вывести в ответе только количество, 0 - стандартный вывод (по умолчанию 0)
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter
func SalesFilter(ctx context.Context, subdomain, apiKey string, inputParams *SalesFilterParams) (*SalesFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/filter", subdomain)

	// Параметры запроса

	params := make(map[string]string, len(inputParams.Manager)+
		len(inputParams.Type)+
		len(inputParams.Stage)+
		len(inputParams.ByIDs)+
		len(inputParams.SliceFields)+
		1+ // limit
		len(inputParams.Fields)*2)

	// search
	if inputParams.Search != "" {
		params["params[search]"] = inputParams.Search
	}
	// manager
	addSliceToParams("manager", params, inputParams.Manager)
	// type
	addSliceToParams("type", params, inputParams.Type)
	// stage
	addSliceToParams("stage", params, inputParams.Stage)
	// by_ids
	addSliceToParams("by_ids", params, inputParams.ByIDs)
	// slice_fields
	addSliceToParams("slice_fields", params, inputParams.SliceFields)
	// limit
	switch l := inputParams.Limit; {
	case l == 0, l >= 500:
		params["params[limit]"] = "500"
	default:
		params["params[limit]"] = strconv.FormatUint(uint64(l), 10)
	}
	// fields
	var count int
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[fields][%d][id]", count)] = fmt.Sprint(k)
		params[fmt.Sprintf("params[fields][%d][value]", count)] = v
		count++
		// TODO: Добавить внешнюю функцию обработки value под формат php
	}

	// Получение ответа

	var resp SalesFilterResponse
	if err := request(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
