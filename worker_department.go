package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: https://www.intrumnet.com/api/#worker-department
func WorkerDepartment(ctx context.Context, subdomain, apiKey string) (*WorkerDepartmentResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/worker/department", subdomain)

	resp := new(WorkerDepartmentResponse)
	if err := rawRequest(ctx, apiKey, methodURL, nil, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
