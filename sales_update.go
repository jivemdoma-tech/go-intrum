package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-update
type SalesUpdateParams struct {
	Fields map[uint64]any
	// TODO: Добавить больше параметров
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-update
func SalesUpdate(ctx context.Context, subdomain, apiKey string, timeoutSec int, inputParams map[uint64]*SalesUpdateParams) (*SalesUpdateResponse, error) {
	var (
		primaryURL string = fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/update", subdomain)
		backupURL  string = fmt.Sprintf("http://%s.intrumnet.com:80/sharedapi/sales/update", subdomain)
	)

	// Параметры запроса

	params := make(map[string]string, getParamsSize(inputParams))

	var count1 int
	for saleID, saleParams := range inputParams {
		// TODO: Унифицировать добавление параметров (напр. id) внешней функцией
		params[fmt.Sprintf("params[%d][id]", count1)] = strconv.FormatUint(saleID, 10)

		var count2 int
		for k, v := range saleParams.Fields {
			params[fmt.Sprintf("params[%d][fields][%d][id]", count1, count2)] = fmt.Sprint(k)
			params[fmt.Sprintf("params[%d][fields][%d][value]", count1, count2)] = fmt.Sprint(v)
			count2++
		}

		count1++
	}

	// Получение ответа

	var resp SalesUpdateResponse

	if err := rawRequest(ctx, primaryURL, apiKey, timeoutSec, params, &resp); err != nil {
		if err := rawRequest(ctx, backupURL, apiKey, timeoutSec, params, &resp); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}
