package intrumgo

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter
type SalesFilterParams struct {
	Manager     []uint64 // Массив ID ответственных
	Type        []uint16 // Массив ID типов сделок
	Stage       []uint16 // Массив ID стадий сделок
	ByIDs       []uint64 // Получение сделок по массиву ID
	SliceFields []uint64 // Массив ID дополнительных полей, которые будут в ответе (по умолчанию выводятся все)
	Limit       uint16   // Число записей в выборке (Макс. 500)

	// TODO: Добавить больше параметров запроса
	// Search         string       // Поисковая строка
	// Customer       uint32       // ID контакта
	// Groups         []uint16     // Массив CRM групп
	// SaleCreatorID  uint32       // ID создателя
	// Page           uint32       // Номер страницы выборки
	// Publish        bool         // 1 - активные, 0 - удаленные, по умолчанию 1
	// ByID           uint32       // Получение сделки по ее id
	// Order          string       // Направление сортировки asc - по возрастанию, desc - по убыванию
	// OrderField     string       // Если в качестве значения указать sale_activity_date выборка будет сортироваться по дате активности, create_date - по дате создания, delete_date - по дате удаления, id - по id
	// Date           [2]time.Time // {from: "2015-10-29", to: "2015-11-19"} выборка за определенный период
	// DateField      string       // Если в качестве значения указать sale_activity_date выборка по параметру активности
	// CountTotal     bool         // Подсчет общего количества найденых записей, 1 - считать, 0 - нет (по умолчанию 0)
	// OnlyCountField bool         // 1 - вывести в ответе только количество, 0 - стандартный вывод (по умолчанию 0)

	/*
			Массив условий поиска, где ключ - ID поля, значение - значение поля
			Для полей с типом integer,decimal,price,time,date,datetime возможно указывать границы:
		    Value: ">= значение" - больше или равно
		    Value: "<= значение" - меньше или равно
		    Value: "значение_1 & значение_2" - между значением 1 и 2
	*/
	Fields map[uint64]any
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter
func SalesFilter(ctx context.Context, subdomain, apiKey string, timeoutSec int, inputParams *SalesFilterParams) (*SalesFilterResponse, error) {
	var (
		primaryURL string = fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/filter", subdomain)
		backupURL  string = fmt.Sprintf("http://%s.intrumnet.com:80/sharedapi/sales/filter", subdomain)
	)

	// Параметры запроса

	params := make(map[string]string, getParamsSize(inputParams))

	// manager
	addSliceToParams(params, "manager", inputParams.Manager)
	// type
	addSliceToParams(params, "type", inputParams.Type)
	// stage
	addSliceToParams(params, "stage", inputParams.Stage)
	// by_ids
	addSliceToParams(params, "by_ids", inputParams.ByIDs)
	// slice_fields
	addSliceToParams(params, "slice_fields", inputParams.SliceFields)
	// fields
	if len(inputParams.Fields) > 0 {
		var count int
		for k, v := range inputParams.Fields {
			params[fmt.Sprintf("params[fields][%d][id]", count)] = fmt.Sprint(k)
			params[fmt.Sprintf("params[fields][%d][value]", count)] = fmt.Sprint(v)
			// TODO: Добавить внешнюю функцию обработки value под формат php
		}
	}
	// limit (макс. 500)
	switch {
	case inputParams.Limit > 500:
		params["params[limit]"] = "500"
	case inputParams.Limit != 0:
		params["params[limit]"] = strconv.FormatUint(uint64(inputParams.Limit), 10)
	}

	// Получение ответа

	var resp SalesFilterResponse

	if err := rawRequest(ctx, primaryURL, apiKey, timeoutSec, params, &resp); err != nil {
		if err := rawRequest(ctx, backupURL, apiKey, timeoutSec, params, &resp); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}
