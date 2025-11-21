package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#purchaser-attach
type PurchaserAttachParams struct {
	IDs []uint64 // ID объектов
}

// PurchaserAttach. Ссылка на метод: https://www.intrumnet.com/api/#purchaser-attach
func PurchaserAttach(ctx context.Context, subdomain, apiKey string, params *PurchaserAttachParams) (*PurchaserAttachResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/purchaser/attach", subdomain)

	// Параметры запроса

	p := make(map[string]string, len(params.IDs))
	// id
	addSliceToParams(p, "ids", params.IDs)

	// Получение ответа

	resp := new(PurchaserAttachResponse)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
