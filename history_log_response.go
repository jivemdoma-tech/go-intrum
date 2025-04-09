package gointrum

import (
	"encoding/json"
	"time"
)

type HistoryLogResponse struct {
	Status string            `json:"status"`
	Data   []*HistoryLogData `json:"data"`
}
type HistoryLogData struct {
	ObjectID   uint64    `json:"object_id,string"`
	PropertyID string    `json:"property_id"`
	Value      string    `json:"value"`
	Current    string    `json:"current"`
	Date       time.Time `json:"date"`
	EmployeeID uint64    `json:"employee_id,string"`
}

func (d *HistoryLogData) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias HistoryLogData

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(d), // Приведение типа к Alias
	}

	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена
	parsedDate, err := time.Parse(datetimeLayout, aux.Date)
	if err != nil {
		return err
	}
	d.Date = parsedDate

	return nil
}
