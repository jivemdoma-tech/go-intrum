package gointrum

// Список возможных ошибок
//
//	"no text content" - не заданы заголовок и описание задачи (поля title, description)
//	"no director" - не задан постановщик задачи (поле director)
//	"no performer" - не задан исполнитель задачи (поле performer)
//	"invalid priority" - значение приоритета вне диапазона [1, 7]
//	"unknown entity type" - неизвестный тип CRM сущности в параметре attaches, допустимые значения: stock, customer, request, sale
//	"invalid entity id" - некорректное значение id прикрепленной CRM сущности в параметре attaches
type TasksCreateResponse struct {
	*Response
	Data TasksData `json:"data"`
}

type TasksData struct {
	ID int64 `json:"id"`
}
