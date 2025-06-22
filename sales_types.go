package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-types
func SalesTypes(ctx context.Context, subdomain, apiKey string) (*SalesTypesResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/types", subdomain)

	// Параметры запроса

	var params map[string]string

	// Получение ответа

	var resp SalesTypesResponse
	if err := request(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
