package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type PurchaserFilterParams struct {
	Groups              []uint32     // массив id CRM групп
	Manager             uint64       // id ответственного
	AdditionlaManagerID []uint64     // массив ID дополнительных ответственных
	CustomerCreatorID   uint64       // id создателя
	ByID                []uint64     // id контакта или массив id контактов
	Search              string       // поисковая строка (может содержать фамилию или имя, email, телефон)
	Date                [2]time.Time // {from: "2015-10-29", to: "2015-11-19"} выборка за определенный период
	Page                uint16       // номер страницы выборки (нумерация с 1)
	Publish             string       // 1 - активные, 0 - удаленные, по умолчанию 1
	Limit               uint32       // число записей в выборке (макс. 500)
	SliceFields         []uint64     // массив id дополнительных полей, которые будут в ответе (по умолчанию если не задано то выводятся все)
	// массив условий поиска по полям [{id:id свойства,value: значение},{...}] для полей с типом integer,decimal,price,time,date,datetime возможно указывать границы:
	// 	value: '>= значение' - больше или равно
	//	value: '<= значение' - меньше или равно
	// 	value: 'значение_1 & значение_2' - между значением 1 и 2
	Fields map[uint64]string
	// TODO
	// Marktype // массив id типов
	// NatType // одно из значений подтипа physface - Юрлицо, jurface - Физлицо, по умолчанию выводятся все
	// IndexFields // индексировать массив полей по id свойства, 1 - да, 0 - нет, (по умолчанию 0)
	// Order // направление сортировки asc - по возрастанию, desc - по убыванию
	// OrderField //
	// DateField string       //  если в качестве значения указать customer_activity_date выборка будет сортироваться по дате активности;
	// create_date - по дате создания, delete_date - по дате удаления, id - по id
	// CountTotal uint64 // подсчет общего количества найденых записей, 1 - считать, 0 - нет (по умолчанию 0)
	// OnlyCountField // 1 - вывести в ответе только количество, 0 - стандартный вывод (по умолчанию 0)
}

func PurchaserFilter(ctx context.Context, subdomain, apiKey string, params *PurchaserFilterParams) (*PurchaserFilterResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/purchaser/filter", subdomain)

	// Параметры запроса

	p := make(map[string]string, 8)

	// groups
	addSliceToParams("groups", p, params.Groups)
	// manager
	if params.Manager != 0 {
		p["params[manager]"] = strconv.FormatUint(params.Manager, 10)
	}
	// additional_manager_id
	addSliceToParams("additional_manager_id", p, params.AdditionlaManagerID)
	// customer_creator_id
	if params.CustomerCreatorID != 0 {
		p["params[customer_creator_id ]"] = strconv.FormatUint(params.CustomerCreatorID, 10)
	}
	// byid
	addSliceToParams("byid", p, params.ByID)
	// search
	if params.Search != "" {
		p["params[search]"] = params.Search
	}
	// date
	if !params.Date[0].IsZero() {
		p["params[date][from]"] = params.Date[0].Format(datetimeLayout)
	}
	if !params.Date[1].IsZero() {
		p["params[date][to]"] = params.Date[1].Format(datetimeLayout)
	}
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
	switch l := params.Limit; l {
	case 0:
		p["params[limit]"] = "500"
	default:
		p["params[limit]"] = strconv.FormatUint(uint64(params.Limit), 10)
	}
	// slice_fields
	addSliceToParams("slice_fields", p, params.SliceFields)
	// fields
	fieldCount := 0
	for k, v := range params.Fields {
		p[fmt.Sprintf("params[fields][%d][id]", fieldCount)] = fmt.Sprint(k)
		p[fmt.Sprintf("params[fields][%d][value]", fieldCount)] = v
		fieldCount++
	}

	// Получение ответа

	resp := new(PurchaserFilterResponse)
	if err := rawRequest(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
