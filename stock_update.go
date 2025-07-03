package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

type StockUpdateParams struct {
	// ID существующего объекта
	//	! ОБЯЗАТЕЛЬНО !
	ID uint64

	Parent int64  // ID категории объекта
	Name   string // Наименование объекта

	// ID гл. ответственного
	//	Ввод -1 удаляет гл. ответственного
	Author int64
	// ID доп. ответственных
	//	Ввод []int64{} удаляет доп. ответственных
	AdditionalAuthor []int64
	// ID контакта, прикрепленного к объекту
	//	Ввод -1 открепляет контакт
	RelatedWithCustomer int64

	// Доп. поля
	//	Key: ID поля
	//	Value: Значение поля
	//		"{знач1},{знач2}..." - для полей типа 'multiselect'
	Fields map[int64]string

	// TODO: Добавить больше параметров запроса
	// Проблема конечно в том что нормальной документации нет
	// и приходится вычленять параметры из примеров...
}

// Ссылка на метод: https://www.intrumnet.com/api/#stock-update
//
//	! ВНИМАНИЕ ! Ограничение 1 запрос == 1 объект
func StockUpdate(ctx context.Context, subdomain, apiKey string, inParams StockUpdateParams) (*StockUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/update", subdomain)

	// Обязательность ввода параметров
	if inParams.ID == 0 {
		return nil, returnErrBadParams(methodURL)
	}

	// Параметры запроса
	p := make(map[string]string, 8+
		len(inParams.AdditionalAuthor)+
		len(inParams.Fields)*2)

	// id
	p["params[0][id]"] = strconv.FormatUint(inParams.ID, 10)
	// parent
	if v := inParams.Parent; v > 0 {
		p["params[0][parent]"] = strconv.FormatInt(v, 10)
	}
	// name
	switch v := inParams.Name; {
	case v == " ":
		p["params[0][name]"] = ""
	case v != "":
		p["params[0][name]"] = v
	}
	// author
	switch v := inParams.Author; {
	case v > 0:
		p["params[0][author]"] = strconv.FormatInt(v, 10)
	case v < 0:
		p["params[0][author]"] = "0"
	}
	// additional_author
	switch vSlice := inParams.AdditionalAuthor; {
	case vSlice == nil:
		break
	case len(vSlice) == 0:
		p["params[0][additional_author]"] = "false"
	default:
		for i, v := range vSlice {
			if v == 0 {
				continue
			}
			k := fmt.Sprintf("params[0][additional_author][%d]", i)
			p[k] = strconv.FormatInt(v, 10)
		}
	}
	// related_with_customer
	switch v := inParams.RelatedWithCustomer; {
	case v > 0:
		p["params[0][related_with_customer]"] = strconv.FormatInt(v, 10)
	case v < 0:
		p["params[0][related_with_customer]"] = "0"
	}
	// fields
	countFields := 0
	for k, v := range inParams.Fields {
		if k <= 0 || v == "" {
			continue
		}
		p[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
		switch v {
		case " ":
			p[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = ""
		default:
			p[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		}
		countFields++
	}

	// Запрос
	resp := new(StockUpdateResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
