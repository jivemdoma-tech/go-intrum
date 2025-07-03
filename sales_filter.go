package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type SalesFilterParams struct {
	// Получение сделок по массиву ID
	ByIDs []int64

	// Поисковая строка
	Search string

	// Массив ID типов сделок
	Type []int64

	// Массив ID стадий сделок
	Stage []int64

	// Массив ID ответственных
	Manager []int64

	// Массив CRM групп
	Groups []int64

	// ID создателя
	SaleCreatorID int64

	// ID прикрепленного контакта
	Customer int64

	// Массив условий поиска.
	//	Key: ID поля
	//	Value: Значение поля
	// Для полей с типом integer, decimal, price, time, date, datetime возможно указывать границы:
	//	Value: ">= {значение}" - больше или равно
	//	Value: "<= {значение}" - меньше или равно
	//	Value: "{значение_1} & {значение_2}" - между значением 1 и 2
	Fields map[int64]string

	// Массив ID дополнительных полей, которые будут в ответе (по умолчанию выводятся все)
	SliceFields []int64

	// Направление сортировки (asc - по возрастанию, desc - по убыванию)
	Order string

	// ID поля, по которому нужно сделать сортировку. Если в качестве значения указать:
	// 	"sale_activity_date" - сортировка по дате активности
	// 	"create_date" - сортировка по дате создания
	// 	"delete_date" - сортировка по дате удаления
	OrderField string

	// Выборка за определенный период
	Date [2]time.Time

	// Если в качестве значения указать sale_activity_date, то выборка по параметру последней активности
	DateField string

	// (bool) "1" - активные | "0" - удаленные | "ignore" - вывод всех (по умолчанию "1")
	Publish string

	// Число записей в выборке (По умолчанию 500)
	Limit uint64

	// Номер страницы выборки (нумерация с 1)
	Page int64

	// TODO
	// CountTotal     string // (bool) Подсчет общего количества найденых записей
	// OnlyCountField string // (bool) Вывести в ответе только количество
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter
func SalesFilter(ctx context.Context, subdomain, apiKey string, inParams *SalesFilterParams) (*SalesFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/filter", subdomain)

	// Параметры запроса
	p := make(map[string]string, 8+
		len(inParams.ByIDs)+
		len(inParams.Type)+
		len(inParams.Stage)+
		len(inParams.Manager)+
		len(inParams.Groups)+
		len(inParams.SliceFields)+
		len(inParams.Fields)*2)

	// Параметры запроса

	// byid + by_ids
	switch {
	case len(inParams.ByIDs) == 1:
		addToParams(p, "byid", inParams.ByIDs[0])
	case len(inParams.ByIDs) >= 2:
		addSliceToParams(p, "by_ids", inParams.ByIDs)
	}
	// search
	addToParams(p, "search", inParams.Search)
	// type
	addSliceToParams(p, "type", inParams.Type)
	// stage
	addSliceToParams(p, "stage", inParams.Stage)
	// manager
	addSliceToParams(p, "manager", inParams.Manager)
	// groups
	addSliceToParams(p, "groups", inParams.Groups)
	// sale_creator_id
	if v := inParams.SaleCreatorID; v > 0 {
		addToParams(p, "sale_creator_id", v)
	}
	// customer
	if v := inParams.Customer; v > 0 {
		addToParams(p, "customer", v)
	}
	// fields
	fieldsCount := 0
	for k, v := range inParams.Fields {
		if k == 0 || v == "" {
			continue
		}
		p[fmt.Sprintf("params[fields][%d][id]", fieldsCount)] = strconv.FormatInt(k, 10)
		p[fmt.Sprintf("params[fields][%d][value]", fieldsCount)] = v
		fieldsCount++
	}
	// slice_fields
	addSliceToParams(p, "slice_fields", inParams.SliceFields)
	// order
	switch v := inParams.Order; v {
	case "asc", "desc":
		addToParams(p, "order", v)
	}
	// order_field
	switch v := inParams.OrderField; v {
	case "sale_activity_date", "create_date", "delete_date":
		addToParams(p, "order_field", v)
	default:
		if _, err := strconv.ParseUint(v, 10, 64); err == nil {
			addToParams(p, "order_field", v)
		}
	}
	// date
	if !inParams.Date[0].IsZero() {
		p["params[date][from]"] = inParams.Date[0].Format(DateLayout)
	}
	if !inParams.Date[1].IsZero() {
		p["params[date][to]"] = inParams.Date[1].Format(DateLayout)
	}
	// date_field
	addToParams(p, "date_field", inParams.DateField)
	// publish
	addBoolStringToParams(p, "publish", inParams.Publish)
	// limit
	switch v := inParams.Limit; {
	case v == 0, v >= 500:
		addToParams(p, "limit", "500")
	default:
		addToParams(p, "limit", v)
	}
	// page
	if v := inParams.Page; v >= 1 {
		addToParams(p, "page", v)
	}

	// Запрос
	resp := new(SalesFilterResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
