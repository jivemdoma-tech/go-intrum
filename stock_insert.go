package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: 	http://domainname.intrumnet.com:81/sharedapi/stock/insert
type StockInsertParams struct {
	Parent              int64   // ID категории объекта // Обязательно
	Name                string  // Название объекта
	Author              int64   // ID ответственного
	AdditionalAuthor    []int64 // Массив ID дополнительных ответственных
	RelatedWithCustomer int64   // ID контакта, прикрепленного к объекту
	GroupID             int64   // ID группы
	Copy                int64   // Родительский объект группы
	// Дополнительные поля
	//
	// 	Ключ int64 == ID поля
	// 	Значение any == Значение поля
	//		"знач1,знач2,знач3" (Для значений с типом "множественный выбор")
	Fields       map[int64]string
	FieldsCoords map[int64]CoordsVal // Поле с координатами (относится к fields)
	FieldsFiles  map[int64][]string  // Файлы, в массиве указывать название файла на сервере интрум (относится к fileds)
}

type CoordsVal struct {
	Lat float64 // Широта
	Lon float64 // Долгота
}

// Ссылка на метод: 	http://domainname.intrumnet.com:81/sharedapi/stock/insert
//
// Ограничение 1 запрос == 1 объект
func StockInsert(ctx context.Context, subdomain, apiKey string, inputParams *StockInsertParams) (*StockInsertResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/insert", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.Parent == 0:
		return nil, fmt.Errorf("error create request for method stock insert: parent param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 8+
		len(inputParams.Fields)*2)

	// parent
	params["params[0][parent]"] = strconv.FormatInt(inputParams.Parent, 10)
	// name
	if inputParams.Name != "" {
		params["params[0][name]"] = inputParams.Name
	}
	// author
	if inputParams.Author != 0 {
		params["params[0][author]"] = strconv.FormatInt(inputParams.Author, 10)
	}
	// additional_author
	for i, v := range inputParams.AdditionalAuthor {
		params[fmt.Sprintf("params[0][additional_author][%d]", i)] = strconv.FormatInt(v, 10)
	}
	// related_with_customer
	if inputParams.RelatedWithCustomer != 0 {
		params["params[0][related_with_customer]"] = strconv.FormatInt(inputParams.RelatedWithCustomer, 10)
	}
	// group_id
	if inputParams.GroupID != 0 {
		params["params[0][group_id]"] = strconv.FormatInt(inputParams.GroupID, 10)
	}
	// copy
	if inputParams.Copy != 0 {
		params["params[0][copy]"] = strconv.FormatInt(inputParams.Copy, 10)
	}

	countFields := 0
	// fields
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}
	// fieldsCoords
	for k, v := range inputParams.FieldsCoords {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value][lat]", countFields)] = strconv.FormatFloat(v.Lat, 'f', 10, 64)
		params[fmt.Sprintf("params[0][fields][%d][value][lon]", countFields)] = strconv.FormatFloat(v.Lon, 'f', 10, 64)
		countFields++
	}
	// fieldsFiles
	for k, fileNames := range inputParams.FieldsFiles {
		for _, fileName := range fileNames {
			params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatInt(k, 10)
			params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = fileName
			countFields++
		}
	}

	// Получение ответа

	var resp StockInsertResponse
	if err := request(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
