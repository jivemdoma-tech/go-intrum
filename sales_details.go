package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-details
type SalesDetailsParams struct {
	IDs []uint64 // ID объектов
}

func SalesDetails(ctx context.Context, subdomain, apiKey string, params *SalesDetailsParams) (*SalesDetailsResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/details", subdomain)

	// Параметры запроса

	p := make(map[string]string, len(params.IDs))

	// ids
	addSliceToParams(p, "ids", params.IDs)

	// Получение ответа

	resp := new(SalesDetailsResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
