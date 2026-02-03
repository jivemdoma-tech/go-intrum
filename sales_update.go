package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-update
type SalesUpdateParams struct {
	ID            uint64 // ID существующего объекта // ! Обязательно
	SalesStatusID uint64 // ID стадии сделки

	// Дополнительные поля
	//
	// 	Ключ uint64 == ID поля
	// 	Значение any == Значение поля
	//		"знач1,знач2,знач3" (Для значений с типом "множественный выбор")
	Fields map[uint64]string // Дополнительные поля

	// TODO: Добавить больше параметров запроса
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-update
//
// Ограничение 1 запрос == 1 сделка
func SalesUpdate(ctx context.Context, subdomain, apiKey string, inputParams *SalesUpdateParams) (*SalesUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/update", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.ID == 0:
		return nil, fmt.Errorf("error create request for method sales update: id param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 8+
		len(inputParams.Fields)*2)

	// id
	params["params[0][id]"] = strconv.FormatUint(inputParams.ID, 10)
	// sales_status_id
	if inputParams.SalesStatusID != 0 {
		params["params[0][sales_status_id]"] = strconv.FormatUint(inputParams.SalesStatusID, 10)
	}
	// fields
	countFields := 0
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatUint(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}

	// Получение ответа

	var resp SalesUpdateResponse
	if err := request(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
