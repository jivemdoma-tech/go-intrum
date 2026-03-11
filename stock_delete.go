package intrum

import (
	"context"
	"fmt"
	"strconv"
)

// StockDelete - удаление объектов в CRM. Документация: https://www.intrumnet.com/api/example.php#stock-delete
func StockDelete(ctx context.Context, subdomain, apiKey string, ids ...int64) (*StockDeleteResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/delete", subdomain)

	// Валидация
	if err := validateRequestArgs(methodURL, subdomain, apiKey); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, newErrEmptyParams(methodURL)
	}

	// Параметры запроса
	p := make(map[string]string, len(ids))
	for i, id := range ids {
		if id > 0 {
			k, v := fmt.Sprintf("params[%d]", i), strconv.FormatInt(id, 10)
			p[k] = v
		}
	}

	// Запрос
	resp := &StockDeleteResponse{}
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// =====================================================================================================================
// Response
// =====================================================================================================================

type StockDeleteResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    bool   `json:"data,omitempty"`
}

func (r *StockDeleteResponse) GetErrorMessage() string {
	switch {
	case r == nil:
		return ""
	default:
		return ""
	case r.Status != "" && r.Message != "":
		return r.Status + ": " + r.Message
	case r.Message != "":
		return r.Message
	}
}
