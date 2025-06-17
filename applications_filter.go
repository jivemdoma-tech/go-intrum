package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// Ссылка на метод: https://www.intrumnet.com/api/#applications-filter
type ApplicationFilterParams struct {
	Search           string           // поисковая строка (может содержать фамилию, телефон контакта или название заявки)
	Groups           []uint32         // массив CRM групп
	Manager          []uint64         // id ответственного или массив с несколькими id
	RequestCreatorID uint64           // id создателя
	ByID             uint64           // id заявки
	ByIDs            []uint64         // массив ids заявок
	Customer         uint64           // id контакта
	Fields           map[int64]string // массив условий поиска по полям
	Types            []uint32         // массив id типов
	OrderField       string           // если в качестве значения указать request_activity_date выборка будет сортироваться по дате активности
	Order            string           // направление сортировки asc - по возрастанию, desc - по убыванию (сортировка только по дате последней активности)
	Date             [2]time.Time     // {from: "2015-10-29", to: "2015-11-19"} выборка за определенный период
	DateField        string           // если в качестве значения указать request_activity_date выборка по параметру заявки, create_date - по дате создания, delete_date - по дате удаления, id - по id
	Statuses         []string         // массив id статусов
	Page             uint16           // номер страницы выборки (нумерация с 1)
	Publish          string           // 1 - активные, 0 - удаленные, по умолчанию 1
	Limit            uint32           // число записей в выборке (макс. 500)
	SliceFields      []uint64         // массив id дополнительных полей, которые будут в ответе (по умолчанию, если не задано, то выводятся все)
}

func ApplicationFilter(ctx context.Context, subdomain, apiKey string, params *ApplicationFilterParams) (*ApplicationFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/applications/filter", subdomain)

	// Параметры запроса

	p := make(map[string]string, 8)
	// search
	if params.Search != "" {
		p["params[search]"] = params.Search
	}
	// groups
	addSliceToParams("groups", p, params.Groups)
	// manager
	addSliceToParams("manager", p, params.Manager)
	// request_creator_id
	if params.RequestCreatorID != 0 {
		p["params[request_creator_id]"] = strconv.FormatUint(params.RequestCreatorID, 10)
	}
	// byid
	if params.ByID != 0 {
		p["params[byid]"] = strconv.FormatUint(params.ByID, 10)
	}
	// by_ids
	addSliceToParams("by_ids", p, params.ByIDs)
	// customer
	if params.Customer != 0 {
		p["params[customer]"] = strconv.FormatUint(params.Customer, 10)
	}
	// fields
	fieldCount := 0
	for k, v := range params.Fields {
		p[fmt.Sprintf("params[fields][%d][id]", fieldCount)] = fmt.Sprint(k)
		p[fmt.Sprintf("params[fields][%d][value]", fieldCount)] = v
		fieldCount++
	}
	// types
	addSliceToParams("types", p, params.Types)
	// order_field
	if params.OrderField != "" {
		p["params[order_field]"] = params.OrderField
	}
	// order
	if params.Order != "" {
		p["params[order]"] = params.Order
	}
	// date
	if !params.Date[0].IsZero() {
		p["params[date][from]"] = params.Date[0].Format(datetimeLayout)
	}
	if !params.Date[1].IsZero() {
		p["params[date][to]"] = params.Date[1].Format(datetimeLayout)
	}
	// date_field
	if params.DateField != "" {
		p["params[date_field]"] = params.DateField
	}
	// statuses
	addSliceToParams("statuses", p, params.Statuses)
	// page
	if params.Page != 0 {
		p["params[page]"] = strconv.FormatUint(uint64(params.Page), 10)
	}
	// publish
	switch params.Publish {
	case "1", "true":
		p["params[publish]"] = "1"
	case "0", "false":
		p["params[publish]"] = "0"
	case "ignore":
		p["params[publish]"] = "ignore"
	}
	// limit
	// limit
	switch l := params.Limit; {
	case l == 0, l >= 500:
		p["params[limit]"] = "500"
	default:
		p["params[limit]"] = strconv.FormatUint(uint64(l), 10)
	}
	// slice_fields
	addSliceToParams("slice_fields", p, params.SliceFields)
	// Получение ответа

	resp := new(ApplicationFilterResponse)
	if err := rawRequest(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
