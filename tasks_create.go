package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

type TasksCreateParams struct {
	Title       string // заголовок задачи //!Обязательно если нет описания
	Description string // описание //!Обязательно если нет заголовка
	Director    uint64 // id постановщика //!Обязательное поле
	Performer   uint64 // id исполнителя //!Обязательное поле
	Coperformer string // id соисполнителей через запятую без пробелов
	Priority    uint8  // приоритет задачи (1 - низкий, 7 - срочная задача)
	Attaches    string // прикрепленные сущности в виде сущность#id (stock - объекты, customer - клиенты, request - заявки, sale - сделки), соединенных через запятую
	// Список возможных ошибок
	// no text content - не заданы заголовок и описание задачи (поля title, description)
	// no director - не задан постановщик задачи (поле director)
	// no performer - не задан исполнитель задачи (поле performer)
	// invalid priority - значение приоритета вне диапазона [1, 7]
	// unknown entity type - неизвестный тип CRM сущности в параметре attaches, допустимые значения: stock, customer, request, sale
	// invalid entity id - некорректное значение id прикрепленной CRM сущности в параметре attaches
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
	if inputParams.Coperformer != "" {
		params["params[coperformer]"] = inputParams.Coperformer
	}
	// priority
	if inputParams.Priority != 0 {
		params["params[priority]"] = strconv.FormatUint(uint64(inputParams.Priority), 10)
	}
	// attaches
	if inputParams.Attaches != "" {
		params["params[attaches]"] = inputParams.Attaches
	}

	// получение ответа

	resp := new(TasksCreateResponse)
	if err := rawRequest(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil

}
