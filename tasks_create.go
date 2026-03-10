package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TasksCreateParams struct {
	// Заголовок задачи
	//	! ОБЯЗАТЕЛЬНО ! (Если не указан 'Description')
	Title string
	// Описание задачи
	//	! ОБЯЗАТЕЛЬНО ! (Если не указан 'Title')
	Description string
	// Постановщик задачи
	//	! ОБЯЗАТЕЛЬНО !
	Director int64
	// Исполнитель задачи
	//	! ОБЯЗАТЕЛЬНО !
	Performer int64

	// Список ID соисполнителей
	Coperformer []int64

	Terms time.Time // Сроки задачи
	// Приоритет задачи
	//	7 - срочная задача
	//	...
	//	1 - низкий приоритет
	Priority int64

	// Список прикрепленных сущностей
	//	Key: Сущность ("stock" | "customer" | "sale" | "request")
	//	Value: ID
	Attaches map[string]int64

	// TODO
	// Checklist any // Массив пунктов чеклиста задачи
}

// Ссылка на метод: https://www.intrumnet.com/api/#tasks-create
//
//	! ВНИМАНИЕ ! Ограничение 1 запрос == 1 задача
func TasksCreate(ctx context.Context, subdomain, apiKey string, inParams TasksCreateParams) (*TasksCreateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/tasks/create", subdomain)

	// Обязательность ввода параметров
	switch {
	case inParams.Title == "" && inParams.Description == "":
		return nil, newErrEmptyRequiredParams(methodURL)
	case inParams.Director <= 0, inParams.Performer <= 0:
		return nil, newErrEmptyRequiredParams(methodURL)
	}

	// Параметры запроса
	params := make(map[string]string, 8)

	// title
	addToSingularParams(params, "title", inParams.Title)
	// description
	addToSingularParams(params, "description", inParams.Description)
	// director
	addToSingularParams(params, "director", inParams.Director)
	// performer
	addToSingularParams(params, "performer", inParams.Performer)
	// coperformer
	if slice := inParams.Coperformer; len(slice) != 0 {
		// Преобразование в слайс строк
		sliceStr := make([]string, 0, len(slice))
		for _, v := range slice {
			sliceStr = append(sliceStr, strconv.FormatInt(v, 10))
		}
		addToSingularParams(params, "coperformer", strings.Join(sliceStr, ","))
	}
	// terms
	addToSingularParams(params, "terms", inParams.Terms)
	// priority
	switch v := inParams.Priority; v {
	case 1, 2, 3, 4, 5, 6, 7:
		addToSingularParams(params, "priority", v)
	}
	// attaches
	if m := inParams.Attaches; len(m) != 0 {
		s := make([]string, 0, len(m))
		for k, v := range m {
			if k == "" || v <= 0 {
				continue
			}
			switch k {
			case "stock", "customer", "sale", "request":
				s = append(s, k+"#"+strconv.FormatInt(v, 10))
			}
		}
		addToSingularParams(params, "attaches", strings.Join(s, ","))
	}

	// Запрос
	resp := new(TasksCreateResponse)
	if err := request(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
