package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter
type WorkerFilterParams struct {
	// Group       uint16   // ID CRM группы // TODO
	ID          []uint64 // Массив id сотрудников
	DivisionID  []uint16 // Массив id отделов
	SubofficeID []uint16 // Массив id филиалов
	Surname     string   // Фамилия
	Name        string   // Имя
	Email       string   // Email
	Phone       string   // Телефон
	Boss        string   // 1 - Начадьник отдела, 0 - Не начальник отдела, не указано - вывод всех
	SliceFields []uint64 // массив id дополнительных полей, которые будут в ответе (по умолчанию, если не задано, то выводятся все)
	// По умолчанию - 1
	// 	"1" - активные
	// 	"0" - удаленные
	// 	"ignore" - все
	Publish string
	// Массив условий поиска.
	//	Ключ - ID поля
	//	Значение - значение поля
	// Для полей с типом integer,decimal,price,time,date,datetime возможно указывать границы:
	//	Value: ">= {значение}" - больше или равно
	//	Value: "<= {значение}" - меньше или равно
	//	Value: "{значение_1} & {значение_2}" - между значением 1 и 2
	Fields map[uint64]string
}

// Ссылка на метод: https://www.intrumnet.com/api/#worker-filter
func WorkerFilter(ctx context.Context, subdomain, apiKey string, paramsInput *WorkerFilterParams) (*WorkerFilterResponse, error) {
	u := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/worker/filter", subdomain)

	// Параметры запроса

	paramsResult := make(map[string]string, 12+len(paramsInput.Fields)+len(paramsInput.SliceFields))
	// id
	addSliceToParams("id", paramsResult, paramsInput.ID)
	// division_id
	addSliceToParams("division_id", paramsResult, paramsInput.DivisionID)
	// suboffice_id
	addSliceToParams("suboffice_id", paramsResult, paramsInput.SubofficeID)
	// surname
	if paramsInput.Surname != "" {
		paramsResult["params[surname]"] = paramsInput.Surname
	}
	// name
	if paramsInput.Name != "" {
		paramsResult["params[name]"] = paramsInput.Name
	}
	// email
	if paramsInput.Email != "" {
		paramsResult["params[email]"] = paramsInput.Email
	}
	// phone
	if paramsInput.Phone != "" {
		paramsResult["params[phone]"] = paramsInput.Phone
	}
	// publish
	if paramsInput.Publish != "" {
		paramsResult["params[publish]"] = paramsInput.Publish
	}
	// boss
	switch paramsInput.Boss {
	case "1", "0":
		paramsResult["params[boss]"] = paramsInput.Boss
	}
	// slice_fields
	addSliceToParams("slice_fields", paramsResult, paramsInput.SliceFields)
	// fields
	var count int
	for k, v := range paramsInput.Fields {
		paramsResult[fmt.Sprintf("params[fields][%d][id]", count)] = fmt.Sprint(k)
		paramsResult[fmt.Sprintf("params[fields][%d][value]", count)] = v
		count++
	}

	// Обработка ответа

	resp := new(WorkerFilterResponse)
	if err := rawRequest(ctx, apiKey, u, paramsResult, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
