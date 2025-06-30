package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-delete
func SalesDelete(ctx context.Context, subdomain, apiKey string, saleIDs ...int64) (*SalesDeleteResp, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/delete", subdomain)

	// Обязательность ввода параметров
	if len(saleIDs) == 0 {
		return nil, returnErrBadParams(methodURL)
	}

	// Параметры запроса
	p := make(map[string]string, len(saleIDs))
	for i, id := range saleIDs {
		if id > 0 {
			k := fmt.Sprintf("params[%d]", i)
			p[k] = strconv.FormatInt(id, 10)
		}
	}

	// Запрос
	resp := new(SalesDeleteResp)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
