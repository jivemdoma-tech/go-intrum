package gointrum

// import (
// 	"context"
// 	"fmt"
// )

// // Ссылка на метод: https://www.intrumnet.com/api/#sales-details
// type SalesDetailsParams struct {
// 	IDs []uint64 // ID объектов
// }

// func SalesAttach(ctx context.Context, subdomain, apiKey string, params *SalesDetailsParams) (*SalesDetailsParams, error) {
// 	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/details", subdomain)

// 	// Параметры запроса

// 	p := make(map[string]string, len(params.IDs))
// 	// id
// 	addSliceToParams("id", p, params.IDs)

// 	// Получение ответа

// 	resp := new(SalesDetailsResponse)
// 	if err := rawRequest(ctx, apiKey, methodURL, p, resp); err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }
