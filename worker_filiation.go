package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#worker-filiation
func WorkerFiliation(ctx context.Context, subdomain, apiKey string) (*WorkerFiliationResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/worker/filiation", subdomain)

	resp := new(WorkerFiliationResponse)
	if err := rawRequest(ctx, apiKey, methodURL, nil, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
