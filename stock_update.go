package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/#stock-update
type StockUpdateParams struct {
	ID                  uint64   // ID существующего объекта // ! Обязательно
	Parent              uint16   // ID категории объекта
	Name                string   // Наименования объекта
	Author              uint64   // ID ответственного
	AdditionalAuthor    []uint64 // Массив ID доп. ответственных
	RelatedWithCustomer uint64   // ID контакта, прикрепленного к объекту
	GroupID             uint16   // Связь с группой объектов. Подробнее о группах: https://www.intrumnet.com/wiki/gruppirovka_produktov___obektov__zhilye_kompleksy__kottedzhnye_poselki_-207
	Copy                uint64   // Родительский объект группы
	// Дополнительные поля
	//
	// 	Ключ uint64 == ID поля
	// 	Значение any == Значение поля
	//		"знач1,знач2,знач3" (Для значений с типом "множественный выбор")
	Fields map[uint64]string

	// TODO: Добавить больше параметров запроса
}

// Ссылка на метод: https://www.intrumnet.com/api/#stock-update
//
// Ограничение 1 запрос == 1 сделка
func StockUpdate(ctx context.Context, subdomain, apiKey string, inputParams *StockUpdateParams) (*StockUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/update", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.ID == 0:
		return nil, fmt.Errorf("error create request for method stock update: id param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 8+
		len(inputParams.Fields)*2)

	// id
	params["params[0][id]"] = strconv.FormatUint(inputParams.ID, 10)
	// parent
	if inputParams.Parent != 0 {
		params["params[0][parent]"] = strconv.FormatUint(uint64(inputParams.Parent), 10)
	}
	// name
	if inputParams.Name != "" {
		params["params[0][name]"] = inputParams.Name
	}
	// author
	if inputParams.Author != 0 {
		params["params[0][author]"] = strconv.FormatUint(inputParams.Author, 10)
	}
	// additional_author
	for i, v := range inputParams.AdditionalAuthor {
		params[fmt.Sprintf("params[0][additional_author][%d]", i)] = strconv.FormatUint(v, 10)
	}
	// related_with_customer
	if inputParams.RelatedWithCustomer != 0 {
		params["params[0][related_with_customer]"] = strconv.FormatUint(inputParams.RelatedWithCustomer, 10)
	}
	// group_id
	if inputParams.GroupID != 0 {
		params["params[0][group_id]"] = strconv.FormatUint(uint64(inputParams.GroupID), 10)
	}
	// copy
	if inputParams.Copy != 0 {
		params["params[0][copy]"] = strconv.FormatUint(inputParams.Copy, 10)
	}
	// fields
	countFields := 0
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatUint(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}

	// Получение ответа

	var resp StockUpdateResponse
	if err := rawRequest(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
