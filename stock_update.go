package intrum

import (
	"context"
	"fmt"
	"strconv"
)

// TODO: Реализовать оставшиеся поля StockUpdateParams.
//  Списка полей нет, т.к. нет полноценной документации по методу, только примеры.
// TODO: Реализовать в StockUpdateParams.Fields изменение полей типов: file, attach

// StockUpdate - редактирование объекта в CRM. Документация: https://www.intrumnet.com/api/#stock-update
func StockUpdate(ctx context.Context, subdomain, apiKey string, p *StockUpdateParams) (*StockUpdateResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/stock/update", subdomain)

	// Валидация
	if err := validateRequestArgs(methodURL, subdomain, apiKey); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, newErrEmptyParams(methodURL)
	}

	// Обязательные поля
	if p.ID <= 0 {
		return nil, newErrEmptyRequiredParams(methodURL)
	}

	// Запрос
	resp := &StockUpdateResponse{}
	if err := request(ctx, apiKey, methodURL, p.params(), resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// =====================================================================================================================
// Request
// =====================================================================================================================

// StockUpdateParams - параметры запроса StockUpdate.
//
// Обязательные поля:
//   - ID
//
// Основные параметры запроса:
//   - ID:                  Id существующего объекта.
//   - Name:                Редактирование названия объекта. Передайте new("") для удаления.
//   - Manager:             Редактирование id главного ответственного. Передайте new(0) для удаления.
//   - AdditionalManagers:  Редактирование массива id доп. ответственных. Передайте []int64{} для удаления.
//   - RelatedWithCustomer: Редактирование id прикрепленного контакта. Передайте new(0) для удаления.
//   - Fields:              Редактирование массива id полей и значений. Для удаления поля передайте nil по ключу.
type StockUpdateParams struct {
	ID                  int64   // Id существующего объекта.
	Category            int64   // Редактирование id категории объекта.
	Name                *string // Редактирование названия объекта. Передайте new("") для удаления.
	Manager             *int64  // Редактирование id главного ответственного. Передайте new(0) для удаления.
	AdditionalManagers  []int64 // Редактирование массива id доп. ответственных. Передайте []int64{} для удаления.
	RelatedWithCustomer *int64  // Редактирование id прикрепленного контакта. Передайте new(0) для удаления.
	// Fields: Редактирование массива id полей и значений.
	//
	// Для типа (multiselect) возможно указывать несколько вариантов:
	//  "{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ},{ЗНАЧЕНИЕ}".
	//
	// Не работает для типов: point, file, attach.
	Fields      map[int64]*string
	FieldsPoint map[int64]*Point // Аналогично Fields для типа "point".
}

// params возвращает параметры запроса в формате map[string]string (с эффективным выделением памяти).
func (p StockUpdateParams) params() map[string]string {
	// Выделение памяти
	size := 5 // Поля с простыми типами
	size += len(p.AdditionalManagers)
	size += len(p.Fields) * 2
	size += len(p.FieldsPoint) * 3
	paramsMap := make(map[string]string, size)

	// id
	paramsMap["params[0][id]"] = strconv.FormatInt(p.ID, 10)
	// parent
	if v := p.Category; v > 0 {
		paramsMap["params[0][parent]"] = strconv.FormatInt(v, 10)
	}
	// name
	if pV := p.Name; pV != nil {
		paramsMap["params[0][name]"] = *pV
	}
	// author
	if pV := p.Manager; pV != nil {
		v := max(*pV, 0)
		paramsMap["params[0][author]"] = strconv.FormatInt(v, 10)
	}
	// additional_author
	switch vSlice := p.AdditionalManagers; {
	case vSlice == nil:
		break
	case len(vSlice) == 0:
		paramsMap["params[0][additional_author]"] = "false"
	default:
		for i, v := range vSlice {
			if v <= 0 {
				continue
			}
			k, v := fmt.Sprintf("params[0][additional_author][%d]", i), strconv.FormatInt(v, 10)
			paramsMap[k] = v
		}
	}
	// related_with_customer
	if pV := p.RelatedWithCustomer; pV != nil {
		v := max(*pV, 0)
		paramsMap["params[0][related_with_customer]"] = strconv.FormatInt(v, 10)
	}
	// fields
	fieldsCount := 0
	for id, pV := range p.Fields {
		if id <= 0 {
			continue
		}
		// ID
		paramsMap[fmt.Sprintf("params[0][fields][%d][id]", fieldsCount)] = strconv.FormatInt(id, 10)
		// Value
		paramsMap[fmt.Sprintf("params[0][fields][%d][value]", fieldsCount)] = func() string {
			if pV == nil {
				return ""
			}
			return *pV
		}()

		fieldsCount++
	}
	// fields (point)
	for id, point := range p.FieldsPoint {
		if id <= 0 {
			continue
		}
		// ID
		paramsMap[fmt.Sprintf("params[0][fields][%d][id]", fieldsCount)] = strconv.FormatInt(id, 10)
		// Value
		switch latStr, lonStr := point.StringLat(), point.StringLon(); {
		case latStr == "" || lonStr == "":
			paramsMap[fmt.Sprintf("params[0][fields][%d][value]", fieldsCount)] = ""
		default:
			paramsMap[fmt.Sprintf("params[0][fields][%d][value][lat]", fieldsCount)] = latStr
			paramsMap[fmt.Sprintf("params[0][fields][%d][value][lon]", fieldsCount)] = lonStr
		}

		fieldsCount++
	}

	return paramsMap
}

// =====================================================================================================================
// Response
// =====================================================================================================================

type StockUpdateResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    bool   `json:"data,omitempty"`
}

func (r *StockUpdateResponse) GetErrorMessage() string {
	switch {
	case r == nil:
		return ""
	default:
		return ""
	case r.Status != "" && r.Message != "":
		return r.Status + ": " + r.Message
	case r.Message != "":
		return r.Message
	}
}
