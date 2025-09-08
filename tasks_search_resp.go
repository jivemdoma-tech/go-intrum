package gointrum

import (
	"encoding/json"
	"time"
)

type TasksSearchResp struct {
	*Response
	Data Data `json:"data"`
}

type TasksSearchData struct {
	Tasks []Task `json:"tasks"`
	Count int64  `json:"count"`
	Pages int64  `json:"pages"`
}

type Task struct {
	ID          int64            `json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Priority    int64            `json:"priority"`
	Author      int64            `json:"author"`
	Director    int64            `json:"director"`
	Performer   int64            `json:"performer"`
	Coperformer []int64          `json:"coperformer"`
	Attaches    map[int64]string `json:"attaches"`
	// Terms       time.Time `json:"terms"`
	// Checklist   []interface{} `json:"checklist"`
	// Files       []interface{} `json:"files"`
	// Tags        []interface{} `json:"tags"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type (
		Alias Task
	)

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		Attaches  []struct {
			ID   int64  `json:"id"`
			Type string `json:"type"`
		} `json:"attaches"`
	}{
		Alias: (*Alias)(t), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedDate, err := time.Parse(DatetimeLayout, aux.CreatedAt)
	switch err {
	case nil:
		t.CreatedAt = parsedDate
	default:
		t.CreatedAt = time.Time{}
	}

	t.Attaches = func() map[int64]string {
		if len(aux.Attaches) == 0 {
			return nil
		}

		result := make(map[int64]string, len(aux.Attaches))
		for _, attach := range aux.Attaches {
			result[attach.ID] = attach.Type
		}

		return result
	}()

	return nil
}
