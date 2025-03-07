package intrumgo

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-types
func SalesTypes(ctx context.Context, subdomain, apiKey string, timeoutSec int) (*SalesTypesResponse, error) {
	var (
		primaryURL string = fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/types", subdomain)
		backupURL  string = fmt.Sprintf("http://%s.intrumnet.com:80/sharedapi/sales/types", subdomain)
	)

	// Параметры запроса

	params := make(map[string]string, 0)

	// Получение ответа

	var resp SalesTypesResponse

	if err := rawRequest(ctx, primaryURL, apiKey, timeoutSec, params, &resp); err != nil {
		if err := rawRequest(ctx, backupURL, apiKey, timeoutSec, params, &resp); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}
