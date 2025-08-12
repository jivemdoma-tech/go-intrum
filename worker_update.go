package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

type WorkerUpdateParams struct {
	// ID сотрудника
	//	! ОБЯЗАТЕЛЬНО !
	ID uint64

	// Доп. поля
	//	Key: ID поля
	//	Value: Значение поля
	//		"{знач1},{знач2}..." - для полей типа 'multiselect'
	Fields map[int64]string

	// TODO: Добавить больше параметров запроса
}

// Ссылка на метод: https://www.intrumnet.com/api/#worker-update
//
//	! ВНИМАНИЕ ! Ограничение 1 запрос == 1 объект
func WorkerUpdate(ctx context.Context, subdomain, apiKey string, inParams WorkerUpdateParams) (*WorkerUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/worker/update", subdomain)

	// Обязательность ввода параметров
	if inParams.ID == 0 {
		return nil, returnErrBadParams(methodURL)
	}

	// Параметры запроса
	p := make(map[string]string, 8+
		len(inParams.Fields)*2)

	// id
	p["params[0][id]"] = strconv.FormatUint(inParams.ID, 10)
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
	resp := new(WorkerUpdateResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
