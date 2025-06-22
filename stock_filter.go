package gointrum

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type StockFilterParams struct {
	// ID типа объекта
	// 	! ОБЯЗАТЕЛЬНО (Если не указан "ByIDs")
	Type uint64

	// Массив условий поиска.
	//	Key: ID поля
	//	Value: Значение поля
	// Для полей с типом integer, decimal, price, time, date, datetime возможно указывать границы:
	//	Value: ">= {значение}" - больше или равно
	//	Value: "<= {значение}" - меньше или равно
	//	Value: "{значение_1} & {значение_2}" - между значением 1 и 2
	Fields map[uint64]string

	ByIDs               []uint64 // Массив ID объектов (Все объекты из массива должны быть одного типа)
	Category            uint64   // ID категории объекта
	Nested              string   // (bool) Включить вложенные категории
	Search              string   // Поисковая строка. Может содержать имя объекта или вхождения в поля с типами text, select, multiselect (Полнотекстовый поиск)
	Manager             []uint64 // Массив ID ответственных
	Groups              []uint64 // Массив CRM групп
	StockCreatorID      uint64   // ID создателя
	IndexFields         string   // (bool) Индексировать массив полей по ID свойства
	RelatedWithCustomer uint64   // ID контакта, связанного с объектом
	Order               string   // Направление сортировки (asc - по возрастанию, desc - по убыванию)
	// ID поля, по которому нужно сделать сортировку. Если в качестве значения указать:
	// 	"stock_activity_date" - сортировка по дате активности
	// 	"date_add" - сортировка по дате создания
	// 	"date_delete" - сортировка по дате удаления
	// 	"ID" - по ID
	OrderField    string
	Date          [2]time.Time // Выборка за определенный период
	DateField     string       // Если в качестве значения указать stock_activity_date, то выборка по параметру последней активности (в этом случае период выборки нужно передавать в параметре date)
	Page          uint16       // Номер страницы выборки (например, 2 страница с limit 500 на каждой, нумерация page начиная с 1)
	Publish       string       // (bool) "1" - активные | "0" - удаленные | "ignore" - вывод всех (по умолчанию "1")
	Limit         uint64       // Число записей в выборке (По умолчанию 500)
	OnlyPrimaryID string       // (bool) Вывести в ответе только ID объектов
	SliceFields   []uint64     // Массив id дополнительных полей, которые будут в ответе (по умолчанию если не задано то выводятся все)

	// TODO

	// CountTotal     string // (bool) Подсчет общего количества найденых записей
	// OnlyCountField string // (bool) Вывести в ответе только количество
	// Log            string // Фильтр по истории изменений

	// SumField       uint64 // ID поля, которое нужно просуммировать. В ответе будет сумма значений поля результатов выборки (переменная: sum_field) и их число (count_field). Опция работает только для числовых полей (целое, число, цена)

	// GroupID        uint64 // ID группы для группированных объектов
	// Copy           uint64 // ID Родителя группы для группированных объектов
	// ObjectGroups   uint64 // Число записей в выборке, по умолчанию 20, макс. 500
}

// Ссылка на метод: https://www.intrumnet.com/api/#stock-search
func StockFilter(ctx context.Context, subdomain, apiKey string, inputParams *StockFilterParams) (*StockFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/filter", subdomain)

	// Обязательность ввода параметров
	if inputParams.Type == 0 && len(inputParams.ByIDs) == 0 {
		u, _ := url.ParseRequestURI(methodURL)
		return nil, fmt.Errorf("failed to create request for method %s: %s", u.Path, statusBadParams)
	}

	// Параметры запроса
	p := make(map[string]string, 8+
		len(inputParams.ByIDs)+
		len(inputParams.Manager)+
		len(inputParams.Groups)+
		len(inputParams.SliceFields)+
		len(inputParams.Fields)*2)

	// type
	addToParams(p, "type", inputParams.Type)
	// byid + by_ids
	switch {
	case len(inputParams.ByIDs) == 1:
		addToParams(p, "byid", inputParams.ByIDs[0])
	case len(inputParams.ByIDs) >= 2:
		addSliceToParams(p, "by_ids", inputParams.ByIDs)
	}
	// category
	addToParams(p, "category", inputParams.Category)
	// nested
	switch v := inputParams.Nested; v {
	case "true", "false", "1", "0":
		addToParams(p, "nested", v)
	}
	// search
	addToParams(p, "search", inputParams.Search)
	// manager
	addSliceToParams(p, "manager", inputParams.Manager)
	// groups
	addSliceToParams(p, "groups", inputParams.Groups)
	// stock_creator_id
	addToParams(p, "stock_creator_id", inputParams.StockCreatorID)
	// fields
	fieldsCount := 0
	for k, v := range inputParams.Fields {
		if k == 0 || v == "" {
			continue
		}
		p[fmt.Sprintf("params[fields][%d][id]", fieldsCount)] = strconv.FormatUint(k, 10)
		p[fmt.Sprintf("params[fields][%d][value]", fieldsCount)] = v
		fieldsCount++
	}
	// index_fields
	switch v := inputParams.IndexFields; v {
	case "true", "false", "1", "0":
		addToParams(p, "index_fields", v)
	}
	// related_with_customer
	addToParams(p, "related_with_customer", inputParams.RelatedWithCustomer)
	// order
	switch v := inputParams.Order; v {
	case "asc", "desc":
		addToParams(p, "order", v)
	}
	// order_field
	switch v := inputParams.OrderField; v {
	case "stock_activity_date", "date_add", "date_delete", "ID":
		addToParams(p, "order_field", v)
	default:
		if _, err := strconv.ParseUint(v, 10, 64); err == nil {
			addToParams(p, "order_field", v)
		}
	}
	// date
	if !inputParams.Date[0].IsZero() {
		p["params[date][from]"] = inputParams.Date[0].Format(DatetimeLayout)
	}
	if !inputParams.Date[1].IsZero() {
		p["params[date][to]"] = inputParams.Date[1].Format(DatetimeLayout)
	}
	// date_field
	addToParams(p, "date_field", inputParams.DateField)
	// page
	addToParams(p, "page", inputParams.Page)
	// publish
	switch v := inputParams.Publish; v {
	case "true", "false", "1", "0", "ignore":
		addToParams(p, "publish", v)
	}
	// limit
	switch v := inputParams.Limit; {
	case v == 0, v >= 500:
		addToParams(p, "limit", "500")
	default:
		addToParams(p, "limit", v)
	}
	// only_primary_id
	switch v := inputParams.OnlyPrimaryID; v {
	case "true", "false", "1", "0":
		addToParams(p, "only_primary_id", v)
	}
	// slice_fields
	addSliceToParams(p, "slice_fields", inputParams.SliceFields)

	// Запрос
	resp := new(StockFilterResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
