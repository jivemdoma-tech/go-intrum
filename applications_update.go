package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#applications-update
type ApplicationsUpdateParams struct {
	ID                   uint64   // ID существующей заявки // ! Обязательно
	EmployeeID           uint64   // ID ответственного
	AdditionalEmployeeID []uint64 // Массив ID доп ответственных
	CustomersID          uint64   // ID контакта
	RequestName          string   // Имя заявки
	// Дополнительные поля
	//
	// 	Ключ uint64 == ID поля
	// 	Значение any == Значение поля
	//		"знач1,знач2,знач3" (Для значений с типом "множественный выбор")
	Fields map[uint64]string

	// TODO: Добавить больше параметров запроса
}

// Ссылка на метод: https://www.intrumnet.com/api/#applications-update
//
// Ограничение 1 запрос == 1 заявка
func ApplicationsUpdate(ctx context.Context, subdomain, apiKey string, inputParams *ApplicationsUpdateParams) (*ApplicationsUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/applications/update", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.ID == 0:
		return nil, fmt.Errorf("failed to create request for method applications update: id param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 4+len(inputParams.Fields))

	// id
	params["params[0][id]"] = strconv.FormatUint(inputParams.ID, 10)
	// employee_id
	if inputParams.EmployeeID != 0 {
		params["params[0][employee_id]"] = strconv.FormatUint(inputParams.EmployeeID, 10)
	}
	// additional_employee_id
	for i, v := range inputParams.AdditionalEmployeeID {
		params[fmt.Sprintf("params[0][additional_employee_id][%d]", i)] = strconv.FormatUint(v, 10)
	}
	// customers_id
	if inputParams.CustomersID != 0 {
		params["params[0][customers_id]"] = strconv.FormatUint(inputParams.CustomersID, 10)
	}
	// request_name
	if inputParams.RequestName != "" {
		params["params[0][request_name]"] = inputParams.RequestName
	}
	// fields
	countFields := 0
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatUint(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}

	// Получение ответа

	resp := new(ApplicationsUpdateResponse)
	if err := rawRequest(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
