package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// Ссылка на метод: https://www.intrumnet.com/api/#stock-search

type StockFilterParams struct {
	Type                uint32       // ID типа объекта (Обязательное поле, если не указаны ByID/ByIDs)
	ByID                uint64       // ID объекта
	ByIDs               []uint64     // Массив ID объектов (Все объекты из массива должны быть одного типа)
	Category            uint32       // ID категории объекта
	Nested              bool         // Включить вложенные категории
	Search              string       // Поисковая строка может содержать имя объекта или вхождения в поля с типами text,select,multiselect (полнотекстовый поиск)
	Manager             []uint64     // Массив ID ответственных
	Groups              []uint32     // Массив CRM групп
	StockCreatorID      uint64       // ID создателя
	IndexFields         bool         // Индексировать массив полей по ID свойства
	RelatedWithCustomer uint64       // ID контакта, связанного с объектом
	Order               string       // Направление сортировки (asc - по возрастанию, desc - по убыванию)
	OrderField          uint64       // ID поля, по которому нужно сделать сортировку (если в качестве значения указать stock_activity_date выборка будет сортироваться по дате активности; date_add - по дате создания, date_delete - по дате удаления, ID - по ID)
	Date                [2]time.Time // Выборка за определенный период
	DateField           string       // Если в качестве значения указать stock_activity_date, то выборка по параметру последней активности (в этом случае период выборки нужно передавать в параметре date)
	Page                uint16       // Номер страницы выборки (например, 2 страница с limit 500 на каждой, нумерация page начиная с 1)
	Publish             string       // "1" - активные | "0" - удаленные | "ignore" - вывод всех (по умолчанию "1")
	Limit               uint32       // Число записей в выборке, по умолчанию 20, макс. 500
	GroupID             uint32       // ID группы для группированных объектов
	Copy                uint64       // ID Родителя группы для группированных объектов
	ObjectGroups        uint32       // Число записей в выборке, по умолчанию 20, макс. 500
	CountTotal          bool         // Подсчет общего количества найденых записей
	OnlyPrimaryID       bool         // Вывести в ответе только ID объектов
	OnlyCountField      bool         // Вывести в ответе только количество
	SliceFields         []uint64     // Массив id дополнительных полей, которые будут в ответе (по умолчанию если не задано то выводятся все)
	SumField            uint64       // ID поля, которое нужно просуммировать. В ответе будет сумма значений поля результатов выборки (переменная: sum_field) и их число (count_field). Опция работает только для числовых полей (целое, число, цена)
	// Log // TODO

	// Массив условий поиска.
	//	Ключ - ID поля
	//	Значение - значение поля
	// Для полей с типом integer,decimal,price,time,date,datetime возможно указывать границы:
	//	Value: ">= {значение}" - больше или равно
	//	Value: "<= {значение}" - меньше или равно
	//	Value: "{значение_1} & {значение_2}" - между значением 1 и 2
	Fields map[uint64]string
}

// Ссылка на метод: https://www.intrumnet.com/api/#stock-search
func StockFilter(ctx context.Context, subdomain, apiKey string, inputParams *StockFilterParams) (*StockFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/filter", subdomain)

	// Параметры запроса

	params := make(map[string]string, 8)

	// type
	if inputParams.Type != 0 {
		params["params[type]"] = strconv.FormatUint(uint64(inputParams.Type), 10)
	}
	// byid
	if inputParams.ByID != 0 {
		params["params[byid]"] = strconv.FormatUint(inputParams.ByID, 10)
	}
	// by_ids
	addSliceToParams("by_ids", params, inputParams.ByIDs)
	// category
	if inputParams.Category != 0 {
		params["params[category]"] = strconv.FormatUint(uint64(inputParams.Category), 10)
	}
	// nested
	switch inputParams.Nested {
	case true:
		params["params[nested]"] = "1"
	default:
		params["params[nested]"] = "0"
	}
	// search
	if inputParams.Search != "" {
		params["params[search]"] = inputParams.Search
	}
	// manager
	addSliceToParams("manager", params, inputParams.Manager)
	// groups
	addSliceToParams("groups", params, inputParams.Groups)
	// stock_creator_id
	if inputParams.StockCreatorID != 0 {
		params["params[stock_creator_id]"] = strconv.FormatUint(inputParams.StockCreatorID, 10)
	}
	// fields
	fieldCount := 0
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[fields][%d][id]", fieldCount)] = fmt.Sprint(k)
		params[fmt.Sprintf("params[fields][%d][value]", fieldCount)] = v
		fieldCount++
	}
	// index_fields
	if inputParams.IndexFields {
		params["params[index_fields]"] = "1"
	}
	// related_with_customer
	if inputParams.RelatedWithCustomer != 0 {
		params["params[related_with_customer]"] = strconv.FormatUint(inputParams.RelatedWithCustomer, 10)
	}
	// order
	if inputParams.Order != "" {
		params["params[order]"] = inputParams.Order
	}
	// order_field
	if inputParams.OrderField != 0 {
		params["params[order_field]"] = strconv.FormatUint(inputParams.OrderField, 10)
	}
	// date
	if !inputParams.Date[0].IsZero() && !inputParams.Date[1].IsZero() {
		params["params[date][from]"] = inputParams.Date[0].Format(datetimeLayout)
		params["params[date][to]"] = inputParams.Date[1].Format(datetimeLayout)
	}
	// date_field
	if inputParams.DateField != "" {
		params["params[date_field]"] = inputParams.DateField
	}
	// page
	if inputParams.Page != 0 {
		params["params[page]"] = strconv.FormatUint(uint64(inputParams.Page), 10)
	}
	// publish
	switch inputParams.Publish {
	case "1", "true":
		params["params[publish]"] = "1"
	case "0", "false":
		params["params[publish]"] = "0"
	case "ignore":
		params["params[publish]"] = "ignore"
	}
	// limit
	switch l := inputParams.Limit; {
	case l == 0, l >= 500:
		params["params[limit]"] = "500"
	default:
		params["params[limit]"] = strconv.FormatUint(uint64(l), 10)
	}
	// group_id
	if inputParams.GroupID != 0 {
		params["params[group_id]"] = strconv.FormatUint(uint64(inputParams.GroupID), 10)
	}
	// copy
	if inputParams.Copy != 0 {
		params["params[copy]"] = strconv.FormatUint(inputParams.Copy, 10)
	}
	// object_groups
	switch {
	case inputParams.ObjectGroups > 500:
		params["params[object_groups]"] = "500"
	case inputParams.ObjectGroups != 0:
		params["params[object_groups]"] = strconv.FormatUint(uint64(inputParams.ObjectGroups), 10)
	}
	// count_total
	if inputParams.CountTotal {
		params["params[count_total]"] = "1"
	}
	// only_primary_id
	if inputParams.OnlyPrimaryID {
		params["params[only_primary_id]"] = "1"
	}
	// only_count_field
	if inputParams.OnlyCountField {
		params["params[only_count_field]"] = "1"
	}
	// slice_fields
	addSliceToParams("slice_fields", params, inputParams.SliceFields)
	// sum_field
	if inputParams.SumField != 0 {
		params["params[sum_field]"] = strconv.FormatUint(inputParams.SumField, 10)
	}
	// log // TODO

	// Получение ответа

	var resp StockFilterResponse
	if err := rawRequest(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
