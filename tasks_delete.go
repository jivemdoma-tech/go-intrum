package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-delete
func TasksDelete(ctx context.Context, subdomain, apiKey string, taskID int64) (*TasksDeleteResp, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/tasks/delete", subdomain)

	if taskID == 0 {
		return nil, returnErrBadParams(methodURL)
	}

	// Параметры запроса
	p := map[string]string{
		"id": strconv.FormatInt(taskID, 10),
	}

	// Запрос
	resp := new(TasksDeleteResp)
	if err := request(ctx, apiKey, methodURL, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
