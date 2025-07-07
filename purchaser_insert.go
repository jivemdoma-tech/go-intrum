package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

type PurchaserInsertParams struct {
	Name                string   // Имя
	Surname             string   // Фамилия
	Secondname          string   // Отчество
	ManagerID           int64    // ID ответственного
	AdditionalManagerID []int64  // Массив ID дополнительных ответственных
	Marktype            int64    // Тип
	Phone               []string // номер телефона, без добавления комментариев и мессенджеров
	//email - массив email адресов
	//fields
	//TODO fields and email
}

//Ссылка на метод: https://www.intrumnet.com/api/#purchaser-insert

// Ограничение 1 запрос == 1 контакт
func PurchaserInsert(ctx context.Context, subdomain, apiKey string, inputParams *PurchaserInsertParams) (*PurchaserInsertResp, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/purchaser/insert", subdomain)

	//Параметры запроса

	params := make(map[string]string, 8)

	//name
	if inputParams.Name != "" {
		params["params[0][name]"] = inputParams.Name
	}

	//surname
	if inputParams.Surname != "" {
		params["params[0][surname]"] = inputParams.Surname
	}

	//secondname
	if inputParams.Secondname != "" {
		params["params[0][secondname]"] = inputParams.Secondname
	}

	//manager_id
	if inputParams.ManagerID > 0 {
		params["params[0][manager_id]"] = strconv.FormatInt(inputParams.ManagerID, 10)
	}

	//additional_manager_id
	for i, v := range inputParams.AdditionalManagerID {
		if v > 0 {
			params[fmt.Sprintf("params[0][additional_manager_id][%d]", i)] = strconv.FormatInt(v, 10)
		}
	}

	//marktype
	if v := inputParams.Marktype; v > 0 {
		params["params[0][marktype]"] = strconv.FormatInt(v, 10)
	}

	//phone
	for i, v := range inputParams.Phone {
		if v != "" {
			params[fmt.Sprintf("params[0][phone][%d]", i)] = v
		}
	}

	//Получение ответа

	resp := new(PurchaserInsertResp)
	if err := request(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
