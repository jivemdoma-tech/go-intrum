package gointrum

import (
	"encoding/json"
	"time"
)

type SalesGetByChangeStageResponse struct {
	Status string                     `json:"status"`
	Data   *SalesGetByChangeStageData `json:"data"`
}
type SalesGetByChangeStageData struct {
	List []*SalesGetByChangeStageDataList `json:"list"`
}
type SalesGetByChangeStageDataList struct {
	SaleID     uint64    `json:"sale_id,string"`
	SaleTypeID uint16    `json:"sale_type_id,string"`
	ToStage    uint16    `json:"to_stage,string"`
	FromStage  uint16    `json:"from_stage,string"`
	Date       time.Time `json:"date"`
}

func (s *SalesGetByChangeStageDataList) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias SalesGetByChangeStageDataList

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(s), // Приведение типа к Alias
	}

	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена
	parsedTime, err := time.Parse(datetimeLayout, aux.Date)
	if err != nil {
		return err
	}
	s.Date = parsedTime

	return nil
}
