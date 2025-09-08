package gointrum

import (
	"context"
	"fmt"
)

type TasksSearchParams struct {
	// Текст для поиска в заголовке задачи
	// 	! ОБЯЗАТЕЛЬНО !
	Title string
	Limit int64 // Лимит задач в выборке (По умолчанию - 100, макс. значение - 1000)
	Page  int64 // Номер страницы выдачи (По умолчанию - 0)
}

// Ссылка на метод: https://www.intrumnet.com/api/#tasks-search
func TasksSearch(ctx context.Context, subdomain, apiKey string, inParams TasksSearchParams) (*TasksSearchResp, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/tasks/search", subdomain)

	// Обязательность ввода параметров
	if inParams.Title == "" {
		return nil, returnErrBadParams(methodURL)
	}

	// Параметры запроса
	p := make(map[string]string, 3)

	// title
	addToParams(p, "title", inParams.Title)
	// limit
	switch v := inParams.Limit; {
	case v == 0:
		addToParams(p, "limit", "100")
	case v >= 1000:
		addToParams(p, "limit", "1000")
	default:
		addToParams(p, "limit", v)
	}
	// page
	addToParams(p, "page", inParams.Page)

	// Запрос
	resp := new(TasksSearchResp)
	if err := requestFn(ctx, apiKey, subdomain, p, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
