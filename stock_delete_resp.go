package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/example.php#stock-delete
func StockDelete(ctx context.Context, subdomain, apiKey string, stockIDs ...int64) (*StockDeleteResp, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/delete", subdomain)

	// Обязательность ввода параметров
	if len(stockIDs) == 0 {
		return nil, returnErrBadParams(methodURL)
	}

	// Параметры запроса
	p := make(map[string]string, len(stockIDs))
	for i, id := range stockIDs {
		if id > 0 {
			k := fmt.Sprintf("params[%d]", i)
			p[k] = strconv.FormatInt(id, 10)
		}
	}

	// Запрос
	resp := new(StockDeleteResp)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
