package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#stock-attach
type StockAttachParams struct {
	ID []uint64 // ID объектов
}

// StockAttach. Ссылка на метод: https://www.intrumnet.com/api/#stock-attach
func StockAttach(ctx context.Context, subdomain, apiKey string, params *StockAttachParams) (*StockAttachResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/attach", subdomain)

	// Параметры запроса

	p := make(map[string]string, len(params.ID))
	// id
	addSliceToParams("id", p, params.ID)

	// Получение ответа

	resp := new(StockAttachResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
