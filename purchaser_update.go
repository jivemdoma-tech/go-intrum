package gointrum

import (
	"context"
	"fmt"
	"strconv"
)

// Ссылка на метод: https://www.intrumnet.com/api/example.php#purchaser-update
type PurchaserUpdateParams struct {
	ID uint64 // ID контакта // ! Обязательно
	// Параметр работает интуитивно и очень интересно. Первый элемент массива - гл. ответственный, остальные - доп. ответственные.
	// 	Передача {0, n...} удаляет главного ответственного.
	// 	Передача {n, 0} удаляет доп. ответственных. Передайте {1, 0} чтобы пропустить гл. ответственного.
	// 	Передача {0, 0} удаляет всех ответственных.
	Authors []uint64
	Fields  map[uint64]string
	//Surname string // Фамилия
	//Name    string // Имя

	// TODO: Добавить больше параметров запроса
}

// Ссылка на метод: https://www.intrumnet.com/api/example.php#purchaser-update
//
// Ограничение 1 запрос == 1 заявка
func PurchaserUpdate(ctx context.Context, subdomain, apiKey string, inputParams *PurchaserUpdateParams) (*PurchaserUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/purchaser/update", subdomain)

	// Обязательность параметров
	switch {
	case inputParams.ID == 0:
		return nil, fmt.Errorf("failed to create request for method purchaser update: id param is required")
	}

	// Параметры запроса

	params := make(map[string]string, 4+len(inputParams.Fields))

	// id
	params["params[0][id]"] = strconv.FormatUint(inputParams.ID, 10)

	// author + additional_autor
	if len(inputParams.Authors) != 0 {
		// Гл. ответственный
		switch primary := inputParams.Authors[0]; primary {
		// Удаление
		case 0:
			params["params[0][author]"] = ""
		// Пропуск
		case 1:
			break
		// Изменение
		default:
			params["params[0][author]"] = strconv.FormatUint(inputParams.Authors[0], 10)
		}

		// Доп. ответственные
		if len(inputParams.Authors) >= 2 {
			switch additional := inputParams.Authors[1:]; {
			// Удаление
			case len(additional) == 1 && additional[0] <= 0:
				params["params[0][additional_author]"] = ""
			// Изменение
			default:
				for i, v := range additional {
					if v > 0 {
						k := fmt.Sprintf("params[0][additional_author][%d]", i)
						params[k] = strconv.FormatUint(v, 10)
					}
				}
				addSliceToParams(params, "additional_author", additional)
			}
		}
	}

	// fields
	countFields := 0
	for k, v := range inputParams.Fields {
		params[fmt.Sprintf("params[0][fields][%d][id]", countFields)] = strconv.FormatUint(k, 10)
		params[fmt.Sprintf("params[0][fields][%d][value]", countFields)] = v
		countFields++
	}

	// Получение ответа

	resp := new(PurchaserUpdateResponse)
	if err := request(ctx, apiKey, methodURL, params, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
