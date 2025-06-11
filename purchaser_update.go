package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/example.php#purchaser-update
type PurchaserUpdateParams struct {
	ID uint64 // ID контакта // ! Обязательно
	//Surname string // Фамилия
	//Name    string // Имя
	Fields map[uint64]string

	// TODO: Добавить больше параметров запроса
}

// Ссылка на метод: https://www.intrumnet.com/api/example.php#purchaser-update
//
// Ограничение 1 запрос == 1 заявка
func PurchaserUpdate(ctx context.Context, subdomain, apiKey string, inputParams *PurchaserUpdateParams) (*PurchaserUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/purchaser/update", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.ID == 0:
		return nil, fmt.Errorf("failed to create request for method purchaser update: id param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 4+len(inputParams.Fields))

	// id
	params["params[0][id]"] = strconv.FormatUint(inputParams.ID, 10)
	// fields
	countFields := 0
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatUint(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}

	// Получение ответа

	resp := new(PurchaserUpdateResponse)
	if err := rawRequest(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
