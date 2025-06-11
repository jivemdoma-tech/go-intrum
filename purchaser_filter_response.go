package gointrum

import (
	"encoding/json"
	"strconv"
	"time"
)

type PurchaserFilterResponse struct {
	*Response
	Data *PurchaserFilterData `json:"data,omitempty"`
}

type PurchaserFilterData struct {
	List  []Purchaser `json:"list"`
	Count bool        `json:"count"`
}

type Purchaser struct {
	ID                   uint64                 `json:"id,string"`
	GroupID              string                 `json:"group_id"`
	Name                 string                 `json:"name"`
	Surname              string                 `json:"surname"`
	Secondname           string                 `json:"secondname"`
	ManagerID            uint64                 `json:"manager_id,string"`
	Phone                []*Phone               `json:"phone"`
	Address              string                 `json:"address"`
	CreateDate           time.Time              `json:"create_date"`
	Comment              string                 `json:"comment"`
	CustomerActivityType string                 `json:"customer_activity_type"`
	CustomerActivityDate time.Time              `json:"customer_activity_date"`
	CustomerCreatorID    uint64                 `json:"customer_creator_id,string"`
	Fields               map[uint64]*StockField `json:"fields"`
	EmployeeID           uint64                 `json:"employee_id,string"`
	AdditionalManagerID  []uint64               `json:"additional_manager_id"`
	AdditionalEmployeeID []uint64               `json:"additional_employee_id"`
	// Markname             string                     `json:"markname"`
	// Marktype             string                     `json:"marktype"`
	// Nattype              string                     `json:"nattype"`
	// Email                []interface{}              `json:"email"`

}

type PurchaserField struct {
	Datatype string `json:"datatype"`
	Value    any    `json:"value"`
}

type Phone struct {
	Phone   string `json:"phone"`
	Comment string `json:"comment"`
}

func (s *Purchaser) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias Purchaser

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		CreateDate           string            `json:"create_date"`
		CustomerActivityDate string            `json:"customer_activity_date"`
		AdditionalManagerID  []string          `json:"additional_manager_id"`
		AdditionalEmployeeID []string          `json:"additional_employee_id"`
		Fields               []*PurchaserField `json:"fields"`
	}{
		Alias: (*Alias)(s), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена

	parsedDate, err := time.Parse(datetimeLayout, aux.CreateDate)
	if err != nil {
		return err
	}
	s.CreateDate = parsedDate

	parsedDate, err = time.Parse(datetimeLayout, aux.CustomerActivityDate)
	if err != nil {
		return err
	}
	s.CustomerActivityDate = parsedDate

	// Замена массивов

	newSlice := make([]uint64, 0, len(aux.AdditionalManagerID))
	for _, v := range aux.AdditionalManagerID {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	s.AdditionalManagerID = newSlice

	newSlice = make([]uint64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	s.AdditionalEmployeeID = newSlice

	return nil
}
