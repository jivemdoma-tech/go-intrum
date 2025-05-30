package gointrum

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type TasksCreateParams struct {
	Title       string            // заголовок задачи //! Обязательно если нет описания
	Description string            // описание //! Обязательно если нет заголовка
	Director    uint64            // id постановщика //! Обязательное поле
	Performer   uint64            // id исполнителя //! Обязательное поле
	Coperformer []uint64          // Массив id соисполнителей
	Priority    uint8             // приоритет задачи (1 - низкий, 7 - срочная задача)
	Attaches    map[string]uint64 // прикрепленные сущности в виде сущность#id (stock | customer | sale | request) через запятую
}

// Ссылка на метод: https://www.intrumnet.com/api/#tasks-create
//
// Ограничение 1 запрос == 1 задача
func TasksCreate(ctx context.Context, subdomain, apiKey string, inputParams *TasksCreateParams) (*TasksCreateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/tasks/create", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.Title == "" && inputParams.Description == "":
		return nil, fmt.Errorf("failed to create request for method tasks create: title or description is required")
	case inputParams.Director == 0:
		return nil, fmt.Errorf("failed to create request for method tasks create: director id is required")
	case inputParams.Performer == 0:
		return nil, fmt.Errorf("failed to create request for method tasks create: performer id is required")
	}

	// Параметры

	params := make(map[string]string, 8)

	// title
	if inputParams.Title != "" {
		params["params[title]"] = inputParams.Title
	}
	// description
	if inputParams.Description != "" {
		params["params[description]"] = inputParams.Description
	}
	// director
	params["params[director]"] = strconv.FormatUint(inputParams.Director, 10)
	// performer
	params["params[performer]"] = strconv.FormatUint(inputParams.Performer, 10)
	// coperformer
	if len(inputParams.Coperformer) != 0 {
		// Преобразование в слайс строк
		sliceStr := make([]string, 0, len(inputParams.Coperformer))
		for _, v := range inputParams.Coperformer {
			sliceStr = append(sliceStr, strconv.FormatUint(v, 10))
		}

		params["params[performer]"] = strings.Join(sliceStr, ",") // Объединение слайса строк через ,
	}
	// priority
	if inputParams.Priority != 0 {
		params["params[priority]"] = strconv.FormatUint(uint64(inputParams.Priority), 10)
	}
	// attaches
	if len(inputParams.Attaches) != 0 {
		s := make([]string, 0, len(inputParams.Attaches))
		for k, v := range inputParams.Attaches {
			s = append(s, k+"#"+strconv.FormatUint(v, 10)) // Форматирует в строку формата 'entity#123456'
		}

		params["params[attaches]"] = strings.Join(s, ",") // Объединение слайса строк через ,
	}

	// Получение ответа

	resp := new(TasksCreateResponse)
	if err := rawRequest(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil

}
