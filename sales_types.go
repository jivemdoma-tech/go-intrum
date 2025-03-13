package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-types
func SalesTypes(ctx context.Context, subdomain, apiKey string) (*SalesTypesResponse, error) {
	var u string = fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/types", subdomain)

	// Параметры запроса

	params := make(map[string]string, 0)

	// Получение ответа

	var resp SalesTypesResponse
	if err := rawRequest(ctx, apiKey, u, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
