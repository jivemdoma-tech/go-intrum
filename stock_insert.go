package intrumgo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Ссылка на метод: 	http://domainname.intrumnet.com:81/sharedapi/stock/insert
type StockInsertParams struct {
	Parent              uint64   // ID категории объекта // Обязательно
	Name                string   // Название объекта
	Author              uint64   // ID ответственного
	AdditionalAuthor    []uint64 // Массив ID дополнительных ответственных
	RelatedWithCustomer uint64   // ID контакта, прикрепленного к объекту
	GroupID             uint16   // ID группы
	Copy                uint64   // Родительский объект группы

	// Дополнительные поля
	//
	// 	Ключ uint64 == ID поля
	// 	Значение any == Значение поля
	//		[]any (Для значений с типом "множественных выбор")
	Fields map[uint64]any
}

// Ссылка на метод: 	http://domainname.intrumnet.com:81/sharedapi/stock/insert
func StockInsert(ctx context.Context, subdomain, apiKey string, inputParams []*StockInsertParams) (*StockInsertResponse, error) {
	var (
		primaryURL string = fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/insert", subdomain)
		backupURL  string = fmt.Sprintf("http://%s.intrumnet.com:80/sharedapi/stock/insert", subdomain)
	)

	// Параметры запроса

	params := make(map[string]string, getParamsSize(inputParams))

	for objectCount, objectParams := range inputParams {
		if objectParams.Parent == 0 {
			log.Println("error create request for method stock insert: parent param in required")
			continue
		}

		// TODO: унифицировать добавление параметров внешней функцией

		// parent
		params[fmt.Sprintf("params[%d][parent]", objectCount)] = strconv.FormatUint(objectParams.Parent, 10)

		// name
		params[fmt.Sprintf("params[%d][name]", objectCount)] = objectParams.Name

		// author
		params[fmt.Sprintf("params[%d][author]", objectCount)] = strconv.FormatUint(objectParams.Author, 10)

		// additional_author
		addAuthorStr := make([]string, 0, len(objectParams.AdditionalAuthor))
		for _, id := range objectParams.AdditionalAuthor {
			addAuthorStr = append(addAuthorStr, strconv.FormatUint(id, 10))
		}
		params[fmt.Sprintf("params[%d][additional_author]", objectCount)] = strings.Join(addAuthorStr, ",")

		// related_with_customer
		params[fmt.Sprintf("params[%d][related_with_customer]", objectCount)] = strconv.FormatUint(objectParams.RelatedWithCustomer, 10)

		// group_id
		params[fmt.Sprintf("params[%d][group_id]", objectCount)] = strconv.FormatUint(uint64(objectParams.GroupID), 10)

		// copy
		params[fmt.Sprintf("params[%d][copy]", objectCount)] = strconv.FormatUint(objectParams.Copy, 10)

		// Fields
		fieldCount := 0
		for k, v := range objectParams.Fields {
			params[fmt.Sprintf("params[%d][fields][%d][id]", objectCount, fieldCount)] = fmt.Sprint(k)
			params[fmt.Sprintf("params[%d][fields][%d][value]", objectCount, fieldCount)] = fmt.Sprint(v)
			fieldCount++
		}

		objectCount++
	}

	// Получение ответа

	var resp StockInsertResponse

	if err := rawRequest(ctx, primaryURL, apiKey, params, &resp); err != nil {
		if err := rawRequest(ctx, backupURL, apiKey, params, &resp); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}
