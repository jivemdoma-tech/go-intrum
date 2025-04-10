package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// Ссылка на метод: https://www.intrumnet.com/api/#history
type HistoryLogParams struct {
	ObjectType string                   // Обязательно. Одно из значений: stock | customer | sale | request
	ObjectID   []uint64                 // Массив ID объекта
	EmployeeID []uint64                 // Массив ID сотрудников, проводивших изменения
	Date       [2]time.Time             // Выборка за определенный период
	Log        [][]*HistoryLogParamsLog // Массив условий
}
type HistoryLogParamsLog struct {
	PropertyID string       // ID свойства
	Date       [2]time.Time // Выборка за определенный период
	Value      string       // Предыдущее значение. Одно из значений: @any | @empty | @not-empty
	Current    string       // Текущее значение. Одно из значений: @any | @empty | @not-empty
}

// Ссылка на метод: https://www.intrumnet.com/api/#history
func HistoryLog(ctx context.Context, subdomain, apiKey string, inputParams *HistoryLogParams) (*HistoryLogResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/history/log", subdomain)

	// Обязательность параметров
	if inputParams.ObjectType == "" {
		return nil, fmt.Errorf("error create request for method history logs: object_type param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 8)

	// object_type
	params["params[object_type]"] = inputParams.ObjectType
	// object_id
	for i, id := range inputParams.ObjectID {
		params[fmt.Sprintf("params[object_id][%d]", i)] = strconv.FormatUint(id, 10)
	}
	// employee_id
	for i, id := range inputParams.EmployeeID {
		params[fmt.Sprintf("params[employee_id][%d]", i)] = strconv.FormatUint(id, 10)
	}
	// date
	if !inputParams.Date[0].IsZero() {
		params["params[date][from]"] = inputParams.Date[0].Format(dateLayout)
	}
	if !inputParams.Date[1].IsZero() {
		params["params[date][to]"] = inputParams.Date[1].Format(dateLayout)
	}
	// log
	for i, logParamsSlice := range inputParams.Log {
		for j, logParams := range logParamsSlice {
			// property_id
			if logParams.PropertyID != "" {
				params[fmt.Sprintf("params[log][%d][%d][property_id]", i, j)] = logParams.PropertyID
			}
			// date
			if !logParams.Date[0].IsZero() {
				params[fmt.Sprintf("params[log][%d][%d][date][from]", i, j)] = logParams.Date[0].Format(dateLayout)
			}
			if !logParams.Date[1].IsZero() {
				params[fmt.Sprintf("params[log][%d][%d][date][to]", i, j)] = logParams.Date[1].Format(dateLayout)
			}
			// value
			if logParams.Value != "" {
				params[fmt.Sprintf("params[log][%d][%d][value]", i, j)] = logParams.Value
			}
			// current
			if logParams.Current != "" {
				params[fmt.Sprintf("params[log][%d][%d][current]", i, j)] = logParams.Current
			}
		}
	}

	// Получение ответа

	resp := new(HistoryLogResponse)
	if err := rawRequest(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
